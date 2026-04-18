# Less math (gameplay helpers)

These built-ins replace repeated **distance**, **normalize**, **spawn ring**, **camera-relative WASD**, and **terrain height snap** patterns so loops stay readable without hand-written `SQRT`, `COS`/`SIN`, or yaw decomposition.

**CGO** is required for the full runtime (`INPUT.MOVEDIR`, `INPUT.MOUSEDELTA`, `TERRAIN.SNAPY`, `WORLD.SETCENTERENTITY`, `ENTITY.GETXZ`, etc.); stub builds may no-op or error.

See the refactored sample: [`examples/terrain_chase/main.mb`](../../examples/terrain_chase/main.mb).

Project stance (**helpers + full `MATH.*`**): [Easy language](../EASY_LANGUAGE.md).

## Core Workflow

- Camera-relative WASD: `INPUT.MOVEDIR(camYaw)` → returns `[dx, dz]` ready for entity movement.
- Terrain snap: `TERRAIN.SNAPY(terrain, x, z, offset)` → sets entity Y each frame.
- World streaming: `WORLD.SETCENTERENTITY(entity)` → keeps chunks loaded around the player.
- Facing a target: `ENTITY.FACETARGET(e, targetX, targetZ)` → smooth yaw rotation.

---

## `MATH.CIRCLEPOINT(cx, cz, radius, i, count)` → `(x, z)`

Places **`count`** points on a full circle around **`(cx, cz)`** with radius **`radius`**. Index **`i`** is **1-based** (`1 .. count`); angle is `(i-1) * 2π / count`.

**Alias:** `CIRCLEPOINT` (same arguments).

```basic
ex, ez = MATH.CIRCLEPOINT(TCX, TCZ, 42.0, i, N_ENEMY)
```

---

## `VEC2.DIST(x1, y1, x2, y2)` → `float`

Euclidean distance between **`(x1,y1)`** and **`(x2,y2)`**. Overload: two **`VEC2`** handles (same name).

---

## `VEC2.DISTSQ(x1, y1, x2, y2)` → `float`

Squared distance — use for **comparisons** (e.g. `IF VEC2.DISTSQ(...) < r*r`) to avoid `SQRT`.

---

## `VEC3.DIST` / `VEC3.DISTSQ`

**`VEC3.DIST(x1, y1, z1, x2, y2, z2)`** → **float** — 3D distance without building vec3 handles. Overload: **two vec3 handles** (same behavior as **`VEC3.Distance`**).

**`VEC3.DISTSQ(x1, y1, z1, x2, y2, z2)`** → **float** — squared distance for **radius / proximity** tests without `SQRT`.

---

## `VEC2.MOVE_TOWARD(x, y, tx, ty, maxDist)` → `(x, y)`

Steps from **`(x,y)`** toward **`(tx,ty)`** by at most **`maxDist`** (zero length is safe).

---

## `VEC2.NORMALIZE(x, y)` → `(x, y)`

Unit vector in the **`(x,y)`** direction; returns **`(0,0)`** when length is zero (no manual guard).

---

## `VEC2.PUSHOUT(x, z, cx, cz, minRadius)` → `(x, z)`

If **`(x,z)`** is inside radius **`minRadius`** of **`(cx,cz)`**, pushes it outward to the circle boundary. Unchanged if already outside. If the point is on the center (degenerate), pushes along **+X**.

```basic
bx, bz = VEC2.PUSHOUT(bx, bz, TCX, TCZ, BEACON_SAFE_RADIUS)
```

---

## `INPUT.MOVEDIR(yaw, speed)` → `(stepX, stepZ)`

Reads **W/A/S/D**, normalizes the **forward/strafe** mix, rotates by **`yaw`** (radians) into world **XZ**, multiplies by **`speed`**. Matches the usual third-person convention: forward **`(sin(yaw), cos(yaw))`**, right **`(cos(yaw), -sin(yaw))`**.

```basic
stepX, stepZ = INPUT.MOVEDIR(camYaw, MOVE_SPD * dt)
```

---

## `INPUT.MOUSEDELTA()` → `(dx, dy)`

Single call replacing **`INPUT.MOUSEXSPEED()`** + **`INPUT.MOUSEYSPEED()`** (frame mouse delta from the underlying input layer).

---

## `MATH.LERPANGLE(a, b, t)` → `float`

Shortest-path interpolation between angles **`a`** and **`b`** (radians), factor **`t`** in **`[0,1]`** (implementation uses `atan2`/`sin`/`cos` so **350°** eases toward **10°** the short way).

---

