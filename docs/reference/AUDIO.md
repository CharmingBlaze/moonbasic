# Audio Commands

Load and play **sound effects** and **streaming music** through Raylib. **CGO** builds required for full audio (see [BUILDING.md](../BUILDING.md)).

Page shape: [DOC_STYLE_GUIDE.md](../DOC_STYLE_GUIDE.md) (**WAVE pattern**).

## Core Workflow

**`AUDIO.INIT()`** once at startup → **`AUDIO.LOADSOUND`** / **`AUDIO.LOADMUSIC`** → each frame **`AUDIO.UPDATEMUSIC`** for music → **`AUDIO.PLAY`** / **`AUDIO.STOP`** as needed → **`SOUND.FREE`** / **`MUSIC.FREE`**, then **`AUDIO.CLOSE()`** before exit.

For **spatial** pan/falloff, call **`AUDIO.LISTENERCAMERA(cam)`** each frame before **`EmitSound(sound, entity)`** (flat global — see manifest). **`Listener(cam)`** is an alias of **`AUDIO.LISTENERCAMERA`**. **`Load3DSound(path)`** loads the same buffers as **`AUDIO.LOADSOUND`** for scripts pairing listener + emit.

---

### `AUDIO.INIT()`
Initializes the audio device. Call before other **`AUDIO.*`** / **`SOUND.*`** / **`MUSIC.*`** commands.

- **Returns**: (None)

---

### `AUDIO.CLOSE()`
Closes the device and releases audio resources.

---

## Sound effects

### `AUDIO.LOADSOUND(path)`
Loads a short sound (**`.wav`**, **`.ogg`**, …) into memory. Returns a **sound handle**.

- **Arguments**:
    - `path`: (String) File path relative to working directory.
- **Returns**: (Handle) The new sound handle.
- **Example**:
    ```basic
    jumpSfx = AUDIO.LOADSOUND("jump.wav")
    ```

---

### `AUDIO.PLAY(handle)` / `STOP` / `PAUSE` / `RESUME`
Playback control for sounds and music.

- **Arguments**:
    - `handle`: (Handle) The sound or music stream to control.
- **Returns**: (Handle) The sound/music handle (for chaining).

---

### `AUDIO.SETSOUNDVOLUME(handle, v)` / `SETSOUNDPITCH` / `SETSOUNDPAN`
Per-sound mix settings.

- **Arguments**:
    - `handle`: (Handle) The sound to modify.
    - `v/p/pan`: (Float) Volume (0-1), Pitch (0-2), or Pan (-1 to 1).
- **Returns**: (Handle) The sound handle (for chaining).

---

### `SOUND.FREE(handle)`
Unloads a sound and releases its heap slot.

---

### `AUDIO.LISTENERCAMERA(cam)` / `Listener(cam)`
Sets the spatial listener from a **3D** camera handle.

- **Arguments**:
    - `cam`: (Handle) The 3D camera to track.
- **Returns**: (None)

---

### `EmitSound(sound, entity)`
Plays a sound once with distance falloff and stereo pan based on an entity's position.

- **Arguments**:
    - `sound`: (Handle) The sound to play.
    - `entity`: (Handle) The world entity providing the position.
- **Returns**: (None)

---

## Music (streaming)

### `AUDIO.LOADMUSIC(path)`
Returns a **music handle** (streamed from disk).

- **Arguments**:
    - `path`: (String) File path.
- **Returns**: (Handle) The music handle.

---

### `AUDIO.UPDATEMUSIC(handle)`
**Must be called every frame** while music should advance.

- **Arguments**:
    - `handle`: (Handle) The music stream to update.
- **Returns**: (None)

---

### `AUDIO.SETMUSICVOLUME(handle, v)` / `SETMUSICPITCH`
Streaming mix settings.

- **Returns**: (Handle) The music handle (for chaining).

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
