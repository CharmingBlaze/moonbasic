# Kinematic Commands

`KINEMATIC.*` and `KINEMATICREF.*` commands are documented in [BODYREF.md](BODYREF.md).

## Quick Reference

| Command | Description |
|---|---|
| `KINEMATIC.CREATE(shapeHandle)` | Create a kinematic body from a shape |
| `KINEMATICREF.SETVELOCITY(ref, vx, vy, vz)` | Set velocity for moving body |
| `KINEMATICREF.UPDATE(ref)` | Resolve movement and collisions |

Use `BODYREF.*` commands to set position, rotation, layer, and collision on the returned handle.

## Extended Command Reference

| Command | Description |
|--------|-------------|
| `KINEMATIC.MAKE(shape)` | Deprecated alias of `KINEMATIC.CREATE`. |

## See also

- [BODYREF.md](BODYREF.md) — full kinematic / static body docs
- [STATIC.md](STATIC.md) — static bodies
- [BODY3D.md](BODY3D.md) — dynamic bodies
