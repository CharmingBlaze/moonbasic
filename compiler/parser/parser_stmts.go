package parser

import (
	"moonbasic/compiler/arena"
	"moonbasic/compiler/ast"
	"moonbasic/compiler/token"
)

func containsStop(stop []token.TokenType, t token.TokenType) bool {
	for _, s := range stop {
		if s == t {
			return true
		}
	}
	return false
}

func (p *Parser) parseStmt() (ast.Stmt, error) {
	p.skipNewlines()
	switch p.cur().Type {
	case token.NEWLINE:
		p.advance()
		return nil, nil
	case token.IF:
		return p.parseIf()
	case token.WHILE:
		return p.parseWhile()
	case token.FOR:
		return p.parseFor()
	case token.REPEAT:
		return p.parseRepeat()
	case token.DO:
		return p.parseDo()
	case token.EXIT:
		return p.parseExit()
	case token.CONTINUE:
		return p.parseContinue()
	case token.SELECT:
		return p.parseSelect()
	case token.DIM, token.REDIM:
		return p.parseDim()
	case token.LOCAL:
		return p.parseLocal()
	case token.GLOBAL:
		return p.parseGlobal()
	case token.CONST:
		return p.parseConst()
	case token.STATIC:
		return p.parseStatic()
	case token.SWAP:
		return p.parseSwap()
	case token.ERASE:
		return p.parseErase()
	case token.GOTO:
		return p.parseGoto()
	case token.GOSUB:
		return p.parseGosub()
	case token.RETURN:
		return p.parseReturn()
	case token.INCLUDE:
		return p.parseInclude()
	case token.DELETE:
		line, col := p.cur().Line, p.cur().Col
		p.advance()
		e, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		return arena.Make(p.ar, ast.DeleteStmt{Expr: e, Line: line, Col: col}), nil
	case token.END:
		line, col := p.cur().Line, p.cur().Col
		p.advance()
		return arena.Make(p.ar, ast.EndProgramStmt{Line: line, Col: col}), nil
	case token.DOT:
		line, col := p.cur().Line, p.cur().Col
		p.advance()
		name, err := p.expectIdent()
		if err != nil {
			return nil, err
		}
		return arena.Make(p.ar, ast.LabelNode{Name: name, Line: line, Col: col}), nil
	case token.IDENT:
		name := p.cur().Lit
		line, col := p.cur().Line, p.cur().Col
		p.advance()
		return p.parseStmtAfterIdent(name, line, col)
	default:
		return nil, p.failf("unexpected token at start of statement: %s", p.cur().Type.String())
	}
}

