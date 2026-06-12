# Project workflow — new, modules, test, package, and help

> Create a game folder, split source with `IMPORT`, validate with check/test, ship with package, and query APIs with `HELP`.

**Namespaces:** `PROJECT` · `PACKAGE` · `MODULE` · `HELP` · `TEST` · `TEMPLATE` (CLI + language)

**Overview:** [11-TOOLING.md](../11-TOOLING.md) · **Install:** [GitHub Releases](https://github.com/CharmingBlaze/moonbasic-compiler/releases/latest)

---

## Table of contents

- [At a glance](#at-a-glance)
- [When to use this system](#when-to-use-this-system)
- [Choose the right tool](#choose-the-right-tool)
- [Core workflow](#core-workflow)
- [Create and run](#create-and-run)
- [Modules](#modules)
- [Check and test](#check-and-test)
- [Package and ship](#package-and-ship)
- [HELP in scripts](#help-in-scripts)
- [Templates](#templates)
- [Common mistakes](#common-mistakes)
- [See also](#see-also)

---

## At a glance

| Idea | Detail |
|------|--------|
| **You get** | `moonbasic new`, `run`/`build`, `test`, `package`, `IMPORT`, `HELP()` |
| **You need first** | `moonbasic` + `moonrun` from Releases zip |
| **Typical use** | Every project from day one through shipping |
| **Not for** | Building the compiler — [DEVELOPER.md](../../DEVELOPER.md) |

**Why tooling exists:** Same layout for every beginner → paths, debug config, and ship steps stay predictable.

---

## When to use this system

**Use when:**

- Starting any new game (`moonbasic new`).
- Splitting `player.mb` / `levels.mb` (`IMPORT`).
- CI or pre-release (`moonbasic test`, `--check`).
- Distributing to players (`moonbasic package`).

**Skip when:**

- Single-file throwaway scratch — still use `--check` before sharing.

---

## Choose the right tool

| I want to… | Use | Not |
|------------|-----|-----|
| New folder layout | `moonbasic new MyGame` | Hand-copy random structure |
| Starter genre code | `moonbasic new --template platformer` | Empty `main.mb` only |
| Play game | `moonrun main.mb` | `go run` (maintainer build) |
| Syntax only | `moonbasic --check main.mb` | `moonrun` |
| Split files | `IMPORT "player.mb"` | One 5000-line file |
| Ask arity in REPL/script | `HELP("ENTITY.SETPOS")` | Guess parameters |
| Automated tests | `moonbasic test` | Manual click only |
| Zip for friends | `moonbasic package windows` | Email raw `.mb` only |

---

## Core workflow

1. **`moonbasic new MyGame`** — scaffold `main.mb`, `assets/`, configs.  
   **Why:** Standard paths for assets and IDE.

2. **Develop** — edit `main.mb`, add `IMPORT` modules as game grows.

3. **Check often** — `moonbasic --check main.mb`.  
   **Why:** Fast feedback ([COMPILER-ERRORS.md](COMPILER-ERRORS.md)).

4. **Run** — `moonrun main.mb` (or `moonbasic run` if CLI delegates to moonrun).  
   **Why:** Full window, audio, GPU.

5. **Test** — `moonbasic test` when you add `tests/` scripts.  
   **Why:** Regression before package.

6. **Package** — `moonbasic package windows` or `linux`.  
   **Why:** Folder with `moonrun`, your `.mb`, and `assets/`.

---

## Create and run

```bash
moonbasic new MyGame
cd MyGame
moonrun main.mb
```

| Command | Why |
|---------|-----|
| `moonbasic new <Name>` | Scaffold project |
| `moonrun main.mb` | Play with full runtime |
| `moonbasic --check main.mb` | Compile without window |
| `moonbasic build main.mb` | Emit `.mbc` bytecode |

See [00-START.md](../00-START.md) for foundation loop inside `main.mb`.

---

## Modules

```basic
IMPORT "player.mb"
IMPORT "weapons.mb"

APP.OPEN(800, 600, "Modular")
; player.mb defines InitPlayer(), called here
InitPlayer()
```

**Why `IMPORT`:** One namespace per file; `moonbasic --check` validates all imports.

Rules: paths relative to project; use quotes; forward slashes in scripts.

---

## Check and test

```bash
moonbasic --check main.mb
moonbasic test
```

| Command | Why |
|---------|-----|
| `--check` | Parser + semantic errors |
| `--check --strict-deprecated` | Migration warnings |
| `moonbasic test` | Run test scripts in project |

Details: [DEBUG-AND-TESTING.md](DEBUG-AND-TESTING.md), [COMPILER-ERRORS.md](COMPILER-ERRORS.md).

---

## Package and ship

```bash
moonbasic package windows
moonbasic package linux
moonbasic pack
```

| Command | Why |
|---------|-----|
| `package windows` / `linux` | Player folder with runtime + game |
| `pack` | Archive bundle (layout per CLI help) |

Ship `assets/`, save paths, and `options.json` beside `main.mb` so [FILES-AND-JSON.md](FILES-AND-JSON.md) paths resolve.

---

## HELP in scripts

```basic
HELP("ENTITY.SETPOSITION")
```

**Why:** Prints arity and aliases when coding without IDE — complements LSP (`moonbasic --lsp` + VSIX).

---

## Templates

```bash
moonbasic new --template 3d MyGame
moonbasic new --template platformer MyGame
moonbasic new --template ui MyGame
moonbasic new --template physics MyGame
```

**Why:** Genre-specific `main.mb` starter — faster than blank file.

---

## Common mistakes

| Mistake | Fix |
|---------|-----|
| Run from wrong folder | `cd` to project root; assets paths are relative |
| `IMPORT` path wrong | Quote path; check case on disk for filename |
| Ship only `.mb` | Package includes `moonrun` + assets |
| Skip `--check` in CI | Add check step before release |
| `HELP` wrong string | Use full `NAMESPACE.COMMAND` |

---

## See also

- [00-START.md](../00-START.md) — first project
- [GETTING_STARTED.md](../../GETTING_STARTED.md) — install details
- [DEBUG-AND-TESTING.md](DEBUG-AND-TESTING.md) — test folder patterns
- [ASSETS-PIPELINE.md](ASSETS-PIPELINE.md) — ship `assets.json`
