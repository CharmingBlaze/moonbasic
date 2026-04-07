# Instant-game / QOL shortcuts (`runtime/mbgame`)

**Curriculum / designed API surface:** [moonbasic-command-set/README.md](moonbasic-command-set/README.md) (memory-aware tables + QOL section).

The **`mbgame`** module registers **DarkBASIC-style** short names and helpers: **`SCREENW`**, **`SCREENH`**, **`DT`**, **`MX`**, **`MY`**, **`ENDGAME`**, **`ELAPSED`**, **`FRAMECOUNT`**, collision and movement math, easing, noise, **`CONFIG.*`**, timers, **`GAME.*`** volume/screen flash, etc.

- **Registry keys** are **one per uppercase name** — implementations live only in **`runtime/mbgame`** (do not re-register the same dotted key elsewhere).
- Source uses **dotted** or **namespace** style where applicable; names are **case-agnostic** (normalized to uppercase).

## Input / window / time shortcuts

| Spec-style name | Registry / usage | Notes |
|-----------------|------------------|--------|
| Screen size | **`SCREENW`**, **`SCREENH`** | Raylib render size when CGO is on. |
| Frame delta | **`DT`** | Use with **`Time.Delta()`** pattern; **`DT`** mirrors it for short scripts. |
| Mouse | **`MX`**, **`MY`**, **`MWHEEL`**, etc. | See **`INPUT.*`** for full input. |
| Frames / time | **`FRAMECOUNT`**, **`ELAPSED`** | **`ELAPSED`** is wall seconds since module init (`t0`), not necessarily **`Time.Get`**. |
| Exit | **`ENDGAME`** | Stops the VM (**`TerminateVM`**), not “close window then quit” alone. |

## Related docs

- **[COLOR.md](COLOR.md)** — **`RGB`**, **`ARGB`**, mixing.  
- **[COLLISION.md](COLLISION.md)** — **`BOXCOLLIDE`**, **`AABBCOLLIDE`**, …  
- **[RAYCAST.md](RAYCAST.md)** — **`RAY.*`**, **`RAY2D.*`** (picking and 2D ray tests).  
- **[BLITZ3D.md](BLITZ3D.md)** — **`Camera.Turn`**, **`Entity.Create`**, **`KeyHit`**, **`JoyX`**, …  
- **[GAMEHELPERS.md](GAMEHELPERS.md)** — **`BOXTOPLAND`** (sphere landing Y).  
- **[MATH.md](MATH.md)** — **`MOVEX`**, **`MOVEZ`**, **`IIF`**, **`IIF$`**.  
- **[INPUT.md](INPUT.md)** — **`Input.Axis`**, **`Input.AxisDeg`**.  
- **[MOVEMENT.md](MOVEMENT.md)** — **`WRAPVALUE`**, **`NEWXVALUE`**, …  
- **[EASING.md](EASING.md)** — **`EASEIN`**, **`EASELERP`**, …  
- **[PROCEDURAL.md](PROCEDURAL.md)** — **`PERLIN`**, **`RNDRANGE`**, …  
- **[CONFIG.md](CONFIG.md)** — **`CONFIG.LOAD`**, **`CONFIG.GET`**, …  
- **[TIMER.md](TIMER.md)** — **`TIMER.NEW`** vs **`TIMER.MAKE`**.  
