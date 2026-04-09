package lsp

import (
	"strings"

	"moonbasic/compiler/pipeline"
)

// documentSymbols returns LSP DocumentSymbol[] (flat list) for outline view.
func (s *server) documentSymbols(uri, text string) []any {
	path := filePathFromURI(uri)
	if path == "" {
		path = "buffer.mb"
	}
	syms, err := pipeline.DocumentSymbols(path, text)
	if err != nil {
		return nil
	}
	out := make([]any, 0, len(syms))
	for _, m := range syms {
		name, _ := m["name"].(string)
		kind, _ := m["kind"].(int)
		line := bestEffortDefinitionLine(text, name)
		out = append(out, map[string]any{
			"name": name,
			"kind": kind,
			"range": map[string]any{
				"start": map[string]int{"line": line, "character": 0},
				"end":   map[string]int{"line": line, "character": minInt(80, len(name)+10)},
			},
			"selectionRange": map[string]any{
				"start": map[string]int{"line": line, "character": 0},
				"end":   map[string]int{"line": line, "character": len(name)},
			},
		})
	}
	return out
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// definitionLocation returns TextDocumentDefinition result (Location) or nil.
func (s *server) definitionLocation(uri, text string, line, character int) any {
	lines := strings.Split(text, "\n")
	if line < 0 || line >= len(lines) {
		return nil
	}
	ln := lines[line]
	word := identifierAt(ln, character)
	if word == "" {
		return nil
	}
	// Builtins: delegate to hover manifest only for dotted keys
	if strings.Contains(word, ".") {
		return nil
	}
	l := bestEffortDefinitionLine(text, word)
	if l < 0 {
		return nil
	}
	return map[string]any{
		"uri": uri,
		"range": map[string]any{
			"start": map[string]int{"line": l, "character": 0},
			"end":   map[string]int{"line": l, "character": len(lines[l])},
		},
	}
}

func identifierAt(line string, col int) string {
	if col < 0 || col >= len(line) {
		return ""
	}
	start := col
	for start > 0 && (isIdentRune(rune(line[start-1]))) {
		start--
	}
	end := col
	for end < len(line) && isIdentRune(rune(line[end])) {
		end++
	}
	if start >= end {
		return ""
	}
	return strings.ToUpper(strings.TrimSpace(line[start:end]))
}

func isIdentRune(r rune) bool {
	return (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '#' || r == '$' || r == '_' || r == '?'
}

// bestEffortDefinitionLine finds a line declaring name (FUNCTION, TYPE, CONST, or first assignment).
func bestEffortDefinitionLine(src, name string) int {
	u := strings.ToUpper(strings.TrimSpace(name))
	if u == "" {
		return -1
	}
	lines := strings.Split(src, "\n")
	for i, ln := range lines {
		ul := strings.ToUpper(strings.TrimSpace(ln))
		if strings.HasPrefix(ul, "FUNCTION "+u) || strings.HasPrefix(ul, "FUNCTION "+u+"(") {
			return i
		}
		if strings.HasPrefix(ul, "TYPE "+u) {
			return i
		}
		if strings.HasPrefix(ul, "CONST "+u) || strings.HasPrefix(ul, "CONST "+u+" ") {
			return i
		}
		if strings.HasPrefix(ul, "DIM "+u) || strings.HasPrefix(ul, "DIM "+u+"(") {
			return i
		}
	}
	// Global assignment: "name = " or "name="
	for i, ln := range lines {
		ul := strings.ToUpper(strings.TrimSpace(ln))
		if strings.HasPrefix(ul, u+"=") || strings.HasPrefix(ul, u+" =") {
			return i
		}
	}
	return 0
}
