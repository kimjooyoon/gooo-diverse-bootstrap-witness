// Package verifier interprets the metacode contract and compares the two
// independent wire-producing paths. Neither path imports this package.
package verifier

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/kimjooyoon/gooo-diverse-bootstrap-witness/internal/patha"
	"github.com/kimjooyoon/gooo-diverse-bootstrap-witness/internal/pathb"
	"github.com/kimjooyoon/gooo-diverse-bootstrap-witness/internal/wire"
)

const (
	statusClosed  = "CLOSED"
	statusUnknown = "UNKNOWN"
	statusRefuted = "REFUTED"
)

func LoadMeta(path string) (wire.Meta, []byte, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return wire.Meta{}, nil, err
	}
	var meta wire.Meta
	if err := json.Unmarshal(raw, &meta); err != nil {
		return wire.Meta{}, nil, fmt.Errorf("decode metacode: %w", err)
	}
	if err := ValidateMeta(meta); err != nil {
		return wire.Meta{}, nil, err
	}
	return meta, raw, nil
}

func ValidateMeta(meta wire.Meta) error {
	if meta.Schema != "gooo.diverse-bootstrap-witness/v1" || meta.Authority != "metacode" || meta.ContractID == "" {
		return errors.New("metacode schema, authority, or contract_id is invalid")
	}
	if meta.Language.SourceExtension != ".gooo" || meta.Semantic.IRSchema == "" {
		return errors.New("metacode language or semantic schema is invalid")
	}
	for _, rule := range []string{"program", "binding", "emission", "effect"} {
		if meta.Language.Keywords[rule] == "" {
			return fmt.Errorf("metacode keyword %q is missing", rule)
		}
	}
	if len(meta.GenerationGraph) != 5 {
		return errors.New("metacode must define exactly five generation graph nodes")
	}
	if len(meta.Independence.Paths) != 2 || len(meta.Independence.AllowedSharedPackages) != 1 || meta.Independence.AllowedSharedPackages[0] != "internal/wire" || !meta.Independence.RequireCIImportIntersectionCheck || !meta.Independence.RequireDistinctSources {
		return errors.New("metacode independence predicate is incomplete")
	}
	if !meta.CanonicalComparison.ArtifactCannotSubstitute || !meta.TerminalTracePolicy.ReasonDriftIsFailure {
		return errors.New("metacode identity separation or trace policy is incomplete")
	}
	if len(meta.Resolution.StatusPrecedence) != 3 || meta.Resolution.StatusPrecedence[0] != statusRefuted || meta.Resolution.StatusPrecedence[1] != statusUnknown || meta.Resolution.StatusPrecedence[2] != statusClosed {
		return errors.New("metacode status precedence is not REFUTED > UNKNOWN > CLOSED")
	}
	if !sameStrings(meta.Resolution.UnknownFields, []string{"stage", "step", "reason", "unknown_class", "next_operation", "blocked_by"}) {
		return errors.New("metacode unknown field contract is incomplete")
	}
	if len(meta.FixedCases) != 6 {
		return errors.New("metacode must define exactly six fixed cases")
	}
	ids := make(map[string]bool, len(meta.FixedCases))
	expected := map[string]string{
		"independent-convergence":            statusClosed,
		"comment-only-canonical-convergence": statusClosed,
		"injected-self-propagating-semantic": statusRefuted,
		"terminal-reason-drift":              statusRefuted,
		"missing-diverse-path":               statusUnknown,
		"byte-identical-replay":              statusClosed,
	}
	for _, fixed := range meta.FixedCases {
		if ids[fixed.ID] {
			return fmt.Errorf("duplicate fixed case %q", fixed.ID)
		}
		ids[fixed.ID] = true
		if expected[fixed.ID] != fixed.ExpectedStatus {
			return fmt.Errorf("fixed case %q has unexpected status %q", fixed.ID, fixed.ExpectedStatus)
		}
		if fixed.Source == "" || !strings.HasSuffix(fixed.Source, ".gooo") {
			return fmt.Errorf("fixed case %q has invalid source", fixed.ID)
		}
		if fixed.ExpectedStatus == statusUnknown {
			if fixed.Unknown == nil || !unknownComplete(fixed.Unknown) {
				return fmt.Errorf("fixed case %q must define all UNKNOWN fields", fixed.ID)
			}
		}
	}
	for id := range expected {
		if !ids[id] {
			return fmt.Errorf("missing fixed case %q", id)
		}
	}
	if len(meta.OptionalInputs) != 1 || meta.OptionalInputs[0].Name != "gooo-two-generation-bootstrap" || meta.OptionalInputs[0].Tag != "v0.1.1" || meta.OptionalInputs[0].RequiredGate != 0 || meta.OptionalInputs[0].Digest == "" {
		return errors.New("the prior bootstrap must remain one digest-pinned optional input with required_gate=0")
	}
	return nil
}

