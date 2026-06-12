# Particles — fire, smoke, sparks, and bursts

> GPU particle emitters for juice: muzzle flashes, dust, magic trails, and ambient atmosphere.

**Namespaces:** `PARTICLE` · **Status:** Shipped

**Commands:** [COMMAND_REGISTRY.md#2d-world](../COMMAND_REGISTRY.md#2d-world) · **Overview:** [07-2D-WORLD.md](../07-2D-WORLD.md#particle-system)

---

## Table of contents

- [At a glance](#at-a-glance)
- [When to use this system](#when-to-use-this-system)
- [Choose the right tool](#choose-the-right-tool)
- [Core workflow](#core-workflow)
- [Key commands](#key-commands)
- [Full example](#full-example)
- [Common mistakes](#common-mistakes)
- [Memory notes](#memory-notes)
- [See also](#see-also)

---

## At a glance

| Idea | Detail |
|------|--------|
| **You get** | Emitter handles with rate, lifetime, color, size, gravity, texture |
| **You need first** | Render loop ([GAME-LOOP-AND-RENDERING.md](GAME-LOOP-AND-RENDERING.md)) |
| **Typical games** | Shooters, RPG spells, racing dust |
| **Not for** | Persistent world geometry — use entities or terrain |

**Why particles:** Thousands of cheap quads with blended alpha beat spawning entities for every spark.

---

## When to use this system

**Use when:**

- Short-lived visual effects (explosion, hit flash).
- Continuous ambience (torch flicker, rain streaks).
- One-shot bursts (`SETBURST`) on gun fire.

**Skip when:**

- Solid debris that collides — use `ENTITY` + physics ([COLLISION-3D.md](COLLISION-3D.md)).
- UI sparkles — `DRAW` / `GUI` layers.

---

## Choose the right tool

| I want to… | Use | Not |
|------------|-----|-----|
| Fire / smoke column | `PARTICLE.CREATE` + `SETRATE` | Many small sprites |
| Colored sparks | `SETCOLOR` / `SETCOLOREND` fade | White-only texture |
| World-space fountain | `SETPOS` + `UPDATE` in 3D pass | Screen-only `DRAW.RECT` |
| Soft round puff | Soft alpha texture + `SETBILLBOARD` | Hard square PNG |
| Footstep dust | Short `SETLIFETIME`, low `SETRATE` | Permanent emitter |

---

## Core workflow

1. **Create** — `fx = PARTICLE.CREATE()`.  
   **Why:** Allocates emitter slot.

2. **Texture** — `tex = TEXTURE.LOAD("spark.png")` → `PARTICLE.SETTEXTURE(fx, tex)`.  
   **Why:** Particle appearance (soft circle works well).

3. **Tune** — `SETRATE`, `SETLIFETIME`, `SETSPEED`, `SETCOLOR`, `SETGRAVITY`, `SETSPREAD`.  
   **Why:** Art direction and motion feel.

4. **Place** — `PARTICLE.SETPOS(fx, x, y, z)` at effect origin.

5. **Loop** — `PARTICLE.UPDATE(fx, APP.DELTA())` then `PARTICLE.DRAW(fx)` inside `RENDER.BEGIN`…`END`.  
   **Why:** Simulation and draw are separate steps.

6. **Control** — `PARTICLE.PLAY(fx)` / `STOP(fx)`; `FREE(fx)` on unload.

---

## Key commands

| Command | Why |
|---------|-----|
| `PARTICLE.CREATE()` | New emitter |
| `PARTICLE.SETTEXTURE(fx, tex)` | Sprite for each particle |
| `PARTICLE.SETRATE(fx, n)` | Spawns per second |
| `PARTICLE.SETLIFETIME(fx, min, max)` | How long each lives |
| `PARTICLE.SETSPEED(fx, min, max)` | Initial speed range |
| `PARTICLE.SETCOLOR` / `SETCOLOREND` | Start/end tint |
| `PARTICLE.SETGRAVITY(fx, y)` | Fall speed |
| `PARTICLE.SETSPREAD(fx, degrees)` | Cone randomness |
| `PARTICLE.SETBURST(fx, count)` | One-shot burst |
| `PARTICLE.UPDATE(fx, dt)` | Simulate |
| `PARTICLE.DRAW(fx)` | Render in 3D pass |
| `PARTICLE.PLAY` / `STOP` | Emit on/off |
| `PARTICLE.FREE(fx)` | Release GPU buffers |

**Aliases:** `PARTICLE.MAKE` → `CREATE`; `SETPOSITION` → `SETPOS`.

---

## Full example

```basic
APP.OPEN(800, 600, "Particles")
APP.SETFPS(60)

cam = CAMERA.CREATE()
CAMERA.SETACTIVE(cam)
CAMERA.SETPOS(cam, 0, 3, -10)
CAMERA.LOOKAT(cam, 0, 0, 0)

tex = TEXTURE.LOAD("assets/soft.png")
fx = PARTICLE.CREATE()
PARTICLE.SETTEXTURE(fx, tex)
PARTICLE.SETRATE(fx, 60)
PARTICLE.SETLIFETIME(fx, 0.4, 1.2)
PARTICLE.SETSPEED(fx, 1, 3)
PARTICLE.SETCOLOR(fx, 255, 200, 80, 255)
PARTICLE.SETCOLOREND(fx, 255, 80, 20, 0)
PARTICLE.SETPOS(fx, 0, 1, 0)
PARTICLE.PLAY(fx)

WHILE NOT APP.SHOULDCLOSE()
    PARTICLE.UPDATE(fx, APP.DELTA())
    RENDER.CLEAR(10, 12, 18)
    RENDER.BEGIN(cam)
    SCENE.DRAW()
    PARTICLE.DRAW(fx)
    RENDER.END()
    RENDER.FRAME()
WEND

PARTICLE.FREE(fx)
TEXTURE.FREE(tex)
APP.CLOSE()
```

---

## Common mistakes

| Mistake | Fix |
|---------|-----|
| No particles visible | Call `UPDATE` + `DRAW` each frame; `PLAY` emitter |
| Wrong layer | Draw inside `RENDER.BEGIN` for world-space |
| Forgot delta on `UPDATE` | Pass `APP.DELTA()` for stable simulation |
| Square harsh edges | Use soft alpha texture; fade with `SETCOLOREND` |
| Emitter never stops | `STOP` or `FREE` on level exit |

---

## Memory notes

- `PARTICLE.FREE` on level unload — stops GPU buffer growth.
- `STOP` stops spawning; existing particles die out per lifetime.

---

## See also

- [MESHES-MODELS-MATERIALS.md](MESHES-MODELS-MATERIALS.md) — `TEXTURE.LOAD` for particle image
- [AUDIO-FEEDBACK.md](AUDIO-FEEDBACK.md) — pair boom SFX with burst
- [07-2D-WORLD.md](../07-2D-WORLD.md) — particles alongside sprites/terrain
