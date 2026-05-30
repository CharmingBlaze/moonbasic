# Path Commands

Read waypoints from a path handle returned by `NAV.FINDPATH` / `NAV.GETPATH`.

## Core Workflow

1. `path = NAV.FINDPATH(nav, sx, sy, sz, tx, ty, tz)` — get a path handle.
2. `PATH.ISVALID(path)` — check it succeeded.
3. Iterate `0` to `PATH.NODECOUNT(path) - 1`, reading `PATH.NODEX/Y/Z(path, i)`.
4. `PATH.FREE(path)` when done.

---

## Commands

### `PATH.ISVALID(path)` 

Returns `TRUE` if the path is valid (a route was found).

---

### `PATH.NODECOUNT(path)` 

Returns the number of waypoints in the path.

---

### `PATH.NODEX(path, index)` / `PATH.NODEY(path, index)` / `PATH.NODEZ(path, index)` 

Returns the world X, Y, or Z coordinate of waypoint at `index` (0-based).

---

### `PATH.FREE(path)` 

Frees the path handle.

---

## Full Example

Following a found path waypoint by waypoint.

```basic
nav  = NAV.CREATE()
NAV.SETGRID(nav, 20, 20, 1.0, -10, 0, -10)
NAV.BUILD(nav)

path = NAV.FINDPATH(nav, -8, 0, -8, 8, 0, 8)

IF PATH.ISVALID(path) THEN
    n = PATH.NODECOUNT(path)
    PRINT "Path waypoints: " + STR(n)
    FOR i = 0 TO n - 1
        PRINT STR(i) + ": " + STR(PATH.NODEX(path, i)) + ", " + STR(PATH.NODEZ(path, i))
    NEXT i
ELSE
    PRINT "No path found"
END IF

PATH.FREE(path)
NAV.FREE(nav)
```

---

## See also

- [NAV.md](NAV.md) — navmesh and grid construction
- [NAVAGENT.md](NAVAGENT.md) — agents that follow paths automatically
