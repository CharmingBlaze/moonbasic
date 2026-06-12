# moonBASIC API Reference

Complete reference for all moonBASIC commands organized by namespace. Every command uses the `Namespace.Method()` pattern as the canonical API style.

---

## Quick Start

```basic
Window.Open(1280, 720, "My Game")
Window.SetFPS(60)

WHILE NOT Window.ShouldClose()
    Render.Clear(25, 25, 40)
    Draw.Text("Hello moonBASIC!", 10, 10, 24, 255, 255, 255, 255)
    Render.Frame()
WEND

Window.Close()
```

---

## Core Namespaces

| Namespace | Description | Doc |
|-----------|-------------|-----|
| **WINDOW** | Window creation, FPS, fullscreen, monitors | [WINDOW.md](WINDOW.md) |
| **RENDER** | Frame lifecycle, clear, blend, depth, screenshots, post-processing | [RENDER.md](RENDER.md) |
| **DRAW** | Immediate-mode 2D shapes, text, and 3D primitives | [DRAW.md](DRAW.md) |
| **INPUT** | Keyboard, mouse, gamepad, touch, gestures | [INPUT.md](INPUT.md) |
| **CAMERA** | 3D/2D cameras, orbit, FPS mode, shake, culling | [CAMERA.md](CAMERA.md) |
| **AUDIO** | Sounds, music, streams, waves, spatial 3D audio | [AUDIO.md](AUDIO.md) |

## World & Scene

| Namespace | Description | Doc |
|-----------|-------------|-----|
| **WORLD** | Gravity, fog, sky, streaming, time scale, flash/hitstop | [WORLD.md](WORLD.md) |
| **ENTITY** | Game entities with position, rotation, physics, rendering | [ENTITY.md](ENTITY.md) |
| **TERRAIN** | Heightmap terrain generation and chunk streaming | TERRAIN.md |
| **LIGHT** | Point, directional, and spot lights | LIGHT.md |
| **SKY** | Procedural sky with day/night cycle | SKY.md |
| **WEATHER** | Rain, snow, fog systems | WEATHER.md |
| **CLOUD** | Volumetric cloud rendering | CLOUD.md |
| **WATER** | Water plane with waves and reflection | WATER.md |

## 3D Assets

