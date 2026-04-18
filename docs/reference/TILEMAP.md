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

### `TILEMAP.DRAW(handle, ox, oy)`
Renders the tilemap with a pixel offset (scrolling).

- **Arguments**:
    - `handle`: (Handle) The tilemap.
    - `ox, oy`: (Float) Pixel offsets.
- **Returns**: (None)

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

## Full Example: Scrolling Platformer Map

```basic
WINDOW.OPEN(960, 540, "Tilemap Demo")
WINDOW.SETFPS(60)

map = TILEMAP.LOAD("assets/maps/level1.tmx")

px = 100
py = 100
pvx = 0
pvy = 0
on_ground = 0
TILE_SIZE = 16

WHILE NOT WINDOW.SHOULDCLOSE()
    dt = TIME.DELTA()

    IF INPUT.KEYDOWN(KEY_A) THEN pvx = pvx - 400 * dt
    IF INPUT.KEYDOWN(KEY_D) THEN pvx = pvx + 400 * dt
    pvx = pvx * 0.85

    IF on_ground AND INPUT.KEYPRESSED(KEY_SPACE) THEN pvy = -500

    pvy = pvy + 900 * dt
    px = px + pvx * dt
    py = py + pvy * dt

    on_ground = 0

    tx = INT(px / TILE_SIZE)
    ty = INT((py + 24) / TILE_SIZE)
    IF TILEMAP.ISSOLID(map, tx, ty) THEN
        py = ty * TILE_SIZE - 24
        pvy = 0
        on_ground = 1
    ENDIF

    ty_top = INT((py - 16) / TILE_SIZE)
    IF TILEMAP.ISSOLID(map, tx, ty_top) THEN
        py = (ty_top + 1) * TILE_SIZE + 16
        pvy = 0
    ENDIF

    cam_x = INT(px) - 480
    cam_y = INT(py) - 270

    RENDER.CLEAR(40, 60, 80)
    TILEMAP.DRAW(map, -cam_x, -cam_y)
    DRAW.RECTANGLE(INT(px) - cam_x - 8, INT(py) - cam_y - 16, 16, 28, 255, 200, 80, 255)
    RENDER.FRAME()
WEND

TILEMAP.FREE(map)
WINDOW.CLOSE()
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
| `TILEMAP.DRAWLAYER(map, layer, ox, oy)` | Draw a specific layer by index. |

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
