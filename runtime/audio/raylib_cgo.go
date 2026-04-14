//go:build cgo || (windows && !cgo)

package mbaudio

import (
	"moonbasic/runtime"
	"moonbasic/vm/value"
)

func raylibAudioOpen() { initAudioOnce() }

func raylibAudioClose() { closeAudioOnce() }

func (m *Module) audioInit(args []value.Value) (value.Value, error) {
	if len(args) != 0 {
		return value.Nil, runtime.Errorf("AUDIO.INIT expects 0 arguments")
	}
	initAudioOnce()
	return value.Nil, nil
}

func (m *Module) audioClose(args []value.Value) (value.Value, error) {
	if len(args) != 0 {
		return value.Nil, runtime.Errorf("AUDIO.CLOSE expects 0 arguments")
	}
	closeAudioOnce()
	return value.Nil, nil
}

// Shutdown implements runtime.Module.
func (m *Module) Shutdown() {}
