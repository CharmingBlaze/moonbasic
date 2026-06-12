# Animation — skeletal clips and state machines

> Play walk/run clips on imported models and drive multi-state animation with `ANIM.*` helpers.

**Namespaces:** `ENTITY` (clip playback) · `MODEL` · `ANIM` · **Status:** Partial (dual API)

**Commands:** [COMMAND_REGISTRY.md#2d-world](../COMMAND_REGISTRY.md#2d-world) (ANIM section) · **Overview:** [07-2D-WORLD.md](../07-2D-WORLD.md#animation-system)

---

## Table of contents

- [At a glance](#at-a-glance)
- [When to use this system](#when-to-use-this-system)
- [Choose the right tool](#choose-the-right-tool)
- [Core workflow — skeletal clips](#core-workflow--skeletal-clips)
- [FSM with ANIM.*](#fsm-with-anim)
- [Key commands](#key-commands)
- [Full example](#full-example)
- [Common mistakes](#common-mistakes)
- [See also](#see-also)

---

## At a glance

| Idea | Detail |
|------|--------|
| **You get** | Play/stop clips on entities; optional state machine transitions |
| **You need first** | `MODEL.LOAD` with skeleton ([MESHES-MODELS-MATERIALS.md](MESHES-MODELS-MATERIALS.md)) |
| **Typical games** | Characters, doors, machinery |
| **Not for** | 2D sprite flipbooks — animate `SPRITE` frames manually |

**Why two APIs:** **Entity clips** are one-shot or loop playback. **ANIM FSM** groups clips into states (idle → walk → jump) with rules.

---

## When to use this system

**Use when:**

- Model file contains animation tracks (`hero.glb`).
- Character needs idle/walk/run based on speed.
- Scripted props (open door, spinning wheel) use embedded clips.

**Skip when:**

- Static props — no `PLAYANIM`.
- 2D tile characters — `SPRITE` + timer frame index.
- Procedural motion — `ENTITY.TURN` / physics only.

---

## Choose the right tool

| I want to… | Use | Not |
|------------|-----|-----|
| Play one clip | `ENTITY.PLAYANIM(ent, "Walk")` | Manual bone matrices |
| List clip names | `MODEL.ANIMCOUNT` / `GETANIMNAME` | Guess string names |
| Idle vs run vs jump | `ANIM.DEFINE` + `ANIM.UPDATE` | `IF` chains calling `PLAYANIM` every frame |
| Move on slopes | `CHAR.*` ([CHARACTER-3D-WALKING.md](CHARACTER-3D-WALKING.md)) | Animation only — no locomotion |
| 2D walk cycle | Sprite sheet timing | Skeletal `PLAYANIM` |

---

## Core workflow — skeletal clips

1. **Load model** — `MODEL.LOAD("hero.glb")` → `ENTITY.SETMODEL`.  
   **Why:** Clips live on the model asset.

2. **Discover clips** — `MODEL.ANIMCOUNT(model)` and `GETANIMNAME(model, i)`.  
   **Why:** Exporter names vary (`Walk` vs `walking`).

3. **Play** — `ENTITY.PLAYANIM(hero, "Walk")` or index via `ENTITY.PLAY`.  
   **Why:** Starts GPU skinning for that track.

4. **Stop** — `ENTITY.STOPANIM(hero)` when needed.

5. **Each frame** — entity draw via `SCENE.DRAW` (skinning updates with scene).

**Aliases:** Checklist `ANIM.PLAY` on entity → `ENTITY.PLAYANIM`.

---

## FSM with ANIM.*

For characters with several clips on one **animated entity**:

1. Load skeleton — `ENTITY.LOADANIMATEDMESH("hero.glb")` (or `SETMODEL` with clips).
2. `ANIM.DEFINE(hero, "idle", startFrame, endFrame, fps, loop)` — register states.
3. `ANIM.ADDTRANSITION(hero, "idle", "run", "moving")` — param-driven transitions.
4. Each frame: set `ANIM.SETPARAM(hero, "moving", 1)` when speed > 0, then `ANIM.UPDATE(hero, APP.DELTA())`.

**Why FSM:** Avoids restarting `PLAYANIM` every frame and handles blend windows.

See [reference/ANIM.md](../../reference/ANIM.md) for full example and arity.

---

## Key commands

| Command | Why |
|---------|-----|
| `ENTITY.PLAYANIM(ent, name)` | Loop or play clip by name |
| `ENTITY.PLAY(ent, index)` | Play clip by index |
| `ENTITY.STOPANIM(ent)` | Stop current clip |
| `MODEL.ANIMCOUNT(model)` | How many clips |
| `MODEL.GETANIMNAME(model, i)` | Clip name for UI / debug |
| `ANIM.DEFINE(ent, …)` / `ADDTRANSITION` | State machine on entity |
| `ANIM.SETPARAM(ent, name, value)` | Drive transitions |
| `ANIM.UPDATE(ent, dt)` | Tick transitions each frame |

---

## Full example

```basic
APP.OPEN(800, 600, "Animation")
APP.SETFPS(60)

cam = CAMERA.CREATE()
CAMERA.SETACTIVE(cam)
CAMERA.SETPOS(cam, 0, 2, -6)
CAMERA.LOOKAT(cam, 0, 0, 0)

mdl = MODEL.LOAD("assets/hero.glb")
hero = ENTITY.CREATE("Hero")
ENTITY.SETMODEL(hero, mdl)
hero.pos(0, 0, 3)

; Discover first clip name in debug
n = MODEL.ANIMCOUNT(mdl)
IF n > 0 THEN
    DEBUG.LOG("Clip 0: " + MODEL.GETANIMNAME(mdl, 0))
    ENTITY.PLAYANIM(hero, MODEL.GETANIMNAME(mdl, 0))
ENDIF

WHILE NOT APP.SHOULDCLOSE()
    RENDER.CLEAR(20, 22, 30)
    RENDER.BEGIN(cam)
    SCENE.DRAW()
    RENDER.END()
    RENDER.FRAME()
WEND

MODEL.FREE(mdl)
APP.CLOSE()
```

Check: `moonbasic --check` · Run: `moonrun`

---

## Common mistakes

| Mistake | Fix |
|---------|-----|
| Clip name wrong | Log `GETANIMNAME`; names are case-insensitive in scripts |
| Model static | File has no skeleton — re-export from DCC |
| `PLAYANIM` every frame | Restarts clip — use FSM or call once on state change |
| No model on entity | `SETMODEL` before `PLAYANIM` |
| Expect 2D sprite anim | Use sprite sheet + `TIMER` |

---

## See also

- [MESHES-MODELS-MATERIALS.md](MESHES-MODELS-MATERIALS.md) — load GLB
- [CHARACTER-3D-WALKING.md](CHARACTER-3D-WALKING.md) — movement + clips together
- [reference/ANIM.md](../../reference/ANIM.md) — FSM details
