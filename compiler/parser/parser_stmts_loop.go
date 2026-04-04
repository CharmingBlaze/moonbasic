package parser

import (
	"moonbasic/compiler/arena"
	"moonbasic/compiler/ast"
	"moonbasic/compiler/token"
)

func (p *Parser) parseWhile() (ast.Stmt, error) {
	line, col := p.cur().Line, p.cur().Col
	p.advance()
	cond, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	if p.cur().Type != token.NEWLINE {
		return nil, p.failf("expected newline after WHILE condition")
	}
	p.advance()
	body, err := p.parseStmtBlockUntil([]token.TokenType{token.WEND, token.ENDWHILE})
	if err != nil {
		return nil, err
	}
	if p.cur().Type == token.WEND || p.cur().Type == token.ENDWHILE {
		p.advance()
	} else {
		return nil, p.failf("expected WEND")
	}
	return arena.Make(p.ar, ast.WhileNode{Cond: cond, Body: body, Line: line, Col: col}), nil
}

func (p *Parser) parseFor() (ast.Stmt, error) {
	line, col := p.cur().Line, p.cur().Col
	p.advance()
	vname, err := p.expectIdent()
	if err != nil {
		return nil, err
	}
	if err := p.expect(token.EQ); err != nil {
		return nil, err
	}
	if p.cur().Type == token.EACH {
		p.advance()
		if err := p.expect(token.LPAREN); err != nil {
			return nil, err
		}
		tname, err := p.expectIdent()
		if err != nil {
			return nil, err
		}
		if err := p.expect(token.RPAREN); err != nil {
			return nil, err
		}
		if p.cur().Type != token.NEWLINE {
			return nil, p.failf("expected newline after EACH")
		}
		p.advance()
		p.defineAssignedName(vname)
		body, err := p.parseStmtBlockUntil([]token.TokenType{token.NEXT})
		if err != nil {
			return nil, err
		}
		if err := p.expect(token.NEXT); err != nil {
			return nil, err
		}
		return arena.Make(p.ar, ast.EachStmt{Var: vname, TypeName: tname, Body: body, Line: line, Col: col}), nil
	}
	from, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	if p.cur().Type != token.TO && p.cur().Type != token.DOWNTO {
		return nil, p.failf("expected TO or DOWNTO")
	}
	p.advance()
	to, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	var step ast.Expr
	if p.cur().Type == token.STEP {
		p.advance()
		step, err = p.parseExpr()
		if err != nil {
			return nil, err
		}
	}
	if p.cur().Type != token.NEWLINE {
		return nil, p.failf("expected newline after FOR header")
	}
	p.advance()
	p.defineAssignedName(vname)
	body, err := p.parseStmtBlockUntil([]token.TokenType{token.NEXT})
	if err != nil {
		return nil, err
	}
	if err := p.expect(token.NEXT); err != nil {
		return nil, err
	}
	// Optional NEXT var only on the same line as NEXT (do not skip newlines).
	if p.cur().Type == token.IDENT {
		p.advance()
	}
	return arena.Make(p.ar, ast.ForNode{Var: vname, From: from, To: to, Step: step, Body: body, Line: line, Col: col}), nil
}

func (p *Parser) parseRepeat() (ast.Stmt, error) {
	line, col := p.cur().Line, p.cur().Col
	p.advance()
	if p.cur().Type != token.NEWLINE {
		return nil, p.failf("expected newline after REPEAT")
	}
	p.advance()
	body, err := p.parseStmtBlockUntil([]token.TokenType{token.UNTIL})
	if err != nil {
		return nil, err
	}
	if err := p.expect(token.UNTIL); err != nil {
		return nil, err
	}
	cond, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	return arena.Make(p.ar, ast.RepeatNode{Body: body, Condition: cond, Line: line, Col: col}), nil
}

func (p *Parser) parseDo() (ast.Stmt, error) {
	line, col := p.cur().Line, p.cur().Col
	p.advance() // DO
	p.skipNewlines()
	if p.cur().Type == token.WHILE {
		p.advance()
		cond, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		if p.cur().Type != token.NEWLINE {
			return nil, p.failf("expected newline after DO WHILE condition")
		}
		p.advance()
		body, err := p.parseStmtBlockUntil([]token.TokenType{token.LOOP})
		if err != nil {
			return nil, err
		}
		if err := p.expect(token.LOOP); err != nil {
			return nil, err
		}
		return arena.Make(p.ar, ast.DoLoopNode{Kind: ast.DoPreWhile, Cond: cond, Body: body, Line: line, Col: col}), nil
	}
	p.skipNewlines()
	body, err := p.parseStmtBlockUntil([]token.TokenType{token.LOOP})
	if err != nil {
		return nil, err
	}
	if err := p.expect(token.LOOP); err != nil {
		return nil, err
	}
	p.skipNewlines()
	switch p.cur().Type {
	case token.WHILE:
		p.advance()
		cond, err2 := p.parseExpr()
		if err2 != nil {
			return nil, err2
		}
		return arena.Make(p.ar, ast.DoLoopNode{Kind: ast.DoPostWhile, Cond: cond, Body: body, Line: line, Col: col}), nil
	case token.UNTIL:
		p.advance()
		cond, err2 := p.parseExpr()
		if err2 != nil {
			return nil, err2
		}
		return arena.Make(p.ar, ast.DoLoopNode{Kind: ast.DoPostUntil, Cond: cond, Body: body, Line: line, Col: col}), nil
	default:
		return nil, p.failf("expected WHILE or UNTIL after LOOP")
	}
}

func (p *Parser) parseExit() (ast.Stmt, error) {
	line, col := p.cur().Line, p.cur().Col
	p.advance()
	p.skipNewlines()
	var target string
	switch p.cur().Type {
	case token.FOR:
		p.advance()
		target = "FOR"
	case token.WHILE:
		p.advance()
		target = "WHILE"
	case token.REPEAT:
		p.advance()
		target = "REPEAT"
	case token.DO:
		p.advance()
		target = "DO"
	case token.FUNCTION:
		p.advance()
		target = "FUNCTION"
	default:
		return nil, p.failf("EXIT must be followed by FOR, WHILE, REPEAT, DO, or FUNCTION")
	}
	return arena.Make(p.ar, ast.ExitStmt{Target: target, Line: line, Col: col}), nil
}

func (p *Parser) parseContinue() (ast.Stmt, error) {
	line, col := p.cur().Line, p.cur().Col
	p.advance()
	p.skipNewlines()
	var target string
	switch p.cur().Type {
	case token.FOR:
		p.advance()
		target = "FOR"
	case token.WHILE:
		p.advance()
		target = "WHILE"
	case token.REPEAT:
		p.advance()
		target = "REPEAT"
	case token.DO:
		p.advance()
		target = "DO"
	default:
		return nil, p.failf("CONTINUE must be followed by FOR, WHILE, REPEAT, or DO")
	}
	return arena.Make(p.ar, ast.ContinueStmt{Target: target, Line: line, Col: col}), nil
}
