# Tilemap Commands

Commands for loading and rendering Tiled (`.tmx`) tilemaps.

moonBASIC has built-in support for the [Tiled map editor](https://www.mapeditor.org/).
Export your maps in `.tmx` format (XML), place the tileset image next to the
`.tmx` file, and load it with `Tilemap.Load()`.

## Core Workflow

1. **Design your map** in Tiled. Use a single external tileset image (PNG).
2. **Mark collision** by adding a layer named `"collision"` in Tiled — any
   non-zero tile on that layer will be treated as solid.
3. **Load**: Call `Tilemap.Load()` once, before the main loop.
4. **Draw**: Call `Tilemap.Draw()` each frame to render all layers.
5. **Collide**: Use `Tilemap.IsSolid()` to check if a tile position is blocked.
6. **Free**: Call `Tilemap.Free()` when you are done.

---

### `TILEMAP.LOAD(path)`
Loads a Tiled `.tmx` map and its tileset texture.

- **Arguments**:
    - `path`: (String) File path to the `.tmx` file.
- **Returns**: (Handle) The new tilemap handle.
- **Example**:
    ```basic
    map = TILEMAP.LOAD("level1.tmx")
    ```

---

### `TILEMAP.DRAW(handle)`
Renders all tile layers at the map origin (no scroll offset parameter).

- **Arguments**:
    - `handle`: (Handle) The tilemap.
- **Returns**: (None)

For a scrolling view, offset your player/sprites with **`Camera2D`**, draw the map, then draw entities in screen space — or keep the map small enough to fit the window. Runnable sample: [`examples/tilemap/main.mb`](../../examples/tilemap/main.mb).

---

### `TILEMAP.ISSOLID(handle, tx, ty)`
Returns `TRUE` if the tile at grid coordinates `(tx, ty)` has collision.

- **Returns**: (Boolean)

---

### `TILEMAP.GETTILE(handle, layerName, tx, ty)` / `SETTILE`
Accesses specific tiles on a named layer.

- **Returns**: (Integer) The tile GID for `GETTILE`.

---

### `TILEMAP.FREE(handle)`
Releases the tilemap and its texture from memory.

---

## Full Example

See **[`examples/tilemap/main.mb`](../../examples/tilemap/main.mb)** (TMX + collision layer + player square). Minimal loop:

```basic
map = TILEMAP.LOAD("examples/tilemap/assets/level1.tmx")
TILEMAP.SETTILESIZE(map, 16, 16)

WHILE NOT WINDOW.SHOULDCLOSE()
    ; move player, then:
    tileX = INT(px / 16)
    tileY = INT(py / 16)
    IF NOT TILEMAP.ISSOLID(map, tileX, tileY) THEN
        ; apply movement
    ENDIF

    RENDER.CLEAR(20, 24, 32)
    TILEMAP.DRAW(map)
    DRAW.RECTANGLE(INT(px) - 6, INT(py) - 6, 12, 12, 255, 220, 80, 255)
    RENDER.FRAME()
WEND

TILEMAP.FREE(map)
```

### Platformer-style collision (manual)

For jumping and floor snaps, sample tile coordinates from pixel position (same idea as a platformer):

```basic
tx = INT(px / TILE_SIZE)
ty = INT((py + 24) / TILE_SIZE)
IF TILEMAP.ISSOLID(map, tx, ty) THEN
    py = ty * TILE_SIZE - 24
    pvy = 0
    on_ground = 1
ENDIF
```

---

## Extended Command Reference

### Map info

| Command | Description |
|--------|-------------|
| `TILEMAP.WIDTH(map)` | Returns map width in tiles. |
| `TILEMAP.HEIGHT(map)` | Returns map height in tiles. |
| `TILEMAP.LAYERCOUNT(map)` | Returns number of tile layers. |
| `TILEMAP.SETTILESIZE(map, w, h)` | Override displayed tile pixel size. |

### Layer drawing

| Command | Description |
|--------|-------------|
| `TILEMAP.DRAWLAYER(map, layerIndex)` | Draw a single layer by index (no scroll offset). |

### Collision

| Command | Description |
|--------|-------------|
| `TILEMAP.SETCOLLISION(map, layer, bool)` | Enable/disable collision on a layer. |
| `TILEMAP.COLLISIONAT(map, tileX, tileY)` | Returns `TRUE` if tile at grid position is solid. |
| `TILEMAP.ISSOLIDCATEGORY(map, tileX, tileY, category)` | Returns `TRUE` if tile matches collision category. |
| `TILEMAP.MERGECOLLISIONLAYER(map, layer)` | Merge a layer into the physics collision mesh. |

## See also

- [PHYSICS2D.md](PHYSICS2D.md) — 2D collision bodies
- [SPRITE.md](SPRITE.md) — sprite animation