## `MATH.REMAP(value, inMin, inMax, outMin, outMax)` → `float`

Linearly maps **`value`** from the range **`[inMin, inMax]`** to **`[outMin, outMax]`**. If **`inMin` = `inMax`**, returns **`outMin`** (safe default).

**Alias:** `REMAP`.

---

## `MATH.INVERSE_LERP(a, b, x)` → `float`

Returns **`(x - a) / (b - a)`** — how far **`x`** is between **`a`** and **`b`** (not clamped). If **`a` = `b`**, returns **`0`**. Often paired with **`MATH.LERP`** to convert UI or sensor ranges.

**Alias:** `INVERSE_LERP`.

---

## `MATH.SATURATE(x)` → `float`

Clamps **`x`** to **`[0, 1]`** (handy for turning arbitrary values into blend factors).

**Alias:** `SATURATE`.

---

## `MATH.APPROACH(current, target, maxStep)` → `float`

Moves **`current`** toward **`target`** by at most **`maxStep`**, clamping so it does not overshoot. Handy for zoom distance or any scalar that should ease to a target without picking a lerp factor.

---

## `Terrain.SnapY(terrain, entity, yOffset)`

Reads the entity’s **world XZ**, samples the heightfield, sets **Y** to **`height + yOffset`**, and writes the global position. Replaces the two-line **`GetHeight` + reposition** pattern for grounded props.

```basic
Terrain.SnapY(terrain, player, 0.55)
```

## `Terrain.Place(terrain, entity, x, z, yOffset)`

Sets **X** and **Z**, then snaps **Y** from the heightfield — same end result as **`PositionEntity`** (with ground align) plus **`Terrain.SnapY`**, in one call. See [TERRAIN.md](TERRAIN.md).

---

## `Entity.GetXZ(entity)` → `(x, z)`

World **X** and **Z** without pulling **Y** when you only need ground-plane coordinates.

---

## `World.SetCenterEntity(entity)`

Sets streaming center from the entity’s world **XZ** (same role as **`World.SetCenter(x,z)`** for terrain streaming). See [WORLD.md](WORLD.md).

---

## More table helpers (XZ / 2D / easing)

See **[GAME_MATH_HELPERS.md](GAME_MATH_HELPERS.md)** for **`HDIST` / `HDISTSQ`** (horizontal 3D distance), **`DIST2D` / `DISTSQ2D`**, **`YAWFROMXZ`**, **`ANGLEDIFFRAD`**, and **`SMOOTHERSTEP`**.

---

## Full Example

Camera-relative movement with LERP approach.

```basic
WINDOW.OPEN(960, 540, "LessMath Demo")
WINDOW.SETFPS(60)

cam = CAMERA.CREATE()
CAMERA.SETPOS(cam, 0, 8, -12)
CAMERA.SETTARGET(cam, 0, 0, 0)

px = 0.0  pz = 0.0
camYaw = 0.0

WHILE NOT WINDOW.SHOULDCLOSE()
    dt = TIME.DELTA()
    f  = INPUT.AXIS(KEY_S, KEY_W)
    s  = INPUT.AXIS(KEY_A, KEY_D)

    camYaw = camYaw + INPUT.ORBIT(KEY_Q, KEY_E, 90, dt)

    dx = SIN(camYaw) * f + COS(camYaw) * s
    dz = COS(camYaw) * f - SIN(camYaw) * s
    px = LERP(px, px + dx * 4 * dt, 0.2)
    pz = LERP(pz, pz + dz * 4 * dt, 0.2)

    CAMERA.SETORBIT(cam, px, 0, pz, camYaw, 0.5, 10)

    RENDER.CLEAR(30, 40, 60)
    RENDER.BEGIN3D(cam)
        DRAW.SPHERE(px, 0.5, pz, 0.4, 80, 200, 255, 255)
        DRAW.GRID(20, 1.0)
    RENDER.END3D()
    RENDER.FRAME()
WEND

CAMERA.FREE(cam)
WINDOW.CLOSE()
```

---

## See also

- [MATH.md](MATH.md) — core math
- [GAME_MATH_HELPERS.md](GAME_MATH_HELPERS.md) — 2D/3D gameplay math shortcuts
- [VEC_QUAT.md](VEC_QUAT.md) — `VEC2.*` overview
- [INPUT.md](INPUT.md) — keyboard/mouse
- [TERRAIN.md](TERRAIN.md) — heightfield
- [ENTITY.md](ENTITY.md) — entity position helpers
- [GAMEHELPERS.md](GAMEHELPERS.md) — other gameplay shortcuts
