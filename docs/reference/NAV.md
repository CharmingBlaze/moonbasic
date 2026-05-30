# Nav Commands

Navmesh and navigation grid construction, pathfinding, and AI movement integration. `NAV.*` builds and queries navigation data; [NAVAGENT.md](NAVAGENT.md) drives agents along the resulting paths.

## Core Workflow

**Navmesh (3D):**
1. `NAV.CREATE()` — allocate a navmesh handle.
2. `NAV.SETGRID` / `NAV.ADDTERRAIN` / `NAV.ADDOBSTACLE` — define the walkable world.
3. `NAV.BUILD(nav)` — bake navigation data.
4. `NAV.GETPATH(nav, sx, sy, sz, tx, ty, tz)` or `NAV.FINDPATH(...)` — query a path.
5. Read path with `PATH.*` commands.

**Grid (2D/top-down):**
1. `NAV.CREATE()` → returns an int id.
2. `NAV.SETGRID(id, cols, rows, cellSize, ox, oz)` — define the grid.
3. `NAV.BUILD(id)` — bake.
4. `NAV.GOTO(entityId, x, y, z)` to move entities directly.
5. `NAV.UPDATE(entityId)` each frame.

---

## Creation

### `NAV.CREATE()` 

Creates a new navmesh (returns handle) or a new nav grid (returns int id, depending on overload).

---

### `NAV.FREE(navHandle)` 

Frees a navmesh handle.

---

## Grid Setup

### `NAV.SETGRID(nav, cols, rows, cellSize, ox, oy, oz)` 

Defines grid dimensions: `cols` × `rows` cells of `cellSize` world units, with origin offset `(ox, oy, oz)`.

---

### `NAV.ADDTERRAIN(nav, meshHandle)` 

Registers a terrain mesh as walkable surface.

---

### `NAV.ADDOBSTACLE(nav, entityHandle)` 

Marks an entity's bounding box as a navigation obstacle.

---

### `NAV.BUILD(nav)` 

Bakes the navigation data. Must be called after all terrain/obstacles are added.

---

### `NAV.BAKE(nav, agentRadius, agentHeight)` 

Alternate bake command accepting agent dimensions for clearance calculations. Returns the nav handle.

---

## Path Queries

### `NAV.FINDPATH(nav, sx, sy, sz, tx, ty, tz)` 

Returns a **path handle** from `(sx, sy, sz)` to `(tx, ty, tz)`. Read with `PATH.*`.

---

### `NAV.GETPATH(nav, sx, sy, sz, tx, ty, tz)` 

Alias of `FINDPATH`. Returns a **path handle**.

---

### `NAV.ISREACHABLE(nav, sx, sy, sz, tx, ty, tz)` 

Returns `TRUE` if `(tx, ty, tz)` is reachable from `(sx, sy, sz)`.

---

## Entity Movement (Grid)

### `NAV.GOTO(entityId, x, y, z [, speed [, tolerance]])` 

Moves `entityId` toward world position on the grid. Optional `speed` and arrival `tolerance`.

---

### `NAV.UPDATE(entityId)` 

Advances `entityId` along its current path. Call each frame.

---

### `NAV.CHASE(entityId, targetId, speed, stopDist)` 

Moves `entityId` to follow `targetId` at `speed`, stopping within `stopDist`.

---

### `NAV.PATROL(entityId, x1, z1, x2, z2, speed)` 

Makes `entityId` patrol between two points at `speed`.

---

## Debug

### `NAV.DEBUGDRAW(nav)` 

Renders the navmesh or grid as wireframe for debugging. Call inside `RENDER.BEGIN3D` / `RENDER.END3D`.

---

## Full Example

Enemy chasing the player on a navmesh.

```basic
WINDOW.OPEN(960, 540, "Nav Demo")
WINDOW.SETFPS(60)

PHYSICS3D.START()

nav = NAV.CREATE()
NAV.SETGRID(nav, 32, 32, 1.0, -16, 0, -16)
NAV.BUILD(nav)

player = ENTITY.CREATECUBE(0.8)
ENTITY.SETPOS(player, 0, 0.4, 0)

enemy = ENTITY.CREATECUBE(0.8)
ENTITY.SETCOLOR(enemy, 255, 80, 80)
ENTITY.SETPOS(enemy, 8, 0.4, 8)

cam = CAMERA.CREATE()
CAMERA.SETPOS(cam, 0, 20, -5)
CAMERA.SETTARGET(cam, 0, 0, 5)

px = 0.0 : pz = 0.0

WHILE NOT WINDOW.SHOULDCLOSE()
    dt = TIME.DELTA()

    IF INPUT.KEYDOWN(KEY_D) THEN px = px + 4 * dt
    IF INPUT.KEYDOWN(KEY_A) THEN px = px - 4 * dt
    IF INPUT.KEYDOWN(KEY_S) THEN pz = pz + 4 * dt
    IF INPUT.KEYDOWN(KEY_W) THEN pz = pz - 4 * dt
    ENTITY.SETPOS(player, px, 0.4, pz)

    NAV.CHASE(enemy, player, 2.5, 0.5)
    NAV.UPDATE(enemy)

    ENTITY.UPDATE(dt)
    PHYSICS3D.UPDATE()

    RENDER.CLEAR(20, 25, 35)
    RENDER.BEGIN3D(cam)
        ENTITY.DRAWALL()
        NAV.DEBUGDRAW(nav)
        DRAW.GRID(32, 1.0)
    RENDER.END3D()
    RENDER.FRAME()
WEND

NAV.FREE(nav)
PHYSICS3D.STOP()
WINDOW.CLOSE()
```

---

## See also

- [NAVAGENT.md](NAVAGENT.md) — handle-based nav agents
- [PATH.md](PATH.md) — reading path waypoints
- [STEER.md](STEER.md) — steering forces
- [BTREE.md](BTREE.md) — AI behaviour trees
