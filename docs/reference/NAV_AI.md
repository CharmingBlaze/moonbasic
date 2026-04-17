# Navigation, steering, and behavior trees (`NAV.*`, `PATH.*`, `NAVAGENT.*`, `STEER.*`, `BTREE.*`)

Grid-based pathfinding on the **XZ** plane, lightweight steering forces as **Vec3 handles**, moving agents, and a small **behavior-tree** runner that calls your own **`FUNCTION`**s.

**Build:** These natives are registered only in **CGO** builds (`runtime/mbnav/register_cgo.go`). With **`CGO_ENABLED=0`**, every call fails with a stub error telling you to enable **CGO**.

**Related:** Model bounds for **`NAV.ADDTERRAIN`** / **`NAV.ADDOBSTACLE`** come from loaded models (**`MODEL.*`** handles).

Page shape: [DOC_STYLE_GUIDE.md](../DOC_STYLE_GUIDE.md) (**WAVE pattern** where sections use **`### \`SIGNATURE\``**).

## Core Workflow

**`NAV.MAKE()`** → **`NAV.SETGRID`** → **`NAV.ADDTERRAIN`** / **`NAV.ADDOBSTACLE`** (optional) → **`NAV.BUILD`** → **`NAV.FINDPATH`** → **`PATH.*`** queries → **`PATH.FREE`**. Spawn agents with **`NAVAGENT.CREATE(nav)`**, drive with **`NAVAGENT.MOVETO`** / **`NAVAGENT.UPDATE`**, or combine with **`STEER.*`** forces into **`NAVAGENT.APPLYFORCE`**. Behavior trees: **`BTREE.CREATE`** → **`BTREE.ADDCONDITION`** / **`BTREE.ADDACTION`** → **`BTREE.RUN`**.

---

## `NAV.*` — navigation grid

### `NAV.MAKE()` 

Creates a nav grid handle with a default **64×64** cell layout. Call **`NAV.SETGRID`** before serious use.

---

### `NAV.FREE(nav)` 

Frees the nav object.

---

### `NAV.SETGRID(nav, gw, gh, cellSize, originX, originZ)` 

