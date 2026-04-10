# Terrain performance and loading

Command reference for `TERRAIN.*`, `CHUNK.*`, and world streaming lives in **[docs/reference/TERRAIN.md](reference/TERRAIN.md)**.

## Main-thread safety

Raylib mesh creation (`GenMeshHeightmap` and uploads) runs on the **main thread** with an active OpenGL context. **Large procedural work** (huge heightmaps, tight loops that edit every cell, or loading many chunks at once) can block the UI long enough for Windows to show **“Not Responding.”**

**Do this instead:**

- Treat heavy generation as an explicit **loading phase**: keep **`RENDER.FRAME`** running so the OS sees a responsive window.
- Use **`WINDOW.SETLOADINGMODE(true)`** while meshes are still building: the engine keeps polling **`WINDOW.SHOULDCLOSE`** and related events, but **`TERRAIN.DRAW`** is skipped so you are not stacking expensive draws on top of mesh work. Show a simple clear color, a flat floor, or a “Loading…” UI until terrain is ready, then call **`WINDOW.SETLOADINGMODE(false)`**.
- Use **`TERRAIN.SETMESHBUILDBUDGET(terrain, n)`** with a small **`n`** (for example **1–4**): each **`WORLD.UPDATE`** tick schedules at most **`n`** chunk mesh builds (each build still runs **`GenMeshHeightmap`** to completion for that chunk). The budget is **per chunk per tick**, not a time slice inside a single C call—so one slow chunk can still stall a frame, which is why **`SETASYNCMESHBUILD`**, **`SETLOADINGMODE`**, and smaller **`TERRAIN.SETCHUNKSIZE`** matter on weak GPUs. **`0`** means unlimited builds per tick (default).
- Use **`TERRAIN.SETASYNCMESHBUILD(terrain, true)`** to run the **CPU** part of each chunk heightmap (grayscale sampling into a buffer) on a **background goroutine**. The **`GenMeshHeightmap`** call and GPU material setup still run on the **main thread** when the job is drained (during **`WORLD.UPDATE`**), so the window can keep processing frames between jobs. Pair with **`SETMESHBUILDBUDGET`** and **`WINDOW.SETLOADINGMODE`** for large worlds.

The engine also calls **`PollInputEvents()`** while draining mesh jobs and during large terrain brush edits so Windows is less likely to mark the app as **“Not Responding”** during heavy work.

Procedural generation in **MoonBASIC script** is still interpreted on the main thread; spreading work with **`SETMESHBUILDBUDGET`**, **`SETASYNCMESHBUILD`**, and **`SETLOADINGMODE`** avoids freezing the window while GPU meshes catch up.

## Chunking

Internal terrain is built from **chunks** of height samples. Very large chunks mean fewer draw calls but **more vertices per mesh** and longer **blocking** `GenMeshHeightmap` calls.

**Practical guidance:**

- Prefer **moderate chunk sizes** (the default internal chunk size is **64**; adjust with **`TERRAIN.SETCHUNKSIZE`** to match your world scale).
- For broad compatibility (including integrated GPUs), avoid relying on a **single** enormous heightfield mesh; use **streaming** (`WORLD.*` / **`CHUNK.*`**) so only nearby chunks hold GPU meshes.

## Flat-floor fallback (LOD 0)

A **large untextured quad** (**`FLAT`**) or **`GRID3`** ground is a good **immediate** stand-in while heightfield data loads or meshes build. **Keep this LOD0 floor visible until the heightfield is fully ready** (for example while **`WORLD.ISREADY(terrain)`** is false): if you switch to **`TERRAIN.DRAW`** too early, collision and visuals may assume a complete surface while some chunks still have no mesh, and actors can intersect empty space. The samples **`examples/terrain_stress/main.mb`** (status text) and **`examples/terrain_async/main.mb`** (spinner + LOD swap) draw **`FLAT` + `GRID3`** during loading mode, then **`TERRAIN.DRAW`** after **`WORLD.ISREADY`**. Swap to full terrain only once **`WINDOW.SETLOADINGMODE(false)`** and your readiness check pass.

## See also

- [reference/TERRAIN.md](reference/TERRAIN.md) — API details
- [reference/WORLD.md](reference/WORLD.md) — `World.Update`, preload, streaming center
- [ARCHITECTURE.md](../ARCHITECTURE.md) — open-world stack
