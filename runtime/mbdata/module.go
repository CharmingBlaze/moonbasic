// Package mbdata implements DATA.* compression, encoding, and hashing (pure Go).
package mbdata

// Module registers DATA.* builtins.
type Module struct{}

// NewModule creates the module.
func NewModule() *Module { return &Module{} }

func (m *Module) Reset() {}


