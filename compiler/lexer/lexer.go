// Package lexer implements the moonBASIC tokeniser.
//
// All identifiers are normalised to uppercase at scan time. String literal
// contents are never uppercased. [ and ] map to LPAREN and RPAREN. CRLF is
// normalised to LF before scanning.
//
// Sub-files (lexer_ident.go, lexer_number.go, lexer_string.go) rely on this
// Lexer API only:
//
//	Fields:  input, file, pos, line, col, lineStart
//	Methods: peek() byte, advance() byte, eof() bool, currentLineText() string
//
// Multi-word END (lexer_ident.go expandEndKeyword): after keyword END, the next
// non-space word is inspected so END IF emits token.ENDIF (not END then IF), and
// END FUNCTION → ENDFUNCTION, END WHILE → WEND, END SELECT → ENDSELECT, END TYPE → ENDTYPE;
// bare END stays token.END.
//
// Colon (token.COLON) is emitted for ':'; the parser chains statements on one line.
// Brackets '[' and ']' are mapped to LPAREN and RPAREN (same as '(' and ')').
//
// Tokens use the shared schema: token.Token{Type, Lit, Line, Col}.
package lexer

import (
	"fmt"

	"moonbasic/compiler/errors"
	"moonbasic/compiler/intern"
	"moonbasic/compiler/token"
)

// Lexer holds scanner state for one source file.
//
// lexer_ident.go / lexer_number.go / lexer_string.go require:
//   - input, lineStart: buffer and start offset of the current line (for errors)
//   - file: logical path for MoonError
//   - pos, line, col: cursor position (1-based line/col)
//   - peek, advance, eof, currentLineText
type Lexer struct {
	input     string
	file      string
	pos       int
	line      int
	col       int
	lineStart int
	names     *intern.Table // identifier / keyword spellings (nil = no interning)
}

// normalizeLineEndings converts CRLF and lone CR to LF before scanning.
func normalizeLineEndings(input string) string {
	out := make([]byte, 0, len(input))
	for i := 0; i < len(input); i++ {
		if input[i] == '\r' {
			out = append(out, '\n')
			if i+1 < len(input) && input[i+1] == '\n' {
				i++
			}
			continue
		}
		out = append(out, input[i])
	}
	return string(out)
}

// New builds a lexer. file is used only in error messages.
func New(file, input string) *Lexer {
	return NewWithIntern(file, input, intern.New())
}

// NewWithIntern builds a lexer using the given intern table (may be nil).
func NewWithIntern(file, input string, tab *intern.Table) *Lexer {
	in := normalizeLineEndings(input)
	return &Lexer{
		input:     in,
		file:      file,
		pos:       0,
		line:      1,
		col:       1,
		lineStart: 0,
		names:     tab,
	}
}

// File returns the logical source path.
func (l *Lexer) File() string { return l.file }

// Scan tokenises the full source; stops at EOF or first error.
func Scan(file, input string) ([]token.Token, error) {
	lex := New(file, input) // fresh intern table per Scan / compile unit
	var toks []token.Token
	for {
		tok, err := lex.NextToken()
		if err != nil {
			return nil, err
		}
		toks = append(toks, tok)
		if tok.Type == token.EOF {
			break
		}
	}
	return toks, nil
}

// NextToken returns the next token (Type, Lit, Line, Col) or an error.
func (l *Lexer) NextToken() (token.Token, error) {
	for !l.eof() {
		c := l.peek()
		if c != ' ' && c != '\t' {
			break
		}
		l.advance()
	}

	if l.eof() {
		return token.Token{Type: token.EOF, Lit: "", Line: l.line, Col: l.col}, nil
	}

	ch := l.peek()

	if ch == '\n' {
		line, col := l.line, l.col
		l.advance()
		return token.Token{Type: token.NEWLINE, Lit: "\n", Line: line, Col: col}, nil
	}

	if ch == ';' {
		for !l.eof() && l.peek() != '\n' {
			l.advance()
		}
		return l.NextToken()
	}

	if ch == '"' {
		return l.lexString()
	}

	if ch == '[' {
		line, col := l.line, l.col
		l.advance()
		return token.Token{Type: token.LPAREN, Lit: "[", Line: line, Col: col}, nil
	}
	if ch == ']' {
		line, col := l.line, l.col
		l.advance()
		return token.Token{Type: token.RPAREN, Lit: "]", Line: line, Col: col}, nil
	}

	if isIdentStart(ch) {
		return l.lexIdent()
	}

	if isDigit(ch) {
		return l.lexNumber()
	}
	if ch == '-' && !l.eof() && isDigit(l.peek()) {
		return l.lexNumber()
	}

	return l.lexOp()
}

