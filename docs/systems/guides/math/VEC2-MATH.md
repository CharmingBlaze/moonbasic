# 2D vector math — VEC2 handles and screen/world directions

> Create, normalize, rotate, and interpolate 2D directions — the right layer when scalar `DIST2D` is not enough.

**Namespaces:** `VEC2` · **Status:** Shipped

**Commands:** [COMMAND_REGISTRY.md#data](../../COMMAND_REGISTRY.md#data) · **Scalars:** [MATH-2D-GAMEPLAY.md](MATH-2D-GAMEPLAY.md) · **Deep reference:** [reference/VEC2.md](../../../reference/VEC2.md)

---

## Table of contents

- [At a glance](#at-a-glance)
- [When to use VEC2](#when-to-use-vec2)
- [Choose the right tool](#choose-the-right-tool)
- [Core workflow](#core-workflow)
- [Key commands](#key-commands)
- [Patterns](#patterns)
- [Full example — arena movement](#full-example--arena-movement)
- [Memory notes](#memory-notes)
- [Common mistakes](#common-mistakes)
- [See also](#see-also)

---

## At a glance

| Idea | Detail |
|------|--------|
| **You get** | Heap handles for x,y; normalize, rotate, lerp, pushout, move_toward |
| **You need first** | 2D positions as floats or sprites |
| **Typical games** | Top-down, twin-stick, 2D physics overlap resolution |
| **Angles** | **Radians** for `ROTATE` and `ANGLE` |

**Why handles:** Reusable direction objects; scalar overloads (`LENGTH(x,y)`) avoid allocation for one-off checks.

---

## When to use VEC2

**Use when:**

- Diagonal keyboard move must stay constant speed (normalize then scale).
- Rotate thrust vector each frame.
- Push player out of circle collider (`PUSHOUT`).
- `COLLISION.*2D` tests pass `VEC2` center/size handles.

**Skip when:**

- One distance compare per frame — `MATH.DISTSQ2D` is enough.
- Pure 3D gameplay — [VEC3-MATH.md](VEC3-MATH.md).

---

## Choose the right tool

| I want to… | Use | Not |
|------------|-----|-----|
| One distance check | `MATH.DIST2D` | `VEC2.CREATE` each time |
| Unit move direction | `VEC2.NORMALIZE` | Divide x,y by hand |
| Slide along circle edge | `VEC2.PUSHOUT` | Manual trig |
| Smooth chase point | `VEC2.MOVE_TOWARD` | Lerp scalars separately |
| 2D box overlap API | `VEC2.CREATE` + `COLLISION.*2D` | Four loose floats |

---

## Core workflow

1. **Create** — `v = VEC2.CREATE(x, y)` when you need a handle.
2. **Operate** — `ADD`, `SUB`, `MUL`, `NORMALIZE`, `ROTATE`.
3. **Read** — `VEC2.X(v)`, `VEC2.Y(v)`.
4. **Free** — `VEC2.FREE(v)` when done (or reuse one frame).

**Scalar overloads:** `VEC2.LENGTH(x, y)`, `VEC2.DIST(x1,y1,x2,y2)`, `VEC2.NORMALIZE(x, y)` return results without long-lived handles when possible.

---

## Key commands

| Command | Why |
|---------|-----|
| `VEC2.CREATE(x, y)` | New 2D vector handle |
| `VEC2.ADD` / `SUB` / `MUL` | Vector arithmetic |
| `VEC2.LENGTH(v)` | Magnitude |
| `VEC2.NORMALIZE(v)` | Unit direction |
| `VEC2.DIST` / `DISTSQ` | Gap between points |
| `VEC2.DOT` | Not in registry — use angle between via `ANGLE` |
| `VEC2.ANGLE(a, b)` | Signed angle between vectors (rad) |
| `VEC2.ROTATE(v, radians)` | Spin vector |
| `VEC2.LERP(a, b, t)` | Blend positions |
| `VEC2.MOVE_TOWARD(x,y,tx,ty,step)` | Cap step toward target |
| `VEC2.PUSHOUT(x,y,cx,cy,r)` | Resolve circle penetration |
| `VEC2.TRANSFORMMAT4(v, mat)` | Apply camera matrix |
| `VEC2.FREE(v)` | Release handle |

---

## Patterns

### Constant-speed diagonal move

```basic
dx = 0
dy = 0
IF INPUT.KEYDOWN(KEY_RIGHT) THEN dx = 1
IF INPUT.KEYDOWN(KEY_LEFT)  THEN dx = -1
IF INPUT.KEYDOWN(KEY_DOWN)  THEN dy = 1
IF INPUT.KEYDOWN(KEY_UP)    THEN dy = -1

IF dx <> 0 OR dy <> 0 THEN
    vel = VEC2.CREATE(dx, dy)
    nrm = VEC2.NORMALIZE(vel)
    px = px + VEC2.X(nrm) * speed * APP.DELTA()
    py = py + VEC2.Y(nrm) * speed * APP.DELTA()
    VEC2.FREE(nrm)
    VEC2.FREE(vel)
ENDIF
```

### Compare distance without sqrt

```basic
dSq = VEC2.DISTSQ(px, py, ex, ey)
IF dSq < aggro * aggro THEN ...
```

---

## Full example — arena movement

See [reference/VEC2.md](../../../reference/VEC2.md) for the arena `PUSHOUT` demo. Minimal loop:

```basic
APP.OPEN(800, 600, "VEC2")
APP.SETFPS(60)
px = 400
py = 300

WHILE NOT APP.SHOULDCLOSE()
    dx = 0
    dy = 0
    IF INPUT.KEYDOWN(KEY_D) THEN dx = 1
    IF INPUT.KEYDOWN(KEY_A) THEN dx = -1
    IF INPUT.KEYDOWN(KEY_S) THEN dy = 1
    IF INPUT.KEYDOWN(KEY_W) THEN dy = -1

    IF dx <> 0 OR dy <> 0 THEN
        n = VEC2.NORMALIZE(dx, dy)
        px = px + VEC2.X(n) * 200 * APP.DELTA()
        py = py + VEC2.Y(n) * 200 * APP.DELTA()
        VEC2.FREE(n)
    ENDIF

    out = VEC2.PUSHOUT(px, py, 400, 300, 180)
    px = VEC2.X(out)
    py = VEC2.Y(out)
    VEC2.FREE(out)

    DRAW.CIRCLEOUTLINE(400, 300, 180, 60, 60, 100, 255)
    DRAW.CIRCLE(INT(px), INT(py), 10, 100, 180, 255, 255)
    RENDER.FRAME()
WEND
APP.CLOSE()
```

---

## Memory notes

- `VEC2.FREE` handles you created with `CREATE` / `NORMALIZE` returning new handles.
- Prefer scalar overloads in hot loops if profiling shows heap pressure.

---

## Common mistakes

| Mistake | Fix |
|---------|-----|
| Normalize (0,0) | Check length > epsilon first |
| Forget `FREE` in loops | Free each frame or use scalar overloads |
| `ROTATE` in degrees | Convert with `MATH.DEG2RAD` |
| Y axis flip | Screen Y down vs world Y up — flip input |
| `MAKE` vs `CREATE` | `MAKE` deprecated — use `CREATE` |

---

## See also

- [MATH-2D-GAMEPLAY.md](MATH-2D-GAMEPLAY.md) — scalar distances and angles
- [COLLISION-2D.md](../COLLISION-2D.md) — `VEC2` in overlap tests
- [SPRITES-TILEMAPS-2D.md](../SPRITES-TILEMAPS-2D.md) — sprite positions
