# Scatter, Prop & Spawner Commands

Batch placement of decorative and gameplay objects across the world.

Page shape follows [DOC_STYLE_GUIDE.md](../DOC_STYLE_GUIDE.md) (**WAVE pattern**).

## Core Workflow

**Scatter** — Create a scatter set from a model, apply it to populate the world, draw all instances each frame, and free when done.

**Prop** — Place individual static props at specific positions, draw all each frame, and free individually.

**Spawner** — Create a spawner that produces entity instances at runtime.

---

### `SCATTER.CREATE(modelHandle)` 

Creates a scatter set for the given model. Returns a scatter handle.

- `modelHandle`: The model to scatter instances of.

---

### `SCATTER.APPLY(scatterHandle)` 

Applies the scatter set, populating the world with instances based on the set's configuration.

---

### `SCATTER.DRAWALL()` 

Draws all active scatter instances. Call once per frame during the render pass.

---

### `SCATTER.FREE(scatterHandle)` 

Frees the scatter set and all its instances.

---

### `PROP.PLACE(model, x, y, z)` 

Places a static prop at the given world position. Returns a prop handle.

- `model`: Model handle or asset reference.
- `x`, `y`, `z`: World position.

---

### `PROP.DRAWALL()` 

Draws all placed props. Call once per frame during the render pass.

---

### `PROP.FREE(propHandle)` 

Removes and frees a placed prop.

---

### `SPAWNER.MAKE(prefab, x, y, z)` 

Creates a spawner at the given position that produces instances of the given prefab. Returns a spawner handle.

- `prefab`: Entity prefab or model handle.
- `x`, `y`, `z`: Spawn origin.

---

## Full Example

This example scatters trees across a terrain and places a few props manually.

```basic
tree_model = MODEL.LOAD("tree.glb")

; Scatter trees
trees = SCATTER.CREATE(tree_model)
SCATTER.APPLY(trees)

; Place a few static props
rock = MODEL.LOAD("rock.glb")
r1 = PROP.PLACE(rock, 10.0, 0.0, 5.0)
r2 = PROP.PLACE(rock, -3.0, 0.0, 12.0)

; Create a spawner for enemies
enemy_prefab = MODEL.LOAD("enemy.glb")
spawner = SPAWNER.MAKE(enemy_prefab, 0.0, 1.0, 0.0)

WHILE NOT WINDOW.SHOULDCLOSE()
    RENDER.BEGINFRAME()
    RENDER.BEGINMODE3D(cam)
    SCATTER.DRAWALL()
    PROP.DRAWALL()
    RENDER.ENDMODE3D()
    RENDER.ENDFRAME()
WEND

; Cleanup
PROP.FREE(r1)
PROP.FREE(r2)
SCATTER.FREE(trees)
```
