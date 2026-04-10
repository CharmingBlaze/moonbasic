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

## Loading & Freeing

### `Tilemap.Load(path$)`

Loads a Tiled `.tmx` file and its associated tileset image. Returns a handle to
the tilemap.

- `path`: Path to the `.tmx` file. The tileset image must be in the same
  directory (or the path specified inside the `.tmx`).

```basic
map = Tilemap.Load("assets/maps/level1.tmx")
```

---

### `Tilemap.Free(handle)`

Unloads the tilemap texture and frees all associated data.

---

## Drawing

### `Tilemap.Draw(handle, offsetX, offsetY)`

Draws all tile layers of the map, shifted by a pixel offset. Use the offset
to implement scrolling.

- `offsetX`, `offsetY`: Pixel offset applied to the entire map (can be negative
  to scroll right/down).

```basic
; Scroll the map so the player is centered
cam_x = INT(player_x# - 480)
cam_y = INT(player_y# - 270)
Tilemap.Draw(map, -cam_x, -cam_y)
```

---

### `Tilemap.DrawLayer(handle, layerName$, offsetX, offsetY)`

Draws a single named layer. Use this when you need to draw some layers behind
the player and others in front.

- `layerName`: The layer name exactly as it appears in Tiled (e.g. `"ground"`,
  `"decoration"`).

```basic
; Draw background layers first
Tilemap.DrawLayer(map, "ground", -cam_x, -cam_y)
Tilemap.DrawLayer(map, "decoration", -cam_x, -cam_y)

; Draw the player here

; Draw foreground layer on top of the player
Tilemap.DrawLayer(map, "foreground", -cam_x, -cam_y)
```

---

## Collision

### `Tilemap.IsSolid(handle, tileX, tileY)`

Returns `TRUE` if the tile at grid position `(tileX, tileY)` has collision.
Collision data is loaded from a layer named `"collision"` in the Tiled map.
Any non-zero tile in that layer counts as solid.

- `tileX`, `tileY`: Zero-based tile grid coordinates.

```basic
; Convert world pixel position to tile grid coordinates
tile_x = INT(player_x# / 16)
tile_y = INT(player_y# / 16)

IF Tilemap.IsSolid(map, tile_x, tile_y + 2) THEN
    ; Player is standing on solid ground
    on_ground = TRUE
ENDIF
```

---

### `Tilemap.CollisionAt(handle, tileX, tileY)`

Returns the raw collision category value at the specified tile position (0 if
walkable). Use this for more complex collision filtering with bit masking.

---

### `Tilemap.SetCollision(handle, tileX, tileY, category)`

Sets the collision category of a tile at runtime. `0` = walkable.

---

## Tile Data

### `Tilemap.GetTile(handle, layerName$, tileX, tileY)`

Returns the tile GID (global ID) at a specific grid position on a named layer.
Returns `0` for empty tiles.

```basic
tile_id = Tilemap.GetTile(map, "ground", 5, 3)
```

---

### `Tilemap.SetTile(handle, layerName$, tileX, tileY, gid)`

Changes a tile in a layer at runtime. Set `gid = 0` to erase.

```basic
; Destroy a block when the player hits it
Tilemap.SetTile(map, "ground", tile_x, tile_y, 0)
```

---

## Map Information

### `Tilemap.Width(handle)` / `Tilemap.Height(handle)`

Returns the map dimensions in tiles.

### `Tilemap.LayerCount(handle)`

Returns the number of tile layers.

### `Tilemap.LayerName(handle, index)`

Returns the name of the layer at `index` (0-based).

### `Tilemap.SetTileSize(handle, width, height)`

Overrides the display pixel size of each tile. Useful for scaling up pixel-art
tiles without changing the `.tmx` file.

```basic
; Display 16×16 tiles as 32×32 on screen
Tilemap.SetTileSize(map, 32, 32)
```

---

## Full Example: Scrolling Platformer Map

```basic
Window.Open(960, 540, "Tilemap Demo")
Window.SetFPS(60)

map = Tilemap.Load("assets/maps/level1.tmx")

px# = 100
py# = 100
pvx# = 0
pvy# = 0
on_ground = 0
TILE_SIZE = 16

WHILE NOT Window.ShouldClose()
    dt# = Time.Delta()

    ; --- INPUT ---
    IF Input.KeyDown(KEY_A) THEN pvx# = pvx# - 400 * dt#
    IF Input.KeyDown(KEY_D) THEN pvx# = pvx# + 400 * dt#
    pvx# = pvx# * 0.85

    IF on_ground AND Input.KeyPressed(KEY_SPACE) THEN pvy# = -500

    ; --- PHYSICS ---
    pvy# = pvy# + 900 * dt#
    px# = px# + pvx# * dt#
    py# = py# + pvy# * dt#

    ; --- TILE COLLISION ---
    on_ground = 0

    ; Floor check (below feet)
    tx = INT(px# / TILE_SIZE)
    ty = INT((py# + 24) / TILE_SIZE)
    IF Tilemap.IsSolid(map, tx, ty) THEN
        py# = ty * TILE_SIZE - 24
        pvy# = 0
        on_ground = 1
    ENDIF

    ; Ceiling check (above head)
    ty_top = INT((py# - 16) / TILE_SIZE)
    IF Tilemap.IsSolid(map, tx, ty_top) THEN
        py# = (ty_top + 1) * TILE_SIZE + 16
        pvy# = 0
    ENDIF

    ; --- CAMERA ---
    cam_x = INT(px#) - 480
    cam_y = INT(py#) - 270

    ; --- DRAW ---
    Render.Clear(40, 60, 80)
    Tilemap.Draw(map, -cam_x, -cam_y)
    Draw.Rectangle(INT(px#) - cam_x - 8, INT(py#) - cam_y - 16, 16, 28, 255, 200, 80, 255)
    Render.Frame()
WEND

Tilemap.Free(map)
Window.Close()
```
