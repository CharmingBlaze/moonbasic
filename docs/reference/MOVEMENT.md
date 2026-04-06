# Movement helpers (`NEWXVALUE`, `WRAPVALUE`, …)

Registered from **`runtime/mbgame`** (movement / clamp / wrap helpers used for arcade-style motion). Exact names and arities are in **`compiler/builtinmanifest/commands.json`** and **`runtime/mbgame`** `register_*` files.

Use together with **`Time.Delta()`** or **`DT`** for frame-rate independent motion.

**Camera-relative XZ stepping** (yaw + forward/strafe) is implemented as top-level **`MOVEX`** / **`MOVEZ`**, or **`MOVESTEPX`** / **`MOVESTEPZ`** when you already have **`speed`** and **`dt`** — see **[MATH.md](MATH.md)** (and **`Input.Axis`** / **`Input.AxisDeg`** in **[INPUT.md](INPUT.md)**).
