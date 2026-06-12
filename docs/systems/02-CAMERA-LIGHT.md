# Camera and light systems

> View projection and scene lighting — bracket your 3D draw pass with a camera and illuminate entities with point, directional, and spot lights.

**All commands:** [COMMAND_REGISTRY.md#camera-light](COMMAND_REGISTRY.md#camera-light)

**Deep guides:** [guides/CAMERA-AND-INPUT.md](guides/CAMERA-AND-INPUT.md) · [guides/LIGHTING.md](guides/LIGHTING.md)

**See also:** [01-CORE](01-CORE.md) · [03-ASSETS](03-ASSETS.md) · [reference/CAMERA.md](../reference/CAMERA.md) · [reference/LIGHT.md](../reference/LIGHT.md)

**Case:** Command names are **case-insensitive** (`camera.create` = `CAMERA.CREATE`).

---

## Table of contents

- [CAMERA system](#camera-system)
- [LIGHT system](#light-system)
- [Full example](#full-example)
- [Memory notes](#memory-notes)
- [See also](#see-also)

---

## CAMERA system

3D cameras control where the player looks and how world space maps to the screen.

### Core workflow

1. `CAMERA.CREATE()` — new camera handle.
2. `CAMERA.SETACTIVE(cam)` — default for `RENDER.BEGIN()` with no argument.
3. `CAMERA.SETPOS` / `SETTARGET` / `LOOKAT` — aim the view.
4. Each frame: `RENDER.BEGIN(cam)` … draw … `RENDER.END()`.
5. `CAMERA.FREE(cam)` when done.

---

### `CAMERA.CREATE()`

Creates a perspective 3D camera.

**Returns:** `handle`

**Example:**

```basic
cam = CAMERA.CREATE()
```

---

### `CAMERA.SETACTIVE(camera)`

Sets the camera used when `RENDER.BEGIN()` is called without arguments.

| Argument | Type | Description |
|----------|------|-------------|
| camera | handle | Camera handle |

**Returns:** nothing

**Example:**

```basic
CAMERA.SETACTIVE(cam)
```

---

### `CAMERA.SETPOS(camera, x, y, z)`

Places the camera eye in world space.

| Argument | Type | Description |
|----------|------|-------------|
| camera | handle | Camera |
| x, y, z | float | World position |

**Returns:** `handle` — camera (for chaining)

**Aliases:** `CAMERA.SETPOSITION`, handle `.pos(x,y,z)`

**Example:**

```basic
CAMERA.SETPOS(cam, 0, 2, -8)
```

---

### `CAMERA.SETTARGET(camera, x, y, z)` / `CAMERA.LOOKAT(camera, x, y, z)`

Points the camera at a world position.

| Argument | Type | Description |
|----------|------|-------------|
| camera | handle | Camera |
| x, y, z | float | Look-at point |

**Returns:** `handle`

**Example:**

```basic
CAMERA.LOOKAT(cam, 0, 0, 0)
```

---

### `CAMERA.SETFOV(camera, degrees)`

Sets horizontal field of view.

| Argument | Type | Description |
|----------|------|-------------|
| camera | handle | Camera |
| degrees | float | FOV (e.g. 70) |

**Returns:** `handle`

**Example:**

```basic
CAMERA.SETFOV(cam, 70)
```

---

### `CAMERA.FOLLOW(camera, entity, offsetX, offsetY, offsetZ)`

Third-person style follow helper (orbit / lag depending on setup).

| Argument | Type | Description |
|----------|------|-------------|
| camera | handle | Camera |
| entity | handle | Entity to follow |
| offsetX, offsetY, offsetZ | float | Offset from target |

**Returns:** `handle`

**Aliases:** `CAMERA.FOLLOWENTITY`

**Example:**

```basic
CAMERA.FOLLOW(cam, player, 0, 3, -8)
```

---

### `CAMERA.BEGIN(camera)` / `CAMERA.END()`

Starts or ends a camera-specific 3D pass. Prefer **`RENDER.BEGIN(cam)`** / **`RENDER.END()`** in game loops.

**Returns:** nothing

**Example:**

```basic
RENDER.BEGIN(cam)
SCENE.DRAW()
RENDER.END()
```

---

### `CAMERA.GETRAY(camera, screenX, screenY)`

Builds a world ray from screen pixels (picking helper).

**Returns:** ray data — see [PICK system](05-PHYSICS.md#pick-system)

**Aliases:** checklist `CAMERA.SCREENRAY`

**Example:**

```basic
; Use PICK.SCREENCAST(cam) for full hit workflow
```

---

### `CAMERA.WORLDTOSCREEN(camera, x, y, z)` / `CAMERA.SCREENTOWORLD(camera, x, y)`

Convert between world and screen coordinates.

**Returns:** position components or handle array — see [reference/CAMERA.md](../reference/CAMERA.md)

**Example:**

```basic
; HUD markers: project entity world pos to screen
```

---

### `CAMERA.FREE(camera)`

Releases the camera handle.

| Argument | Type | Description |
|----------|------|-------------|
| camera | handle | Camera to free |

**Returns:** nothing

**Example:**

```basic
CAMERA.FREE(cam)
```

---

## LIGHT system

Lights are handles attached to the lighting pass. Use point, directional, and spot types.

### Core workflow

1. `LIGHT.CREATEPOINT` / `CREATEDIRECTIONAL` / `CREATESPOT` (or `LIGHT.CREATE("point")`).
2. `LIGHT.SETPOS` or `SETDIR` — place or aim the light.
3. `LIGHT.SETCOLOR`, `SETINTENSITY`, `SETRANGE` as needed.
4. `LIGHT.SETSHADOW(light, true)` for shadow-casting sun lights.
5. `LIGHT.FREE(light)` when removing.

---

### `LIGHT.CREATEPOINT()` / `LIGHT.CREATEDIRECTIONAL()` / `LIGHT.CREATESPOT()`

Creates a typed light handle.

**Returns:** `handle`

**Aliases:** `LIGHT.CREATE("point"|"directional"|"spot")`

**Example:**

```basic
lamp = LIGHT.CREATEPOINT()
sun = LIGHT.CREATEDIRECTIONAL()
```

---

### `LIGHT.SETPOS(light, x, y, z)` / `LIGHT.SETDIR(light, x, y, z)`

Position a point/spot light or set a directional light’s direction vector.

| Argument | Type | Description |
|----------|------|-------------|
| light | handle | Light |
| x, y, z | float | Position or direction |

**Returns:** `handle`

**Example:**

```basic
LIGHT.SETPOS(lamp, 0, 5, -3)
LIGHT.SETDIR(sun, -1, -2, -1)
```

---

### `LIGHT.SETCOLOR(light, r, g, b [, a])`

Sets RGB color (0–255) and optional intensity multiplier.

| Argument | Type | Description |
|----------|------|-------------|
| light | handle | Light |
| r, g, b | int | Color components |
| a | float | Optional strength multiplier |

**Returns:** `handle`

**Example:**

```basic
LIGHT.SETCOLOR(lamp, 255, 220, 180)
```

---

### `LIGHT.SETINTENSITY(light, value)`

Sets brightness scalar.

| Argument | Type | Description |
|----------|------|-------------|
| light | handle | Light |
| value | float | Intensity |

**Returns:** `handle`

**Example:**

```basic
LIGHT.SETINTENSITY(lamp, 2.0)
```

---

### `LIGHT.SETRANGE(light, distance)`

Maximum reach for point and spot lights.

| Argument | Type | Description |
|----------|------|-------------|
| light | handle | Light |
| distance | float | Range in world units |

**Returns:** `handle`

**Example:**

```basic
LIGHT.SETRANGE(lamp, 15)
```

---

### `LIGHT.SETSHADOW(light, enabled)`

Enables or disables shadow casting for this light.

| Argument | Type | Description |
|----------|------|-------------|
| light | handle | Light |
| enabled | bool | true = cast shadows |

**Returns:** `handle`

**Aliases:** checklist `SETSHADOWS`

**Example:**

```basic
LIGHT.SETSHADOW(sun, true)
```

---

### `LIGHT.FREE(light)`

Destroys the light handle.

**Returns:** nothing

**Example:**

```basic
LIGHT.FREE(lamp)
```

---

## Full example

```basic
APP.OPEN(1280, 720, "Camera + Light")
APP.SETFPS(60)

cam = CAMERA.CREATE()
CAMERA.SETACTIVE(cam)
CAMERA.SETPOS(cam, 0, 2, -8)
CAMERA.LOOKAT(cam, 0, 0, 0)

cube = ENTITY.CREATECUBE(2, 2, 2)
cube.pos(0, 0, 5)

lamp = LIGHT.CREATEPOINT()
LIGHT.SETPOS(lamp, 0, 5, -3)
LIGHT.SETCOLOR(lamp, 255, 240, 200)
LIGHT.SETINTENSITY(lamp, 1.5)

WHILE NOT APP.SHOULDCLOSE()
    cube.turn(0, 45 * APP.DELTA(), 0)
    RENDER.CLEAR(25, 28, 35)
    RENDER.BEGIN(cam)
    SCENE.DRAW()
    RENDER.END()
    RENDER.FRAME()
WEND

CAMERA.FREE(cam)
LIGHT.FREE(lamp)
ENTITY.FREE(cube)
APP.CLOSE()
```

Check: `moonbasic --check` on your script. Run with **`moonrun`**.

---

## Memory notes

- Call **`CAMERA.FREE`** and **`LIGHT.FREE`** when discarding handles.
- Lights do not automatically follow entity transforms unless you update `SETPOS` each frame or parent a light entity.
- Shadow quality ties to `RENDER.SETSHADOWMAPSIZE` — see [reference/LIGHT.md](../reference/LIGHT.md).

---

## See also

- [01-CORE](01-CORE.md) — `RENDER.BEGIN` / `SCENE.DRAW`
- [05-PHYSICS](05-PHYSICS.md) — `PICK.SCREENCAST`
- [reference/CAMERA_LIGHT_RENDER.md](../reference/CAMERA_LIGHT_RENDER.md) — combined render pass map
