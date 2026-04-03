// Package runtime implements the moonBASIC native command layer.
package runtime

// InitCore registers the non-namespaced global moonBASIC commands (PRINT, string helpers).
func (r *Registry) InitCore() {
	registerConsoleIO(r)
	registerStringBuiltins(r)
	registerHostArgv(r)
}
