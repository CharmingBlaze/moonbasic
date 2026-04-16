# VEHICLE

**Conventions:** [STYLE_GUIDE.md](../../STYLE_GUIDE.md), [API_CONVENTIONS.md](API_CONVENTIONS.md) — use **`VEHICLE.*`** registry keys in new examples; parameter names below are plain identifiers (no Blitz-style **`#`** suffixes).

The **`VEHICLE`** namespace provides a high-level raycast-based vehicle simulation: cars, trucks, and other wheeled vehicles with suspension and traction helpers.

## Commands

### `VEHICLE.CREATE(entity, wheelCount)`
Creates a new vehicle simulation bound to the specified chassis **entity**.
* `entity`: Numeric **entity id** of the chassis (not a raw model handle — spawn or reference an entity first; see [ENTITY.md](ENTITY.md)).
* `wheelCount`: Number of wheels.

### `VEHICLE.SETWHEEL(vehicle, index, ox, oy, oz, radius)`
Configures one wheel relative to the chassis.
* `vehicle`: Vehicle handle returned from **`VEHICLE.CREATE`**.
* `index`: Wheel index from **0** to **wheelCount − 1**.
* `ox`, `oy`, `oz`: Local offset from the chassis center.
* `radius`: Wheel radius.

### `VEHICLE.CONTROL(vehicle, throttle, steer, brake)`
Applies control inputs.
* `throttle`: Acceleration input (-1 to 1).
* `steer`: Steering input (-1 to 1).
* `brake`: Braking input (0 to 1).

### `VEHICLE.STEP(dt)`
Advances the simulation by **`dt`** seconds. Call once per frame with **`TIME.DELTA()`** (or **`DT()`**).

## Example: simple car setup

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

> [!TIP]
> Use **`LEVEL.STATIC`** on ground entities so vehicle raycasts can detect the floor.