func RunConformance(metaPath, repoRoot, artifactDir, reportPath string) error {
	meta, rawMeta, err := LoadMeta(metaPath)
	if err != nil {
		return err
	}
	if err := requireOutside(repoRoot, artifactDir); err != nil {
		return fmt.Errorf("artifact output boundary: %w", err)
	}
	if err := os.MkdirAll(artifactDir, 0o755); err != nil {
		return err
	}
	report := wire.ConformanceReport{
		Schema:         "gooo.conformance/v1",
		ContractDigest: Digest(rawMeta),
		SuiteStatus:    statusClosed,
		FixedCaseCount: len(meta.FixedCases),
	}
	var mismatchRefuted bool
	var mismatchUnknown bool
	for _, fixed := range meta.FixedCases {
		observation, err := runCase(meta, repoRoot, artifactDir, fixed)
		if err != nil {
			return fmt.Errorf("case %s: %w", fixed.ID, err)
		}
		report.Cases = append(report.Cases, observation)
		if observation.ActualStatus != fixed.ExpectedStatus {
			mismatchRefuted = mismatchRefuted || observation.ActualStatus == statusRefuted || fixed.ExpectedStatus == statusRefuted
			mismatchUnknown = mismatchUnknown || observation.ActualStatus == statusUnknown || fixed.ExpectedStatus == statusUnknown
		}
	}
	report.SuiteStatus = resolve(mismatchRefuted, mismatchUnknown, !mismatchRefuted && !mismatchUnknown)
	if report.SuiteStatus != statusClosed {
		return fmt.Errorf("fixed-case expectations did not close: %s", report.SuiteStatus)
	}
	return WriteJSONOutside(reportPath, repoRoot, report)
}

