# Transform — 3D transform matrices

**`Transform.*`** is the recommended API for **4×4 transformation matrices**: where a mesh sits in the world (position, rotation, scale) and how you combine those into one matrix for `Mesh.Draw`.

The name matches what game engines call a **transform**—not raw “linear algebra,” and not easy to confuse with **`Material.*`** (shaders/textures) or **`MATRIX.*`** (other engine handles).

Every 3D object is drawn with a transform handle. `Body3D.GetMatrix()` and related APIs return the same kind of handle, which you can pass straight to `Mesh.Draw`.

---

## Creating transforms

### `Transform.Identity()`

New identity matrix (no translation, no rotation, scale 1). Returns a handle—usual starting point.

```basic
xform = Transform.Identity()
```

---

### `Transform.Translation(x, y, z)`

Matrix that moves an object to a world position (same idea as the old `Mat4.FromTranslation`).

```basic
xform = Transform.Translation(5, 0, -3)
Mesh.Draw(cube, mat, xform)
```

---

### `Transform.Rotation(rx, ry, rz)`

Rotation from XYZ Euler angles **in radians**. Returns a new handle. Order: XYZ.

```basic
xform = Transform.Rotation(0, PI() / 2, 0) ; 90° around Y
```

---

### `Transform.Scale(sx, sy, sz)`

Scale matrix.

```basic
xform = Transform.Scale(2, 2, 2)
```

---

### `Transform.LookAt(eyeX, eyeY, eyeZ, atX, atY, atZ, upX, upY, upZ)`

View matrix: eye at `(eyeX,eyeY,eyeZ)` looking at `(atX,atY,atZ)` with `up`. For custom cameras or billboards.

---

### `Transform.Perspective(fovY, aspect, near, far)` / `Transform.Ortho(left, right, bottom, top, near, far)`

Projection matrices. Usually the **Camera** handles this; exposed for custom pipelines.

---

## Changing an existing transform (in place)

Prefer these in the game loop so you **reuse one handle** instead of allocating every frame.

### `Transform.SetRotation(handle, rx, ry, rz)`

Sets rotation (radians); keeps translation and scale.

```basic
angle = angle + Time.Delta()
Transform.SetRotation(xform, 0, angle, 0)
```

---

## Combining and inspecting

### `Transform.Multiply(a, b)` / `Transform.Inverse(handle)` / `Transform.Transpose(handle)`

Standard matrix ops. `Multiply` order matters.

```basic
t = Transform.Translation(5, 0, 0)
r = Transform.Rotation(0, angle, 0)
combined = Transform.Multiply(t, r)
```

---

### `Transform.GetElement(handle, row, col)`

Single entry (rows/cols **0–3**).

---

### `Transform.ApplyX(handle, x, y, z)` / `ApplyY` / `ApplyZ`

Multiplies the matrix by the point `(x,y,z)` and returns that **world** X, Y, or Z component. Use this to move a local point into world space (or read a component after the multiply).

```basic
phys_mat = Body3D.GetMatrix(cube_body)
x = Transform.ApplyX(phys_mat, 0, 0, 0)
y = Transform.ApplyY(phys_mat, 0, 0, 0)
z = Transform.ApplyZ(phys_mat, 0, 0, 0)
```

---

## `Transform.Free(handle)`

Release the handle when you are done. If you reuse one matrix all frame with `SetRotation`, allocate once and free once at shutdown—not every frame.

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
