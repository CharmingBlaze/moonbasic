# Lighting — point, directional, and spot lights

> Illuminate 3D entities with colored lights, control intensity and range, and optionally cast shadows.

**Namespaces:** `LIGHT` · **Status:** Shipped · **Platform:** full runtime

**Commands:** [COMMAND_REGISTRY.md#camera-light](../COMMAND_REGISTRY.md#camera-light) · **Overview:** [02-CAMERA-LIGHT.md](../02-CAMERA-LIGHT.md)

---

## Table of contents

- [At a glance](#at-a-glance)
- [When to use this system](#when-to-use-this-system)
- [Choose the right tool](#choose-the-right-tool)
- [Core workflow](#core-workflow)
- [Light types](#light-types)
- [Key commands](#key-commands)
- [Full example](#full-example)
- [Common mistakes](#common-mistakes)
- [Memory notes](#memory-notes)
- [See also](#see-also)

---

## At a glance

| Idea | Detail |
|------|--------|
| **You get** | Point, directional, and spot lights with color, intensity, range, shadows |
| **You need first** | Window loop + camera + entities ([GAME-LOOP-AND-RENDERING.md](GAME-LOOP-AND-RENDERING.md)) |
| **Typical games** | 3D adventures, horror, outdoor scenes |
| **Not for** | Flat 2D sprite games — art is pre-lit in texture |

**Why lights exist:** Materials react to light direction and color. Without `LIGHT.*`, meshes look flat gray unless you bake lighting into textures.

---

## When to use this system

**Use when:**

- 3D models should read depth and form (indoor, sunset, flashlight).
- You need a sun (directional) or torches (points) that move.
- You want shadow silhouettes for mood or gameplay cues.

**Skip when:**

- Pure 2D `SPRITE.DRAW` stack — use bright textures instead.
- Stylized flat shading — set `MATERIAL.SETCOLOR` only ([MESHES-MODELS-MATERIALS.md](MESHES-MODELS-MATERIALS.md)).

---

## Choose the right tool

| I want to… | Use | Not |
|------------|-----|-----|
| Sun / moon | `LIGHT.CREATEDIRECTIONAL` | Point light at infinity hack |
| Lamp, fire, bulb | `LIGHT.CREATEPOINT` | Directional for local glow |
| Flashlight, stage spot | `LIGHT.CREATESPOT` | Point with tiny range |
| Follow player torch | `LIGHT.SETPOS` each frame | Static light at origin |
| Outdoor day | 1 directional + ambient clear color | 50 point lights |

---

## Core workflow

1. **Create** — `LIGHT.CREATEPOINT`, `CREATEDIRECTIONAL`, or `CREATESPOT`.  
   **Why:** Registers GPU light slot.

2. **Place** — `LIGHT.SETPOS` or `LIGHT.SETDIR`.  
   **Why:** Directional uses direction vector; point/spot use position.

3. **Tune** — `SETCOLOR`, `SETINTENSITY`, `SETRANGE` (point/spot).  
   **Why:** Art direction and readability.

4. **Optional shadows** — `LIGHT.SETSHADOW(light, true)` on key lights only.  
   **Why:** Shadows cost GPU — one or two shadow casters usually enough.

5. **Draw loop** — lights apply automatically inside `RENDER.BEGIN` → `SCENE.DRAW`.  
   **Why:** No per-entity light bind calls in typical workflow.

6. **Unload** — `LIGHT.FREE(light)` on level change.

---

## Light types

| Type | API | Best for |
|------|-----|----------|
| **Directional** | `CREATEDIRECTIONAL` | Sun — parallel rays, no position falloff |
| **Point** | `CREATEPOINT` | Torches, pickups — omni from a point |
| **Spot** | `CREATESPOT` | Headlights, narrow cones |

---

## Key commands

| Command | Why |
|---------|-----|
| `LIGHT.SETPOS(light, x, y, z)` | World position (point, spot) |
| `LIGHT.SETDIR(light, x, y, z)` | Aim direction (directional, spot) |
| `LIGHT.SETCOLOR(light, r, g, b [, a])` | Tint (warm torch, cold moon) |
| `LIGHT.SETINTENSITY(light, value)` | Brightness multiplier |
| `LIGHT.SETRANGE(light, distance)` | Point/spot falloff distance |
| `LIGHT.SETSHADOW(light, enabled)` | Shadow map for this light |
| `LIGHT.FREE(light)` | Release on level unload |

---

## Full example

**Runnable:** [examples/guides/lighting.mb](../../../examples/guides/lighting.mb)

```basic
; Check: moonbasic --check examples/guides/lighting.mb
; Run:   moonrun examples/guides/lighting.mb

APP.OPEN(800, 600, "Lighting demo")
APP.SETFPS(60)

cam = CAMERA.CREATE()
CAMERA.SETACTIVE(cam)
CAMERA.SETPOS(cam, 0, 3, -10)
CAMERA.LOOKAT(cam, 0, 0, 0)

sun = LIGHT.CREATEDIRECTIONAL(-0.5, -1, -0.3, 255, 240, 220, 1.2)
torch = LIGHT.CREATEPOINT(2, 2, 0, 255, 180, 80, 2.0)

floor = ENTITY.CREATECUBE(20, 0.2, 20)
floor.pos(0, -1, 0)
pillar = ENTITY.CREATECUBE(1, 4, 1)
pillar.pos(0, 1, 0)

WHILE NOT APP.SHOULDCLOSE()
    RENDER.CLEAR(15, 18, 25)
    RENDER.BEGIN(cam)
    SCENE.DRAW()
    RENDER.END()
    RENDER.FRAME()
WEND

LIGHT.FREE(torch)
LIGHT.FREE(sun)
APP.CLOSE()
```

---

## Common mistakes

| Mistake | Fix |
|---------|-----|
| Everything black | Add directional or point light; check `SETINTENSITY` > 0 |
| No shadows | `SETSHADOW(true)` on one key light; mesh needs materials |
| Too many shadow lights | Limit to 1–2 — performance cliff |
| Torch stuck at origin | `SETPOS` each frame to follow entity |
| Lights in 2D-only game | Unnecessary — skip `LIGHT.*` |

---

## Memory notes

- `LIGHT.FREE` when leaving level — light slots are finite.
- Shadow maps allocate GPU memory per shadow-casting light.

---

## See also

- [MESHES-MODELS-MATERIALS.md](MESHES-MODELS-MATERIALS.md) — materials respond to lights
- [02-CAMERA-LIGHT.md](../02-CAMERA-LIGHT.md) — camera + light overview
- [reference/LIGHT.md](../../reference/LIGHT.md) — exhaustive overload list
