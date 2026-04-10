# Full-screen transitions (`TRANSITION.*`)

Fullscreen **fade** and **wipe** overlays drawn at the end of the frame (Raylib), used standalone or with [SCENE.LoadWithTransition](SCENE.md). **CGO** only (`runtime/mbtransition/transition_cgo.go`); a frame hook must be registered on the window module so `transitionDraw` runs each frame.

---

## State machine

1. Call **`TRANSITION.FADEOUT`** or **`TRANSITION.WIPE`** with a duration.
2. Each frame, the hook advances time until progress reaches **1** → **`TRANSITION.ISDONE`** becomes `TRUE`.
3. For fade-out, the screen ends fully covered by the transition color; for **fade-in**, the overlay clears back to transparent and the mode returns to idle.

---

## Commands

### `Transition.FadeOut(seconds#)` / `Transition.FadeIn(seconds#)`

Fades the overlay **in** (blocking the view) or **out** (revealing the scene). Duration must be positive.

### `Transition.IsDone()` → bool

`TRUE` when the current transition segment has finished.

### `Transition.Wipe(direction$, seconds#)`

Wipe overlay using the current transition **color**. `direction` (case-insensitive):

- `left` — bar grows from the left edge.
- `right` — from the right.
- `up` or `top` — from the top.
- `down` or `bottom` — from the bottom.
- Anything else — full-screen fill.

### `Transition.SetColor(r, g, b, a)`

Overlay color (components clamped **0–255**). Default is set in the CGO init (typically opaque black).

---

## Scene integration

[SCENE.LoadWithTransition](SCENE.md) calls `FADEOUT` or `WIPE`, waits for `ISDONE`, runs the scene loader, then `FADEIN` for `"fade"` wipes that need a fade-in phase.