func runCase(meta wire.Meta, repoRoot, artifactDir string, fixed wire.FixedCase) (wire.CaseObservation, error) {
	sourcePath := filepath.Join(repoRoot, filepath.FromSlash(fixed.Source))
	source, err := os.ReadFile(sourcePath)
	if err != nil {
		return wire.CaseObservation{}, err
	}
	caseDir := filepath.Join(artifactDir, fixed.ID)
	if err := os.MkdirAll(caseDir, 0o755); err != nil {
		return wire.CaseObservation{}, err
	}
	pathAResult, err := patha.Generate(meta, source, "normal")
	if err != nil {
		return wire.CaseObservation{}, fmt.Errorf("path-a: %w", err)
	}
	pathAObservation, err := savePath(caseDir, "path-a", pathAResult)
	if err != nil {
		return wire.CaseObservation{}, err
	}
	observation := wire.CaseObservation{
		ID:             fixed.ID,
		ExpectedStatus: fixed.ExpectedStatus,
		SourceDigest:   Digest(source),
		PathA:          pathAObservation,
		PathB:          wire.PathObservation{Available: fixed.PathBAvailable},
		Identity: wire.IdentityIndicators{
			SemanticIdentity:      false,
			ArtifactByteIdentity:  false,
			TerminalTraceIdentity: false,
		},
	}
	if fixed.ReplayPathA {
		replay, replayErr := patha.Generate(meta, source, "normal")
		if replayErr != nil {
			return wire.CaseObservation{}, fmt.Errorf("path-a replay: %w", replayErr)
		}
		replayIR, _ := canonical(replay.IR)
		firstIR, _ := canonical(pathAResult.IR)
		replayTrace, _ := canonical(replay.Trace)
		firstTrace, _ := canonical(pathAResult.Trace)
		replayIdentity := wire.IdentityIndicators{
			SemanticIdentity:      bytes.Equal(firstIR, replayIR),
			ArtifactByteIdentity:  bytes.Equal(pathAResult.ArtifactBytes, replay.ArtifactBytes),
			TerminalTraceIdentity: bytes.Equal(firstTrace, replayTrace),
		}
		observation.ReplayIdentity = &replayIdentity
	}
	if !fixed.PathBAvailable {
		observation.ActualStatus = statusUnknown
		observation.Reason = "independent path-b is unavailable"
		observation.Unknown = fixed.Unknown
		return observation, nil
	}
	pathBResult, err := pathb.Generate(meta, source, fixed.PathBVariant)
	if err != nil {
		return wire.CaseObservation{}, fmt.Errorf("path-b: %w", err)
	}
	pathBObservation, err := savePath(caseDir, "path-b", pathBResult)
	if err != nil {
		return wire.CaseObservation{}, err
	}
	observation.PathB = pathBObservation
	aIR, err := canonical(pathAResult.IR)
	if err != nil {
		return wire.CaseObservation{}, err
	}
	bIR, err := canonical(pathBResult.IR)
	if err != nil {
		return wire.CaseObservation{}, err
	}
	aTrace, err := canonical(pathAResult.Trace)
	if err != nil {
		return wire.CaseObservation{}, err
	}
	bTrace, err := canonical(pathBResult.Trace)
	if err != nil {
		return wire.CaseObservation{}, err
	}
	observation.Identity = wire.IdentityIndicators{
		SemanticIdentity:      bytes.Equal(aIR, bIR),
		ArtifactByteIdentity:  bytes.Equal(pathAResult.ArtifactBytes, pathBResult.ArtifactBytes),
		TerminalTraceIdentity: bytes.Equal(aTrace, bTrace),
	}
	if !observation.Identity.SemanticIdentity {
		observation.ActualStatus = statusRefuted
		observation.Reason = "canonical semantic IR differs"
		return observation, nil
	}
	if !observation.Identity.TerminalTraceIdentity {
		observation.ActualStatus = statusRefuted
		observation.Reason = "terminal reason/effect trace differs"
		return observation, nil
	}
	if !observation.Identity.ArtifactByteIdentity {
		observation.ActualStatus = statusRefuted
		observation.Reason = "generated artifact bytes differ"
		return observation, nil
	}
	if observation.ReplayIdentity != nil && (!observation.ReplayIdentity.SemanticIdentity || !observation.ReplayIdentity.ArtifactByteIdentity || !observation.ReplayIdentity.TerminalTraceIdentity) {
		observation.ActualStatus = statusRefuted
		observation.Reason = "byte-identical replay did not reproduce all identities"
		return observation, nil
	}
	observation.ActualStatus = statusClosed
	observation.Reason = "independent semantic, artifact, terminal-trace, and replay identities converged"
	return observation, nil
}

