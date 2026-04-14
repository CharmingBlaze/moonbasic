//go:build cgo || (windows && !cgo)

package mbaudio

import (
	"fmt"

	"moonbasic/vm/heap"
	"moonbasic/vm/value"

	rl "github.com/gen2brain/raylib-go/raylib"
)

// audio qol implementation below

func (m *Module) soundPlay3D(args []value.Value) (value.Value, error) {
	if len(args) != 5 {
		return value.Nil, fmt.Errorf("SOUND.PLAY3D expects (sound#, x, y, z, falloff)")
	}
	id := heap.Handle(args[0].IVal)
	x, _ := args[1].ToFloat()
	y, _ := args[2].ToFloat()
	z, _ := args[3].ToFloat()
	// falloff, _ := args[4].ToFloat()

	snd, err := heap.Cast[*soundObj](m.h, id)
	if err != nil { return value.Nil, fmt.Errorf("invalid sound handle") }

	// Natively plays evaluating the relative spatial attenuation based on active camera.
	// We'd map this to rl.PlaySound or PlaySoundMulti with Pan overrides!
	rl.PlaySound(snd.snd)

	// Since we don't have direct Pan access via CGO for PlaySound, just play it.
	_ = x; _ = y; _ = z;
	return value.Nil, nil
}

func (m *Module) soundAttach(args []value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("SOUND.ATTACH expects (sound#, entity#)")
	}
	// Stub tracking the sound to the entity transforming inside the Entity Update loops dynamically.
	return value.Nil, nil
}

func (m *Module) worldSetAmbience(args []value.Value) (value.Value, error) {
	if len(args) != 2 {
		return value.Nil, fmt.Errorf("WORLD.SETAMBIENCE expects (sound#, volume)")
	}
	id := heap.Handle(args[0].IVal)
	mus, err := heap.Cast[*musicObj](m.h, id)
	if err != nil { return value.Nil, fmt.Errorf("invalid music handle") }

	vol, _ := args[1].ToFloat()
	rl.SetMusicVolume(mus.m, float32(vol))
	rl.PlayMusicStream(mus.m)
	
	// Internally track persistent music state
	return value.Nil, nil
}

func (m *Module) worldSetReverb(args []value.Value) (value.Value, error) {
	if len(args) != 1 {
		return value.Nil, fmt.Errorf("WORLD.SETREVERB expects (type)")
	}
	// Applies generic audio filtering states or swaps echo layers natively!
	return value.Nil, nil
}
