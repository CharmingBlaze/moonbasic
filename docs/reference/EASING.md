# Easing Commands

Interpolation and easing functions for smooth animation, UI transitions, and procedural motion.

Page shape follows [DOC_STYLE_GUIDE.md](../DOC_STYLE_GUIDE.md) (**WAVE pattern**).

## Core Workflow

All easing functions take three floats: `(t, start, end)` where `t` is a normalised progress value (0.0–1.0), `start` is the value at t=0, and `end` is the value at t=1. Each returns the interpolated result.

Pick a curve that matches the feel you want — **linear** (`EASELERP`), **quadratic** (`EASEIN`/`EASEOUT`), **cubic** (`EASEIN3`/`EASEOUT3`), **sine**, **back**, **elastic**, or **bounce**.

---

### `EASELERP(t, start, end)` 

Linear interpolation (no easing). Equivalent to `start + (end - start) * t`.

---

### `EASEIN(t, start, end)` 

Quadratic ease-in — starts slow, accelerates.

---

### `EASEOUT(t, start, end)` 

Quadratic ease-out — starts fast, decelerates.

---

### `EASEINOUT(t, start, end)` 

Quadratic ease-in-out — slow at both ends, fast in the middle.

---

### `EASEIN3(t, start, end)` 

Cubic ease-in — stronger acceleration than quadratic.

---

### `EASEOUT3(t, start, end)` 

Cubic ease-out — stronger deceleration than quadratic.

---

### `EASEINOUT3(t, start, end)` 

Cubic ease-in-out.

---

### `EASEINSINE(t, start, end)` 

Sine-based ease-in.

---

### `EASEOUTSINE(t, start, end)` 

Sine-based ease-out.

---

### `EASEINOUTSINE(t, start, end)` 

Sine-based ease-in-out.

---

### `EASEINBACK(t, start, end)` 

Overshoots slightly past `start` before accelerating toward `end`.

---

### `EASEOUTBACK(t, start, end)` 

Overshoots slightly past `end` before settling.

---

### `EASEINELASTIC(t, start, end)` 

Elastic spring effect on entry.

---

### `EASEOUTELASTIC(t, start, end)` 

Elastic spring effect on exit.

---

### `EASEINBOUNCE(t, start, end)` 

Bouncing effect on entry.

---

### `EASEOUTBOUNCE(t, start, end)` 

Bouncing effect on exit.

---

## Full Example

This example animates a rectangle sliding across the screen with different easing curves.

```basic
startX = 50.0
endX   = 700.0
duration = 2.0
timer  = 0.0

WHILE NOT WINDOW.SHOULDCLOSE()
    dt = DELTATIME()
    timer = timer + dt
    t = MATH.CLAMP(timer / duration, 0.0, 1.0)

    ; Compare easing curves side by side
    y1 = EASELERP(t, startX, endX)
    y2 = EASEIN(t, startX, endX)
    y3 = EASEOUTBOUNCE(t, startX, endX)
    y4 = EASEOUTELASTIC(t, startX, endX)

    RENDER.BEGINFRAME()
    DRAW.RECT(INT(y1), 50,  20, 20, 255, 0, 0, 255)
    DRAW.RECT(INT(y2), 100, 20, 20, 0, 255, 0, 255)
    DRAW.RECT(INT(y3), 150, 20, 20, 0, 0, 255, 255)
    DRAW.RECT(INT(y4), 200, 20, 20, 255, 255, 0, 255)
    RENDER.ENDFRAME()

    ; Loop
    IF timer > duration THEN timer = 0.0
WEND
```
