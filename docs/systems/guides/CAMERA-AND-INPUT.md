# Camera and player controls

> Aim the view, orbit around a target, follow the hero, and tie **WASD** to **camera-relative** movement.

**Namespaces:** `CAMERA` · `INPUT` · `ACTION` · **Status:** Shipped

**Commands:** [COMMAND_REGISTRY.md#camera-light](../COMMAND_REGISTRY.md#camera-light) · [04-INPUT.md](../04-INPUT.md) · [02-CAMERA-LIGHT.md](../02-CAMERA-LIGHT.md)

---

## Table of contents

- [At a glance](#at-a-glance)
- [Choose your camera mode](#choose-your-camera-mode)
- [Static aim camera](#static-aim-camera)
- [Orbit camera (third person)](#orbit-camera-third-person)
- [Follow camera](#follow-camera)
- [First-person look](#first-person-look)
- [Input → movement pipeline](#input--movement-pipeline)
- [Common mistakes](#common-mistakes)
- [See also](#see-also)

---

## At a glance

| Mode | Typical commands | Game type |
|------|------------------|-----------|
| **Fixed look-at** | `SETPOS` + `LOOKAT` | Showcase, simple demos |
| **Orbit** | `CAMERA.ORBIT`, mouse delta | Third-person action |
| **Follow** | `CAMERA.FOLLOW` / `FOLLOWENTITY` | Chase cam |
| **FPS** | `INPUT.MOUSEDELTA` + yaw/pitch on camera | Shooters |

**Why camera is separate from entity:** The **view** (projection) can move without moving the **hero** mesh — essential for orbit and FPS.

---

## Choose your camera mode

| I want… | Start with |
|---------|------------|
| Debug cube spin | Fixed `LOOKAT` at origin |
| Mario-style orbit | `ORBIT` + RMB or `USEMOUSEORBIT` |
| Behind-player follow | `FOLLOW(entity, offset…)` |
| FPS | Lock mouse + rotate camera entity |

---

## Static aim camera

**Why:** Simplest — one position, stare at gameplay focal point.

```basic
cam = CAMERA.CREATE()
CAMERA.SETACTIVE(cam)
CAMERA.SETPOS(cam, 0, 4, -10)
CAMERA.LOOKAT(cam, 0, 0, 0)
```

Call each frame before `RENDER.BEGIN(cam)`.

---

## Orbit camera (third person)

**Why:** Player sees their character; mouse orbits around target.

```basic
CAMERA.ORBIT(cam, hero, 8.0)           ; distance
CAMERA.USEMOUSEORBIT(cam, true)        ; RMB drag — see reference
```

Or manual: read `INPUT.MOUSEDELTA_X/Y`, adjust yaw/pitch, `CAMERA.ORBITAROUND`.

Sample: [`examples/terrain_chase`](../../../examples/terrain_chase/main.mb) (orbit + WASD).

---

## Follow camera

**Why:** Camera lags behind moving target — good for racers and action.

```basic
CAMERA.FOLLOW(cam, hero, 0, 3, -8)     ; offset from target
; or
CAMERA.FOLLOWENTITY(cam, hero, 10, 3, 5)
```

Call **every frame** after hero moves.

---

## First-person look

**Why:** Mouse drives view direction; body may be invisible or weapon-only.

```basic
INPUT.LOCKMOUSE(true)
yaw = yaw - INPUT.MOUSEDELTA_X() * 0.003
pitch = pitch - INPUT.MOUSEDELTA_Y() * 0.003
CAMERA.SETROT(cam, pitch, yaw, 0)
```

Pair with `CHAR.MOVEWITHCAMERA` for walk direction ([CHARACTER-3D-WALKING.md](CHARACTER-3D-WALKING.md)).

---

## Input → movement pipeline

**Why layers:** Raw keys change per OS; **actions** stay stable in gameplay code.

1. **Bind once:** `ACTION.MAPKEY("Forward", KEY_W)` …
2. **Each frame:** `ACTION.DOWN("Forward")` → apply velocity.
3. **Scale by** `APP.DELTA()` so speed is frame-independent.
4. **Camera-relative:** `INPUT.MOVEDIR()` or `CHAR.MOVEWITHCAMERA`.

```basic
speed = 5
IF ACTION.DOWN("Forward") THEN hero.move(0, 0, speed * APP.DELTA())
```

See [04-INPUT.md](../04-INPUT.md).

---

## Common mistakes

| Mistake | Fix |
|---------|-----|
| Forget `SETACTIVE` | `RENDER.BEGIN()` uses wrong camera |
| Orbit without updating each frame | Call orbit/follow in loop |
| Movement without `DELTA()` | Speed tied to FPS |
| Mouse FPS without lock | Cursor leaves window |

---

## See also

- [02-CAMERA-LIGHT.md](../02-CAMERA-LIGHT.md)
- [CHARACTER-3D-WALKING.md](CHARACTER-3D-WALKING.md)
- [reference/CAMERA.md](../../reference/CAMERA.md)
