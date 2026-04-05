//go:build !cgo

package mbaudio

import (
	"fmt"

	"moonbasic/runtime"
	"moonbasic/vm/value"
)

const hint = "AUDIO.* requires CGO: set CGO_ENABLED=1 and install a C compiler, then rebuild"

func raylibAudioOpen()  {}
func raylibAudioClose() {}

// Register implements runtime.Module.
func (m *Module) Register(r runtime.Registrar) {
	stub := func(name string) runtime.BuiltinFn {
		return func(rt *runtime.Runtime, args ...value.Value) (value.Value, error) {
			return value.Nil, fmt.Errorf("%s: %s", name, hint)
		}
	}
	r.Register("AUDIO.INIT", "audio", stub("AUDIO.INIT"))
	r.Register("AUDIO.CLOSE", "audio", stub("AUDIO.CLOSE"))
	r.Register("AUDIO.LOADMUSIC", "audio", stub("AUDIO.LOADMUSIC"))
	r.Register("AUDIO.LOADSOUND", "audio", stub("AUDIO.LOADSOUND"))
	r.Register("AUDIO.PLAY", "audio", stub("AUDIO.PLAY"))
	r.Register("AUDIO.STOP", "audio", stub("AUDIO.STOP"))
	r.Register("AUDIO.PAUSE", "audio", stub("AUDIO.PAUSE"))
	r.Register("AUDIO.RESUME", "audio", stub("AUDIO.RESUME"))
	names := []string{
		"AUDIO.PLAYVARYSOUND", "AUDIO.PLAYRNDSOUND",
		"AUDIO.UPDATEMUSIC", "MUSIC.FREE",
		"AUDIO.SETSOUNDVOLUME", "AUDIO.SETSOUNDPITCH", "AUDIO.SETSOUNDPAN",
		"AUDIO.SETMUSICVOLUME", "AUDIO.SETMUSICPITCH", "AUDIO.SETMASTERVOLUME",
		"AUDIO.ISSOUNDPLAYING", "AUDIO.ISMUSICPLAYING",
		"AUDIO.GETMUSICLENGTH", "AUDIO.GETMUSICTIME", "AUDIO.SEEKMUSIC",
		"AUDIOSTREAM.MAKE", "AUDIOSTREAM.UPDATE", "AUDIOSTREAM.ISREADY", "AUDIOSTREAM.ISPLAYING",
		"AUDIOSTREAM.PLAY", "AUDIOSTREAM.PAUSE", "AUDIOSTREAM.RESUME", "AUDIOSTREAM.STOP",
		"AUDIOSTREAM.SETVOLUME", "AUDIOSTREAM.SETPITCH", "AUDIOSTREAM.SETPAN", "AUDIOSTREAM.FREE",
		"WAVE.LOAD", "WAVE.COPY", "WAVE.CROP", "WAVE.FORMAT", "WAVE.EXPORT", "WAVE.FREE",
		"SOUND.FROMWAVE", "SOUND.FREE",
	}
	for _, n := range names {
		r.Register(n, "audio", stub(n))
	}
}

// Shutdown implements runtime.Module.
func (m *Module) Shutdown() {}
