# KinematicRef Commands

`KINEMATICREF.*` commands are documented in [BODYREF.md](BODYREF.md) under the **KINEMATICREF Commands** section.

## Quick Reference

| Command | Description |
|---|---|
| `KINEMATICREF.SETVELOCITY(ref, vx, vy, vz)` | Set velocity for a moving kinematic body |
| `KINEMATICREF.UPDATE(ref)` | Resolve movement and collisions this frame |

Bodies are created with `KINEMATIC.CREATE(shapeHandle)` — see [KINEMATIC.md](KINEMATIC.md).

## See also

- [BODYREF.md](BODYREF.md) — full kinematic / static body docs
- [KINEMATIC.md](KINEMATIC.md) — `KINEMATIC.CREATE`
- [SHAPE.md](SHAPE.md) — shape handles
