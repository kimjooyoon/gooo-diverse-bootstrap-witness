// Package pathb is a second complete Gooo parser/lowerer/emitter/executor
// path. Its implementation is intentionally independent from patha.
package pathb

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/kimjooyoon/gooo-diverse-bootstrap-witness/internal/wire"
)

type statement struct {
	kind  string
	name  string
	value string
}

type programModel struct {
	name         string
	declarations map[string]string
	emits        []string
	effects      []statement
}

func Generate(meta wire.Meta, source []byte, variant string) (wire.GeneratedResult, error) {
	statements, err := readStatements(meta, source)
	if err != nil {
		return wire.GeneratedResult{}, err
	}
	model, err := assemble(meta, statements)
	if err != nil {
		return wire.GeneratedResult{}, err
	}
	ir := makeIR(meta, model, variant)
	trace := runExecutor(meta, ir, variant)
	artifact := makeArtifact(ir, trace)
	return wire.GeneratedResult{IR: ir, ArtifactBytes: artifact, Trace: trace}, nil
}

func readStatements(meta wire.Meta, source []byte) ([]statement, error) {
	rules := meta.Language.Keywords
	if rules["program"] == "" || rules["binding"] == "" || rules["emission"] == "" || rules["effect"] == "" {
		return nil, errors.New("path-b metacode is missing a keyword rule")
	}
	reader := bufio.NewScanner(bytes.NewReader(source))
	var result []statement
	line := 0
	for reader.Scan() {
		line++
		text := strings.TrimSpace(trimComment(reader.Text(), meta.Language.CommentPrefixes))
		if text == "" {
			continue
		}
		fields := strings.Fields(text)
		if len(fields) == 0 {
			continue
		}
		switch fields[0] {
		case rules["program"]:
			if len(fields) != 2 || !isIdentifier(fields[1]) {
				return nil, fmt.Errorf("path-b line %d: malformed program", line)
			}
			result = append(result, statement{kind: "program", name: fields[1]})
		case rules["binding"], rules["effect"]:
			name, value, parseErr := parseAssignment(text, fields[0])
			if parseErr != nil || !isIdentifier(name) {
				return nil, fmt.Errorf("path-b line %d: malformed assignment", line)
			}
			kind := "binding"
			if fields[0] == rules["effect"] {
				kind = "effect"
			}
			result = append(result, statement{kind: kind, name: name, value: value})
		case rules["emission"]:
			if len(fields) != 2 || !isIdentifier(fields[1]) {
				return nil, fmt.Errorf("path-b line %d: malformed emission", line)
			}
			result = append(result, statement{kind: "emission", name: fields[1]})
		default:
			return nil, fmt.Errorf("path-b line %d: unsupported rule %q", line, fields[0])
		}
	}
	if err := reader.Err(); err != nil {
		return nil, err
	}
	return result, nil
}

func assemble(meta wire.Meta, statements []statement) (programModel, error) {
	model := programModel{declarations: map[string]string{}}
	for _, item := range statements {
		switch item.kind {
		case "program":
			if model.name != "" {
				return programModel{}, errors.New("path-b has more than one program")
			}
			model.name = item.name
		case "binding":
			if _, exists := model.declarations[item.name]; exists {
				return programModel{}, fmt.Errorf("path-b duplicate binding %q", item.name)
			}
			model.declarations[item.name] = item.value
		case "emission":
			model.emits = append(model.emits, item.name)
		case "effect":
			model.effects = append(model.effects, item)
		}
	}
	if model.name == "" {
		return programModel{}, errors.New("path-b source has no program")
	}
	return model, nil
}

