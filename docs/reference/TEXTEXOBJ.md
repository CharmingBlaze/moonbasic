# TextExObj Commands

`TEXTEXOBJ.*` commands are documented in [TEXTDRAW.md](TEXTDRAW.md) under the **TEXTEXOBJ Commands** section.

Custom-font retained text objects with float position, size, and spacing. Requires a `FONT.LOAD` handle.

## Quick Reference

| Command | Description |
|---|---|
| `TEXTEXOBJ.POS(handle, x, y)` | Float screen position |
| `TEXTEXOBJ.SIZE(handle, size)` | Float font size |
| `TEXTEXOBJ.SPACING(handle, sp)` | Character spacing |
| `TEXTEXOBJ.COLOR(handle, r, g, b, a)` | Text color |
| `TEXTEXOBJ.SETTEXT(handle, text)` | Update string |
| `TEXTEXOBJ.DRAW(handle)` | Draw this frame |
| `TEXTEXOBJ.FREE(handle)` | Free handle |

## See also

- [TEXTDRAW.md](TEXTDRAW.md) — both text draw object types
- [FONT.md](FONT.md) — font loading
