//go:build cgo || (windows && !cgo)

package mbaudio

import (
	"fmt"
	"math/rand"

	rl "github.com/gen2brain/raylib-go/raylib"

	"moonbasic/vm/heap"
	"moonbasic/vm/value"
)

// variety playback implementation below

func (m *Module) audioPlayVarySound(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) != 3 {
		return value.Nil, fmt.Errorf("AUDIO.PLAYVARYSOUND expects (sound, minPitch#, maxPitch#)")
	}
	if args[0].Kind != value.KindHandle {
		return value.Nil, fmt.Errorf("AUDIO.PLAYVARYSOUND: sound handle required")
	}
	initAudioOnce()
	hid := heap.Handle(args[0].IVal)
	so, err := heap.Cast[*soundObj](m.h, hid)
	if err != nil {
		return value.Nil, fmt.Errorf("AUDIO.PLAYVARYSOUND: expected sound handle")
	}
	lo, ok1 := args[1].ToFloat()
	hi, ok2 := args[2].ToFloat()
	if !ok1 {
		if i, ok := args[1].ToInt(); ok {
			lo = float64(i)
			ok1 = true
		}
	}
	if !ok2 {
		if i, ok := args[2].ToInt(); ok {
			hi = float64(i)
			ok2 = true
		}
	}
	if !ok1 || !ok2 {
		return value.Nil, fmt.Errorf("AUDIO.PLAYVARYSOUND: pitch range must be numeric")
	}
	if hi < lo {
		lo, hi = hi, lo
	}
	p := lo
	if hi > lo {
		p = lo + rand.Float64()*(hi-lo)
	}
	rl.SetSoundPitch(so.snd, float32(p))
	rl.PlaySound(so.snd)
	return value.Nil, nil
}

func (m *Module) audioPlayRndSound(args []value.Value) (value.Value, error) {
	if err := m.requireHeap(); err != nil {
		return value.Nil, err
	}
	if len(args) < 1 {
		return value.Nil, fmt.Errorf("AUDIO.PLAYRNDSOUND expects at least one sound handle")
	}
	for i := range args {
		if args[i].Kind != value.KindHandle {
			return value.Nil, fmt.Errorf("AUDIO.PLAYRNDSOUND: all arguments must be sound handles")
		}
		if _, err := heap.Cast[*soundObj](m.h, heap.Handle(args[i].IVal)); err != nil {
			return value.Nil, fmt.Errorf("AUDIO.PLAYRNDSOUND: expected sound handles")
		}
	}
	initAudioOnce()
	ix := rand.Intn(len(args))
	so, _ := heap.Cast[*soundObj](m.h, heap.Handle(args[ix].IVal))
	rl.PlaySound(so.snd)
	return value.Nil, nil
}
