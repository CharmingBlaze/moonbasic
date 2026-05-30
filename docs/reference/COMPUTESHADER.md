# ComputeShader Commands

GPU compute shader dispatch: load GLSL compute programs, create GPU buffers, set uniforms, and dispatch parallel work groups.

Requires **fullruntime** with compute shader support (OpenGL 4.3+ / Vulkan backend).

## Core Workflow

1. `COMPUTESHADER.LOAD(path)` — load a `.glsl` compute shader.
2. `COMPUTESHADER.BUFFERMAKE(sizeBytes)` — allocate a GPU SSBO.
3. `COMPUTESHADER.SETBUFFER(shader, binding, buffer)` — bind buffer to a slot.
4. `COMPUTESHADER.SETINT` / `COMPUTESHADER.SETFLOAT` — push uniforms.
5. `COMPUTESHADER.DISPATCH(shader, groupsX, groupsY, groupsZ)` — execute.
6. Read back results by binding the buffer to a render pass.
7. `COMPUTESHADER.FREE` / `COMPUTESHADER.BUFFERFREE` on exit.

---

## Shader

### `COMPUTESHADER.LOAD(path)` 

Loads a compute shader from a GLSL source file. Returns a **shader handle**.

---

### `COMPUTESHADER.FREE(shader)` 

Frees the compute shader.

---

## Buffers

### `COMPUTESHADER.BUFFERMAKE(sizeBytes)` 

Allocates a GPU Shader Storage Buffer Object (SSBO) of `sizeBytes`. Returns a **buffer handle**.

---

### `COMPUTESHADER.BUFFERFREE(buffer)` 

Frees the GPU buffer.

---

## Setup

### `COMPUTESHADER.SETBUFFER(shader, bindingIndex, buffer)` 

Binds a buffer to the shader's SSBO binding point `bindingIndex` (matches `layout(binding=N)` in GLSL).

---

### `COMPUTESHADER.SETINT(shader, uniformName, value)` 

Sets an integer uniform by name.

---

### `COMPUTESHADER.SETFLOAT(shader, uniformName, value)` 

Sets a float uniform by name.

---

## Dispatch

### `COMPUTESHADER.DISPATCH(shader, groupsX, groupsY, groupsZ)` 

Dispatches the compute shader with `groupsX × groupsY × groupsZ` work groups. The total invocations = groups × `local_size_x/y/z` defined in the shader.

---

## Full Example

Parallel float-array squaring on the GPU.

```glsl
; assets/square.glsl
#version 430
layout(local_size_x = 64) in;
layout(std430, binding = 0) buffer Data { float values[]; };
void main() {
    uint i = gl_GlobalInvocationID.x;
    values[i] = values[i] * values[i];
}
```

```basic
WINDOW.OPEN(960, 540, "ComputeShader Demo")
WINDOW.SETFPS(60)

cs     = COMPUTESHADER.LOAD("assets/square.glsl")
buf    = COMPUTESHADER.BUFFERMAKE(256 * 4)   ; 256 floats
COMPUTESHADER.SETBUFFER(cs, 0, buf)

; dispatch: 4 groups × 64 local = 256 invocations
COMPUTESHADER.DISPATCH(cs, 4, 1, 1)

PRINT "Compute pass done"

COMPUTESHADER.BUFFERFREE(buf)
COMPUTESHADER.FREE(cs)
WINDOW.CLOSE()
```

---

## See also

- [SHADER.md](SHADER.md) — vertex/fragment shaders
- [RENDERTARGET.md](RENDERTARGET.md) — off-screen render targets
- [MATERIAL.md](MATERIAL.md) — assigning shaders to materials
