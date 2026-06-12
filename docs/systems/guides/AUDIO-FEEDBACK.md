# Audio — feedback, music, and 3D sound

> Hear jumps, hits, and ambience; attach sounds to the world or listener camera.

**Namespaces:** `AUDIO` · `SOUND` · **Status:** Shipped

**Commands:** [COMMAND_REGISTRY.md#audio](../COMMAND_REGISTRY.md#audio) · [06-AUDIO.md](../06-AUDIO.md)

---

## Table of contents

- [At a glance](#at-a-glance)
- [When to use which audio API](#when-to-use-which-audio-api)
- [One-shot sounds (2D)](#one-shot-sounds-2d)
- [Music streaming](#music-streaming)
- [3D positional audio](#3d-positional-audio)
- [Asset pack sounds](#asset-pack-sounds)
- [Full example](#full-example)
- [Common mistakes](#common-mistakes)
- [See also](#see-also)

---

## At a glance

| Type | API | Why |
|------|-----|-----|
| UI click, jump bark | `AUDIO.PLAYSOUND` | Same volume everywhere |
| Level music | `AUDIO.LOADMUSIC` + `UPDATEMUSIC` each frame | Streams long files |
| Explosion at point | `SOUND.PLAY3D` / `AUDIO3D` aliases | Volume by distance |
| Packaged game | `ASSET.SOUND("jump")` | No hard-coded paths |

**Why `UPDATEMUSIC`:** Streaming decoders need per-frame pump — unlike short WAV one-shots.

---

## When to use which audio API

| Need | Use |
|------|-----|
| HUD click | `AUDIO.PLAYSOUND` |
| Background loop | `AUDIO.PLAYMUSIC` + `UPDATEMUSIC` |
| Footsteps at feet | `SOUND.PLAY3D` at entity position |
| Head-relative mix | `AUDIO.LISTENERCAMERA(cam)` each frame |

---

## One-shot sounds (2D)

```basic
jump = AUDIO.LOADSOUND("audio/jump.wav")
IF ACTION.HIT("Jump") THEN AUDIO.PLAYSOUND(jump)
```

| Command | Why |
|---------|-----|
| `AUDIO.LOADSOUND(path)` | Load sample to memory |
| `AUDIO.PLAYSOUND(handle)` | Fire once |
| `AUDIO.SETVOLUME(handle, 0.8)` | Per-sample gain |
| `AUDIO.FREESOUND(handle)` | Release on level unload |

---

## Music streaming

```basic
theme = AUDIO.LOADMUSIC("audio/theme.ogg")
AUDIO.PLAYMUSIC(theme)
WHILE playing
    AUDIO.UPDATEMUSIC(theme)    ; required each frame
    ; ... game ...
WEND
AUDIO.STOPMUSIC(theme)
```

**Why not `PLAYSOUND` for theme:** Long OGG/MP3 streams; music API handles buffering.

---

## 3D positional audio

**Why listener:** Left/right balance depends on **camera** (or hero) position.

```basic
AUDIO.LISTENERCAMERA(cam)       ; each frame
boom = AUDIO.LOADSOUND("audio/boom.wav")
SOUND.PLAY3D(boom, ex, ey, ez) ; world position
```

Checklist names `AUDIO3D.*` map to `SOUND.*` / `AUDIO.LISTENER*` — see [06-AUDIO.md](../06-AUDIO.md).

---

## Asset pack sounds

**Why:** Ship one `assets.json` instead of scattering paths in code.

```basic
ASSET.LOADPACK("assets/assets.json")
jump = ASSET.SOUND("jump")
AUDIO.PLAYSOUND(jump)
```

---

## Full example

```basic
APP.OPEN(640, 480, "Audio")
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

`moonrun` required.

---

## Common mistakes

| Mistake | Fix |
|---------|-----|
| No `UPDATEMUSIC` | Music stalls or silent |
| 3D sound without listener update | Wrong panning |
| Load same path every jump | Load once, play many |
| Forget `FREESOUND` on reload | Leak on level change |

---

## See also

- [ASSETS-PIPELINE.md](ASSETS-PIPELINE.md)
- [03-ASSETS.md](../03-ASSETS.md)
