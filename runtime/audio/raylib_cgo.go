//go:build cgo || (windows && !cgo)

package mbaudio

import (
	"moonbasic/runtime"
	"moonbasic/vm/value"
)

func raylibAudioOpen() { initAudioOnce() }

func raylibAudioClose() { closeAudioOnce() }

// Register implements runtime.Module.
func (m *Module) Register(r runtime.Registrar) {
	r.Register("AUDIO.INIT", "audio", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if len(args) != 0 {
			return value.Nil, runtime.Errorf("AUDIO.INIT expects 0 arguments")
		}
		initAudioOnce()
		return value.Nil, nil
	}))
	r.Register("AUDIO.CLOSE", "audio", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		if len(args) != 0 {
			return value.Nil, runtime.Errorf("AUDIO.CLOSE expects 0 arguments")
		}
		closeAudioOnce()
		return value.Nil, nil
	}))
	m.registerSound(r)
	m.registerMusic(r)
	m.registerStreamWaveSound(r)
	m.registerPlayback(r)
	m.registerVarietyPlayback(r)
	m.registerAudioProps(r)
	m.registerSpatialAudio(r)
	registerAudioQoLAPI(m, r)
}

// Shutdown implements runtime.Module.
func (m *Module) Shutdown() {}
