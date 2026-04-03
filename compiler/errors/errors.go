// Package errors provides structured moonBASIC compiler and runtime errors.
package errors

import (
	"fmt"
	"strings"
)

// Category classifies the error source.
type Category string

const (
	Lexer    Category = "Lexer Error"
	Parse    Category = "Parse Error"
	TypeErr  Category = "Type Error"
	Runtime  Category = "Runtime Error"
	CodeGen  Category = "CodeGen Error"
)

// MoonError is a structured error with source location and optional hint.
type MoonError struct {
	Category   Category
	File       string
	Line       int
	Col        int
	Message    string
	SourceLine string
	Hint       string
}

func (e *MoonError) Error() string {
	return Format(e)
}

// Format renders the canonical moonBASIC error block.
// The caret aligns with column Col (1-based) on the source line, after the "  N | " gutter.
func Format(e *MoonError) string {
	var b strings.Builder
	fmt.Fprintf(&b, "[moonBASIC] %s in %s line %d col %d:\n", e.Category, e.File, e.Line, e.Col)
	fmt.Fprintf(&b, "  %s\n\n", e.Message)
	if e.SourceLine != "" {
		gutter := fmt.Sprintf("  %d | ", e.Line)
		fmt.Fprintf(&b, "%s%s\n", gutter, e.SourceLine)
		col := e.Col
		if col < 1 {
			col = 1
		}
		under := strings.Repeat(" ", len(gutter)+col-1)
		fmt.Fprintf(&b, "%s^\n", under)
	}
	if e.Hint != "" {
		fmt.Fprintf(&b, "  Hint: %s\n", e.Hint)
	}
	return b.String()
}

// NewLexerError builds a lexer error.
func NewLexerError(file string, line, col int, message, sourceLine, hint string) *MoonError {
	return &MoonError{Category: Lexer, File: file, Line: line, Col: col, Message: message, SourceLine: sourceLine, Hint: hint}
}

// NewParseError builds a parse error.
func NewParseError(file string, line, col int, message, sourceLine, hint string) *MoonError {
	return &MoonError{Category: Parse, File: file, Line: line, Col: col, Message: message, SourceLine: sourceLine, Hint: hint}
}

// NewTypeError builds a type error.
func NewTypeError(file string, line, col int, message, sourceLine, hint string) *MoonError {
	return &MoonError{Category: TypeErr, File: file, Line: line, Col: col, Message: message, SourceLine: sourceLine, Hint: hint}
}

// NewRuntimeError builds a runtime error.
func NewRuntimeError(file string, line, col int, message, sourceLine, hint string) *MoonError {
	return &MoonError{Category: Runtime, File: file, Line: line, Col: col, Message: message, SourceLine: sourceLine, Hint: hint}
}

// NewCodeGenError builds a codegen error.
func NewCodeGenError(file string, line, col int, message, sourceLine, hint string) *MoonError {
	return &MoonError{Category: CodeGen, File: file, Line: line, Col: col, Message: message, SourceLine: sourceLine, Hint: hint}
}

// IsMoonError reports whether err is a *MoonError.