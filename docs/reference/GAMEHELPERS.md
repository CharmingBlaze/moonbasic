# Game helpers (collision and patterns)

Small built-ins and idioms that keep gameplay code short without a full physics engine.

---

## `BOXTOPLAND`

`BOXTOPLAND(px#, py#, pz#, pvy#, pr#, bx#, by#, bz#, bw#, bh#, bd#)` → **float**

Returns **`0.0`** when there is **no** top landing this frame. Otherwise returns the **sphere center Y** to snap to (`box top + pr`).

- Only meaningful when **`pvy <= 0`** (falling or resting). If moving upward, returns **`0.0`**.
- Horizontal test: sphere center must be within the box footprint, expanded by **`pr`** on X/Z.
- Vertical test: **feet** (`py - pr`) must sit in a small band below/around the **top** of the box (`by + bh/2`).

Typical use:

```basic
landY# = BOXTOPLAND(px#, py#, pz#, pvy#, pr#, bx#, by#, bz#, bw#, bh#, bd#)
IF landY# > 0.0 THEN
    py# = landY# : pvy# = 0.0 : on_ground? = TRUE
ENDIF
```

---

## `LANDBOXES` / `LANDBOX`

`LANDBOXES(px#, py#, pz#, pvy#, pr#, plx#, ply#, plz#, plw#, plh#, pld#, count)` → **float**

**`LANDBOX`** is an **alias** — same arguments and return value.

Runs the same test as **`BOXTOPLAND`** for **`count`** boxes given as **six parallel float arrays** (centre `x,y,z` and size `w,h,d`). Returns the **largest** positive snap Y among all boxes, or **`0.0`** if none apply. Use this instead of a **`FOR`** loop when platforms are stored as parallel **`DIM`** arrays.

Implementation note: it is equivalent to **`BOXTOPLAND`** per index — not a full physics engine. **`TYPE`** platform rows still use a loop or manual **`BOXTOPLAND`** unless you keep parallel arrays for collision.

---

## `PLAYER.MOVERELATIVE`

`PLAYER.MOVERELATIVE(camYaw#, forward#, strafe#, speed#, dt#)` → **handle** (2-float array **`[deltaX, deltaZ]`**)

Same math as **`MOVESTEPX`** and **`MOVESTEPZ`** combined. **Free** the returned array with **`ERASE`** when you are done (each frame if you allocate every frame). For hot loops, **`MOVESTEPX`/`MOVESTEPZ`** avoid the extra heap array.

---

## Simple physics without a physics engine

Gravity and integration are only a few lines. Keep **`dt#`** from **`Time.Delta()`** or **`DT()`** (both are **clamped** by default so tab-switch spikes do not explode simulation).

```basic
CONST GRAVITY# = -26.0

; Each frame:
vel_y# = vel_y# + GRAVITY# * dt#
pos_y# = pos_y# + vel_y# * dt#

; Ground check (flat floor at y = radius):
IF pos_y# < radius# THEN
    pos_y# = radius#
    vel_y# = 0.0
    on_ground? = TRUE
ENDIF
```

For **one-shot** actions (jump, shoot), use **`Input.KeyPressed`** or the flat **`KEYPRESSED`** helper, not **`KeyDown`** / **`KEYDOWN`**, which fire every frame the key is held.

See also: [INPUT.md](INPUT.md) (keyboard table), [CAMERA.md](CAMERA.md) (`Camera.OrbitAround` for third-person orbit).

---

## Third-person orbit input (`ORBITYAWDELTA` / `ORBITPITCHDELTA` / `ORBITDISTDELTA`)

These **`GAME.*`** helpers (short names without the `GAME.` prefix also work) wrap **right-mouse drag** and **mouse wheel** together with the same **Q/E yaw** math as **`Input.Orbit`**. They return **plain floats** each frame — **no heap handles**, nothing to **`ERASE`**. Use them to update your **`camYaw#`**, **`camPitch#`**, and **`camDist#`**, then call **`Camera.SetOrbit`** (or **`Camera.OrbitAround`**) yourself.

| Command | Returns | Meaning |
|--------|---------|--------|
| **`ORBITYAWDELTA(dt#, mouseSens#, negKey, posKey, degPerSec#)`** | radians | **Keyboard:** same as **`Input.Orbit(negKey, posKey, degPerSec#, dt#)`** (degrees/sec → radians). **Mouse:** if **right button** is down, adds **`MDX * mouseSens`** (typically `mouseSens` ≈ `0.004`–`0.006`). |
| **`ORBITPITCHDELTA(mouseSens#)`** | radians | If **right button** is down: **`-MDY * mouseSens`**. Otherwise **`0`**. |
| **`ORBITDISTDELTA(wheelScale#)`** | world units | **`-MWHEEL * wheelScale`** — add to your orbit distance (scroll **up** moves the eye **closer** when **`wheelScale`** is positive). |

Clamp **`pitch`** and **`dist`** in your script after adding deltas (the helpers do not clamp).

Typical frame (see **`examples/mario64/main_orbit_simple.mb`**):

```basic
camYaw# = camYaw# + ORBITYAWDELTA(dt#, 0.0048, KEY_Q, KEY_E, 72.0)
camPitch# = camPitch# + ORBITPITCHDELTA(0.0048)
camDist# = camDist# + ORBITDISTDELTA(0.85)
; … clamp pitch & dist, then Camera.SetOrbit(cam, tx, ty, tz, camYaw#, camPitch#, camDist#)
```

That example is structured for reading **top to bottom**: one **`CONST`** block (world bounds, orbit tuning, colours), parallel **`DIM`** rows for **`LANDBOXES`**, a single loop section for input → physics → **`Camera.SetOrbit`** → draw, then **`ERASE ALL`** (see [MEMORY.md](../MEMORY.md)).

**Memory:** no allocations — see [MEMORY.md](../MEMORY.md) (game orbit helpers).
