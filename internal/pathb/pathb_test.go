package pathb

import (
	"bytes"
	"testing"

	"github.com/kimjooyoon/gooo-diverse-bootstrap-witness/internal/wire"
)

func pathBMeta() wire.Meta {
	return wire.Meta{
		Language: wire.LanguageSpec{
			CommentPrefixes: []string{"#", "//"},
			Keywords:        map[string]string{"program": "program", "binding": "const", "emission": "emit", "effect": "effect"},
		},
		Semantic: wire.SemanticSpec{IRSchema: "gooo.semantic-ir/v1"},
	}
}

func TestPathBSortsDeclarationMapForCanonicalIR(t *testing.T) {
	result, err := Generate(pathBMeta(), []byte("program x\nconst z = \"last\"\nconst a = \"first\"\n"), "normal")
	if err != nil {
		t.Fatal(err)
	}
	if len(result.IR.Bindings) != 2 || result.IR.Bindings[0].Name != "a" || result.IR.Bindings[1].Name != "z" {
		t.Fatalf("unexpected bindings: %#v", result.IR.Bindings)
	}
}

func TestPathBInjectionChangesSemanticBytes(t *testing.T) {
	clean, err := Generate(pathBMeta(), []byte("program x\nconst message = \"hello\"\n"), "normal")
	if err != nil {
		t.Fatal(err)
	}
	injected, err := Generate(pathBMeta(), []byte("program x\nconst message = \"hello\"\n"), "inject-self-propagating")
	if err != nil {
		t.Fatal(err)
	}
	if bytes.Equal(clean.ArtifactBytes, injected.ArtifactBytes) || clean.IR.Bindings[0].Value == injected.IR.Bindings[0].Value {
		t.Fatal("path-b injection was not observable")
	}
}
