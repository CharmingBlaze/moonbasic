# Navigation, steering, and behavior trees (`NAV.*`, `PATH.*`, `NAVAGENT.*`, `STEER.*`, `BTREE.*`)

Grid-based pathfinding on the XZ plane, lightweight steering forces as **Vec3 handles**, moving agents, and a small **behavior-tree** runner that calls your own `FUNCTION`s.

**Build:** These natives are registered only in **CGO** builds (`runtime/mbnav/register_cgo.go`). With `CGO_ENABLED=0`, every call fails with a stub error telling you to enable CGO.

**Related:** Model bounds for `NAV.ADDTERRAIN` / `NAV.ADDOBSTACLE` come from loaded models (`Model` handles).

---

## `NAV.*` — navigation grid

### `Nav.Make()`

Creates a nav grid handle with a default **64×64** cell layout. Call `NAV.SETGRID` before serious use.

### `Nav.Free(nav)`

Frees the nav object.

### `Nav.SetGrid(nav, gw, gh, cellSize, originX, originZ)`

Resizes the grid: `gw`/`gh` are tile counts (1–4096), `cellSize` is world units per cell (> 0), `(originX, originZ)` is the world origin of cell `(0,0)`. Clears blocked flags and marks the nav **not built** until `NAV.BUILD`.

### `Nav.AddTerrain(nav, modelHandle)` / `Nav.AddObstacle(nav, modelHandle)`

Uses the model’s **axis-aligned bounding box** in world space:

- **AddTerrain** — marks that XZ footprint **walkable** and sets ground height from `Min.Y` across that region.
- **AddObstacle** — marks those cells **blocked**.

### `Nav.Build(nav)`

Marks the nav data ready (sets the internal `built` flag). Call after editing terrain/obstacles and before path queries.

### `Nav.FindPath(nav, sx, sy, sz, tx, ty, tz)`

Runs A* on the grid from start to target world position. Returns a **`Path`** handle (may be invalid if no path). Call `PATH.FREE` when done.

---

## `PATH.*` — path result

### `Path.IsValid(path)` → bool

Whether the search produced a usable path.

### `Path.NodeCount(path)` → int

Number of waypoints (0 if invalid).

### `Path.NodeX(path, index)` / `Path.NodeY` / `Path.NodeZ`

World coordinates of waypoint `index` (0-based). Errors if the index is out of range.

### `Path.Free(path)`

Releases the path handle.

---

## `NAVAGENT.*` — agent on the nav mesh

Create with `NavAgent.Make(nav)` — ties the agent to that **`nav`** handle.

| Command | Role |
|--------|------|
| `NavAgent.SetPos(agent, x, y, z)` | Teleport position. |
| `NavAgent.SetSpeed(agent, speed)` | Max speed (≥ 0). |
| `NavAgent.SetMaxForce(agent, maxForce)` | Caps acceleration from `APPLYFORCE`. |
| `NavAgent.ApplyForce(agent, fx, fy, fz)` | Adds to velocity, then clamps speed. |
| `NavAgent.MoveTo(agent, tx, ty, tz)` | Plans a path with `NAV.FINDPATH`; on success stores waypoints and clears velocity. |
| `NavAgent.Update(agent, dt)` | Advances along waypoints at `speed`, or integrates velocity with damping when no path. |
| `NavAgent.IsAtDestination(agent)` → bool | `TRUE` when there is no active `MoveTo` destination. |
| `NavAgent.X` / `.Y` / `.Z` | Current position. |
| `NavAgent.Free(agent)` | Frees the agent. |

---

## `STEER.*` — steering forces (Vec3 handles)

Steering helpers return **`VEC3`-style handles** (three floats) meant to be combined with `NAVAGENT.APPLYFORCE` or your own logic. Create groups with `Steer.GroupMake`, add agents with `Steer.GroupAdd(group, agent)`.

| Command | Arguments | Result |
|--------|-----------|--------|
| `Steer.Seek` | `(agent, tx, ty, tz)` | Vector toward target. |
| `Steer.Flee` | `(agent, tx, ty, tz)` | Vector away from target. |
| `Steer.Arrive` | `(agent, tx, ty, tz, slowingRadius)` | Seek with speed ramp inside radius. |
| `Steer.Wander` | `(agent, speed, jitterRadius)` | Pseudo-random direction from agent id. |
| `Steer.Flock` | `(selfAgent, group, cohesion, separation, alignment)` | Blended boids-style force. |
| `Steer.AvoidObstacles` | `(agent, radius)` | Repulsion from **blocked** nav cells near the agent. |
| `Steer.FollowPath` | `(agent, path)` | Seeks the nearest waypoint on the path. |

`Steer.GroupClear(group)` empties the group.

---

## `BTREE.*` — behavior tree (user functions)

### `BTree.Make()` / `BTree.Free(bt)`

Allocates a tree whose root is a **sequence** node.

### `BTree.Sequence(bt)` → handle

Returns the same handle (reserved for fluent style; the runtime keeps a single root sequence).

### `BTree.AddCondition(bt, functionName)` / `BTree.AddAction(bt, functionName)`

Appends a child to the root **sequence**. On `BTREE.RUN`, children run in order:

- **Condition** / **action** — invokes the named **user function** with one argument: the **agent handle** passed to `RUN`.
- The function must return a value interpreted as boolean success for conditions; sequence stops on first failure.

### `BTree.Run(bt, agentHandle, dt)`

Walks the tree; `dt` is reserved. User functions are resolved via the VM’s user-function invoker (same mechanism as `SCENE.*` loaders).

---

## Minimal sketch

```basic
; Pseudocode — requires CGO, loaded models, and valid grid setup
nav = Nav.Make()
Nav.SetGrid nav, 64, 64, 1.0, 0.0, 0.0
; ... AddTerrain / AddObstacle with model handles ...
Nav.Build nav
path = Nav.FindPath(nav, x0, y0, z0, x1, y1, z1)
IF Path.IsValid(path) THEN
    PRINT Path.NodeCount(path)
ENDIF
Path.Free path
Nav.Free nav
```
