# 2D Physics Commands

Commands for creating and controlling a 2D physics simulation using Box2D.

**Registry names** use the **`PHYSICS2D.*`** / **`BODY2D.*`** / **`JOINT2D.*`** prefixes (e.g. **`PHYSICS2D.START`**, **`BODY2D.CREATE`**, deprecated **`BODY2D.MAKE`**). A full teaching-oriented map (Designed → Implementation → memory) is in [moonbasic-command-set/physics-2d.md](moonbasic-command-set/physics-2d.md).

## Core Workflow

1.  **Initialize**: Start the physics world with **`PHYSICS2D.START`** (optional gravity args).
2.  **Create Bodies**: Define and create physics bodies (**`BODY2D.CREATE`**, **`BODY2D.ADDRECT`**, **`BODY2D.COMMIT`**; deprecated **`BODY2D.MAKE`** alias).
3.  **Update**: Advance the simulation each frame with **`PHYSICS2D.STEP`**.
4.  **Synchronize**: Use **`BODY2D.X`** and **`BODY2D.Y`** to read positions for your visual shapes.
5.  **Cleanup**: Shut down the world with **`PHYSICS2D.STOP`**.

---

## World Management

### `PHYSICS2D.START([gx, gy])`
Initializes the 2D physics world. Default gravity is `(0, 500)` if omitted.

### `PHYSICS2D.STOP()`
Shuts down the 2D physics simulation and frees internal buffers.

---

### `PHYSICS2D.STEP()`
Advances the 2D simulation (call once per frame).

### `PHYSICS2D.SETGRAVITY(x, y)`
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
WINDOW.OPEN(800, 600, "2D Physics Example")
WINDOW.SETFPS(60)

PHYSICS2D.START()
PHYSICS2D.SETGRAVITY(0, 500) ; Positive Y is down in 2D

floorDef = BODY2D.CREATE("STATIC")
BODY2D.ADDRECT(floorDef, 800, 50)
floorBody = BODY2D.COMMIT(floorDef, 400, 575)

boxDef = BODY2D.CREATE("DYNAMIC")
BODY2D.ADDRECT(boxDef, 40, 40)
boxBody = BODY2D.COMMIT(boxDef, 400, 100)

WHILE NOT WINDOW.SHOULDCLOSE()
    PHYSICS2D.STEP()

    RENDER.CLEAR(10, 20, 30)
    CAMERA2D.BEGIN()
        box_x = BODY2D.X(boxBody)
        box_y = BODY2D.Y(boxBody)

        DRAW.RECTANGLE(INT(BODY2D.X(floorBody)) - 400, INT(BODY2D.Y(floorBody)) - 25, 800, 50, 100, 100, 100, 255)
        DRAW.RECTANGLE(INT(box_x) - 20, INT(box_y) - 20, 40, 40, 200, 50, 50, 255)

    CAMERA2D.END()
    RENDER.FRAME()
WEND

PHYSICS2D.STOP()
WINDOW.CLOSE()
```
