# Transform — 3D transform matrices

**`Transform.*`** is the recommended API for **4×4 transformation matrices**: where a mesh sits in the world (position, rotation, scale) and how you combine those into one matrix for `Mesh.Draw`.

The name matches what game engines call a **transform**—not raw “linear algebra,” and not easy to confuse with **`Material.*`** (shaders/textures) or **`MATRIX.*`** (other engine handles).

Every 3D object is drawn with a transform handle. `Body3D.GetMatrix()` and related APIs return the same kind of handle, which you can pass straight to `Mesh.Draw`.

---

### `Transform.Identity()`
Returns a new identity matrix handle.

### `Transform.Translate(x, y, z)`
Returns a translation matrix handle. Alias: `Transform.Translation()`.

### `Transform.Rotate(p, y, r)`
Returns a rotation matrix handle from Euler angles in **radians**. Alias: `Transform.Rotation()`.

### `Transform.Scale(sx, sy, sz)`
Returns a non-uniform scale matrix handle.

---

### `Transform.Multiply(a, b)`
Returns a new matrix handle representing the product of two matrices.

### `Transform.Invert(handle)`
Returns the inverse of the given matrix.

### `Transform.Transpose(handle)`
Returns the transpose of the given matrix.

---

### `Transform.Free(handle)`
Releases the matrix from the heap and frees its memory.

---

## Full example: spinning cube

```basic
Window.Open(960, 540, "Transform demo")
Window.SetFPS(60)

cam = Camera.Make()
cam.SetPos(0, 3, 10)
cam.SetTarget(0, 0, 0)
cam.SetFOV(45)

cube = Mesh.MakeCube(2, 2, 2)
mat  = Material.MakeDefault()
Material.SetColor(mat, MATERIAL_MAP_ALBEDO, 100, 180, 255, 255)

xform = Transform.Identity()
angle = 0.0

WHILE NOT Window.ShouldClose()
    angle = angle + 1.2 * Time.Delta()
    Transform.SetRotation(xform, angle * 0.5, angle, angle * 0.3)

    Render.Clear(12, 14, 22)
    cam.Begin()
        Mesh.Draw(cube, mat, xform)
        Draw.Grid(20, 1.0)
    cam.End()
    Render.Frame()
WEND

Transform.Free(xform)
Window.Close()
```

---

## Legacy: `Mat4.*`

Older samples use **`Mat4.*`** with names like **`FromTranslation`**. That API is still supported; see [MAT4.md](MAT4.md) for the exact spellings. New code should prefer **`Transform.*`**.
