package include

import (
	"fmt"
	"os"
	"path/filepath"

	"moonbasic/compiler/arena"
	"moonbasic/compiler/ast"
	"moonbasic/compiler/parser"
)

// Expand processes INCLUDE directives: merges types, functions, and statement lists;
// resolves paths relative to the including file; detects cycles and duplicate includes.
func Expand(hostFile string, prog *ast.Program) (*ast.Program, error) {
	return ExpandWithArena(hostFile, prog, nil)
}

// ExpandWithArena is like Expand; when ar is non-nil, included files parse into the same arena.
func ExpandWithArena(hostFile string, prog *ast.Program, ar *arena.Arena) (*ast.Program, error) {
	hostAbs, err := filepath.Abs(hostFile)
	if err != nil {
		return nil, err
	}
	seen := make(map[string]struct{})
	return expandProgram(hostAbs, prog, nil, seen, ar)
}

func expandProgram(hostAbs string, prog *ast.Program, stack []string, seen map[string]struct{}, ar *arena.Arena) (*ast.Program, error) {
	if stackContains(stack, hostAbs) {
		return nil, fmt.Errorf("[moonBASIC] Error: circular INCLUDE detected: %s", formatCircularChain(stack, hostAbs))
	}
	stack = append(stack, hostAbs)

	var prefTypes []*ast.TypeDef
	var prefFuncs []*ast.FunctionDef

	stmts, err := expandStmtSlice(hostAbs, prog.Stmts, stack, seen, &prefTypes, &prefFuncs, ar)
	if err != nil {
		return nil, err
	}
	prog.Stmts = stmts

	for i := range prog.Functions {
		body, err := expandStmtSlice(hostAbs, prog.Functions[i].Body, stack, seen, &prefTypes, &prefFuncs, ar)
		if err != nil {
			return nil, err
		}
		prog.Functions[i].Body = body
	}

	prog.Types = append(prefTypes, prog.Types...)
	prog.Functions = append(prefFuncs, prog.Functions...)
	return prog, nil
}

func expandStmtSlice(hostAbs string, stmts []ast.Stmt, stack []string, seen map[string]struct{}, prefT *[]*ast.TypeDef, prefF *[]*ast.FunctionDef, ar *arena.Arena) ([]ast.Stmt, error) {
	var out []ast.Stmt
	for _, s := range stmts {
		inc, ok := s.(*ast.IncludeNode)
		if !ok {
			out = append(out, s)
			continue
		}
		childPath, err := Resolve(hostAbs, inc.Path)
		if err != nil {
			return nil, err
		}
		childAbs, err := filepath.Abs(childPath)
		if err != nil {
			return nil, err
		}

		data, err := os.ReadFile(childPath)
		childAbsFinal := childAbs
		if err != nil {
			altPath, altData, altErr := TryOpenInclude(inc.Path)
			if altErr != nil {
				return nil, fmt.Errorf("INCLUDE %q: %w", inc.Path, err)
			}
			childPath = altPath
			childAbsFinal, err = filepath.Abs(altPath)
			if err != nil {
				return nil, err
			}
			data = altData
		}
		if stackContains(stack, childAbsFinal) {
			return nil, fmt.Errorf("[moonBASIC] Error: circular INCLUDE detected: %s", formatCircularChain(stack, childAbsFinal))
		}
		if _, dup := seen[childAbsFinal]; dup {
			return nil, fmt.Errorf("[moonBASIC] Error: duplicate INCLUDE of %q", inc.Path)
		}
		seen[childAbsFinal] = struct{}{}

		sub, err := parser.ParseSourceWithArena(childPath, string(data), ar)
		if err != nil {
			return nil, err
		}
		sub, err = expandProgram(childAbsFinal, sub, stack, seen, ar)
		if err != nil {
			return nil, err
		}
		*prefT = append(*prefT, sub.Types...)
		*prefF = append(*prefF, sub.Functions...)
		out = append(out, sub.Stmts...)
	}
	return out, nil
}
