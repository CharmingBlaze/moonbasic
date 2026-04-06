// Package ast defines abstract syntax tree nodes for moonBASIC.
package ast

import (
	"fmt"
	"strings"
)

// Stmt is a statement node.
type Stmt interface {
	stmt()
	String() string
}

// Expr is an expression node.
type Expr interface {
	expr()
	String() string
}

// Program is the root: ordered top-level statements plus type and function definitions.
type Program struct {
	Stmts     []Stmt
	Functions []*FunctionDef
	Types     []*TypeDef
}

// FunctionDef defines a user function.
type FunctionDef struct {
	Name   string
	Params []Param
	Body   []Stmt
	Line   int
	Col    int
}

// Param is a formal parameter with optional type suffix.
type Param struct {
	Name   string
	Suffix string // "", "#", "$", "?"
}

// TypeDef is a user-defined TYPE ... FIELD ... ENDTYPE.
type TypeDef struct {
	Name   string
	Fields []string
	Line   int
	Col    int
}

// AssignNode is name = expr (suffix embedded in Name e.g. "X#").
// Global is set when the statement was parsed as GLOBAL name = expr.
type AssignNode struct {
	Name   string
	Expr   Expr
	Global bool
	Line   int
	Col    int
}

func (n *AssignNode) stmt() {}
func (n *AssignNode) String() string {
	return fmt.Sprintf("Assign(%s, %s)", n.Name, n.Expr.String())
}

// IndexAssignNode is arr(i[,...]) = expr (parenthesized or bracketed indices).
type IndexAssignNode struct {
	Array string
	Index []Expr
	Expr  Expr
	Line  int
	Col   int
}

func (n *IndexAssignNode) stmt() {}
func (n *IndexAssignNode) String() string {
	return fmt.Sprintf("IndexAssign(%s[...], %s)", n.Array, n.Expr.String())
}

// FieldAssignNode is obj.field = expr.
type FieldAssignNode struct {
	Object string
	Field  string
	Expr   Expr
	Line   int
	Col    int
}

func (n *FieldAssignNode) stmt() {}
func (n *FieldAssignNode) String() string {
	return fmt.Sprintf("FieldAssign(%s.%s, %s)", n.Object, n.Field, n.Expr.String())
}

// FieldAccessNode is obj.field in expression context.
type FieldAccessNode struct {
	Object string
	Field  string
	Line   int
	Col    int
}

func (n *FieldAccessNode) expr() {}
func (n *FieldAccessNode) String() string {
	return fmt.Sprintf("FieldAccess(%s.%s)", n.Object, n.Field)
}

// CallStmtNode is a bare user function call.
type CallStmtNode struct {
	Name string
	Args []Expr
	Line int
	Col  int
}

func (n *CallStmtNode) stmt() {}
func (n *CallStmtNode) String() string {
	return fmt.Sprintf("CallStmt(%s(...))", n.Name)
}

// NamespaceCallStmt is NS.METHOD(args).
type NamespaceCallStmt struct {
	NS     string
	Method string
	Args   []Expr
	Line   int
	Col    int
}

func (n *NamespaceCallStmt) stmt() {}
func (n *NamespaceCallStmt) String() string {
	return fmt.Sprintf("NsCall(%s.%s)", n.NS, n.Method)
}

// HandleCallStmt is handleVar.METHOD(args).
type HandleCallStmt struct {
	Receiver string
	Method   string
	Args     []Expr
	Line     int
	Col      int
}

func (n *HandleCallStmt) stmt() {}
func (n *HandleCallStmt) String() string {
	return fmt.Sprintf("HandleCall(%s.%s)", n.Receiver, n.Method)
}

// IfNode represents IF / ELSEIF / ELSE / ENDIF.
type IfNode struct {
	Cond   Expr
	Then   []Stmt
	ElseIf []ElseIfClause
	Else   []Stmt
	Line   int
	Col    int
}

// ElseIfClause is ELSEIF cond THEN body.
type ElseIfClause struct {
	Cond Expr
	Body []Stmt
}

func (n *IfNode) stmt()          {}
func (n *IfNode) String() string { return "If(...)" }