func (p *Parser) parseStmtAfterIdent(name string, line, col int) (ast.Stmt, error) {
	p.skipNewlines()
	if p.cur().Type == token.DOT {
		p.advance()
		field, err := p.expectIdent()
		if err != nil {
			return nil, err
		}
		p.skipNewlines()
		switch p.cur().Type {
		case token.EQ:
			p.advance()
			e, err := p.parseExpr()
			if err != nil {
				return nil, err
			}
			return arena.Make(p.ar, ast.FieldAssignNode{Object: name, Field: field, Expr: e, Line: line, Col: col}), nil
		case token.LPAREN:
			args, err := p.parseArgList()
			if err != nil {
				return nil, err
			}
			if p.sym.IsVar(name) {
				return arena.Make(p.ar, ast.HandleCallStmt{Receiver: name, Method: field, Args: args, Line: line, Col: col}), nil
			}
			return arena.Make(p.ar, ast.NamespaceCallStmt{NS: name, Method: field, Args: args, Line: line, Col: col}), nil
		default:
			return nil, p.failf("expected '=' or '(' after field %s", field)
		}
	}
	if p.cur().Type == token.LPAREN {
		if p.sym.IsVar(name) && !p.sym.IsFunction(name) {
			args, err := p.parseArgList()
			if err != nil {
				return nil, err
			}
			p.skipNewlines()
			switch p.cur().Type {
			case token.DOT:
				p.advance()
				field, errDot := p.expectIdent()
				if errDot != nil {
					return nil, errDot
				}
				p.skipNewlines()
				switch p.cur().Type {
				case token.EQ:
					p.advance()
					rhs, err2 := p.parseExpr()
					if err2 != nil {
						return nil, err2
					}
					return arena.Make(p.ar, ast.IndexFieldAssignNode{Array: name, Index: args, Field: field, Expr: rhs, Line: line, Col: col}), nil
				case token.PLUSEQ, token.MINUSEQ, token.STAREQ, token.SLASHEQ:
					op := p.cur().Type
					p.advance()
					rhs, err2 := p.parseExpr()
					if err2 != nil {
						return nil, err2
					}
					var binOp string
					switch op {
					case token.PLUSEQ:
						binOp = "+"
					case token.MINUSEQ:
						binOp = "-"
					case token.STAREQ:
						binOp = "*"
					case token.SLASHEQ:
						binOp = "/"
					}
					load := arena.Make(p.ar, ast.IndexFieldExpr{Array: name, Index: args, Field: field, Line: line, Col: col})
					bin := arena.Make(p.ar, ast.BinopNode{Op: binOp, Left: load, Right: rhs, Line: line, Col: col})
					return arena.Make(p.ar, ast.IndexFieldAssignNode{Array: name, Index: args, Field: field, Expr: bin, Line: line, Col: col}), nil
				default:
					return nil, p.failf("expected '=' or compound assign after %s(...).%s", name, field)
				}
			case token.EQ:
				p.advance()
				rhs, err2 := p.parseExpr()
				if err2 != nil {
					return nil, err2
				}
				return arena.Make(p.ar, ast.IndexAssignNode{Array: name, Index: args, Expr: rhs, Line: line, Col: col}), nil
			case token.PLUSEQ, token.MINUSEQ, token.STAREQ, token.SLASHEQ:
				op := p.cur().Type
				p.advance()
				rhs, err2 := p.parseExpr()
				if err2 != nil {
					return nil, err2
				}
				var binOp string
				switch op {
				case token.PLUSEQ:
					binOp = "+"
				case token.MINUSEQ:
					binOp = "-"
				case token.STAREQ:
					binOp = "*"
				case token.SLASHEQ:
					binOp = "/"
				}
				idn := arena.Make(p.ar, ast.IdentNode{Name: name, Line: line, Col: col})
				load := arena.Make(p.ar, ast.IndexExpr{Base: idn, Index: args, Line: line, Col: col})
				bin := arena.Make(p.ar, ast.BinopNode{Op: binOp, Left: load, Right: rhs, Line: line, Col: col})
				return arena.Make(p.ar, ast.IndexAssignNode{Array: name, Index: args, Expr: bin, Line: line, Col: col}), nil
			default:
				return nil, p.failf("expected '=' or compound assign after subscript")
			}
		}
		args, err := p.parseArgList()
		if err != nil {
			return nil, err
		}
		return arena.Make(p.ar, ast.CallStmtNode{Name: name, Args: args, Line: line, Col: col}), nil
	}
	if p.cur().Type == token.EQ {
		p.advance()
		e, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		p.defineAssignedName(name)
		return arena.Make(p.ar, ast.AssignNode{Name: name, Expr: e, Line: line, Col: col}), nil
	}
	if p.cur().Type == token.PLUSEQ || p.cur().Type == token.MINUSEQ || p.cur().Type == token.STAREQ || p.cur().Type == token.SLASHEQ {
		op := p.cur().Type
		p.advance()
		rhs, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		p.defineAssignedName(name)
		var binOp string
		switch op {
		case token.PLUSEQ:
			binOp = "+"
		case token.MINUSEQ:
			binOp = "-"
		case token.STAREQ:
			binOp = "*"
		case token.SLASHEQ:
			binOp = "/"
		}
		left := arena.Make(p.ar, ast.IdentNode{Name: name, Line: line, Col: col})
		bin := arena.Make(p.ar, ast.BinopNode{Op: binOp, Left: left, Right: rhs, Line: line, Col: col})
		return arena.Make(p.ar, ast.AssignNode{Name: name, Expr: bin, Line: line, Col: col}), nil
	}
	return nil, p.failf("expected '=', compound assign, '.' or '(' after %s", name)
}

