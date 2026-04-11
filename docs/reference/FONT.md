# Font Commands

Commands for loading and drawing with custom fonts.

## Core Workflow

1.  **Load Font**: Use `Font.Load()` to load a `.ttf` or `.otf` file. Store the returned handle.
2.  **Draw Text**: In the main loop, use `Draw.TextFont()` with the font handle to draw text.
3.  **Free Font**: When you are done, call `Font.Free()` to unload the font from memory.

---

### `Font.Load(path)`
Loads a `.ttf` or `.otf` font file from disk. Returns a **font handle**.

### `Font.Free(handle)`
Unloads a font from memory and releases its heap slot.

---

### `Draw.TextFont(handle, text, x, y, size, spacing, r, g, b, a)`
Draws text using a specific font handle. This must be called within a **`Camera2D.Begin()`** / **`Camera2D.End()`** block.
- `handle`: The handle of the loaded font.
- `text`: The string to draw.
- `x`, `y`: Screen position.
- `size`: Font size in pixels.
- `spacing`: Extra spacing between characters.
- `r, g, b, a`: Color components (0-255).

---

## Full Example

This example assumes you have a font file named `my_font.ttf` in the same directory as your script.

```basic
Window.Open(800, 600, "Font Example")
Window.SetFPS(60)

; 1. Load the font
my_font = Font.Load("my_font.ttf")

IF my_font = 0 THEN
    PRINT "Error: Could not load my_font.ttf"
    Window.Close()
    SYSTEM.EXIT()
ENDIF

WHILE NOT Window.ShouldClose()
    Render.Clear(50, 60, 70)

    Camera2D.Begin()
        ; 2. Draw text using the loaded font
        Draw.TextFont(my_font, "Hello, moonBASIC!", 100, 200, 48, 2, 255, 200, 100, 255)

        ; You can still use the default font for other text
        Draw.Text("This is the default system font.", 100, 300, 20, 200, 200, 200, 255)
    Camera2D.End()

    Render.Frame()
WEND

; 3. Free the font
Font.Free(my_font)
Window.Close()
```