// WhileNode is WHILE ... WEND.
type WhileNode struct {
	Cond Expr
	Body []Stmt
	Line int
	Col  int
}

func (n *WhileNode) stmt()          {}
func (n *WhileNode) String() string { return "While(...)" }

// ForNode is FOR var = from TO to [STEP step] ... NEXT.
type ForNode struct {
	Var  string
	From Expr
	To   Expr
	Step Expr
	Body []Stmt
	Line int
	Col  int
}

func (n *ForNode) stmt()          {}
func (n *ForNode) String() string { return "For(...)" }

// RepeatNode is REPEAT ... UNTIL.
type RepeatNode struct {
	Body      []Stmt
	Condition Expr
	Line      int
	Col       int
}

func (n *RepeatNode) stmt()          {}
func (n *RepeatNode) String() string { return "Repeat(...)" }

// DoLoopKind selects DO/LOOP variant (Raylib-style BASIC).
type DoLoopKind int

const (
	// DoPostWhile: DO ... LOOP WHILE cond (body runs at least once).
	DoPostWhile DoLoopKind = iota
	// DoPostUntil: DO ... LOOP UNTIL cond (exit when cond true).
	DoPostUntil
	// DoPreWhile: DO WHILE cond ... LOOP (may skip body).
	DoPreWhile
)

// DoLoopNode is DO ... LOOP with WHILE or UNTIL (post-test) or DO WHILE ... LOOP (pre-test).
type DoLoopNode struct {
	Kind DoLoopKind
	Cond Expr
	Body []Stmt
	Line int
	Col  int
}

func (n *DoLoopNode) stmt()          {}
func (n *DoLoopNode) String() string { return "DoLoop(...)" }

// ExitStmt is EXIT FOR | EXIT WHILE | EXIT REPEAT | EXIT DO | EXIT FUNCTION.
type ExitStmt struct {
	Target string // FOR, WHILE, REPEAT, DO, FUNCTION
	Line   int
	Col    int
}

func (n *ExitStmt) stmt()          {}
func (n *ExitStmt) String() string { return "Exit(" + n.Target + ")" }

// ContinueStmt is CONTINUE FOR | WHILE | REPEAT | DO.
type ContinueStmt struct {
	Target string // FOR, WHILE, REPEAT, DO
	Line   int
	Col    int
}

func (n *ContinueStmt) stmt()          {}
func (n *ContinueStmt) String() string { return "Continue(" + n.Target + ")" }

// SelectNode is SELECT expr ... CASE ... DEFAULT ... ENDSELECT.
type SelectNode struct {
	Expr    Expr
	Cases   []CaseClause
	Default []Stmt
	Line    int
	Col     int
}

// CaseClause is CASE value: body (value expression; body until next CASE/DEFAULT/ENDSELECT).
type CaseClause struct {
	Value Expr
	Body  []Stmt
}

func (n *SelectNode) stmt()          {}
func (n *SelectNode) String() string { return "Select(...)" }

// ReturnNode is RETURN [expr].
type ReturnNode struct {
	Expr Expr
	Line int
	Col  int
}

func (n *ReturnNode) stmt()          {}
func (n *ReturnNode) String() string { return "Return(...)" }

// GotoNode is GOTO label.
type GotoNode struct {
	Label string
	Line  int
	Col   int
}

func (n *GotoNode) stmt()          {}
func (n *GotoNode) String() string { return fmt.Sprintf("Goto(%s)", n.Label) }

// GosubNode is GOSUB label.
type GosubNode struct {
	Label string
	Line  int
	Col   int
}

func (n *GosubNode) stmt()          {}
func (n *GosubNode) String() string { return fmt.Sprintf("Gosub(%s)", n.Label) }

// LabelNode is .label.
type LabelNode struct {
	Name string
	Line int
	Col  int
}

func (n *LabelNode) stmt()          {}
func (n *LabelNode) String() string { return fmt.Sprintf("Label(.%s)", n.Name) }

// DimNode is DIM name(dim...) or REDIM [PRESERVE] name(dim...), or DIM name AS Type(dim...) for typed handle arrays.
type DimNode struct {
	Name     string
	ElemType string // non-empty: array of heap instances of this TYPE
	Dims     []Expr
	IsRedim  bool
	Preserve bool // REDIM always preserves data in moonBASIC; PRESERVE is accepted for readability
	Line     int
	Col      int
}

