// Expression precedence (loose → tight): OR < XOR < AND < NOT < comparisons < +- < */MOD < ^ < unary- < postfix < primary.
// Canonical table: ARCHITECTURE.md §7; long-form: compiler/errors/MoonBasic.md (Expression precedence).
package parser

import (
	"moonbasic/compiler/arena"
	"moonbasic/compiler/ast"
	"moonbasic/compiler/token"
)

func (p *Parser) parseExpr() (ast.Expr, error) {
	return p.parseOr()
}

func (p *Parser) parseOr() (ast.Expr, error) {
	left, err := p.parseXor()
	if err != nil {
		return nil, err
	}
	for p.cur().Type == token.OR {
		line, col := p.cur().Line, p.cur().Col
		p.advance()
		right, err := p.parseXor()
		if err != nil {
			return nil, err
		}
		left = arena.Make(p.ar, ast.BinopNode{Op: "OR", Left: left, Right: right, Line: line, Col: col})
	}
	return left, nil
}

func (p *Parser) parseXor() (ast.Expr, error) {
	left, err := p.parseAnd()
	if err != nil {
		return nil, err
	}
	for p.cur().Type == token.XOR {
		line, col := p.cur().Line, p.cur().Col
		p.advance()
		right, err := p.parseAnd()
		if err != nil {
			return nil, err
		}
		left = arena.Make(p.ar, ast.BinopNode{Op: "XOR", Left: left, Right: right, Line: line, Col: col})
	}
	return left, nil
}

func (p *Parser) parseAnd() (ast.Expr, error) {
	left, err := p.parseNot()
	if err != nil {
		return nil, err
	}
	for p.cur().Type == token.AND {
		line, col := p.cur().Line, p.cur().Col
		p.advance()
		right, err := p.parseNot()
		if err != nil {
			return nil, err
		}
		left = arena.Make(p.ar, ast.BinopNode{Op: "AND", Left: left, Right: right, Line: line, Col: col})
	}
	return left, nil
}

func (p *Parser) parseNot() (ast.Expr, error) {
	if p.cur().Type == token.NOT {
		line, col := p.cur().Line, p.cur().Col
		p.advance()
		inner, err := p.parseNot()
		if err != nil {
			return nil, err
		}
		return arena.Make(p.ar, ast.UnaryNode{Op: "NOT", Expr: inner, Line: line, Col: col}), nil
	}
	return p.parseComparison()
}

func (p *Parser) parseComparison() (ast.Expr, error) {
	left, err := p.parseAddSub()
	if err != nil {
		return nil, err
	}
	for {
		var op string
		switch p.cur().Type {
		case token.EQ:
			op = "="
		case token.NEQ:
			op = "<>"
		case token.LT:
			op = "<"
		case token.GT:
			op = ">"
		case token.LTE:
			op = "<="
		case token.GTE:
			op = ">="
		default:
			return left, nil
		}
		line, col := p.cur().Line, p.cur().Col
		p.advance()
		right, err := p.parseAddSub()
		if err != nil {
			return nil, err
		}
		left = arena.Make(p.ar, ast.BinopNode{Op: op, Left: left, Right: right, Line: line, Col: col})
	}
}

func (p *Parser) parseAddSub() (ast.Expr, error) {
	left, err := p.parseMulDivMod()
	if err != nil {
		return nil, err
	}
	for {
		var op string
		switch p.cur().Type {
		case token.PLUS:
			op = "+"
		case token.MINUS:
			op = "-"
		default:
			return left, nil
		}
		line, col := p.cur().Line, p.cur().Col
		p.advance()
		right, err := p.parseMulDivMod()
		if err != nil {
			return nil, err
		}
		left = arena.Make(p.ar, ast.BinopNode{Op: op, Left: left, Right: right, Line: line, Col: col})
	}
}

func (p *Parser) parseMulDivMod() (ast.Expr, error) {
	left, err := p.parsePow()
	if err != nil {
		return nil, err
	}
	for {
		var op string
		switch p.cur().Type {
		case token.STAR:
			op = "*"
		case token.SLASH:
			op = "/"
		case token.MOD:
			line, col := p.cur().Line, p.cur().Col
			p.advance()
			right, err := p.parsePow()
			if err != nil {
				return nil, err
			}
			left = arena.Make(p.ar, ast.BinopNode{Op: "MOD", Left: left, Right: right, Line: line, Col: col})
			continue
		default:
			return left, nil
		}
		line, col := p.cur().Line, p.cur().Col
		p.advance()
		right, err := p.parsePow()
		if err != nil {
			return nil, err
		}
		left = arena.Make(p.ar, ast.BinopNode{Op: op, Left: left, Right: right, Line: line, Col: col})
	}
}

func (p *Parser) parsePow() (ast.Expr, error) {
	left, err := p.parseUnaryMinus()
	if err != nil {
		return nil, err
	}
	if p.cur().Type == token.CARET {
		line, col := p.cur().Line, p.cur().Col
		p.advance()
		right, err := p.parsePow()
		if err != nil {
			return nil, err
		}
		left = arena.Make(p.ar, ast.BinopNode{Op: "^", Left: left, Right: right, Line: line, Col: col})
	}
	return left, nil
}

