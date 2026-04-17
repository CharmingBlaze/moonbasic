# Essential Blitz-style API (implemented in moonBASIC)

Quick map from **familiar Blitz3D / DBPro names** to **moonBASIC** commands. Identifiers are case-insensitive; registry uses **`NAMESPACE.NAME`**. Full parity notes: [BLITZ_COMMAND_INDEX.md](BLITZ_COMMAND_INDEX.md), [BLITZ3D.md](BLITZ3D.md).

## Core Workflow

This page is a compatibility lookup. General workflow:

1. Use Blitz-style names as shown — they alias the registry commands listed.
2. Where a **`NAMESPACE.*`** form is noted, prefer it in new scripts for IDE autocomplete.
3. See [ENTITY.md](ENTITY.md), [CAMERA.md](CAMERA.md), [AUDIO.md](AUDIO.md) for full API docs.

---

## 1. Entity system (the “Blitz feel”)

Entities use **integer ids** (`entity`) or **`CUBE()`** / **`SPHERE()`** handles with dot methods — see [ENTITY.md](ENTITY.md), [BLITZ3D.md](BLITZ3D.md).

| Concept | moonBASIC |
|--------|-----------|
| **Position** | **`Entity.Position()`** / **`Entity.SetPosition()`** |
| **Move** | **`Entity.Move()`** / **`Entity.Translate()`** |
| **Rotate** | **`Entity.SetRotation()`** / **`Entity.Rotate()`** |
| **Turn** | **`Entity.Turn()`** |
| **Scale** | **`Entity.Scale()`** |
| **Parent** | **`Entity.Parent()`** / **`Entity.Unparent()`** |
| **Color** | **`Entity.Color()`** / **`Entity.Alpha()`** |
| **Distance** | **`Entity.Distance()`** |

---

## 2. Meshes & primitives

### `Mesh.MakeCube(w, h, d)` 
Creates a procedural box mesh.

---

### `Mesh.MakeSphere(r, rings, slices)` 
Creates a procedural sphere mesh.

---

### `Mesh.Load(path)` 
Loads a mesh from a file.

---

## 3. Camera & picking

### `Camera.GetRay(cam, sx, sy)` 
Returns a screen-to-world ray handle.

---

### `Camera.Project(cam, wx, wy, wz)` 
Projects 3D point to screen coordinates.

---

### `Camera.LookAt(cam, x, y, z)` 
Aims camera at a world point.

---

## 4. 2D / screen space

### `Sprite.Load(path)` 
Loads a sprite handle.

---

### `Sprite.Draw(id, x, y)` 
Draws sprite at pixel position.

---

### `Sprite.Hit(a, b)` 
Checks for sprite collision (oriented quads matching **`SPRITE.DRAW`** — scale, origin, rotation; see [SPRITE.md](SPRITE.md)).

---

## 5. Logic & game juice

### `CurveValue(dest, cur, speed)` 
Interpolates value toward target.

---

### `CurveAngle(dest, cur, speed)` 
Interpolates angle toward target.

---

### `Rnd(min, max)` 
Returns inclusive random integer.

---

## Full Example

Blitz-style 3D cube with camera and input using familiar names.

```basic
WINDOW.OPEN(800, 600, "Blitz Style")
WINDOW.SETFPS(60)

cam  = Camera.Make()
Camera.SetPos(cam, 0, 4, -8)
Camera.SetTarget(cam, 0, 0, 0)

cube = Entity.CreateCube(1.5)
Entity.SetPos(cube, 0, 0, 0)

WHILE NOT WINDOW.SHOULDCLOSE()
    dt = DT()
    Entity.Turn(cube, 0, 40 * dt, 0)
    Entity.Update(dt)

    RENDER.CLEAR(20, 20, 40)
    RENDER.BEGIN3D(cam)
        Entity.DrawAll()
        DRAW3D.GRID(10, 1.0)
    RENDER.END3D()
    RENDER.FRAME()
WEND

Entity.Free(cube)
Camera.Free(cam)
WINDOW.CLOSE()
```

---

## See also

- [PROGRAMMING.md](../PROGRAMMING.md) — main loop, **`TIME.DELTA`**, **`ENTITY.UPDATE`**, **`RENDER.FRAME`**
- [API_CONSISTENCY.md](../API_CONSISTENCY.md) — every registered name
- [BLITZ3D.md](BLITZ3D.md) — camera, mesh, render pipeline aliases
