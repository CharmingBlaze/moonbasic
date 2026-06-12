# 2D collision — overlap tests and platformers

> Detect when 2D shapes touch: bullets vs enemies, player vs platforms, mouse vs buttons — **without** or **with** a physics engine.

**Namespaces:** `COLLISION` · `VEC2` · `PHYSICS2D` · `BODY2D` · **Status:** Shipped

**Commands:** [COMMAND_REGISTRY.md#physics](COMMAND_REGISTRY.md#physics) · [reference/COLLISION.md](../../reference/COLLISION.md) · [reference/PHYSICS2D.md](../../reference/PHYSICS2D.md)

---

## Table of contents

- [At a glance](#at-a-glance)
- [When to use 2D collision](#when-to-use-2d-collision)
- [Choose the right tool](#choose-the-right-tool)
- [Path A — manual rectangles (why beginners start here)](#path-a--manual-rectangles-why-beginners-start-here)
- [Path B — math helpers (`COLLISION.*`)](#path-b--math-helpers-collision)
- [Path C — Box2D (`PHYSICS2D`)](#path-c--box2d-physics2d)
- [Full example — manual platformer](#full-example--manual-platformer)
- [Common mistakes](#common-mistakes)
- [See also](#see-also)

---

## At a glance

| Approach | Best for | You manage |
|----------|----------|------------|
| **Manual rects** | One hero, few platforms, learning | Position, velocity, gravity yourself |
| **`COLLISION.*`** | Many overlap tests, clean math | Positions/sizes as `VEC2` handles |
| **`PHYSICS2D`** | Stacks, joints, realistic bounce | Bodies; engine integrates motion |

**Why 2D collision matters:** Games need to know “did these two things touch?” before you apply damage, play sound, or stop falling.

---

## When to use 2D collision

**Use when:**

- Top-down or side-view 2D gameplay.
- Hit detection (projectile vs target).
- Point-in-rectangle UI tests (optional — `GUI.*` handles mouse for widgets).

**Use 3D instead when:**

- Camera and entities live in X/Y/Z — see [COLLISION-3D.md](COLLISION-3D.md).

---

## Choose the right tool

| I want to… | Use | Why not the other |
|------------|-----|-------------------|
| Single player + 3 platforms | Manual `IF px > …` | Physics adds setup for no gain |
| Many moving crates bouncing | `PHYSICS2D` + `BODY2D` | Manual = nightmare |
| “Is circle A in box B?” once per frame | `COLLISION.CIRCLEBOX2D` | Physics world overhead |
| Rotated rectangles | `PHYSICS2D` or polygon `BODY2D` | AABB manual math gets messy |
| Line-of-sight 2D | `COLLISION.LINESEGINTERSECT2D` | — |

Deep Box2D walkthrough: [PHYSICS-2D-PLATFORMER.md](PHYSICS-2D-PLATFORMER.md)

---

## Path A — manual rectangles (why beginners start here)

**Why:** No handles, no world init — you see exactly what happens. The repo platformer demo uses this pattern ([`examples/platformer/main.mb`](../../../examples/platformer/main.mb)).

**Idea each frame:**

1. Apply input → velocity.
2. Apply gravity → velocity.
3. Position += velocity × `dt`.
4. **Resolve** overlaps with platforms (snap Y, zero vertical velocity, set `onGround`).

```basic
; Ground at y = 400
IF py >= 398 THEN
    py = 400
    pvy = 0
    onGround = 1
ENDIF

; Floating platform: x 200–520, top at y = 300
IF px > 200 AND px < 520 AND py > 300 THEN
    py = 300
    pvy = 0
    onGround = 1
ENDIF
```

**When to graduate:** More than ~5 moving colliders, or you need rotation/restitution.

---

## Path B — math helpers (`COLLISION.*`)

**Why:** Centralized, tested overlap math; pass **`VEC2`** handles instead of eight floats.

| Command | Tests |
|---------|--------|
| `COLLISION.BOXOVERLAP2D(posA, sizeA, posB, sizeB)` | Two AABB rects |
| `COLLISION.CIRCLEOVERLAP2D(c1, r1, c2, r2)` | Two circles |
| `COLLISION.POINTINBOX2D(point, boxPos, boxSize)` | Point inside rect |
| `COLLISION.CIRCLEBOX2D(center, r, boxPos, boxSize)` | Circle vs AABB |
| `COLLISION.LINESEGINTERSECT2D(a1, a2, b1, b2)` | Segment cross |

**Example:**

```basic
pa = VEC2.CREATE(px, py)
sa = VEC2.CREATE(28, 28)
bp = VEC2.CREATE(200, 320)
bs = VEC2.CREATE(320, 24)

IF COLLISION.BOXOVERLAP2D(pa, sa, bp, bs) THEN
  ; standing on platform
ENDIF
```

**Why `VEC2`:** Reuse position handles from movement code; fewer argument mistakes.

Legacy globals (`BOXCOLLIDE`, …) still work — prefer `COLLISION.*` in new code.

---

## Path C — Box2D (`PHYSICS2D`)

**Why:** The engine integrates forces, friction, and resting contact. You read `BODY2D.X/Y` each frame and draw sprites at those coordinates.

**Minimal workflow:**

```basic
PHYSICS2D.START()
PHYSICS2D.SETGRAVITY(0, 500)    ; Y down in 2D

def = BODY2D.CREATE("dynamic")
BODY2D.ADDRECT(def, 28, 28)
player = BODY2D.COMMIT(def, 120, 360)

WHILE NOT APP.SHOULDCLOSE()
    PHYSICS2D.STEP()
    px = BODY2D.X(player)
    py = BODY2D.Y(player)
    ; draw at px, py
WEND
PHYSICS2D.STOP()
```

See [PHYSICS-2D-PLATFORMER.md](PHYSICS-2D-PLATFORMER.md) for static floors, jumps, and queries.

---

## Full example — manual platformer

```basic
WINDOW.OPEN(960, 540, "2D collision manual")
WINDOW.SETFPS(60)

px = 120
py = 360
pvx = 0
pvy = 0
ong = 0

WHILE NOT WINDOW.SHOULDCLOSE()
    dt = TIME.DELTA()
    IF INPUT.KEYDOWN(KEY_A) THEN pvx = pvx - 520 * dt
    IF INPUT.KEYDOWN(KEY_D) THEN pvx = pvx + 520 * dt
    pvx = pvx * 0.88
    IF ong = 1 AND INPUT.KEYDOWN(KEY_SPACE) THEN pvy = -420

    px = px + pvx * dt
    py = py + pvy * dt
    pvy = pvy + 980 * dt
    ong = 0

    IF py > 400 THEN py = 400
    IF py >= 398 THEN pvy = 0
    IF py >= 398 THEN ong = 1

    IF px > 200 AND px < 520 AND py > 300 THEN
        py = 300
        pvy = 0
        ong = 1
    ENDIF

    RENDER.CLEAR(30, 40, 55)
    DRAW.RECTANGLE(0, 440, 960, 100, 50, 120, 70, 255)
    DRAW.RECTANGLE(200, 320, 320, 24, 90, 70, 40, 255)
    DRAW.RECTANGLE(INT(px) - 14, INT(py) - 28, 28, 28, 255, 200, 80, 255)
    RENDER.FRAME()
WEND
WINDOW.CLOSE()
```

`moonrun` required.

---

## Common mistakes

| Mistake | Fix |
|---------|-----|
| Test collision before moving | Order: move → then resolve overlap |
| Forget `onGround` before jump | Only jump when platform resolution set ground |
| Mix screen Y-up math with engine Y-down | 2D draw: Y grows downward on screen |
| Use `PHYSICS2D` for one static box | Manual or `COLLISION.*` is simpler |
| Read body position before `PHYSICS2D.STEP` | Step first, then read `BODY2D.X/Y` |

---

## See also

- [PHYSICS-2D-PLATFORMER.md](PHYSICS-2D-PLATFORMER.md)
- [ENTITY-SYSTEM.md](ENTITY-SYSTEM.md) — 3D objects
- [examples/platformer](../../../examples/platformer/main.mb)
