package lexer

import (
	"strings"

	"moonbasic/compiler/errors"
	"moonbasic/compiler/token"
)

func (l *Lexer) lexString() (token.Token, error) {
	startLine, startCol := l.line, l.col
	if l.peek() != '"' {
		return token.Token{}, errors.NewLexerError(l.file, l.line, l.col,
			"expected '\"' to start string",
			l.currentLineText(),
			"Start string literals with a double quote.")
	}
	l.advance() // opening "
	var b strings.Builder
	for !l.eof() {
		c := l.peek()
		if c == '"' {
			l.advance()
			return token.Token{Type: token.STRING, Lit: b.String(), Line: startLine, Col: startCol}, nil
		}
		if c == '\\' {
			l.advance()
			if l.eof() {
				break
			}
			switch l.peek() {
			case '"':
				l.advance()
				b.WriteByte('"')
			case '\\':
				l.advance()
				b.WriteByte('\\')
			default:
				b.WriteByte('\\')
			}
			continue
		}
		if c == '\n' {
			return token.Token{}, errors.NewLexerError(l.file, startLine, startCol,
				"unterminated string literal",
				l.currentLineText(),
				"Close the string with \" before end of line.")
		}
		b.WriteByte(l.advance())
	}
	return token.Token{}, errors.NewLexerError(l.file, startLine, startCol,
		"unterminated string literal",
		l.currentLineText(),
		"Close the string with \" before end of file.")
}