func (n *DimNode) stmt()          {}
func (n *DimNode) String() string { return fmt.Sprintf("Dim(%s)", n.Name) }

// IncludeNode is INCLUDE "path".
type IncludeNode struct {
	Path string
	Line int
	Col  int
}

func (n *IncludeNode) stmt()          {}
func (n *IncludeNode) String() string { return fmt.Sprintf("Include(%q)", n.Path) }

// LocalDeclNode is LOCAL name [= expr] or list — simplified as single name per line.
type LocalDeclNode struct {
	Name string
	Init Expr
	Line int
	Col  int
}

func (n *LocalDeclNode) stmt()          {}
func (n *LocalDeclNode) String() string { return fmt.Sprintf("Local(%s)", n.Name) }

// ConstDeclNode is CONST name = expr.
type ConstDeclNode struct {
	Name string
	Expr Expr
	Line int
	Col  int
}

func (n *ConstDeclNode) stmt()          {}
func (n *ConstDeclNode) String() string { return fmt.Sprintf("Const(%s)", n.Name) }

// StaticDeclNode is STATIC name [= expr] inside a FUNCTION.
type StaticDeclNode struct {
	Name string
	Init Expr
	Line int
	Col  int
}

func (n *StaticDeclNode) stmt()          {}
func (n *StaticDeclNode) String() string { return fmt.Sprintf("Static(%s)", n.Name) }

// SwapStmt is SWAP a, b.
type SwapStmt struct {
	A, B string
	Line int
	Col  int
}

func (n *SwapStmt) stmt()          {}
func (n *SwapStmt) String() string { return fmt.Sprintf("Swap(%s,%s)", n.A, n.B) }

// EraseStmt is ERASE arr — frees heap array and sets variable to NULL.
type EraseStmt struct {
	Name string
	Line int
	Col  int
}

func (n *EraseStmt) stmt()          {}
func (n *EraseStmt) String() string { return fmt.Sprintf("Erase(%s)", n.Name) }

// NewNode is NEW(TypeName).
type NewNode struct {
	TypeName string
	Line     int
	Col      int
}

func (n *NewNode) expr()          {}
func (n *NewNode) String() string { return fmt.Sprintf("New(%s)", n.TypeName) }

// DeleteStmt is DELETE expr.
type DeleteStmt struct {
	Expr Expr
	Line int
	Col  int
}

func (n *DeleteStmt) stmt()          {}
func (n *DeleteStmt) String() string { return "Delete(...)" }

// EachNode is FOR var = EACH(Type) ... NEXT (represented as ForEachStmt).
type EachStmt struct {
	Var      string
	TypeName string
	Body     []Stmt
	Line     int
	Col      int
}

func (n *EachStmt) stmt()          {}
func (n *EachStmt) String() string { return fmt.Sprintf("Each(%s in %s)", n.Var, n.TypeName) }

// ExprStatement wraps an expression used as a statement (rare).
type ExprStmt struct {
	Expr Expr
}

func (n *ExprStmt) stmt()          {}
func (n *ExprStmt) String() string { return n.Expr.String() }

// EndProgramStmt is bare END (terminate program).
type EndProgramStmt struct {
	Line int
	Col  int
}

func (n *EndProgramStmt) stmt()          {}
func (n *EndProgramStmt) String() string { return "END" }

// BinopNode is left op right.
type BinopNode struct {
	Op    string
	Left  Expr
	Right Expr
	Line  int
	Col   int
}

func (n *BinopNode) expr() {}
func (n *BinopNode) String() string {
	return fmt.Sprintf("(%s %s %s)", n.Left.String(), n.Op, n.Right.String())
}

// UnaryNode is op expr.
type UnaryNode struct {
	Op   string
	Expr Expr
	Line int
	Col  int
}

func (n *UnaryNode) expr()          {}
func (n *UnaryNode) String() string { return fmt.Sprintf("(%s %s)", n.Op, n.Expr.String()) }

// IdentNode is a variable reference (name includes suffix).
type IdentNode struct {
	Name string
	Line int
	Col  int
}

