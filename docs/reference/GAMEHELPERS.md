# Game helpers — screen, timing, shortcuts

Convenience globals and aliases live in the **`game`** namespace when registered (see `runtime/mbgame` in full builds). Typical keys:

| Command | Role |
|---------|------|
| `SCREENW()` / `SCREENH()` | Backbuffer width/height |
| `SCREENCX()` / `SCREENCY()` | Screen center |
| `DT()` | Alias for **`TIME.DELTA()`** per frame |
| `MX()` / `MY()` | Mouse position in screen space |

Many programs use **`TIME.DELTA()`** directly (same value as **`DT()`** when the game helpers module is active).

---

## See also

- [TIMER.md](TIMER.md) — wall and simulated timers
- [PROGRAMMING.md](../PROGRAMMING.md) — main loop
