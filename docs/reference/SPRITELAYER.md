# SpriteLayer Commands

`SPRITELAYER.*` commands are documented in [SPRITEBATCH.md](SPRITEBATCH.md) under the **SpriteLayer** section.

## Quick Reference

| Command | Description |
|---|---|
| `SPRITELAYER.CREATE(zDepth)` | Create a depth-sorted layer |
| `SPRITELAYER.ADD(layer, sprite)` | Add a sprite |
| `SPRITELAYER.CLEAR(layer)` | Remove all |
| `SPRITELAYER.SETZ(layer, z)` | Update Z depth |
| `SPRITELAYER.DRAW(layer, x, y)` | Draw at offset |
| `SPRITELAYER.FREE(layer)` | Free handle |

## Extended Command Reference

| Command | Description |
|--------|-------------|
| `SPRITELAYER.MAKE(...)` | Deprecated alias of `SPRITELAYER.CREATE`. |

## See also

- [SPRITEBATCH.md](SPRITEBATCH.md) — full sprite batching and layering docs
