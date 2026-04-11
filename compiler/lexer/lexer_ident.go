package lexer

import (
	"strings"

	"moonbasic/compiler/token"
)

type lexerState struct {
	pos       int
	line      int
	col       int
	lineStart int
}

func (l *Lexer) save() lexerState {
	return lexerState{l.pos, l.line, l.col, l.lineStart}
}

func (l *Lexer) restore(s lexerState) {
	l.pos, l.line, l.col, l.lineStart = s.pos, s.line, s.col, s.lineStart
}

func (l *Lexer) skipSpacesTabs() {
	for !l.eof() {
		c := l.peek()
		if c != ' ' && c != '\t' {
			break
		}
		l.advance()
	}
}

func isIdentStart(c byte) bool {
	return c == '_' || (c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z')
}

func isIdentCont(c byte) bool {
	return isIdentStart(c) || isDigit(c)
}

func (l *Lexer) lexIdent() (token.Token, error) {
	startLine, startCol := l.line, l.col
	var b strings.Builder
	for !l.eof() && isIdentCont(l.peek()) {
		b.WriteByte(l.advance())
	}

	raw := b.String()
	upper := strings.ToUpper(raw)
	kw := token.LookupKeyword(upper)
	if upper == "END" && kw == token.END {
		return l.expandEndKeyword(startLine, startCol)
	}
	tt := token.LookupKeyword(upper)
	canonical := internLit(l, upper)
	if tt == token.IDENT {
		return token.Token{Type: token.IDENT, Lit: canonical, Line: startLine, Col: startCol}, nil
	}
	return token.Token{Type: tt, Lit: canonical, Line: startLine, Col: startCol}, nil
}

// scanIdentUpper reads an identifier (with suffix) without END expansion; returns uppercase lit.
func (l *Lexer) scanIdentUpper() string {
	var b strings.Builder
	for !l.eof() && isIdentCont(l.peek()) {
		b.WriteByte(l.advance())
	}

	return strings.ToUpper(b.String())
}

func internLit(l *Lexer, s string) string {
	if l.names == nil {
		return s
	}
	return l.names.Intern(s)
}

func (l *Lexer) expandEndKeyword(startLine, startCol int) (token.Token, error) {
	saved := l.save()
	l.skipSpacesTabs()
	if l.eof() || l.peek() == '\n' {
		l.restore(saved)
		return token.Token{Type: token.END, Lit: internLit(l, "END"), Line: startLine, Col: startCol}, nil
	}
	if !isIdentStart(l.peek()) {
		l.restore(saved)
		return token.Token{Type: token.END, Lit: internLit(l, "END"), Line: startLine, Col: startCol}, nil
	}
	next := l.scanIdentUpper()
	switch next {
	case "IF":
		return token.Token{Type: token.ENDIF, Lit: internLit(l, "ENDIF"), Line: startLine, Col: startCol}, nil
	case "FUNCTION":
		return token.Token{Type: token.ENDFUNCTION, Lit: internLit(l, "ENDFUNCTION"), Line: startLine, Col: startCol}, nil
	case "WHILE":
		return token.Token{Type: token.WEND, Lit: internLit(l, "WEND"), Line: startLine, Col: startCol}, nil
	case "SELECT":
		return token.Token{Type: token.ENDSELECT, Lit: internLit(l, "ENDSELECT"), Line: startLine, Col: startCol}, nil
	case "TYPE":
		return token.Token{Type: token.ENDTYPE, Lit: internLit(l, "ENDTYPE"), Line: startLine, Col: startCol}, nil
	default:
		l.restore(saved)
		return token.Token{Type: token.END, Lit: internLit(l, "END"), Line: startLine, Col: startCol}, nil
	}
}

func isDigit(c byte) bool { return c >= '0' && c <= '9' }
