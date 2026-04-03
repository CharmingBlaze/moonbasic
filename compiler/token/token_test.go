package token

import "testing"

func TestLookupKeyword(t *testing.T) {
	if LookupKeyword("IF") != IF {
		t.Fatal()
	}
	if LookupKeyword("WHILE") != WHILE {
		t.Fatal()
	}
	if LookupKeyword("UNKNOWNVAR") != IDENT {
		t.Fatal()
	}
}

func TestTokenTypeString(t *testing.T) {
	if EOF.String() != "EOF" {
		t.Fatalf("got %q", EOF.String())
	}
}
