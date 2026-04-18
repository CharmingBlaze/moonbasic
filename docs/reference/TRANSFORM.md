# Transform Commands

4×4 transformation matrices for positioning, rotating, and scaling 3D objects.

Page shape follows [DOC_STYLE_GUIDE.md](../DOC_STYLE_GUIDE.md) (**WAVE pattern**).

## Core Workflow

1. Create a matrix with `TRANSFORM.IDENTITY`, `TRANSFORM.TRANSLATION`, `TRANSFORM.ROTATION`, or `TRANSFORM.SCALE`.
2. Combine with `TRANSFORM.MULTIPLY`.
3. Pass to `MESH.DRAW` or other 3D commands.
4. Free with `TRANSFORM.FREE`.

For legacy `MAT4.*` naming see [MAT4.md](MAT4.md).

---

### `TRANSFORM.IDENTITY()`
Returns a new identity matrix.

- **Returns**: (Handle)

---

### `TRANSFORM.TRANSLATION(x, y, z)` / `ROTATION` / `SCALE`
Creates a transformation matrix.

- **Arguments**:
    - `x, y, z`: (Float) Translation or scale components.
    - `p, y, r`: (Float) Euler angles (pitch, yaw, roll).
- **Returns**: (Handle) The new matrix handle.
- **Example**:
    ```basic
    m = TRANSFORM.TRANSLATION(0, 10, 0)
    ```

---

### `TRANSFORM.MULTIPLY(a, b)`
Returns the product of two matrices.

- **Returns**: (Handle) The new matrix handle.

---

### `TRANSFORM.FREE(handle)`
Releases the matrix handle from memory.

---

## Full Example

Registry keys: **`TRANSFORM.*`**, **`MESH.*`**, **`MATERIAL.*`**, **`RENDER.Begin3D`**, **`DRAW.GRID`**, **`TIME.DELTA`**. The material map index **`0`** is albedo (**`MATERIAL_MAP_ALBEDO`** at runtime when globals are seeded).

```basic
WINDOW.OPEN(960, 540, "Transform demo")
WINDOW.SETFPS(60)

cam = CAMERA.CREATE()
cam.SETPOS(0, 3, 10)
cam.SETTARGET(0, 0, 0)
cam.SETFOV(45)

cube = MESH.CREATECUBE(2, 2, 2)
mat = MATERIAL.CREATEDEFAULT()
MATERIAL.SETCOLOR(mat, 0, 100, 180, 255, 255)

xform = TRANSFORM.IDENTITY()
angle = 0.0

WHILE NOT WINDOW.SHOULDCLOSE()
    angle = angle + 1.2 * TIME.DELTA()
    TRANSFORM.SETROTATION(xform, angle * 0.5, angle, angle * 0.3)

    RENDER.CLEAR(12, 14, 22)
    RENDER.Begin3D(cam)
        MESH.DRAW(cube, mat, xform)
        DRAW.GRID(20, 1.0)
    RENDER.END3D()
    RENDER.FRAME()
WEND

TRANSFORM.FREE(xform)
WINDOW.CLOSE()
```

---

## Extended Command Reference

### Matrix operations

| Command | Description |
|--------|-------------|
| `TRANSFORM.INVERSE(xform)` | Returns the inverse of the transform matrix. |
| `TRANSFORM.TRANSPOSE(xform)` | Returns the transpose of the transform matrix. |
| `TRANSFORM.GETELEMENT(xform, row, col)` | Returns a single matrix element at `[row, col]`. |

### Build from parameters

| Command | Description |
|--------|-------------|
| `TRANSFORM.PERSPECTIVE(fov, aspect, near, far)` | Build a perspective projection matrix. |
| `TRANSFORM.ORTHO(left, right, bottom, top, near, far)` | Build an orthographic projection matrix. |
| `TRANSFORM.LOOKAT(ex,ey,ez, tx,ty,tz, ux,uy,uz)` | Build a view (look-at) matrix. |

### Apply to coordinates

| Command | Description |
|--------|-------------|
| `TRANSFORM.APPLYX(xform, x, y, z)` | Returns X component after multiplying `[x,y,z]` by matrix. |
| `TRANSFORM.APPLYY(xform, x, y, z)` | Returns Y component. |
| `TRANSFORM.APPLYZ(xform, x, y, z)` | Returns Z component. |

## See also

- [MAT4.md](MAT4.md) — lower-level 4×4 matrix operations
- [CAMERA.md](CAMERA.md) — `CAMERA.GETMATRIX`

