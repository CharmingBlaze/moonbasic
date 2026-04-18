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
Loads GLSL shaders from file paths. Returns a shader handle.

- **Arguments**:
    - `vertexPath`: (String) Path to vertex shader.
    - `fragmentPath`: (String) Path to fragment shader.
- **Returns**: (Handle) The new shader handle.
- **Example**:
    ```basic
    sh = SHADER.LOAD("water.vert", "water.frag")
    ```

---

### `SHADER.FREE(shaderHandle)`
Unloads the shader and releases its heap slot.

---

### `SHADER.GETLOC(shaderHandle, uniformName)`
Returns the location index of a uniform variable.

- **Arguments**:
    - `shaderHandle`: (Handle) The shader to query.
    - `uniformName`: (String) Variable name in GLSL.
- **Returns**: (Integer) Location index (cache this for performance).

---

### `SHADER.SETFLOAT(shaderHandle, uniformName, value)` / `SETINT`
Sets a scalar uniform value.

- **Returns**: (Handle) The shader handle (for chaining).

---

### `SHADER.SETVEC2(shaderHandle, uniformName, x, y)` / `SETVEC3` / `SETVEC4`
Sets vector uniform values.

- **Returns**: (Handle) The shader handle (for chaining).

---

### `SHADER.SETTEXTURE(shaderHandle, uniformName, textureHandle)`
Binds a texture to a shader sampler.

- **Returns**: (Handle) The shader handle (for chaining).

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
