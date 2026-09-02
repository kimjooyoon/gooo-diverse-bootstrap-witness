package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/kimjooyoon/gooo-diverse-bootstrap-witness/internal/verifier"
)

func main() {
	if len(os.Args) < 2 {
		fail(errors.New("command is required: audit, observe, conformance, or evidence"))
	}
	switch os.Args[1] {
	case "audit":
		runAudit(os.Args[2:])
	case "conformance":
		runConformance(os.Args[2:])
	case "observe":
		runObserve(os.Args[2:])
	case "evidence":
		runEvidence(os.Args[2:])
	default:
		fail(fmt.Errorf("unknown command %q", os.Args[1]))
	}
}

func runObserve(args []string) {
	flags := flag.NewFlagSet("observe", flag.ContinueOnError)
	meta := flags.String("meta", ".gooo/diverse-bootstrap.gooo", "authoritative .gooo metacode")
	repoRoot := flags.String("repo-root", ".", "repository root")
	source := flags.String("source", "", "caller-selected .gooo input")
	witness := flags.String("witness", "", "declared witness identity")
	if err := flags.Parse(args); err != nil {
		fail(err)
	}
	if err := verifier.Observe(*meta, *repoRoot, *source, *witness); err != nil {
		fail(err)
	}
}

func runAudit(args []string) {
	flags := flag.NewFlagSet("audit", flag.ContinueOnError)
	meta := flags.String("meta", ".gooo/diverse-bootstrap.gooo", "authoritative .gooo metacode")
	if err := flags.Parse(args); err != nil {
		fail(err)
	}
	if _, _, err := verifier.LoadMeta(*meta); err != nil {
		fail(err)
	}
	fmt.Println("semantic-audit=CLOSED")
}

func runConformance(args []string) {
	flags := flag.NewFlagSet("conformance", flag.ContinueOnError)
	meta := flags.String("meta", ".gooo/diverse-bootstrap.gooo", "authoritative .gooo metacode")
	repoRoot := flags.String("repo-root", ".", "repository root containing semantic fixtures")
	artifactDir := flags.String("artifact-dir", "", "caller-owned output directory for generated artifacts")
	report := flags.String("report", "", "caller-owned conformance report")
	if err := flags.Parse(args); err != nil {
		fail(err)
	}
	if err := verifier.RunConformance(*meta, *repoRoot, *artifactDir, *report); err != nil {
		fail(err)
	}
}

func runEvidence(args []string) {
	flags := flag.NewFlagSet("evidence", flag.ContinueOnError)
	meta := flags.String("meta", ".gooo/diverse-bootstrap.gooo", "authoritative .gooo metacode")
	repoRoot := flags.String("repo-root", ".", "repository root")
	conformance := flags.String("conformance", "", "conformance report")
	artifactDir := flags.String("artifact-dir", "", "generated artifact directory")
	runtimeDir := flags.String("runtime-dir", "", "generated artifact runtime output directory")
	evidence := flags.String("evidence", "", "caller-owned evidence output")
	buildMetric := flags.String("build-metric", "", "Go build elapsed/RSS metric")
	testMetric := flags.String("test-metric", "", "Go test elapsed/RSS metric")
	conformanceMetric := flags.String("conformance-metric", "", "conformance elapsed/RSS metric")
	generatedBuildMetric := flags.String("generated-build-metric", "", "generated artifact build elapsed/RSS metric")
	generatedRunMetric := flags.String("generated-run-metric", "", "generated artifact run elapsed/RSS metric")
	testJSON := flags.String("test-json", "", "go test -json output")
	witnessMetrics := flags.String("witness-metrics", "", "caller-owned per-witness metric directory")
	if err := flags.Parse(args); err != nil {
		fail(err)
	}
	err := verifier.BuildEvidence(*meta, *repoRoot, *conformance, *artifactDir, *runtimeDir, *evidence, *buildMetric, *testMetric, *conformanceMetric, *generatedBuildMetric, *generatedRunMetric, *testJSON, *witnessMetrics)
	if err != nil {
		fail(err)
	}
}

func fail(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	os.Exit(1)
}
