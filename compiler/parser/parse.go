package parser

import (
	"strings"

	"moonbasic/compiler/arena"
	"moonbasic/compiler/ast"
	"moonbasic/compiler/lexer"
)

func normalized(src string) string {
	n := strings.ReplaceAll(src, "\r\n", "\n")
	return strings.ReplaceAll(n, "\r", "\n")
}

// SplitLines normalises line endings the same way as ParseSource and splits into lines.
func SplitLines(src string) []string {
	return strings.Split(normalized(src), "\n")
}

// ParseSource lexes and parses a complete source file, returning its AST.
func ParseSource(file, src string) (*ast.Program, error) {
	return ParseSourceWithArena(file, src, nil)
}

// ParseSourceWithArena parses using ar for AST allocations (nil = heap only).
func ParseSourceWithArena(file, src string, ar *arena.Arena) (*ast.Program, error) {
	norm := normalized(src)
	lines := strings.Split(norm, "\n")

	toks, err := lexer.Scan(file, norm)
	if err != nil {
		return nil, err
	}
	p := NewParserWithArena(file, toks, lines, ar)
	return p.ParseProgram()
}

// MustParse is for tests; panics on error.
func MustParse(file, src string) *ast.Program {
	prog, err := ParseSource(file, src)
	if err != nil {
		panic(err)
	}
	return prog
}

// StmtCount returns the number of top-level statements in the program.
func StmtCount(p *ast.Program) int {
	return len(p.Stmts)
}
