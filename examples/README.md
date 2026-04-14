# moonBASIC examples

Runnable sample programs. Work from the **repository root** so working-directory defaults (e.g. `rpg_save.json`) match the comments in each file.

## Requirements

- **CGO** enabled and a C toolchain (same as building Raylib).
- See [docs/BUILDING.md](../docs/BUILDING.md) for toolchains.

## Compile vs run (important)

The **default** repo entrypoint (`go run .` without build tags) **only compiles** `.mb` → `.mbc` next to the source. It does **not** open a window.

**Compile to bytecode** (writes `examples/spin_cube/main.mbc`):

```bash
CGO_ENABLED=1 go run . examples/spin_cube/main.mb
```

```powershell
$env:CGO_ENABLED="1"
go run . examples\spin_cube\main.mb
```

**Run the game** (opens a window) — use the full runtime:

```bash
CGO_ENABLED=1 go run -tags fullruntime ./cmd/moonrun examples/spin_cube/main.mb
```

```powershell
$env:CGO_ENABLED="1"
go run -tags fullruntime ./cmd/moonrun examples\spin_cube\main.mb
```

Or build `moonrun` once (`go build -tags fullruntime -o moonrun ./cmd/moonrun`) and run `.\moonrun examples\spin_cube\main.mb`.

If you use a **pre-built** `moonbasic` / `moonrun` from [Releases](https://github.com/CharmingBlaze/moonbasic/releases): `moonbasic` compiles to `.mbc`; **`moonrun`** runs `.mb` / `.mbc`.

More commands: [docs/DEVELOPER.md](../docs/DEVELOPER.md), [Makefile](../Makefile), [scripts/dev.ps1](../scripts/dev.ps1).

## Index

| Folder | Description |
|--------|-------------|
| [gui_basics](gui_basics/main.mb) | raygui: window box, label, button |
| [gui_theme](gui_theme/main.mb) | `GUI.THEMEAPPLY` — embedded official raygui `.rgs` themes |
| [gui_form](gui_form/main.mb) | Text field, slider, checkbox, tab bar |
| [gui_counter](gui_counter/main.mb) | raygui + optional TTF; small interactive demo |
| [game_math_helpers](game_math_helpers/main.mb) | Compile-only exercise of **`HDIST`**, **`YAWFROMXZ`**, **`SMOOTHERSTEP`**, … ([GAME_MATH_HELPERS](../docs/reference/GAME_MATH_HELPERS.md)) |
| [spin_cube](spin_cube/main.mb) | 3D camera, lit cube, **`Transform.*`** matrix + grid, resource cleanup |
| [sphere_drop](sphere_drop/main.mb) | **Jolt `PHYSICS3D`**: click to spawn colored spheres on a platform, orbit camera, cull fallen balls (needs CGO + Jolt libs; see [JOLT_WINDOWS_PARITY](../docs/JOLT_WINDOWS_PARITY.md) on Windows) |
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

- [Getting started](../docs/GETTING_STARTED.md) — install, first window
- [Programming guide](../docs/PROGRAMMING.md) — game loop, namespaces, platforms
- [Examples (narrative)](../docs/EXAMPLES.md) — same ideas with inline explanations
- [GUI reference](../docs/reference/GUI.md) — `GUI.*` / raygui
