# Game systems reference

moonBASIC is organized into **40 game systems** — window, rendering, entities, physics, audio, saves, and tooling. This section documents each system for beginners and matches the foundation checklist in [FINAL_POLISH_SYSTEMS.md](../FINAL_POLISH_SYSTEMS.md).

**New to moonBASIC?** Read [../BEGIN_HERE.md](../BEGIN_HERE.md) and [00-START.md](00-START.md) first — install, `moonrun` vs `moonbasic`, and **why** each loop command exists.

**Deep guides (all 40 systems + multiplayer):** [GUIDES.md](GUIDES.md) — **24 topic guides**

**Every command (arity + returns):** [COMMAND_REGISTRY.md](COMMAND_REGISTRY.md) · **Entire engine:** [API_CONSISTENCY.md](../API_CONSISTENCY.md)

**Style:** Every page follows [DOCUMENTATION_STYLE_GUIDE.md](../DOCUMENTATION_STYLE_GUIDE.md).

**Case:** All command names are **case-insensitive** (`app.open` = `APP.OPEN`).

---

## Table of contents

### Recommended build order

0. [00 — Start](00-START.md) — why these commands; foundation loop explained
0b. **[Topic guides](GUIDES.md)** — entity, 2D/3D collision, UI, multiplayer (detailed how/why)
1. [01 — Core](01-CORE.md) — APP, RENDER, SCENE, ENTITY
2. [02 — Camera & light](02-CAMERA-LIGHT.md) — CAMERA, LIGHT
3. [03 — Assets](03-ASSETS.md) — MESH, MODEL, MATERIAL, TEXTURE, ASSET
4. [04 — Input](04-INPUT.md) — INPUT, ACTION
5. [10 — Debug & timer](10-DEBUG-TIMER.md) — DEBUG, ERROR, TIMER
6. [03 — Assets](03-ASSETS.md) — ASSET pack
7. [03 — Assets](03-ASSETS.md) — MODEL, ANIMATION (on entities)
8. [05 — Physics](05-PHYSICS.md) — PICK, COLLISION, PHYSICS, BODY
9. [06 — Audio](06-AUDIO.md) — AUDIO, AUDIO3D
10. [07 — 2D & world](07-2D-WORLD.md) — SPRITE, TILEMAP, TERRAIN, PARTICLE
11. [08 — UI & text](08-UI-TEXT.md) — UI, FONT, TEXT
12. [09 — Data](09-DATA.md) — SAVE, FILE, JSON, MATH, VEC3
13. [11 — Tooling](11-TOOLING.md) — PROJECT, PACKAGE, MODULE, HELP, TEST, TEMPLATE

### All systems (alphabetical)

