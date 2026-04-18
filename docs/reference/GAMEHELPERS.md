# Game Helper Commands

Small built-ins for box landing, camera orbit, 2D movers, and camera-relative movement.

Page shape follows [DOC_STYLE_GUIDE.md](../DOC_STYLE_GUIDE.md) (**WAVE pattern**).

## Core Workflow

Use `BOXTOPLAND` / `LANDBOXES` for platform snapping, `ORBITYAWDELTA` / `ORBITPITCHDELTA` / `ORBITDISTDELTA` for third-person orbit input, and `PLAYER2D.*` / `MOVEENTITY2D` for XZ ground-plane movement. Combine with `CAMERA.SETORBIT` for a complete third-person loop.

---

## `BOXTOPLAND`

`BOXTOPLAND(px, py, pz, pvy, pr, bx, by, bz, bw, bh, bd)` → **float**

Returns **`0.0`** when there is **no** top landing this frame. Otherwise returns the **sphere center Y** to snap to (`box top + pr`).

- Only meaningful when **`pvy <= 0`** (falling or resting). If moving upward, returns **`0.0`**.
- Horizontal test: sphere center must be within the box footprint, expanded by **`pr`** on X/Z.
- Vertical test: **feet** (`py - pr`) must sit in a small band below/around the **top** of the box (`by + bh/2`).

Typical use:

```basic
landY = BOXTOPLAND(px, py, pz, pvy, pr, bx, by, bz, bw, bh, bd)
IF landY > 0.0 THEN
    py = landY : pvy = 0.0 : on_ground = TRUE
ENDIF
```

---

## `LANDBOXES` / `LANDBOX`

`LANDBOXES(px, py, pz, pvy, pr, plx, ply, plz, plw, plh, pld, count)` → **float**

**`LANDBOX`** is an **alias** — same arguments and return value.

Runs the same test as **`BOXTOPLAND`** for **`count`** boxes given as **six parallel float arrays** (centre `x,y,z` and size `w,h,d`). Returns the **largest** positive snap Y among all boxes, or **`0.0`** if none apply. Use this instead of a **`FOR`** loop when platforms are stored as parallel **`DIM`** arrays.

Implementation note: it is equivalent to **`BOXTOPLAND`** per index — not a full physics engine. **`TYPE`** platform rows still use a loop or manual **`BOXTOPLAND`** unless you keep parallel arrays for collision.

---

## `PLAYER.MOVERELATIVE`

`PLAYER.MOVERELATIVE(camYaw, forward, strafe, speed, dt)` → **handle** (2-float array **`[deltaX, deltaZ]`**)

Same math as **`MOVESTEPX`** and **`MOVESTEPZ`** combined. **Free** the returned array with **`ERASE`** when you are done (each frame if you allocate every frame). For hot loops, **`MOVESTEPX`/`MOVESTEPZ`** avoid the extra heap array.

---

## Simple physics without a physics engine

Gravity and integration are only a few lines. Keep **`dt`** from **`TIME.DELTA()`** or **`DT()`** (both are **clamped** by default so tab-switch spikes do not explode simulation).

```basic
CONST GRAVITY = -26.0

; Each frame:
vel_y = vel_y + GRAVITY * dt
pos_y = pos_y + vel_y * dt

; Ground check (flat floor at y = radius):
IF pos_y < radius THEN
    pos_y = radius
    vel_y = 0.0
    on_ground = TRUE
ENDIF
```

For **one-shot** actions (jump, shoot), use **`INPUT.KEYPRESSED`** or the flat **`KEYPRESSED`** helper, not **`INPUT.KEYDOWN`**, which is **TRUE** every frame the key is held.

See also: [INPUT.md](INPUT.md) (keyboard table), [CAMERA.md](CAMERA.md) (**`CAMERA.ORBITAROUND`** for third-person orbit).

---

## Third-person orbit input (`ORBITYAWDELTA` / `ORBITPITCHDELTA` / `ORBITDISTDELTA`)

These **`GAME.*`** helpers (short names without the `GAME.` prefix also work) wrap **right-mouse drag** and **mouse wheel** together with the same **Q/E yaw** math as **`INPUT.ORBIT`**. They return **plain floats** each frame — **no heap handles**, nothing to **`ERASE`**. Use them to update your **`camYaw`**, **`camPitch`**, and **`camDist`**, then call **`CAMERA.SETORBIT`** (or **`CAMERA.ORBITAROUND`**) yourself.

