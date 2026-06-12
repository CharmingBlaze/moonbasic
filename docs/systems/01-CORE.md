# Core systems: APP, RENDER, SCENE, ENTITY

> Window, frame loop, scene graph, and the entity object model — the foundation of every moonBASIC game.

**Why this page first:** Nothing else runs until you **open a window**, **clear** each frame, **draw** entities, and **present** with `RENDER.FRAME`. Entities are how the engine tracks “things in the world.” See [00-START.md](00-START.md) for a line-by-line foundation example.

**All commands in this chapter:** [COMMAND_REGISTRY.md#core-window-time](COMMAND_REGISTRY.md#core-window-time) (full arity list).

**Deep guides:** [guides/GAME-LOOP-AND-RENDERING.md](guides/GAME-LOOP-AND-RENDERING.md) (APP, RENDER, SCENE) · [guides/ENTITY-SYSTEM.md](guides/ENTITY-SYSTEM.md)

**See also:** [00-START](00-START.md) · [02-CAMERA-LIGHT](02-CAMERA-LIGHT.md) · [PROGRAMMING.md](../PROGRAMMING.md) · [reference/WINDOW.md](../reference/WINDOW.md)

---

## Table of contents

- [APP system](#app-system)
- [RENDER system](#render-system)
- [SCENE system](#scene-system)
- [ENTITY system](#entity-system)
- [Full example](#full-example)
- [See also](#see-also)

---

## APP system

Handles the application window, timing, and main-loop state. Checklist name **`APP.*`**; canonical internals use **`WINDOW.*`** and **`TIME.*`**.

### `APP.OPEN(width, height, title)`

Opens the game window and initializes the graphics runtime.

| Argument | Type | Description |
|----------|------|-------------|
| width | int | Window width in pixels |
| height | int | Window height in pixels |
| title | string | Window title shown in the OS shell |

**Returns:** nothing

**Aliases:** `WINDOW.OPEN`

**Example:**

```basic
APP.OPEN(1280, 720, "My Game")
```

---

### `APP.CLOSE()`

Closes the window and releases display resources.

**Returns:** nothing

**Aliases:** `WINDOW.CLOSE`

**Example:**

```basic
APP.CLOSE()
```

---

### `APP.SHOULDCLOSE()`

Returns whether the user requested exit (close button or platform quit).

**Returns:** `bool` — true when the loop should end

**Aliases:** `WINDOW.SHOULDCLOSE`

**Example:**

```basic
WHILE NOT APP.SHOULDCLOSE()
    ; game loop
WEND
```

---

### `APP.SETFPS(fps)`

Sets the target frames-per-second cap.

| Argument | Type | Description |
|----------|------|-------------|
| fps | int | Target FPS (e.g. 60) |

**Returns:** nothing

**Aliases:** `WINDOW.SETFPS`

**Example:**

```basic
APP.SETFPS(60)
```

---

### `APP.GETFPS()`

Returns the current measured frames per second.

**Returns:** `float`

**Aliases:** `TIME.GETFPS`, `WINDOW.GETFPS`

**Example:**

```basic
fps = APP.GETFPS()
```

---

### `APP.WIDTH()` / `APP.HEIGHT()`

Returns the window client size in pixels.

**Returns:** `int`

**Aliases:** `WINDOW.WIDTH`, `WINDOW.HEIGHT`

**Example:**

```basic
w = APP.WIDTH()
h = APP.HEIGHT()
```

---

### `APP.TIME()` / `APP.DELTA()`

`APP.TIME()` — elapsed seconds since the program started. `APP.DELTA()` — seconds since the last frame (use for movement).

**Returns:** `float`

**Aliases:** `TIME.GET`, `TIME.DELTA`

**Example:**

```basic
speed = 200
cube.move(speed * APP.DELTA(), 0, 0)
```

---

### `APP.VERSION()`

Returns the moonBASIC release version string.

**Returns:** `string`

**Aliases:** `SYSTEM.VERSION`

**Example:**

```basic
PRINT APP.VERSION()
```

---

## RENDER system

Frame clearing, 3D pass bracketing, and presentation.

### `RENDER.CLEAR(r, g, b [, a])`

Clears the framebuffer (and depth) to a color. Call once per frame before drawing.

| Argument | Type | Description |
|----------|------|-------------|
| r, g, b | int | RGB 0–255 |
| a | int | Optional alpha 0–255 |

**Returns:** nothing

**Example:**

```basic
RENDER.CLEAR(20, 24, 32)
```

---

### `RENDER.BEGIN([camera])`

Starts the 3D rendering pass. With no arguments, uses the active camera from `CAMERA.SETACTIVE`. With a camera handle, uses that camera.

**Returns:** nothing

**Aliases:** `RENDER.BEGIN3D`

**Example:**

```basic
RENDER.BEGIN(cam)
; draw 3D content
RENDER.END()
```

---

### `RENDER.END()`

Ends the current 3D pass.

**Returns:** nothing

**Aliases:** `RENDER.END3D`

**Example:**

```basic
RENDER.END()
```

---

### `RENDER.FRAME()`

Presents the frame to the screen. Call after all drawing for the frame.

**Returns:** nothing

**Example:**

```basic
RENDER.FRAME()
```

---

### `RENDER.SETBACKGROUND(r, g, b)`

Sets the clear color for the next frame (alias of clear-color path).

| Argument | Type | Description |
|----------|------|-------------|
| r, g, b | int | RGB 0–255 |

**Returns:** nothing

**Example:**

```basic
RENDER.SETBACKGROUND(18, 20, 28)
```

---

### `RENDER.SETWIREFRAME(enabled)`

Toggles wireframe 3D drawing for debugging.

| Argument | Type | Description |
|----------|------|-------------|
| enabled | bool | true = wireframe |

**Returns:** nothing

**Example:**

```basic
RENDER.SETWIREFRAME(true)
```

---

### `RENDER.SCREENSHOT(path)`

Saves the current framebuffer to an image file (PNG).

| Argument | Type | Description |
|----------|------|-------------|
| path | string | Output file path |

**Returns:** nothing

**Example:**

```basic
RENDER.SCREENSHOT("screenshots/level1.png")
```

---

## SCENE system

Groups entities for load/save and batched drawing.

### `SCENE.REGISTER(name)`

Creates or registers a named scene handle.

| Argument | Type | Description |
|----------|------|-------------|
| name | string | Scene name |

**Returns:** `handle` — scene handle

**Note:** Checklist `SCENE.CREATE` maps to register/switch workflow — see [reference/SCENE.md](../reference/SCENE.md).

**Example:**

```basic
level = SCENE.REGISTER("Level1")
```

---

### `SCENE.SWITCH(scene)`

Activates a scene for updates and drawing.

| Argument | Type | Description |
|----------|------|-------------|
| scene | handle | Scene handle |

**Returns:** nothing

**Aliases:** checklist `SCENE.SETACTIVE`

**Example:**

```basic
SCENE.SWITCH(level)
```

---

### `SCENE.DRAW()`

Draws the active scene’s entities. Call inside `RENDER.BEGIN` / `RENDER.END`.

**Returns:** nothing

**Example:**

```basic
RENDER.BEGIN(cam)
SCENE.DRAW()
RENDER.END()
```

---

### `SCENE.SAVESCENE(scene, path)` / `SCENE.LOADSCENE(path)`

Save or load scene data to disk.

| Argument | Type | Description |
|----------|------|-------------|
| scene | handle | Scene (save only) |
| path | string | File path |

**Returns:** nothing

**Aliases:** checklist `SCENE.SAVE`, `SCENE.LOAD`

**Example:**

```basic
SCENE.SAVESCENE(level, "levels/level1.scene")
SCENE.LOADSCENE("levels/level1.scene")
```

---

### `SCENE.CLEARSCENE()`

Removes entities from the active scene without closing the window.

**Returns:** nothing

**Aliases:** checklist `SCENE.CLEAR`

**Example:**

```basic
SCENE.CLEARSCENE()
```

---

## ENTITY system

The main object model — transforms, hierarchy, visibility, tags, and drawing.

### `ENTITY.CREATE(name)`

Creates an empty entity.

| Argument | Type | Description |
|----------|------|-------------|
| name | string | Debug / scene name |

**Returns:** `handle`

**Example:**

```basic
player = ENTITY.CREATE("Player")
```

---

### `ENTITY.CREATECUBE(w, h, d)`

Creates a box entity with default mesh.

| Argument | Type | Description |
|----------|------|-------------|
| w, h, d | float | Box dimensions |

**Returns:** `handle`

**Example:**

```basic
crate = ENTITY.CREATECUBE(2, 2, 2)
```

---

### `ENTITY.CREATESPHERE(radius)` / `ENTITY.CREATEPIVOT(name)`

Sphere primitive or empty pivot for hierarchy.

**Returns:** `handle`

**Example:**

```basic
ball = ENTITY.CREATESPHERE(1)
pivot = ENTITY.CREATEPIVOT("Pivot")
```

---

### `ENTITY.FREE(entity)`

Destroys an entity and frees its resources.

| Argument | Type | Description |
|----------|------|-------------|
| entity | handle | Entity to destroy |

**Returns:** nothing

**Aliases:** checklist `ENTITY.DESTROY`

**Example:**

```basic
ENTITY.FREE(crate)
```

---

### `ENTITY.COPY(entity)`

Clones an entity.

**Returns:** `handle` — new entity

**Aliases:** checklist `ENTITY.CLONE`

**Example:**

```basic
copy = ENTITY.COPY(player)
```

---

### Transform commands

| Command | Description |
|---------|-------------|
| `ENTITY.SETPOS(ent, x, y, z)` | Set world position |
| `ENTITY.SETROT(ent, p, y, r)` | Set rotation (degrees) |
| `ENTITY.SETSCALE(ent, x, y, z)` | Set scale |
| `ENTITY.MOVE(ent, dx, dy, dz)` | Translate by delta |
| `ENTITY.TURN(ent, dp, dy, dr)` | Rotate by delta |
| `ENTITY.X/Y/Z(ent)` | Read position components |
| `ENTITY.LOOKAT(ent, target)` | Face another entity |
| `ENTITY.POINTAT(ent, x, y, z)` | Face a world point |

**Handle shortcuts:** `ent.pos(x,y,z)`, `ent.turn(...)`, `ent.move(...)`

**Example:**

```basic
ENTITY.SETPOS(player, 0, 1, 5)
player.turn(0, 45 * APP.DELTA(), 0)
```

---

### Hierarchy and state

| Command | Description |
|---------|-------------|
| `ENTITY.PARENT(child, parent)` | Set parent (checklist `SETPARENT`) |
| `ENTITY.PARENTCLEAR(child)` | Clear parent (`CLEARPARENT`) |
| `ENTITY.PARENTGET(child)` | Get parent handle |
| `ENTITY.CHILDCOUNT(parent)` | Number of children |
| `ENTITY.GETCHILD(parent, index)` | Child by index |
| `ENTITY.SHOW/HIDE(ent)` | Visibility |
| `ENTITY.SETVISIBLE(ent, bool)` | Visibility flag |
| `ENTITY.VISIBLE(ent)` | Read visibility (`ISVISIBLE`) |
| `ENTITY.SETNAME/GETNAME` | Name string |
| `ENTITY.SETTAG/HASTAG` | Tag strings for queries |

**Example:**

```basic
ENTITY.PARENT(sword, player)
ENTITY.SETTAG(enemy, "hostile")
```

---

## Full example

```basic
; Foundation loop — see examples/foundation/main.mb
APP.OPEN(1280, 720, "Core Demo")
APP.SETFPS(60)

cam = CAMERA.CREATE()
CAMERA.SETACTIVE(cam)
CAMERA.SETPOS(cam, 0, 2, -8)
CAMERA.LOOKAT(cam, 0, 0, 0)

cube = ENTITY.CREATECUBE(2, 2, 2)
cube.pos(0, 0, 5)

WHILE NOT APP.SHOULDCLOSE()
    cube.turn(0, 60 * APP.DELTA(), 0)
    RENDER.CLEAR(20, 20, 30)
    RENDER.BEGIN(cam)
    SCENE.DRAW()
    RENDER.END()
    RENDER.FRAME()
WEND

APP.CLOSE()
```

---

## Memory notes

- Call `ENTITY.FREE` (or `ent.free()`) when removing entities permanently.
- `APP.CLOSE` / `WINDOW.CLOSE` after the main loop ends.
- Scene clear does not free GPU assets from textures/models — use `TEXTURE.FREE`, `MODEL.FREE` separately.

---

## See also

- [02-CAMERA-LIGHT](02-CAMERA-LIGHT.md)
- [reference/ENTITY.md](../reference/ENTITY.md)
- [reference/RENDER.md](../reference/RENDER.md)
- [reference/WINDOW.md](../reference/WINDOW.md)