| # | System | Doc | Status |
|---|--------|-----|--------|
| 1 | APP | [01-CORE](01-CORE.md#app-system) | Shipped (aliases) |
| 2 | RENDER | [01-CORE](01-CORE.md#render-system) | Shipped |
| 3 | SCENE | [01-CORE](01-CORE.md#scene-system) | Shipped |
| 4 | ENTITY | [01-CORE](01-CORE.md#entity-system) | Shipped |
| 5 | CAMERA | [02-CAMERA-LIGHT](02-CAMERA-LIGHT.md#camera-system) | Shipped |
| 6 | LIGHT | [02-CAMERA-LIGHT](02-CAMERA-LIGHT.md#light-system) | Shipped |
| 7 | MESH | [03-ASSETS](03-ASSETS.md#mesh-system) | Shipped |
| 8 | MODEL | [03-ASSETS](03-ASSETS.md#model-system) | Shipped |
| 9 | MATERIAL | [03-ASSETS](03-ASSETS.md#material-system) | Shipped |
| 10 | TEXTURE | [03-ASSETS](03-ASSETS.md#texture-system) | Shipped |
| 11 | ANIMATION | [07-2D-WORLD](07-2D-WORLD.md#animation-system) | Partial |
| 12 | INPUT | [04-INPUT](04-INPUT.md#input-system) | Shipped |
| 13 | ACTION | [04-INPUT](04-INPUT.md#action-system) | Shipped |
| 14 | PHYSICS | [05-PHYSICS](05-PHYSICS.md#physics-system) | Shipped |
| 15 | BODY | [05-PHYSICS](05-PHYSICS.md#body-system) | Shipped (aliases) |
| 16 | COLLISION | [05-PHYSICS](05-PHYSICS.md#collision-system) | Partial |
| 17 | PICK | [05-PHYSICS](05-PHYSICS.md#pick-system) | Shipped |
| 18 | AUDIO | [06-AUDIO](06-AUDIO.md#audio-system) | Shipped |
| 19 | AUDIO3D | [06-AUDIO](06-AUDIO.md#audio3d-system) | Partial |
| 20 | UI | [08-UI-TEXT](08-UI-TEXT.md#ui-system) | Partial |
| 21 | FONT / TEXT | [08-UI-TEXT](08-UI-TEXT.md) | Shipped |
| 22 | SPRITE | [07-2D-WORLD](07-2D-WORLD.md#sprite-system) | Shipped |
| 23 | TILEMAP | [07-2D-WORLD](07-2D-WORLD.md#tilemap-system) | Shipped |
| 24 | TERRAIN | [07-2D-WORLD](07-2D-WORLD.md#terrain-system) | Shipped |
| 25 | PARTICLE | [07-2D-WORLD](07-2D-WORLD.md#particle-system) | Shipped |
| 26 | TIMER | [10-DEBUG-TIMER](10-DEBUG-TIMER.md#timer-system) | Shipped |
| 27 | SAVE | [09-DATA](09-DATA.md#save-system) | Shipped |
| 28 | ASSET | [03-ASSETS](03-ASSETS.md#asset-system) | Shipped |
| 29 | FILE | [09-DATA](09-DATA.md#file-system) | Shipped |
| 30 | JSON | [09-DATA](09-DATA.md#json-system) | Shipped |
| 31 | MATH | [09-DATA](09-DATA.md#math-system) | Shipped |
| 32 | VEC3 | [09-DATA](09-DATA.md#vec3-system) | Shipped |
| 33 | DEBUG | [10-DEBUG-TIMER](10-DEBUG-TIMER.md#debug-system) | Shipped |
| 34 | ERROR | [10-DEBUG-TIMER](10-DEBUG-TIMER.md#error-system) | Shipped (compiler) |
| 35 | PROJECT | [11-TOOLING](11-TOOLING.md#project-system) | Shipped |
| 36 | PACKAGE | [11-TOOLING](11-TOOLING.md#package-system) | Partial |
| 37 | MODULE | [11-TOOLING](11-TOOLING.md#module-system) | Shipped |
| 38 | HELP | [11-TOOLING](11-TOOLING.md#help-system) | Partial |
| 39 | TEST | [11-TOOLING](11-TOOLING.md#test-system) | Shipped |
| 40 | TEMPLATE | [11-TOOLING](11-TOOLING.md#template-system) | Shipped |

---

## Foundation loop

Minimum 3D loop (from checklist):

```basic
APP.OPEN(1280, 720, "Test")
cam = CAMERA.CREATE()
CAMERA.SETACTIVE(cam)
CAMERA.SETPOS(cam, 0, 2, -8)
CAMERA.LOOKAT(cam, 0, 0, 0)

cube = ENTITY.CREATECUBE(2, 2, 2)
cube.pos(0, 0, 5)

WHILE NOT APP.SHOULDCLOSE()
    cube.turn(0, 60 * APP.DELTA(), 0)
    RENDER.CLEAR(20, 20, 30)
    RENDER.BEGIN(cam)
    SCENE.DRAW()
    RENDER.END()
    RENDER.FRAME()
WEND

APP.CLOSE()
```

Full sample: [examples/foundation/main.mb](../../examples/foundation/main.mb)

Run: `moonrun examples/foundation/main.mb` · Check: `moonbasic --check examples/foundation/main.mb`

---

## See also

- [Documentation hub](../README.md)
- [COMMAND_REGISTRY.md](COMMAND_REGISTRY.md) — complete command list (beginner namespaces)
- [DOCUMENTATION_STYLE_GUIDE.md](../DOCUMENTATION_STYLE_GUIDE.md) — page format for these guides
- [FINAL_POLISH_SYSTEMS.md](../FINAL_POLISH_SYSTEMS.md) — shipped vs planned checklist
- [COMMANDS.md](../COMMANDS.md) — full namespace index
- [reference/](../reference/) — engine deep reference
