//go:build cgo || (windows && !cgo)

package mbaudio

import (
	"fmt"

	rl "github.com/gen2brain/raylib-go/raylib"

	"moonbasic/runtime"
	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

func (m *Module) registerSound(r runtime.Registrar) {
	r.Register("AUDIO.LOADSOUND", "audio", m.soundLoad)

	// Global shorthands (Easy Mode)
	r.Register("LOADSOUND", "audio", m.soundLoad)
	r.Register("FREESOUND", "audio", func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
		if len(args) != 1 || args[0].Kind != value.KindHandle {
			return value.Nil, fmt.Errorf("FREESOUND expects 1 handle (sound)")
		}
		if rt != nil && rt.Heap != nil {
			rt.Heap.Free(heap.Handle(args[0].IVal))
		}
		return value.Nil, nil
	})
}

func (m *Module) soundLoad(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 1 || args[0].Kind != value.KindString {
		return value.Nil, fmt.Errorf("AUDIO.LOADSOUND expects 1 string path")
	}
	path, err := rt.ArgString(args, 0)
	if err != nil {
		return value.Nil, err
	}
	initAudioOnce()
	snd := rl.LoadSound(path)
	id, err := m.h.Alloc(&soundObj{snd: snd, gain: 1, pan: 0})
	if err != nil {
		rl.UnloadSound(snd)
		return value.Nil, err
	}
	return value.FromHandle(id), nil
}
