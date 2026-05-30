# Steer Commands

Steering behaviour forces for autonomous agent movement: seek, flee, arrive, wander, flock, obstacle avoidance, and path following. Returns force vectors that you apply to a physics body or position each frame.

## Core Workflow

1. Create a group with `STEER.GROUPMAKE()` if using flock behaviour.
2. Add agents to the group with `STEER.GROUPADD`.
3. Each frame: call the desired steer command (e.g. `STEER.SEEK`) — it returns a force handle.
4. Apply the force to your agent's physics body or position.
5. `STEER.GROUPCLEAR(group)` when restructuring.

---

## Groups (Flock)

### `STEER.GROUPMAKE()` 

Creates a flock group handle. Returns a **group handle**.

---

### `STEER.GROUPADD(group, agentHandle)` 

Adds an agent to a flock group. Used by `STEER.FLOCK` to compute cohesion/separation/alignment.

---

### `STEER.GROUPCLEAR(group)` 

Removes all agents from the group.

---

## Behaviours

### `STEER.SEEK(agent, tx, ty, tz)` 

Returns a force handle steering `agent` toward target `(tx, ty, tz)` at full speed.

---

### `STEER.FLEE(agent, tx, ty, tz)` 

Returns a force handle steering `agent` away from `(tx, ty, tz)`.

---

### `STEER.ARRIVE(agent, tx, ty, tz, slowRadius)` 

Like `SEEK` but decelerates within `slowRadius` world units of the target. Returns a force handle.

---

### `STEER.WANDER(agent, circleRadius, circleOffset)` 

Returns a force handle that produces smooth random wandering. `circleRadius` and `circleOffset` tune the wander circle in front of the agent.

---

### `STEER.FLOCK(agent, group, separationWeight, alignmentWeight, cohesionWeight)` 

Returns a combined force handle for flocking behaviour. Weights tune each of the three components (separation, alignment, cohesion).

---

### `STEER.AVOIDOBSTACLES(agent, avoidRadius)` 

Returns a force handle that steers `agent` away from nearby static obstacles within `avoidRadius`.

---

### `STEER.FOLLOWPATH(agent, pathHandle)` 

Returns a force handle steering `agent` along a pre-computed `pathHandle` (from `NAV.*` or manually built).

---

## Full Example

Three enemies flocking toward the player.

```basic
WINDOW.OPEN(960, 540, "Steer Demo")
WINDOW.SETFPS(60)

cam = CAMERA.CREATE()
CAMERA.SETPOS(cam, 0, 20, -20)
CAMERA.SETTARGET(cam, 0, 0, 0)

group = STEER.GROUPMAKE()

enemies = ARRAY.MAKE(3)
meshes  = ARRAY.MAKE(3)
FOR i = 0 TO 2
    e = NAVAGENT.CREATE(0)
    NAVAGENT.SETPOS(e, RNDF(-8, 8), 0, RNDF(-8, 8))
    NAVAGENT.SETSPEED(e, 4.0)
    STEER.GROUPADD(group, e)
    ARRAY.SET(enemies, i, e)
    m = ENTITY.CREATECUBE(0.8)
    ARRAY.SET(meshes, i, m)
NEXT i

px = 0.0
pz = 0.0

WHILE NOT WINDOW.SHOULDCLOSE()
    dt = TIME.DELTA()
    IF INPUT.KEYDOWN(KEY_RIGHT) THEN px = px + 5 * dt
    IF INPUT.KEYDOWN(KEY_LEFT)  THEN px = px - 5 * dt
    IF INPUT.KEYDOWN(KEY_DOWN)  THEN pz = pz + 5 * dt
    IF INPUT.KEYDOWN(KEY_UP)    THEN pz = pz - 5 * dt

    FOR i = 0 TO 2
        e = ARRAY.GET(enemies, i)
        f = STEER.ARRIVE(e, px, 0, pz, 2.0)
        ; apply force as velocity nudge
        vx = VEC3.X(f) * dt
        vz = VEC3.Z(f) * dt
        ex = NAVAGENT.X(e) + vx
        ez = NAVAGENT.Z(e) + vz
        NAVAGENT.SETPOS(e, ex, 0, ez)
        ENTITY.SETPOS(ARRAY.GET(meshes, i), ex, 0, ez)
        VEC3.FREE(f)
    NEXT i

    ENTITY.UPDATE(dt)
    RENDER.CLEAR(20, 25, 35)
    RENDER.BEGIN3D(cam)
        ENTITY.DRAWALL()
        DRAW3D.SPHERE(px, 0, pz, 0.5, 80, 200, 80, 255)
        DRAW3D.GRID(20, 1.0)
    RENDER.END3D()
    RENDER.FRAME()
WEND

WINDOW.CLOSE()
```

---

## See also

- [NAVAGENT.md](NAVAGENT.md) — navigation agents with path following
- [BTREE.md](BTREE.md) — behaviour tree AI logic
- [PHYSICS3D.md](PHYSICS3D.md) — apply forces to rigid bodies