func savePath(caseDir, name string, result wire.GeneratedResult) (wire.PathObservation, error) {
	irBytes, err := canonical(result.IR)
	if err != nil {
		return wire.PathObservation{}, err
	}
	traceBytes, err := canonical(result.Trace)
	if err != nil {
		return wire.PathObservation{}, err
	}
	artifactName := name + ".generated.go"
	irName := name + ".ir.json"
	traceName := name + ".trace.json"
	if err := os.WriteFile(filepath.Join(caseDir, artifactName), result.ArtifactBytes, 0o644); err != nil {
		return wire.PathObservation{}, err
	}
	if err := os.WriteFile(filepath.Join(caseDir, irName), append(irBytes, '\n'), 0o644); err != nil {
		return wire.PathObservation{}, err
	}
	if err := os.WriteFile(filepath.Join(caseDir, traceName), append(traceBytes, '\n'), 0o644); err != nil {
		return wire.PathObservation{}, err
	}
	return wire.PathObservation{
		Available:      true,
		SemanticDigest: Digest(irBytes),
		ArtifactDigest: Digest(result.ArtifactBytes),
		TraceDigest:    Digest(traceBytes),
		ArtifactPath:   filepath.ToSlash(filepath.Join(filepath.Base(caseDir), artifactName)),
		IRPath:         filepath.ToSlash(filepath.Join(filepath.Base(caseDir), irName)),
		TracePath:      filepath.ToSlash(filepath.Join(filepath.Base(caseDir), traceName)),
	}, nil
}

func BuildEvidence(metaPath, repoRoot, conformancePath, artifactDir, runtimeDir, evidencePath, buildMetric, testMetric, conformanceMetric, generatedBuildMetric, generatedRunMetric, testJSON string) error {
	meta, rawMeta, err := LoadMeta(metaPath)
	if err != nil {
		return err
	}
	var conformance wire.ConformanceReport
	conformanceBytes, err := os.ReadFile(conformancePath)
	if err != nil {
		return err
	}
	if err := json.Unmarshal(conformanceBytes, &conformance); err != nil {
		return err
	}
	if conformance.SuiteStatus != statusClosed || len(conformance.Cases) != len(meta.FixedCases) {
		return errors.New("conformance report is not closed or has the wrong fixed-case count")
	}
	inventory, err := measureInventory(repoRoot, meta.MeasurementPolicy.InventoryExcludePaths)
	if err != nil {
		return err
	}
	generated, err := measureFiles(artifactDir)
	if err != nil {
		return err
	}
	runtime, err := wireRuntimeMetrics(buildMetric, testMetric, conformanceMetric, generatedBuildMetric, generatedRunMetric)
	if err != nil {
		return err
	}
	tests, err := parseTestJSON(testJSON)
	if err != nil {
		return err
	}
	runtimeExecution := wire.RuntimeExecutionMetrics{}
	for index := range conformance.Cases {
		fixed := &conformance.Cases[index]
		for _, pathName := range []string{"path-a", "path-b"} {
			available := pathName == "path-a" && fixed.PathA.Available || pathName == "path-b" && fixed.PathB.Available
			if !available {
				continue
			}
			runtimeExecution.AvailableArtifacts++
			stdoutPath := filepath.Join(runtimeDir, fixed.ID, pathName+".stdout")
			if _, err := os.Stat(stdoutPath); err != nil {
				runtimeExecution.UnknownArtifacts++
				return fmt.Errorf("missing generated artifact runtime output %s", stdoutPath)
			}
			runtimeExecution.ExecutedArtifacts++
		}
		if fixed.PathA.Available && fixed.PathB.Available {
			aOutput, readAErr := os.ReadFile(filepath.Join(runtimeDir, fixed.ID, "path-a.stdout"))
			bOutput, readBErr := os.ReadFile(filepath.Join(runtimeDir, fixed.ID, "path-b.stdout"))
			if readAErr != nil || readBErr != nil {
				return fmt.Errorf("runtime output cannot be read for %s", fixed.ID)
			}
			equal := bytes.Equal(aOutput, bOutput)
			fixed.Identity.RuntimeIdentity = &equal
			if fixed.ExpectedStatus == statusClosed && !equal {
				return fmt.Errorf("closed case %s has different generated runtime output", fixed.ID)
			}
		}
	}
	inside := 0
	if isInside(repoRoot, artifactDir) || isInside(repoRoot, runtimeDir) || isInside(repoRoot, evidencePath) {
		inside = 1
	}
	evidence := wire.Evidence{
		Schema:             "gooo.evidence/v1",
		ContractDigest:     Digest(rawMeta),
		SuiteStatus:        statusClosed,
		Conformance:        conformance,
		Inventory:          inventory,
		Runtime:            runtime,
		Tests:              tests,
		GeneratedArtifacts: generated,
		RuntimeExecution:   runtimeExecution,
		Authority: wire.AuthorityMetrics{
			RepositoryWrites:                0,
			InputRepositoryWrites:           0,
			GeneratedOutputInsideRepository: inside,
			AutoCommitAuthority:             0,
			AutoPushAuthority:               0,
			AutoMergeAuthority:              0,
		},
	}
	return WriteJSONOutside(evidencePath, repoRoot, evidence)
}

