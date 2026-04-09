//go:build !cgo && !windows

package mbfile

import "moonbasic/runtime"

func (m *Module) registerFileBlitz(_ runtime.Registrar) {}
