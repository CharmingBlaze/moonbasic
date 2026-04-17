# Mat4 Commands

4×4 matrix creation, manipulation, and point transformation for 3D math.

Page shape follows [DOC_STYLE_GUIDE.md](../DOC_STYLE_GUIDE.md) (**WAVE pattern**).

> **Compatibility note:** `MAT4.*` is fully supported. New code may also use `TRANSFORM.*` aliases — see [TRANSFORM.md](TRANSFORM.md).

## Core Workflow

1. Create a matrix with a factory (`MAT4.IDENTITY`, `MAT4.FROMROTATION`, `MAT4.PERSPECTIVE`, etc.).
2. Combine matrices with `MAT4.MULTIPLY`.
3. Transform points with `MAT4.TRANSFORMX` / `Y` / `Z`.
4. Free with `MAT4.FREE` when no longer needed.

---

### `MAT4.IDENTITY()` 

Returns a new identity matrix handle.

---

### `MAT4.FROMROTATION(pitch, yaw, roll)` 

Creates a rotation matrix from Euler angles (degrees).

---

### `MAT4.ROTATION(pitch, yaw, roll)` 

Alias for `MAT4.FROMROTATION`.

---

### `MAT4.SETROTATION(matHandle, pitch, yaw, roll)` 

Overwrites the rotation component of an existing matrix.

---

### `MAT4.FROMSCALE(sx, sy, sz)` 

Creates a scale matrix.

---

### `MAT4.FROMTRANSLATION(tx, ty, tz)` 

Creates a translation matrix.

---

### `MAT4.LOOKAT(eyeX, eyeY, eyeZ, targetX, targetY, targetZ, upX, upY, upZ)` 

Creates a look-at view matrix.

---

### `MAT4.PERSPECTIVE(fovY, aspect, near, far)` 

Creates a perspective projection matrix.

- `fovY`: Vertical field of view in degrees.

---

### `MAT4.ORTHO(left, right, bottom, top, near, far)` 

Creates an orthographic projection matrix.

---

### `MAT4.MULTIPLY(a, b)` 

Returns the product of two matrices (`a × b`). The result is a new handle.

---

### `MAT4.INVERSE(matHandle)` 

Returns the inverse of the matrix. The result is a new handle.

---

### `MAT4.TRANSPOSE(matHandle)` 

Returns the transpose of the matrix. The result is a new handle.

---

### `MAT4.GETELEMENT(matHandle, row, col)` 

Returns the float value at the given row and column (0-based).

---

### `MAT4.TRANSFORMX(matHandle, x, y, z)` 

Returns the X component after transforming the point `(x, y, z)` by the matrix.

---

### `MAT4.TRANSFORMY(matHandle, x, y, z)` 

Returns the Y component after transforming the point.

---

### `MAT4.TRANSFORMZ(matHandle, x, y, z)` 

Returns the Z component after transforming the point.

---

### `MAT4.FREE(matHandle)` 

Frees the matrix handle from memory.

---

## Full Example

This example builds a model matrix, applies it to a point, and prints the result.

```basic
; Build a transform: translate then rotate
t = MAT4.FROMTRANSLATION(10.0, 0.0, 5.0)
r = MAT4.FROMROTATION(0.0, 45.0, 0.0)
m = MAT4.MULTIPLY(t, r)

; Transform the origin
px = MAT4.TRANSFORMX(m, 0.0, 0.0, 0.0)
py = MAT4.TRANSFORMY(m, 0.0, 0.0, 0.0)
pz = MAT4.TRANSFORMZ(m, 0.0, 0.0, 0.0)
PRINT "Transformed: " + STR(px) + ", " + STR(py) + ", " + STR(pz)

; Cleanup
MAT4.FREE(t)
MAT4.FREE(r)
MAT4.FREE(m)
```
