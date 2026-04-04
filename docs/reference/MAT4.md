# Mat4 — Matrix Math Commands

Commands for creating and manipulating 4×4 transformation matrices. These are
the backbone of 3D positioning, rotation, and scaling in moonBASIC.

Every 3D object (mesh, model) is drawn at a position defined by a matrix.
`Body3D.GetMatrix()` and `CharController` functions also return matrices that
can be passed directly to `Mesh.Draw()`.

---

## Creating Matrices

### `Mat4.Identity()`

Creates a new identity matrix (no translation, no rotation, scale = 1). Returns
a handle. This is the starting point for most transforms.

```basic
xform = Mat4.Identity()
```

---

### `Mat4.FromTranslation(x#, y#, z#)`

Creates a matrix that places an object at the given world position. Returns a
handle.

```basic
; Place a mesh at (5, 0, -3)
xform = Mat4.FromTranslation(5, 0, -3)
Mesh.Draw(cube, mat, xform)
```

---

### `Mat4.FromRotation(rx#, ry#, rz#)`

Creates a rotation matrix from XYZ Euler angles **in radians**. Returns a
handle. The rotation is applied in XYZ order.

```basic
; Rotate 90° around Y
xform = Mat4.FromRotation(0, PI() / 2, 0)
```

Alias: `Mat4.Rotation(rx#, ry#, rz#)`

---

### `Mat4.FromScale(sx#, sy#, sz#)`

Creates a scale matrix. Returns a handle.

```basic
; Scale object to twice its size
xform = Mat4.FromScale(2, 2, 2)
```

---

### `Mat4.LookAt(eyeX#, eyeY#, eyeZ#, atX#, atY#, atZ#, upX#, upY#, upZ#)`

Creates a view matrix that positions an eye at `(eyeX, eyeY, eyeZ)` looking
toward `(atX, atY, atZ)`, with `up` as the up vector. Useful for custom camera
or billboard math.

---

### `Mat4.Perspective(fovY#, aspect#, near#, far#)`

Creates a perspective projection matrix. Normally you don't need this directly
as the camera handles it, but it's available for custom rendering pipelines.

---

### `Mat4.Ortho(left#, right#, bottom#, top#, near#, far#)`

Creates an orthographic projection matrix.

---

## Modifying Existing Matrices

These commands **mutate** an existing matrix handle in place, avoiding an extra
allocation. Prefer these in the game loop.

### `Mat4.SetRotation(handle, rx#, ry#, rz#)`

Updates the rotation component of an existing matrix (radians). Translation and
scale are preserved.

```basic
; Rotate a spinning cube — no new handle allocated each frame
angle# = angle# + Time.Delta()
Mat4.SetRotation(xform, 0, angle#, 0)
```

---

## Matrix Operations

### `Mat4.Multiply(a, b)`

Multiplies two matrices together and returns a new handle. Order matters:
`Multiply(translation, rotation)` produces a rotated-then-translated result.

```basic
; Combine translation and rotation
t = Mat4.FromTranslation(5, 0, 0)
r = Mat4.FromRotation(0, angle#, 0)
combined = Mat4.Multiply(t, r)
```

---

### `Mat4.Inverse(handle)`

Returns the inverse of a matrix as a new handle. Useful for converting world-space
coordinates to object-space.

---

### `Mat4.Transpose(handle)`

Returns the transposed matrix as a new handle.

---

## Reading Matrix Values

### `Mat4.GetElement(handle, row, col)`

Returns the float value at the specified row and column (0-indexed).

### `Mat4.TransformX(handle)` / `Mat4.TransformY(handle)` / `Mat4.TransformZ(handle)`

Returns the world-space X, Y, or Z position encoded in the matrix's translation
column. Useful for reading back a position from a physics body matrix.

```basic
phys_mat = Body3D.GetMatrix(cube_body)
x# = Mat4.TransformX(phys_mat)
y# = Mat4.TransformY(phys_mat)
z# = Mat4.TransformZ(phys_mat)
PRINT "Cube is at " + STR$(x#) + ", " + STR$(y#) + ", " + STR$(z#)
```

---

## Freeing

### `Mat4.Free(handle)`

Releases the matrix memory. Call this when a matrix is no longer needed to
prevent heap growth. In tight game loops where you reuse a matrix each frame
with `SetRotation`, you only need one handle — don't allocate and free each frame.

---

## Full Example: Spinning, Textured Cube

```basic
Window.Open(960, 540, "Mat4 Demo")
Window.SetFPS(60)

cam = Camera.Make()
cam.SetPos(0, 3, 10)
cam.SetTarget(0, 0, 0)
cam.SetFOV(45)

cube = Mesh.MakeCube(2, 2, 2)
mat  = Material.MakeDefault()
Material.SetColor(mat, MATERIAL_MAP_ALBEDO, 100, 180, 255, 255)

; Allocate one matrix and reuse it every frame
xform = Mat4.Identity()

angle# = 0.0

WHILE NOT Window.ShouldClose()
    angle# = angle# + 1.2 * Time.Delta()

    ; Mutate the existing matrix — no allocation
    Mat4.SetRotation(xform, angle# * 0.5, angle#, angle# * 0.3)

    Render.Clear(12, 14, 22)
    cam.Begin()
        Mesh.Draw(cube, mat, xform)
        Draw.Grid(20, 1.0)
    cam.End()
    Render.Frame()
WEND

Mat4.Free(xform)
Window.Close()
```

---

## Full Example: Using Physics Matrix

```basic
; Physics bodies give you a matrix directly — no manual position math needed

Physics3D.Start()
Physics3D.SetGravity(0, -10, 0)

cam = Camera.Make()
cam.SetPos(0, 8, 20)
cam.SetTarget(0, 0, 0)

floor_def = Body3D.Make("static")
Body3D.AddBox(floor_def, 40, 1, 40)
floor_body = Body3D.Commit(floor_def, 0, 0, 0)

cube_def = Body3D.Make("dynamic")
Body3D.AddBox(cube_def, 2, 2, 2)
cube_body = Body3D.Commit(cube_def, 0, 12, 0)
Body3D.SetMass(cube_body, 1.0)

floor_mesh = Mesh.MakeCube(40, 1, 40)
cube_mesh  = Mesh.MakeCube(2, 2, 2)
mat        = Material.MakeDefault()

Window.Open(960, 540, "Physics Matrix Demo")
Window.SetFPS(60)

WHILE NOT Window.ShouldClose()
    Physics3D.Step()

    Render.Clear(10, 12, 20)
    cam.Begin()
        ; Body3D.GetMatrix returns a matrix ready to use directly
        Mesh.Draw(floor_mesh, mat, Body3D.GetMatrix(floor_body))
        Mesh.Draw(cube_mesh,  mat, Body3D.GetMatrix(cube_body))
        Draw.Grid(40, 1.0)
    cam.End()
    Render.Frame()
WEND

Physics3D.Stop()
Window.Close()
```