func (p *Parser) parseUnaryMinus() (ast.Expr, error) {
	if p.cur().Type == token.MINUS {
		line, col := p.cur().Line, p.cur().Col
		p.advance()
		inner, err := p.parseUnaryMinus()
		if err != nil {
			return nil, err
		}
		return arena.Make(p.ar, ast.UnaryNode{Op: "-", Expr: inner, Line: line, Col: col}), nil
	}
	return p.parsePostfix()
}

func (p *Parser) parsePostfix() (ast.Expr, error) {
	base, err := p.parsePrimary()
	if err != nil {
		return nil, err
	}
	return p.parsePostfixChain(base)
}

func (p *Parser) parsePostfixChain(base ast.Expr) (ast.Expr, error) {
	for {
		switch p.cur().Type {
		case token.LPAREN:
			openBracket := p.cur().Lit == "["
			args, err := p.parseArgList()
			if err != nil {
				return nil, err
			}
			if id, ok := base.(*ast.IdentNode); ok {
				if openBracket || (p.sym.IsVar(id.Name) && !p.sym.IsFunction(id.Name)) {
					base = arena.Make(p.ar, ast.IndexExpr{Base: id, Index: args, Line: id.Line, Col: id.Col})
				} else {
					base = arena.Make(p.ar, ast.CallExprNode{Name: id.Name, Args: args, Line: id.Line, Col: id.Col})
				}
			} else {
				base = arena.Make(p.ar, ast.IndexExpr{Base: base, Index: args, Line: p.cur().Line, Col: p.cur().Col})
			}
		case token.DOT:
			id, ok := base.(*ast.IdentNode)
			if !ok {
				return nil, p.failf("method call requires simple identifier receiver")
			}
			p.advance()
			meth, err := p.expectIdent()
			if err != nil {
				return nil, err
			}
			if p.cur().Type == token.LPAREN {
				args, err := p.parseArgList()
				if err != nil {
					return nil, err
				}
				if p.sym.IsVar(id.Name) {
					base = arena.Make(p.ar, ast.HandleCallExpr{Receiver: id.Name, Method: meth, Args: args, Line: id.Line, Col: id.Col})
				} else {
					base = arena.Make(p.ar, ast.NamespaceCallExpr{NS: id.Name, Method: meth, Args: args, Line: id.Line, Col: id.Col})
				}
			} else {
				// No parens = Field Access (e.g. p.name)
				base = arena.Make(p.ar, ast.FieldAccessNode{Object: id.Name, Field: meth, Line: id.Line, Col: id.Col})
			}
		default:
			return base, nil
		}
	}
}

func (p *Parser) parsePrimary() (ast.Expr, error) {
	p.skipNewlines()
	switch t := p.cur(); t.Type {
	case token.INT:
		p.advance()
		ex, err := p.parseNumberLit(t.Lit, t.Line, t.Col)
		return ex, err
	case token.FLOAT:
		p.advance()
		f, err := p.parseNumberLit(t.Lit, t.Line, t.Col)
		return f, err
	case token.STRING:
		p.advance()
		return arena.Make(p.ar, ast.StringLitNode{Value: t.Lit, Line: t.Line, Col: t.Col}), nil
	case token.TRUE:
		p.advance()
		return arena.Make(p.ar, ast.BoolLitNode{Value: true, Line: t.Line, Col: t.Col}), nil
	case token.FALSE:
		p.advance()
		return arena.Make(p.ar, ast.BoolLitNode{Value: false, Line: t.Line, Col: t.Col}), nil
	case token.NULL:
		p.advance()
		return arena.Make(p.ar, ast.NullLitNode{Line: t.Line, Col: t.Col}), nil
	case token.NEW:
		p.advance()
		if err := p.expect(token.LPAREN); err != nil {
			return nil, err
		}
		name, err := p.expectIdent()
		if err != nil {
			return nil, err
		}
		if err := p.expect(token.RPAREN); err != nil {
			return nil, err
		}
		return arena.Make(p.ar, ast.NewNode{TypeName: name, Line: t.Line, Col: t.Col}), nil
	case token.LPAREN:
		p.advance()
		inner, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		if err := p.expect(token.RPAREN); err != nil {
			return nil, err
		}
		return arena.Make(p.ar, ast.GroupedExpr{Inner: inner}), nil
	case token.IDENT:
		p.advance()
		return arena.Make(p.ar, ast.IdentNode{Name: t.Lit, Line: t.Line, Col: t.Col}), nil
	default:
		return nil, p.failf("unexpected token in expression: %s", t.Type.String())
	}
}

func (p *Parser) parseArgList() ([]ast.Expr, error) {
	if err := p.expect(token.LPAREN); err != nil {
		return nil, err
	}
	var args []ast.Expr
	p.skipNewlines()
	if p.cur().Type == token.RPAREN {
		p.advance()
		return args, nil
	}
	for {
		e, err := p.parseExpr()
		if err != nil {
			return nil, err
		}
		args = append(args, e)
		p.skipNewlines()
		if p.cur().Type == token.COMMA {
			p.advance()
			p.skipNewlines()
			continue
		}
		break
	}
	if err := p.expect(token.RPAREN); err != nil {
		return nil, err
	}
	return args, nil
}
