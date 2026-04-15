# Kinematic Character Controller (KCC)

High-fidelity character movement with iterative "Move and Slide" collision resolution. On **Windows and Linux** desktop **`fullruntime`** builds with **CGO + Jolt**, KCC is **Jolt CharacterVirtual** end-to-end. **`!cgo`** or non-desktop builds use stubs (no real solver). **Documentation policy:** describe **Windows** before **Linux**; see [DEVELOPER.md](DEVELOPER.md#platform-priority-windows-then-linux).

## Core Workflow

1.  **Initialize**: Call `CHAR.CREATE(entity, radius#, height#)`.
2.  **Configuration**: Set step height with `CHAR.SETSTEP(entity, height#)`.
3.  **Movement**: Apply intentional movement using `CHAR.MOVEWITHCAMERA()` or `CHAR.MOVE()`.
4.  **Feedback**: Check grounding with `CHAR.ISGROUNDED()`.

```basic
player = CreateCapsule()
CHAR.CREATE(player, 0.5, 2.0)
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

### `CHARACTER.CREATE` (`Character.*`)

**Entity-bound only:** **`Character.Create(entity, radius#, height#)`** — same role as **`CHAR.CREATE`** (deprecated **`CHAR.MAKE`**): visuals and KCC stay in sync. Requires **desktop `fullruntime` with CGO + Jolt** (see [CHARACTER.md](reference/CHARACTER.md)).

### `CHAR.CREATE(entity, radius#, height#)`

Allocates a Kinematic Character Controller for the entity (**Jolt `CharacterVirtual`** on Windows and Linux when built with CGO + Jolt).

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

## Implementation Details (Jolt CharacterVirtual)

On **Windows and Linux** with **CGO + Jolt**, **`CHAR.*` / `PLAYER.*` / `CHARACTER.*`** drive **Jolt `CharacterVirtual`** (move-and-slide, stairs, slope limits) via the charcontroller bridge. The engine syncs entity transforms from the capsule each frame.

Stub builds (**`!cgo`** or non-desktop) do **not** run this solver; they surface a single **CGO + Jolt required** error for KCC commands (see [PLAYER.md](reference/PLAYER.md) / runtime player stubs).
