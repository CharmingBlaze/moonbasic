# Program / runtime

| Designed | moonBASIC | Notes |
|----------|------------|-------|
| **AppTitle(title)** | **`Window.SetTitle()`** | Updates the window title. |
| **SetFPS(fps)** | **`Window.SetFPS()`** | Sets the target frame rate. |
| **DeltaTime()** | **`Time.Delta()`** | Seconds since last frame. |
| **TimeMs()** | **`Time.TickCount()`** | Total runtime in milliseconds. |
| **Date()** | **`Date()`** | Formatted date string. |
| **Time()** | **`Time()`** | Formatted time string. |
| **Sleep(ms)** | **`Sleep()`** | Alias for `Wait()`. |
| **End()** | **`EndGame()`** | Terminates the VM and program. |

See [blitz-engine.md](blitz-engine.md) for all flat Blitz-style globals registered by the engine facade.
