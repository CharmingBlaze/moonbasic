# Box2D Commands

Low-level Box2D world and body API. For the high-level 2D physics system see [BODY2D.md](BODY2D.md) and [PHYSICS2D.md](PHYSICS2D.md).

`BOX2D.*` gives direct access to the underlying Box2D world for cases where `BODY2D.*` is insufficient.

## Core Workflow

1. `BOX2D.WORLDCREATE(gravX, gravY)` — create a Box2D world.
2. `BOX2D.BODYCREATE(x, y, type)` — add a body.
3. `BOX2D.FIXTUREBOX` or `BOX2D.FIXTURECIRCLE` — attach a shape.
4. `BOX2D.WORLDSTEP(dt, velocityIter, positionIter)` — step.

---

## Commands

### `BOX2D.WORLDCREATE(gravX, gravY)` 

Creates a Box2D world with gravity vector `(gravX, gravY)`. Returns a world handle.

---

### `BOX2D.WORLDSTEP(world, dt, velocityIterations, positionIterations)` 

Steps the Box2D world. Typical values: velocity=8, position=3.

---

### `BOX2D.BODYCREATE(world, x, y, type)` 

Creates a body at `(x, y)`. `type`: `0`=static, `1`=kinematic, `2`=dynamic. Returns a body handle.

---

### `BOX2D.FIXTUREBOX(body, hx, hy, density, friction)` 

Attaches a box fixture (half-extents `hx`, `hy`) to `body`.

---

### `BOX2D.FIXTURECIRCLE(body, radius)` 

Attaches a circle fixture to `body`.

---

## Full Example

```basic
world = BOX2D.WORLDCREATE(0, -10)
floor = BOX2D.BODYCREATE(world, 0, -5, 0)   ; static
BOX2D.FIXTUREBOX(floor, 20, 0.5, 0, 0.3)

ball = BOX2D.BODYCREATE(world, 0, 5, 2)     ; dynamic
BOX2D.FIXTURECIRCLE(ball, 0.5)

WHILE NOT WINDOW.SHOULDCLOSE()
    BOX2D.WORLDSTEP(world, TIME.DELTA(), 8, 3)
    RENDER.FRAME()
WEND
```

---

## See also

- [BODY2D.md](BODY2D.md) — high-level 2D rigid body API
- [PHYSICS2D.md](PHYSICS2D.md) — `PHYSICS2D.START` / `STEP`
- [JOINT2D.md](JOINT2D.md) — 2D constraints
