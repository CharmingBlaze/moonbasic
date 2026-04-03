package parser

import (
	"strings"

	"moonbasic/compiler/arena"
	"moonbasic/compiler/ast"
	"moonbasic/compiler/symtable"
	"moonbasic/compiler/token"
)

// parseFunctionDef parses: FUNCTION name(params) body ENDFUNCTION
func (p *Parser) parseFunctionDef() (*ast.FunctionDef, error) {
	line, col := p.cur().Line, p.cur().Col
	p.advance() // consume FUNCTION
	name, err := p.expectIdent()
	if err != nil {
		return nil, err
	}
	if err := p.expect(token.LPAREN); err != nil {
		return nil, err
	}
	var params []ast.Param
	p.skipNewlines()
	if p.cur().Type != token.RPAREN {
		for {
			pname, err2 := p.expectIdent()
			if err2 != nil {
				return nil, err2
			}
			suf := ""
			if len(pname) > 0 {
				last := pname[len(pname)-1]
				if last == '#' || last == '$' || last == '?' {
					suf = string(last)
				}
			}
			params = append(params, ast.Param{Name: pname, Suffix: suf})
			p.skipNewlines()
			if p.cur().Type == token.COMMA {
				p.advance()
				p.skipNewlines()
				continue
			}
			break
		}
	}
	if err := p.expect(token.RPAREN); err != nil {
		return nil, err
	}
	if p.cur().Type != token.NEWLINE {
		return nil, p.failf("expected newline after FUNCTION header")
	}
	p.advance()

	// Enter a fresh scope for the function body.
	// We preserve predeclared funcs/types for forward references.
	savedSym := p.sym
	p.sym = symtable.New()
	for k := range savedSym.Funcs() {
		p.sym.PredeclareFunction(k)
	}
	for k := range savedSym.Types() {
		p.sym.PredeclareType(k)
	}
	p.sym.PushScope()
	for _, par := range params {
		p.sym.DefineParam(par.Name)
	}
	p.sym.DefineFunction(name)

	savedFn := p.FuncName
	p.FuncName = strings.ToUpper(name)
	defer func() { p.FuncName = savedFn }()

	var body []ast.Stmt
	for {
		p.skipNewlines()
		if p.cur().Type == token.ENDFUNCTION {
			break
		}
		if p.cur().Type == token.EOF {
			return nil, p.failf("unexpected EOF inside FUNCTION %s", name)
		}
		s, err2 := p.parseStmt()
		if err2 != nil {
			return nil, err2
		}
		if s != nil {
			body = append(body, s)
		}
	}
	p.advance() // consume ENDFUNCTION
	p.sym.PopScope()
	p.sym = savedSym
	p.sym.DefineFunction(name)

	return arena.Make(p.ar, ast.FunctionDef{Name: name, Params: params, Body: body, Line: line, Col: col}), nil
}

// parseTypeDef parses: TYPE name FIELD ... ENDTYPE
func (p *Parser) parseTypeDef() (*ast.TypeDef, error) {
	line, col := p.cur().Line, p.cur().Col
	p.advance() // consume TYPE
	tname, err := p.expectIdent()
	if err != nil {
		return nil, err
	}
	if p.cur().Type != token.NEWLINE {
		return nil, p.failf("expected newline after TYPE name")
	}
	p.advance()
	var fields []string
	for {
		p.skipNewlines()
		if p.cur().Type == token.ENDTYPE {
			p.advance()
			break
		}
		if p.cur().Type != token.FIELD {
			return nil, p.failf("expected FIELD or ENDTYPE, got %s", p.cur().Type.String())
		}
		p.advance()
		for {
			fname, err2 := p.expectIdent()
			if err2 != nil {
				return nil, err2
			}
			fields = append(fields, fname)
			if p.cur().Type == token.COMMA {
				p.advance()
				p.skipNewlines()
				continue
			}
			break
		}
		if p.cur().Type == token.NEWLINE {
			p.advance()
		}
	}
	p.sym.DefineType(tname)
	return arena.Make(p.ar, ast.TypeDef{Name: tname, Fields: fields, Line: line, Col: col}), nil
}
