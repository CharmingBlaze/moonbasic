# Raycasts (2D and 3D)

Ray–primitive tests return whether a half-line from an **origin** along a **direction** hits geometry, and where. **Distance** is measured along the ray from the origin in **world units** after the direction is normalized (see each section).

---

## 3D (`RAY.*`)

Implemented with Raylib (**requires CGO** on typical Unix builds; on **non-Windows** platforms without CGO, `RAY.*` is unavailable — see build errors that mention enabling CGO). **`RAY2D.*`** does not depend on CGO.

### Ray handle

- **`RAY.MAKE(ox, oy, oz, dx, dy, dz)`** — returns a **handle** to a ray. Origin `(ox,oy,oz)`; direction `(dx,dy,dz)` (need not be unit length; Raylib normalizes internally for collision).
- **`RAY.FREE(handle)`** — frees the ray handle.

**Screen picking:** **`Camera.GetRay`**, **`Camera.GetViewRay`**, and **`Camera.MouseRay`** return a **6-float array handle**: origin then direction — same layout as the six floats passed to **`RAY.MAKE`**. See [CAMERA.md](CAMERA.md).

### Hit queries

Each test has a family of commands with the same arguments except for the suffix:

| Suffix | Meaning |
|--------|---------|
| **`_HIT`** | **`TRUE`** if there is an intersection in the forward half-line (`t ≥ 0` per Raylib). |
| **`_DISTANCE`** | Distance from origin to hit along the ray (or **`0.0`** if no hit). |
| **`_POINTX`**, **`_POINTY`**, **`_POINTZ`** | World position of the hit (or **`0.0`** if no hit). |
| **`_NORMALX`**, **`_NORMALY`**, **`_NORMALZ`** | Surface normal at the hit (or **`0.0`** if none). |

**Sphere** — `RAY.HITSPHERE_* (ray, cx, cy, cz, r)`  
**Axis-aligned box** — `RAY.HITBOX_* (ray, minX, minY, minZ, maxX, maxY, maxZ)`  
**Plane** — `RAY.HITPLANE_* (ray, nx, ny, nz, d)` — plane `nx*x + ny*y + nz*z + d = 0`  
**Triangle** — `RAY.HITTRIANGLE_* (ray, ax, ay, az, bx, by, bz, cx, cy, cz)`  
**Mesh** — `RAY.HITMESH_* (ray, mesh, transform)` — mesh and transform handles  
**Model** — `RAY.HITMODEL_* (ray, model)` — model handle  

Use **`RAY.HITSPHERE_HIT`** (etc.) first; read distance and point only when **`_HIT`** is true.

---

## 2D (`RAY2D.*`)

Pure math in **`runtime/mbcollision/ray2d.go`** — available **with or without CGO**. For each shape, the ray is **`(ox, oy, dx, dy)`**; **`dx`/`dy` must not be zero** (length is normalized internally).

### Circle

**`RAY2D.HITCIRCLE_* (ox, oy, dx, dy, cx, cy, r)`**

Ray vs circle centre **`(cx,cy)`** radius **`r`** (negative **`r`** is treated as absolute). Returns the **first** forward intersection **`t ≥ 0`** (smaller positive root of the quadratic). If the origin is inside the circle, the hit is at **`t = 0`**.

### Axis-aligned rectangle

**`RAY2D.HITRECT_* (ox, oy, dx, dy, minX, minY, maxX, maxY)`**

Slab test on the AABB. **`min*`**/**`max*`** are swapped if given in reverse order. If the origin lies inside the box, **`_DISTANCE`** is **`0`** and the point is the origin.

### Segment

**`RAY2D.HITSEGMENT_* (ox, oy, dx, dy, x1, y1, x2, y2)`**

Ray vs finite segment **`(x1,y1)`–`(x2,y2)`**. Parallel ray and segment yields no hit.

### 2D result fields

| Command | Returns |
|---------|---------|
| **`_HIT`** | **`TRUE`** / **`FALSE`** |
| **`_DISTANCE`** | **`t`** along the normalized ray from the origin, or **`0.0`** if no hit |
| **`_POINTX`**, **`_POINTY`** | Hit coordinates, or **`0.0`** if no hit |

There are no 2D normal commands in this set; derive a normal from the primitive if needed.

---

## Related

- [CAMERA.md](CAMERA.md) — **`Camera.GetRay`**, **`Camera.MouseRay`**
- [COLLISION.md](COLLISION.md) — **`BOXCOLLIDE`**, **`CIRCLECOLLIDE`**, … (overlap tests, not rays)
