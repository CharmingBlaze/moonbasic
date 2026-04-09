//go:build !cgo && !windows

package mblight

import "moonbasic/runtime"

func (m *Module) registerPointLightBlitz(_ runtime.Registrar) {}
