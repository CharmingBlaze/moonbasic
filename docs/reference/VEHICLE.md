# Vehicle Commands

High-level raycast-based vehicle simulation with suspension and traction helpers.

Page shape follows [DOC_STYLE_GUIDE.md](../DOC_STYLE_GUIDE.md) (**WAVE pattern**).

## Core Workflow

1. Create a vehicle with `VEHICLE.CREATE`, binding it to a chassis entity.
2. Configure wheels with `VEHICLE.SETWHEEL`.
3. Each frame, apply inputs with `VEHICLE.CONTROL` and step with `VEHICLE.STEP`.

---

### `VEHICLE.CREATE(entity, wheelCount)` 
Creates a new vehicle simulation bound to the specified chassis **entity**.
* `entity`: Numeric **entity id** of the chassis (not a raw model handle — spawn or reference an entity first; see [ENTITY.md](ENTITY.md)).
* `wheelCount`: Number of wheels.

---

### `VEHICLE.SETWHEEL(vehicle, index, ox, oy, oz, radius)` 
Configures one wheel relative to the chassis.
* `vehicle`: Vehicle handle returned from **`VEHICLE.CREATE`**.
* `index`: Wheel index from **0** to **wheelCount − 1**.
* `ox`, `oy`, `oz`: Local offset from the chassis center.
* `radius`: Wheel radius.

---

### `VEHICLE.CONTROL(vehicle, throttle, steer, brake)` 
Applies control inputs.
* `throttle`: Acceleration input (-1 to 1).
* `steer`: Steering input (-1 to 1).
* `brake`: Braking input (0 to 1).

---

### `VEHICLE.STEP(dt)` 
Advances the simulation by **`dt`** seconds. Call once per frame with **`TIME.DELTA()`** (or **`DT()`**).

## Full Example

```basic
WORLD.SETUP(-9.81)

carEnt = 1
v = VEHICLE.CREATE(carEnt, 4)

VEHICLE.SETWHEEL(v, 0, -1.0, 0,  1.5, 0.4)
VEHICLE.SETWHEEL(v, 1,  1.0, 0,  1.5, 0.4)
VEHICLE.SETWHEEL(v, 2, -1.0, 0, -1.5, 0.4)
VEHICLE.SETWHEEL(v, 3,  1.0, 0, -1.5, 0.4)

WHILE NOT WINDOW.SHOULDCLOSE()
    steer = INPUT.AXIS(KEY_LEFT, KEY_RIGHT)
    gas = INPUT.AXIS(KEY_UP, KEY_DOWN)
    VEHICLE.CONTROL(v, gas, steer, 0.0)
    VEHICLE.STEP(TIME.DELTA())
    RENDER.FRAME()
WEND
```

---

## Extended Command Reference

| Command | Description |
|--------|-------------|
| `VEHICLE.MAKE(entity)` | Deprecated alias of `VEHICLE.CREATE`. |
| `VEHICLE.SETTHROTTLE(v, value)` | Set throttle -1.0 (reverse) to 1.0 (forward). |
| `VEHICLE.SETSTEER(v, angle)` | Set steering angle in degrees. |
| `VEHICLE.SETTUNING(v, key, value)` | Set a tuning parameter by name (e.g. `"maxSpeed"`). |
| `VEHICLE.WHEELX(v, wheel)` / `WHEELY` / `WHEELZ` | World position of wheel `wheel` per axis. |

## See also

- [PHYSICS3D.md](PHYSICS3D.md) — physics world setup
- [BODY3D.md](BODY3D.md) — rigid body for chassis

> [!TIP]
> Use **`LEVEL.STATIC`** on ground entities so vehicle raycasts can detect the floor.