func makeIR(meta wire.Meta, model programModel, variant string) wire.SemanticIR {
	mutation := findMutation(meta, variant)
	keys := make([]string, 0, len(model.declarations))
	for key := range model.declarations {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	ir := wire.SemanticIR{Schema: meta.Semantic.IRSchema, Program: model.name}
	for _, key := range keys {
		value := model.declarations[key]
		if mutation.Effect == "replace-binding-value" && key == mutation.Target {
			value = mutation.Value
		}
		ir.Bindings = append(ir.Bindings, wire.Binding{Name: key, Value: value})
	}
	for _, emission := range model.emits {
		ir.Emissions = append(ir.Emissions, emission)
		if _, exists := model.declarations[emission]; !exists {
			ir.Diagnostics = append(ir.Diagnostics, wire.Diagnostic{Code: "unbound-emission", Message: emission})
		}
	}
	for _, effect := range model.effects {
		ir.Effects = append(ir.Effects, wire.Binding{Name: effect.name, Value: effect.value})
	}
	return ir
}

func runExecutor(meta wire.Meta, ir wire.SemanticIR, variant string) wire.TerminalTrace {
	reason := "complete:program:" + ir.Program
	mutation := findMutation(meta, variant)
	if mutation.Effect == "replace-terminal-reason" {
		reason = mutation.Value + ir.Program
	}
	trace := wire.TerminalTrace{Schema: "gooo.terminal-trace/v1", TerminalReason: reason}
	for _, effect := range ir.Effects {
		trace.Effects = append(trace.Effects, wire.EffectEvent{Kind: "declared-effect", Name: effect.Name, Value: effect.Value})
	}
	for _, name := range ir.Emissions {
		value := ""
		for _, binding := range ir.Bindings {
			if binding.Name == name {
				value = binding.Value
				break
			}
		}
		trace.Effects = append(trace.Effects, wire.EffectEvent{Kind: "emission", Name: name, Value: value})
	}
	return trace
}

func findMutation(meta wire.Meta, variant string) wire.MutationSpec {
	for _, mutation := range meta.SemanticKernel.Evaluation.Mutations {
		if mutation.ID == variant {
			return mutation
		}
	}
	return wire.MutationSpec{}
}

func makeArtifact(ir wire.SemanticIR, trace wire.TerminalTrace) []byte {
	var output strings.Builder
	output.WriteString("package main\n\nimport \"fmt\"\n\nfunc main() {\n")
	fmt.Fprintf(&output, "\tfmt.Println(%s)\n", strconv.Quote("GOOO_PROGRAM\t"+ir.Program))
	for _, binding := range ir.Bindings {
		fmt.Fprintf(&output, "\tfmt.Println(%s)\n", strconv.Quote("GOOO_BINDING\t"+binding.Name+"\t"+binding.Value))
	}
	for _, emission := range ir.Emissions {
		fmt.Fprintf(&output, "\tfmt.Println(%s)\n", strconv.Quote("GOOO_EMIT\t"+emission))
	}
	for _, effect := range trace.Effects {
		fmt.Fprintf(&output, "\tfmt.Println(%s)\n", strconv.Quote("GOOO_EFFECT\t"+effect.Kind+"\t"+effect.Name+"\t"+effect.Value))
	}
	fmt.Fprintf(&output, "\tfmt.Println(%s)\n", strconv.Quote("GOOO_TRACE\t"+trace.TerminalReason))
	output.WriteString("}\n")
	return []byte(output.String())
}

func parseAssignment(line, keyword string) (string, string, error) {
	body := strings.TrimSpace(strings.TrimPrefix(line, keyword))
	parts := strings.SplitN(body, "=", 2)
	if len(parts) != 2 {
		return "", "", errors.New("assignment missing equals")
	}
	value, err := strconv.Unquote(strings.TrimSpace(parts[1]))
	if err != nil {
		return "", "", err
	}
	return strings.TrimSpace(parts[0]), value, nil
}

func trimComment(line string, prefixes []string) string {
	quoted := false
	escape := false
	for index := 0; index < len(line); index++ {
		if line[index] == '\\' && quoted && !escape {
			escape = true
			continue
		}
		if line[index] == '"' && !escape {
			quoted = !quoted
		}
		escape = false
		if quoted {
			continue
		}
		for _, prefix := range prefixes {
			if strings.HasPrefix(line[index:], prefix) {
				return line[:index]
			}
		}
	}
	return line
}

func isIdentifier(value string) bool {
	if value == "" {
		return false
	}
	for position, char := range value {
		if position == 0 {
			if (char < 'a' || char > 'z') && (char < 'A' || char > 'Z') && char != '_' {
				return false
			}
			continue
		}
		if (char < 'a' || char > 'z') && (char < 'A' || char > 'Z') && (char < '0' || char > '9') && char != '_' {
			return false
		}
	}
	return true
}
