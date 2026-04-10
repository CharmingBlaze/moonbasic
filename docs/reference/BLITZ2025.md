# Blitz-style “2025” API map

MoonBASIC uses **dotted command names** (`ENTITY.*`, `CAMERA.*`, …). This page maps common Blitz-style names to **built-in equivalents** (aliases share the same implementation; no duplicate semantics).

For a **line-by-line checklist** of classic Blitz3D/BlitzPlus command names (Plot, CreateCube, CameraZoom, LoadTexture, …) and their moonBASIC counterparts, see **[BLITZ_COMMAND_INDEX.md](BLITZ_COMMAND_INDEX.md)**.

---

## 1. Scene and world

| Intent | MoonBASIC command | Notes |
|--------|-------------------|--------|
| Save entities to JSON | **`ENTITY.SAVESCENE(path)`** or **`SCENE.SAVESCENE`** | Same handler; format v1 (see below). |
| Load entities from JSON | **`ENTITY.LOADSCENE(path)`** or **`SCENE.LOADSCENE`** | Clears the world first, then spawns entities. |
| Clear all entities | **`ENTITY.CLEARSCENE`** or **`SCENE.CLEARSCENE`** | Resets entity ids and groups. |
| Scene flow / transitions | **`SCENE.REGISTER`**, **`SCENE.LOAD`**, … | Existing [scene](SCENE.md) system for scripted scene switches. |

**Scene file format (v1):** JSON with `v: 1` and `e: [...]` entity records (kind, transform, color, optional `path` for models). Primitives and file-backed models reload; procedural meshes without a path are recreated as a default cube mesh.

**Groups (runtime only):** **`ENTITY.GROUPCREATE`**, **`GROUPADD`**, **`GROUPREMOVE`**. **`ENTITY.ENTITIESINGROUP`** returns an **array handle** of entity ids, or **NULL** when empty.

**Spatial queries:** **`ENTITY.ENTITIESINRADIUS`**, **`ENTITY.ENTITIESINBOX`** — return a **1D float array** of ids (use **`ARRAYLEN`** / indexing). **NULL** means no matches.

---

## 2. Math and vectors

| Intent | MoonBASIC |
|--------|-----------|
| Vec3(x,y,z) | **`VEC3.MAKE`** or alias **`VEC3.VEC3`** |
| VecAdd / Sub / Scale / Normalize / Dot / Cross / Length | **`VEC3.ADD`**, **`VEC3.SUB`**, **`VEC3.MUL`**, **`VEC3.NORMALIZE`**, **`VEC3.DOT`**, **`VEC3.CROSS`**, **`VEC3.LENGTH`** — or **`VEC3.VECADD`**, **`VECSUB`**, **`VECSCALE`**, … |
| WrapAngle, Lerp, SmoothStep | Top-level **`WRAPANGLE`**, **`LERP`**, **`SMOOTHSTEP`** or **`MATH.*`** variants (see manifest). |

---

## 3. Camera

| Intent | MoonBASIC |
|--------|-----------|
| Set target to point | **`CAMERA.SETTARGET`** |
| Set target to entity | **`CAMERA.SETTARGETENTITY(cam, entity)`** |
| Third-person follow | **`CAMERA.FOLLOW`** (world target) or **`CAMERA.FOLLOWENTITY`** / **`CAMERA.CAMERAFOLLOW`** (entity + dist/height/smooth) |
| Orbit | **`CAMERA.SETORBIT`**, **`CAMERA.ORBIT`**, **`CAMERA.ORBITENTITY`** |
| Zoom FOV | **`CAMERA.ZOOM`** |
| Screen ray / pick | **`CAMERA.GETRAY(cam, sx, sy)`** or alias **`CAMERA.PICK`** (same as **`GETRAY`**) |
| Shake | **`CAMERA.SHAKE(cam, amount, duration)`** — applied during **`CAMERA.BEGIN`** |

---

## 4. Physics (Jolt, Linux + CGO)

