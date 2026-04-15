package mbaudio

import (
	"moonbasic/runtime"
	"moonbasic/vm/value"
)

func (m *Module) Register(r runtime.Registrar) {
	// Raylib backend / LifeCycle
	r.Register("AUDIO.INIT", "audio", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		return m.audioInit(args)
	}))
	r.Register("AUDIO.CLOSE", "audio", runtime.AdaptLegacy(func(args []value.Value) (value.Value, error) {
		return m.audioClose(args)
	}))

	// Playback (dispatch_cgo.go)
	r.Register("AUDIO.PLAY", "audio", runtime.AdaptLegacy(m.audioPlay))
	r.Register("AUDIO.STOP", "audio", runtime.AdaptLegacy(m.audioStop))
	r.Register("AUDIO.PAUSE", "audio", runtime.AdaptLegacy(m.audioPause))
	r.Register("AUDIO.RESUME", "audio", runtime.AdaptLegacy(m.audioResume))
	r.Register("PLAYSOUND", "audio", runtime.AdaptLegacy(m.audioPlay))
	r.Register("SOUNDVOLUME", "audio", runtime.AdaptLegacy(m.audioSetVolume))

	// Loading (sound_cgo.go, music_cgo.go)
	r.Register("AUDIO.LOADSOUND", "audio", m.soundLoad)
	r.Register("LOADSOUND", "audio", m.soundLoad)
	r.Register("FREESOUND", "audio", m.soundFree)
	r.Register("AUDIO.LOADMUSIC", "audio", m.musicLoad)

	// Props (props_cgo.go)
	r.Register("AUDIO.SETSOUNDVOLUME", "audio", runtime.AdaptLegacy(m.setSoundVolume))
	r.Register("AUDIO.SETSOUNDPITCH", "audio", runtime.AdaptLegacy(m.setSoundPitch))
	r.Register("AUDIO.SETSOUNDPAN", "audio", runtime.AdaptLegacy(m.setSoundPan))
	r.Register("AUDIO.SETMUSICVOLUME", "audio", runtime.AdaptLegacy(m.setMusicVolume))
	r.Register("AUDIO.SETMUSICPITCH", "audio", runtime.AdaptLegacy(m.setMusicPitch))
	r.Register("AUDIO.SETMASTERVOLUME", "audio", runtime.AdaptLegacy(m.setMasterVolume))
	r.Register("AUDIO.ISSOUNDPLAYING", "audio", runtime.AdaptLegacy(m.isSoundPlaying))
	r.Register("AUDIO.ISMUSICPLAYING", "audio", runtime.AdaptLegacy(m.isMusicPlaying))
	r.Register("AUDIO.GETMUSICLENGTH", "audio", runtime.AdaptLegacy(m.getMusicLength))
	r.Register("AUDIO.GETMUSICTIME", "audio", runtime.AdaptLegacy(m.getMusicTime))

	// Variety (variety_cgo.go)
	r.Register("AUDIO.PLAYVARYSOUND", "audio", runtime.AdaptLegacy(m.audioPlayVarySound))
	r.Register("AUDIO.PLAYRNDSOUND", "audio", runtime.AdaptLegacy(m.audioPlayRndSound))

	// QoL API (audio_qol_cgo.go)
	r.Register("SOUND.PLAY3D", "audio", runtime.AdaptLegacy(m.soundPlay3D))
	r.Register("SOUND.ATTACH", "audio", runtime.AdaptLegacy(m.soundAttach))
	r.Register("WORLD.SETAMBIENCE", "audio", runtime.AdaptLegacy(m.worldSetAmbience))
	r.Register("WORLD.SETREVERB", "audio", runtime.AdaptLegacy(m.worldSetReverb))

	// Spatial (spatial_cgo.go)
	r.Register("AUDIO.LISTENERCAMERA", "audio", m.audioListenerCamera)
	r.Register("Listener", "audio", m.audioListenerCamera)
	r.Register("Load3DSound", "audio", m.soundLoad3D)
	r.Register("SoundVolume", "audio", runtime.AdaptLegacy(m.setSoundVolume))
	r.Register("SoundPitch", "audio", runtime.AdaptLegacy(m.setSoundPitch))

	// Playback Extras
	r.Register("AUDIO.SEEKMUSIC", "audio", runtime.AdaptLegacy(m.seekMusic))
	r.Register("AUDIO.UPDATEMUSIC", "audio", runtime.AdaptLegacy(m.musicUpdate))
	r.Register("MUSIC.FREE", "audio", runtime.AdaptLegacy(m.musicFree))

	// Streams (stream_wave_cgo.go)
	r.Register("AUDIOSTREAM.CREATE", "audio", runtime.AdaptLegacy(m.streamMake))
	r.Register("AUDIOSTREAM.MAKE", "audio", runtime.AdaptLegacy(m.streamMake))
	r.Register("AUDIOSTREAM.UPDATE", "audio", runtime.AdaptLegacy(m.streamUpdate))
	r.Register("AUDIOSTREAM.ISREADY", "audio", runtime.AdaptLegacy(m.streamIsReady))
	r.Register("AUDIOSTREAM.ISPLAYING", "audio", runtime.AdaptLegacy(m.streamIsPlaying))
	r.Register("AUDIOSTREAM.PLAY", "audio", runtime.AdaptLegacy(m.streamPlay))
	r.Register("AUDIOSTREAM.PAUSE", "audio", runtime.AdaptLegacy(m.streamPause))
	r.Register("AUDIOSTREAM.RESUME", "audio", runtime.AdaptLegacy(m.streamResume))
	r.Register("AUDIOSTREAM.STOP", "audio", runtime.AdaptLegacy(m.streamStop))
	r.Register("AUDIOSTREAM.SETVOLUME", "audio", runtime.AdaptLegacy(m.streamSetVolume))
	r.Register("AUDIOSTREAM.SETPITCH", "audio", runtime.AdaptLegacy(m.streamSetPitch))
	r.Register("AUDIOSTREAM.SETPAN", "audio", runtime.AdaptLegacy(m.streamSetPan))
	r.Register("AUDIOSTREAM.FREE", "audio", runtime.AdaptLegacy(m.streamFree))

	// Waves (stream_wave_cgo.go)
	r.Register("WAVE.LOAD", "audio", m.waveLoad)
	r.Register("WAVE.COPY", "audio", runtime.AdaptLegacy(m.waveCopy))
	r.Register("WAVE.CROP", "audio", runtime.AdaptLegacy(m.waveCrop))
	r.Register("WAVE.FORMAT", "audio", runtime.AdaptLegacy(m.waveFormat))
	r.Register("WAVE.EXPORT", "audio", m.waveExport)
	r.Register("WAVE.FREE", "audio", runtime.AdaptLegacy(m.waveFree))

	// Compound
	r.Register("SOUND.FROMWAVE", "audio", runtime.AdaptLegacy(m.soundFromWave))
	r.Register("SOUND.FREE", "audio", m.soundFree)
}
