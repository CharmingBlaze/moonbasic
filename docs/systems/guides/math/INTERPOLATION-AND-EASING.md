# Interpolation and easing — smooth motion and UI

> Move values smoothly with lerp, approach, smoothstep, remap, and wrap — without jitter or overshoot bugs.

**Namespaces:** `MATH` · **Status:** Shipped

**Commands:** [COMMAND_REGISTRY.md#data](../../COMMAND_REGISTRY.md#data) · **2D vectors:** [VEC2-MATH.md](VEC2-MATH.md) (`LERP`) · **Reference:** [reference/MATH.md](../../../reference/MATH.md)

---

## Table of contents

- [At a glance](#at-a-glance)
- [When to use which interpolator](#when-to-use-which-interpolator)
- [Core workflow](#core-workflow)
- [Command guide](#command-guide)
- [Frame-rate independence](#frame-rate-independence)
- [Full example — UI bar + platform](#full-example--ui-bar--platform)
- [Common mistakes](#common-mistakes)
- [See also](#see-also)

---

## At a glance

| Idea | Detail |
|------|--------|
| **You get** | `LERP`, `APPROACH`, `SMOOTHSTEP`, `REMAP`, `INVERSE_LERP`, `PINGPONG`, `WRAP` |
| **Typical use** | Health bar display, camera lag, color pulse, moving platforms |
| **Not for** | Physics velocity — engine integrates forces |

---

## When to use which interpolator

| Goal | Use | Why |
|------|-----|-----|
| Blend two endpoints by % | `MATH.LERP(a, b, t)` | Simple mix; `t` 0→1 |
| Move max N per frame | `MATH.APPROACH(cur, target, step)` | Never overshoot |
| Ease in/out UI | `MATH.SMOOTHSTEP` / `SMOOTHERSTEP` | Zero derivative at edges |
| Exponential decay feel | `MATH.CURVE(cur, target, factor)` | Common “smooth follow” |
| Map slider to range | `MATH.REMAP(v, in0, in1, out0, out1)` | One formula |
| Find t from value | `MATH.INVERSE_LERP(a, b, v)` | Reverse lerp |
| Oscillate 0…length | `MATH.PINGPONG(t, len)` | Platforms, sway |
| Wrap index / torus | `MATH.WRAP(v, min, max)` | Array loop, wrap world |
| Clamp 0…1 | `MATH.SATURATE(v)` | Before smoothstep |

**Angles:** `MATH.LERPANGLE` — [ANGLES-AND-ROTATION.md](ANGLES-AND-ROTATION.md).

---

## Core workflow

1. **Pick model** — approach for chase speed cap; lerp for fixed blend factor; smoothstep for timed UI.
2. **Scale by delta** — `step = maxSpeed * APP.DELTA()` for approach per frame.
3. **Saturate** — `t = MATH.SATURATE(t)` before smoothstep.
4. **Display vs sim** — lerp **shown** health toward **real** health for juice.

---

## Command guide

```basic
; Approach — never overshoot target
hpDisplay = MATH.APPROACH(hpDisplay, hpReal, 50 * APP.DELTA())

; Lerp — fixed blend (often scale t by dt for chase)
camX = MATH.LERP(camX, targetX, 5 * APP.DELTA())

; Smooth UI fade
alpha = MATH.SMOOTHSTEP(0, 1, MATH.SATURATE(timer / duration))

; Remap input axis -1..1 to 0..255
bright = MATH.REMAP(stick, -1, 1, 0, 255)

; Ping-pong platform offset
offset = MATH.PINGPONG(APP.TIME() * 2, 4)
```

---

## Frame-rate independence

| Pattern | Formula |
|---------|---------|
| Approach speed | `APPROACH(v, target, rate * APP.DELTA())` |
| Lerp chase | `LERP(v, target, k * APP.DELTA())` — clamp k*dt ≤ 1 |
| Timed tween | `t = elapsed / duration` then `SMOOTHSTEP(0,1,t)` |

Without `DELTA()`, 120 FPS moves twice as fast as 60 FPS.

---

## Full example — UI bar + platform

```basic
APP.OPEN(500, 300, "Easing")
APP.SETFPS(60)

hpReal = 100
hpShow = 100
pulseT = 0
platX = 100

WHILE NOT APP.SHOULDCLOSE()
    IF INPUT.KEYHIT(KEY_SPACE) THEN hpReal = hpReal - MATH.RAND(5, 20)
    hpReal = MATH.CLAMP(hpReal, 0, 100)

    hpShow = MATH.APPROACH(hpShow, hpReal, 40 * APP.DELTA())
    pulseT = pulseT + APP.DELTA()
    pulse = MATH.SMOOTHERSTEP(0, 1, MATH.PINGPONG(pulseT, 1))
    platX = 100 + MATH.PINGPONG(APP.TIME() * 80, 200)

    barW = INT(MATH.REMAP(hpShow, 0, 100, 0, 300))
    DRAW.RECTANGLE(50, 40, barW, 16, 80, 200, 80)
    DRAW.RECTANGLE(INT(platX), 200, 60, 12, INT(100 + 155 * pulse), 100, 200)
    RENDER.FRAME()
WEND
APP.CLOSE()
```

---

## Common mistakes

| Mistake | Fix |
|---------|-----|
| Lerp factor > 1 | Clamp `t` or use approach |
| No delta on approach | Multiply step by `APP.DELTA()` |
| Smoothstep outside 0..1 | `SATURATE` input first |
| Lerp angles raw | Use `LERPANGLE` |
| Same as physics | Interpolation is kinematic display |

---

## See also

- [MATH-2D-GAMEPLAY.md](MATH-2D-GAMEPLAY.md) — chase with lerp on scalars
- [UI-MENUS.md](../UI-MENUS.md) — HUD animation
- [reference/GAME_MATH_HELPERS.md](../../../reference/GAME_MATH_HELPERS.md)
