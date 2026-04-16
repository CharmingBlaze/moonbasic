# Audio Commands

Load and play **sound effects** and **streaming music** through Raylib. **CGO** builds required for full audio (see [BUILDING.md](../BUILDING.md)).

Page shape: [DOC_STYLE_GUIDE.md](../DOC_STYLE_GUIDE.md) (**WAVE pattern**).

## Core Workflow

**`AUDIO.INIT()`** once at startup → **`AUDIO.LOADSOUND`** / **`AUDIO.LOADMUSIC`** → each frame **`AUDIO.UPDATEMUSIC`** for music → **`AUDIO.PLAY`** / **`AUDIO.STOP`** as needed → **`SOUND.FREE`** / **`MUSIC.FREE`**, then **`AUDIO.CLOSE()`** before exit.

For **spatial** pan/falloff, call **`AUDIO.LISTENERCAMERA(cam)`** each frame before **`EmitSound(sound, entity)`** (flat global — see manifest). **`Listener(cam)`** is an alias of **`AUDIO.LISTENERCAMERA`**. **`Load3DSound(path)`** loads the same buffers as **`AUDIO.LOADSOUND`** for scripts pairing listener + emit.

---

### `AUDIO.INIT()`

Initializes the audio device. Call before other **`AUDIO.*`** / **`SOUND.*`** / **`MUSIC.*`** commands.

---

### `AUDIO.CLOSE()`

Closes the device and releases audio resources.

---

## Sound effects

### `AUDIO.LOADSOUND(path)`

Loads a short sound (**`.wav`**, **`.ogg`**, …) into memory. Returns a **sound handle**.

---

### `AUDIO.PLAY(handle)` / `AUDIO.STOP(handle)` / `AUDIO.PAUSE(handle)` / `AUDIO.RESUME(handle)`

Playback control (sounds and music share **`AUDIO.PLAY`** where the manifest overloads allow).

---

### `AUDIO.SETSOUNDVOLUME(handle, v)` / `AUDIO.SETSOUNDPITCH(handle, p)` / `AUDIO.SETSOUNDPAN(handle, pan)`

Per-sound mix. Flat aliases **`SoundVolume`**, **`SoundPitch`** map to these.

---

### `SOUND.FREE(handle)`

Unloads a sound and releases its heap slot.

---

### `AUDIO.LISTENERCAMERA(cam)` / `Listener(cam)`

Sets the spatial listener from a **3D** camera handle (**`CAMERA.CREATE`**). Call **each frame** before **`EmitSound`** so pan and falloff stay correct.

---

### `Load3DSound(path)`

Same buffers as **`AUDIO.LOADSOUND`**; naming reflects use with **`Listener`** + **`EmitSound`** (see manifest description).

---

### `EmitSound(sound, entity)`

Plays once with distance falloff and stereo pan (see runtime). Requires a valid **entity** id and an up-to-date listener.

---

## Music (streaming)

### `AUDIO.LOADMUSIC(path)`

Returns a **music handle** (streamed from disk).

---

### `AUDIO.UPDATEMUSIC(handle)`

**Must be called every frame** while music should advance.

---

### `AUDIO.SETMUSICVOLUME(handle, v)` / `AUDIO.SETMUSICPITCH(handle, p)` / `AUDIO.GETMUSICLENGTH(handle)` / `AUDIO.GETMUSICTIME(handle)` / `AUDIO.SEEKMUSIC(handle, t)`

Streaming controls — see [API_CONSISTENCY.md](../API_CONSISTENCY.md) for arities.

---

### `MUSIC.FREE(handle)`

Unloads a music stream.

---

## Full Example

This example assumes **`jump.wav`** and **`theme.mp3`** next to the program.

```basic
WINDOW.OPEN(800, 600, "Audio Example")
WINDOW.SETFPS(60)

AUDIO.INIT()

jumpSfx = AUDIO.LOADSOUND("jump.wav")
bgMusic = AUDIO.LOADMUSIC("theme.mp3")
AUDIO.PLAY(bgMusic)

WHILE NOT WINDOW.SHOULDCLOSE()
    AUDIO.UPDATEMUSIC(bgMusic)

    IF INPUT.KEYPRESSED(KEY_SPACE) THEN
        AUDIO.PLAY(jumpSfx)
    ENDIF

    RENDER.CLEAR(40, 40, 40)
    CAMERA2D.BEGIN()
        DRAW.TEXT("Press SPACE to play a sound!", 190, 200, 20, 255, 255, 255, 255)
    CAMERA2D.END()
    RENDER.FRAME()
WEND

SOUND.FREE(jumpSfx)
MUSIC.FREE(bgMusic)
AUDIO.CLOSE()
WINDOW.CLOSE()
```

---

## See also

- [WAVE.md](WAVE.md) — raw **`WAVE.*`** samples → **`SOUND.FROMWAVE`**
