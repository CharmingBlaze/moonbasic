# SpriteGroup Commands

`SPRITEGROUP.*` commands are documented in [SPRITEBATCH.md](SPRITEBATCH.md) under the **SpriteGroup** section.

## Quick Reference

| Command | Description |
|---|---|
| `SPRITEGROUP.CREATE()` | Create a sprite group |
| `SPRITEGROUP.ADD(group, sprite)` | Add a sprite |
| `SPRITEGROUP.REMOVE(group, sprite)` | Remove a sprite |
| `SPRITEGROUP.CLEAR(group)` | Remove all sprites |
| `SPRITEGROUP.DRAW(group, x, y)` | Draw all at offset |
| `SPRITEGROUP.FREE(group)` | Free the handle |

## Extended Command Reference

| Command | Description |
|--------|-------------|
| `SPRITEGROUP.MAKE(...)` | Deprecated alias of `SPRITEGROUP.CREATE`. |

## See also

- [SPRITEBATCH.md](SPRITEBATCH.md) — full sprite batching and layering docs
