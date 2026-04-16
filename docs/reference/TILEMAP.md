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

### `Tilemap.Load(path)`
Loads a Tiled `.tmx` file and its associated tileset image. Returns a **tilemap handle**.

### `Tilemap.Free(handle)`
Unloads the tilemap texture and frees all associated data from memory.

---

### `Tilemap.Draw(handle, offsetX, offsetY)`
Draws all tile layers of the map, shifted by a pixel offset. Use the offset to implement scrolling.

### `Tilemap.IsSolid(handle, tileX, tileY)`
Returns `TRUE` if the tile at grid position `(tileX, tileY)` has collision. Collision data is loaded from a layer named `"collision"` in the Tiled map.

---

### `Tilemap.GetTile(handle, layerName, tileX, tileY)`
Returns the tile GID (global ID) at a specific grid position on a named layer. Returns `0` for empty tiles.

### `Tilemap.SetTile(handle, layerName, x, y, id)`
Changes a tile in a layer at runtime. Set `id = 0` to erase.

---

### `Tilemap.Width(handle)` / `Tilemap.Height(handle)`
Returns map dimensions in tiles.

### `Tilemap.TileWidth(handle)` / `Tilemap.TileHeight(handle)`
Returns tile dimensions in pixels.

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
