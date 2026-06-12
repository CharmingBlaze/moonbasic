# Deep topic guides — how and why

> Narrative guides for **every beginner system** and common game problems. Each page explains **which commands to use**, **in what order**, and **why** — not just signatures.

**New to moonBASIC?** Read [../BEGIN_HERE.md](../BEGIN_HERE.md) and [00-START.md](00-START.md) first.

**Every command (arity):** [COMMAND_REGISTRY.md](COMMAND_REGISTRY.md) · **Entire engine:** [API_CONSISTENCY.md](../API_CONSISTENCY.md)

**Style:** [DOCUMENTATION_STYLE_GUIDE.md](../DOCUMENTATION_STYLE_GUIDE.md)

---

## Table of contents

- [All 40 systems → guide map](#all-40-systems--guide-map)
- [Guide index by topic](#guide-index-by-topic)
- [Math & vectors (deep library)](#math--vectors-deep-library)
- [How to pick a guide](#how-to-pick-a-guide)
- [Suggested learning paths](#suggested-learning-paths)
- [See also](#see-also)

---

## All 40 systems → guide map

Every system in [README.md](README.md) has a deep guide (some guides cover multiple related namespaces).

| System | Deep guide |
|--------|------------|
| APP | [Game loop & rendering](guides/GAME-LOOP-AND-RENDERING.md) |
| RENDER | [Game loop & rendering](guides/GAME-LOOP-AND-RENDERING.md) |
| SCENE | [Game loop & rendering](guides/GAME-LOOP-AND-RENDERING.md) |
| ENTITY | [Entity system](guides/ENTITY-SYSTEM.md) |
| CAMERA | [Camera & input](guides/CAMERA-AND-INPUT.md) |
| LIGHT | [Lighting](guides/LIGHTING.md) |
| MESH | [Meshes, models, materials](guides/MESHES-MODELS-MATERIALS.md) |
| MODEL | [Meshes, models, materials](guides/MESHES-MODELS-MATERIALS.md) |
| MATERIAL | [Meshes, models, materials](guides/MESHES-MODELS-MATERIALS.md) |
| TEXTURE | [Meshes, models, materials](guides/MESHES-MODELS-MATERIALS.md) |
| ANIM | [Animation](guides/ANIMATION.md) |
| INPUT | [Camera & input](guides/CAMERA-AND-INPUT.md) |
| ACTION | [Camera & input](guides/CAMERA-AND-INPUT.md) |
| PHYSICS | [3D collision & physics](guides/COLLISION-3D.md) |
| BODY | [3D collision & physics](guides/COLLISION-3D.md) |
| COLLISION | [2D collision](guides/COLLISION-2D.md) · [3D collision](guides/COLLISION-3D.md) |
| PICK | [3D collision & physics](guides/COLLISION-3D.md) |
| AUDIO | [Audio & feedback](guides/AUDIO-FEEDBACK.md) |
| AUDIO3D | [Audio & feedback](guides/AUDIO-FEEDBACK.md) |
| UI | [UI & menus](guides/UI-MENUS.md) |
| FONT / TEXT | [UI & menus](guides/UI-MENUS.md) |
| SPRITE | [Sprites & tilemaps](guides/SPRITES-TILEMAPS-2D.md) |
| TILEMAP | [Sprites & tilemaps](guides/SPRITES-TILEMAPS-2D.md) |
| TERRAIN | [Terrain & open world](guides/TERRAIN-OPEN-WORLD.md) |
| PARTICLE | [Particles](guides/PARTICLES.md) |
| TIMER | [Debug & testing](guides/DEBUG-AND-TESTING.md) |
| SAVE | [Save & progress](guides/SAVE-AND-PROGRESS.md) |
| ASSET | [Assets pipeline](guides/ASSETS-PIPELINE.md) |
| FILE | [Files & JSON](guides/FILES-AND-JSON.md) |
| JSON | [Files & JSON](guides/FILES-AND-JSON.md) |
| MATH | [Math hub](guides/MATH-AND-VECTORS.md) → [math/](guides/math/README.md) |
| VEC2 | [2D vector math](guides/math/VEC2-MATH.md) |
| VEC3 | [3D vector math](guides/math/VEC3-MATH.md) |
| DEBUG | [Debug & testing](guides/DEBUG-AND-TESTING.md) |
| ERROR | [Compiler errors](guides/COMPILER-ERRORS.md) |
| PROJECT | [Project workflow](guides/PROJECT-WORKFLOW.md) |
| PACKAGE | [Project workflow](guides/PROJECT-WORKFLOW.md) |
| MODULE | [Project workflow](guides/PROJECT-WORKFLOW.md) |
| HELP | [Project workflow](guides/PROJECT-WORKFLOW.md) |
| TEST | [Debug & testing](guides/DEBUG-AND-TESTING.md) · [Project workflow](guides/PROJECT-WORKFLOW.md) |
| TEMPLATE | [Project workflow](guides/PROJECT-WORKFLOW.md) |

**Also documented (engine extras):** [3D walking characters](guides/CHARACTER-3D-WALKING.md) (`CHAR.*`), [2D physics platformer](guides/PHYSICS-2D-PLATFORMER.md) (`PHYSICS2D`), [Multiplayer](guides/MULTIPLAYER.md), [Networking low-level](guides/NETWORKING-LOW-LEVEL.md).

---

## Guide index by topic

### Core loop and world

| Guide | Systems | Start here if… |
|-------|---------|----------------|
| [Game loop & rendering](guides/GAME-LOOP-AND-RENDERING.md) | APP, RENDER, SCENE | Window, clear, draw, present |
| [Entity system](guides/ENTITY-SYSTEM.md) | ENTITY | Objects in the world |
| [Camera & input](guides/CAMERA-AND-INPUT.md) | CAMERA, INPUT, ACTION | View + controls |
| [Lighting](guides/LIGHTING.md) | LIGHT | 3D looks flat or gray |
| [3D walking characters](guides/CHARACTER-3D-WALKING.md) | CHAR | Humanoid on slopes |

### Visual assets

| Guide | Systems | Start here if… |
|-------|---------|----------------|
| [Meshes, models, materials](guides/MESHES-MODELS-MATERIALS.md) | MESH, MODEL, MATERIAL, TEXTURE | Primitives, GLB, materials |
| [Assets pipeline](guides/ASSETS-PIPELINE.md) | ASSET | `assets.json` manifest |
| [Animation](guides/ANIMATION.md) | ANIM, MODEL clips | Walk/run on skeleton |
| [Particles](guides/PARTICLES.md) | PARTICLE | Fire, smoke, sparks |

### 2D and world scale

| Guide | Systems | Start here if… |
|-------|---------|----------------|
| [Sprites & tilemaps](guides/SPRITES-TILEMAPS-2D.md) | SPRITE, TILEMAP | 2D levels, TMX |
| [Terrain & open world](guides/TERRAIN-OPEN-WORLD.md) | TERRAIN | Large outdoor 3D |
| [2D collision](guides/COLLISION-2D.md) | COLLISION (2D) | Platformer hits |
| [2D physics platformer](guides/PHYSICS-2D-PLATFORMER.md) | PHYSICS2D, BODY2D | Many 2D rigid bodies |
| [3D collision & physics](guides/COLLISION-3D.md) | PHYSICS, BODY, PICK | 3D stacks, mouse pick |

### Audio, UI, data

| Guide | Systems | Start here if… |
|-------|---------|----------------|
| [Audio & feedback](guides/AUDIO-FEEDBACK.md) | AUDIO, AUDIO3D | SFX, music, 3D sound |
| [UI & menus](guides/UI-MENUS.md) | UI, FONT, TEXT | Menus, HUD |
| [Save & progress](guides/SAVE-AND-PROGRESS.md) | SAVE | Level select, hi-score |
| [Files & JSON](guides/FILES-AND-JSON.md) | FILE, JSON | Config files, nested data |
| [Math hub](guides/MATH-AND-VECTORS.md) | MATH, VEC2, VEC3 | Pick focused math guide |

### Math & vectors (deep library)

| Guide | Start here if… |
|-------|----------------|
| [math/README.md](guides/math/README.md) | Index of all math guides |
| [2D game math](guides/math/MATH-2D-GAMEPLAY.md) | `DIST2D`, screen aim, X/Y |
| [3D game math](guides/math/MATH-3D-GAMEPLAY.md) | `HDIST`, yaw, XZ movement |
| [2D vector math](guides/math/VEC2-MATH.md) | Normalize, rotate, pushout |
| [3D vector math](guides/math/VEC3-MATH.md) | Dot, cross, reflect |
| [Interpolation & easing](guides/math/INTERPOLATION-AND-EASING.md) | Lerp, smoothstep, approach |
| [Angles & rotation](guides/math/ANGLES-AND-ROTATION.md) | Wrap, `LERPANGLE`, quats |
| [Randomness](guides/math/RANDOMNESS-AND-PROCEDURE.md) | Loot, dice, seeds |

### Networking and quality

| Guide | Systems | Start here if… |
|-------|---------|----------------|
| [Multiplayer](guides/MULTIPLAYER.md) | SERVER, CLIENT, RPC | Two `moonrun` processes |
| [Networking (mid level)](guides/NETWORKING-LOW-LEVEL.md) | NET, PEER | Custom packets |
| [Debug & testing](guides/DEBUG-AND-TESTING.md) | DEBUG, TIMER, TEST | Watches, timers, tests |
| [Compiler errors](guides/COMPILER-ERRORS.md) | ERROR | `--check` diagnostics |
| [Project workflow](guides/PROJECT-WORKFLOW.md) | PROJECT, PACKAGE, MODULE, HELP, TEMPLATE | New project → ship |

**Overview pages (shorter):** [01-CORE](01-CORE.md) … [11-TOOLING](11-TOOLING.md)

---

## How to pick a guide

```
No window yet?                 → GAME-LOOP-AND-RENDERING
3D objects?                    → ENTITY-SYSTEM
View / WASD?                   → CAMERA-AND-INPUT
Flat gray 3D?                  → LIGHTING
Walk on slopes (humanoid)?     → CHARACTER-3D-WALKING
Prototype mesh / GLB?          → MESHES-MODELS-MATERIALS
Ship art manifest?             → ASSETS-PIPELINE
Skeletal clips?                → ANIMATION
Fire / smoke?                  → PARTICLES
2D TMX level?                  → SPRITES-TILEMAPS-2D
Big outdoor map?               → TERRAIN-OPEN-WORLD
2D platforms, few objects?     → COLLISION-2D
2D crates and joints?          → PHYSICS-2D-PLATFORMER
3D physics / mouse pick?       → COLLISION-3D
Sound and music?               → AUDIO-FEEDBACK
Menus / HUD?                   → UI-MENUS
Save game slots?               → SAVE-AND-PROGRESS
config.json?                   → FILES-AND-JSON
2D dist / screen aim?          → math/MATH-2D-GAMEPLAY
3D XZ / yaw?                   → math/MATH-3D-GAMEPLAY
2D vectors?                    → math/VEC2-MATH
3D vectors?                    → math/VEC3-MATH
Smooth UI / camera?            → math/INTERPOLATION-AND-EASING
Angles / spin?                 → math/ANGLES-AND-ROTATION
Loot / dice?                   → math/RANDOMNESS-AND-PROCEDURE
Two moonrun processes?         → MULTIPLAYER
Custom net protocol?           → NETWORKING-LOW-LEVEL
--check failed?                → COMPILER-ERRORS
moonbasic new / package?       → PROJECT-WORKFLOW
```

**Rule:** Use the **simplest** tool that works, then upgrade when pain appears (manual rects → Box2D → 3D Jolt).

---

## Suggested learning paths

| Game type | Read order |
|-----------|------------|
| **3D demo cube** | [00-START](00-START.md) → GAME-LOOP → [01-CORE](01-CORE.md) → CAMERA-AND-INPUT → LIGHTING |
| **3D character action** | ENTITY → CAMERA-AND-INPUT → CHARACTER-3D-WALKING → COLLISION-3D → AUDIO |
| **2D platformer** | COLLISION-2D or PHYSICS-2D-PLATFORMER → SPRITES-TILEMAPS → UI-MENUS |
| **LAN co-op** | MULTIPLAYER → (optional) NETWORKING-LOW-LEVEL |
| **RPG with saves** | ENTITY → SAVE-AND-PROGRESS → ASSETS-PIPELINE → FILES-AND-JSON → UI-MENUS |
| **Ship to players** | PROJECT-WORKFLOW → COMPILER-ERRORS → DEBUG-AND-TESTING |

---

## See also

- [README.md](README.md) — 40-system build order
- [COMMAND_REGISTRY.md](COMMAND_REGISTRY.md) — grouped overload lists
- [../reference/BEGINNER_FULL_STACK.md](../reference/BEGINNER_FULL_STACK.md) — gameplay shortcuts
