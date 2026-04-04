# Render Commands

Commands for drawing, clearing the screen, and managing the render state.

## The Rendering Loop

All drawing in moonBASIC happens inside the main game loop. Each frame, you must:

1.  **Clear the screen**: Call `Render.Clear()` to wipe the previous frame.
2.  **Enter a drawing mode**: Use `Render.BeginMode2D()` or `cam.Begin()` for 3D.
3.  **Call drawing commands**: Use `Draw.*`, `Mesh.Draw`, etc.
4.  **End the mode**: Use `Render.EndMode2D()` or `cam.End()`.
5.  **Present the frame**: Call `Render.Frame()` to show the result on screen.

```basic
WHILE NOT Window.ShouldClose()
    ; 1. Clear
    Render.Clear(10, 20, 30)

    ; 2. Begin Mode
    Render.BeginMode2D()
        ; 3. Draw
        Draw.Rect(100, 100, 50, 50, 255, 255, 255, 255)
    ; 4. End Mode
    Render.EndMode2D()

    ; 5. Present
    Render.Frame()
WEND
```

---

## Core Commands

### `Render.Clear(r, g, b, [a])`

Clears the entire screen to a specific color. The alpha component `a` is optional.

- `r`, `g`, `b`, `a`: The red, green, blue, and alpha components of the color (0-255).

---

### `Render.Frame()`

Swaps the back buffer to the front, displaying everything drawn since the last `Render.Clear()` call. This should be called once at the very end of the main loop.

---

## Drawing Modes

### `Render.BeginMode2D()` / `Render.EndMode2D()`

Enters and exits 2D drawing mode. The coordinate system starts at `(0,0)` in the top-left corner. All `Draw.*` commands for 2D shapes, text, and textures must be within this block.

---

### `cam.Begin()` / `cam.End()`

Enters and exits 3D drawing mode using a specific camera's perspective. All 3D drawing commands like `Mesh.Draw` and `Model.Draw` must be within this block. See the [Camera Reference](CAMERA.md) for details.

---

### `Render.BeginShader(shaderHandle)` / `Render.EndShader()`

Applies a custom shader to all subsequent drawing commands until `Render.EndShader()` is called. See the [Shader Reference](SHADER.md) for details.

---

## Utilities

### `Render.DrawFPS(x, y)`

Draws the current frames per second at the specified screen coordinates.

- `x`, `y`: The position to draw the text.

---

### `Render.Screenshot(filePath$)`

Saves a screenshot of the current frame to a `.png` file.

- `filePath$`: The path to save the file.

---

## State Management

### `Render.SetBlend(mode)`

Sets the blending mode for drawing. This is useful for effects like transparency or additive lighting.

- `mode`: An integer representing the blend mode. Common values include:
    - `BLEND_ALPHA`: Standard transparency.
    - `BLEND_ADDITIVE`: Brightens colors by adding them together.
