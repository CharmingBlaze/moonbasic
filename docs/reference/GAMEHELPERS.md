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

## `LANDBOXES`

`LANDBOXES(px#, py#, pz#, pvy#, pr#, plx#, ply#, plz#, plw#, plh#, pld#, count)` → **float**

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
