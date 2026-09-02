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
	if meta.Schema != "gooo.diverse-bootstrap-witness/v2" || meta.Authority != "metacode" || meta.ContractID != "gooo-diverse-bootstrap-witness/v2" {
		return errors.New("metacode schema, authority, or contract_id is invalid")
	}
	if meta.ContractEvolution.PreviousSchema != "gooo.diverse-bootstrap-witness/v1" || meta.ContractEvolution.PreviousContractID != "gooo-diverse-bootstrap-witness/v1" || meta.ContractEvolution.PreviousDenominator != 6 || meta.ContractEvolution.CurrentDenominator != 14 || !meta.ContractEvolution.AppendOnly || len(meta.ContractEvolution.AddedFixedCaseIDs) != 8 || len(meta.ContractEvolution.RetiredFixedCaseIDs) != 0 {
		return errors.New("contract evolution must append eight cases to the released six-case denominator")
	}
	if meta.Language.SourceExtension != ".gooo" || meta.Semantic.IRSchema == "" {
		return errors.New("metacode language or semantic schema is invalid")
	}
	for _, rule := range []string{"program", "binding", "emission", "effect"} {
		if meta.Language.Keywords[rule] == "" {
			return fmt.Errorf("metacode keyword %q is missing", rule)
		}
	}
	if meta.Semantic.SourceAuthority != "this .gooo semantic kernel; Go is bounded evaluator/runtime only" || !sameStrings(meta.Semantic.ObservationFields, []string{"canonical_ir_digest", "decision_digest", "provenance_digest"}) {
		return errors.New("semantic source authority or observation fields are incomplete")
	}
	if meta.SemanticKernel.Source.Path == "" || !validDigest(meta.SemanticKernel.Source.Digest) || meta.SemanticKernel.Source.DigestScope != "exact-file-sha256" || meta.SemanticKernel.CanonicalInput.ID == "" || meta.SemanticKernel.CanonicalInput.Path == "" || !validDigest(meta.SemanticKernel.CanonicalInput.Digest) || meta.SemanticKernel.CanonicalInput.DigestScope != "exact-file-sha256" {
		return errors.New("semantic kernel source or canonical input identity is incomplete")
	}
	if meta.SemanticKernel.ExpectedObservation.Status != statusClosed || !meta.SemanticKernel.ExpectedObservation.SemanticIdentity || !meta.SemanticKernel.ExpectedObservation.DecisionIdentity || !meta.SemanticKernel.ExpectedObservation.ProvenanceIdentity || !sameStrings(meta.SemanticKernel.ExpectedObservation.RequiredDigests, []string{"canonical_ir_digest", "decision_digest", "provenance_digest"}) {
		return errors.New("expected canonical semantic observation is incomplete")
	}
	if meta.SemanticKernel.Evaluation.Authority != "semantic-kernel-source" || len(meta.SemanticKernel.Evaluation.Operations) < 5 || len(meta.SemanticKernel.Evaluation.Mutations) != 2 {
		return errors.New("semantic kernel evaluation declaration is incomplete")
	}
	if len(meta.GenerationGraph) != 5 {
		return errors.New("metacode must define exactly five generation graph nodes")
	}
	if len(meta.Independence.Paths) != 2 || len(meta.Independence.AllowedSharedPackages) != 1 || meta.Independence.AllowedSharedPackages[0] != "internal/wire" || !meta.Independence.RequireCIImportIntersectionCheck || !meta.Independence.RequireDistinctSources || !meta.Independence.RequireIndependentWitnesses || meta.Independence.SameLineageStatus != statusRefuted || meta.Independence.DuplicateArtifactStatus != statusRefuted || !sameStrings(meta.Independence.IdentityFields, []string{"id", "role", "implementation", "source_digest", "lineage_id", "toolchain_id", "input_identity"}) {
		return errors.New("metacode independence predicate is incomplete")
	}
	if meta.WitnessPlan.RequiredIndependentCount != 2 || !completeWitness(meta.WitnessPlan.Stage0Reference) || !completeWitness(meta.WitnessPlan.GeneratedCurrent) || meta.WitnessPlan.Stage0Reference.LineageID == meta.WitnessPlan.GeneratedCurrent.LineageID || meta.WitnessPlan.Stage0Reference.Implementation == meta.WitnessPlan.GeneratedCurrent.Implementation || len(meta.WitnessPlan.ProofClaims) != 3 || !meta.WitnessPlan.Trilemma.Acknowledged || meta.WitnessPlan.Trilemma.Foundation == "" || meta.WitnessPlan.Trilemma.Coherence == "" || meta.WitnessPlan.Trilemma.Regression == "" {
		return errors.New("independent witness plan or trilemma acknowledgement is incomplete")
	}
	for index, claim := range meta.WitnessPlan.ProofClaims {
		expectedChoice := []string{"FOUNDATION", "COHERENCE", "REGRESSION"}[index]
		if claim.ID == "" || claim.Claim == "" || claim.ProofChoice != expectedChoice || claim.IndependenceBasis == "" || len(claim.RequiredEvidence) == 0 {
			return fmt.Errorf("proof claim %d is incomplete", index+1)
		}
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
	if len(meta.Resolution.RefutedReasons) != 5 || !sameStrings(meta.MeasurementPolicy.WitnessRuntimeFields, []string{"wall_ms", "peak_rss_kib", "build_ms", "test_ms"}) || meta.MeasurementPolicy.ImprovementMissingStatus != statusUnknown || len(meta.MeasurementPolicy.ImprovementPairFields) != 6 || meta.MeasurementPolicy.ScoreAggregation != "forbidden" || meta.MeasurementPolicy.EstimatedRate != "forbidden" {
		return errors.New("refutation, witness measurement, or improvement policy is incomplete")
	}
	if len(meta.FixedCases) != meta.ContractEvolution.CurrentDenominator {
		return errors.New("metacode fixed-case denominator is invalid")
	}
	ids := make(map[string]bool, len(meta.FixedCases))
	expected := map[string]string{
		"independent-convergence":            statusClosed,
		"comment-only-canonical-convergence": statusClosed,
		"injected-self-propagating-semantic": statusRefuted,
		"terminal-reason-drift":              statusRefuted,
		"missing-diverse-path":               statusUnknown,
		"byte-identical-replay":              statusClosed,
		"stage0-current-convergence":         statusClosed,
		"missing-witness-identity":            statusUnknown,
		"missing-input-identity":              statusUnknown,
		"witness-disagreement":                statusRefuted,
		"self-approval-cycle":                 statusRefuted,
		"forged-digest":                       statusRefuted,
		"same-lineage-replay":                 statusRefuted,
		"frozen-bootstrap-mismatch":           statusRefuted,
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
		if fixed.ScenarioClass == "" || fixed.IdentityMode == "" {
			return fmt.Errorf("fixed case %q is missing scenario or identity mode", fixed.ID)
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
	if !sameStrings(meta.ContractEvolution.AddedFixedCaseIDs, []string{"stage0-current-convergence", "missing-witness-identity", "missing-input-identity", "witness-disagreement", "self-approval-cycle", "forged-digest", "same-lineage-replay", "frozen-bootstrap-mismatch"}) {
		return errors.New("contract evolution added case list is not append-only and ordered")
	}
	if len(meta.OptionalInputs) != 1 || meta.OptionalInputs[0].Name != "gooo-two-generation-bootstrap" || meta.OptionalInputs[0].Tag != "v0.1.1" || meta.OptionalInputs[0].RequiredGate != 0 || meta.OptionalInputs[0].Digest == "" {
		return errors.New("the prior bootstrap must remain one digest-pinned optional input with required_gate=0")
	}
	return nil
}

// Observe runs exactly one declared witness against one caller-selected input.
// It emits only a digest line; all generated material remains caller-owned.
func Observe(metaPath, repoRoot, sourcePath, witnessID string) error {
	meta, rawMeta, err := LoadMeta(metaPath)
	if err != nil {
		return err
	}
	if _, _, err := verifyDeclaredIdentities(meta, repoRoot); err != nil {
		return err
	}
	source, err := os.ReadFile(sourcePath)
	if err != nil {
		return err
	}
	var result wire.GeneratedResult
	switch witnessID {
	case meta.WitnessPlan.Stage0Reference.ID:
		result, err = patha.Generate(meta, source, "normal")
	case meta.WitnessPlan.GeneratedCurrent.ID:
		result, err = pathb.Generate(meta, source, "normal")
	default:
		return fmt.Errorf("unknown witness %q", witnessID)
	}
	if err != nil {
		return err
	}
	ir, err := canonical(result.IR)
	if err != nil {
		return err
	}
	contractDigest := Digest(rawMeta)
	decision := decisionDigest(result.Trace)
	provenance := provenanceDigest(Digest(source), contractDigest, Digest(ir), decision)
	observation := observationDigest(Digest(ir), decision, provenance)
	fmt.Printf("witness_id=%s input_digest=%s observation_digest=%s\n", witnessID, Digest(source), observation)
	return nil
}

func RunConformance(metaPath, repoRoot, artifactDir, reportPath string) error {
	meta, rawMeta, err := LoadMeta(metaPath)
	if err != nil {
		return err
	}
	kernelDigest, canonicalInputDigest, err := verifyDeclaredIdentities(meta, repoRoot)
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
		Schema:               "gooo.conformance/v2",
		ContractVersion:      meta.Schema,
		ContractDigest:       Digest(rawMeta),
		KernelSourceDigest:   kernelDigest,
		CanonicalInputDigest: canonicalInputDigest,
		SuiteStatus:          statusClosed,
		FixedCaseCount:       len(meta.FixedCases),
		Denominator:          denominatorFor(meta),
		Bootstrap:            bootstrapObservation(meta, kernelDigest, canonicalInputDigest),
	}
	var mismatchRefuted bool
	var mismatchUnknown bool
	for _, fixed := range meta.FixedCases {
		observation, err := runCase(meta, repoRoot, artifactDir, Digest(rawMeta), fixed)
		if err != nil {
			return fmt.Errorf("case %s: %w", fixed.ID, err)
		}
		observation.ProofReceipts = caseProofReceipts(meta, observation)
		report.Cases = append(report.Cases, observation)
		if observation.ActualStatus != fixed.ExpectedStatus {
			mismatchRefuted = mismatchRefuted || observation.ActualStatus == statusRefuted || fixed.ExpectedStatus == statusRefuted
			mismatchUnknown = mismatchUnknown || observation.ActualStatus == statusUnknown || fixed.ExpectedStatus == statusUnknown
		}
	}
	report.SuiteStatus = resolve(mismatchRefuted, mismatchUnknown, !mismatchRefuted && !mismatchUnknown)
	canonicalFound := false
	for _, observation := range report.Cases {
		if observation.ID != meta.SemanticKernel.CanonicalInput.ID {
			continue
		}
		canonicalFound = true
		expected := meta.SemanticKernel.ExpectedObservation
		if observation.ActualStatus != expected.Status || observation.Identity.SemanticIdentity != expected.SemanticIdentity || observation.Identity.DecisionIdentity != expected.DecisionIdentity || observation.Identity.ProvenanceIdentity != expected.ProvenanceIdentity {
			return fmt.Errorf("canonical input did not produce its declared semantic observation")
		}
	}
	if !canonicalFound {
		return errors.New("canonical input case is not present in the fixed denominator")
	}
	report.ProofReceipts = proofReceipts(meta, report.Cases)
	if report.SuiteStatus != statusClosed {
		return fmt.Errorf("fixed-case expectations did not close: %s", report.SuiteStatus)
	}
	return WriteJSONOutside(reportPath, repoRoot, report)
}

func runCase(meta wire.Meta, repoRoot, artifactDir, contractDigest string, fixed wire.FixedCase) (wire.CaseObservation, error) {
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
	pathAObservation, err := savePath(caseDir, "path-a", pathAResult, Digest(source), contractDigest)
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
			DecisionIdentity:      false,
			ProvenanceIdentity:   false,
			ArtifactByteIdentity:  false,
			TerminalTraceIdentity: false,
		},
	}
	observation.Witnesses = append(observation.Witnesses, witnessObservation(meta.WitnessPlan.Stage0Reference, pathAObservation, Digest(source), ""))
	if fixed.ReplayPathA {
		replay, replayErr := patha.Generate(meta, source, "normal")
		if replayErr != nil {
			return wire.CaseObservation{}, fmt.Errorf("path-a replay: %w", replayErr)
		}
		replayIR, _ := canonical(replay.IR)
		firstIR, _ := canonical(pathAResult.IR)
		replayTrace, _ := canonical(replay.Trace)
		firstTrace, _ := canonical(pathAResult.Trace)
		replayDecision := decisionDigest(replay.Trace)
		firstDecision := decisionDigest(pathAResult.Trace)
		replayProvenance := provenanceDigest(Digest(source), contractDigest, Digest(replayIR), replayDecision)
		firstProvenance := provenanceDigest(Digest(source), contractDigest, Digest(firstIR), firstDecision)
		replayIdentity := wire.IdentityIndicators{
			SemanticIdentity:      bytes.Equal(firstIR, replayIR),
			DecisionIdentity:      firstDecision == replayDecision,
			ProvenanceIdentity:   firstProvenance == replayProvenance,
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
	pathBObservation, err := savePath(caseDir, "path-b", pathBResult, Digest(source), contractDigest)
	if err != nil {
		return wire.CaseObservation{}, err
	}
	observation.PathB = pathBObservation
	observation.Witnesses = append(observation.Witnesses, witnessObservation(meta.WitnessPlan.GeneratedCurrent, pathBObservation, Digest(source), fixed.IdentityMode))
	if fixed.IdentityMode == "same-lineage" {
		observation.Witnesses[1].LineageID = observation.Witnesses[0].LineageID
	}
	if fixed.SelfApprovalCycle {
		observation.ActualStatus = statusRefuted
		observation.Reason = "self-approval cycle: a witness cannot approve the evaluator that approves itself"
		return observation, nil
	}
	if fixed.DigestClaim != "" && fixed.DigestClaim != pathBObservation.SemanticDigest {
		observation.ActualStatus = statusRefuted
		observation.Reason = "forged digest claim does not match the generated witness observation"
		return observation, nil
	}
	if fixed.FrozenBootstrapMismatch && fixed.FrozenBootstrapDigest != meta.WitnessPlan.Stage0Reference.SourceDigest {
		observation.ActualStatus = statusRefuted
		observation.Reason = "frozen bootstrap witness identity does not match the current semantic kernel"
		return observation, nil
	}
	if fixed.IdentityMode == "missing-lineage-toolchain-input" || fixed.IdentityMode == "missing-input" {
		observation.ActualStatus = statusUnknown
		observation.Reason = "witness lineage, toolchain, or input identity is incomplete"
		observation.Unknown = fixed.Unknown
		return observation, nil
	}
	if !independentWitnesses(observation.Witnesses[0], observation.Witnesses[1]) {
		observation.Witnesses[0].Independent = false
		observation.Witnesses[1].Independent = false
		observation.ActualStatus = statusRefuted
		observation.Reason = "witnesses are not independent: same source, lineage, or evaluator identity"
		return observation, nil
	}
	observation.Witnesses[0].Independent = true
	observation.Witnesses[1].Independent = true
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
		DecisionIdentity:      pathAObservation.DecisionDigest == pathBObservation.DecisionDigest,
		ProvenanceIdentity:   pathAObservation.ProvenanceDigest == pathBObservation.ProvenanceDigest,
		ArtifactByteIdentity:  bytes.Equal(pathAResult.ArtifactBytes, pathBResult.ArtifactBytes),
		TerminalTraceIdentity: bytes.Equal(aTrace, bTrace),
	}
	if !observation.Identity.SemanticIdentity {
		observation.ActualStatus = statusRefuted
		observation.Reason = "canonical semantic IR differs"
		return observation, nil
	}
	if !observation.Identity.DecisionIdentity {
		observation.ActualStatus = statusRefuted
		observation.Reason = "canonical decision digest differs"
		return observation, nil
	}
	if !observation.Identity.ProvenanceIdentity {
		observation.ActualStatus = statusRefuted
		observation.Reason = "canonical provenance digest differs"
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
	observation.Reason = "two independent witnesses agree on canonical IR, decision, provenance, artifact, and terminal trace"
	return observation, nil
}

func savePath(caseDir, name string, result wire.GeneratedResult, sourceDigest, contractDigest string) (wire.PathObservation, error) {
	irBytes, err := canonical(result.IR)
	if err != nil {
		return wire.PathObservation{}, err
	}
	decision := decisionDigest(result.Trace)
	provenance := provenanceDigest(sourceDigest, contractDigest, Digest(irBytes), decision)
	observation := observationDigest(Digest(irBytes), decision, provenance)
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
		Available:         true,
		SemanticDigest:    Digest(irBytes),
		DecisionDigest:    decision,
		ProvenanceDigest:  provenance,
		ObservationDigest: observation,
		ArtifactDigest:    Digest(result.ArtifactBytes),
		TraceDigest:       Digest(traceBytes),
		ArtifactPath:      filepath.ToSlash(filepath.Join(filepath.Base(caseDir), artifactName)),
		IRPath:            filepath.ToSlash(filepath.Join(filepath.Base(caseDir), irName)),
		TracePath:         filepath.ToSlash(filepath.Join(filepath.Base(caseDir), traceName)),
	}, nil
}

