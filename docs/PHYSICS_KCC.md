# Kinematic Character Controller (KCC)

High-fidelity character movement with iterative "Move and Slide" collision resolution. This system provides professional-grade movement across **Windows** (host solver) and **Linux** (Jolt CharacterVirtual). **Documentation policy:** describe **Windows** before **Linux**; see [DEVELOPER.md](DEVELOPER.md#platform-priority-windows-then-linux).

## Core Workflow

1.  **Initialize**: Call `CHAR.MAKE(entity, radius#, height#)`.
2.  **Configuration**: Set step height with `CHAR.SETSTEP(entity, height#)`.
3.  **Movement**: Apply intentional movement using `CHAR.MOVEWITHCAMERA()` or `CHAR.MOVE()`.
4.  **Feedback**: Check grounding with `CHAR.ISGROUNDED()`.

```basic
player = CreateCapsule()
CHAR.MAKE(player, 0.5, 2.0)
CHAR.SETSTEP(player, 0.4)

WHILE NOT Window.ShouldClose()
    moveF = Axis(KEY_S, KEY_W)
    moveS = Axis(KEY_A, KEY_D)
    CHAR.MOVEWITHCAMERA(player, cam, moveF, moveS, 10.0)
    
    IF CHAR.ISGROUNDED(player) AND KeyHit(KEY_SPACE) THEN CHAR.JUMP(player, 9.0)
    
    ENTITY.UPDATE(Time.Delta())
    Render.Frame()
WEND
```

---

## Creation & Base Ops

### Polymorphic `CHARACTER.CREATE` (`Character.*`)

**Two shapes** (see [CHARACTER.md](reference/CHARACTER.md) for full detail and platform notes):

- **`Character.Create(x#, y#, z#)`** — **Standalone** kinematic character at world position; **host KCC** (e.g. Windows `fullruntime`) assigns a **virtual id** (negative integer, from a descending counter starting at **-1000**) so physics runs **without** a bound **EntityRef**.
- **`Character.Create(entity, radius#, height#)`** — **Entity-bound** capsule; same role as **`CHAR.MAKE`**: visuals and KCC stay in sync.

On **Linux + Jolt**, **`CHARACTER.CREATE`** is currently **entity-bound only**; use **`CHAR.MAKE`** / **`PLAYER.CREATE`** for Jolt if you spawn from coordinates by positioning an entity first.

### `CHAR.MAKE(entity, radius#, height#)`

Allocates a Kinematic Character Controller for the entity. On **Windows**, this enables the **host** iterative solver; on **Linux**, it allocates a Jolt **`CharacterVirtual`**.

- **`entity`**: Entity handle or id.
- **`radius#`**: Horizontal radius of the character capsule.
- **`height#`**: Total vertical height of the capsule.

### `PLAYER.TELEPORT(entity, x#, y#, z#)`
Instantly teleports the character to a world position and resets velocities.

---

## Movement

### `CHAR.MOVE(entity, dx#, dz#, speed#)`
Applies world-space horizontal movement.
- `dx#, dz#`: Normalized direction vector.
- `speed#`: Travel speed in units per second (scaled by `dt` internally).

### `CHAR.MOVEWITHCAMERA(entity, cam, fwd#, side#, speed#)`
Applies movement relative to a camera's orientation.
- `cam`: The camera handle.
- `fwd#, side#`: Input axes (typically from `AXIS()`).

---

## Interaction

### `CHAR.ISGROUNDED(entity [, coyote#])`
Returns TRUE if the character is supported by a floor surface.
- `coyote#`: Optional grace period in seconds after leaving a ledge (Batch 7+).

### `CHAR.JUMP(entity, force#)`
Applies an upward vertical impulse.

---

## Implementation Details (Move-and-Slide)

The MoonBASIC Host KCC implements an **Iterative Collision Solver** inspired by Godot's `CharacterBody3D`. 

1.  **Intent**: The character moves by `velocity * dt`.
2.  **Resolution**: The solver detects overlaps with static geometry using sphere-vs-AABB checks.
3.  **Sliding**: If a collision is found, the character is depenetrated, and the velocity vector is projected (slid) onto the surface plane to maintain momentum.
4.  **Iteration**: Steps 2-3 repeat up to 4 times per frame to handle complex corners and narrow corridors without jitter.
5.  **Snapping**: A predictive floor sensor snaps the character to the geometry when moving down slopes, preventing accidental "flight" on shallow descents.
