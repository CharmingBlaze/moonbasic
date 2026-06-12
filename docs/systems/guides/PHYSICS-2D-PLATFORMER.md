# 2D physics platformer — Box2D workflow

> Use **`PHYSICS2D`** and **`BODY2D`** when platforms, crates, and the player should **move and collide automatically** — beyond manual rectangle math.

**Namespaces:** `PHYSICS2D` · `BODY2D` · `JOINT2D` · `PLAYER2D` · **Status:** Shipped

**Commands:** [COMMAND_REGISTRY.md#physics](../COMMAND_REGISTRY.md#physics) · [reference/PHYSICS2D.md](../../reference/PHYSICS2D.md)

**Compare:** Manual platforms → [COLLISION-2D.md](COLLISION-2D.md)

---

## Table of contents

- [At a glance](#at-a-glance)
- [When to use PHYSICS2D](#when-to-use-physics2d)
- [Core workflow](#core-workflow)
- [Building bodies](#building-bodies)
- [Moving the player](#moving-the-player)
- [Collision queries](#collision-queries)
- [Tuning](#tuning)
- [Full example sketch](#full-example-sketch)
- [Common mistakes](#common-mistakes)
- [See also](#see-also)

---

## At a glance

| Step | Command | Why |
|------|---------|-----|
| Start world | `PHYSICS2D.START()` | Creates Box2D world |
| Gravity | `PHYSICS2D.SETGRAVITY(0, 500)` | Y-down screen coords |
| Build body | `BODY2D.CREATE` → shapes → `COMMIT` | Definition vs runtime body |
| Simulate | `PHYSICS2D.STEP()` | Once per frame |
| Read pose | `BODY2D.X/Y/ROT` | Sync sprite position |

**Default gravity** if omitted: `(0, 500)` — positive Y is **down** on screen.

---

## When to use PHYSICS2D

**Use when:**

- Many dynamic objects, crates, ragdolls (2D).
- Joints (doors, bridges) — `JOINT2D.*`.
- You want friction/restitution from material properties.

**Use manual collision when:**

- One character and a handful of platforms — [COLLISION-2D.md](COLLISION-2D.md) is simpler.

---

## Core workflow

1. `PHYSICS2D.START([gx, gy])`
2. Create **static** ground: `BODY2D.CREATE("static")` → `ADDRECT` → `COMMIT`
3. Create **dynamic** player: `BODY2D.CREATE("dynamic")` → shape → `COMMIT`
4. Loop:
   - Read input → `BODY2D.APPLYFORCE` / `SETVEL` on player
   - `PHYSICS2D.STEP()`
   - Draw sprites at `BODY2D.X/Y`
5. `PHYSICS2D.STOP()`

**Why `COMMIT`:** Shapes attach to a **definition**; `COMMIT` inserts the body into the world at `(x,y)`.

---

## Building bodies

```basic
PHYSICS2D.START()

; Static floor
floorDef = BODY2D.CREATE("static")
BODY2D.ADDRECT(floorDef, 960, 40)
floor = BODY2D.COMMIT(floorDef, 480, 520)

; Dynamic player box
playerDef = BODY2D.CREATE("dynamic")
BODY2D.ADDRECT(playerDef, 28, 28)
player = BODY2D.COMMIT(playerDef, 120, 400)
BODY2D.SETFRICTION(player, 0.2)
BODY2D.SETRESTITUTION(player, 0)
```

| Type | Why |
|------|-----|
| `static` | Never moves — terrain |
| `dynamic` | Simulated mass |
| `kinematic` | Moved by script, pushes dynamics |

Shapes: `ADDRECT`, `ADDCIRCLE`, `ADDPOLYGON` — see manifest for arity.

---

## Moving the player

**Why forces vs teleport:** `SETPOS` every frame fights the solver — use impulses for jumps, forces for run.

```basic
dt = TIME.DELTA()
IF INPUT.KEYDOWN(KEY_A) THEN BODY2D.APPLYFORCE(player, -800, 0)
IF INPUT.KEYDOWN(KEY_D) THEN BODY2D.APPLYFORCE(player, 800, 0)
IF INPUT.KEYHIT(KEY_SPACE) THEN BODY2D.APPLYIMPULSE(player, 0, -300)

PHYSICS2D.STEP()
px = BODY2D.X(player)
py = BODY2D.Y(player)
```

**Handle shortcuts:** `player.applyForce(fx,fy)` · `player.applyImpulse(ix,iy)`

Optional helpers: `PLAYER2D.*` — see [reference/PHYSICS2D.md](../../reference/PHYSICS2D.md).

---

## Collision queries

After `PHYSICS2D.STEP`, per-body:

- Contact lists / overlap queries — see **Per-Body Collision Queries** in [reference/PHYSICS2D.md](../../reference/PHYSICS2D.md)
- `RAY2D.*` for line casts in the 2D world

**Why after STEP:** Contacts exist after the solver runs.

---

## Tuning

| Command | Effect |
|---------|--------|
| `PHYSICS2D.SETGRAVITY(x,y)` | Fall speed |
| `BODY2D.SETFRICTION` | Slide vs stick |
| `BODY2D.SETRESTITUTION` | Bounce 0..1 |
| `BODY2D.SETMASS` | Heavy vs light |
| `PHYSICS2D.SETITERATIONS` | Accuracy vs CPU (if exposed) |

Match step to frame rate: one `STEP()` per frame with `WINDOW.SETFPS(60)` is the common case.

---

## Full example sketch

```basic
WINDOW.OPEN(960, 540, "Box2D platformer")
WINDOW.SETFPS(60)
PHYSICS2D.START(0, 600)

fd = BODY2D.CREATE("static")
BODY2D.ADDRECT(fd, 900, 50)
BODY2D.COMMIT(fd, 480, 500)

pd = BODY2D.CREATE("dynamic")
BODY2D.ADDRECT(pd, 30, 30)
hero = BODY2D.COMMIT(pd, 200, 300)

WHILE NOT WINDOW.SHOULDCLOSE()
    dt = TIME.DELTA()
    IF INPUT.KEYDOWN(KEY_D) THEN BODY2D.APPLYFORCE(hero, 400, 0)
    IF INPUT.KEYHIT(KEY_SPACE) THEN BODY2D.APPLYIMPULSE(hero, 0, -250)
    PHYSICS2D.STEP()
    hx = BODY2D.X(hero)
    hy = BODY2D.Y(hero)
    RENDER.CLEAR(40, 44, 52)
    DRAW.RECTANGLE(0, 480, 960, 60, 60, 100, 70, 255)
    DRAW.RECTANGLE(INT(hx)-15, INT(hy)-15, 30, 30, 255, 180, 80, 255)
    RENDER.FRAME()
WEND

BODY2D.FREE(hero)
PHYSICS2D.STOP()
WINDOW.CLOSE()
```

---

## Common mistakes

| Mistake | Fix |
|---------|-----|
| Draw before `STEP` | Step then read position |
| `COMMIT` without shapes | Add rect/circle first |
| Huge impulse every frame | Use force for run, impulse once for jump |
| Forget `BODY2D.FREE` | Free bodies when removing |
| Wrong gravity sign | Screen Y down → positive gravity Y |

---

## See also

- [COLLISION-2D.md](COLLISION-2D.md)
- [examples/platformer](../../../examples/platformer/main.mb) — manual version
- `moonbasic new --template platformer`
