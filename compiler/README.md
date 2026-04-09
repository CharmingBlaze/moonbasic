# MoonBASIC Compiler Architecture

Welcome to the MoonBASIC compiler source code! This directory contains all the elements required to translate MoonBASIC source code into bytecode that the MooonBASIC Virtual Machine (VM) can execute.

If you're a new developer looking to contribute to the engine, modify language features, or understand how compilation works, this document is your starting point.

---

## The Compilation Pipeline

The journey of MoonBASIC source code to bytecode involves several distinct phases. This orchestration is primarily handled in `pipeline/compile.go`.

### 1. Lexical Analysis (`lexer/`)
The `Lexer` breaks raw source code text into a stream of `token.Token` structures. It turns keywords like `IF`, variable names like `playerX`, and strings into categorized tokens.

### 2. Parsing (`parser/`)
The parser takes the stream of tokens and constructs an **Abstract Syntax Tree (AST)**. This is a hierarchical tree representation of the code logic. Nodes for the AST are defined in `ast/`. 

### 3. Includes & Expansion (`include/`)
Any `INCLUDE "file.mb"` statements are resolved at the AST level. The included source is parsed and stitched into the primary AST.

### 4. Semantic Analysis & Symbol Tables (`symtable/`, `semantic/`)
Before generating code, the compiler checks the AST to make sure everything makes sense.
- **Symbol Table Builder (`symtable/`)**: Collects all functions, variables, and custom types so the compiler knows what exists.
- **Semantic Analyzer (`semantic/`)**: Validates type usage (e.g. checking you aren't trying to do math on strings unless doing concatenation), verifies variable scopes, checks argument counts in function calls, and resolves implicit declarations.

### 5. Optimization (Optional - `opt/`)
Basic constant folding or other minor optimizations happen at the AST level before lowering.

### 6. Code Generation (`codegen/`)
The final compiler phase takes the validated AST and walks through it, generating VM bytecode instructions (e.g. `OP_LOAD`, `OP_ADD`, `OP_CALL`). The result is an `opcode.Program` struct containing the binary instructions.

**Builtin calls (`OpCallBuiltin`):** Statement and expression calls to built-ins (e.g. `MoveEntity(...)`, `DrawEntities()`) evaluate each argument into consecutive temporary registers starting at a per-statement base; the opcode records the first argument register and packs `(argCount << 24) | nameIndex` so the VM can invoke the registered handler with that register window. See [`codegen/codegen_calls.go`](codegen/codegen_calls.go) and [`codegen/codegen_expr.go`](codegen/codegen_expr.go) (`CallExprNode`). Zero-argument builtins may be written as `Name` or `Name()` when `Name` is listed in [`builtinmanifest`](builtinmanifest/) with arity 0 ([`parser_stmts.go`](parser/parser_stmts.go) `parseStmtAfterIdent`).

---

## Core Tenets & Subsystems

- **`intern` (String Interning):** To save memory and increase speed, variable names and identifiers are interned using the `intern/` package. Identical strings point to the same memory location.
- **`arena` (Memory Pooling):** To reduce garbage collection overhead during rapid parsing passes, temporary objects might be allocated using an arena.
- **`types` (Type System):** The representations for MoonBASIC variable types (Int, Float, String, Function, Handle, Array) are standardized here.
- **`errors` (Standardized Errors):** Every compilation error uses `errors.CompileError` to ensure it reports a specific File, Line number, and helpful message.

---

## How to Add a New Built-in Command

If you want to add a new command like `DRAW_CIRCLE x, y, r`:

1. **Token (if it's a structural keyword like `FOR`):** Add it to `token/token.go` and `lexer/keywords.go`. *(Note: Most game commands are just plain identifiers / library functions, not structural tokens. You might not need to touch the lexer unless it's a core language primitive!)*
2. **Built-in Manifest (`builtinmanifest/`):** For standard library features (like math or rendering), MoonBASIC relies on a manifest to register the VM handlers and signatures. Ensure the command name exists.
3. **VM Registration:** Don't forget that after adding it to the compiler, the runtime VM (`vm/`) must have a corresponding Go function registered to handle the dispatched call.

## The AST Visitor Pattern

If you are writing a new compiler pass (such as a linter, structural analyzer, or optimizer), look at the `ast.Visitor` interface. You can create a struct that satisfies `Visitor` and call `astNode.Accept(visitor)`. This allows you to walk the tree without writing a massive `switch` block.
