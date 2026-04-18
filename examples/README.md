# moonBASIC examples

Runnable sample programs. Work from the **repository root** so working-directory defaults (e.g. `rpg_save.json`) match the comments in each file.

---

## Run examples with a pre-built release (recommended)

1. Download **`moonbasic`** and **`moonrun`** from **[GitHub Releases](https://github.com/CharmingBlaze/moonbasic/releases/latest)** — use the **full runtime** archive (Windows or Linux) so you have **`moonrun`** for games with a window.
2. Get the **`examples/`** sources: **clone** this repository or **download the repo as a ZIP** from GitHub. You need the files on disk; the release zip does not bundle `examples/`.
3. Open a terminal at the **repository root** (the folder that contains `examples/` and `go.mod`). Put the release binaries on your **`PATH`**, or call them with a full path.

**Check** a sample (no window):

```bash
moonbasic --check examples/spin_cube/main.mb
```

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
- **`moonrun`** — runs `.mb` or `.mbc` with the full engine. Release **`moonrun`** compiles in-process when needed; you do **not** need Go or GCC on your machine for these zips.

More on installing and using the compiler: **[docs/GETTING_STARTED.md](../docs/GETTING_STARTED.md)**. If a DLL or GPU message appears on Windows/Linux, see **`README-RELEASE.txt`** inside the full-runtime archive.

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
| [sphere_drop](sphere_drop/main.mb) | **Jolt `PHYSICS3D`**: click to spawn colored spheres on a platform, orbit camera, cull fallen balls (desktop full runtime; building from source on Windows may need [JOLT_WINDOWS_PARITY](../docs/JOLT_WINDOWS_PARITY.md)) |
| [pong](pong/main.mb) | 2D rectangles + default-font HUD |
| [platformer](platformer/main.mb) | Simple platform collision |
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

<details>
<summary><strong>Contributors: compile with <code>go run</code> from a dev tree</strong></summary>

The default repo entrypoint (`go run .` without build tags) **only compiles** `.mb` → `.mbc`. It does **not** open a window.

**Compile to bytecode** (writes `examples/spin_cube/main.mbc`):

```bash
CGO_ENABLED=1 go run . examples/spin_cube/main.mb
```

```powershell
$env:CGO_ENABLED="1"
go run . examples\spin_cube\main.mb
```

**Run the game** (opens a window) — full runtime:

```bash
CGO_ENABLED=1 go run -tags fullruntime ./cmd/moonrun examples/spin_cube/main.mb
```

```powershell
$env:CGO_ENABLED="1"
go run -tags fullruntime ./cmd/moonrun examples\spin_cube\main.mb
```

Or build `moonrun` once (`go build -tags fullruntime -o moonrun ./cmd/moonrun`) and run `moonrun examples/spin_cube/main.mb`.

Requires **CGO** and a **C toolchain** (same as building Raylib). See **[docs/BUILDING.md](../docs/BUILDING.md)** for toolchains, **[docs/DEVELOPER.md](../docs/DEVELOPER.md)**, [Makefile](../Makefile), [scripts/dev.ps1](../scripts/dev.ps1).

</details>
