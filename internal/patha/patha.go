// Package patha is one complete Gooo parser/lowerer/emitter/executor path.
// It intentionally has no dependency on pathb or on a shared evaluator.
package patha

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

type sourceProgram struct {
	name      string
	bindings  []sourceBinding
	emissions []string
	effects   []sourceEffect
}

type sourceBinding struct {
	name  string
	value string
}

type sourceEffect struct {
	name  string
	value string
}

func Generate(meta wire.Meta, source []byte, variant string) (wire.GeneratedResult, error) {
	return generate(meta, source, variant, "", "", "")
}

// GenerateCase carries the authoritative case labels into the IR that is
// emitted for the verifier and retained in release evidence.
func GenerateCase(meta wire.Meta, source []byte, variant, caseID, proofChoice, indicatorClass string) (wire.GeneratedResult, error) {
	return generate(meta, source, variant, caseID, proofChoice, indicatorClass)
}

func generate(meta wire.Meta, source []byte, variant, caseID, proofChoice, indicatorClass string) (wire.GeneratedResult, error) {
	parsed, err := parse(meta, source)
	if err != nil {
		return wire.GeneratedResult{}, err
	}
	lowered, err := lower(meta, parsed, caseID, proofChoice, indicatorClass)
	if err != nil {
		return wire.GeneratedResult{}, err
	}
	trace := execute(lowered)
	artifact, err := emit(lowered, trace)
	if err != nil {
		return wire.GeneratedResult{}, err
	}
	return wire.GeneratedResult{IR: lowered, ArtifactBytes: artifact, Trace: trace}, nil
}

func parse(meta wire.Meta, source []byte) (sourceProgram, error) {
	if meta.Language.Keywords["program"] == "" || meta.Language.Keywords["binding"] == "" || meta.Language.Keywords["emission"] == "" || meta.Language.Keywords["effect"] == "" {
		return sourceProgram{}, errors.New("path-a metacode is missing a keyword rule")
	}
	var program sourceProgram
	seenProgram := false
	scanner := bufio.NewScanner(bytes.NewReader(source))
	lineNumber := 0
	for scanner.Scan() {
		lineNumber++
		line := strings.TrimSpace(removeComment(scanner.Text(), meta.Language.CommentPrefixes))
		if line == "" {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) == 0 {
			continue
		}
		switch fields[0] {
		case meta.Language.Keywords["program"]:
			if len(fields) != 2 || !identifier(fields[1]) || seenProgram {
				return sourceProgram{}, fmt.Errorf("path-a line %d: invalid program declaration", lineNumber)
			}
			program.name = fields[1]
			seenProgram = true
		case meta.Language.Keywords["binding"]:
			name, value, err := assignment(line, meta.Language.Keywords["binding"])
			if err != nil || !identifier(name) {
				return sourceProgram{}, fmt.Errorf("path-a line %d: invalid binding", lineNumber)
			}
			program.bindings = append(program.bindings, sourceBinding{name: name, value: value})
		case meta.Language.Keywords["emission"]:
			if len(fields) != 2 || !identifier(fields[1]) {
				return sourceProgram{}, fmt.Errorf("path-a line %d: invalid emission", lineNumber)
			}
			program.emissions = append(program.emissions, fields[1])
		case meta.Language.Keywords["effect"]:
			name, value, err := assignment(line, meta.Language.Keywords["effect"])
			if err != nil || !identifier(name) {
				return sourceProgram{}, fmt.Errorf("path-a line %d: invalid effect", lineNumber)
			}
			program.effects = append(program.effects, sourceEffect{name: name, value: value})
		default:
			return sourceProgram{}, fmt.Errorf("path-a line %d: unknown rule %q", lineNumber, fields[0])
		}
	}
	if err := scanner.Err(); err != nil {
		return sourceProgram{}, err
	}
	if !seenProgram {
		return sourceProgram{}, errors.New("path-a source has no program")
	}
	return program, nil
}

