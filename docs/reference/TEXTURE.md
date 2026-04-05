# Texture Commands

Commands for loading and managing texture resources, which are images stored in GPU memory for fast drawing. For **CPU-side** pixel editing first, see **[IMAGE.md](IMAGE.md)** and **`Texture.FromImage`** below.

## Core Workflow

1.  **Load**: Use `Texture.Load()` to load an image from a file into a texture. Store the returned handle.
2.  **Draw**: In the main loop, use `Draw.Texture()` with the handle to draw the texture.
3.  **Free**: When you are done with the texture (e.g., when a level ends or the program closes), call `Texture.Free()` to unload it from the GPU.

---

### `Texture.Load(filePath$)`

Loads an image file from disk into GPU memory as a texture. It's efficient to load textures once at the beginning of your program or level, not every frame.

- `filePath$`: The path to the image file (e.g., `.png`, `.jpg`).

Returns a handle to the texture resource.

---

### `Texture.Free(textureHandle)`

Unloads a texture from GPU memory, freeing up resources. This is a crucial step to prevent memory leaks.

- `textureHandle`: The handle of the texture to free.

---

### `Texture.FromImage(imageHandle)`

Creates a texture from an in-memory `Image` handle. This is an advanced feature for when you need to generate or modify image data in your code before drawing it.

- `imageHandle`: The handle to an image resource.

---

### `Texture.GenWhiteNoise(width, height)`

Generates and returns a handle to a procedural white noise texture. Useful for screen effects.

- `width`, `height`: The dimensions of the texture to generate.

---

## Full Example: Load and Draw a Texture

This example assumes you have an image file named `character.png` in the same directory as your script.

```basic
Window.Open(800, 600, "Texture Example")
Window.SetFPS(60)

; 1. Load the texture once
char_tex = Texture.Load("character.png")

; Check if loading failed (e.g., file not found)
IF char_tex = 0 THEN
    PRINT "Error: Could not load character.png"
    Window.Close()
    SYSTEM.EXIT()
ENDIF

WHILE NOT Window.ShouldClose()
    Render.Clear(20, 30, 40)

    Render.BeginMode2D()
        ; 2. Draw the texture every frame
        Draw.Text("My Character:", 100, 150, 20, 255, 255, 255, 255)
        Draw.Texture(char_tex, 100, 180)
    Render.EndMode2D()

    Render.Frame()
WEND

; 3. Free the texture before exiting
Texture.Free(char_tex)
Window.Close()
```

---

## Render targets and atlases

Off-screen **`RenderTexture`** usage and multi-stage post-processing are implemented inside engine modules (e.g. shadow maps, post stack). For **sprite sheets**, prefer **[ATLAS.md](ATLAS.md)** (`ATLAS.*`) so multiple logical sprites share one GPU texture.

---