func (n *IdentNode) expr()          {}
func (n *IdentNode) String() string { return n.Name }

// IntLitNode is an integer literal.
type IntLitNode struct {
	Value int64
	Line  int
	Col   int
}

func (n *IntLitNode) expr()          {}
func (n *IntLitNode) String() string { return fmt.Sprintf("%d", n.Value) }

// FloatLitNode is a float literal.
type FloatLitNode struct {
	Value float64
	Lit   string
	Line  int
	Col   int
}

func (n *FloatLitNode) expr()          {}
func (n *FloatLitNode) String() string { return n.Lit }

// StringLitNode is a string literal.
type StringLitNode struct {
	Value string
	Line  int
	Col   int
}

func (n *StringLitNode) expr()          {}
func (n *StringLitNode) String() string { return fmt.Sprintf("%q", n.Value) }

// BoolLitNode is TRUE or FALSE.
type BoolLitNode struct {
	Value bool
	Line  int
	Col   int
}

func (n *BoolLitNode) expr()          {}
func (n *BoolLitNode) String() string { return fmt.Sprintf("%v", n.Value) }

// NullLitNode is NULL.
type NullLitNode struct {
	Line int
	Col  int
}

func (n *NullLitNode) expr()          {}
func (n *NullLitNode) String() string { return "NULL" }

// CallExprNode is user function call in expression context.
type CallExprNode struct {
	Name string
	Args []Expr
	Line int
	Col  int
}

func (n *CallExprNode) expr()          {}
func (n *CallExprNode) String() string { return fmt.Sprintf("%s(...)", n.Name) }

// NamespaceCallExpr is NS.METHOD(args) in expression context.
type NamespaceCallExpr struct {
	NS     string
	Method string
	Args   []Expr
	Line   int
	Col    int
}

func (n *NamespaceCallExpr) expr() {}
func (n *NamespaceCallExpr) String() string {
	return fmt.Sprintf("%s.%s(...)", n.NS, n.Method)
}

// HandleCallExpr is recv.METHOD(args) in expression context.
type HandleCallExpr struct {
	Receiver string
	Method   string
	Args     []Expr
	Line     int
	Col      int
}

func (n *HandleCallExpr) expr() {}
func (n *HandleCallExpr) String() string {
	return fmt.Sprintf("%s.%s(...)", n.Receiver, n.Method)
}

// IndexExpr is base(index...) using paren or bracket (lexer normalised).
type IndexExpr struct {
	Base  Expr
	Index []Expr
	Line  int
	Col   int
}

func (n *IndexExpr) expr()          {}
func (n *IndexExpr) String() string { return fmt.Sprintf("%s[...]", n.Base.String()) }

// IndexFieldExpr is arr(idx...).field — read a field on an array element (handle).
type IndexFieldExpr struct {
	Array string
	Index []Expr
	Field string
	Line  int
	Col   int
}

func (n *IndexFieldExpr) expr()          {}
func (n *IndexFieldExpr) String() string { return fmt.Sprintf("%s(...).%s", n.Array, n.Field) }

// IndexFieldAssignNode is arr(idx...).field = expr.
type IndexFieldAssignNode struct {
	Array string
	Index []Expr
	Field string
	Expr  Expr
	Line  int
	Col   int
}

func (n *IndexFieldAssignNode) stmt()          {}
func (n *IndexFieldAssignNode) String() string { return fmt.Sprintf("IndexFieldAssign(%s.%s)", n.Array, n.Field) }

// GroupedExpr is ( expr ).
type GroupedExpr struct {
	Inner Expr
}

func (n *GroupedExpr) expr()          {}
func (n *GroupedExpr) String() string { return fmt.Sprintf("(%s)", n.Inner.String()) }

// PrettyPrint writes an indented tree for debugging.
func PrettyPrint(p *Program) string {
	var b strings.Builder
	for _, t := range p.Types {
		fmt.Fprintf(&b, "TYPE %s ...\n", t.Name)
	}
	for _, f := range p.Functions {
		fmt.Fprintf(&b, "FUNCTION %s(...)\n", f.Name)
	}
	for _, s := range p.Stmts {
		fmt.Fprintf(&b, "%s\n", s.String())
	}
	return b.String()
}
