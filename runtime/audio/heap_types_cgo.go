//go:build cgo

package mbaudio

import (
	rl "github.com/gen2brain/raylib-go/raylib"

	"moonbasic/vm/heap"
)

type audioStreamObj struct {
	s       rl.AudioStream
	release heap.ReleaseOnce
}

func (o *audioStreamObj) TypeName() string { return "AudioStream" }

func (o *audioStreamObj) TypeTag() uint16 { return heap.TagAudioStream }

func (o *audioStreamObj) Free() {
	o.release.Do(func() { rl.UnloadAudioStream(o.s) })
}

type waveObj struct {
	w       rl.Wave
	release heap.ReleaseOnce
}

func (o *waveObj) TypeName() string { return "Wave" }

func (o *waveObj) TypeTag() uint16 { return heap.TagWave }

func (o *waveObj) Free() {
	o.release.Do(func() { rl.UnloadWave(o.w) })
}

type soundObj struct {
	snd     rl.Sound
	release heap.ReleaseOnce
}

func (o *soundObj) TypeName() string { return "Sound" }

func (o *soundObj) TypeTag() uint16 { return heap.TagSound }

func (o *soundObj) Free() {
	o.release.Do(func() { rl.UnloadSound(o.snd) })
}

type musicObj struct {
	m       rl.Music
	release heap.ReleaseOnce
}

func (o *musicObj) TypeName() string { return "Music" }

func (o *musicObj) TypeTag() uint16 { return heap.TagMusic }

func (o *musicObj) Free() {
	o.release.Do(func() { rl.UnloadMusicStream(o.m) })
}
