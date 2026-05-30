# EntityRef Commands

Handle-based helpers for entity physics state: grounding, jumping, and navigation updates.

Page shape follows [DOC_STYLE_GUIDE.md](../DOC_STYLE_GUIDE.md) (**WAVE pattern**).

## Core Workflow

`ENTITYREF.*` commands operate on an entity handle that has physics or navigation attached. Use them in your game loop to query grounding state, trigger jumps, and update navigation.

For the full entity API see [ENTITY.md](ENTITY.md). For character controllers see [CHARACTER.md](CHARACTER.md) and [CHARACTER_PHYSICS.md](CHARACTER_PHYSICS.md).

---

### `ENTITYREF.ISGROUNDED(entityHandle)` 

Returns `TRUE` if the entity is currently on the ground (physics probe succeeded).

---

### `ENTITYREF.JUMP(entityHandle, force)` 

Applies a vertical jump impulse to the entity.

- `entityHandle`: Entity with physics attached.
- `force`: Upward impulse strength.

---

### `ENTITYREF.NAVUPDATE(entityHandle)` 

Advances the entity's navigation agent one tick along its current path. Call each frame for entities using `ENTITY.NAVTO` or `ENT.NAVTO`.

---

## Full Example

This example checks grounding state and allows the player entity to jump.

```basic
player = ENTITY.LOAD("player.glb")
ENTITY.ADDPHYSICS(player)
ENTITY.NAVTO(player, 50.0, 0.0, 50.0)

WHILE NOT WINDOW.SHOULDCLOSE()
    ; Update navigation
    ENTITYREF.NAVUPDATE(player)

    ; Jump when grounded and space pressed
    IF ENTITYREF.ISGROUNDED(player) AND INPUT.KEYPRESSED(KEY_SPACE)
        ENTITYREF.JUMP(player, 8.0)
    END IF

    RENDER.BEGINFRAME()
    RENDER.BEGINMODE3D(cam)
    ENTITY.DRAW(player)
    RENDER.ENDMODE3D()
    RENDER.ENDFRAME()
WEND
```
