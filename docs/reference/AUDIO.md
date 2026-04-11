# Audio Commands

Commands for loading and playing sound and music.

## Device Management

### `Audio.Init()`

Initializes the audio device. This must be called before any other audio commands.

### `Audio.Close()`

Closes the audio device and releases all audio resources.

---

## Sound Effects

Sounds are loaded completely into memory, making them fast to play. Ideal for short, repeatable effects like jumps or explosions.

### `Audio.LoadSound(path)`

Loads a sound effect from a file (e.g., `.wav`, `.ogg`). Returns a **sound handle**.

### Spatial / Blitz-style 3D helpers

- **`Listener(cameraHandle)`** / **`Audio.ListenerCamera()`** — sets the virtual listener from a **`Camera.Make()`** handle (position + horizontal forward). Call **each frame** before **`Entity.EmitSound()`** so pan and falloff stay correct.
- **`Load3DSound(path)`** — same buffers as **`Audio.LoadSound()`**; the “3D” path is for scripts that pair it with **`Listener`** + **`Entity.EmitSound()`**.
- **`Entity.EmitSound(sound, entity)`** — plays once with **quadratic distance falloff** (max distance ≈ 80 world units) and **stereo pan** from the horizontal angle to the source. Restores each sound’s last **`Audio.SetSoundVolume()`** / **`Audio.SetSoundPan()`** after the play call.
- **`Audio.SoundVolume()`** / **`Audio.SoundPitch()`** — aliases of **`Audio.SetSoundVolume()`** / **`Audio.SetSoundPitch()`**.

Raylib does not expose a full OpenAL-style HRTF; this is a lightweight **pan + attenuation** model.

### `Audio.Play(handle)`

Plays a loaded sound effect. Multiple instances of the same sound can overlap.

### `Audio.Stop(handle)`

Stops playback of a sound or music handle.

### `Audio.PlayVarySound(sound, minPitch, maxPitch)`

Picks a **uniform random** pitch between **`minPitch`** and **`maxPitch`**, applies it with **`Audio.SetSoundPitch()`**, then plays the sound. Pitch stays on the sound object until changed again.

### `Audio.PlayRndSound(sound1, sound2, ...)`

Plays **one** of the given sound handles, chosen uniformly at random (two to four overloads are listed in the manifest). All arguments must be **sound** handles.

### `Sound.Free(handle)`

Unloads a sound effect from memory and releases its heap slot.

---

## Music

Music is streamed from the file on disk, which uses less memory. Ideal for long background tracks.

### `Audio.LoadMusic(path)`

Loads a music file for streaming (e.g., `.mp3`, `.ogg`). Returns a **music handle**.

### `Audio.UpdateMusic(handle)`

Updates the music stream buffer. **Must be called every frame** for music to play correctly.

### `Audio.Play(handle)`

Starts playing a music track.

### `Audio.Stop(handle)` / `Audio.Pause(handle)` / `Audio.Resume(handle)`

Controls music playback.

### `Music.Free(handle)`

Unloads a music stream and releases its resources.

---

## Full Example

This example assumes you have `jump.wav` and `theme.mp3` in the same directory.

```basic
Window.Open(800, 600, "Audio Example")
Window.SetFPS(60)

; 1. Initialize audio
Audio.Init()

; 2. Load assets
jump_sfx = Audio.LoadSound("jump.wav")
bg_music = Audio.LoadMusic("theme.mp3")

; 3. Play music
Audio.Play(bg_music)

WHILE NOT Window.ShouldClose()
    ; 4. Update music stream every frame
    Audio.UpdateMusic(bg_music)

    ; Play sound on key press
    IF Input.KeyPressed(KEY_SPACE) THEN
        Audio.Play(jump_sfx)
    ENDIF

    Render.Clear(40, 40, 40)
    Camera2D.Begin()
        Draw.Text("Press SPACE to play a sound!", 190, 200, 20, 255, 255, 255, 255)
    Camera2D.End()
    Render.Frame()
WEND

; 5. Cleanup
Sound.Free(jump_sfx)
Music.Free(bg_music)
Audio.Close()
Window.Close()
```
