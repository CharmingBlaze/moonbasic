# Animation 3D Commands

Skeletal animation playback through the entity pipeline or manual model handles.

Page shape follows [DOC_STYLE_GUIDE.md](../DOC_STYLE_GUIDE.md) (**WAVE pattern**).

## Core Workflow

**Entity path (recommended):** load with `ENTITY.LOADANIMATEDMESH`, play with `ENTITY.PLAY` / `ENTITY.PLAYNAME`, advance with `ENTITY.UPDATE(dt)` each frame.

**Model path (manual):** load with `MODEL.LOADANIMATIONS`, play with `MODEL.PLAYIDX`, advance with `MODEL.UPDATEANIM(dt)`.

No built-in cross-fade — switching clips snaps to the new pose.

---

### `ENTITY.LOADANIMATEDMESH(path)`
Loads a 3D model with embedded skeletal animations.

- **Arguments**:
    - `path`: (String) File path to the model.
- **Returns**: (Integer) The new entity ID.
- **Example**:
    ```basic
    hero = ENTITY.LOADANIMATEDMESH("hero.glb")
    ```

---

### `ENTITY.PLAYNAME(entity, name)` / `PLAY`
Starts playing a specific animation clip by name or index.

- **Arguments**:
    - `entity`: (Integer) The character entity.
    - `name`: (String) Clip name from the model file.
- **Returns**: (Integer) The entity ID (for chaining).

---

### `ENTITY.UPDATE(dt)`
Advances the skeletal pose and bone matrices for all entities.

- **Returns**: (None)

---

### `ENTITY.FINDBONE(entity, name)`
Returns a "socket" entity that follows a specific bone's transform.

- **Arguments**:
    - `entity`: (Integer) The source entity.
    - `name`: (String) Bone name (e.g., "Hand_R").
- **Returns**: (Integer) A new hidden entity ID tracking the bone.

**Loop:** Call **`ENTITY.UPDATE(TIME.DELTA())`** each frame so physics and **`UpdateModelAnimation`** / bone matrices run before **`ENTITY.DRAW`** / **`DRAWALL`**.

---

## Path A — Entities (Blitz-style, recommended for games)

| Command | Role |
|--------|------|
| **`ENTITY.LOADANIMATEDMESH(path [, parent])`** | Loads geometry + embedded/bundled animations from a file Raylib supports (e.g. glTF). |
| **`ENTITY.UPDATE(dt)`** | Per frame: advances physics, then **skinned pose** for every entity with clips. Pass **`TIME.DELTA()`** (or your fixed step). |
| **`ENTITY.ANIMATE(entity [, mode, speed])`** | **`mode`**: **`0`–`1`** = loop over the active sub-range, **`2`** = ping-pong, **`3`+** = clamp/hold at end. **`speed`** scales how fast **`animTime`** advances (see below). |
| **`ENTITY.SETANIMINDEX(entity, idx)`** | Select clip by index; resets **`animTime`** to **0**. (**`ENTITY.PLAY`** does the same and sets playback speed to **1**.) |
| **`ENTITY.ANIMINDEX(entity)`** | Current clip index (pair with **`ENTITY.ANIMCOUNT`**). |
| **`ENTITY.ANIMCOUNT(entity)`** | Number of loaded clips. |
| **`ENTITY.EXTRACTANIMSEQ(entity, startFrame, endFrame)`** | Restrict playback to an **inclusive** frame sub-range **within** the current clip (useful when one file stores many actions in one timeline). |
| **`ENTITY.SETANIMTIME` / `ENTITY.ANIMTIME`** | Set or read continuous **animation time** (drives frame selection; not always an integer). |
| **`ENTITY.ANIMLENGTH(entity)`** | Length of the **current** clip in **frames** (Raylib **`FrameCount`** for that clip). |
| **`ENTITY.FINDBONE(entity, name)`** | Returns a **hidden** entity id whose transform follows a **bone** — parent props or weapons to it (see [ENTITY.md](ENTITY.md)). |

**Time scaling:** Inside **`ENTITY.UPDATE`**, internal time advances as **`animTime += dt * animSpeed * 30`**. So **`speed`** from **`ENTITY.ANIMATE`** is a multiplier on that default “~30 units per second” feel; tune **`speed`** and **`dt`** together for your asset.

