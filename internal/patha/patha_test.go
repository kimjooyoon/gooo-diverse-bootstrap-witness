package patha

import (
	"bytes"
	"testing"

	"github.com/kimjooyoon/gooo-diverse-bootstrap-witness/internal/wire"
)

func pathAMeta() wire.Meta {
	return wire.Meta{
		Language: wire.LanguageSpec{
			CommentPrefixes: []string{"#", "//"},
			Keywords:        map[string]string{"program": "program", "binding": "const", "emission": "emit", "effect": "effect"},
		},
		Semantic: wire.SemanticSpec{IRSchema: "gooo.semantic-ir/v1"},
	}
}

func TestCommentDoesNotChangeCanonicalIR(t *testing.T) {
	first, err := Generate(pathAMeta(), []byte("program x\nconst b = \"two\"\nconst a = \"one\"\nemit a\n"), "normal")
	if err != nil {
		t.Fatal(err)
	}
	second, err := Generate(pathAMeta(), []byte("# note\nprogram x\nconst b = \"two\" // note\nconst a = \"one\"\nemit a\n"), "normal")
	if err != nil {
		t.Fatal(err)
	}
	if string(first.ArtifactBytes) != string(second.ArtifactBytes) || !bytes.Equal(first.ArtifactBytes, second.ArtifactBytes) {
		t.Fatal("comment changed path-a artifact")
	}
}

func TestPathAProducesTerminalTrace(t *testing.T) {
	result, err := Generate(pathAMeta(), []byte("program x\neffect audit = \"ok\"\n"), "normal")
	if err != nil {
		t.Fatal(err)
	}
	if result.Trace.TerminalReason != "complete:program:x" || len(result.Trace.Effects) != 1 {
		t.Fatalf("unexpected trace: %#v", result.Trace)
	}
}
