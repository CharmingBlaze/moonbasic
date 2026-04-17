# Shader Commands

Load, configure, and apply custom GLSL shaders to the render pipeline.

Page shape follows [DOC_STYLE_GUIDE.md](../DOC_STYLE_GUIDE.md) (**WAVE pattern**).

## Core Workflow

1. Load a shader from vertex/fragment source files with `SHADER.LOAD`.
2. Look up uniform locations with `SHADER.GETLOC`.
3. Set uniform values each frame with `SHADER.SETFLOAT`, `SHADER.SETVEC3`, etc.
4. Free the shader when done with `SHADER.FREE`.

For compute shaders, see the `COMPUTESHADER.*` namespace.

---

### `SHADER.LOAD(vertexPath, fragmentPath)` 

Loads GLSL vertex and fragment shaders from file paths. Returns a shader handle.

- `vertexPath`: Path to the `.vs` / `.vert` file.
- `fragmentPath`: Path to the `.fs` / `.frag` file.

---

### `SHADER.FREE(shaderHandle)` 

Unloads the shader from GPU memory and releases its heap slot.

---

### `SHADER.GETLOC(shaderHandle, uniformName)` 

Returns the integer location of a uniform variable by name. Use this to cache locations for per-frame updates.

---

### `SHADER.SETFLOAT(shaderHandle, uniformName, value)` 

Sets a float uniform value in the shader.

---

### `SHADER.SETINT(shaderHandle, uniformName, value)` 

Sets an integer uniform value in the shader.

---

### `SHADER.SETVEC2(shaderHandle, uniformName, x, y)` 

Sets a 2-component vector uniform.

---

### `SHADER.SETVEC3(shaderHandle, uniformName, x, y, z)` 

Sets a 3-component vector uniform.

---

### `SHADER.SETVEC4(shaderHandle, uniformName, x, y, z, w)` 

Sets a 4-component vector uniform.

---

### `SHADER.SETTEXTURE(shaderHandle, uniformName, textureHandle)` 

Binds a texture handle to a shader uniform sampler.

---

## Full Example

This example loads a custom shader, sets a time uniform each frame, and draws with it.

```basic
sh = SHADER.LOAD("custom.vs", "custom.fs")

WHILE NOT WINDOW.SHOULDCLOSE()
    t = TIME.GET()
    SHADER.SETFLOAT(sh, "uTime", t)

    RENDER.BEGINFRAME()
    RENDER.BEGINSHADER(sh)
    DRAW.RECT(100, 100, 200, 200, 255, 255, 255, 255)
    RENDER.ENDSHADER()
    RENDER.ENDFRAME()
WEND

SHADER.FREE(sh)
```
