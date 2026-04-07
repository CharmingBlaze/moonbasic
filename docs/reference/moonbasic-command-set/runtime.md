# Program / runtime

| Designed | Implementation | Memory / notes |
|----------|----------------|----------------|
| **AppTitle (title$)** | **Flat:** **`APPTITLE`** → **`WINDOW.SETTITLE`**; or title in **`WINDOW.OPEN(…, title$)`** | Window string only. |
| **SetFPS (fps)** | **Flat:** **`SETFPS`** → **`WINDOW.SETFPS`** | Caps frame rate. |
| **DeltaTime ()** | **Flat:** **`DELTATIME`** → **`TIME.DELTA`**; shortcut **`DT`** / **`GAME.DT`** | Float seconds. |
| **TimeMS ()** | **Flat:** **`TIMEMS`** → **`TICKCOUNT`** (ms); or monotonic wall **`TIME.GET`** | |
| **Sleep (ms)** | **`SLEEP`** / **`WAIT`** (core) | Int = ms, float = seconds. |
| **End ()** | **Flat:** **`FINISH`** → **`ENDGAME`** — `END` is a reserved keyword | **`ENDGAME`** stops VM; free handles with **`ERASE ALL`** if needed — [MEMORY.md](../../MEMORY.md). |

See [blitz-engine.md](blitz-engine.md) for all flat Blitz-style globals registered by the engine facade.
