# Gameplay Intent API Reference

The **Gameplay Intent API** is a high-level command manifest designed to bridge the gap between raw physics and professional game feel. It provides a Kinematic Character Controller (KCC), autonomous navigation, and simplified combat logic that works identically on Windows and Linux.

## 1. Kinematic Character Controller (CHAR.*)

Replaces standard "Rigid Body" physics for players and humanoid NPCs. These commands use collision-aware sweeps instead of forces, eliminating jitter and "bunny hopping."

**Naming:** **`CHAR.CREATE`** is canonical; **`CHAR.MAKE`** is a deprecated alias (same handler as **`PLAYER.CREATE`**).

| Command | Role | Description |
| :--- | :--- | :--- |
| `CHAR.CREATE(e, r, h)` | Setup | Initializes the KCC on entity `e` with radius `r` and height `h`. Disables scripted gravity. |
| `CHAR.SETSTEP(e, h)` | Stairs | Sets the maximum height (world units) the character can automatically step over. |
| `CHAR.SETSLOPE(e, deg)` | Slopes | Prevents the character from climbing surfaces steeper than `deg`. |
| `CHAR.STICK(e, dist)` | Glue | Keeps the character "glued" to the floor when moving down slopes within `dist`. |
| `CHAR.JUMP(e, pwr)` | Jump | Applies an instant upward velocity peak to the character. |
| `CHAR.MOVE(e, dx, dz, spd)`| Move | Direct XZ movement in world space at speed `spd`. |
| `CHAR.MOVEWITHCAM(e, c, f, s, spd)` | Controls | Moves character relative to camera `c` orientation using forward/strafe axes. |
| `CHAR.ISGROUNDED(e)` | State | Returns true if the character is currently touching a floor surface. |

## 2. Navigation Intent (NAV.*)

Simplifies AI movement. These commands interact with the KCC to handle obstacle avoidance and pathing at a high level.

| Command | Usage | Description |
| :--- | :--- | :--- |
| `NAV.GOTO(e, x, z, spd)` | Pathing | Move entity `e` to world XZ coordinates at speed `spd`. Includes soft-stop damping. |
| `NAV.CHASE(e, target, gap, spd)` | IA | Follow `target` entity until within `gap` world units. |
| `NAV.PATROL(e, ax, az, bx, bz, spd)` | Logic | Ping-pong between two XZ points indefinitely. |
| `NAV.UPDATE(e)` | Tick | Advances navigation state (called automatically by `UPDATEPHYSICS` but can be forced). |

## 3. Combat & Logic (ENT.*)

Streamlined entity management for RPG and Action games.

| Command | Description |
| :--- | :--- |
| `ENT.SETHP(e, cur, max)`| Initialize health state on entity. |
| `ENT.DAMAGE(e, amount)` | Reduce HP; triggers squashing and optional red-flash effects. |
| `ENT.SETTEAM(e, id)` | Set numeric team for AI targeting / friendly-fire logic. |
| `ENT.GETNEAREST(e, radius, tag)` | Find the closest entity within `radius` that has a specific tag. |
| `ENT.TWEEN(e, x, y, z, dur)` | Smoothly interpolate an entity to a target position. |
| `ENT.FADE(e, alpha, dur)` | Smoothly interpolate entity transparency (cinematic deaths/spawns). |
| `ENT.WOBBLE(e, pwr, spd)` | Add a continuous floating/bobbing effect (pickups/targets). |

## 4. World Interaction (WORLD.*)

High-level helpers for camera and environment targeting.

| Command | Role | Description |
| :--- | :--- | :--- |
| `WORLD.MOUSEFLOOR(cam, y)` | Targeting | Returns [x, y, z] where mouse ray hits plane at height `y`. |
| `WORLD.MOUSEPICK(cam)` | Selection | Returns the entity ID under the mouse cursor using physics raycast. |
| `WORLD.TOSCREEN(e, cam)` | UI | Returns [x, y] screen coordinates for an entity's 3D position. |
| `WORLD.HITSTOP(dur)` | Feedback | Pauses simulation for `dur` seconds to emphasize hits/impacts. |
| `WORLD.SHAKE(pwr, dur)` | Feedback | Shakes the active camera with intensity `pwr`. |

---
> [!TIP]
> Always call `UPDATEPHYSICS()` once per frame in your main loop to advance the KCC and Navigation simulations.
