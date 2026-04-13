# VEHICLE

The `VEHICLE` namespace provides a high-level raycast-based vehicle simulation system. It simplifies the creation of cars, trucks, and other wheeled vehicles by automating suspension and traction calculations.

## Commands

### VEHICLE.CREATE(entity, wheelCount)
Creates a new vehicle simulation bound to the specified entity (chassis).
* `entity`: Handle to the entity that serves as the vehicle chassis.
* `wheelCount`: The number of wheels the vehicle will have.

### VEHICLE.SETWHEEL(vehicle, index, ox#, oy#, oz#, radius#)
Configures a specific wheel relative to the chassis.
* `vehicle`: Handle to the vehicle.
* `index`: Index of the wheel (0 to wheelCount-1).
* `ox#, oy#, oz#`: Local offset from the chassis center.
* `radius#`: Radius of the wheel.

### VEHICLE.CONTROL(vehicle, throttle#, steer#, brake#)
Applies control inputs to the vehicle.
* `throttle#`: Acceleration input (-1 to 1).
* `steer#`: Steering input (-1 to 1).
* `brake#`: Braking input (0 to 1).

### VEHICLE.STEP(dt#)
Advances the vehicle simulation by a time step. Usually called once per frame with `DeltaTime()`.
* `dt#`: The time delta in seconds.

## Example: Simple Car Setup

```basic
; Load Chassis
car = Model.Load("models/car_chassis.glb")
World.Setup(-9.81)

; Create Vehicle with 4 wheels
v = Vehicle.Create(car, 4)

; Configure Wheels (Front Left, Front Right, Rear Left, Rear Right)
Vehicle.SetWheel(v, 0, -1.0, 0,  1.5, 0.4)
Vehicle.SetWheel(v, 1,  1.0, 0,  1.5, 0.4)
Vehicle.SetWheel(v, 2, -1.0, 0, -1.5, 0.4)
Vehicle.SetWheel(v, 3,  1.0, 0, -1.5, 0.4)

Repeat
    ; Basic Driving Input
    steer = Input.AxisX()
    gas = Input.AxisY()
    
    Vehicle.Control(v, gas, steer, 0)
    Vehicle.Step(DeltaTime())
    
    RenderWorld()
Until KeyHit(1)
```

> [!TIP]
> Use `LEVEL.STATIC` on your ground entities to ensure the vehicle raycasts can detect the floor correctly.
