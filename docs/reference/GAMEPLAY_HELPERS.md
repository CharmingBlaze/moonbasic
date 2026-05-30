# Gameplay Helpers

A collection of beginner-friendly building blocks for common game patterns: entity distance/movement, proximity triggers, and camera-relative controls.

Page shape follows [DOC_STYLE_GUIDE.md](../DOC_STYLE_GUIDE.md) (**WAVE pattern**).

## Core Workflow

1. **Proximity:** Use `ENTITY.WITHINRADIUS` or `ENTITY.DISTANCE` to detect nearby targets.
2. **Movement:** Drive behavior using `ENTITY.MOVETOWARD` or `ENTITY.MOVEWITHCAMERA`.
3. **Alignment:** Ensure entities face their targets with `ENTITY.TURNTOWARD` or `ENTITY.LOOKAT`.
4. **Camera:** Link the view to gameplay using `CAMERA.ORBITENTITY`.

---

## 1. Proximity and Distance

| Command | Arguments | Returns | Notes |
|---------|-----------|---------|-------|
| `ENTITY.DISTANCE(a, b)` | Integer, Integer | Float | 3D world distance between `a` and `b`. |
| `ENTITY.DISTANCETO(a, b)`| Integer, Integer | Float | Same as `DISTANCE`. |
| `ENTITY.WITHINRADIUS(a, b, r)`| Integer, Integer, Float| Boolean | `TRUE` if `b` is within `r` of `a`. |

---

## 2. Targeted Movement

| Command | Arguments | Returns | Notes |
|---------|-----------|---------|-------|
| `ENTITY.MOVETOWARD(e, target, s)`| Integer, Integer, Float | Handle | Moves entity `e` toward entity `target`. |
| `ENTITY.MOVETOWARD(e, x, z, s)` | Integer, Float, Float, Float| Handle | Moves entity `e` toward world `(x, z)`. |
| `ENTITY.TURNTOWARD(e, x, z, s)` | Integer, Float, Float, Float| Handle | Smoothly rotates `e` to face `(x, z)`. |
| `ENTITY.LOOKAT(e, x, z)` | Integer, Float, Float | Handle | Instantly faces `(x, z)`. |

---

## 3. Character and Camera

| Command | Arguments | Returns | Notes |
|---------|-----------|---------|-------|
| `ENTITY.MOVEWITHCAMERA(e, cam, f, s, speed)`| Int, Handle, Float...| Handle | Moves `e` relative to camera (WASD style). |
| `CAMERA.ORBITENTITY(cam, e, y, p, d)` | Handle, Int, Float...| None | Third-person orbit setup. |

---

## Full Example

This example creates a player that moves relative to the camera and a "follower" entity that trails the player.

```basic
WINDOW.OPEN(1280, 720, "Gameplay Helpers Demo")
cam = CAMERA.CREATE()
player = ENTITY.CREATECUBE(1.0).SETPOS(0, 0.5, 0)
follower = ENTITY.CREATESPHERE(0.5, 16).SETPOS(5, 0.5, 5)

WHILE NOT WINDOW.SHOULDCLOSE()
    dt = TIME.DELTA()
    
    ; 1. Player movement (WASD)
    forward = INPUT.AXIS(KEY_S, KEY_W)
    strafe = INPUT.AXIS(KEY_A, KEY_D)
    player.MOVEWITHCAMERA(cam, forward, strafe, 6.0)

    ; 2. Follower logic
    IF NOT ENTITY.WITHINRADIUS(follower, player, 2.0)
        follower.MOVETOWARD(player, 4.0)
        follower.TURNTOWARD(player.X(), player.Z(), 5.0)
    END IF

    ; 3. Camera setup
    CAMERA.ORBITENTITY(cam, player, 0, 0.5, 10.0)
    
    RENDER.CLEAR(30, 30, 40)
    RENDER.BEGIN3D(cam)
        ENTITY.DRAWALL()
        DRAW3D.GRID(50, 1.0)
    RENDER.END3D()
    RENDER.FRAME()
WEND
```

## See also

- [BEGINNER_FULL_STACK.md](BEGINNER_FULL_STACK.md) — high-level bridge commands
- [ENTITY.md](ENTITY.md) — full entity API reference
- [CAMERA.md](CAMERA.md) — camera modes and projection
