# Player2D Commands

Simple 2D player handle for screen-space movement, clamping, and position queries. Designed for top-down and side-scroller prototypes.

## Core Workflow

1. `PLAYER2D.CREATE()` — allocate a player handle.
2. `PLAYER2D.SETPOS(handle, x, z)` — place the player.
3. Each frame: `PLAYER2D.MOVE(handle, dx, dz, speed, dt, friction)` — move with optional friction.
4. `PLAYER2D.CLAMP(handle, minX, minZ, maxX, maxZ)` or `PLAYER2D.KEEPINBOUNDS(handle)` — constrain.
5. `PLAYER2D.GETX(handle)` / `PLAYER2D.GETZ(handle)` — read position.
6. `PLAYER2D.FREE(handle)` when done.

---

## Commands

### `PLAYER2D.CREATE()` 

Creates a 2D player handle at origin. Returns a **player handle**.

---

### `PLAYER2D.SETPOS(handle, x, z)` 

Teleports the player to `(x, z)`.

---

### `PLAYER2D.MOVE(handle, dx, dz, speed, dt, friction)` 

Moves the player by `(dx, dz) × speed × dt`, with optional friction deceleration each frame.

---

### `PLAYER2D.CLAMP(handle, minX, minZ, maxX, maxZ)` 

Restricts the player position to the given axis-aligned rectangle.

---

### `PLAYER2D.KEEPINBOUNDS(handle)` 

Clamps the player to the window bounds.

---

### `PLAYER2D.GETX(handle)` / `PLAYER2D.GETZ(handle)` 

Returns the X or Z position as a float.

---

### `PLAYER2D.GETPOS(handle)` 

Returns the XZ position as a Vec2 handle.

---

### `PLAYER2D.FREE(handle)` 

Frees the player handle.

---

## Full Example

Top-down movement with screen boundary clamping.

```basic
WINDOW.OPEN(800, 600, "Player2D Demo")
WINDOW.SETFPS(60)

p = PLAYER2D.CREATE()
PLAYER2D.SETPOS(p, 400, 300)

WHILE NOT WINDOW.SHOULDCLOSE()
    dt = TIME.DELTA()
    dx = 0.0 : dz = 0.0
    IF INPUT.KEYDOWN(KEY_RIGHT) THEN dx =  1
    IF INPUT.KEYDOWN(KEY_LEFT)  THEN dx = -1
    IF INPUT.KEYDOWN(KEY_DOWN)  THEN dz =  1
    IF INPUT.KEYDOWN(KEY_UP)    THEN dz = -1

    PLAYER2D.MOVE(p, dx, dz, 200, dt, 0.0)
    PLAYER2D.KEEPINBOUNDS(p)

    RENDER.CLEAR(20, 20, 40)
    DRAW.CIRCLE(INT(PLAYER2D.GETX(p)), INT(PLAYER2D.GETZ(p)), 16, 80, 200, 255, 255)
    RENDER.FRAME()
WEND

PLAYER2D.FREE(p)
WINDOW.CLOSE()
```

---

## Extended Command Reference

| Command | Description |
|--------|-------------|
| `PLAYER2D.MAKE(...)` | Deprecated alias of `PLAYER2D.CREATE`. |
| `PLAYER2D.SETPOSITION(p, x, y)` | Alias of `PLAYER2D.SETPOS`. |

---

## See also

- [PLAYER.md](PLAYER.md) — 3D player controller
- [CHARCONTROLLER.md](CHARCONTROLLER.md) — 3D capsule controller
- [MOUSE.md](MOUSE.md) — mouse position for top-down aiming
