# Static Commands

`STATIC.*` commands are documented in [BODYREF.md](BODYREF.md).

## Quick Reference

| Command | Description |
|---|---|
| `STATIC.CREATE(shapeHandle)` | Create an immovable physics body from a shape |

Use `BODYREF.*` commands to set position, rotation, layer, and collision on the returned handle.

## Extended Command Reference

| Command | Description |
|--------|-------------|
| `STATIC.MAKE(shape)` | Deprecated alias of `STATIC.CREATE`. |

## See also

- [BODYREF.md](BODYREF.md) — full static / kinematic body docs
- [SHAPE.md](SHAPE.md) — shape handles
- [BODY3D.md](BODY3D.md) — dynamic bodies
