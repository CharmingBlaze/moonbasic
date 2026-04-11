# Shader Commands

Commands for loading and using custom shaders.

---

### `Shader.Load(vs, fs)`
Loads GLSL vertex and fragment shaders from file paths. Returns a **shader handle**.

### `Shader.Free(handle)`
Unloads the shader from GPU memory and releases its heap slot.

---

### `Shader.GetLoc(handle, name)`
Returns the location index of a uniform variable by its string name.

### `Shader.SetFloat(handle, name, value)`
Sets a float uniform value in the shader.

### `Shader.SetVec3(handle, name, x, y, z)`
Sets a 3-component vector uniform value in the shader.

### `Shader.SetTexture(handle, name, tex)`
Binds a texture handle to a shader uniform sampler.

```basic
sh = Shader.Load("custom.vs", "custom.fs")
t = TIME.GET()
Shader.SetFloat(sh, "uTime", t)
Shader.Free(sh)
```