func (p *Parser) parseIf() (ast.Stmt, error) {
	line, col := p.cur().Line, p.cur().Col
	p.advance()
	cond, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	if err := p.expect(token.THEN); err != nil {
		return nil, err
	}
	var then []ast.Stmt
	multiline := false
	if p.cur().Type == token.NEWLINE {
		multiline = true
		p.advance()
		then, err = p.parseStmtBlockUntil([]token.TokenType{token.ELSEIF, token.ELSE, token.ENDIF})
		if err != nil {
			return nil, err
		}
	} else {
		then, err = p.parseSingleLineIfBody()
		if err != nil {
			return nil, err
		}
	}
	var elseif []ast.ElseIfClause
	var els []ast.Stmt
	if multiline {
		for p.cur().Type == token.ELSEIF {
			p.advance()
			ec, err2 := p.parseExpr()
			if err2 != nil {
				return nil, err2
			}
			if err2 := p.expect(token.THEN); err2 != nil {
				return nil, err2
			}
			var body []ast.Stmt
			if p.cur().Type == token.NEWLINE {
				p.advance()
				body, err2 = p.parseStmtBlockUntil([]token.TokenType{token.ELSEIF, token.ELSE, token.ENDIF})
			} else {
				body, err2 = p.parseSingleLineIfBody()
			}
			if err2 != nil {
				return nil, err2
			}
			elseif = append(elseif, ast.ElseIfClause{Cond: ec, Body: body})
		}
		if p.cur().Type == token.ELSE {
			p.advance()
			if p.cur().Type == token.NEWLINE {
				p.advance()
				els, err = p.parseStmtBlockUntil([]token.TokenType{token.ENDIF})
			} else {
				els, err = p.parseSingleLineIfBody()
			}
			if err != nil {
				return nil, err
			}
		}
		if err := p.expect(token.ENDIF); err != nil {
			return nil, err
		}
	}
	return arena.Make(p.ar, ast.IfNode{Cond: cond, Then: then, ElseIf: elseif, Else: els, Line: line, Col: col}), nil
}

// parseSingleLineIfBody reads statements until newline (optional colon between); consumes final newline.
func (p *Parser) parseSingleLineIfBody() ([]ast.Stmt, error) {
	var out []ast.Stmt
	for {
		if p.cur().Type == token.NEWLINE || p.cur().Type == token.EOF {
			break
		}
		s, err := p.parseStmt()
		if err != nil {
			return nil, err
		}
		if s != nil {
			out = append(out, s)
		}
		if p.cur().Type == token.COLON {
			p.advance()
			continue
		}
		break
	}
	if p.cur().Type == token.NEWLINE {
		p.advance()
	}
	return out, nil
}

func (p *Parser) parseStmtBlockUntil(stop []token.TokenType) ([]ast.Stmt, error) {
	var out []ast.Stmt
	for {
		p.skipNewlines()
		if containsStop(stop, p.cur().Type) || p.cur().Type == token.EOF {
			break
		}
		s, err := p.parseStmt()
		if err != nil {
			return nil, err
		}
		if s != nil {
			out = append(out, s)
		}
	}
	return out, nil
}

