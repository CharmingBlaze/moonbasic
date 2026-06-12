# Audio systems: AUDIO and AUDIO3D

> Sound effects, music streams, and spatial audio attached to the world or camera.

**All commands:** [COMMAND_REGISTRY.md#audio](COMMAND_REGISTRY.md#audio)

**Deep guide:** [guides/AUDIO-FEEDBACK.md](guides/AUDIO-FEEDBACK.md)

**See also:** [03-ASSETS](03-ASSETS.md) · [reference/AUDIO.md](../reference/AUDIO.md) · [reference/SOUND.md](../reference/SOUND.md)

---

## Table of contents

- [AUDIO system](#audio-system)
- [AUDIO3D system](#audio3d-system)
- [Full example](#full-example)
- [Memory notes](#memory-notes)
- [See also](#see-also)

---

## AUDIO system

2D sound playback and music management.

### Core workflow

1. `AUDIO.LOADSOUND(path)` or `ASSET.SOUND(id)` after `ASSET.LOADPACK`.
2. `AUDIO.PLAYSOUND(handle)` for one-shots; `AUDIO.STOPSOUND` to stop.
3. Music: `AUDIO.LOADMUSIC`, `AUDIO.PLAYMUSIC`, `AUDIO.UPDATEMUSIC` each frame, `AUDIO.STOPMUSIC`.
4. `AUDIO.FREESOUND` / free music handle when done.

---

### `AUDIO.LOADSOUND(path)`

Loads a WAV/OGG sample.

| Argument | Type | Description |
|----------|------|-------------|
| path | string | Audio file path |

**Returns:** `handle`

**Example:**

```basic
jump = AUDIO.LOADSOUND("audio/jump.wav")
```

---

### `AUDIO.PLAYSOUND(sound)` / `AUDIO.STOPSOUND(sound)`

Play or stop a loaded sound.

**Example:**

```basic
IF ACTION.HIT("Jump") THEN AUDIO.PLAYSOUND(jump)
```

---

### `AUDIO.SETVOLUME(sound, volume)`

Volume 0.0–1.0 for a sound handle.

**Example:**

```basic
AUDIO.SETVOLUME(jump, 0.8)
```

---

### Music

| Command | Description |
|---------|-------------|
| `AUDIO.LOADMUSIC(path)` | Load streaming music |
| `AUDIO.PLAYMUSIC(music)` | Start playback |
| `AUDIO.UPDATEMUSIC(music)` | Pump stream each frame |
| `AUDIO.STOPMUSIC(music)` | Stop music |
| `AUDIO.SETMUSICVOLUME(music, v)` | Music volume |

**Example:**

```basic
theme = AUDIO.LOADMUSIC("audio/theme.ogg")
AUDIO.PLAYMUSIC(theme)
WHILE NOT APP.SHOULDCLOSE()
    AUDIO.UPDATEMUSIC(theme)
    ; ... frame ...
WEND
```

---

## AUDIO3D system

Checklist **`AUDIO3D.*`** maps to spatial **`SOUND.*`** and listener helpers.

### Core workflow

1. Load sound (`AUDIO.LOADSOUND` or `ASSET.SOUND`).
2. `SOUND.PLAY3D(sound, x, y, z)` or attach to entity.
3. `AUDIO.LISTENERCAMERA(cam)` or `AUDIO3D.SETLISTENER(camera)` each frame.
4. Set range / falloff on emitters.

---

### Spatial aliases

| Checklist / alias | Canonical |
|-------------------|-----------|
| `AUDIO3D.PLAYAT(snd, x, y, z)` | `SOUND.PLAY3D` |
| `AUDIO3D.ATTACH(snd, entity)` | `SOUND.ATTACH` / entity-bound play |
| `AUDIO3D.SETLISTENER(camera)` | `AUDIO.LISTENERCAMERA` |
| `AUDIO3D.SETRANGE(snd, dist)` | `SOUND.SETRANGE` |
| `AUDIO.PLAYSOUND` style | `AUDIO.PLAY` |

**Example:**

```basic
boom = AUDIO.LOADSOUND("audio/explosion.wav")
AUDIO.LISTENERCAMERA(cam)
SOUND.PLAY3D(boom, 10, 0, 5)
```

See [reference/SOUND.md](../reference/SOUND.md) for full 3D audio API.

---

## Full example

```basic
APP.OPEN(640, 480, "Audio")
APP.SETFPS(60)

jump = AUDIO.LOADSOUND("assets/jump.wav")
cam = CAMERA.CREATE()
CAMERA.SETACTIVE(cam)

WHILE NOT APP.SHOULDCLOSE()
    AUDIO.LISTENERCAMERA(cam)
    IF INPUT.KEYHIT(KEY_SPACE) THEN AUDIO.PLAYSOUND(jump)

    RENDER.CLEAR(0, 0, 0)
    RENDER.FRAME()
WEND

AUDIO.FREESOUND(jump)
APP.CLOSE()
```

---

## Memory notes

- **`AUDIO.FREESOUND`** when unloading levels; **`ASSET.UNLOAD`** frees packed sounds too.
- Call **`AUDIO.UPDATEMUSIC`** every frame while music plays.

---

## See also

- [03-ASSETS](03-ASSETS.md) — `ASSET.SOUND("jump")`
- [examples/rpg](../examples/rpg/main.mb) — gameplay with assets
