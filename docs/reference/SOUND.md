# Sound Commands

Extended sound handle commands: create a sound from a wave, play it spatially in 3D, pin it to an entity, and free it. For general audio loading and playback see [AUDIO.md](AUDIO.md).

## Core Workflow

1. `AUDIO.INIT()` — start the audio device.
2. `wave = WAVE.LOAD("sfx.wav")` — load raw wave data.
3. `snd = SOUND.FROMWAVE(wave)` — create a sound handle.
4. `SOUND.PLAY3D(snd, x, y, z, maxDist)` — play spatially, or `SOUND.ATTACH(snd, entity)` to track an entity.
5. `SOUND.FREE(snd)` → `WAVE.UNLOAD(wave)` → `AUDIO.CLOSE()` on exit.

---

## Commands

### `SOUND.FROMWAVE(waveHandle)` 

Creates a sound handle from a loaded wave handle (from `WAVE.LOAD`). Returns a **sound handle**.

---

### `SOUND.PLAY3D(sound, x, y, z, maxDist)` 

Plays the sound spatially at world position `(x, y, z)`. Volume attenuates to zero at `maxDist`. Requires an active audio listener position (set via `AUDIO.SETLISTENER` or camera position).

---

### `SOUND.ATTACH(sound, entityHandle)` 

Pins the sound emitter to an entity. Position updates automatically each frame with the entity.

---

### `SOUND.FREE(sound)` 

Frees the sound handle.

---

## Full Example

Footstep sound pinned to a moving entity.

```basic
WINDOW.OPEN(960, 540, "Sound Demo")
WINDOW.SETFPS(60)

AUDIO.INIT()
wave  = WAVE.LOAD("assets/footstep.wav")
snd   = SOUND.FROMWAVE(wave)

player = ENTITY.CREATECUBE(1.0)
ENTITY.SETPOS(player, 0, 0.5, 0)
SOUND.ATTACH(snd, player)

cam = CAMERA.CREATE()
CAMERA.SETPOS(cam, 0, 5, -10)
CAMERA.SETTARGET(cam, 0, 0, 0)

px = 0.0
WHILE NOT WINDOW.SHOULDCLOSE()
    dt = TIME.DELTA()
    IF INPUT.KEYDOWN(KEY_RIGHT) THEN
        px = px + 4 * dt
        IF INPUT.KEYPRESSED(KEY_RIGHT) THEN SOUND.PLAY3D(snd, px, 0, 0, 20)
    END IF
    ENTITY.SETPOS(player, px, 0.5, 0)
    ENTITY.UPDATE(dt)

    RENDER.CLEAR(20, 25, 35)
    RENDER.BEGIN3D(cam)
        ENTITY.DRAWALL()
        DRAW.GRID(20, 1.0)
    RENDER.END3D()
    RENDER.FRAME()
WEND

SOUND.FREE(snd)
WAVE.UNLOAD(wave)
AUDIO.CLOSE()
WINDOW.CLOSE()
```

---

## See also

- [AUDIO.md](AUDIO.md) — `AUDIO.LOAD`, `AUDIO.PLAY`, `AUDIO.STOP`
- [MUSIC.md](MUSIC.md) — background music streaming
