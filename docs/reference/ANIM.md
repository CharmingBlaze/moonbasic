# Anim Commands

State-machine animation controller for entities with skeletal animations. Define named animation clips, add transition conditions, update each frame, and drive via parameters.

Requires an entity loaded with `ENTITY.LOADANIMATEDMESH` (or equivalent animated model).

## Core Workflow

1. `ANIM.DEFINE(entity, name, startFrame, endFrame, fps, loop)` — register a clip.
2. `ANIM.ADDTRANSITION(entity, fromState, toState, paramName)` — add a condition-based transition.
3. `ANIM.SETPARAM(entity, paramName, value)` — set a parameter to trigger transitions.
4. Each frame: `ANIM.UPDATE(entity, dt)` — advance the state machine.

---

## Define Clips

### `ANIM.DEFINE(entity, name, startFrame, endFrame, fps, loop)` 

Registers a named animation clip on an entity. `name` is the state name (e.g. `"idle"`, `"run"`). `startFrame` and `endFrame` are frame indices in the loaded animation data. `fps` is playback rate. `loop` is `TRUE` or `FALSE`.

---

## Transitions

### `ANIM.ADDTRANSITION(entity, fromState, toState, paramName)` 

Adds a transition from `fromState` to `toState` that fires when `paramName` is truthy (non-zero). `fromState` may be `"*"` to transition from any state.

---

## Parameters

### `ANIM.SETPARAM(entity, paramName, value)` 

Sets a named animation parameter. Use string `paramName` and any value — the transition system checks truthiness or equality against conditions.

---

## Update

### `ANIM.UPDATE(entity, dt)` 

Advances the animation state machine by `dt` seconds. Evaluates transitions and blends frames. Call every frame.

---

## Full Example

Character that transitions from idle to run when moving.

```basic
WINDOW.OPEN(960, 540, "Anim Demo")
WINDOW.SETFPS(60)

cam = CAMERA.CREATE()
CAMERA.SETPOS(cam, 0, 3, -7)
CAMERA.SETTARGET(cam, 0, 1.5, 0)

char = ENTITY.LOADANIMATEDMESH("assets/character.glb")
ENTITY.SETPOS(char, 0, 0, 0)

ANIM.DEFINE(char, "idle", 0,  59,  24, TRUE)
ANIM.DEFINE(char, "run",  60, 119, 30, TRUE)
ANIM.DEFINE(char, "jump", 120, 149, 24, FALSE)

ANIM.ADDTRANSITION(char, "idle", "run",  "moving")
ANIM.ADDTRANSITION(char, "run",  "idle", "idle")
ANIM.ADDTRANSITION(char, "*",    "jump", "jump")

px = 0.0

WHILE NOT WINDOW.SHOULDCLOSE()
    dt = TIME.DELTA()
    moving = 0

    IF INPUT.KEYDOWN(KEY_RIGHT) THEN px = px + 3 * dt : moving = 1
    IF INPUT.KEYDOWN(KEY_LEFT)  THEN px = px - 3 * dt : moving = 1
    IF INPUT.KEYPRESSED(KEY_SPACE) THEN
        ANIM.SETPARAM(char, "jump", 1)
    ELSE
        ANIM.SETPARAM(char, "jump", 0)
    END IF

    ANIM.SETPARAM(char, "moving", moving)
    ANIM.SETPARAM(char, "idle",   1 - moving)
    ANIM.UPDATE(char, dt)

    ENTITY.SETPOS(char, px, 0, 0)
    ENTITY.UPDATE(dt)

    RENDER.CLEAR(25, 30, 40)
    RENDER.BEGIN3D(cam)
        ENTITY.DRAWALL()
        DRAW3D.GRID(10, 1.0)
    RENDER.END3D()
    RENDER.FRAME()
WEND

ENTITY.FREE(char)
CAMERA.FREE(cam)
WINDOW.CLOSE()
```

---

## See also

- [ENTITY.md](ENTITY.md) — entity creation and `ENTITY.LOADANIMATEDMESH`
- [BTREE.md](BTREE.md) — behaviour trees driving animation params
- [PLAYER.md](PLAYER.md) — player with integrated animation