**Drawing:** Use **`ENTITY.DRAWALL`** for sorted scene draw, or **`ENTITY.DRAW`** / **`DrawEntity`** for a single id, after **`ENTITY.UPDATE`** so bones and sockets match the updated pose.

---

## Path B — `MODEL` handles (manual playback)

For a **`MODEL.LOAD`** handle (not an entity id), load clips from a **sidecar** or shared file Raylib accepts:

| Command | Role |
|--------|------|
| **`MODEL.LOADANIMATIONS(model, path)`** | Load/replace animation set; previous set is unloaded. |
| **`MODEL.PLAYIDX(model, idx)`** | Start clip **`idx`** from frame **0**. |
| **`MODEL.UPDATEANIM(model, dt)`** | Advance playback; call each frame with **`TIME.DELTA()`** while playing. |
| **`MODEL.STOP` / `MODEL.LOOP` / `MODEL.SETSPEED`** | Stop, toggle loop, speed multiplier (playback advances at ~60 base FPS × speed × **`dt`**). |
| **`MODEL.ANIMCOUNT` / `MODEL.ANIMNAME(model, idx)`** | Introspection. |
| **`MODEL.GETFRAME` / `MODEL.TOTALFRAMES`** | Current frame index and length of **active** clip. |
| **`MODEL.ISPLAYING` / `MODEL.ANIMDONE`** | One-shot clips: **`ANIMDONE`** when non-loop playback has reached the end. |

**`MODEL.PLAY(model, name)`** is **not** implemented — use **`PLAYIDX`** after resolving names with **`ANIMNAME`**.

Details and pitfalls: [MODEL.md](MODEL.md) (animation subsection).

---

## Choosing a path

| Use **entities** when… | Use **`MODEL.*`** when… |
|------------------------|-------------------------|
| You already use **`ENTITY.POSITION`**, **`ENTITY.PARENT`**, **`ENTITY.DRAWALL`**, collisions, or **bone sockets**. | You draw with **`MODEL.DRAW`** / custom transforms and do not need the entity store. |
| You want **one** **`ENTITY.UPDATE`** to drive motion + skinning. | You prefer explicit **`MODEL.UPDATEANIM`** next to your own transform code. |

---

## Limitations (current runtime)

- **Single active pose** per mesh — no skeletal **cross-fade** between clips (Raylib updates one clip + frame index).
- **Root motion** is not extracted from clips; move the **entity** or **model** yourself if your asset bakes displacement into bones only.
- **Clip naming** on the **`MODEL`** path: resolve by index with **`MODEL.ANIMNAME`**; there is no string-based **`PLAY`** yet.

---

## Full Example

```basic
WINDOW.OPEN(1280, 720, "Animation Demo")
WINDOW.SETFPS(60)

cam = CAMERA.CREATE()
CAMERA.SETPOS(cam, 0, 3, -6)
CAMERA.SETTARGET(cam, 0, 1, 0)

hero = ENTITY.LOADANIMATEDMESH("character.glb")
ENTITY.SETPOS(hero, 0, 0, 0, TRUE)
ENTITY.PLAYNAME(hero, "walk")

WHILE NOT WINDOW.SHOULDCLOSE()
    ENTITY.UPDATE(TIME.DELTA())

    RENDER.CLEAR(40, 50, 60)
    RENDER.BEGIN3D(cam)
        ENTITY.DRAWALL()
    RENDER.END3D()
    RENDER.FRAME()
WEND

ENTITY.FREE(hero)
CAMERA.FREE(cam)
WINDOW.CLOSE()
```

---

## See also

- [ENTITY.md](ENTITY.md) — bone sockets, brushes, **`ENTITY.DRAWALL`**
- [MODEL.md](MODEL.md) — **`MODEL.LOAD`**, materials, **`MODEL.UPDATEANIM`**
- [MEMORY.md](../MEMORY.md) — ownership of **`LoadModelAnimations`**
- [GAME_ENGINE_PATTERNS.md](GAME_ENGINE_PATTERNS.md) — entity + camera patterns
