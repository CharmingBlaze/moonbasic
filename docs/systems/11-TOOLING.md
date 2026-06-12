# Tooling systems: PROJECT, PACKAGE, MODULE, HELP, TEST, TEMPLATE

> Create projects, split source files, check and ship games, and get built-in command help.

**Install:** [GitHub Releases](https://github.com/CharmingBlaze/moonbasic-compiler/releases/latest) — full runtime includes **`moonbasic`** and **`moonrun`**.

**Deep guide:** [guides/PROJECT-WORKFLOW.md](guides/PROJECT-WORKFLOW.md)

**See also:** [GETTING_STARTED.md](../GETTING_STARTED.md) · [PROGRAMMING.md](../PROGRAMMING.md)

---

## Table of contents

- [PROJECT system](#project-system)
- [PACKAGE system](#package-system)
- [MODULE system](#module-system)
- [HELP system](#help-system)
- [TEST system](#test-system)
- [TEMPLATE system](#template-system)
- [Quick reference](#quick-reference)
- [See also](#see-also)

---

## PROJECT system

Scaffold and run a game folder with `main.mb`.

### `moonbasic new <ProjectName>`

Creates a new directory with `main.mb`, `assets/`, `.vscode/launch.json`, and `README.md`.

```bash
moonbasic new MyGame
cd MyGame
```

### `moonbasic new --template <name> <ProjectName>`

Starter templates: **`empty`**, **`3d`**, **`platformer`**, **`ui`**, **`physics`**.

```bash
moonbasic new --template 3d My3DGame
```

### Run your project

From the project folder (where `main.mb` lives):

```bash
moonrun main.mb
```

Optional project helpers (same folder):

| Command | Description |
|---------|-------------|
| `moonbasic run` | Run `main.mb` via project runner |
| `moonbasic build` | Compile to `build/main.mbc` |

Prefer **`moonrun main.mb`** when `moonrun` is on your `PATH` from the full runtime install.

### Check and compile

| Command | Description |
|---------|-------------|
| `moonbasic --check main.mb` | Parse and type-check (no window) |
| `moonbasic main.mb` | Write `main.mbc` next to source |
| `moonbasic --lsp` | Language server for editors |

---

## PACKAGE system

Ship a folder players can extract and run.

### `moonbasic package windows` / `moonbasic package linux`

Packages bytecode and assets for the chosen platform layout (ships `.mbc` + assets; add **`moonrun`** from the same release when distributing).

```bash
cd MyGame
moonbasic build
moonbasic package windows
```

### Folder bundle (manual)

Copy from the **full runtime** zip into your game folder:

- `moonrun` / `moonrun.exe`
- Your `main.mb` or `build/main.mbc`
- `assets/`

Players run:

```bash
moonrun main.mb
```

See [GETTING_STARTED.md](../GETTING_STARTED.md) — **Ship your game**.

### `moonbasic pack`

Creates a distributable archive from a script (see `moonbasic pack --help`).

---

## MODULE system

Split code across multiple `.mb` files.

### `IMPORT "path.mb"`

Merges another source file into the compile unit. Only **`.mb`** files are valid targets — not markdown docs.

**Example:**

```basic
IMPORT "player.mb"
IMPORT "ui/hud.mb"
```

Compile or check the **entry** file only:

```bash
moonbasic --check main.mb
```

The compiler expands all `IMPORT` directives before analysis.

See [reference/INCLUDE.md](../reference/INCLUDE.md) and [LANGUAGE.md](../LANGUAGE.md).

---

## HELP system

Built-in command documentation in the console.

### `HELP("NAMESPACE.COMMAND")`

Prints manifest help for a builtin when available.

**Example:**

```basic
HELP("ENTITY.SETPOSITION")
HELP("CAMERA")
```

**Status:** Partial — coverage grows with the manifest; deep reference remains in [COMMANDS.md](../COMMANDS.md) and [systems/README.md](README.md).

---

## TEST system

### `moonbasic test`

Runs compiler and runtime smoke checks from the moonBASIC install (validates language and foundation examples when present).

```bash
moonbasic test
```

Use **`moonbasic --check`** on your own scripts in CI:

```bash
moonbasic --check main.mb
```

---

## TEMPLATE system

Same as **`moonbasic new --template`** — starter projects without hand-writing boilerplate.

| Template | Use case |
|----------|----------|
| `empty` | Minimal window loop |
| `3d` | Camera + cube + lighting |
| `platformer` | 2D movement starter |
| `ui` | GUI widgets demo |
| `physics` | Physics world stub |

```bash
moonbasic new --template platformer MyPlatformer
cd MyPlatformer
moonrun main.mb
```

---

## Quick reference

| Goal | Command |
|------|---------|
| New project | `moonbasic new MyGame` |
| 3D starter | `moonbasic new --template 3d MyGame` |
| Run game | `moonrun main.mb` |
| Check only | `moonbasic --check main.mb` |
| Bytecode | `moonbasic main.mb` → `main.mbc` |
| Build to folder | `moonbasic build` |
| Package | `moonbasic package windows` |
| Split code | `IMPORT "other.mb"` |
| Command help | `HELP("RENDER.CLEAR")` |
| Smoke tests | `moonbasic test` |

---

## See also

- [systems/README.md](README.md) — all 40 game systems
- [DOCUMENTATION_STYLE_GUIDE.md](../DOCUMENTATION_STYLE_GUIDE.md) — how these docs are written
- [editors/vscode-moonbasic](../editors/vscode-moonbasic/README.md) — VSIX from Releases
