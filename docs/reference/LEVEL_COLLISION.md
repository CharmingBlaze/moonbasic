# Level Collision Commands

Commands for managing environment physics and automated static mesh collisions.

## Core Workflow

1. **Load**: Load your level model using `LEVEL.LOAD()`.
2. **Mark**: (Optional) Use `ENTITY.SETSTATIC()` to mark specific props as part of the environment.
3. **Cook**: Call `LEVEL.AUTOCOLLIDE()` to batch-create high-performance collisions for all static objects.

```basic
; 1. Setup the world
LEVEL.SETUP(-28)

; 2. Load the main level
castle = LEVEL.LOAD("models/castle.glb")
ENTITY.SETSTATIC(castle, TRUE)

; 3. Load extra props
gate = MODEL.LOAD("models/gate.glb")
ENTITY.SETSTATIC(gate, TRUE)

; 4. Bake all collisions in one go
LEVEL.AUTOCOLLIDE()
```

---

## Static Collision

### `LEVEL.STATIC(entity)` 
Generates a persistent static mesh collision body for a specific entity.
- `entity`: The ID or handle of the entity to process.
- **Workflow**: This command is best for dynamic environment pieces that are added after the main level load.

---

### `LEVEL.AUTOCOLLIDE()` 
The **Easy Mode** way to setup a world. It scans every active entity and creates collisions for anything marked as static that has a valid 3D model.
- **Performance**: This command is extremely efficient and should be called once the initial level loading is finished.
- **Returns**: The total count of entities that were successfully "cooked" into the physics world.

---

## Handle Methods (Easy Mode)

You can also trigger these operations directly on entity handles:

### `ENTITY.SETSTATIC(toggle)` 
Marks or unmarks an entity as "Static" for the next `AUTOCOLLIDE` pass.

---

### `ENTITY.SETCOLLISIONMESH()` 
Immediately bakes a static mesh collision for this entity. Equivalent to `LEVEL.STATIC(entity)`.

```basic
hero = MODEL.LOAD("castle.glb")
hero.SetCollisionMesh()
```

---

## Platform Notes

| Feature | Windows + Linux (CGO + Jolt) | Stub / no Jolt |
|---------|-------------------------------|----------------|
| Automated Baking | Jolt collision pipeline where enabled | Limited / errors |
| Performance | Native Jolt | N/A |
| Stability | Same code path on both desktop OSes | Compile-safe stubs |

---

## Consistency Check
Signatures match `compiler/builtinmanifest/commands.json`.
- `LEVEL.STATIC(any)`
- `LEVEL.AUTOCOLLIDE()`
- `ENTITY.SETSTATIC(any, any)`

---

## Full Example

Loading a level and baking static collision in one pass.

```basic
WINDOW.OPEN(960, 540, "Level Collision Demo")
WINDOW.SETFPS(60)

PHYSICS3D.START()
PHYSICS3D.SETGRAVITY(0, -10, 0)

LEVEL.SETUP(-10)
castle = LEVEL.LOAD("assets/castle.glb")
ENTITY.SETSTATIC(castle, TRUE)

barrel = MODEL.LOAD("assets/barrel.glb")
ENTITY.SETSTATIC(barrel, TRUE)
ENTITY.SETPOS(barrel, 4, 0, 2)

LEVEL.AUTOCOLLIDE()

cam = CAMERA.CREATE()
CAMERA.SETPOS(cam, 0, 8, -14)
CAMERA.SETTARGET(cam, 0, 2, 0)

WHILE NOT WINDOW.SHOULDCLOSE()
    PHYSICS3D.UPDATE()
    ENTITY.UPDATE(TIME.DELTA())

    RENDER.CLEAR(80, 100, 130)
    RENDER.BEGIN3D(cam)
        ENTITY.DRAWALL()
        DRAW.GRID(20, 1.0)
    RENDER.END3D()
    RENDER.FRAME()
WEND

PHYSICS3D.STOP()
WINDOW.CLOSE()
```

---

## See also

- [LEVEL.md](LEVEL.md) — `LEVEL.LOAD`, `LEVEL.SETUP`, `LEVEL.BINDSCRIPT`
- [PHYSICS3D.md](PHYSICS3D.md) — physics world setup
- [ENTITY.md](ENTITY.md) — `ENTITY.SETSTATIC`
