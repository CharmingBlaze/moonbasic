# Audio Commands

Commands for loading, playing, and controlling sounds, music, audio streams, and raw wave data. moonBASIC audio is backed by Raylib's audio system and supports spatial 3D sound, pitch variation, and streaming.

## Core Concepts

- **Sound** — Short audio clip loaded fully into memory. Best for effects (gunshots, footsteps, UI clicks). Loaded with `Audio.LoadSound`.
- **Music** — Streamed from disk, decoded in chunks. Best for background music and long tracks. Loaded with `Audio.LoadMusic`.
- **AudioStream** — Raw PCM audio buffer for procedural/generated audio.
- **Wave** — Raw waveform data in memory for offline processing (crop, format, export).
- **Audio.Init()** — Must be called before any audio commands if the window module hasn't initialized audio automatically.

---

## Initialization

### `Audio.Init()`

Initializes the audio device. Called automatically when `Window.Open` is used, but can be called manually for console-mode programs that need audio.

**How it works:** Opens the default audio output device via Raylib's `InitAudioDevice()`. Safe to call multiple times — subsequent calls are no-ops.

```basic
Audio.Init()
```

---

### `Audio.Close()`

Closes the audio device and frees all audio resources.

```basic
Audio.Close()
```

---

## Sound Loading & Playback

### `Audio.LoadSound(filePath)`

Loads a sound file into memory. Supports `.wav`, `.ogg`, `.mp3`, `.flac`. Returns a handle.

- `filePath` (string) — Path to the audio file.

**Returns:** `handle`

**How it works:** The entire file is decoded and stored in a memory buffer. This means playback is instant with zero latency, but large files consume significant RAM. Use `Music` for long tracks.

```basic
shootSound = Audio.LoadSound("assets/shoot.wav")
```

---

### `Audio.Play(handle)`

Plays a sound or music handle. For sounds, plays from the beginning. For music, starts or resumes streaming.

- `handle` (handle) — Sound or music handle.

**How it works:** Internally detects whether the handle is a Sound or Music and dispatches to the correct Raylib playback function.

```basic
Audio.Play(shootSound)
```

---

### `Audio.Stop(handle)`

Stops playback of a sound or music handle.

- `handle` (handle) — Sound or music handle.

---

### `Audio.Pause(handle)`

Pauses playback. Use `Audio.Resume` to continue.

- `handle` (handle) — Sound or music handle.

---

### `Audio.Resume(handle)`

Resumes playback from where it was paused.

- `handle` (handle) — Sound or music handle.

---

### `Sound.Free(handle)`

Frees a sound from memory.

- `handle` (handle) — Sound handle.

```basic
Sound.Free(shootSound)
```

---

## Sound Properties

### `Audio.SetSoundVolume(soundHandle, volume)`

Sets the volume of a specific sound.

- `soundHandle` (handle) — Sound handle.
- `volume` (float) — Volume from 0.0 (silent) to 1.0 (full).

```basic
Audio.SetSoundVolume(shootSound, 0.5)
```

---

### `Audio.SetSoundPitch(soundHandle, pitch)`

Sets the playback pitch/speed of a sound. 1.0 = normal, 2.0 = double speed (one octave up), 0.5 = half speed.

- `soundHandle` (handle) — Sound handle.
- `pitch` (float) — Pitch multiplier.

---

### `Audio.SetSoundPan(soundHandle, pan)`

Sets the stereo panning of a sound.

- `soundHandle` (handle) — Sound handle.
- `pan` (float) — Pan from 0.0 (full left) to 1.0 (full right). 0.5 = center.

---

### `Audio.IsSoundPlaying(soundHandle)`

Returns `TRUE` if the sound is currently playing.

**Returns:** `bool`

---

## Sound Variation

### `Audio.PlayVarySound(soundHandle, pitchVariance)`

Plays a sound with a random pitch offset for natural variation. Each call sounds slightly different.

- `soundHandle` (handle) — Sound handle.
- `pitchVariance` (float) — Maximum pitch deviation (e.g., 0.1 for ±10%).

**How it works:** Temporarily adjusts the sound's pitch by a random amount within the variance range, plays it, then restores the original pitch.

```basic
; Footstep sounds with natural variation
Audio.PlayVarySound(footstepSound, 0.15)
```

---

### `Audio.PlayRndSound(soundHandle1, soundHandle2, ...)`

Randomly plays one of the provided sound handles. Useful for randomized effect banks.

---

## Music

### `Audio.LoadMusic(filePath)`

Loads a music file for streaming playback. Supports `.ogg`, `.mp3`, `.wav`, `.flac`.

