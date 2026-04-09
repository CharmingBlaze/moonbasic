package semantic

import (
	"testing"

	"moonbasic/compiler/parser"
)

func TestImplicitDeclaration(t *testing.T) {
	tests := []struct {
		name    string
		src     string
		wantErr bool
	}{
		{"Single Assignment", "x = 5\nPRINT(x)\n", false},
		{"Use Before Assignment", "PRINT(x)\nx = 5\n", true},
		{"Function Local", "FUNCTION f()\n  x = 5\n  PRINT(x)\nENDFUNCTION\nf()\n", false},
		{"Function Parameter", "FUNCTION f(x)\n  PRINT(x)\nENDFUNCTION\nf(10)\n", false},
		{"Global vs Local", "x = 1\nFUNCTION f()\n  PRINT(x)\nENDFUNCTION\n", false},
		{"Read Undeclared Global", "PRINT(unassigned)\n", true},
		{"Write Undeclared Array", "a(0) = 1\n", true}, // Requires DIM
		{"Dim Declares", "DIM a(10)\na(0) = 1\n", false},
		{"For Iterator", "FOR i = 0 TO 10\n PRINT(i)\nNEXT\n", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prog, err := parser.ParseSource("t.mb", tt.src)
			if err != nil {
				t.Fatalf("parse error: %v", err)
			}
			a := DefaultAnalyzer("t.mb", parser.SplitLines(tt.src))
			err = a.Run(prog)
			if (err != nil) != tt.wantErr {
				t.Errorf("got error %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