| Command | Returns | Meaning |
|--------|---------|--------|
| **`ORBITYAWDELTA(dt, mouseSens, negKey, posKey, degPerSec)`** | radians | **Keyboard:** same as **`INPUT.ORBIT(negKey, posKey, degPerSec, dt)`** (degrees/sec → radians). **Mouse:** if **right button** is down, adds **`MDX * mouseSens`** (typically `mouseSens` ≈ `0.004`–`0.006`). |
| **`ORBITPITCHDELTA(mouseSens)`** | radians | If **right button** is down: **`-MDY * mouseSens`**. Otherwise **`0`**. |
| **`ORBITDISTDELTA(wheelScale)`** | world units | **`-MWHEEL * wheelScale`** — add to your orbit distance (scroll **up** moves the eye **closer** when **`wheelScale`** is positive). |

Clamp **`pitch`** and **`dist`** in your script after adding deltas (the helpers do not clamp).

Typical frame (see **`examples/mario64/main_orbit_simple.mb`**):

```basic
camYaw = camYaw + ORBITYAWDELTA(dt, 0.0048, KEY_Q, KEY_E, 72.0)
camPitch = camPitch + ORBITPITCHDELTA(0.0048)
camDist = camDist + ORBITDISTDELTA(0.85)
; … clamp pitch & dist, then CAMERA.SETORBIT(cam, tx, ty, tz, camYaw, camPitch, camDist)
```

That example is structured for reading **top to bottom**: one **`CONST`** block (world bounds, orbit tuning, colours), parallel **`DIM`** rows for **`LANDBOXES`**, a single loop section for input → physics → **`CAMERA.SETORBIT`** → draw, then **`ERASE ALL`** (see [MEMORY.md](../MEMORY.md)).

**Memory:** no allocations — see [MEMORY.md](../MEMORY.md) (game orbit helpers).

---

## Blitz-style English helpers (2D XZ mover + camera yaw)

These read like classic Blitz commands: a **`PLAYER2D`** handle stores **X/Z** on the ground plane. **`MOVEENTITY2D`**, **`MOVEPLAYER`**, **`CLAMPENTITY2D`**, and **`KEEPPLAYERINBOUNDS`** are **aliases** of the same **`PLAYER2D.*`** implementations (pick whichever name reads best in your script).

| Command | Role |
|--------|------|
| **`p = PLAYER2D.CREATE()`** | Create a mover; **`PLAYER2D.FREE p`** or scene **`ERASE ALL`** when done. |
| **`PLAYER2D.SETPOS p, x, z`** | Set world X/Z (e.g. spawn). |
| **`MOVEENTITY2D p, camYaw, f, s, speed, dt`** | Camera-relative move on **XZ** (same math as **`MOVESTEPX`/`MOVESTEPZ`** applied in place). Aliases: **`PLAYER2D.MOVE`**, **`MOVEPLAYER`**. |
| **`CLAMPENTITY2D p, minX, maxX, minZ, maxZ`** | Store bounds and clamp **current** position into the axis-aligned box. Alias: **`PLAYER2D.CLAMP`**. |
| **`KEEPPLAYERINBOUNDS p`** | Clamp again using the **last** bounds from **`CLAMPENTITY2D`** (call after **`MOVEENTITY2D`** each frame). No-op if bounds were never set. Alias: **`PLAYER2D.KEEPINBOUNDS`**. |
| **`PLAYER2D.GETX p`**, **`PLAYER2D.GETZ p`** | Read position for **`BOXTOPLAND`**, rendering, etc. |

Camera **yaw** is still a script variable (e.g. **`camYaw`**). The camera handle is only validated so you do not pass the wrong object:

| Command | Returns | Role |
|--------|---------|------|
| **`CAMERA.TURNLEFT cam, amount`** | **float** (radians) | **`-abs(amount)`** — add to **`camYaw`** to turn left. |
| **`CAMERA.TURNRIGHT cam, amount`** | **float** (radians) | **`+abs(amount)`** — add to **`camYaw`** to turn right. |
| **`CAMERA.ORBITCAMERA cam, mouseSens, keyDegPerSec, dt`** | **float** (radians) | Same as **`FLOAT(INPUT.MOUSEDELTAX()) * mouseSens + INPUT.ORBIT(KEY_Q, KEY_E, keyDegPerSec, dt)`** — add the result to **`camYaw`** each frame. |

