# Terrain and open worlds

> Large **heightfield** ground, sampling height for gameplay, and optional **chunk streaming** for big maps.

**Namespaces:** `TERRAIN` · `WORLD` · `CHUNK` · **Status:** Shipped · **Platform:** full runtime (mesh build)

**Commands:** [COMMAND_REGISTRY.md#2d-world](../COMMAND_REGISTRY.md#2d-world) · [reference/TERRAIN.md](../../reference/TERRAIN.md)

**Sample:** [examples/terrain_chase](../examples/terrain_chase/main.mb)

---

## Table of contents

- [At a glance](#at-a-glance)
- [When to use terrain](#when-to-use-terrain)
- [Simple heightfield](#simple-heightfield)
- [Height queries for gameplay](#height-queries-for-gameplay)
- [Streaming (large worlds)](#streaming-large-worlds)
- [Drawing terrain](#drawing-terrain)
- [Common mistakes](#common-mistakes)
- [See also](#see-also)

---

## At a glance

| Piece | Role |
|-------|------|
| **`TERRAIN.CREATE` / `LOAD`** | Height grid + mesh |
| **`TERRAIN.GETHEIGHT(x, z)`** | Snap player Y to ground |
| **`CHUNK.*` + `WORLD.*`** | Load/unload distant chunks |
| **`TERRAIN.DRAW`** | Render inside 3D pass |

**Why not one giant mesh:** Streaming keeps memory and frame cost bounded on large maps.

---

## When to use terrain

**Use when:**

- Outdoor 3D levels with hills.
- Need `GETHEIGHT` for placing trees/hero.

**Skip when:**

- Single flat plane — `ENTITY.CREATEPLANE` is enough.
- Pure 2D — use [SPRITES-TILEMAPS-2D.md](SPRITES-TILEMAPS-2D.md).

---

## Simple heightfield

```basic
terrain = TERRAIN.CREATE(128, 128)
TERRAIN.FILLPERLIN(terrain, ...)   ; see reference for args
TERRAIN.SETPOS(terrain, 0, 0, 0)
```

Or load heightmap image: `TERRAIN.LOAD("height.png")`.

---

## Height queries for gameplay

**Why:** Hero Y should match ground — not float or sink.

```basic
h = TERRAIN.GETHEIGHT(terrain, px, pz)
ENTITY.SETPOS(hero, px, h + 1, pz)
```

Aliases / helpers: `TERRAIN.SNAPY`, `ENTITY` terrain snap — see [reference/ENTITY.md](../../reference/ENTITY.md).

**Collision:** Terrain mesh may not be a Jolt body in simple samples — use height snap or add static physics separately ([COLLISION-3D.md](COLLISION-3D.md)).

---

## Streaming (large worlds)

**Why:** Only build GPU meshes near the player.

1. `CHUNK.SETRANGE(terrain, loadDist, unloadDist)`
2. Each frame: `WORLD.SETCENTER(x, z)` or follow entity
3. `WORLD.UPDATE()` — loads/unloads chunks

See [reference/WORLD.md](../../reference/WORLD.md) and terrain_chase sample.

---

## Drawing terrain

```basic
RENDER.BEGIN(cam)
TERRAIN.DRAW(terrain)
SCENE.DRAW()          ; props on terrain
RENDER.END()
```

Draw order: often **sky → terrain → entities** ([reference/TERRAIN.md](../../reference/TERRAIN.md)).

---

## Common mistakes

| Mistake | Fix |
|---------|-----|
| No height snap | Hero floats above hills |
| `DRAW` outside 3D pass | Wrong depth/state |
| Forget `WORLD.UPDATE` when streaming | Holes in world |
| Expect Jolt auto-collide | Add snap or physics |

**Run sample:**

```bash
moonrun examples/terrain_chase/main.mb
```

---

## See also

- [CAMERA-AND-INPUT.md](CAMERA-AND-INPUT.md) — orbit + WASD in terrain_chase
- [ENTITY-SYSTEM.md](ENTITY-SYSTEM.md)
