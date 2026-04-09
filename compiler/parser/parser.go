// Package parser builds an AST from moonBASIC tokens.
package parser

import (
	"fmt"
	"strconv"
	"strings"

	"moonbasic/compiler/arena"
	"moonbasic/compiler/ast"
	"moonbasic/compiler/errors"
	"moonbasic/compiler/symtable"
	"moonbasic/compiler/token"
)

// Parser consumes a token slice produced by the lexer.
type Parser struct {
	toks    []token.Token
	i       int
	file    string
	lines   []string
	sym     *symtable.Table
	curLine string
	ar      *arena.Arena
	// FuncName is the enclosing FUNCTION name (uppercase), or "" at module scope.
	FuncName string
}

// NewParser creates a parser. lines are raw source lines for error context (1-based indexing optional).
func NewParser(file string, toks []token.Token, lines []string) *Parser {
	return NewParserWithArena(file, toks, lines, nil)
}

// NewParserWithArena creates a parser that allocates AST nodes from ar (optional).
func NewParserWithArena(file string, toks []token.Token, lines []string, ar *arena.Arena) *Parser {
	return &Parser{
		toks:  toks,
		i:     0,
		file:  file,
		lines: lines,
		sym:   symtable.New(),
		ar:    ar,
	}
}

// ParseProgram parses the full program.
func (p *Parser) ParseProgram() (*ast.Program, error) {
	Predeclare(p.sym, p.toks)
	prog := arena.Make(p.ar, ast.Program{})
	for {
		p.skipNewlines()
		if p.cur().Type == token.EOF {
			break
		}
		switch p.cur().Type {
		case token.FUNCTION:
			f, err := p.parseFunctionDef()
			if err != nil {
				return nil, err
			}
			prog.Functions = append(prog.Functions, f)
		case token.TYPE:
			td, err := p.parseTypeDef()
			if err != nil {
				return nil, err
			}
			prog.Types = append(prog.Types, td)
		default:
			s, err := p.parseStmt()
			if err != nil {
				return nil, err
			}
			if s != nil {
				prog.Stmts = append(prog.Stmts, s)
			}
			// Colon separates statements on one line: x = 1 : y = 2
			for p.cur().Type == token.COLON {
				p.advance()
				p.skipNewlines()
				s2, err2 := p.parseStmt()
				if err2 != nil {
					return nil, err2
				}
				if s2 != nil {
					prog.Stmts = append(prog.Stmts, s2)
				}
			}
		}
	}
	return prog, nil
}

// Predeclare scans tokens for FUNCTION and TYPE names.
func Predeclare(sym *symtable.Table, toks []token.Token) {
	for i := 0; i < len(toks); i++ {
		switch toks[i].Type {
		case token.FUNCTION:
			if i+1 < len(toks) && toks[i+1].Type == token.IDENT {
				sym.PredeclareFunction(toks[i+1].Lit)
			}
		case token.TYPE:
			if i+1 < len(toks) && toks[i+1].Type == token.IDENT {
				sym.PredeclareType(toks[i+1].Lit)
			}
		}
	}
}

func (p *Parser) cur() token.Token {
	if p.i >= len(p.toks) {
		return token.Token{Type: token.EOF}
	}
	return p.toks[p.i]
}

func (p *Parser) advance() {
	if p.i < len(p.toks) {
		p.i++
	}
}

func (p *Parser) save() int {
	return p.i
}

func (p *Parser) restore(i int) {
	p.i = i
}

func (p *Parser) skipNewlines() {
	for p.cur().Type == token.NEWLINE {
		p.advance()
	}
}

func (p *Parser) expect(tt token.TokenType) error {
	p.skipNewlines()
	t := p.cur()
	if t.Type != tt {
		return p.failf("expected %s, got %s", tt.String(), t.Type.String())
	}
	p.advance()
	return nil
}

func (p *Parser) expectIdent() (string, error) {
	p.skipNewlines()
	t := p.cur()
	// Accept both IDENT tokens AND keyword tokens as identifier names.
	// In moonBASIC, keywords can legally appear as method/field names after a dot,
	// e.g. cam.Step(), robot.Free(), Physics3D.Stop().
	if t.Type != token.IDENT && !isKeywordUsableAsIdent(t.Type) {
		return "", p.failf("expected identifier, got %s", t.Type.String())
	}
	p.advance()
	return t.Lit, nil
}

// isKeywordUsableAsIdent returns true for keywords that can appear as method or field names.
// This covers all command-word keywords that a user might name a method after.
func isKeywordUsableAsIdent(tt token.TokenType) bool {
	switch tt {
	case token.IF, token.THEN, token.ELSE, token.ELSEIF, token.ENDIF,
		token.WHILE, token.WEND, token.ENDWHILE,
		token.FOR, token.TO, token.DOWNTO, token.STEP, token.NEXT,
		token.REPEAT, token.UNTIL, token.DO, token.LOOP, token.EXIT, token.CONTINUE,
		token.SELECT, token.CASE, token.DEFAULT, token.ENDSELECT,
		token.FUNCTION, token.ENDFUNCTION, token.RETURN,
		token.TYPE, token.FIELD, token.ENDTYPE,
		token.NEW, token.DELETE, token.EACH,
		token.GOTO, token.GOSUB,
		token.AND, token.OR, token.NOT, token.XOR, token.MOD,
		token.DIM, token.AS, token.REDIM, token.PRESERVE, token.LOCAL, token.GLOBAL, token.CONST,
		token.STATIC, token.SWAP, token.ERASE,
		token.INCLUDE, token.TRUE, token.FALSE, token.NULL, token.END:
		return true
	}
	return false
}

func (p *Parser) sourceLine(line int) string {
	if line < 1 || line > len(p.lines) {
		return ""
	}
	return p.lines[line-1]
}

func (p *Parser) failf(format string, args ...interface{}) error {
	t := p.cur()
	msg := fmt.Sprintf(format, args...)
	me := errors.NewParseError(p.file, t.Line, t.Col, msg, p.sourceLine(t.Line),
		"")
	return me
}

func (p *Parser) defineAssignedName(name string) {
	if p.sym.ScopeDepth() > 0 {
		p.sym.DefineLocal(name)
	} else {
		p.sym.DefineGlobalVar(name)
	}
}

func (p *Parser) failMissingParen(name string) error {
	t := p.cur()
	me := errors.NewParseError(p.file, t.Line, t.Col,
		fmt.Sprintf("Expected '(' after '%s'", name),
		p.sourceLine(t.Line),
		fmt.Sprintf("All commands require parentheses: %s(args)", name))
	return me
}

// parseExprUntilColon parses expressions for single-line IF body (until COLON or NEWLINE at depth 0).
func (p *Parser) parseNumberLit(lit string, line, col int) (ast.Expr, error) {
	if strings.Contains(lit, ".") {
		f, err := strconv.ParseFloat(lit, 64)
		if err != nil {
			return nil, fmt.Errorf("float: %w", err)
		}
		return arena.Make(p.ar, ast.FloatLitNode{Value: f, Lit: lit, Line: line, Col: col}), nil
	}
	v, err := strconv.ParseInt(lit, 10, 64)
	if err != nil {
		return nil, fmt.Errorf("int: %w", err)
	}
	return arena.Make(p.ar, ast.IntLitNode{Value: v, Line: line, Col: col}), nil
}