Resizes the grid: **`gw`**/**`gh`** are tile counts (**1–4096**), **`cellSize`** is world units per cell (**> 0**), **`(originX, originZ)`** is the world origin of cell **`(0,0)`**. Clears blocked flags and marks the nav **not built** until **`NAV.BUILD`**.

---

### `NAV.ADDTERRAIN(nav, modelHandle)` / `NAV.ADDOBSTACLE(nav, modelHandle)` 

Uses the model’s **axis-aligned bounding box** in world space:

- **`ADDTERRAIN`** — marks that **XZ** footprint **walkable** and sets ground height from **`Min.Y`** across that region.
- **`ADDOBSTACLE`** — marks those cells **blocked**.

---

### `NAV.BUILD(nav)` 

Marks the nav data ready (sets the internal **built** flag). Call after editing terrain/obstacles and before path queries.

---

### `NAV.FINDPATH(nav, sx, sy, sz, tx, ty, tz)` 

Runs A* on the grid from start to target world position. Returns a **`Path`** handle (may be invalid if no path). Call **`PATH.FREE`** when done.

---

## `PATH.*` — path result

### `PATH.ISVALID(path)` → bool 

Whether the search produced a usable path.

---

### `PATH.NODECOUNT(path)` → int 

Number of waypoints (**0** if invalid).

---

### `PATH.NODEX(path, index)` / `PATH.NODEY(path, index)` / `PATH.NODEZ(path, index)` 

World coordinates of waypoint **`index`** (**0**-based). Errors if the index is out of range.

---

### `PATH.FREE(path)` 

Releases the path handle.

---

## `NAVAGENT.*` — agent on the nav mesh

Create with **`NAVAGENT.CREATE(nav)`** — ties the agent to that **`nav`** handle. **`NAVAGENT.MAKE`** is a deprecated alias.

| Command | Role |
|--------|------|
| **`NAVAGENT.SETPOS(agent, x, y, z)`** | Teleport position. |
| **`NAVAGENT.SETSPEED(agent, speed)`** | Max speed (≥ **0**). |
| **`NAVAGENT.SETMAXFORCE(agent, maxForce)`** | Caps acceleration from **`APPLYFORCE`**. |
| **`NAVAGENT.APPLYFORCE(agent, fx, fy, fz)`** | Adds to velocity, then clamps speed. |
| **`NAVAGENT.MOVETO(agent, tx, ty, tz)`** | Plans a path with **`NAV.FINDPATH`**; on success stores waypoints and clears velocity. |
| **`NAVAGENT.UPDATE(agent, dt)`** | Advances along waypoints at **`speed`**, or integrates velocity with damping when no path. |
| **`NAVAGENT.ISATDESTINATION(agent)`** → bool | **`TRUE`** when there is no active **`MOVETO`** destination. |
| **`NAVAGENT.X`** / **`NAVAGENT.Y`** / **`NAVAGENT.Z`** | Current position. |
| **`NAVAGENT.FREE(agent)`** | Frees the agent. |

---

## `STEER.*` — steering forces (Vec3 handles)

Steering helpers return **`VEC3`**-style handles (three floats) meant to be combined with **`NAVAGENT.APPLYFORCE`** or your own logic. Create groups with **`STEER.GROUPMAKE`**, add agents with **`STEER.GROUPADD(group, agent)`**.

| Command | Arguments | Result |
|--------|-----------|--------|
| **`STEER.SEEK`** | `(agent, tx, ty, tz)` | Vector toward target. |
| **`STEER.FLEE`** | `(agent, tx, ty, tz)` | Vector away from target. |
| **`STEER.ARRIVE`** | `(agent, tx, ty, tz, slowingRadius)` | Seek with speed ramp inside radius. |
| **`STEER.WANDER`** | `(agent, speed, jitterRadius)` | Pseudo-random direction from agent id. |
| **`STEER.FLOCK`** | `(selfAgent, group, cohesion, separation, alignment)` | Blended boids-style force. |
| **`STEER.AVOIDOBSTACLES`** | `(agent, radius)` | Repulsion from **blocked** nav cells near the agent. |
| **`STEER.FOLLOWPATH`** | `(agent, path)` | Seeks the nearest waypoint on the path. |

**`STEER.GROUPCLEAR(group)`** empties the group.

---

## `BTREE.*` — behavior tree (user functions)

### `BTREE.CREATE()` / `BTREE.FREE(bt)` 

Allocates a tree whose root is a **sequence** node. **`BTREE.MAKE`** is a deprecated alias of **`BTREE.CREATE`**.

---

### `BTREE.SEQUENCE(bt)` → handle 

Returns the same handle (reserved for fluent style; the runtime keeps a single root sequence).

---

### `BTREE.ADDCONDITION(bt, functionName)` / `BTREE.ADDACTION(bt, functionName)` 

Appends a child to the root **sequence**. On **`BTREE.RUN`**, children run in order:

- **Condition** / **action** — invokes the named **user function** with one argument: the **agent handle** passed to **`RUN`**.
- The function must return a value interpreted as boolean success for conditions; sequence stops on first failure.

---

### `BTREE.RUN(bt, agentHandle, dt)` 

Walks the tree; **`dt`** is reserved. User functions are resolved via the VM’s user-function invoker (same mechanism as **`SCENE.*`** loaders).

---

## Full Example (sketch)

```basic
; Pseudocode — requires CGO, loaded models, and valid grid setup
nav = NAV.MAKE()
NAV.SETGRID(nav, 64, 64, 1.0, 0.0, 0.0)
; ... NAV.ADDTERRAIN / NAV.ADDOBSTACLE with model handles ...
NAV.BUILD(nav)
path = NAV.FINDPATH(nav, x0, y0, z0, x1, y1, z1)
IF PATH.ISVALID(path) THEN
    PRINT PATH.NODECOUNT(path)
ENDIF
PATH.FREE(path)
NAV.FREE(nav)
```
