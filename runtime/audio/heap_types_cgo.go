//go:build cgo || (windows && !cgo)

package mbaudio

import (
	rl "github.com/gen2brain/raylib-go/raylib"

	"moonbasic/vm/heap"
)

type audioStreamObj struct {
	s       rl.AudioStream
	volume  float32 // last AUDIOSTREAM.SETVOLUME (default 1; raylib has no GetAudioStreamVolume)
	pitch   float32 // last AUDIOSTREAM.SETPITCH (default 1)
	pan     float32 // last AUDIOSTREAM.SETPAN (default 0)
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
	snd       rl.Sound
	gain      float32 // last AUDIO.SETSOUNDVOLUME (default 1)
	pan       float32 // last AUDIO.SETSOUNDPAN (default 0)
	spatial3D bool    // Load3DSound marks true (same buffer; spatial mix in EmitSound)
	release   heap.ReleaseOnce
}

func (o *soundObj) TypeName() string { return "Sound" }

func (o *soundObj) TypeTag() uint16 { return heap.TagSound }

func (o *soundObj) Free() {
	o.release.Do(func() { rl.UnloadSound(o.snd) })
}

type musicObj struct {
	m       rl.Music
	pitch   float32 // last AUDIO.SETMUSICPITCH (default 1; raylib has no GetMusicPitch)
	release heap.ReleaseOnce
}

func (o *musicObj) TypeName() string { return "Music" }

func (o *musicObj) TypeTag() uint16 { return heap.TagMusic }

func (o *musicObj) Free() {
	o.release.Do(func() { rl.UnloadMusicStream(o.m) })
}
