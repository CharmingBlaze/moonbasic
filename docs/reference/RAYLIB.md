# Raylib Commands

Direct Raylib C API bindings. Use these when you need raw Raylib control that moonBASIC's higher-level wrappers don't expose. For most tasks, prefer the `WINDOW.*`, `RENDER.*`, `DRAW.*`, `INPUT.*`, `CAMERA.*`, `TEXTURE.*`, `SHADER.*` namespaces.

## Core Workflow

```basic
RAYLIB.INITWINDOW(960, 540, "Title")
RAYLIB.SETTARGETFPS(60)
WHILE NOT RAYLIB.WINDOWSHOULDCLOSE()
    RAYLIB.BEGINFRAME()
    RAYLIB.CLEARBACKGROUND(20, 25, 35)
    RAYLIB.DRAWRECTANGLE(10, 10, 100, 40, 80, 160, 255, 255)
    RAYLIB.ENDFRAME()
WEND
RAYLIB.CLOSEWINDOW()
```

---

## Window

### `RAYLIB.INITWINDOW(width, height, title)` 

Opens the Raylib window. Direct alias of `WINDOW.OPEN`.

---

### `RAYLIB.CLOSEWINDOW()` 

Closes the window. Alias of `WINDOW.CLOSE`.

---

### `RAYLIB.WINDOWSHOULDCLOSE()` 

Returns `TRUE` when the window close button is pressed.

---

### `RAYLIB.GETFRAMEBUFFERWIDTH()` / `RAYLIB.GETFRAMEBUFFERHEIGHT()` 

Returns the framebuffer dimensions.

---

## Frame

### `RAYLIB.BEGINFRAME()` / `RAYLIB.ENDFRAME()` 

Begin and end a render frame. Alias of `RENDER.CLEAR` + `RENDER.FRAME` pattern.

---

### `RAYLIB.CLEARBACKGROUND(r, g, b)` 

Clears the background color.

---

### `RAYLIB.DRAWFPS(x, y)` 

Draws the current FPS counter at screen position.

---

### `RAYLIB.SETTARGETFPS(fps)` / `RAYLIB.GETFPS()` / `RAYLIB.GETTIME()` 

FPS and time utilities.

---

## Draw

### `RAYLIB.DRAWRECTANGLE(x, y, w, h, r, g, b, a)` 

Draws a filled rectangle.

---

### `RAYLIB.DRAWCIRCLE(x, y, radius, r, g, b, a)` 

Draws a filled circle.

---

### `RAYLIB.DRAWTEXTURE(tex, x, y, w, h, r, g, b, a)` 

Draws a texture.

---

### `RAYLIB.DRAWLINE3D(x1, y1, z1, x2, y2, z2, r, g, b, a)` 

3D line.

---

### `RAYLIB.DRAWCUBE(x, y, z, w, h, d, r, g, b, a)` 

3D filled cube.

---

### `RAYLIB.DRAWSPHERE(x, y, z, radius, r, g, b, a)` 

3D sphere.

---

## Texture & Model

### `RAYLIB.LOADTEXTURE(path)` 

Loads a texture. Returns a texture handle. Alias of `TEXTURE.LOAD`.

---

### `RAYLIB.UNLOADTEXTURE(tex)` 

Unloads a texture handle.

---

### `RAYLIB.LOADMODEL(path)` 

Loads a 3D model. Returns a model handle.

---

### `RAYLIB.DRAWMODEL(model, x, y, z, sx, sy, sz, scale)` 

Draws a model at world position with scale.

---

## Shader

### `RAYLIB.LOADSHADER(vsPath, fsPath)` 

Loads a vertex + fragment shader. Returns a shader handle.

---

### `RAYLIB.BEGINSHADERMODE(shader)` / `RAYLIB.ENDSHADERMODE()` 

Wraps draw calls in a custom shader.

---

## Camera

### `RAYLIB.SETCAMERAMODE(cam, mode)` 

Sets a Raylib camera update mode.

---

### `RAYLIB.UPDATECAMERA(cam, mode)` 

Updates a Raylib camera each frame.

---

## Input

### `RAYLIB.ISKEYDOWN(key)` / `RAYLIB.ISKEYPRESSED(key)` / `RAYLIB.ISKEYRELEASED(key)` 

Raw Raylib key queries. Prefer `INPUT.KEYDOWN` / `INPUT.KEYPRESSED`.

---

### `RAYLIB.GETMOUSEX()` / `RAYLIB.GETMOUSEY()` / `RAYLIB.ISMOUSEBUTTONDOWN(btn)` 

Raw mouse queries.

---

## Full Example

Raw Raylib window with a spinning cube.

```basic
RAYLIB.INITWINDOW(960, 540, "Raylib Raw Demo")
RAYLIB.SETTARGETFPS(60)

t = 0.0
WHILE NOT RAYLIB.WINDOWSHOULDCLOSE()
    t = t + 1.0 / RAYLIB.GETFPS()

    RAYLIB.BEGINFRAME()
    RAYLIB.CLEARBACKGROUND(20, 25, 35)
    RAYLIB.DRAWCUBE(0, 0, 0, 2, 2, 2, 80, 160, 255, 255)
    RAYLIB.DRAWFPS(10, 10)
    RAYLIB.ENDFRAME()
WEND

RAYLIB.CLOSEWINDOW()
```

---

## See also

- [WINDOW.md](WINDOW.md) — window management
- [RENDER.md](RENDER.md) — frame loop
- [DRAW.md](DRAW.md) — 2D/3D drawing
- [INPUT.md](INPUT.md) — input system
- [SHADER.md](SHADER.md) — shader management