func BuildEvidence(metaPath, repoRoot, conformancePath, artifactDir, runtimeDir, evidencePath, buildMetric, testMetric, conformanceMetric, generatedBuildMetric, generatedRunMetric, testJSON, witnessMetricsPath string) error {
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
	if conformance.Schema != "gooo.conformance/v2" || conformance.ContractVersion != meta.Schema || conformance.ContractDigest != Digest(rawMeta) || conformance.SuiteStatus != statusClosed || len(conformance.Cases) != len(meta.FixedCases) || conformance.Denominator.FixedCaseCount != len(meta.FixedCases) || !conformance.Denominator.AppendOnly || conformance.Bootstrap.KernelSourceDigest != meta.SemanticKernel.Source.Digest || conformance.Bootstrap.CanonicalInputDigest != meta.SemanticKernel.CanonicalInput.Digest || conformance.Bootstrap.RequiredIndependentCount != meta.WitnessPlan.RequiredIndependentCount {
		return errors.New("conformance report is not closed or has the wrong fixed-case count")
	}
	if _, _, err := verifyDeclaredIdentities(meta, repoRoot); err != nil {
		return err
	}
	for _, observation := range conformance.Cases {
		if observation.ActualStatus != observation.ExpectedStatus {
			return fmt.Errorf("conformance case %s does not match its declared status", observation.ID)
		}
		if observation.ActualStatus == statusUnknown && (observation.Unknown == nil || !unknownComplete(observation.Unknown)) {
			return fmt.Errorf("conformance case %s has incomplete UNKNOWN evidence", observation.ID)
		}
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
	witnessRuntime, err := parseWitnessRuntimeMetrics(witnessMetricsPath, meta)
	if err != nil {
		return err
	}
	inside := 0
	if isInside(repoRoot, artifactDir) || isInside(repoRoot, runtimeDir) || isInside(repoRoot, evidencePath) {
		inside = 1
	}
	if inside != 0 {
		return errors.New("generated artifacts, runtime output, and evidence must remain caller-owned and outside the repository")
	}
	evidence := wire.Evidence{
		Schema:             "gooo.evidence/v2",
		ContractVersion:    meta.Schema,
		ContractDigest:     Digest(rawMeta),
		SuiteStatus:        statusClosed,
		Conformance:        conformance,
		Bootstrap:          conformance.Bootstrap,
		Inventory:          inventory,
		Runtime:            runtime,
		WitnessRuntime:     witnessRuntime,
		Tests:              tests,
		GeneratedArtifacts: generated,
		RuntimeExecution:   runtimeExecution,
		Improvement:        improvementUnknown(),
		Authority: wire.AuthorityMetrics{
			RepositoryWrites:                0,
			InputRepositoryWrites:           0,
			GeneratedOutputInsideRepository: inside,
			CrossProjectRequiredGates:       0,
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

func parseWitnessRuntimeMetrics(root string, meta wire.Meta) ([]wire.WitnessRuntimeObservation, error) {
	if root == "" {
		return nil, errors.New("witness metrics directory is required")
	}
	items := []struct {
		identity wire.WitnessIdentity
		name     string
		role     string
	}{
		{identity: meta.WitnessPlan.Stage0Reference, name: "path-a", role: "stage0-reference"},
		{identity: meta.WitnessPlan.GeneratedCurrent, name: "path-b", role: "generated-current"},
	}
	result := make([]wire.WitnessRuntimeObservation, 0, len(items))
	for _, item := range items {
		observe, err := parseMetric(filepath.Join(root, item.name+".observe.metric"))
		if err != nil {
			return nil, err
		}
		build, err := parseMetric(filepath.Join(root, item.name+".build.metric"))
		if err != nil {
			return nil, err
		}
		test, err := parseMetric(filepath.Join(root, item.name+".test.metric"))
		if err != nil {
			return nil, err
		}
		result = append(result, wire.WitnessRuntimeObservation{WitnessID: item.identity.ID, Role: item.role, ScenarioID: meta.SemanticKernel.CanonicalInput.ID, InputDigest: meta.SemanticKernel.CanonicalInput.Digest, ToolchainID: item.identity.ToolchainID, WallMS: observe.WallMS, PeakRSSKiB: observe.PeakRSSKiB, BuildMS: build.WallMS, TestMS: test.WallMS})
	}
	return result, nil
}

func improvementUnknown() wire.ImprovementObservation {
	return wire.ImprovementObservation{Status: statusUnknown, Improvement: nil, MatchedPair: false, Unknown: &wire.UnknownRecord{Stage: "measurement", Step: "match-before-after-witness-pair", Reason: "no exact before-after witness pair with the same scenario, input, contract, toolchain, witness, and trial identity is available", UnknownClass: "missing-exact-improvement-pair", NextOperation: "provide-exact-matched-before-after-pair-and-rerun", BlockedBy: "missing-matched-improvement-evidence"}}
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
	return record != nil && record.Stage != "" && record.Step != "" && record.Reason != "" && record.UnknownClass != "" && record.NextOperation != "" && record.BlockedBy != ""
}

func validDigest(value string) bool {
	return len(value) == len("sha256:")+64 && strings.HasPrefix(value, "sha256:") && !strings.Contains(value, "TO_BE_FILLED")
}

func completeWitness(identity wire.WitnessIdentity) bool {
	return identity.ID != "" && identity.Role != "" && identity.Implementation != "" && identity.SourcePath != "" && validDigest(identity.SourceDigest) && identity.LineageID != "" && identity.ToolchainID != "" && identity.InputIdentity != "" && identity.Available && identity.RequiredGate == 1
}

func verifyDeclaredIdentities(meta wire.Meta, repoRoot string) (string, string, error) {
	readDigest := func(relative, expected string) (string, error) {
		if filepath.IsAbs(relative) {
			return "", fmt.Errorf("identity path must be repository-relative: %s", relative)
		}
		data, err := os.ReadFile(filepath.Join(repoRoot, filepath.FromSlash(relative)))
		if err != nil {
			return "", err
		}
		actual := Digest(data)
		if actual != expected {
			return "", fmt.Errorf("identity digest mismatch for %s: declared %s observed %s", relative, expected, actual)
		}
		return actual, nil
	}
	kernelDigest, err := readDigest(meta.SemanticKernel.Source.Path, meta.SemanticKernel.Source.Digest)
	if err != nil {
		return "", "", fmt.Errorf("kernel source identity: %w", err)
	}
	inputDigest, err := readDigest(meta.SemanticKernel.CanonicalInput.Path, meta.SemanticKernel.CanonicalInput.Digest)
	if err != nil {
		return "", "", fmt.Errorf("canonical input identity: %w", err)
	}
	for _, identity := range []wire.WitnessIdentity{meta.WitnessPlan.Stage0Reference, meta.WitnessPlan.GeneratedCurrent} {
		if _, err := readDigest(identity.SourcePath, identity.SourceDigest); err != nil {
			return "", "", fmt.Errorf("witness %s identity: %w", identity.ID, err)
		}
	}
	return kernelDigest, inputDigest, nil
}

func witnessObservation(identity wire.WitnessIdentity, path wire.PathObservation, inputDigest, mode string) wire.WitnessObservation {
	result := wire.WitnessObservation{
		ID:                identity.ID,
		Role:              identity.Role,
		Implementation:    identity.Implementation,
		SourceDigest:      identity.SourceDigest,
		LineageID:         identity.LineageID,
		ToolchainID:       identity.ToolchainID,
		InputIdentity:     identity.InputIdentity,
		InputDigest:       inputDigest,
		SemanticDigest:    path.SemanticDigest,
		DecisionDigest:    path.DecisionDigest,
		ProvenanceDigest:  path.ProvenanceDigest,
		ObservationDigest: path.ObservationDigest,
	}
	if mode == "missing-lineage-toolchain-input" {
		result.LineageID = ""
		result.ToolchainID = ""
		result.InputIdentity = ""
		result.InputDigest = ""
	}
	if mode == "missing-input" {
		result.InputIdentity = ""
		result.InputDigest = ""
	}
	return result
}

func independentWitnesses(left, right wire.WitnessObservation) bool {
	return left.ID != "" && right.ID != "" && left.ID != right.ID && left.Role != right.Role && left.Implementation != right.Implementation && left.SourceDigest != "" && right.SourceDigest != "" && left.SourceDigest != right.SourceDigest && left.LineageID != "" && right.LineageID != "" && left.LineageID != right.LineageID && left.ToolchainID != "" && right.ToolchainID != "" && left.InputIdentity != "" && left.InputIdentity == right.InputIdentity && left.InputDigest != "" && left.InputDigest == right.InputDigest
}

func decisionDigest(trace wire.TerminalTrace) string {
	value := struct {
		TerminalReason string             `json:"terminal_reason"`
		Effects        []wire.EffectEvent `json:"effects"`
	}{TerminalReason: trace.TerminalReason, Effects: trace.Effects}
	raw, _ := canonical(value)
	return Digest(raw)
}

func provenanceDigest(sourceDigest, contractDigest, semanticDigest, decision string) string {
	value := struct {
		SourceDigest   string `json:"source_digest"`
		ContractDigest string `json:"contract_digest"`
		InputDigest    string `json:"input_digest"`
		SemanticDigest string `json:"semantic_digest"`
		DecisionDigest string `json:"decision_digest"`
	}{SourceDigest: sourceDigest, ContractDigest: contractDigest, InputDigest: sourceDigest, SemanticDigest: semanticDigest, DecisionDigest: decision}
	raw, _ := canonical(value)
	return Digest(raw)
}

func observationDigest(semantic, decision, provenance string) string {
	value := struct {
		SemanticDigest   string `json:"semantic_digest"`
		DecisionDigest   string `json:"decision_digest"`
		ProvenanceDigest string `json:"provenance_digest"`
	}{SemanticDigest: semantic, DecisionDigest: decision, ProvenanceDigest: provenance}
	raw, _ := canonical(value)
	return Digest(raw)
}

func denominatorFor(meta wire.Meta) wire.Denominator {
	counts := map[string]int{}
	ids := make([]string, 0, len(meta.FixedCases))
	for _, fixed := range meta.FixedCases {
		counts[fixed.ExpectedStatus]++
		ids = append(ids, fixed.ID)
	}
	return wire.Denominator{ID: meta.ContractID + "/denominator", Version: "v2", FixedCaseCount: len(meta.FixedCases), PreviousCaseCount: meta.ContractEvolution.PreviousDenominator, AddedCaseCount: len(meta.ContractEvolution.AddedFixedCaseIDs), RetiredCaseCount: len(meta.ContractEvolution.RetiredFixedCaseIDs), StatusCounts: counts, CaseIDs: ids, AppendOnly: meta.ContractEvolution.AppendOnly}
}

func bootstrapObservation(meta wire.Meta, kernelDigest, inputDigest string) wire.BootstrapObservation {
	return wire.BootstrapObservation{RequiredIndependentCount: meta.WitnessPlan.RequiredIndependentCount, KernelSourceDigest: kernelDigest, CanonicalInputDigest: inputDigest, Stage0Reference: meta.WitnessPlan.Stage0Reference, GeneratedCurrent: meta.WitnessPlan.GeneratedCurrent, OptionalDiverse: meta.WitnessPlan.OptionalDiverse, Trilemma: meta.WitnessPlan.Trilemma}
}

func caseProofReceipts(meta wire.Meta, observation wire.CaseObservation) []wire.ProofReceipt {
	result := make([]wire.ProofReceipt, 0, len(meta.WitnessPlan.ProofClaims))
	for _, claim := range meta.WitnessPlan.ProofClaims {
		status := statusClosed
		if observation.ActualStatus == statusUnknown && claim.ProofChoice == "COHERENCE" {
			status = statusUnknown
		}
		if observation.ActualStatus == statusRefuted && claim.ProofChoice == "REGRESSION" {
			status = statusClosed
		}
		value := struct {
			CaseID string `json:"case_id"`
			Claim  string `json:"claim"`
			Status string `json:"status"`
		}{CaseID: observation.ID, Claim: claim.ID, Status: status}
		raw, _ := canonical(value)
		result = append(result, wire.ProofReceipt{ClaimID: claim.ID, Claim: claim.Claim, ProofChoice: claim.ProofChoice, Status: status, IndependenceBasis: claim.IndependenceBasis, EvidenceDigest: Digest(raw)})
	}
	return result
}

func proofReceipts(meta wire.Meta, cases []wire.CaseObservation) []wire.ProofReceipt {
	result := make([]wire.ProofReceipt, 0, len(meta.WitnessPlan.ProofClaims))
	for _, claim := range meta.WitnessPlan.ProofClaims {
		statuses := make([]string, 0, len(cases))
		for _, observation := range cases {
			for _, receipt := range observation.ProofReceipts {
				if receipt.ClaimID == claim.ID {
					statuses = append(statuses, receipt.Status)
				}
			}
		}
		status := statusClosed
		for _, observed := range statuses {
			if observed == statusUnknown {
				status = statusUnknown
			}
		}
		value := struct {
			ClaimID string   `json:"claim_id"`
			Statuses []string `json:"statuses"`
		}{ClaimID: claim.ID, Statuses: statuses}
		raw, _ := canonical(value)
		result = append(result, wire.ProofReceipt{ClaimID: claim.ID, Claim: claim.Claim, ProofChoice: claim.ProofChoice, Status: status, IndependenceBasis: claim.IndependenceBasis, EvidenceDigest: Digest(raw)})
	}
	return result
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
