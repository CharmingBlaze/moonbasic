# moonBASIC examples

Runnable sample programs. Work from the **repository root** so working-directory defaults (e.g. `rpg_save.json`) match the comments in each file.

### New project (no clone required)

If you have **`moonbasic`** from [Releases](https://github.com/CharmingBlaze/moonbasic-compiler/releases/latest):

```bash
moonbasic new MyGame
cd MyGame
moonrun main.mb
```

To run the samples below, clone this repo or download the source ZIP — release archives do not include `examples/`.

---

## Guide examples (documented tutorials)

Runnable copies of **Full example** blocks in `docs/systems/guides/`:

| Folder | Run |
|--------|-----|
| [guides/](guides/README.md) | `moonrun examples/guides/game_loop.mb` |
| [guides/math/](guides/README.md) | `moonrun examples/guides/math/math_2d_chase.mb` |

Check any file: `moonbasic --check examples/guides/game_loop.mb`

---

## Run examples (compiled distribution only)

Use the **official compiled distribution** from **[GitHub Releases](https://github.com/CharmingBlaze/moonbasic-compiler/releases/latest)** — not `go run`, not a locally built `moonrun`. Extract the **full runtime** zip/tar.gz so you have **`moonbasic`** and **`moonrun`** next to each other.

1. Download and extract the **full runtime** archive for your OS (Windows or Linux x64).
2. Get the **`examples/`** folder: **clone** this repo or **download the repository ZIP** from GitHub (release archives do not include `examples/`).
3. Open a terminal at the **repository root** (the folder that contains `examples/`). Either add the folder where you extracted **`moonbasic`** / **`moonrun`** to your **`PATH`**, or invoke them with a **full path** to those `.exe` / binaries.

**Check** all samples (no window, from repo root):

```bash
# PowerShell — every examples/**/*.mb file
Get-ChildItem -Recurse examples -Filter *.mb | ForEach-Object { go run . --check $_.FullName }
# Or with installed moonbasic:
moonbasic --check examples/spin_cube/main.mb
```

All **57** example scripts in `examples/` are kept in sync with `moonbasic --check`.

**Run** a sample (opens a window):

```bash
moonrun examples/spin_cube/main.mb
```

Windows (same idea):

```bat
moonbasic.exe --check examples\spin_cube\main.mb
moonrun.exe examples\spin_cube\main.mb
```

- **`moonbasic`** — compiles `.mb` → `.mbc` next to the source; use **`moonbasic path/to/main.mb`** if you only want bytecode.
- **`moonrun`** — runs `.mb` or `.mbc` with the full engine. The **distribution** `moonrun` compiles in-process when needed.

More: **[docs/GETTING_STARTED.md](../docs/GETTING_STARTED.md)** (including **Ship your game** for packaging your own `.mb` / `.mbc` for players). OS/DLL notes: **`README-RELEASE.txt`** in the full-runtime archive.

---

## Index

| Folder | Description |
|--------|-------------|
| [gui_basics](gui_basics/main.mb) | raygui: window box, label, button |
| [gui_theme](gui_theme/main.mb) | `GUI.THEMEAPPLY` — embedded official raygui `.rgs` themes |
| [gui_form](gui_form/main.mb) | Text field, slider, checkbox, tab bar |
| [gui_counter](gui_counter/main.mb) | raygui + optional TTF; small interactive demo |
| [game_math_helpers](game_math_helpers/main.mb) | Compile-only exercise of **`HDIST`**, **`YAWFROMXZ`**, **`SMOOTHERSTEP`**, … ([GAME_MATH_HELPERS](../docs/reference/GAME_MATH_HELPERS.md)) |
| [spin_cube](spin_cube/main.mb) | 3D camera, lit cube, **`Transform.*`** matrix + grid, resource cleanup |
| [sphere_drop](sphere_drop/main.mb) | **Jolt `PHYSICS3D`**: click to spawn colored spheres on a platform, orbit camera, cull fallen balls (**Windows** or **Linux** full runtime) |
| [pong](pong/main.mb) | 2D rectangles + default-font HUD |
| [platformer](platformer/main.mb) | Simple platform collision |
| [tilemap](tilemap/README.md) | Tiled TMX load, draw, collision layer |
| [gamepad](gamepad/README.md) | Controller axes + buttons |
| [fps](fps/main.mb) | Top-down arena + moving targets |
| [racing](racing/main.mb) | Top-down car + checkpoints / lap counter |
| [rpg](rpg/main.mb) | Tile-style movement + JSON save on exit |
| [mario64](mario64/README.md) | 3D hop — `main.mb` (Draw3D), `main_entities.mb` (**`MoveEntity`**, **`EntityHitsType`**, **`TranslateEntity`** — see [ENTITY.md](../docs/reference/ENTITY.md)), plus teaching variants in **`README.md`** |
| [terrain_chase](terrain_chase/README.md) | Procedural **`Terrain.*`** + **`World.*`** streaming, **`Camera.OrbitEntity`**, WASD, slow XZ-seeking enemies ([TERRAIN.md](../docs/reference/TERRAIN.md)) |
| [high_fidelity](high_fidelity/modern_template.mb) | Blitz-style 3D template (any resolution): **`Graphics3D`**, **`SetMSAA`**, **`UpdatePhysics`**, **`RENDER.Begin3D`**, **`DrawEntities`** (see [GETTING_STARTED](../docs/GETTING_STARTED.md)) |

## Fonts and assets

These demos use **`Draw.Text`** for on-screen text so you do **not** need a `.ttf` in `assets/`. For your own projects, add fonts under `assets/fonts/` and use `Font.Load` (see [FONT reference](../docs/reference/FONT.md)).

## Documentation

- [Getting started](../docs/GETTING_STARTED.md) — install, use the compiler, first window
- [Programming guide](../docs/PROGRAMMING.md) — game loop, namespaces, platforms
- [Examples (narrative)](../docs/EXAMPLES.md) — same ideas with inline explanations
- [GUI reference](../docs/reference/GUI.md) — `GUI.*` / raygui