func wireRuntimeMetrics(buildMetric, testMetric, conformanceMetric, generatedBuildMetric, generatedRunMetric string) (wire.RuntimeMetrics, error) {
	build, err := parseMetric(buildMetric)
	if err != nil {
		return wire.RuntimeMetrics{}, err
	}
	test, err := parseMetric(testMetric)
	if err != nil {
		return wire.RuntimeMetrics{}, err
	}
	conformance, err := parseMetric(conformanceMetric)
	if err != nil {
		return wire.RuntimeMetrics{}, err
	}
	generatedBuild, err := parseMetric(generatedBuildMetric)
	if err != nil {
		return wire.RuntimeMetrics{}, err
	}
	generatedRun, err := parseMetric(generatedRunMetric)
	if err != nil {
		return wire.RuntimeMetrics{}, err
	}
	return wire.RuntimeMetrics{Build: build, Test: test, Conformance: conformance, GeneratedBuild: generatedBuild, GeneratedRun: generatedRun}, nil
}

func measureInventory(root string, exclusions []string) (wire.InventoryMetrics, error) {
	excluded := map[string]bool{}
	for _, item := range exclusions {
		excluded[filepath.Clean(item)] = true
	}
	var metrics wire.InventoryMetrics
	err := filepath.WalkDir(root, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() {
			if path != root && entry.Name() == ".git" {
				return filepath.SkipDir
			}
			if path != root {
				metrics.Directories++
			}
			return nil
		}
		if !entry.Type().IsRegular() {
			return nil
		}
		relative, err := filepath.Rel(root, path)
		if err != nil {
			return err
		}
		if excluded[filepath.Clean(relative)] {
			return nil
		}
		metrics.RegularFiles++
		suffix := filepath.Ext(path)
		if suffix != ".go" && suffix != ".gooo" {
			return nil
		}
		content, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		lines := physicalLines(content)
		if suffix == ".go" {
			metrics.GoFiles++
			metrics.GoPhysicalLines += lines
		} else {
			metrics.GoooFiles++
			metrics.GoooPhysicalLines += lines
		}
		return nil
	})
	return metrics, err
}

func measureFiles(root string) (wire.GeneratedMetrics, error) {
	var metrics wire.GeneratedMetrics
	err := filepath.WalkDir(root, func(path string, entry os.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if entry.IsDir() {
			return nil
		}
		if !entry.Type().IsRegular() {
			return nil
		}
		info, err := entry.Info()
		if err != nil {
			return err
		}
		metrics.Files++
		metrics.Bytes += info.Size()
		return nil
	})
	return metrics, err
}

