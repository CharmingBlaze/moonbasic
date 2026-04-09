//go:build !cgo && !windows

package mbfile

import "moonbasic/runtime"

func (m *Module) registerBankTransfer(r runtime.Registrar) {
	_ = m
	_ = r
}
