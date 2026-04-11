# Audio Commands

Commands for loading and playing sound and music.

## Core Workflow

1.  **Initialize**: Call `Audio.Init()` once at the start of your program.
2.  **Load Assets**: Use `Audio.LoadSound()` for short effects and `Audio.LoadMusic()` for background tracks.
3.  **Control Playback**: Use `Audio.Play()`, `Audio.Stop()`, etc., to control the audio in your game logic.
4.  **Update Music**: For streaming music, you must call `Audio.UpdateMusic()` each frame.
5.  **Cleanup**: Unload sounds with `Sound.Free()` and `Music.Free()`, then call `Audio.Close()` before the program exits.

---

## Device Management

### `Audio.Init()`

Initializes the audio device. This must be called before any other audio commands.

### `Audio.Close()`

Closes the audio device and releases all audio resources.

---

## Sound Effects

Sounds are loaded completely into memory, making them fast to play. Ideal for short, repeatable effects like jumps or explosions.

### `Audio.LoadSound(filePath)`

Loads a sound effect from a file (e.g., `.wav`, `.ogg`). Returns a handle.

### Spatial / Blitz-style 3D helpers

- **`Listener(cameraHandle)`** / **`AUDIO.LISTENERCAMERA`** — sets the virtual listener from a **`CAMERA.MAKE`** handle (position + horizontal forward). Call **each frame** before **`EmitSound`** so pan and falloff stay correct.
- **`Load3DSound(path)`** — same buffers as **`AUDIO.LOADSOUND`**; the “3D” path is for scripts that pair it with **`Listener`** + **`EmitSound`**.
- **`EmitSound(sound, entity)`** (registered on **`ENTITY.*`**) — plays once with **quadratic distance falloff** (max distance ≈ 80 world units) and **stereo pan** from the horizontal angle to the source. Restores each sound’s last **`AUDIO.SETSOUNDVOLUME`** / **`AUDIO.SETSOUNDPAN`** after the play call.
- **`SoundVolume`** / **`SoundPitch`** — aliases of **`AUDIO.SETSOUNDVOLUME`** / **`AUDIO.SETSOUNDPITCH`**.

Raylib does not expose a full OpenAL-style HRTF; this is a lightweight **pan + attenuation** model.

### `Audio.Play(soundHandle)`

Plays a loaded sound effect. Multiple instances of the same sound can overlap.

### `Audio.PlayVarySound(sound, minPitch, maxPitch)`

Picks a **uniform random** pitch between **`minPitch`** and **`maxPitch`**, applies it with **`Audio.SetSoundPitch`**, then plays the sound. Pitch stays on the sound object until changed again.

### `Audio.PlayRndSound(sound1, sound2, ...)`

Plays **one** of the given sound handles, chosen uniformly at random (two to four overloads are listed in the manifest). All arguments must be **sound** handles.

### `Sound.Free(soundHandle)`

Unloads a sound from memory.

---

## Music

Music is streamed from the file on disk, which uses less memory. Ideal for long background tracks.

### `Audio.LoadMusic(filePath)`

Loads a music file to be streamed (e.g., `.mp3`, `.ogg`). Returns a handle.

### `Audio.UpdateMusic(musicHandle)`

Updates the buffer for a streaming music track. This **must** be called every frame for music to play correctly.

### `Audio.Play(musicHandle)`

Starts playing a music track.

### `Audio.Stop(musicHandle)` / `Audio.Pause(musicHandle)` / `Audio.Resume(musicHandle)`

Controls music playback.

### `Music.Free(musicHandle)`

Unloads a music stream.

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