- `filePath` (string) — Path to the music file.

**Returns:** `handle`

**How it works:** Opens the file for streaming — only a small buffer is decoded at a time. You **must** call `Audio.UpdateMusic` every frame to keep the stream fed.

```basic
bgMusic = Audio.LoadMusic("assets/theme.ogg")
Audio.Play(bgMusic)
```

---

### `Audio.UpdateMusic(musicHandle)`

Updates the music stream buffer. **Must be called every frame** while music is playing, or playback will stutter and stop.

- `musicHandle` (handle) — Music handle.

```basic
WHILE NOT Window.ShouldClose()
    Audio.UpdateMusic(bgMusic)
    ; ... game loop ...
WEND
```

---

### `Audio.SeekMusic(musicHandle, position)`

Seeks to a position in the music stream.

- `musicHandle` (handle) — Music handle.
- `position` (float) — Position in seconds.

---

### `Audio.SetMusicVolume(musicHandle, volume)`

Sets the volume of a music stream.

- `musicHandle` (handle) — Music handle.
- `volume` (float) — Volume 0.0–1.0.

---

### `Audio.SetMusicPitch(musicHandle, pitch)`

Sets the playback pitch of a music stream.

- `musicHandle` (handle) — Music handle.
- `pitch` (float) — Pitch multiplier.

---

### `Audio.IsMusicPlaying(musicHandle)`

Returns `TRUE` if the music stream is currently playing.

**Returns:** `bool`

---

### `Audio.GetMusicLength(musicHandle)`

Returns the total length of the music stream in seconds.

**Returns:** `float`

---

### `Audio.GetMusicTime(musicHandle)`

Returns the current playback position in seconds.

**Returns:** `float`

```basic
progress = Audio.GetMusicTime(bgMusic) / Audio.GetMusicLength(bgMusic)
Draw.Rectangle(10, 700, progress * 400, 10, 100, 255, 100, 255)
```

---

### `Music.Free(musicHandle)`

Frees a music stream.

- `musicHandle` (handle) — Music handle.

---

## Master Volume

### `Audio.SetMasterVolume(volume)`

Sets the global master volume that scales all audio output.

- `volume` (float) — Volume 0.0–1.0.

```basic
Audio.SetMasterVolume(0.8)
```

---

## 3D Spatial Audio

### `Audio.ListenerCamera(cameraHandle)`

Sets the audio listener position and orientation to match a camera. All 3D sounds will be spatialized relative to this camera.

- `cameraHandle` (handle) — Camera handle.

**How it works:** Updates the Raylib audio listener position and forward vector to match the camera's world position and look direction.

```basic
Audio.ListenerCamera(cam)
```

---

### `Sound.Play3D(soundHandle, x, y, z)`

Plays a sound at a 3D world position. Volume and panning are automatically adjusted based on the listener camera position.

- `soundHandle` (handle) — Sound handle.
- `x`, `y`, `z` (float) — World position.

```basic
; Explosion at world position
Sound.Play3D(explosionSound, 10, 0, -5)
```

---

### `Sound.Attach(soundHandle, entityHandle)`

Attaches a sound to an entity so it follows the entity's position automatically.

- `soundHandle` (handle) — Sound handle.
- `entityHandle` (handle) — Entity handle.

---

### `World.SetAmbience(soundHandle, volume)`

Sets an ambient background sound that plays continuously.

---

### `World.SetReverb(reverbType)`

Sets the reverb environment type for spatial audio.

---

## Audio Streams

Audio streams let you push raw PCM samples for procedural audio generation.

### `AudioStream.Create(sampleRate, sampleSize, channels)`

Creates a raw audio stream buffer.

- `sampleRate` (int) — Sample rate (e.g., 44100).
- `sampleSize` (int) — Bits per sample (8, 16, 32).
- `channels` (int) — Number of channels (1 = mono, 2 = stereo).

**Returns:** `handle`

---

### `AudioStream.Update(streamHandle, data)`

Pushes audio data into the stream buffer.

---

### `AudioStream.IsReady(streamHandle)` / `AudioStream.IsPlaying(streamHandle)`

Query stream state.

**Returns:** `bool`

---

### `AudioStream.Play(streamHandle)` / `AudioStream.Pause(streamHandle)` / `AudioStream.Resume(streamHandle)` / `AudioStream.Stop(streamHandle)`

Control stream playback.

---

### `AudioStream.SetVolume(streamHandle, volume)` / `AudioStream.SetPitch(streamHandle, pitch)` / `AudioStream.SetPan(streamHandle, pan)`

Set stream properties.

