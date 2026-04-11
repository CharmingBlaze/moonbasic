package lexer

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"moonbasic/compiler/token"
)

func TestLexReferenceProgram(t *testing.T) {
	_, file, _, _ := runtime.Caller(0)
	dir := filepath.Dir(file)
	refPath := filepath.Join(dir, "..", "..", "testdata", "reference.mbc")
	src, err := os.ReadFile(refPath)
	if err != nil {
		t.Fatal(err)
	}
	l := New(refPath, string(src))
	var toks []token.Token
	for {
		tok, err := l.NextToken()
		if err != nil {
			t.Fatal(err)
		}
		toks = append(toks, tok)
		if tok.Type == token.EOF {
			break
		}
	}
	if len(toks) < 100 {
		t.Fatalf("expected many tokens, got %d", len(toks))
	}
	if toks[0].Type != token.IDENT || toks[0].Lit != "WINDOW" {
		t.Fatalf("first token: %+v", toks[0])
	}
	if toks[1].Type != token.DOT {
		t.Fatal()
	}
	if toks[2].Type != token.IDENT || toks[2].Lit != "OPEN" {
		t.Fatal()
	}
	// Find WHILE
	foundWhile := false
	for _, tk := range toks {
		if tk.Type == token.WHILE {
			foundWhile = true
			break
		}
	}
	if !foundWhile {
		t.Fatal("expected WHILE")
	}
	// dt name (no suffix)
	foundDt := false
	for _, tk := range toks {
		if tk.Type == token.IDENT && tk.Lit == "DT" {
			foundDt = true
			break
		}
	}
	if !foundDt {
		t.Fatal("expected DT ident")
	}
}

func TestEndIfExpansion(t *testing.T) {
	l := New("t.mbc", "END IF\n")
	tok, err := l.NextToken()
	if err != nil {
		t.Fatal(err)
	}
	if tok.Type != token.ENDIF {
		t.Fatalf("got %v", tok)
	}
}

func TestEndFunctionTwoWords(t *testing.T) {
	l := New("t.mbc", "END FUNCTION\n")
	tok, err := l.NextToken()
	if err != nil {
		t.Fatal(err)
	}
	if tok.Type != token.ENDFUNCTION {
		t.Fatalf("got %v want ENDFUNCTION", tok)
	}
}

func TestEndBare(t *testing.T) {
	l := New("t.mbc", "END\nPRINT")
	t1, err := l.NextToken()
	if err != nil {
		t.Fatal(err)
	}
	if t1.Type != token.END {
		t.Fatalf("got %v", t1)
	}
	t2, err := l.NextToken()
	if err != nil {
		t.Fatal(err)
	}
	if t2.Type != token.NEWLINE {
		t.Fatalf("got %v", t2)
	}
	t3, err := l.NextToken()
	if err != nil {
		t.Fatal(err)
	}
	if t3.Type != token.IDENT || t3.Lit != "PRINT" {
		t.Fatalf("got %v", t3)
	}
}

func TestBracketArray(t *testing.T) {
	l := New("t.mbc", "arr[0]")
	var types []token.TokenType
	for {
		tok, err := l.NextToken()
		if err != nil {
			t.Fatal(err)
		}
		types = append(types, tok.Type)
		if tok.Type == token.EOF {
			break
		}
	}
	// IDENT LPAREN INT RPAREN as [ maps to LPAREN
	want := []token.TokenType{token.IDENT, token.LPAREN, token.INT, token.RPAREN, token.EOF}
	if len(types) != len(want) {
		t.Fatalf("got %v", types)
	}
	for i := range want {
		if types[i] != want[i] {
			t.Fatalf("at %d: got %v want %v", i, types[i], want[i])
		}
	}
}
