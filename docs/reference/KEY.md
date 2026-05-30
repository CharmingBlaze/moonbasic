# Key Commands

Supplementary keyboard query helpers. Prefer `INPUT.KEYDOWN` / `INPUT.KEYPRESSED` for new code — `KEY.*` provides aliases for specific patterns.

## Commands

### `KEY.DOWN(handle, keyCode)` 

Returns `TRUE` while `keyCode` is held, tested against an input context handle.

---

### `KEY.HIT(keyCode)` 

Returns `TRUE` on the first frame `keyCode` is pressed. Alias of `INPUT.KEYPRESSED`.

---

### `KEY.UP(keyCode)` 

Returns `TRUE` on the first frame `keyCode` is released. Alias of `INPUT.KEYRELEASED`.

---

## See also

- [INPUT.md](INPUT.md) — full input system
- [ACTION.md](ACTION.md) — named action bindings
- [GAME.md](GAME.md) — `GAME.KEYDOWN`, `GAME.KEYPRESSED`