// lexOp handles punctuation and operators not handled above.
func (l *Lexer) lexOp() (token.Token, error) {
	ch := l.peek()
	line, col := l.line, l.col

	emit := func(tt token.TokenType, lit string) (token.Token, error) {
		return token.Token{Type: tt, Lit: lit, Line: line, Col: col}, nil
	}

	advanceEmit := func(tt token.TokenType, lit string) (token.Token, error) {
		l.advance()
		return emit(tt, lit)
	}

	switch ch {
	case '(':
		return advanceEmit(token.LPAREN, "(")
	case ')':
		return advanceEmit(token.RPAREN, ")")
	case ',':
		return advanceEmit(token.COMMA, ",")
	case ':':
		return advanceEmit(token.COLON, ":")
	case '.':
		return advanceEmit(token.DOT, ".")
	case '#':
		return advanceEmit(token.HASH, "#")
	case '$':
		return advanceEmit(token.DOLLAR, "$")
	case '?':
		return advanceEmit(token.QUESTION, "?")
	case '^':
		return advanceEmit(token.CARET, "^")
	case '=':
		l.advance()
		return emit(token.EQ, "=")
	case '+':
		l.advance()
		if !l.eof() && l.peek() == '=' {
			l.advance()
			return emit(token.PLUSEQ, "+=")
		}
		return emit(token.PLUS, "+")
	case '-':
		l.advance()
		if !l.eof() && l.peek() == '=' {
			l.advance()
			return emit(token.MINUSEQ, "-=")
		}
		return emit(token.MINUS, "-")
	case '*':
		l.advance()
		if !l.eof() && l.peek() == '=' {
			l.advance()
			return emit(token.STAREQ, "*=")
		}
		return emit(token.STAR, "*")
	case '/':
		l.advance()
		if !l.eof() && l.peek() == '=' {
			l.advance()
			return emit(token.SLASHEQ, "/=")
		}
		return emit(token.SLASH, "/")
	case '<':
		l.advance()
		if !l.eof() && l.peek() == '>' {
			l.advance()
			return emit(token.NEQ, "<>")
		}
		if !l.eof() && l.peek() == '=' {
			l.advance()
			return emit(token.LTE, "<=")
		}
		return emit(token.LT, "<")
	case '>':
		l.advance()
		if !l.eof() && l.peek() == '=' {
			l.advance()
			return emit(token.GTE, ">=")
		}
		return emit(token.GT, ">")
	}

	return token.Token{}, errors.NewLexerError(
		l.file, line, col,
		fmt.Sprintf("unexpected character %q", ch),
		l.currentLineText(),
		"Remove or replace the unrecognised character",
	)
}

func (l *Lexer) eof() bool {
	return l.pos >= len(l.input)
}

// peek returns the current byte, or 0 at EOF.
func (l *Lexer) peek() byte {
	if l.eof() {
		return 0
	}
	return l.input[l.pos]
}

// advance consumes one byte and updates line, col, lineStart.
func (l *Lexer) advance() byte {
	if l.eof() {
		return 0
	}
	ch := l.input[l.pos]
	l.pos++
	if ch == '\n' {
		l.line++
		l.col = 1
		l.lineStart = l.pos
	} else {
		l.col++
	}
	return ch
}

// currentLineText returns the source line containing the cursor (for errors).
func (l *Lexer) currentLineText() string {
	end := l.lineStart
	for end < len(l.input) && l.input[end] != '\n' {
		end++
	}
	return l.input[l.lineStart:end]
}
