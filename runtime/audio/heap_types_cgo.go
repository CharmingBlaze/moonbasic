//go:build cgo

package mbaudio

import (
	rl "github.com/gen2brain/raylib-go/raylib"

	"moonbasic/vm/heap"
)

type audioStreamObj struct {
	s rl.AudioStream
}

func (o *audioStreamObj) TypeName() string { return "AudioStream" }

func (o *audioStreamObj) TypeTag() uint16 { return heap.TagAudioStream }

func (o *audioStreamObj) Free() {
	rl.UnloadAudioStream(o.s)
}

type waveObj struct {
	w rl.Wave
}

func (o *waveObj) TypeName() string { return "Wave" }

func (o *waveObj) TypeTag() uint16 { return heap.TagWave }

func (o *waveObj) Free() {
	rl.UnloadWave(o.w)
}

type soundObj struct {
	snd rl.Sound
}

func (o *soundObj) TypeName() string { return "Sound" }

func (o *soundObj) TypeTag() uint16 { return heap.TagSound }

func (o *soundObj) Free() {
	rl.UnloadSound(o.snd)
}
