# Font Commands

Commands for loading and drawing with custom fonts.

## Core Workflow

1.  **Load Font**: Use `Font.Load()` to load a `.ttf` or `.otf` file. Store the returned handle.
2.  **Draw Text**: In the main loop, use `Draw.TextFont()` with the font handle to draw text.
3.  **Free Font**: When you are done, call `Font.Free()` to unload the font from memory.

---

### `Font.Load(filePath$)`

Loads a font file from disk. It's best to load fonts once at the start of your program.

- `filePath$`: The path to the font file (e.g., `.ttf`, `.otf`).

Returns a handle to the font resource.

---

### `Font.Free(fontHandle)`

Unloads a font from memory. This is important to prevent memory leaks.

- `fontHandle`: The handle of the font to free.

---

### `Draw.TextFont(fontHandle, text$, x, y, size, spacing, r, g, b, a)`

Draws text using a loaded font. This must be called within a **`Camera2D.Begin()`** / **`Camera2D.End()`** block (or between **`Camera2D.Begin(cam)`** / **`Camera2D.End()`** when using a 2D camera handle).

- `fontHandle`: The handle of the font to use.
- `text$`: The string to draw.
- `x`, `y`: The top-left position to start drawing.
- `size`: The font size.
- `spacing`: The spacing between characters.
- `r`, `g`, `b`, `a`: The color and alpha of the text.

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
