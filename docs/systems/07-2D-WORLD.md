# 2D and world systems: SPRITE, TILEMAP, TERRAIN, PARTICLE, ANIMATION

> 2D sprites, tile maps, heightfield terrain, particle effects, and entity/model animation.

**All commands:** [COMMAND_REGISTRY.md#2d-world](COMMAND_REGISTRY.md#2d-world)

**Deep guides:** [guides/SPRITES-TILEMAPS-2D.md](guides/SPRITES-TILEMAPS-2D.md) · [guides/TERRAIN-OPEN-WORLD.md](guides/TERRAIN-OPEN-WORLD.md) · [guides/PARTICLES.md](guides/PARTICLES.md) · [guides/ANIMATION.md](guides/ANIMATION.md)

**See also:** [03-ASSETS](03-ASSETS.md) · [reference/SPRITE.md](../reference/SPRITE.md) · [reference/TERRAIN.md](../reference/TERRAIN.md)

---

## Table of contents

- [SPRITE system](#sprite-system)
- [TILEMAP system](#tilemap-system)
- [TERRAIN system](#terrain-system)
- [PARTICLE system](#particle-system)
- [ANIMATION system](#animation-system)
- [Full example](#full-example)
- [Memory notes](#memory-notes)
- [See also](#see-also)

---

## SPRITE system

2D images in world or screen space.

### Core workflow

1. `SPRITE.LOAD(path)` or texture from `TEXTURE.LOAD` / `ASSET.TEXTURE`.
2. `SPRITE.SETPOSITION`, `SETROTATION`, `SETSCALE`.
3. `SPRITE.DRAW(sprite)` each frame (2D pass).

---

### `SPRITE.LOAD(path)` / `SPRITE.CREATE(texture)`

Creates a sprite from file or texture handle.

**Returns:** `handle`

**Example:**

```basic
icon = SPRITE.LOAD("assets/icon.png")
SPRITE.SETPOSITION(icon, 100, 200)
```

---

### Transform and draw

| Command | Description |
|---------|-------------|
| `SPRITE.SETPOSITION(s, x, y)` | Screen or world 2D position |
| `SPRITE.SETROTATION(s, degrees)` | Rotation |
| `SPRITE.SETSCALE(s, sx, sy)` | Scale |
| `SPRITE.DRAW(s)` | Draw this frame |
| `SPRITE.FREE(s)` | Release handle |

---

## TILEMAP system

Tiled map loading and tile layers.

### `TILEMAP.LOAD(path)`

Loads a TMX (or supported) tile map.

**Returns:** `handle`

**Example:**

```basic
map = TILEMAP.LOAD("levels/level1.tmx")
```

---

### `TILEMAP.DRAW(map)` / `TILEMAP.GETTILE` / `TILEMAP.SETTILE`

| Command | Description |
|---------|-------------|
| `TILEMAP.DRAW(map)` | Render visible layers |
| `TILEMAP.GETTILE(map, layer, x, y)` | Read tile id |
| `TILEMAP.SETTILE(map, layer, x, y, id)` | Write tile id |

See [examples/tilemap](../examples/tilemap/README.md).

---

## TERRAIN system

Heightfield terrain for 3D worlds.

### Core workflow

1. `TERRAIN.CREATE(width, depth)` or `TERRAIN.LOAD(path)`.
2. `TERRAIN.GETHEIGHT(terrain, x, z)` for gameplay.
3. `TERRAIN.DRAW(terrain)` inside `RENDER.BEGIN`.

---

### Key commands

| Command | Description |
|---------|-------------|
| `TERRAIN.CREATE(w, d)` | Empty height grid |
| `TERRAIN.LOAD(path)` | Load heightmap / terrain file |
| `TERRAIN.SETTEXTURE(terrain, tex)` | Ground texture |
| `TERRAIN.SETHEIGHT(terrain, x, z, h)` | Edit height |
| `TERRAIN.GETHEIGHT(terrain, x, z)` | Sample height |
| `TERRAIN.DRAW(terrain)` | Draw mesh |

**Example:**

```basic
terrain = TERRAIN.CREATE(128, 128)
h = TERRAIN.GETHEIGHT(terrain, playerX, playerZ)
```

See [examples/terrain_chase](../examples/terrain_chase/README.md).

---

## PARTICLE system

Emitters for fire, smoke, sparks.

### Core workflow

1. `PARTICLE.CREATE()` — emitter handle.
2. `PARTICLE.SETTEXTURE`, `SETRATE`, `SETLIFETIME`, `SETSPEED`, colors.
3. `PARTICLE.PLAY(fx)` / `PARTICLE.STOP(fx)`; `PARTICLE.UPDATE` / `DRAW` each frame.

---

### Key commands

| Command | Description |
|---------|-------------|
| `PARTICLE.CREATE()` | New emitter |
| `PARTICLE.SETTEXTURE(fx, tex)` | Particle image |
| `PARTICLE.SETRATE(fx, n)` | Spawns per second |
| `PARTICLE.SETLIFETIME(fx, seconds)` | Particle life |
| `PARTICLE.SETSPEED(fx, min, max)` | Speed range |
| `PARTICLE.PLAY(fx)` / `STOP(fx)` | Start / stop |
| `PARTICLE.UPDATE(fx)` | Simulate |
| `PARTICLE.DRAW(fx)` | Render |
| `PARTICLE.FREE(fx)` | Release |

**Example:**

```basic
fx = PARTICLE.CREATE()
PARTICLE.SETRATE(fx, 50)
PARTICLE.PLAY(fx)
```

---

## ANIMATION system

**Status:** Dual API — skeletal clips on entities and FSM-style `ANIM.*` helpers.

### Entity / model animation

| Command | Description |
|---------|-------------|
| `ENTITY.PLAYANIM(ent, name)` | Play clip on model |
| `ENTITY.STOPANIM(ent)` | Stop clip |
| `MODEL.ANIMCOUNT` / `GETANIMNAME` | List clips |

**Aliases:** checklist `ANIM.PLAY` on entity → `ENTITY.PLAYANIM`

---

### FSM animation (`ANIM.*`)

| Command | Description |
|---------|-------------|
| `ANIM.DEFINE(stateMachine, …)` | Define states |
| `ANIM.UPDATE(sm)` | Tick per frame |

See [reference/ANIM.md](../reference/ANIM.md) for state machine details.

---

## Full example

```basic
; 2D sprite overlay on a 3D loop
APP.OPEN(800, 600, "2D + Particles")
APP.SETFPS(60)

cam = CAMERA.CREATE()
CAMERA.SETACTIVE(cam)
CAMERA.SETPOS(cam, 0, 5, -12)

fx = PARTICLE.CREATE()
PARTICLE.SETRATE(fx, 40)
PARTICLE.PLAY(fx)

WHILE NOT APP.SHOULDCLOSE()
    PARTICLE.UPDATE(fx)
    RENDER.CLEAR(10, 12, 18)
    RENDER.BEGIN(cam)
    SCENE.DRAW()
    PARTICLE.DRAW(fx)
    RENDER.END()
    RENDER.FRAME()
WEND

PARTICLE.FREE(fx)
APP.CLOSE()
```

---

## Memory notes

- `SPRITE.FREE`, `TILEMAP.FREE`, `TERRAIN.FREE`, `PARTICLE.FREE` on level unload.
- Particle emitters stop allocating when `STOP` is called; `FREE` releases GPU buffers.

---

## See also

- [03-ASSETS](03-ASSETS.md) — textures for sprites and particles
- [examples/platformer](../examples/platformer/main.mb) — 2D gameplay
