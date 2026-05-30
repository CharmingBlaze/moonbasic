# moonBASIC examples

Runnable sample programs. Run from the **repository root** (the folder that contains `examples/`).

---

## Prerequisites

1. Install the **full runtime** from **[GitHub Releases](https://github.com/CharmingBlaze/moonbasic/releases/latest)**.
2. Clone **this repository** or download its ZIP — release archives do not include `examples/`.

Check:

```bash
moonrun --version
```

---

## Quick start (no clone)

If you only have the release binaries:

```bash
moonbasic new MyGame
cd MyGame
moonrun main.mb
```

---

## Run an example

From the repo root:

```bash
moonbasic --check examples/spin_cube/main.mb
moonrun examples/spin_cube/main.mb
```

Windows:

```bat
moonrun.exe examples\spin_cube\main.mb
```

Add the folder containing `moonbasic` / `moonrun` to **PATH**, or use full paths to the binaries.

---

## Sample index

| Folder | What it shows |
|--------|----------------|
| [spin_cube/](spin_cube/) | Minimal 3D spinning cube |
| [platformer/](platformer/) | 2D platformer basics |
| [pong/](pong/) | Simple 2D game |
| [tilemap/](tilemap/) | Tiled map load and collision |
| [gamepad/](gamepad/) | Controller input |
| [gamejam/](gamejam/) | Zero-asset jam demo (built-in sprites/sounds) |
| [fps/](fps/) | First-person style sample |
| [rpg/](rpg/) | Top-down RPG-style sample |

See each folder's README for controls and notes.

---

## More help

- [Getting started](../docs/GETTING_STARTED.md)
- [Programming guide](../docs/PROGRAMMING.md)
- [Command reference](../docs/COMMANDS.md)
