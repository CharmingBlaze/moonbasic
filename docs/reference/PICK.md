# Pick Commands

Staged ray-cast API for Jolt-backed entity picking. Stages a ray (origin, direction, options), casts it, then reads back hit data. Higher-level than `RAY.*` — works with entity collision layers.

Requires **CGO + Jolt** (Windows/Linux fullruntime).

## Core Workflow

1. Stage the ray with `PICK.ORIGIN` + `PICK.DIRECTION` **or** use `PICK.FROMCAMERA` to build from a camera + screen pixel.
2. Optional: `PICK.MAXDIST`, `PICK.LAYERMASK`.
3. `PICK.CAST()` — run the raycast; returns the hit entity id or `0`.
4. Read results: `PICK.HIT`, `PICK.DIST`, `PICK.X/Y/Z`, `PICK.NX/NY/NZ`, `PICK.ENTITY`.

For raw geometry intersection (boxes, spheres, meshes) without entities, use [RAY.md](RAY.md).

---

## Staging

### `PICK.ORIGIN(x, y, z)` 

Sets the ray origin for the next `PICK.CAST`.

---

### `PICK.DIRECTION(dx, dy, dz)` 

Sets the ray direction. Length acts as max travel unless `PICK.MAXDIST` is set explicitly.

---

### `PICK.MAXDIST(distance)` 

Sets an explicit maximum ray travel distance. Normalises the direction internally.

---

### `PICK.LAYERMASK(mask)` 

Bitmask: bit `i` set = accept entities on collision layer `i`. Pass `0` to accept all layers.

---

### `PICK.RADIUS(radius)` 

Reserved for future sphere-cast. Non-zero currently returns an error.

---

### `PICK.FROMCAMERA(camHandle, screenX, screenY)` 

Stages origin and direction from a camera handle and screen-space pixel coordinates. Sets a default `MAXDIST` if not already set.

---

## Casting

### `PICK.CAST()` 

Runs the staged Jolt raycast. Returns the hit entity id, or `0` if no hit.

---

### `PICK.SCREENCAST(camHandle, screenX, screenY)` 

Convenience: `FROMCAMERA` + `CAST` in one call. Returns entity id or `0`.

---

## Results

### `PICK.HIT()` 

Returns `TRUE` if the last `PICK.CAST` or `PICK.SCREENCAST` produced a hit.

---

### `PICK.ENTITY()` 

Returns the entity id of the last hit (same as `PICK.CAST` return value; reads linked `BODY3D` entities only).

---

### `PICK.DIST()` 

Returns the distance along the ray to the hit point.

---

### `PICK.X()` / `PICK.Y()` / `PICK.Z()` 

Returns the world-space hit point coordinates.

---

### `PICK.NX()` / `PICK.NY()` / `PICK.NZ()` 

Returns the surface normal at the hit point.

---

## Full Example

Click to select an entity in a 3D scene.

```basic
WINDOW.OPEN(960, 540, "Pick Demo")
WINDOW.SETFPS(60)

cam  = CAMERA.CREATE()
CAMERA.SETPOS(cam, 0, 5, -10)
CAMERA.SETTARGET(cam, 0, 0, 0)

boxes = ARRAY.MAKE(4)
FOR i = 0 TO 3
    e = ENTITY.CREATECUBE(1.5)
    ENTITY.SETPOS(e, (i - 1.5) * 3, 0, 0)
    ENTITY.ADDPHYSICS(e, "STATIC")
    ARRAY.SET(boxes, i, e)
NEXT i

selected = 0

WHILE NOT WINDOW.SHOULDCLOSE()
    ENTITY.UPDATE(TIME.DELTA())

    IF MOUSE.PRESSED(0)
        hit = PICK.SCREENCAST(cam, MOUSE.X(), MOUSE.Y())
        IF hit > 0 THEN selected = hit
    END IF

    RENDER.CLEAR(20, 25, 35)
    RENDER.BEGIN3D(cam)
        FOR i = 0 TO 3
            e = ARRAY.GET(boxes, i)
            IF e = selected THEN
                ENTITY.COLOR(e, 255, 180, 50)
            ELSE
                ENTITY.COLOR(e, 80, 120, 200)
            END IF
        NEXT i
        ENTITY.DRAWALL()
        DRAW3D.GRID(20, 1.0)
    RENDER.END3D()
    RENDER.FRAME()
WEND

WINDOW.CLOSE()
```

---

## See also

- [RAY.md](RAY.md) — raw geometry intersection (no entity system)
- [ENTITY.md](ENTITY.md) — entity collision layers
- [CAMERA.md](CAMERA.md) — `CAMERA.GETRAY` for custom ray building
