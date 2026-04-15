# 2D Physics Commands

Commands for creating and controlling a 2D physics simulation using Box2D.

**Registry names** use the **`PHYSICS2D.*`** / **`BODY2D.*`** / **`JOINT2D.*`** prefixes (e.g. **`PHYSICS2D.START`**, **`BODY2D.CREATE`**, deprecated **`BODY2D.MAKE`**). A full teaching-oriented map (Designed → Implementation → memory) is in [moonbasic-command-set/physics-2d.md](moonbasic-command-set/physics-2d.md).

## Core Workflow

1.  **Initialize**: Start the physics world with `Physics2D.Start()`.
2.  **Create Bodies**: Define and create physics bodies (`BODY2D.CREATE`, `BODY2D.ADDRECT`, `BODY2D.COMMIT`; deprecated `BODY2D.MAKE` alias).
3.  **Update**: Advance the simulation each frame with `Physics2D.Step()`.
4.  **Synchronize**: Use `Body2D.X()` and `Body2D.Y()` to update the positions of your visual shapes.
5.  **Cleanup**: Shut down the world with `Physics2D.Stop()`.

---

## World Management

### `Physics2D.Start([gx, gy])`
Initializes the 2D physics world. Default gravity is `(0, 500)` if omitted.

### `Physics2D.Stop()`
Shuts down the 2D physics simulation and frees internal buffers.

---

### `Physics2D.Step()`
Advances the 2D simulation (call once per frame).

### `Physics2D.SetGravity(x, y)`
Sets the global gravity vector for the 2D world.

---

### `BODY2D.CREATE(type)` (canonical; deprecated `BODY2D.MAKE`)
Creates a body definition (`"static"`, `"dynamic"`, `"kinematic"`).

### `BODY2D.ADDRECT(def, w, h)`
Adds a rectangle collision shape to the definition.

### `BODY2D.COMMIT(def, x, y)`
Finalizes the body and adds it to the world at the specified position. Returns a **body handle**.

---

### `BODY2D.SETPOS(handle, x, y)` (canonical; deprecated `BODY2D.SETPOSITION`)
Teleports a 2D body to a new position.

### `BODY2D.X(handle)` / `BODY2D.Y(handle)`
Returns the current X or Y coordinate of the body's center.

### `BODY2D.ROT(handle)`
Returns the body rotation in radians.

### `BODY2D.FREE(handle)`
Removes a body from the simulation and frees its memory.

---

## Full Example: Falling Box

```basic
Window.Open(800, 600, "2D Physics Example")
Window.SetFPS(60)

; 1. Initialize Physics World
Physics2D.Start()
Physics2D.SetGravity(0, 500) ; Positive Y is down in 2D

; 2. Create a static floor
floor_def = BODY2D.CREATE("static")
BODY2D.ADDRECT(floor_def, 800, 50)
floor_body = BODY2D.COMMIT(floor_def, 400, 575)

; 3. Create a dynamic box
box_def = BODY2D.CREATE("dynamic")
BODY2D.ADDRECT(box_def, 40, 40)
box_body = BODY2D.COMMIT(box_def, 400, 100)

WHILE NOT Window.ShouldClose()
    ; 4. Update simulation
    Physics2D.Step()

    Render.Clear(10, 20, 30)
    Camera2D.Begin()
        ; 5. Synchronize visuals
        box_x = BODY2D.X(box_body)
        box_y = BODY2D.Y(box_body)
        box_rot = BODY2D.ROT(box_body)

        ; Draw floor
        Draw.Rectangle(INT(BODY2D.X(floor_body)) - 400, INT(BODY2D.Y(floor_body)) - 25, 800, 50, 100, 100, 100, 255)
        ; Draw box (rotation not visually applied without a sprite, but position is correct)
        Draw.Rectangle(INT(box_x) - 20, INT(box_y) - 20, 40, 40, 200, 50, 50, 255)

    Camera2D.End()
    Render.Frame()
WEND

; 6. Cleanup
Physics2D.Stop()
Window.Close()
```