| Intent | MoonBASIC |
|--------|-----------|
| World gravity | **`PHYSICS3D.SETGRAVITY`** or alias **`PHYSICS.SETGRAVITY`** |
| Ray cast | **`PHYSICS3D.RAYCAST`** or **`PHYSICS.RAYCAST`** |
| Character | **`CHARCONTROLLER.MAKE`** or **`CONTROLLER.CREATE`**; **`MOVE`**, **`ISGROUNDED`** / **`CONTROLLER.GROUNDED`**; **`FREE`** |

**`PHYSICS.SPHERECAST` / `BOXCAST` / `ENABLE` / `DISABLE`:** reserved stubs with messages directing to **`RAYCAST`** / **`BODY3D.ACTIVATE`** / **`DEACTIVATE`**. **`CONTROLLER.JUMP`** is not implemented on this controller yet.

---

## 5. Rendering (fog, lights, materials)

| Intent | MoonBASIC |
|--------|-----------|
| Fog color / range | **`FOG.SETCOLOR`**, **`FOG.SETNEAR` / `SETFAR`** or **`FOG.SETRANGE(near, far)`** |
| Lights | **`LIGHT.MAKE`**, **`LIGHT.SETCOLOR`**, **`LIGHT.SETRANGE`**, cones, … |
| Skybox | **`RENDER.SETSKYBOX`** (path) — see [rendering docs](../BUILDING.md) / manifest |
| Material | **`MATERIAL.MAKEDEFAULT`** or **`MATERIAL.CREATE`** (alias) |

---

## 6. Input

| Intent | MoonBASIC |
|--------|-----------|
| Mouse position | **`INPUT.MOUSEX`**, **`INPUT.MOUSEY`** |
| Mouse delta | **`INPUT.MOUSEXSPEED`**, **`INPUT.MOUSEYSPEED`** (Blitz aliases) or **`INPUT.MOUSEDELTAX/Y`** |
| Mouse down / hit | **`INPUT.MOUSEDOWN`**, **`INPUT.MOUSEHIT`** (pressed this frame) |
| Gamepad | **`INPUT.JOYX`**, **`JOYY`**, **`JOYBUTTON`**, **`INPUT.JOYDOWN`** (alias of **`JOYBUTTON`**) |

---

## 7. Audio

Use **`AUDIO.LOADSOUND`**, **`AUDIO.PLAY`**, **`AUDIO.STOP`**, **`AUDIO.SETSOUNDVOLUME`**, etc. Raylib 2D panning exists; full **3D positional** mixing is not exposed as separate builtins here—use engine docs for extensions.

---

## 8. Timing

**`TIME.GET`**, **`TIME.DELTA`**, **`TIMER`** (monotonic seconds). A full **`TimerCreate` / `TimerWait`** coroutine-style API is not implemented; use **`TIME.GET`** and comparisons in your update loop.

---

## 9. Files and JSON

| Intent | MoonBASIC |
|--------|-----------|
| Exists | **`FILE.EXISTS(path)`** |
| Read whole text | **`FILE.READALLTEXT(path)`** |
| Write whole text | **`FILE.WRITEALLTEXT(path, text)`** |
| Load JSON from file | **`JSON.PARSE(path)`** or **`JSON.LOADFILE`** (same behavior: path on disk) |
| Save JSON | **`JSON.TOFILE`** or **`JSON.SAVEFILE`** `(jsonHandle, path)` |

---

## 10. Debug

| Intent | MoonBASIC |
|--------|-----------|
| Log | **`DEBUG.LOG`**, **`DEBUG.PRINT`** |
| Draw line / box (3D, inside **`CAMERA.Begin`/`End`**) | **`DEBUG.DRAWLINE`**, **`DEBUG.DRAWBOX`** |

---

## Related

- [BLITZ3D.md](BLITZ3D.md) — core Blitz-style entity/camera/input aliases  
- [CAMERA.md](CAMERA.md), [INPUT.md](INPUT.md), [PHYSICS3D.md](PHYSICS3D.md) where present  
