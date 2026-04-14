//go:build !cgo && !windows

package input

import (
	"moonbasic/vm/value"
)

func (m *Module) inputAxis(args []value.Value) (value.Value, error) {
	return value.FromFloat(0), nil
}

func (m *Module) inputAxisDeg(args []value.Value) (value.Value, error) {
	return value.FromFloat(0), nil
}

func (m *Module) axisX(args []value.Value) (value.Value, error) { return value.FromFloat(0), nil }
func (m *Module) axisY(args []value.Value) (value.Value, error) { return value.FromFloat(0), nil }
func (m *Module) axisDX(args []value.Value) (value.Value, error) { return value.FromFloat(0), nil }
func (m *Module) axisDY(args []value.Value) (value.Value, error) { return value.FromFloat(0), nil }
func (m *Module) axisDPadX(args []value.Value) (value.Value, error) { return value.FromFloat(0), nil }
func (m *Module) axisDPadY(args []value.Value) (value.Value, error) { return value.FromFloat(0), nil }
