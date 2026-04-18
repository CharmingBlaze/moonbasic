# NavAgent Commands

Autonomous navigation agents that follow paths on a navmesh or grid. Agents use steering to reach destinations, with configurable speed and steering force.

See [NAV_AI.md](NAV_AI.md) and [NAVMESH.md](NAVMESH.md) for navmesh construction.

## Core Workflow

1. Build a navmesh or grid first (see [NAVMESH.md](NAVMESH.md) / [GRID.md](GRID.md)).
2. `NAVAGENT.CREATE(navHandle)` — attach an agent to the nav structure.
3. `NAVAGENT.MOVETO(agent, x, y, z)` — set destination.
4. Each frame: `NAVAGENT.UPDATE(agent, dt)` — advance movement.
5. Read `NAVAGENT.X/Y/Z` to sync visuals.
6. `NAVAGENT.FREE(agent)` when done.

---

## Creation

### `NAVAGENT.CREATE(navHandle)` 

Creates a navigation agent attached to `navHandle` (a navmesh or grid handle). Returns an **agent handle**.

---

## Position

### `NAVAGENT.SETPOS(agent, x, y, z)` 

Teleports the agent to world position `(x, y, z)`.

- *Handle shortcut*: `agent.setPos(x, y, z)`

---

### `NAVAGENT.GETPOS(agent)` 

Returns the current position as a 3-element array `[x, y, z]`.

---

### `NAVAGENT.X(agent)` / `NAVAGENT.Y(agent)` / `NAVAGENT.Z(agent)` 

Returns the individual world coordinate as a float scalar.

---

## Rotation

### `NAVAGENT.SETROT(agent, yawDegrees)` 

Manually overrides the agent Y rotation in degrees.

---

### `NAVAGENT.GETROT(agent)` 

Returns approximate `[pitch, yaw, roll]` in radians derived from path waypoint tangent or steering velocity.

---

## Navigation

### `NAVAGENT.MOVETO(agent, x, y, z)` 

Sets the destination. The agent will steer toward it on each `UPDATE`.

- *Handle shortcut*: `agent.moveTo(x, y, z)`

---

### `NAVAGENT.STOP(agent)` 

Clears the current path and stops the agent. Returns the agent handle.

---

### `NAVAGENT.ISATDESTINATION(agent)` 

Returns `TRUE` when the agent has reached its current destination (within arrival tolerance).

---

## Speed & Force

### `NAVAGENT.SETSPEED(agent, speed)` 

Sets the movement speed in world units per second.

---

### `NAVAGENT.GETSPEED(agent)` 

Returns the current speed.

---

### `NAVAGENT.SETMAXFORCE(agent, force)` 

Sets the maximum steering force. Higher values = sharper turns.

---

### `NAVAGENT.GETMAXFORCE(agent)` 

Returns the current max steering force.

---

### `NAVAGENT.APPLYFORCE(agent, fx, fy, fz)` 

Applies an external steering force this frame (e.g. for avoidance or knockback).

---

## Update

### `NAVAGENT.UPDATE(agent, dt)` 

Advances the agent along its path by `dt` seconds. Call every frame.

---

## Lifetime

### `NAVAGENT.FREE(agent)` 

Destroys the agent handle.

---

## Full Example

An enemy agent patrolling between two waypoints on a navmesh.

```basic
WINDOW.OPEN(960, 540, "NavAgent Demo")
WINDOW.SETFPS(60)

PHYSICS3D.START()

; build nav grid (see NAVMESH.md for full setup)
nav = NAVMESH.BUILD("level.glb")

enemy     = NAVAGENT.CREATE(nav)
enemyMesh = ENTITY.CREATECUBE(1.0)
NAVAGENT.SETPOS(enemy, -5, 0, 0)
NAVAGENT.SETSPEED(enemy, 3.0)

waypoints = ARRAY.MAKE(2)
ARRAY.SET(waypoints, 0, 5)   ; target A x
ARRAY.SET(waypoints, 1, -5)  ; target B x
wpIndex = 0

NAVAGENT.MOVETO(enemy, 5, 0, 0)

cam = CAMERA.CREATE()
CAMERA.SETPOS(cam, 0, 12, -12)
CAMERA.SETTARGET(cam, 0, 0, 0)

WHILE NOT WINDOW.SHOULDCLOSE()
    dt = TIME.DELTA()
    NAVAGENT.UPDATE(enemy, dt)

    IF NAVAGENT.ISATDESTINATION(enemy) THEN
        wpIndex = 1 - wpIndex
        tx = ARRAY.GET(waypoints, wpIndex)
        NAVAGENT.MOVETO(enemy, tx, 0, 0)
    END IF

    ex = NAVAGENT.X(enemy)
    ey = NAVAGENT.Y(enemy)
    ez = NAVAGENT.Z(enemy)
    ENTITY.SETPOS(enemyMesh, ex, ey, ez)
    ENTITY.UPDATE(dt)

    RENDER.CLEAR(20, 25, 35)
    RENDER.BEGIN3D(cam)
        ENTITY.DRAWALL()
        DRAW3D.GRID(20, 1.0)
    RENDER.END3D()
    RENDER.FRAME()
WEND

NAVAGENT.FREE(enemy)
ENTITY.FREE(enemyMesh)
PHYSICS3D.STOP()
WINDOW.CLOSE()
```

---

## Extended Command Reference

| Command | Description |
|--------|-------------|
| `NAVAGENT.SETPOSITION(agent, x,y,z)` | Alias of `NAVAGENT.SETPOS` — teleport agent to position. |

---

## See also

- [NAVMESH.md](NAVMESH.md) — navmesh construction and pathfinding
- [NAV_AI.md](NAV_AI.md) — AI overview and nav grid
- [STEER.md](STEER.md) — steering behaviour forces
- [BTREE.md](BTREE.md) — behaviour trees for AI logic
