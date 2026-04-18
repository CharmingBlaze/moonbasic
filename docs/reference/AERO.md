# Aero Commands

Aerodynamics helpers for physics bodies: lift coefficient, local-Z thrust, and air drag. Intended for vehicles, projectiles, and aircraft simulations.

Requires **CGO + Jolt** physics bodies.

## Core Workflow

1. Create a `BODY3D` dynamic body.
2. `AERO.SETLIFT(body, coefficient)` — add lift proportional to speed.
3. `AERO.SETTHRUST(body, power)` — apply thrust along the body's local +Z axis each step.
4. `AERO.SETDRAG(body, coefficient)` — add velocity-proportional drag.
5. Call in your frame loop before `PHYSICS3D.UPDATE`.

---

## Commands

### `AERO.SETLIFT(body, coefficient)` 

Sets a lift coefficient for a physics body. Lift force is applied upward proportional to the body's speed squared: `lift = coefficient × speed²`. Simulates wing lift or buoyancy.

---

### `AERO.SETTHRUST(body, power)` 

Applies a continuous force along the body's local **+Z** axis. Call each frame to simulate engine thrust. `power` is in Newtons.

---

### `AERO.SETDRAG(body, coefficient)` 

Applies air resistance: a force opposing the body's velocity, proportional to speed. `drag = coefficient × speed`. Higher values slow the body down faster.

---

## Full Example

Simple airplane: thrust forward, lift upward, drag to limit top speed.

```basic
WINDOW.OPEN(960, 540, "Aero Demo")
WINDOW.SETFPS(60)

PHYSICS3D.START()
PHYSICS3D.SETGRAVITY(0, -4, 0)

planeDef = BODY3D.CREATE("DYNAMIC")
BODY3D.ADDBOX(planeDef, 1.5, 0.2, 0.6)
BODY3D.SETMASS(planeDef, 2.0)
plane = BODY3D.COMMIT(planeDef, 0, 5, 0)
BODY3D.LOCKAXIS(plane, 8 + 32)   ; lock X and Z rotation (keep upright)

AERO.SETLIFT(plane,   3.5)   ; wing lift
AERO.SETTHRUST(plane, 15.0)  ; engine
AERO.SETDRAG(plane,   1.2)   ; air resistance

cam = CAMERA.CREATE()
CAMERA.SETPOS(cam, 0, 8, -16)
CAMERA.SETTARGET(cam, 0, 4, 0)

WHILE NOT WINDOW.SHOULDCLOSE()
    dt = TIME.DELTA()

    PHYSICS3D.UPDATE()

    px = BODY3D.X(plane)
    py = BODY3D.Y(plane)
    pz = BODY3D.Z(plane)
    CAMERA.SETTARGET(cam, px, py, pz)
    CAMERA.SETPOS(cam, px, py + 6, pz - 14)

    RENDER.CLEAR(100, 160, 220)
    RENDER.BEGIN3D(cam)
        DRAW3D.CUBE(px, py, pz, 3, 0.4, 1.2, 200, 200, 220, 255)
        DRAW3D.GRID(40, 1.0)
    RENDER.END3D()
    DRAW.TEXT("Speed: " + STR(INT(BODY3D.VELOCITY(plane))), 10, 10, 18, 255, 255, 255, 255)
    RENDER.FRAME()
WEND

BODY3D.FREE(plane)
PHYSICS3D.STOP()
WINDOW.CLOSE()
```

---

## See also

- [BODY3D.md](BODY3D.md) — rigid bodies
- [PHYSICS3D.md](PHYSICS3D.md) — world setup
- [VEHICLE.md](VEHICLE.md) — wheeled vehicle controller
