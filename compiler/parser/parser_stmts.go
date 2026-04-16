package parser

import (
	"moonbasic/compiler/arena"
	"moonbasic/compiler/ast"
	"moonbasic/compiler/builtinmanifest"
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

	// Soft keyword support for namespaced calls: Static.Create(), End.Exit(), etc.
	if isKeywordUsableAsIdent(p.cur().Type) && p.i+1 < len(p.toks) && p.toks[p.i+1].Type == token.DOT {
		name := p.cur().Lit
		line, col := p.cur().Line, p.cur().Col
		p.advance()
		return p.parseStmtAfterIdent(name, line, col)
	}

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
	// Before skipNewlines: bare `DrawEntities` / `DrawEntities()`-style newline must not be eaten,
	// or the newline would be consumed and a following `DrawEntities()` would merge incorrectly.
	if p.isBareZeroArgBuiltinTerminator() && builtinmanifest.Default().HasArityExact(name, 0) {
		return arena.Make(p.ar, ast.CallStmtNode{Name: name, Args: nil, Line: line, Col: col}), nil
	}
	p.skipNewlines()
	// a, b, c = expr  (destructuring assignment)
	if p.cur().Type == token.COMMA {
		names := []string{name}
		for p.cur().Type == token.COMMA {
			p.advance()
			nm, err := p.expectIdent()
			if err != nil {
				return nil, err
			}
			names = append(names, nm)
		}
		if p.cur().Type != token.EQ {
			return nil, p.failf("expected '=' after destructuring targets")
		}
		p.advance()
		e, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		for _, nm := range names {
			p.defineAssignedName(nm)
		}
		return arena.Make(p.ar, ast.MultiAssignNode{Names: names, Expr: e, Line: line, Col: col}), nil
	}
	// name AS type(dim...) — typed array declaration without DIM
	if p.cur().Type == token.AS {
		p.advance()
		typeName, err := p.expectIdent()
		if err != nil {
			return nil, err
		}
		args, err := p.parseArgList()
		if err != nil {
			return nil, err
		}
		p.defineAssignedName(name)
		return arena.Make(p.ar, ast.DimNode{
			Name: name, TypeName: typeName, Dims: args, Line: line, Col: col,
		}), nil
	}
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
			p.skipNewlines()
			if p.cur().Type == token.EQ {
				p.advance()
				rhs, err := p.parseExpr()
				if err != nil {
					return nil, err
				}
				return arena.Make(p.ar, ast.NamespaceAssignNode{NS: name, Method: field, Args: args, Expr: rhs, Line: line, Col: col}), nil
			}
			if p.sym.IsVar(name) {
				return arena.Make(p.ar, ast.HandleCallStmt{Receiver: arena.Make(p.ar, ast.IdentNode{Name: name, Line: line, Col: col}), Method: field, Args: args, Line: line, Col: col}), nil
			}
			return arena.Make(p.ar, ast.NamespaceCallStmt{NS: name, Method: field, Args: args, Line: line, Col: col}), nil
		default:
			return nil, p.failf("expected '=' or '(' after field %s", field)
		}
	}
	if p.cur().Type == token.LPAREN {
		// Parenthesized: could be Array Index Assignment or Function Call
		saved := p.save()
		args, err := p.parseArgList()
		if err != nil {
			return nil, err
		}
		p.skipNewlines()
		// arr(idx...).field = expr
		if p.cur().Type == token.DOT {
			p.advance()
			field, err2 := p.expectIdent()
			if err2 != nil {
				return nil, err2
			}
			p.skipNewlines()
			if p.cur().Type != token.EQ {
				return nil, p.failf("expected '=' after %s(...).%s", name, field)
			}
			p.advance()
			rhs, err3 := p.parseExpr()
			if err3 != nil {
				return nil, err3
			}
			return arena.Make(p.ar, ast.IndexFieldAssignNode{Array: name, Index: args, Field: field, Expr: rhs, Line: line, Col: col}), nil
		}
		if p.cur().Type == token.EQ || p.cur().Type == token.PLUSEQ || p.cur().Type == token.MINUSEQ || p.cur().Type == token.STAREQ || p.cur().Type == token.SLASHEQ {
			// It is an assignment to an array or object method result (if we allowed that, but we don't yet).
			// We treat it as an IndexAssignNode and let the Analyzer flag if 'name' isn't an array.
			p.defineAssignedName(name) // Mark as assigned for the parser's symtable if we want scalar-style implicit

			// Actually, for arrays, we should NOT define it here yet because arrays ALWAYS need DIM.
			// But to parse it correctly:
			if p.cur().Type == token.EQ {
				p.advance()
				rhs, err2 := p.parseExpr()
				if err2 != nil {
					return nil, err2
				}
				return arena.Make(p.ar, ast.IndexAssignNode{Array: name, Index: args, Expr: rhs, Line: line, Col: col}), nil
			}
			// Handle compound op...
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
			load := arena.Make(p.ar, ast.IndexExpr{Base: arena.Make(p.ar, ast.IdentNode{Name: name, Line: line, Col: col}), Index: args, Line: line, Col: col})
			bin := arena.Make(p.ar, ast.BinopNode{Op: binOp, Left: load, Right: rhs, Line: line, Col: col})
			return arena.Make(p.ar, ast.IndexAssignNode{Array: name, Index: args, Expr: bin, Line: line, Col: col}), nil
		}

		// Not an assignment, backtrack and parse as function call
		p.restore(saved)
		args, err = p.parseArgList()
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

// isBareZeroArgBuiltinTerminator is true when the next token cannot continue an assignment/call
// and may end a bare statement (newline, colon, EOF, or block delimiter).
func (p *Parser) isBareZeroArgBuiltinTerminator() bool {
	switch p.cur().Type {
	case token.NEWLINE, token.COLON, token.EOF,
		token.WEND, token.ENDWHILE, token.NEXT, token.UNTIL, token.LOOP,
		token.ENDIF, token.ELSE, token.ELSEIF, token.END,
		token.CASE, token.DEFAULT, token.ENDSELECT:
		return true
	default:
		return false
	}
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
		// Same as ParseProgram: chain statements on one line with colon (e.g. vx = 0 : vz = 0).
		for p.cur().Type == token.COLON {
			p.advance()
			p.skipNewlines()
			s2, err2 := p.parseStmt()
			if err2 != nil {
				return nil, err2
			}
			if s2 != nil {
				out = append(out, s2)
			}
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
	typeName := ""
	if !isRedim && p.cur().Type == token.AS {
		p.advance()
		tn, err2 := p.expectIdent()
		if err2 != nil {
			return nil, err2
		}
		typeName = tn
	}
	args, err := p.parseArgList()
	if err != nil {
		return nil, err
	}
	p.defineAssignedName(name)
	return arena.Make(p.ar, ast.DimNode{
		Name: name, TypeName: typeName, Dims: args, Line: line, Col: col,
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
