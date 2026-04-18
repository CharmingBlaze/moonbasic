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
Returns a new identity matrix.

- **Returns**: (Handle)

---

### `MAT4.ROTATION(p, y, r)` / `TRANSLATION` / `SCALE`
Creates a transformation matrix from Euler angles or components.

- **Arguments**:
    - `p, y, r`: (Float) Euler angles in degrees.
    - `tx, ty, tz`: (Float) Translation components.
    - `sx, sy, sz`: (Float) Scale factors.
- **Returns**: (Handle) The new matrix handle.

---

### `MAT4.MULTIPLY(a, b)`
Returns the product of two matrices.

- **Returns**: (Handle) The new matrix handle.

---

### `MAT4.TRANSFORMX(handle, x, y, z)` / `Y` / `Z`
Transforms a point by the matrix and returns a single component.

- **Returns**: (Float)

---

### `MAT4.FREE(handle)`
Releases the matrix handle from memory.

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