---

### `AudioStream.Free(streamHandle)`

Frees an audio stream.

---

## Wave Commands

Commands for loading and manipulating raw wave audio data in memory.

### `Wave.Load(filePath)`

Loads a wave file (`.wav`) into memory. Returns a handle to the wave data.

- `filePath` (string) — The path to the `.wav` file.

**Returns:** `handle`

---

### `Wave.Copy(waveHandle)`

Creates a new, independent copy of a wave resource.

- `waveHandle` (handle) — Source wave.

**Returns:** `handle`

---

### `Wave.Crop(waveHandle, startFrame, endFrame)`

Crops the wave data to a new length, from `startFrame` to `endFrame`.

- `waveHandle` (handle) — Wave to crop.
- `startFrame` (int) — First frame to keep.
- `endFrame` (int) — Last frame to keep.

---

### `Wave.Format(waveHandle, sampleRate, sampleSize, channels)`

Converts the wave data to a new format.

- `waveHandle` (handle) — Wave to convert.
- `sampleRate` (int) — Target sample rate.
- `sampleSize` (int) — Target bits per sample.
- `channels` (int) — Target channel count.

---

### `Wave.Export(waveHandle, filePath)`

Saves the wave data to a new `.wav` file.

- `waveHandle` (handle) — Wave to export.
- `filePath` (string) — Output file path.

---

### `Wave.Free(waveHandle)`

Frees the wave data from memory.

- `waveHandle` (handle) — Wave to free.

---

### `Sound.FromWave(waveHandle)`

Creates a playable `SOUND` handle from raw `WAVE` data. The resulting sound can be used with `Audio.Play()`.

- `waveHandle` (handle) — Source wave data.

**Returns:** `handle`

---

## Easy Mode Shortcuts

| Shortcut | Maps To |
|----------|---------|
| `LoadSound(path)` | `Audio.LoadSound(path)` |
| `LOADSOUND(path)` | `Audio.LoadSound(path)` |
| `PlaySound(h)` | `Audio.Play(h)` |
| `PLAYSOUND(h)` | `Audio.Play(h)` |
| `FreeSound(h)` | `Sound.Free(h)` |
| `FREESOUND(h)` | `Sound.Free(h)` |
| `SoundVolume(h, v)` | `Audio.SetSoundVolume(h, v)` |
| `SoundPitch(h, p)` | `Audio.SetSoundPitch(h, p)` |
| `Listener(cam)` | `Audio.ListenerCamera(cam)` |

---

## Full Example

A complete audio demo with sound effects, music streaming, and 3D spatial audio.

```basic
Window.Open(1280, 720, "Audio Demo")
Window.SetFPS(60)

; Audio initializes automatically with Window.Open

; Load assets
shootSound = Audio.LoadSound("assets/shoot.wav")
bgMusic = Audio.LoadMusic("assets/theme.ogg")

; Set up 3D listener
cam = Camera.Create()
cam.pos(0, 5, 10)
cam.look(0, 0, 0)
Audio.ListenerCamera(cam)

; Start music
Audio.SetMusicVolume(bgMusic, 0.5)
Audio.Play(bgMusic)

WHILE NOT Window.ShouldClose()
    ; IMPORTANT: Update music stream every frame
    Audio.UpdateMusic(bgMusic)

    ; Shoot on click with pitch variation
    IF Input.MousePressed(0) THEN
        Audio.PlayVarySound(shootSound, 0.1)
    ENDIF

    ; Render
    Render.Clear(20, 20, 40)
    Render.Begin3D(cam)
        Draw.Grid(20, 1.0)
    Render.End3D()

    ; Show music progress
    IF Audio.IsMusicPlaying(bgMusic) THEN
        progress = Audio.GetMusicTime(bgMusic) / Audio.GetMusicLength(bgMusic)
        Draw.Rectangle(10, 700, INT(progress * 400), 8, 100, 255, 100, 255)
        Draw.Text("Music: " + STR(INT(Audio.GetMusicTime(bgMusic))) + "s", 10, 680, 16, 200, 200, 200, 255)
    ENDIF

    Draw.Text("Click to shoot | Music is streaming", 10, 10, 18, 255, 255, 255, 255)
    Render.Frame()
WEND

; Cleanup
Sound.Free(shootSound)
Music.Free(bgMusic)
Camera.Free(cam)
Window.Close()
```

---

## See Also

- [WINDOW](WINDOW.md) — Audio initializes with `Window.Open`
- [CAMERA](CAMERA.md) — Used for spatial audio listener
- [ENTITY](ENTITY.md) — Attach sounds to entities
