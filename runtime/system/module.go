// Package mbsystem implements SYSTEM.* builtins.
package mbsystem

import (
	"os"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

// Module is a thin SYSTEM namespace module.
type Module struct{}

// NewModule creates SYSTEM builtins.
func NewModule() *Module { return &Module{} }

// Register implements runtime.Module.
// SYSTEM.EXIT calls os.Exit and may skip deferred cleanup in the host process.
func (m *Module) Register(reg runtime.Registrar) {
	reg.Register("SYSTEM.EXIT", "system", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		_ = rt
		code := int64(0)
		if len(args) > 1 {
			return value.Nil, runtime.Errorf("SYSTEM.EXIT expects 0 or 1 arguments, got %d", len(args))
		}
		if len(args) == 1 {
			code, _ = args[0].ToInt()
		}
		os.Exit(int(code))
		return value.Nil, nil
	})
	m.registerHost(reg)
}

// Shutdown implements runtime.Module.
func (m *Module) Shutdown() {}

func (m *Module) Reset() {}

