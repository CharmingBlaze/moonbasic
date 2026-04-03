package lexer

import (
	"strconv"
	"strings"

	"moonbasic/compiler/errors"
	"moonbasic/compiler/token"
)

func (l *Lexer) lexNumber() (token.Token, error) {
	startLine, startCol := l.line, l.col
	var b strings.Builder
	neg := false
	if l.peek() == '-' {
		neg = true
		l.advance()
	}
	if !isDigit(l.peek()) {
		return token.Token{}, errors.NewLexerError(l.file, startLine, startCol,
			"invalid number",
			l.currentLineText(),
			"Expected digits after '-'.")
	}
	for !l.eof() && isDigit(l.peek()) {
		b.WriteByte(l.advance())
	}
	isFloat := false
	if l.peek() == '.' && l.pos+1 < len(l.input) && isDigit(l.input[l.pos+1]) {
		isFloat = true
		b.WriteByte(l.advance()) // '.'
		for !l.eof() && isDigit(l.peek()) {
			b.WriteByte(l.advance())
		}
	}
	lit := b.String()
	if neg {
		lit = "-" + lit
	}
	if isFloat {
		f, err := strconv.ParseFloat(lit, 64)
		if err != nil {
			return token.Token{}, errors.NewLexerError(l.file, startLine, startCol,
				"invalid float literal",
				l.currentLineText(),
				err.Error())
		}
		_ = f
		return token.Token{Type: token.FLOAT, Lit: lit, Line: startLine, Col: startCol}, nil
	}
	i, err := strconv.ParseInt(lit, 10, 64)
	if err != nil {
		return token.Token{}, errors.NewLexerError(l.file, startLine, startCol,
			"invalid integer literal",
			l.currentLineText(),
			err.Error())
	}
	_ = i
	return token.Token{Type: token.INT, Lit: lit, Line: startLine, Col: startCol}, nil
}