func (p *Parser) parseSelect() (ast.Stmt, error) {
	line, col := p.cur().Line, p.cur().Col
	p.advance()
	expr, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	if p.cur().Type != token.NEWLINE {
		return nil, p.failf("expected newline after SELECT expr")
	}
	p.advance()
	var cases []ast.CaseClause
	var def []ast.Stmt
	for {
		p.skipNewlines()
		switch p.cur().Type {
		case token.CASE:
			p.advance()
			val, err := p.parseExpr()
			if err != nil {
				return nil, err
			}
			if p.cur().Type != token.NEWLINE {
				return nil, p.failf("expected newline after CASE value")
			}
			p.advance()
			body, err := p.parseStmtBlockUntil([]token.TokenType{token.CASE, token.DEFAULT, token.ENDSELECT})
			if err != nil {
				return nil, err
			}
			cases = append(cases, ast.CaseClause{Value: val, Body: body})
		case token.DEFAULT:
			p.advance()
			if p.cur().Type == token.NEWLINE {
				p.advance()
			}
			def, err = p.parseStmtBlockUntil([]token.TokenType{token.ENDSELECT})
			if err != nil {
				return nil, err
			}
		case token.ENDSELECT:
			p.advance()
			return arena.Make(p.ar, ast.SelectNode{Expr: expr, Cases: cases, Default: def, Line: line, Col: col}), nil
		default:
			return nil, p.failf("expected CASE, DEFAULT, or ENDSELECT")
		}
	}
}

func (p *Parser) parseDim() (ast.Stmt, error) {
	line, col := p.cur().Line, p.cur().Col
	isRedim := p.cur().Type == token.REDIM
	p.advance()
	preserve := true
	if isRedim && p.cur().Type == token.PRESERVE {
		p.advance()
		preserve = true
	}
	name, err := p.expectIdent()
	if err != nil {
		return nil, err
	}
	elemType := ""
	if !isRedim && p.cur().Type == token.AS {
		p.advance()
		tn, err2 := p.expectIdent()
		if err2 != nil {
			return nil, err2
		}
		elemType = tn
	}
	args, err := p.parseArgList()
	if err != nil {
		return nil, err
	}
	p.defineAssignedName(name)
	return arena.Make(p.ar, ast.DimNode{
		Name: name, ElemType: elemType, Dims: args, Line: line, Col: col,
		IsRedim: isRedim, Preserve: preserve,
	}), nil
}

func (p *Parser) parseStatic() (ast.Stmt, error) {
	line, col := p.cur().Line, p.cur().Col
	p.advance()
	if p.FuncName == "" {
		return nil, p.failf("STATIC may only be used inside a FUNCTION")
	}
	name, err := p.expectIdent()
	if err != nil {
		return nil, err
	}
	p.sym.DefineStatic(p.FuncName, name)
	var init ast.Expr
	if p.cur().Type == token.EQ {
		p.advance()
		init, err = p.parseExpr()
		if err != nil {
			return nil, err
		}
	}
	return arena.Make(p.ar, ast.StaticDeclNode{Name: name, Init: init, Line: line, Col: col}), nil
}

func (p *Parser) parseSwap() (ast.Stmt, error) {
	line, col := p.cur().Line, p.cur().Col
	p.advance()
	paren := false
	if p.cur().Type == token.LPAREN {
		paren = true
		p.advance()
	}
	a, err := p.expectIdent()
	if err != nil {
		return nil, err
	}
	if err := p.expect(token.COMMA); err != nil {
		return nil, err
	}
	b, err := p.expectIdent()
	if err != nil {
		return nil, err
	}
	if paren {
		if err := p.expect(token.RPAREN); err != nil {
			return nil, err
		}
	}
	return arena.Make(p.ar, ast.SwapStmt{A: a, B: b, Line: line, Col: col}), nil
}

func (p *Parser) parseErase() (ast.Stmt, error) {
	line, col := p.cur().Line, p.cur().Col
	p.advance()
	paren := false
	if p.cur().Type == token.LPAREN {
		paren = true
		p.advance()
	}
	name, err := p.expectIdent()
	if err != nil {
		return nil, err
	}
	if paren {
		if err := p.expect(token.RPAREN); err != nil {
			return nil, err
		}
	}
	return arena.Make(p.ar, ast.EraseStmt{Name: name, Line: line, Col: col}), nil
}

