// Package mbutil implements UTIL.* path, file, and (when CGO) drag-drop helpers.
package mbutil

// Module registers UTIL builtins.
type Module struct{}

// NewModule creates the module.
func NewModule() *Module { return &Module{} }

func (m *Module) Reset() {}