| Namespace | Description | Doc |
|-----------|-------------|-----|
| **MODEL** | 3D model loading, rendering, animation, LOD, instancing | [MODEL.md](MODEL.md) |
| **MESH** | Low-level mesh creation, operations, and queries | [MODEL.md](MODEL.md#mesh-commands) |
| **MATERIAL** | Material properties, textures, shaders | [MODEL.md](MODEL.md#material-commands) |
| **TEXTURE** | GPU texture loading, async, reload | [TEXTURE.md](TEXTURE.md) |
| **SHADER** | Custom shader loading, uniforms, preset constants | [MODEL.md](MODEL.md#shader-commands) |
| **TRANSFORM** | 3D transform matrices | [MATH.md](MATH.md#vec2--vec3--mat4) |
| **ANIM** | Skeletal animation playback | [MODEL.md](MODEL.md#skeletal-animation) |

## 2D Graphics

| Namespace | Description | Doc |
|-----------|-------------|-----|
| **SPRITE** | 2D sprite creation, animation, layers | [SPRITE.md](SPRITE.md) |
| **FONT** | Custom font loading and text rendering | FONT.md |
| **ATLAS** | Texture atlas / spritesheet management | ATLAS.md |
| **TILEMAP** | Tile-based map rendering | TILEMAP.md |
| **GUI** | raygui-based UI widgets | [GUI.md](GUI.md) |

## Physics

| Namespace | Description | Doc |
|-----------|-------------|-----|
| **PHYSICS3D** | Jolt Physics 3D simulation | [PHYSICS.md](PHYSICS.md) |
| **PHYSICS2D** | Box2D-style 2D physics | [PHYSICS.md](PHYSICS.md#2d-physics) |
| **BODY3D** | 3D rigid body creation and properties | [PHYSICS.md](PHYSICS.md#3d-body-creation) |
| **BODY2D** | 2D rigid body creation and properties | [PHYSICS.md](PHYSICS.md#2d-physics) |
| **SHAPE** | Collision shape creation (box, sphere, capsule) | [PHYSICS.md](PHYSICS.md#collision-shapes) |
| **JOINT** / **JOINT3D** | Physics joint constraints | [PHYSICS.md](PHYSICS.md#joints) |
| **PICK** | Ray picking and collision queries | [PHYSICS.md](PHYSICS.md#raycasting) |

## Player & Character

| Namespace | Description | Doc |
|-----------|-------------|-----|
| **PLAYER** | Player character controller (Jolt KCC) | PLAYER.md |
| **CHARCONTROLLER** | Generic character controller | CHARACTER_PHYSICS.md |
| **NAV** / **NAVAGENT** | Navigation mesh and pathfinding | NAV.md |

## Particles & Effects

| Namespace | Description | Doc |
|-----------|-------------|-----|
| **PARTICLES** | 3D particle emitter system | [PARTICLES.md](PARTICLES.md) |
| **PARTICLE2D** | 2D particle effects | [PARTICLES.md](PARTICLES.md) |
| **DECAL** | Projected decals (bullet holes, blood splatters) | DECAL.md |
| **POST** | Post-processing pipeline (bloom, vignette, etc.) | [POST.md](POST.md) |
| **EFFECT** | Screen effects (SSAO, motion blur, DOF, FXAA) | [POST.md](POST.md) |
| **TWEEN** | Value interpolation and animation | TWEEN.md |

## Data & Utility

| Namespace | Description | Doc |
|-----------|-------------|-----|
| **JSON** | JSON parsing, creation, querying, serialization | [JSON.md](JSON.md) |
| **FILE** | File I/O (read, write, exists, list, directories) | [FILE.md](FILE.md) |
| **STRING** | String manipulation (upper, lower, mid, replace) | STRING.md |
| **MATH** | Math functions, easing, noise, CurveValue, vectors | [MATH.md](MATH.md) |
| **TIME** | Delta time, FPS, milliseconds, timers | [TIME.md](TIME.md) |
| **TIMER** | Countdown and simulation timers | [TIME.md](TIME.md#countdown-timers) |
| **STOPWATCH** | Elapsed time measurement | [TIME.md](TIME.md#stopwatches) |
| **RAND** | Random number generation | RAND.md |
| **NOISE** | Perlin and simplex noise generation | NOISE.md |
| **COLOR** | Color creation and manipulation | COLOR.md |
| **VEC2** / **VEC3** / **QUAT** / **MAT4** | Vector and matrix math | MATH.md |
| **TABLE** | Data table / dictionary | TABLE.md |

## Networking

| Namespace | Description | Doc |
|-----------|-------------|-----|
| **NET** | TCP/UDP networking | NET.md |
| **ENET** | ENet reliable UDP (multiplayer) | NET.md |
| **SERVER** / **CLIENT** | High-level multiplayer | NET.md |
| **RPC** / **LOBBY** / **PEER** | Multiplayer subsystems | NET.md |

## System

| Namespace | Description | Doc |
|-----------|-------------|-----|
| **CONSOLE** | On-screen debug console | SYSTEM.md |
| **SYSTEM** | Performance monitor, clipboard | SYSTEM.md |
| **DEBUG** | Debug drawing and logging | SYSTEM.md |
| **EVENT** | Automation event recording/replay | EVENT.md |
| **COMPUTESHADER** | GPU compute shaders | COMPUTE.md |
| **RENDERTARGET** | Off-screen render targets | RENDERTARGET.md |

## Blitz3D Compatibility

| Shortcut | Maps To |
|----------|---------|
| `Graphics3D(w, h, d)` | `Window.Open(w, h, "")` |
| `CreateCamera()` | `Camera.Create()` |
| `CreateLight()` | `Light.Create()` |
| `LoadMesh(path)` | `Model.Load(path)` |
| `CreateCube(...)` | `Model.CreateCube(...)` |
| `PositionEntity(e, x, y, z)` | `Entity.SetPos(e, x, y, z)` |
| `RotateEntity(e, p, y, r)` | `Entity.SetRot(e, p, y, r)` |
| `ScaleEntity(e, sx, sy, sz)` | `Entity.SetScale(e, sx, sy, sz)` |
| `EntityColor(e, r, g, b)` | `Entity.SetCol(e, r, g, b)` |
| `EntityAlpha(e, a)` | `Entity.SetAlpha(e, a)` |
| `HideEntity(e)` | `Entity.Hide(e)` |
| `ShowEntity(e)` | `Entity.Show(e)` |
| `FreeEntity(e)` | `Entity.Free(e)` |
| `EntityX(e)` | `Entity.GetPos(e)[0]` |
| `EntityY(e)` | `Entity.GetPos(e)[1]` |
| `EntityZ(e)` | `Entity.GetPos(e)[2]` |
| `MoveEntity(e, f, r, u)` | `Entity.Move(e, f, r, u)` |
| `TurnEntity(e, p, y, r)` | `Entity.Turn(e, p, y, r)` |
| `KeyDown(code)` | `Input.KeyDown(code)` |
| `KeyHit(code)` | `Input.KeyPressed(code)` |
| `MouseDown(btn)` | `Input.MouseDown(btn)` |
| `MouseHit(btn)` | `Input.MouseHit(btn)` |
| `MOUSEDX` | `Input.MouseDeltaX()` |
| `MOUSEDY` | `Input.MouseDeltaY()` |
| `UpdateWorld` | `World.Update(dt)` |
| `RenderWorld` | `Camera.Begin + Entity.DrawAll + Camera.End` |
| `Flip` | `Render.Frame()` |
| `Cls` | `Render.Clear()` |
| `DELTATIME` | `Time.Delta()` |
| `MILLISECS` | `Time.Millisecs()` |

---

## API Style Guide

moonBASIC commands are **case-insensitive**. All of these are equivalent:

```basic
Window.Open(1280, 720, "Game")
WINDOW.OPEN(1280, 720, "Game")
window.open(1280, 720, "Game")
```

### Recommended Style

```basic
; Use Namespace.Method — clear, readable, self-documenting
cam = Camera.Create()
cam.pos(0, 10, 20)
cam.look(0, 0, 0)
cam.fov(60)
```

### Naming Conventions

- **Variables** — `camelCase` (e.g., `playerModel`, `mainCamera`)
- **Constants** — `SCREAMING_SNAKE` (e.g., `MAX_SPEED`, `GRAVITY`)
- **Types** — `PascalCase` (e.g., `PlayerData`, `EnemyStats`)
- **Commands** — `Namespace.Method` (canonical)

### Creation Pattern

Use `CREATE` for all object instantiation:

```basic
cam = Camera.Create()
body = Body3D.Create()
light = Light.CreatePoint()
timer = Timer.Create(5.0)
```

### Resource Cleanup

Every `CREATE` or `LOAD` has a matching `FREE`:

```basic
tex = Texture.Load("sprite.png")
; ... use texture ...
Texture.Free(tex)
```

---

## Engine Architecture

moonBASIC compiles `.mb` source files to register-based IR v3 bytecode, which runs on a stack-less virtual machine. The runtime is written in Go and uses:

- **Raylib 5.5** — Windowing, rendering, input, audio (via CGO or purego DLL)
- **Jolt Physics** — 3D rigid body simulation (via vendored jolt-go)
- **ENet** — Reliable UDP networking
- **raygui** — Immediate-mode GUI widgets

The compiler pipeline: `Lexer → Parser → Semantic Analyzer → CodeGen → IR v3 Bytecode → VM`

All commands are registered as Go functions in `runtime/*` packages. The manifest in `compiler/builtinmanifest/commands.json` contains 3000+ command entries that the compiler validates at parse time.
