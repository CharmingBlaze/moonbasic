// Package codegen lowers moonBASIC AST into opcode.Program bytecode.
// It is split into multiple files for modularity:
// - codegen.go: Structural base and orchestration
// - codegen_expr.go: Expression and literal emission
// - codegen_stmts.go: Statement and control flow emission
// - codegen_calls.go: Builtin and user call resolution
package codegen

import (
	"fmt"
	"strings"

	"moonbasic/compiler/ast"
	"moonbasic/compiler/errors"
	"moonbasic/compiler/opt"
	"moonbasic/compiler/symtable"
	"moonbasic/vm/opcode"
)

// CodeGen emits bytecode into a Program.
type CodeGen struct {
	File    string
	Lines   []string
	Prog    *opcode.Program
	Symbols *symtable.Table // Current scope tracker
	err     error

	selectTmpID int
	fnDepth     int    // >0 when emitting a FUNCTION body
	funcName    string // uppercase function name when fnDepth > 0
	loopStack   []loopFrame
}

// New creates a code generator.
func New(file string, lines []string) *CodeGen {
	return &CodeGen{
		File:    file,
		Lines:   lines,
		Prog:    opcode.NewProgram(),
		Symbols: symtable.New(),
	}
}

func (g *CodeGen) lineText(line int) string {
	if line < 1 || line > len(g.Lines) {
		return ""
	}
	return g.Lines[line-1]
}

func (g *CodeGen) codegenError(line, col int, msg, hint string) {
	if g.err != nil {
		return
	}
	g.err = errors.NewCodeGenError(g.File, line, col, msg, g.lineText(line), hint)
}

// argCountFlags encodes call arity for IR v2 (Flags byte, max 255).
func (g *CodeGen) argCountFlags(count int, line, col int) uint8 {
	if count > 255 {
		g.codegenError(line, col, fmt.Sprintf("too many arguments (%d); IR v2 allows at most 255", count), "")
		return 0
	}
	return uint8(count)
}

func arrayKindFlags(name string) uint8 {
	if len(name) == 0 {
		return 0
	}
	switch name[len(name)-1] {
	case '$':
		return 1
	case '?':
		return 2
	default:
		return 0
	}
}

func (g *CodeGen) resolveOrDefineAssignTarget(name string) *symtable.Symbol {
	if g.fnDepth == 0 {
		return g.Symbols.Resolve(name)
	}
	s := g.Symbols.Resolve(name)
	if s != nil {
		return s
	}
	return g.Symbols.DefineLocal(name)
}

// Compile translates a full program AST into an opcode.Program.
func (g *CodeGen) Compile(tree *ast.Program) (*opcode.Program, error) {
	if tree == nil {
		return nil, fmt.Errorf("codegen: nil program")
	}

	// 1. Pre-declare functions so we can resolve OpCallUser vs OpCallBuiltin
	for _, fn := range tree.Functions {
		g.Prog.Functions[fn.Name] = opcode.NewChunk(fn.Name)
	}

	// User-defined TYPE metadata (needed when emitting main: constructors use Prog.Types)
	for _, t := range tree.Types {
		g.Prog.Types[t.Name] = &opcode.TypeDef{
			Name:   t.Name,
			Fields: t.Fields,
		}
	}

	// 2. Main program
	g.loopStack = nil
	for _, st := range tree.Stmts {
		g.emitStmt(g.Prog.Main, st)
		if g.err != nil {
			return nil, g.err
		}
	}
	g.Prog.Main.Emit(opcode.OpHalt, 0, 0, 0)

	// 3. Function bodies
	for _, fn := range tree.Functions {
		ch := g.Prog.Functions[fn.Name]
		g.loopStack = nil

		g.fnDepth = 1
		g.funcName = strings.ToUpper(fn.Name)
		g.Symbols.PushScope()
		// Define parameters in the local scope
		for _, p := range fn.Params {
			g.Symbols.DefineParam(p.Name)
		}

		g.predeclareFunctionBody(fn.Body)

		for _, st := range fn.Body {
			g.emitStmt(ch, st)
			if g.err != nil {
				return nil, g.err
			}
		}

		g.Symbols.PopScope()
		g.fnDepth = 0
		g.funcName = ""
		ch.Emit(opcode.OpReturnVoid, 0, 0, fn.Line)
	}

	opt.OptimizeProgram(g.Prog)
	g.Prog.SourcePath = g.File
	return g.Prog, nil
}

func (g *CodeGen) predeclareFunctionBody(stmts []ast.Stmt) {
	for _, s := range stmts {
		g.predeclareStmt(s)
	}
}

func (g *CodeGen) predeclareStmt(s ast.Stmt) {
	if s == nil {
		return
	}
	switch n := s.(type) {
	case *ast.LocalDeclNode:
		g.Symbols.DefineLocal(n.Name)
	case *ast.StaticDeclNode:
		g.Symbols.DefineStatic(g.funcName, n.Name)
	case *ast.DimNode:
		g.Symbols.DefineLocal(n.Name)
	case *ast.IfNode:
		for _, st := range n.Then {
			g.predeclareStmt(st)
		}
		for _, c := range n.ElseIf {
			for _, st := range c.Body {
				g.predeclareStmt(st)
			}
		}
		for _, st := range n.Else {
			g.predeclareStmt(st)
		}
	case *ast.WhileNode:
		for _, st := range n.Body {
			g.predeclareStmt(st)
		}
	case *ast.ForNode:
		// FOR var implicitly declares var locally if inside a function
		if g.fnDepth > 0 {
			g.Symbols.DefineLocal(n.Var)
		}
		for _, st := range n.Body {
			g.predeclareStmt(st)
		}
	case *ast.RepeatNode:
		for _, st := range n.Body {
			g.predeclareStmt(st)
		}
	case *ast.DoLoopNode:
		for _, st := range n.Body {
			g.predeclareStmt(st)
		}
	case *ast.SelectNode:
		for _, c := range n.Cases {
			for _, st := range c.Body {
				g.predeclareStmt(st)
			}
		}
		for _, st := range n.Default {
			g.predeclareStmt(st)
		}
	}
}
