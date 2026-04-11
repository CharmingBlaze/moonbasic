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
	Accept(v Visitor)
}

// Expr is an expression node.
type Expr interface {
	expr()
	String() string
	Accept(v Visitor)
}

// Visitor allows walking the AST without massive type switches.
// If VisitStmt or VisitExpr returns a non-nil Visitor, children are visited with it.
type Visitor interface {
	VisitStmt(Stmt) Visitor
	VisitExpr(Expr) Visitor
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

// Param is a formal parameter.
type Param struct {
	Name   string
}

// TypeDef is a user-defined TYPE ... FIELD ... ENDTYPE.
type TypeDef struct {
	Name   string
	Fields []string // field names in order
	// FieldTypeHints is parallel to Fields: "INTEGER", "FLOAT", "STRING", user type name, or "" (legacy).
	FieldTypeHints []string
	// FieldIsArray is parallel to Fields: true if "name AS type(dim...)" had dimension parens (not yet supported at runtime).
	FieldIsArray []bool
	Line         int
	Col          int
}

// AssignNode is name = expr.
// Global is set when the statement was parsed as GLOBAL name = expr.
type AssignNode struct {
	Name   string
	Expr   Expr
	Global bool
	Line   int
	Col    int
}

func (n *AssignNode) stmt()            {}
func (n *AssignNode) Accept(v Visitor) { v.VisitStmt(n) }
func (n *AssignNode) String() string {
	return fmt.Sprintf("Assign(%s, %s)", n.Name, n.Expr.String())
}

// MultiAssignNode is a, b, ... = expr (destructuring assignment).
// Current lowering unpacks from a 1-D array-like handle using 1-based indices.
type MultiAssignNode struct {
	Names []string
	Expr  Expr
	Line  int
	Col   int
}

func (n *MultiAssignNode) stmt()            {}
func (n *MultiAssignNode) Accept(v Visitor) { v.VisitStmt(n) }
func (n *MultiAssignNode) String() string {
	return fmt.Sprintf("MultiAssign(%v, %s)", n.Names, n.Expr.String())
}

// IndexAssignNode is arr(i[,...]) = expr (parenthesized or bracketed indices).
type IndexAssignNode struct {
	Array string
	Index []Expr
	Expr  Expr
	Line  int
	Col   int
}

func (n *IndexAssignNode) stmt()            {}
func (n *IndexAssignNode) Accept(v Visitor) { v.VisitStmt(n) }
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

func (n *FieldAssignNode) stmt()            {}
func (n *FieldAssignNode) Accept(v Visitor) { v.VisitStmt(n) }
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

func (n *FieldAccessNode) expr()            {}
func (n *FieldAccessNode) Accept(v Visitor) { v.VisitExpr(n) }
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

func (n *CallStmtNode) stmt()            {}
func (n *CallStmtNode) Accept(v Visitor) { v.VisitStmt(n) }
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

func (n *NamespaceCallStmt) stmt()            {}
func (n *NamespaceCallStmt) Accept(v Visitor) { v.VisitStmt(n) }
func (n *NamespaceCallStmt) String() string {
	return fmt.Sprintf("NsCall(%s.%s)", n.NS, n.Method)
}

// NamespaceAssignNode is NS.METHOD(args) = expr.
type NamespaceAssignNode struct {
	NS     string
	Method string
	Args   []Expr
	Expr   Expr
	Line   int
	Col    int
}

func (n *NamespaceAssignNode) stmt()            {}
func (n *NamespaceAssignNode) Accept(v Visitor) { v.VisitStmt(n) }
func (n *NamespaceAssignNode) String() string {
	return fmt.Sprintf("NsAssign(%s.%s, %s)", n.NS, n.Method, n.Expr.String())
}

// HandleCallStmt is handleVar.METHOD(args).
type HandleCallStmt struct {
	Receiver string
	Method   string
	Args     []Expr
	Line     int
	Col      int
}

func (n *HandleCallStmt) stmt()            {}
func (n *HandleCallStmt) Accept(v Visitor) { v.VisitStmt(n) }
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

func (n *IfNode) stmt()            {}
func (n *IfNode) Accept(v Visitor) { v.VisitStmt(n) }
func (n *IfNode) String() string   { return "If(...)" }

// WhileNode is WHILE ... WEND.
type WhileNode struct {
	Cond Expr
	Body []Stmt
	Line int
	Col  int
}

func (n *WhileNode) stmt()            {}
func (n *WhileNode) Accept(v Visitor) { v.VisitStmt(n) }
func (n *WhileNode) String() string   { return "While(...)" }

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

func (n *ForNode) stmt()            {}
func (n *ForNode) Accept(v Visitor) { v.VisitStmt(n) }
func (n *ForNode) String() string   { return "For(...)" }

// RepeatNode is REPEAT ... UNTIL.
type RepeatNode struct {
	Body      []Stmt
	Condition Expr
	Line      int
	Col       int
}

func (n *RepeatNode) stmt()            {}
func (n *RepeatNode) Accept(v Visitor) { v.VisitStmt(n) }
func (n *RepeatNode) String() string   { return "Repeat(...)" }

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

func (n *DoLoopNode) stmt()            {}
func (n *DoLoopNode) Accept(v Visitor) { v.VisitStmt(n) }
func (n *DoLoopNode) String() string   { return "DoLoop(...)" }

// ExitStmt is EXIT FOR | EXIT WHILE | EXIT REPEAT | EXIT DO | EXIT FUNCTION.
type ExitStmt struct {
	Target string // FOR, WHILE, REPEAT, DO, FUNCTION
	Line   int
	Col    int
}

func (n *ExitStmt) stmt()            {}
func (n *ExitStmt) Accept(v Visitor) { v.VisitStmt(n) }
func (n *ExitStmt) String() string   { return "Exit(" + n.Target + ")" }

// ContinueStmt is CONTINUE FOR | WHILE | REPEAT | DO.
type ContinueStmt struct {
	Target string // FOR, WHILE, REPEAT, DO
	Line   int
	Col    int
}

func (n *ContinueStmt) stmt()            {}
func (n *ContinueStmt) Accept(v Visitor) { v.VisitStmt(n) }
func (n *ContinueStmt) String() string   { return "Continue(" + n.Target + ")" }

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

func (n *SelectNode) stmt()            {}
func (n *SelectNode) Accept(v Visitor) { v.VisitStmt(n) }
func (n *SelectNode) String() string   { return "Select(...)" }

// ReturnNode is RETURN [expr].
type ReturnNode struct {
	Expr Expr
	Line int
	Col  int
}

func (n *ReturnNode) stmt()            {}
func (n *ReturnNode) Accept(v Visitor) { v.VisitStmt(n) }
func (n *ReturnNode) String() string   { return "Return(...)" }

// GotoNode is GOTO label.
type GotoNode struct {
	Label string
	Line  int
	Col   int
}

func (n *GotoNode) stmt()            {}
func (n *GotoNode) Accept(v Visitor) { v.VisitStmt(n) }
func (n *GotoNode) String() string   { return fmt.Sprintf("Goto(%s)", n.Label) }

// GosubNode is GOSUB label.
type GosubNode struct {
	Label string
	Line  int
	Col   int
}

func (n *GosubNode) stmt()            {}
func (n *GosubNode) Accept(v Visitor) { v.VisitStmt(n) }
func (n *GosubNode) String() string   { return fmt.Sprintf("Gosub(%s)", n.Label) }

// LabelNode is .label.
type LabelNode struct {
	Name string
	Line int
	Col  int
}

func (n *LabelNode) stmt()            {}
func (n *LabelNode) Accept(v Visitor) { v.VisitStmt(n) }
func (n *LabelNode) String() string   { return fmt.Sprintf("Label(.%s)", n.Name) }

// DimNode is DIM name(dim...), REDIM, DIM name AS Type(dim...), or name AS Type(dim...) (no DIM).
type DimNode struct {
	Name     string
	TypeName string // "INTEGER", "FLOAT", "STRING", "HANDLE", user TYPE name, or "" for untyped DIM
	Dims     []Expr
	IsRedim  bool
	Preserve bool // REDIM always preserves data in moonBASIC; PRESERVE is accepted for readability
	Line     int
	Col      int
}

func (n *DimNode) stmt()            {}
func (n *DimNode) Accept(v Visitor) { v.VisitStmt(n) }
func (n *DimNode) String() string   { return fmt.Sprintf("Dim(%s)", n.Name) }

// IncludeNode is INCLUDE "path".
type IncludeNode struct {
	Path string
	Line int
	Col  int
}

func (n *IncludeNode) stmt()            {}
func (n *IncludeNode) Accept(v Visitor) { v.VisitStmt(n) }
func (n *IncludeNode) String() string   { return fmt.Sprintf("Include(%q)", n.Path) }

// LocalDeclNode is LOCAL name [= expr] or list — simplified as single name per line.
type LocalDeclNode struct {
	Name string
	Init Expr
	Line int
	Col  int
}

func (n *LocalDeclNode) stmt()            {}
func (n *LocalDeclNode) Accept(v Visitor) { v.VisitStmt(n) }
func (n *LocalDeclNode) String() string   { return fmt.Sprintf("Local(%s)", n.Name) }

// ConstDeclNode is CONST name = expr.
type ConstDeclNode struct {
	Name string
	Expr Expr
	Line int
	Col  int
}

func (n *ConstDeclNode) stmt()            {}
func (n *ConstDeclNode) Accept(v Visitor) { v.VisitStmt(n) }
func (n *ConstDeclNode) String() string   { return fmt.Sprintf("Const(%s)", n.Name) }

// StaticDeclNode is STATIC name [= expr] inside a FUNCTION.
type StaticDeclNode struct {
	Name string
	Init Expr
	Line int
	Col  int
}

func (n *StaticDeclNode) stmt()            {}
func (n *StaticDeclNode) Accept(v Visitor) { v.VisitStmt(n) }
func (n *StaticDeclNode) String() string   { return fmt.Sprintf("Static(%s)", n.Name) }

// SwapStmt is SWAP a, b.
type SwapStmt struct {
	A, B string
	Line int
	Col  int
}

func (n *SwapStmt) stmt()            {}
func (n *SwapStmt) Accept(v Visitor) { v.VisitStmt(n) }
func (n *SwapStmt) String() string   { return fmt.Sprintf("Swap(%s,%s)", n.A, n.B) }

// EraseStmt is ERASE name — frees heap array and sets variable to NULL, or ERASE ALL (see codegen).
type EraseStmt struct {
	Name string
	Line int
	Col  int
}

func (n *EraseStmt) stmt()            {}
func (n *EraseStmt) Accept(v Visitor) { v.VisitStmt(n) }
func (n *EraseStmt) String() string   { return fmt.Sprintf("Erase(%s)", n.Name) }

// NewNode is NEW(TypeName).
type NewNode struct {
	TypeName string
	Line     int
	Col      int
}

func (n *NewNode) expr()            {}
func (n *NewNode) Accept(v Visitor) { v.VisitExpr(n) }
func (n *NewNode) String() string   { return fmt.Sprintf("New(%s)", n.TypeName) }

// DeleteStmt is DELETE expr.
type DeleteStmt struct {
	Expr Expr
	Line int
	Col  int
}

func (n *DeleteStmt) stmt()            {}
func (n *DeleteStmt) Accept(v Visitor) { v.VisitStmt(n) }
func (n *DeleteStmt) String() string   { return "Delete(...)" }

// EachNode is FOR var = EACH(Type) ... NEXT (represented as ForEachStmt).
type EachStmt struct {
	Var      string
	TypeName string
	Body     []Stmt
	Line     int
	Col      int
}

func (n *EachStmt) stmt()            {}
func (n *EachStmt) Accept(v Visitor) { v.VisitStmt(n) }
func (n *EachStmt) String() string   { return fmt.Sprintf("Each(%s in %s)", n.Var, n.TypeName) }

// ExprStatement wraps an expression used as a statement (rare).
type ExprStmt struct {
	Expr Expr
}

func (n *ExprStmt) stmt()            {}
func (n *ExprStmt) Accept(v Visitor) { v.VisitStmt(n) }
func (n *ExprStmt) String() string   { return n.Expr.String() }

// EndProgramStmt is bare END (terminate program).
type EndProgramStmt struct {
	Line int
	Col  int
}

func (n *EndProgramStmt) stmt()            {}
func (n *EndProgramStmt) Accept(v Visitor) { v.VisitStmt(n) }
func (n *EndProgramStmt) String() string   { return "END" }

// BinopNode is left op right.
type BinopNode struct {
	Op    string
	Left  Expr
	Right Expr
	Line  int
	Col   int
}

func (n *BinopNode) expr()            {}
func (n *BinopNode) Accept(v Visitor) { v.VisitExpr(n) }
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

func (n *UnaryNode) expr()            {}
func (n *UnaryNode) Accept(v Visitor) { v.VisitExpr(n) }
func (n *UnaryNode) String() string   { return fmt.Sprintf("(%s %s)", n.Op, n.Expr.String()) }

// IdentNode is a variable reference (name includes suffix).
type IdentNode struct {
	Name string
	Line int
	Col  int
}

func (n *IdentNode) expr()            {}
func (n *IdentNode) Accept(v Visitor) { v.VisitExpr(n) }
func (n *IdentNode) String() string   { return n.Name }

// IntLitNode is an integer literal.
type IntLitNode struct {
	Value int64
	Line  int
	Col   int
}

func (n *IntLitNode) expr()            {}
func (n *IntLitNode) Accept(v Visitor) { v.VisitExpr(n) }
func (n *IntLitNode) String() string   { return fmt.Sprintf("%d", n.Value) }

// FloatLitNode is a float literal.
type FloatLitNode struct {
	Value float64
	Lit   string
	Line  int
	Col   int
}

func (n *FloatLitNode) expr()            {}
func (n *FloatLitNode) Accept(v Visitor) { v.VisitExpr(n) }
func (n *FloatLitNode) String() string   { return n.Lit }

// StringLitNode is a string literal.
type StringLitNode struct {
	Value string
	Line  int
	Col   int
}

func (n *StringLitNode) expr()            {}
func (n *StringLitNode) Accept(v Visitor) { v.VisitExpr(n) }
func (n *StringLitNode) String() string   { return fmt.Sprintf("%q", n.Value) }

// BoolLitNode is TRUE or FALSE.
type BoolLitNode struct {
	Value bool
	Line  int
	Col   int
}

func (n *BoolLitNode) expr()            {}
func (n *BoolLitNode) Accept(v Visitor) { v.VisitExpr(n) }
func (n *BoolLitNode) String() string   { return fmt.Sprintf("%v", n.Value) }

// NullLitNode is NULL.
type NullLitNode struct {
	Line int
	Col  int
}

func (n *NullLitNode) expr()            {}
func (n *NullLitNode) Accept(v Visitor) { v.VisitExpr(n) }
func (n *NullLitNode) String() string   { return "NULL" }

// CallExprNode is user function call in expression context.
type CallExprNode struct {
	Name string
	Args []Expr
	Line int
	Col  int
}

func (n *CallExprNode) expr()            {}
func (n *CallExprNode) Accept(v Visitor) { v.VisitExpr(n) }
func (n *CallExprNode) String() string   { return fmt.Sprintf("%s(...)", n.Name) }

// NamespaceCallExpr is NS.METHOD(args) in expression context.
type NamespaceCallExpr struct {
	NS     string
	Method string
	Args   []Expr
	Line   int
	Col    int
}

func (n *NamespaceCallExpr) expr()            {}
func (n *NamespaceCallExpr) Accept(v Visitor) { v.VisitExpr(n) }
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

func (n *HandleCallExpr) expr()            {}
func (n *HandleCallExpr) Accept(v Visitor) { v.VisitExpr(n) }
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

func (n *IndexExpr) expr()            {}
func (n *IndexExpr) Accept(v Visitor) { v.VisitExpr(n) }
func (n *IndexExpr) String() string   { return fmt.Sprintf("%s[...]", n.Base.String()) }

// IndexFieldExpr is arr(idx...).field — read a field on an array element (handle).
type IndexFieldExpr struct {
	Array string
	Index []Expr
	Field string
	Line  int
	Col   int
}

func (n *IndexFieldExpr) expr()            {}
func (n *IndexFieldExpr) Accept(v Visitor) { v.VisitExpr(n) }
func (n *IndexFieldExpr) String() string   { return fmt.Sprintf("%s(...).%s", n.Array, n.Field) }

// IndexFieldAssignNode is arr(idx...).field = expr.
type IndexFieldAssignNode struct {
	Array string
	Index []Expr
	Field string
	Expr  Expr
	Line  int
	Col   int
}

func (n *IndexFieldAssignNode) stmt()            {}
func (n *IndexFieldAssignNode) Accept(v Visitor) { v.VisitStmt(n) }
func (n *IndexFieldAssignNode) String() string {
	return fmt.Sprintf("IndexFieldAssign(%s.%s)", n.Array, n.Field)
}

// GroupedExpr is ( expr ).
type GroupedExpr struct {
	Inner Expr
}

func (n *GroupedExpr) expr()            {}
func (n *GroupedExpr) Accept(v Visitor) { v.VisitExpr(n) }
func (n *GroupedExpr) String() string   { return fmt.Sprintf("(%s)", n.Inner.String()) }

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