func parseMetric(path string) (wire.RuntimeMetric, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return wire.RuntimeMetric{}, err
	}
	fields := strings.Fields(string(content))
	if len(fields) != 2 {
		return wire.RuntimeMetric{}, fmt.Errorf("metric %s must contain elapsed seconds and peak RSS KiB", path)
	}
	seconds, err := strconv.ParseFloat(fields[0], 64)
	if err != nil {
		return wire.RuntimeMetric{}, err
	}
	rss, err := strconv.ParseInt(fields[1], 10, 64)
	if err != nil {
		return wire.RuntimeMetric{}, err
	}
	return wire.RuntimeMetric{WallMS: int64(math.Round(seconds * 1000)), PeakRSSKiB: rss}, nil
}

func parseTestJSON(path string) (wire.TestMetrics, error) {
	file, err := os.Open(path)
	if err != nil {
		return wire.TestMetrics{}, err
	}
	defer file.Close()
	var metrics wire.TestMetrics
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var event struct {
			Action string `json:"Action"`
			Test   string `json:"Test"`
			Output string `json:"Output"`
		}
		if err := json.Unmarshal(scanner.Bytes(), &event); err != nil {
			metrics.Unknown++
			continue
		}
		switch event.Action {
		case "run":
			if event.Test != "" {
				metrics.Total++
				metrics.Selected++
				metrics.Executed++
			}
		case "output":
			if strings.Contains(event.Output, "(cached)") {
				metrics.Reused++
			}
		case "fail":
			if event.Test != "" {
				metrics.Failed++
			}
		case "skip":
			if event.Test != "" {
				metrics.Unknown++
			}
		}
	}
	return metrics, scanner.Err()
}

func WriteJSONOutside(path, repoRoot string, value any) error {
	if path == "" {
		return errors.New("output path is required")
	}
	if err := requireOutside(repoRoot, path); err != nil {
		return err
	}
	encoded, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	return os.WriteFile(path, append(encoded, '\n'), 0o644)
}

func canonical(value any) ([]byte, error) {
	return json.Marshal(value)
}

func Digest(data []byte) string {
	sum := sha256.Sum256(data)
	return "sha256:" + hex.EncodeToString(sum[:])
}

func physicalLines(data []byte) int {
	if len(data) == 0 {
		return 0
	}
	count := bytes.Count(data, []byte{'\n'})
	if data[len(data)-1] != '\n' {
		count++
	}
	return count
}

func unknownComplete(record *wire.UnknownRecord) bool {
	return record.Stage != "" && record.Step != "" && record.Reason != "" && record.UnknownClass != "" && record.NextOperation != "" && record.BlockedBy != ""
}

func resolve(refuted, unknown, closed bool) string {
	if refuted {
		return statusRefuted
	}
	if unknown {
		return statusUnknown
	}
	if closed {
		return statusClosed
	}
	return statusUnknown
}

func sameStrings(left, right []string) bool {
	if len(left) != len(right) {
		return false
	}
	for index := range left {
		if left[index] != right[index] {
			return false
		}
	}
	return true
}

func requireOutside(repoRoot, output string) error {
	if output == "" {
		return errors.New("caller-owned output path is required")
	}
	root, err := filepath.Abs(repoRoot)
	if err != nil {
		return err
	}
	target, err := filepath.Abs(output)
	if err != nil {
		return err
	}
	relative, err := filepath.Rel(root, target)
	if err != nil {
		return err
	}
	if relative == "." || (relative != ".." && !strings.HasPrefix(relative, ".."+string(filepath.Separator))) {
		return fmt.Errorf("%s is inside repository %s", target, root)
	}
	return nil
}

func isInside(repoRoot, target string) bool {
	root, err := filepath.Abs(repoRoot)
	if err != nil {
		return true
	}
	absolute, err := filepath.Abs(target)
	if err != nil {
		return true
	}
	relative, err := filepath.Rel(root, absolute)
	if err != nil {
		return true
	}
	return relative == "." || (relative != ".." && !strings.HasPrefix(relative, ".."+string(filepath.Separator)))
}
