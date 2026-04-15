# Timers — wall clock vs simulation (`TIMER.*`)

Implemented in **`runtime/mbgame/timers.go`** and **`runtime/mbgame/timer_sim.go`**.

## Wall-clock timers (`TIMER.NEW`)

- **`TIMER.NEW(duration)`** — returns a handle; **`duration`** is seconds from **wall-clock** **`time.Now`**.
- **`TIMER.RESET(timer, duration)`** — reschedules the deadline.
- **`TIMER.REMAINING(timer)`** — seconds until the deadline (works for **both** wall and sim handles; see below).
- **`TIMER.FINISHED(timer)`** — **`TRUE`** after the deadline (wall timers only; use **`TIMER.DONE`** for sim).
- **`TIMER.FREE(timer)`** — frees the handle.

## Simulation timers (`TIMER.CREATE` / deprecated `TIMER.MAKE`)

Delta-time driven (game time), not tied to **`time.Now`**:

| Command | Role |
|---------|------|
| **`TIMER.CREATE(duration)`** (canonical) / **`TIMER.MAKE(duration)`** (deprecated) | Create a stopped timer with a **duration** in seconds. |
| **`TIMER.START(timer)`** | Start from zero. |
| **`TIMER.STOP(timer)`** | Pause advancement. |
| **`TIMER.REWIND(timer)`** | Reset **`elapsed`** to zero. |
| **`TIMER.SETLOOP(timer, loop)`** | If **`TRUE`**, **`TIMER.UPDATE`** wraps and **`TIMER.DONE`** pulses each cycle. |
| **`TIMER.UPDATE(timer, dt)`** | Advance by **`dt`** (non-negative). |
| **`TIMER.DONE(timer)`** | **`TRUE`** for **one** call when a cycle completes (edge-triggered). |
| **`TIMER.FRACTION(timer)`** | **`elapsed/duration`**, clamped to **`0..1`**. |

**`TIMER.REMAINING`** accepts **either** a wall **`TIMER.NEW`** handle **or** a sim **`TIMER.CREATE`** / **`TIMER.MAKE`** handle (remaining sim time = **`duration − elapsed`**).

## Related

- **`STOPWATCH.*`** — elapsed wall time from **`STOPWATCH.NEW`** (`runtime/mbgame/timers.go`).
- **`ELAPSED`** (no namespace) — seconds since the **`mbgame`** module was created (`Module.t0`), not **`Time.Get`** — see **[QOL.md](QOL.md)**.
