# Shape Commands

Create reusable Jolt Physics collision shape handles. Shapes can be queried for their dimensions and passed to `BODY3D.*` or `TRIGGER.*` instead of calling `BODY3D.ADDBOX` etc. inline.

Requires **CGO + Jolt** (Windows/Linux fullruntime).

## Core Workflow

1. `SHAPE.CREATEBOX(hx, hy, hz)` / `SHAPE.CREATESPHERE(r)` / `SHAPE.CREATECAPSULE(r, h)` / `SHAPE.CREATECYLINDER(r, h)` — allocate a shape handle.
2. Pass the handle to `BODY3D.ADDMESH` or `TRIGGER.CREATE` as needed.
3. Query dimensions with `SHAPE.GETWIDTH`, `SHAPE.GETHEIGHT`, `SHAPE.GETDEPTH`, `SHAPE.GETRADIUS`.

---

## Creation

### `SHAPE.CREATEBOX(halfX, halfY, halfZ)` 

Creates a box collision shape. Arguments are **half-extents**. Returns a **shape handle**.

---

### `SHAPE.CREATESPHERE(radius)` 

Creates a sphere shape with the given `radius`. Returns a **shape handle**.

---

### `SHAPE.CREATECAPSULE(radius, height)` 

Creates a capsule shape. `height` is the cylindrical section height (not total height). Returns a **shape handle**.

---

### `SHAPE.CREATECYLINDER(radius, height)` 

Creates a cylinder shape. Returns a **shape handle**.

---

## Inspection

### `SHAPE.GETTYPE(shape)` 

Returns the shape type as an integer: `1` = Box, `2` = Sphere, `3` = Capsule, `4` = Cylinder.

---

### `SHAPE.GETWIDTH(shape)` / `SHAPE.GETSIZEX(shape)` 

Returns the first dimension (half-extent X or radius).

---

### `SHAPE.GETHEIGHT(shape)` / `SHAPE.GETSIZEY(shape)` 

Returns the second dimension (half-extent Y or height).

---

### `SHAPE.GETDEPTH(shape)` / `SHAPE.GETSIZEZ(shape)` 

Returns the third dimension (half-extent Z).

---

### `SHAPE.GETRADIUS(shape)` 

Returns the radius for sphere and capsule shapes. Same as `SHAPE.GETWIDTH` for spheres.

---

## Full Example

Reusable shape handles shared across multiple bodies.

```basic
WINDOW.OPEN(960, 540, "Shape Demo")
WINDOW.SETFPS(60)

PHYSICS3D.START()
PHYSICS3D.SETGRAVITY(0, -10, 0)

; shared sphere shape
sphereShape = SHAPE.CREATESPHERE(0.5)
PRINT "Shape type: " + STR(SHAPE.GETTYPE(sphereShape))   ; 2 = sphere
PRINT "Radius: "     + STR(SHAPE.GETRADIUS(sphereShape)) ; 0.5

; spawn 5 balls using the same shape
FOR i = 1 TO 5
    def  = BODY3D.CREATE("DYNAMIC")
    BODY3D.ADDSPHERE(def, SHAPE.GETRADIUS(sphereShape))
    ball = BODY3D.COMMIT(def, i * 2 - 5, 5, 0)
NEXT i

; static floor
floorDef = BODY3D.CREATE("STATIC")
BODY3D.ADDBOX(floorDef, 15, 0.5, 15)
floor = BODY3D.COMMIT(floorDef, 0, -0.5, 0)

cam = CAMERA.CREATE()
CAMERA.SETPOS(cam, 0, 6, -12)
CAMERA.SETTARGET(cam, 0, 0, 0)

WHILE NOT WINDOW.SHOULDCLOSE()
    PHYSICS3D.UPDATE()
    RENDER.CLEAR(20, 20, 35)
    RENDER.BEGIN3D(cam)
        DRAW3D.GRID(20, 1.0)
    RENDER.END3D()
    RENDER.FRAME()
WEND

PHYSICS3D.STOP()
WINDOW.CLOSE()
```

---

## Extended Command Reference

| Command | Description |
|--------|-------------|
| `SHAPE.MAKEBOX(hx, hy, hz)` | Deprecated alias of `SHAPE.BOX`. |
| `SHAPE.MAKESPHERE(radius)` | Deprecated alias of `SHAPE.SPHERE`. |
| `SHAPE.MAKECAPSULE(radius, height)` | Deprecated alias of `SHAPE.CAPSULE`. |
| `SHAPE.MAKECYLINDER(radius, height)` | Deprecated alias of `SHAPE.CYLINDER`. |

---

## See also

- [BODY3D.md](BODY3D.md) — rigid bodies using these shapes
- [TRIGGER.md](TRIGGER.md) — sensor zones using `TRIGGER.CREATE(shapeHandle)`
- [PHYSICS3D.md](PHYSICS3D.md) — world setup
