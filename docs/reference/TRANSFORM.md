# Transform — 3D transform matrices

**`Transform.*`** (registry **`TRANSFORM.*`**) is the recommended API for **4×4 transformation matrices**: where a mesh sits in the world (position, rotation, scale) and how you combine those into one matrix for **`MESH.DRAW`**.

The name matches what game engines call a **transform**—not raw “linear algebra,” and not easy to confuse with **`Material.*`** (shaders/textures) or **`MATRIX.*`** (other engine handles).

Every 3D object is drawn with a transform handle. Build one from **`TRANSFORM.TRANSLATION`** / **`TRANSFORM.ROTATION`** / **`TRANSFORM.MULTIPLY`**, or use **`CAMERA.GETMATRIX`** / **`MAT4.*`** where applicable — there is no **`BODY3D.GETMATRIX`**; use **`BODY3D.GETPOS`** + **`BODY3D.GETROT`** for rigid bodies (see [PHYSICS3D.md](PHYSICS3D.md)).

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

## Legacy: `Mat4.*`

Older samples use **`Mat4.*`** with names like **`FromTranslation`**. That API is still supported; see [MAT4.md](MAT4.md) for the exact spellings. New code should prefer **`Transform.*`**.