func (p *Parser) parseLocal() (ast.Stmt, error) {
	line, col := p.cur().Line, p.cur().Col
	p.advance()
	name, err := p.expectIdent()
	if err != nil {
		return nil, err
	}
	p.sym.DefineLocal(name)
	var init ast.Expr
	if p.cur().Type == token.EQ {
		p.advance()
		init, err = p.parseExpr()
		if err != nil {
			return nil, err
		}
	}
	return arena.Make(p.ar, ast.LocalDeclNode{Name: name, Init: init, Line: line, Col: col}), nil
}

func (p *Parser) parseGlobal() (ast.Stmt, error) {
	line, col := p.cur().Line, p.cur().Col
	p.advance()
	name, err := p.expectIdent()
	if err != nil {
		return nil, err
	}
	p.sym.DefineGlobalVar(name)
	if p.cur().Type != token.EQ {
		return nil, p.failf("GLOBAL requires '=' initializer")
	}
	p.advance()
	e, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	return arena.Make(p.ar, ast.AssignNode{Name: name, Expr: e, Global: true, Line: line, Col: col}), nil
}

func (p *Parser) parseConst() (ast.Stmt, error) {
	line, col := p.cur().Line, p.cur().Col
	p.advance()
	name, err := p.expectIdent()
	if err != nil {
		return nil, err
	}
	if err := p.expect(token.EQ); err != nil {
		return nil, err
	}
	e, err := p.parseExpr()
	if err != nil {
		return nil, err
	}
	p.sym.DefineConst(name)
	return arena.Make(p.ar, ast.ConstDeclNode{Name: name, Expr: e, Line: line, Col: col}), nil
}

func (p *Parser) parseGoto() (ast.Stmt, error) {
	line, col := p.cur().Line, p.cur().Col
	p.advance()
	lbl, err := p.expectIdent()
	if err != nil {
		return nil, err
	}
	return arena.Make(p.ar, ast.GotoNode{Label: lbl, Line: line, Col: col}), nil
}

func (p *Parser) parseGosub() (ast.Stmt, error) {
	line, col := p.cur().Line, p.cur().Col
	p.advance()
	lbl, err := p.expectIdent()
	if err != nil {
		return nil, err
	}
	return arena.Make(p.ar, ast.GosubNode{Label: lbl, Line: line, Col: col}), nil
}

func (p *Parser) parseReturn() (ast.Stmt, error) {
	line, col := p.cur().Line, p.cur().Col
	p.advance()
	if p.cur().Type == token.LPAREN {
		p.advance()
		if p.cur().Type == token.RPAREN {
			p.advance()
			return arena.Make(p.ar, ast.ReturnNode{Expr: nil, Line: line, Col: col}), nil
		}
		e, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		if err := p.expect(token.RPAREN); err != nil {
			return nil, err
		}
		return arena.Make(p.ar, ast.ReturnNode{Expr: e, Line: line, Col: col}), nil
	}
	var e ast.Expr
	var err error
	if p.cur().Type != token.NEWLINE && p.cur().Type != token.EOF && p.cur().Type != token.ENDIF && p.cur().Type != token.ENDFUNCTION {
		e, err = p.parseExpr()
		if err != nil {
			return nil, err
		}
	}
	return arena.Make(p.ar, ast.ReturnNode{Expr: e, Line: line, Col: col}), nil
}

func (p *Parser) parseInclude() (ast.Stmt, error) {
	line, col := p.cur().Line, p.cur().Col
	p.advance()
	p.skipNewlines()
	t := p.cur()
	if t.Type != token.STRING {
		return nil, p.failf("INCLUDE requires string path")
	}
	p.advance()
	return arena.Make(p.ar, ast.IncludeNode{Path: t.Lit, Line: line, Col: col}), nil
}
