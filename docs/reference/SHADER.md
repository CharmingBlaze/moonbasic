# Shader Commands

Commands for loading and using custom shaders.

---

### `Shader.Load(vertexPath$, fragmentPath$)`

Loads a GLSL vertex and fragment shader from files. Returns a handle to the shader program.

- `vertexPath$`: Path to the vertex shader file (`.vs`).
- `fragmentPath$`: Path to the fragment shader file (`.fs`).

---

### `Render.BeginShader(shaderHandle)`

Begins a custom shader mode. All subsequent drawing will be processed by this shader.

- `shaderHandle`: The handle of the shader to use.

---

### `Render.EndShader()`

Ends the custom shader mode.

---

### `Material.SetShader(materialHandle, shaderHandle)`

Assigns a custom shader to a material. This is the preferred way to apply shaders to 3D models.

- `materialHandle`: The handle of the material.
- `shaderHandle`: The handle of the shader.

---

## Uniforms and lifecycle

Uniform names are **string** arguments (bytecode string table).

| Command | Purpose |
|--------|---------|
| `Shader.Free(shader)` | Unloads shader GPU program (heap handle). |
| `Shader.GetLoc(shader, name$)` | Raylib location index (`int`). |
| `Shader.SetFloat(shader, name$, value#)` | |
| `Shader.SetVec2` / `SetVec3` / `SetVec4` | Name + float components. |
| `Shader.SetInt(shader, name$, value)` | |
| `Shader.SetTexture(shader, name$, textureHandle)` | Binds a `TEXTURE` heap handle. |

```basic
sh = Shader.Load("custom.vs", "custom.fs")
t# = TIME.GET()
Shader.SetFloat(sh, "uTime", t#)
Shader.Free(sh)
```
