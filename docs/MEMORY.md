# Memory Management in moonBASIC

While moonBASIC aims to simplify game development, it's important to understand how it manages memory for complex resources like textures, models, and physics objects.

---

## The Handle System

moonBASIC uses a **handle-based system** for managing resources. When you create a resource (e.g., by loading a texture or creating a mesh), the command doesn't return the raw data. Instead, it returns a **handle**—a unique integer that acts as an ID for that resource.

```basic
; `my_texture` is not the image data itself, but a handle to it.
my_texture = Texture.Load("player.png")

; `cube` is a handle to a mesh resource managed by the engine.
cube = Mesh.MakeCube(2, 2, 2)
```

You use this handle in other commands to refer to the resource you want to work with.

```basic
; Use the handle to draw the texture
Draw.Texture(my_texture, 100, 100)
```

### Why Handles?

This system keeps the BASIC-like language simple and clean. You don't have to deal with complex data structures or pointers directly. The moonBASIC runtime engine manages the actual data in the background; you just tell it which resource to use via its handle.

---

## The Resource Lifecycle: Create, Use, Free

Every resource you create must eventually be freed. If you don't free a resource, it will remain in memory until the program ends, which can lead to high memory usage (a "memory leak").

The typical lifecycle is:

1.  **Create**: Use a `Make...` or `Load...` command to create the resource and get a handle.
2.  **Use**: Pass the handle to other commands (e.g., `Draw`, `SetPos`, `ApplyForce`).
3.  **Free**: When you are finished with the resource, call its corresponding `Free...` command.

### Example

```basic
; 1. Create a texture resource
player_tex = Texture.Load("player.png")

; 2. Use the texture in the main loop
WHILE NOT Window.ShouldClose()
    Render.BeginMode2D()
        Draw.Texture(player_tex, 50, 50)
    Render.EndMode2D()
    Render.Frame()
WEND

; 3. Free the texture resource before the program exits
Texture.Free(player_tex)

Window.Close()
```

**It is crucial to free resources you are no longer using.** A good rule of thumb is to free any resource you create. If you load a texture in a game level, free it when the level is unloaded.

Common `Free` commands include:
- `Texture.Free()`
- `Mesh.Free()`
- `Model.Free()`
- `Sound.Free()`
- `Body3D.Free()`
- `Body2D.Free()`

---

## Low-Level Memory Blocks

For advanced use cases, moonBASIC provides commands for direct memory manipulation through the `MEM` module. This is useful for working with binary data, custom file formats, or interoperating with external libraries.

These commands also use a handle-based system.

- `MEM.MAKE(size)`: Allocates a block of memory of a given size and returns a handle.
- `MEM.GETBYTE(handle, offset)`: Reads a byte from the memory block.
- `MEM.SETBYTE(handle, offset, value)`: Writes a byte to the memory block.
- `MEM.FREE(handle)`: Frees the memory block.

These are advanced features and should be used with care. Like other resources, memory blocks created with `MEM.MAKE` must be freed with `MEM.FREE`.
