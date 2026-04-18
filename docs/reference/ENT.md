# Ent Commands

High-level entity helper shortcuts for gameplay: health/damage, team assignment, nearest-enemy queries, shooting, tween animations, and nav integration. Designed for rapid game prototyping.

For the full entity system see [ENTITY.md](ENTITY.md).

## Core Workflow

1. Assign teams and HP: `ENT.SETTEAM(e, team)`, `ENT.SETHP(e, hp, maxHp)`.
2. Deal damage: `ENT.DAMAGE(e, amount)`.
3. React to death: `ENT.ONDEATH(e, handler)`.
4. Query: `ENT.GETNEAREST(e, radius, filter)`, `ENT.DIST(a, b)`.

---

## Health & Teams

### `ENT.SETHP(entityId, hp, maxHp)` 

Sets the current and maximum HP for an entity.

---

### `ENT.DAMAGE(entityId, amount)` 

Reduces the entity's HP by `amount`. If HP reaches zero, fires the death callback.

---

### `ENT.SETTEAM(entityId, team)` / `ENT.SET_TEAM(entityId, team)` 

Assigns the entity to a team integer (used by `SHOOT` targeting).

---

### `ENT.ONDEATH(entityId, callbackName)` / `ENT.ONDEATH(entityId, callbackId)` 

Registers a function name or id to call when this entity's HP reaches zero.

---

### `ENT.ISALIVE(entityHandle)` 

Returns `TRUE` if the entity has HP > 0.

---

### `ENT.SET_HP(entityHandle, hp)` 

Sets the HP directly via handle.

---

## Query

### `ENT.GETNEAREST(entityId, radius, filterTag)` / `ENT.GET_NEAREST(...)` 

Returns a handle to the nearest entity within `radius` matching `filterTag`. Returns `0` if none found.

---

### `ENT.DIST(entityIdA, entityIdB)` 

Returns the world-space distance between two entities.

---

## Combat

### `ENT.SHOOT(entityId, targetId, speed)` 

Spawns a projectile from `entityId` toward `targetId` at `speed`. Returns the projectile entity id.

---

## Visual Helpers

### `ENT.FADE(entityId, targetAlpha, duration)` 

Fades the entity's alpha to `targetAlpha` over `duration` seconds.

---

### `ENT.WOBBLE(entityId, amount, speed)` 

Applies a sine-wave wobble scale effect.

---

### `ENT.TWEEN(entityId, x, y, z, duration)` 

Smoothly moves the entity to world position `(x, y, z)` over `duration` seconds.

---

## Navigation

### `ENT.NAVTO(entityHandle, x, y, z)` 

Moves the entity to the nav target position using the attached nav agent.

---

## Full Example

Simple arena shooter with health, teams, and nearest-enemy targeting.

```basic
WINDOW.OPEN(960, 540, "Ent Demo")
WINDOW.SETFPS(60)

cam = CAMERA.CREATE()
CAMERA.SETPOS(cam, 0, 16, -16)
CAMERA.SETTARGET(cam, 0, 0, 0)

player = ENTITY.CREATECUBE(1.0)
ENTITY.SETPOS(player, 0, 0.5, 0)
ENT.SETTEAM(player, 1)
ENT.SETHP(player, 100, 100)

FOR i = 1 TO 4
    e = ENTITY.CREATECUBE(1.0)
    ENTITY.SETCOLOR(e, 255, 60, 60)
    ENTITY.SETPOS(e, RNDF(-8, 8), 0.5, RNDF(-8, 8))
    ENT.SETTEAM(e, 2)
    ENT.SETHP(e, 30, 30)
NEXT i

WHILE NOT WINDOW.SHOULDCLOSE()
    dt = TIME.DELTA()

    IF INPUT.KEYPRESSED(KEY_F) THEN
        nearest = ENT.GETNEAREST(player, 10.0, "")
        IF nearest THEN
            ENT.DAMAGE(nearest, 15)
            PRINT "Damaged! HP left: unknown"
        END IF
    END IF

    ENTITY.UPDATE(dt)

    RENDER.CLEAR(20, 25, 35)
    RENDER.BEGIN3D(cam)
        ENTITY.DRAWALL()
        DRAW.GRID(20, 1.0)
    RENDER.END3D()
    DRAW.TEXT("F = attack nearest", 10, 10, 18, 200, 200, 200, 255)
    RENDER.FRAME()
WEND

WINDOW.CLOSE()
```

---

## See also

- [ENTITY.md](ENTITY.md) — full entity system
- [NAVAGENT.md](NAVAGENT.md) — navigation for entities
- [BTREE.md](BTREE.md) — AI behaviour using `ENT.*` queries
