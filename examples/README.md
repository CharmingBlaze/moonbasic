# moonBASIC examples

Runnable sample programs. Run them from the **repository root** so working-directory defaults (e.g. `rpg_save.json`) match the comments in each file.

## Requirements

- **CGO** enabled and a C toolchain (same as building Raylib).
- Build or run the driver from this repo, for example:

```bash
# Linux / macOS
CGO_ENABLED=1 go run . examples/spin_cube/main.mb
```

```powershell
# Windows (PowerShell)
$env:CGO_ENABLED="1"
go run . examples\spin_cube\main.mb
```

If you use a compiled `moonbasic` binary, replace `go run .` with `moonbasic`.

## Index

| Folder | Description |
|--------|-------------|
| [gui_basics](gui_basics/main.mb) | raygui: window box, label, button |
| [gui_theme](gui_theme/main.mb) | `GUI.THEMEAPPLY` — embedded official raygui `.rgs` themes |
| [gui_form](gui_form/main.mb) | Text field, slider, checkbox, tab bar |
| [spin_cube](spin_cube/main.mb) | 3D camera, lit cube, **`Transform.*`** matrix + grid, resource cleanup |
| [pong](pong/main.mb) | 2D rectangles + default-font HUD |
| [platformer](platformer/main.mb) | Simple platform collision |
| [fps](fps/main.mb) | Top-down arena + moving targets |
| [racing](racing/main.mb) | Top-down car + checkpoints / lap counter |
| [rpg](rpg/main.mb) | Tile-style movement + JSON save on exit |
| [mario64](mario64/README.md) | 3D hop — `main.mb` (Draw3D), `main_entities.mb` (**`MoveEntity`**, **`EntityHitsType`**, **`TranslateEntity`** — see [ENTITY.md](../docs/reference/ENTITY.md)), plus teaching variants in **`README.md`** |
| [high_fidelity](high_fidelity/modern_template.mb) | Blitz-style 3D template (any resolution): **`Graphics3D`**, **`SetMSAA`**, **`UpdatePhysics`**, **`RENDER.Begin3D`**, **`DrawEntities`** (see [GETTING_STARTED](../docs/GETTING_STARTED.md)) |

## Fonts and assets

These demos use **`Draw.Text`** for on-screen text so you do **not** need a `.ttf` in `assets/`. For your own projects, add fonts under `assets/fonts/` and use `Font.Load` (see [FONT reference](../docs/reference/FONT.md)).

## Documentation

- [Getting started](../docs/GETTING_STARTED.md) — install, first window
- [Programming guide](../docs/PROGRAMMING.md) — game loop, namespaces, platforms
- [Examples (narrative)](../docs/EXAMPLES.md) — same ideas with inline explanations
- [GUI reference](../docs/reference/GUI.md) — `GUI.*` / raygui