Example:

```basic
p = PLAYER2D.CREATE()
PLAYER2D.SETPOS(p, 0.0, 0.0)
CLAMPENTITY2D(p, -17.0, 17.0, -17.0, 22.0)

camYaw = camYaw + CAMERA.ORBITCAMERA(cam, MOUSE_ORBIT_SENS, 77.0, dt)
f = INPUT.AXIS(KEY_S, KEY_W)
s = INPUT.AXIS(KEY_A, KEY_D)
MOVEENTITY2D(p, camYaw, f, s, MOVE_SPEED, dt)
KEEPPLAYERINBOUNDS(p)
px = PLAYER2D.GETX(p)
pz = PLAYER2D.GETZ(p)
```

---

## Full Example

Third-person orbit with PLAYER2D ground movement and box-landing.

```basic
WINDOW.OPEN(960, 540, "GameHelpers Demo")
WINDOW.SETFPS(60)

cam = CAMERA.CREATE()
CAMERA.SETPOS(cam, 0, 8, -12)

CONST GRAVITY   = -26.0
CONST MOVE_SPEED = 5.0
CONST MOUSE_ORBIT_SENS = 0.005

p   = PLAYER2D.CREATE()
PLAYER2D.SETPOS(p, 0.0, 0.0)
CLAMPENTITY2D(p, -18.0, 18.0, -18.0, 18.0)

camYaw   = 0.0
camPitch = 0.4
camDist  = 10.0
vel_y    = 0.0
pos_y    = 0.5
on_ground = FALSE

; one platform box: centre, size
bx = 0.0 : by = 0.0 : bz = 4.0
bw = 4.0 : bh = 0.5 : bd = 4.0

WHILE NOT WINDOW.SHOULDCLOSE()
    dt = TIME.DELTA()

    camYaw   = camYaw   + ORBITYAWDELTA(dt, MOUSE_ORBIT_SENS, KEY_Q, KEY_E, 72.0)
    camPitch = camPitch + ORBITPITCHDELTA(MOUSE_ORBIT_SENS)
    camDist  = camDist  + ORBITDISTDELTA(0.85)
    camPitch = MAX(0.1, MIN(1.3, camPitch))
    camDist  = MAX(3.0, MIN(20.0, camDist))

    f = INPUT.AXIS(KEY_S, KEY_W)
    s = INPUT.AXIS(KEY_A, KEY_D)
    MOVEENTITY2D(p, camYaw, f, s, MOVE_SPEED, dt)
    KEEPPLAYERINBOUNDS(p)

    vel_y = vel_y + GRAVITY * dt
    pos_y = pos_y + vel_y * dt
    on_ground = FALSE

    landY = BOXTOPLAND(PLAYER2D.GETX(p), pos_y, PLAYER2D.GETZ(p), vel_y, 0.4, bx, by, bz, bw, bh, bd)
    IF landY > 0.0 THEN
        pos_y = landY : vel_y = 0.0 : on_ground = TRUE
    END IF
    IF pos_y < 0.4 THEN
        pos_y = 0.4 : vel_y = 0.0 : on_ground = TRUE
    END IF

    IF on_ground AND INPUT.KEYPRESSED(KEY_SPACE) THEN vel_y = 9.0

    px = PLAYER2D.GETX(p)
    pz = PLAYER2D.GETZ(p)
    CAMERA.SETORBIT(cam, px, pos_y, pz, camYaw, camPitch, camDist)

    RENDER.CLEAR(60, 80, 120)
    RENDER.BEGIN3D(cam)
        DRAW.SPHERE(px, pos_y, pz, 0.4, 80, 200, 255, 255)
        DRAW.CUBE(bx, by, bz, bw, bh, bd, 0, 120, 80, 40, 255)
        DRAW.GRID(20, 1.0)
    RENDER.END3D()
    RENDER.FRAME()
WEND

PLAYER2D.FREE(p)
WINDOW.CLOSE()
```

---

## See also

- [PLAYER2D.md](PLAYER2D.md) — `PLAYER2D.*` reference
- [PLAYER.md](PLAYER.md) — 3D Jolt KCC
- [CAMERA.md](CAMERA.md) — `CAMERA.SETORBIT`, `CAMERA.ORBITAROUND`
- [GAME.md](GAME.md) — `GAME.ORBITPITCHDELTA`, `GAME.ORBITYAWDELTA`
