# Program / runtime

**Conventions:** [STYLE_GUIDE.md](../../../STYLE_GUIDE.md), [API_CONVENTIONS.md](../API_CONVENTIONS.md).

| Designed | Registry (use in new examples) | Notes |
|----------|-------------------------------|-------|
| **AppTitle(title)** | **`WINDOW.SETTITLE(title)`** | |
| **SetFPS(fps)** | **`WINDOW.SETFPS(fps)`** | |
| **DeltaTime()** | **`TIME.DELTA()`** | Seconds since last frame (or **`DT()`**). |
| **TimeMs()** | **`TICKCOUNT()`** | Milliseconds since start (see [TIME.md](../TIME.md)). |
| **Date()** | **`DATE()`** | Wall-clock helpers vary — see manifest. |
| **Time()** | **`TIME()`** | |
| **Sleep(ms)** | **`SLEEP`** / **`WAIT`** | Per your build’s aliases. |
| **End()** | **`ENDGAME`** | Stops the VM where registered. |

See [blitz-engine.md](blitz-engine.md) for flat Blitz-style globals on the engine facade.
