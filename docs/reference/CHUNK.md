# Chunk Commands

Procedural world chunking: divide an infinite world into loading cells, track which chunks are loaded, and generate terrain per chunk.

## Core Workflow

1. Create a chunk manager and configure `CHUNK.SETRANGE`.
2. Each frame: call `CHUNK.GENERATE(manager, cx, cz)` for cells in range — it returns a chunk handle only for newly entered cells.
3. `CHUNK.ISLOADED(manager, cx, cz)` to avoid re-generating.
4. `CHUNK.COUNT(manager)` to monitor active chunk count.

---

## Commands

### `CHUNK.GENERATE(manager, cx, cz)` 

Attempts to generate the chunk at grid coordinates `(cx, cz)`. Returns a **chunk handle** if this chunk was newly created, `0` if already loaded. Use the returned handle to spawn terrain/entities for that cell.

---

### `CHUNK.COUNT(manager)` 

Returns the number of currently loaded chunks.

---

### `CHUNK.SETRANGE(manager, loadRadius, unloadRadius)` 

Sets the load/unload radius in chunk units. Chunks beyond `unloadRadius` from the player are freed. Returns the manager handle for chaining.

---

### `CHUNK.ISLOADED(manager, cx, cz)` 

Returns `TRUE` if the chunk at `(cx, cz)` is currently loaded.

---

## Full Example

Infinite flat terrain with dynamically loaded 32×32 unit chunks.

```basic
WINDOW.OPEN(960, 540, "Chunk Demo")
WINDOW.SETFPS(60)

PHYSICS3D.START()

cam = CAMERA.CREATE()
CAMERA.SETPOS(cam, 0, 20, 0)
CAMERA.SETTARGET(cam, 0, 0, 0)

CHUNK_SIZE = 32
manager = 1    ; placeholder manager id (see NAVMESH/CHUNK system)
CHUNK.SETRANGE(manager, 3, 5)

px = 0.0
pz = 0.0

WHILE NOT WINDOW.SHOULDCLOSE()
    dt = TIME.DELTA()
    IF INPUT.KEYDOWN(KEY_RIGHT) THEN px = px + 10 * dt
    IF INPUT.KEYDOWN(KEY_LEFT)  THEN px = px - 10 * dt
    IF INPUT.KEYDOWN(KEY_DOWN)  THEN pz = pz + 10 * dt
    IF INPUT.KEYDOWN(KEY_UP)    THEN pz = pz - 10 * dt

    ; player chunk coordinates
    pcx = INT(px / CHUNK_SIZE)
    pcz = INT(pz / CHUNK_SIZE)

    FOR cx = pcx - 3 TO pcx + 3
        FOR cz = pcz - 3 TO pcz + 3
            IF NOT CHUNK.ISLOADED(manager, cx, cz)
                ch = CHUNK.GENERATE(manager, cx, cz)
                IF ch THEN
                    ; spawn ground plane for this chunk
                    e = ENTITY.CREATEPLANE(CHUNK_SIZE, CHUNK_SIZE)
                    ENTITY.SETPOS(e, cx * CHUNK_SIZE, 0, cz * CHUNK_SIZE)
                END IF
            END IF
        NEXT cz
    NEXT cx

    CAMERA.SETPOS(cam, px, 20, pz - 5)
    CAMERA.SETTARGET(cam, px, 0, pz)

    RENDER.CLEAR(80, 120, 160)
    RENDER.BEGIN3D(cam)
        ENTITY.DRAWALL()
    RENDER.END3D()
    DRAW.TEXT("Loaded: " + STR(CHUNK.COUNT(manager)), 10, 10, 18, 255, 255, 255, 255)
    RENDER.FRAME()
WEND

PHYSICS3D.STOP()
WINDOW.CLOSE()
```

---

## See also

- [ENTITY.md](ENTITY.md) — spawning entities per chunk
- [NOISE.md](NOISE.md) — procedural terrain height for each chunk
- [NAVMESH.md](NAVMESH.md) — navmesh with streaming chunks
