# 2D sprites and tilemaps

> Build top-down and side-view 2D games with **image sprites** and **Tiled** tile maps.

**Namespaces:** `SPRITE` · `TILEMAP` · `DRAW` · **Status:** Shipped

**Commands:** [COMMAND_REGISTRY.md#2d-world](../COMMAND_REGISTRY.md#2d-world) · [07-2D-WORLD.md](../07-2D-WORLD.md)

---

## Table of contents

- [At a glance](#at-a-glance)
- [Sprites vs tilemaps](#sprites-vs-tilemaps)
- [Sprite workflow](#sprite-workflow)
- [Tilemap workflow](#tilemap-workflow)
- [Collision in 2D levels](#collision-in-2d-levels)
- [Full example sketch](#full-example-sketch)
- [See also](#see-also)

---

## At a glance

| Tool | Best for |
|------|----------|
| **`SPRITE.*`** | Player, enemies, bullets, UI icons |
| **`TILEMAP.*`** | Large static levels from TMX |
| **`DRAW.RECTANGLE`** | Prototype without art files |

**Why tilemaps:** Drawing thousands of tiles as one `DRAW` call is faster than per-tile sprites.

---

## Sprites vs tilemaps

| Use sprites | Use tilemaps |
|-------------|--------------|
| Moving actors | Ground/wall layers |
| Rotating objects | Collision grid from TMX |
| Few large images | Repeating terrain |

---

## Sprite workflow

1. `SPRITE.LOAD(path)` or texture from `ASSET.TEXTURE`.
2. `SPRITE.SETPOSITION(s, x, y)` each frame.
3. `SPRITE.SETROTATION` / `SETSCALE` optional.
4. `SPRITE.DRAW(s)` inside 2D pass (after `RENDER.CLEAR`).
5. `SPRITE.FREE(s)` when done.

```basic
hero = SPRITE.LOAD("assets/hero.png")
WHILE NOT APP.SHOULDCLOSE()
    SPRITE.SETPOSITION(hero, px, py)
    RENDER.CLEAR(40, 44, 52)
    SPRITE.DRAW(hero)
    RENDER.FRAME()
WEND
```

**Why `SETPOSITION` every frame:** Sprites don’t auto-follow physics — you sync from your movement code or `BODY2D.X/Y`.

---

## Tilemap workflow

1. Design level in **Tiled** → export `.tmx`.
2. `map = TILEMAP.LOAD("levels/level1.tmx")`.
3. Each frame: `TILEMAP.DRAW(map)` after clear.
4. Gameplay: `TILEMAP.GETTILE(map, layer, tx, ty)` for walls.
5. `TILEMAP.FREE(map)` on level change.

```basic
map = TILEMAP.LOAD("levels/arena.tmx")
WHILE NOT APP.SHOULDCLOSE()
    RENDER.CLEAR(30, 35, 45)
    TILEMAP.DRAW(map)
    ; draw sprites on top
    RENDER.FRAME()
WEND
```

Sample: [examples/tilemap](../examples/tilemap/README.md).

---

## Collision in 2D levels

| Approach | Guide |
|----------|-------|
| Tile wall layer | Read `GETTILE` — solid if id ≠ 0 |
| Sprite overlap | [COLLISION-2D.md](COLLISION-2D.md) |
| Physics crates | [PHYSICS-2D-PLATFORMER.md](PHYSICS-2D-PLATFORMER.md) |

---

## Full example sketch

Top-down: tilemap floor + sprite hero + manual or `PHYSICS2D` movement.

`moonrun` · `moonbasic --check`

---

## See also

- [COLLISION-2D.md](COLLISION-2D.md)
- [examples/fps](../examples/fps/main.mb) — 2D without tilemap
