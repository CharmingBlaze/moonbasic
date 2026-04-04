# 2D Drawing Commands

Commands for drawing basic shapes, text, and textures in 2D mode.

**Note**: Basic 2D drawing can be done directly between `Render.Clear()` and `Render.Frame()`. Use `Camera2D.Begin()` / `Camera2D.End()` only when you need 2D camera transforms (scrolling, zooming, rotation).

---

## Shapes

### `Draw.Rectangle(x, y, width, height, r, g, b, a)`

Draws a filled rectangle. The `(x, y)` coordinates specify the top-left corner.

- `r`, `g`, `b`, `a`: The color and alpha (transparency) of the rectangle (0-255).

### `Draw.RectangleRounded(x, y, width, height, radius, r, g, b, a)`

Draws a filled rectangle with rounded corners.

- `radius`: The corner radius in pixels.

### `Draw.Circle(centerX, centerY, radius, r, g, b, a)`

Draws a filled circle, centered at `(centerX, centerY)`.

### `Draw.Line(startX, startY, endX, endY, r, g, b, a)`

Draws a line between two points.

---

## Text

### `Draw.Text(text$, x, y, size, r, g, b, a)`

Draws a line of text using the built-in default font.

- `text$`: The string to draw.
- `x`, `y`: The top-left position to start drawing the text.
- `size`: The font size in pixels.

---

## Textures

### `Draw.Texture(textureHandle, x, y, [r, g, b, a])`

Draws a texture at a specific position. The color tint is optional.

- `textureHandle`: The handle of the texture to draw, returned by `Texture.Load()`.
- `x`, `y`: The top-left position to draw the texture.
- `r, g, b, a`: An optional color to tint the texture.

---

### `Draw.TextureNPatch(textureHandle, left, top, right, bottom, x, y, width, height, r, g, b, a)`

Draws a texture as a 9-patch, which allows it to be stretched without distorting the corners. This is ideal for UI elements like buttons and panels.

- `textureHandle`: The handle of the 9-patch texture.
- `left, top, right, bottom`: The size in pixels of the un-stretched borders of the source image.
- `x, y, width, height`: The destination rectangle on the screen.
- `r, g, b, a`: An optional color tint.


---

## Full Example

This example demonstrates how to use several drawing commands together within the main loop.

```basic
Window.Open(800, 600, "2D Drawing Example")
Window.SetFPS(60)

; Load a texture once at the start
player_tex = Texture.Load("player.png")

player_x# = 375
player_y# = 275

WHILE NOT Window.ShouldClose()
    dt# = Time.Delta()

    ; --- UPDATE LOGIC ---
    IF Input.KeyDown(KEY_RIGHT) THEN player_x# = player_x# + 200 * dt#
    IF Input.KeyDown(KEY_LEFT) THEN player_x# = player_x# - 200 * dt#

    ; --- DRAWING ---
    Render.Clear(14, 22, 33)

    ; Draw a background rectangle
    Draw.Rectangle(0, 500, 800, 100, 40, 50, 60, 255)

    ; Draw the player texture
    Draw.Texture(player_tex, INT(player_x#), INT(player_y#), 255, 255, 255, 255)

    ; Draw a health bar
    Draw.Rectangle(10, 10, 200, 20, 200, 0, 0, 255)

    ; Draw text on top of the health bar
    Draw.Text("HEALTH", 15, 12, 16, 255, 255, 255, 255)

    Render.Frame()
WEND

; --- CLEANUP ---
Texture.Free(player_tex)
Window.Close()
```
