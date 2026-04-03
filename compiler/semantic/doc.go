// Package semantic implements Phase 1 of the moonBASIC compiler pipeline: after the
// AST is built, it runs constant folding and static type checks for built-in
// namespace commands using data from builtinmanifest (no runtime import).
package semantic