func lower(meta wire.Meta, parsed sourceProgram, caseID, proofChoice, indicatorClass string) (wire.SemanticIR, error) {
	bindings := append([]sourceBinding(nil), parsed.bindings...)
	sort.SliceStable(bindings, func(i, j int) bool { return bindings[i].name < bindings[j].name })
	ir := wire.SemanticIR{Schema: meta.Semantic.IRSchema, Program: parsed.name, CaseID: caseID, ProofChoice: proofChoice, IndicatorClass: indicatorClass}
	for _, binding := range bindings {
		ir.Bindings = append(ir.Bindings, wire.Binding{Name: binding.name, Value: binding.value})
	}
	for _, name := range parsed.emissions {
		ir.Emissions = append(ir.Emissions, name)
		if !containsBinding(ir.Bindings, name) {
			ir.Diagnostics = append(ir.Diagnostics, wire.Diagnostic{Code: "unbound-emission", Message: name})
		}
	}
	for _, effect := range parsed.effects {
		ir.Effects = append(ir.Effects, wire.Binding{Name: effect.name, Value: effect.value})
	}
	return ir, nil
}

func execute(ir wire.SemanticIR) wire.TerminalTrace {
	trace := wire.TerminalTrace{Schema: "gooo.terminal-trace/v1", TerminalReason: "complete:program:" + ir.Program}
	for _, effect := range ir.Effects {
		trace.Effects = append(trace.Effects, wire.EffectEvent{Kind: "declared-effect", Name: effect.Name, Value: effect.Value})
	}
	for _, emission := range ir.Emissions {
		value := ""
		for _, binding := range ir.Bindings {
			if binding.Name == emission {
				value = binding.Value
				break
			}
		}
		trace.Effects = append(trace.Effects, wire.EffectEvent{Kind: "emission", Name: emission, Value: value})
	}
	return trace
}

func emit(ir wire.SemanticIR, trace wire.TerminalTrace) ([]byte, error) {
	var builder strings.Builder
	builder.WriteString("package main\n\nimport \"fmt\"\n\nfunc main() {\n")
	fmt.Fprintf(&builder, "\tfmt.Println(%s)\n", strconv.Quote("GOOO_PROGRAM\t"+ir.Program))
	for _, binding := range ir.Bindings {
		fmt.Fprintf(&builder, "\tfmt.Println(%s)\n", strconv.Quote("GOOO_BINDING\t"+binding.Name+"\t"+binding.Value))
	}
	for _, emission := range ir.Emissions {
		fmt.Fprintf(&builder, "\tfmt.Println(%s)\n", strconv.Quote("GOOO_EMIT\t"+emission))
	}
	for _, effect := range trace.Effects {
		fmt.Fprintf(&builder, "\tfmt.Println(%s)\n", strconv.Quote("GOOO_EFFECT\t"+effect.Kind+"\t"+effect.Name+"\t"+effect.Value))
	}
	fmt.Fprintf(&builder, "\tfmt.Println(%s)\n", strconv.Quote("GOOO_TRACE\t"+trace.TerminalReason))
	builder.WriteString("}\n")
	return []byte(builder.String()), nil
}

func assignment(line, keyword string) (string, string, error) {
	rest := strings.TrimSpace(strings.TrimPrefix(line, keyword))
	parts := strings.SplitN(rest, "=", 2)
	if len(parts) != 2 {
		return "", "", errors.New("assignment requires equals")
	}
	name := strings.TrimSpace(parts[0])
	rawValue := strings.TrimSpace(parts[1])
	value, err := strconv.Unquote(rawValue)
	if err != nil {
		return "", "", err
	}
	return name, value, nil
}

func removeComment(line string, prefixes []string) string {
	inString := false
	escaped := false
	for index := 0; index < len(line); index++ {
		char := line[index]
		if char == '\\' && inString && !escaped {
			escaped = true
			continue
		}
		if char == '"' && !escaped {
			inString = !inString
		}
		escaped = false
		if inString {
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

func identifier(value string) bool {
	if value == "" || (value[0] < 'a' || value[0] > 'z') && (value[0] < 'A' || value[0] > 'Z') && value[0] != '_' {
		return false
	}
	for _, char := range value[1:] {
		if (char < 'a' || char > 'z') && (char < 'A' || char > 'Z') && (char < '0' || char > '9') && char != '_' {
			return false
		}
	}
	return true
}

func containsBinding(bindings []wire.Binding, name string) bool {
	for _, binding := range bindings {
		if binding.Name == name {
			return true
		}
	}
	return false
}
