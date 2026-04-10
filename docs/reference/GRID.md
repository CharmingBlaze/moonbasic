# Tactical grid (`GRID.*`)

Logical **XZ** cell storage for strategy / tile logic on top of 3D. Data lives in a **flat array** (`cells[ix + iz * width]`). **CGO + Raylib** required for **`GRID.DRAW`** and **`GRID.RAYCAST`**; without CGO, commands return an error.

See also **[TERRAIN.md](TERRAIN.md)** for heightfields and **[TEXTURE.md](TEXTURE.md)** for atlas animation on entities.

---

## Commands

| Command | Arguments | Returns | Notes |
|--------|-----------|---------|--------|
| **`GRID.CREATE`** | `width`, `depth`, `cellSize` | handle | Origin **(0,0)** in world XZ; cell size in world units |
| **`GRID.FREE`** | `grid` | — | |
| **`GRID.SETCELL`** | `grid`, `ix`, `iz`, `type` | — | Opaque **`int32`** per cell (walkable, blocked, … — you define values) |
| **`GRID.GETCELL`** | `grid`, `ix`, `iz` | int | **0** if out of bounds |
| **`GRID.WORLDTOCELL`** | `grid`, `worldX`, `worldZ` | array handle | **`[ix, iz]`** |
| **`GRID.DRAW`** | `grid`, `r`, `g`, `b` **or** `grid`, `r`, `g`, `b`, `a` | — | Debug lines on **Y = 0** or baked **Y** from **`GRID.FOLLOWTERRAIN`** |
| **`GRID.SNAP`** | `grid`, `entity`, `ix`, `iz` | — | Moves entity to **cell center**; **Y** from **`GRID.FOLLOWTERRAIN`** when set |
| **`GRID.GETPATH`** | `grid`, `sx`, `sz`, `ex`, `ez` | array handle | BFS path in world XZ; packed **`[ix0, iz0, ix1, iz1, …]`**; empty if blocked |
| **`GRID.FOLLOWTERRAIN`** | `grid`, `terrain` | — | Precomputes **Y** per cell center from **[`Terrain.GetHeight`](TERRAIN.md)** |
| **`GRID.PLACEENTITY`** | `grid`, `ix`, `iz`, `entity` | — | Records occupant for **`GRID.GETNEIGHBORS`** |
| **`GRID.RAYCAST`** | `grid`, `screenX`, `screenY` | array handle | Cell under mouse ray vs **ground plane**; **`[-1,-1]`** if miss |
| **`GRID.GETNEIGHBORS`** | `grid`, `ix`, `iz`, `radius` | array handle | Entity IDs in **Chebyshev** distance **≤ radius** with **`GRID.PLACEENTITY`** set |

---

## Example

```text
g = GRID.CREATE(32, 32, 2.0)
GRID.FOLLOWTERRAIN(g, myTerrain)
WHILE Window.Open()
    TEXTURE.TICKALL()
    c = GRID.RAYCAST(g, Input.MouseX(), Input.MouseY())
    Begin3D()
        Terrain.Draw(myTerrain)
        GRID.DRAW(g, 255, 255, 255, 80)
    End3D()
WEND
```

Planned / not in runtime yet: **`Terrain.ApplyMap`**, **`Terrain.ApplyTiles`**, **`Image.LoadSequence`**, **`Image.LoadGIF`**, **`Entity.SetAnimation`** — use **`TEXTURE.LOADANIM`** + **`TEXTURE.TICKALL`** for animated water and **`ENTITY.SCROLLMATERIAL`** for UV-scrolling surfaces instead.
