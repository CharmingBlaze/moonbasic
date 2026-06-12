# Compiler errors — catch mistakes before `moonrun`

> Understand moonBASIC diagnostics: file/line errors, unknown commands, arity mistakes, and how to fix them fast.

**Namespaces:** `ERROR` (compiler) · **Status:** Shipped at compile time

**Related runtime:** [DEBUG-AND-TESTING.md](DEBUG-AND-TESTING.md) · **Message catalog:** [ERROR_MESSAGES.md](../../ERROR_MESSAGES.md)

---

## Table of contents

- [At a glance](#at-a-glance)
- [When to use this system](#when-to-use-this-system)
- [Choose the right tool](#choose-the-right-tool)
- [Core workflow](#core-workflow)
- [What the compiler reports](#what-the-compiler-reports)
- [Common error types](#common-error-types)
- [Strict mode](#strict-mode)
- [Full example session](#full-example-session)
- [Common mistakes](#common-mistakes)
- [See also](#see-also)

---

## At a glance

| Idea | Detail |
|------|--------|
| **You get** | Errors before run: bad commands, wrong arg count, typos with suggestions |
| **You need first** | `moonbasic` CLI from [Releases](https://github.com/CharmingBlaze/moonbasic-compiler/releases/latest) |
| **Typical use** | Every save — CI, VS Code task, pre-ship check |
| **Not for** | Runtime logic bugs — use `DEBUG.LOG` / watches |

**Why compile-time errors:** Fixing `ENTITY.SETPOSITON` before `moonrun` saves minutes of “nothing happened” debugging.

---

## When to use this system

**Use when:**

- After editing any `.mb` file — `moonbasic --check main.mb`.
- Before packaging — same check on whole project.
- Learning APIs — errors teach correct command names.
- Deprecation cleanup — `--strict-deprecated`.

**Skip when:**

- Script already checked and you only tweak assets (images/audio).

---

## Choose the right tool

| I want to… | Use | Not |
|------------|-----|-----|
| Syntax / API errors | `moonbasic --check` | `moonrun` only |
| In-game typo while coding | `HELP("ENTITY.SETPOS")` | Guess from memory |
| Runtime position debug | `DEBUG.WATCH` | Compiler |
| Full overload list | [COMMAND_REGISTRY.md](../COMMAND_REGISTRY.md) | Trial and error |
| LSP in editor | `moonbasic --lsp` + VSIX | Manual check only |

---

## Core workflow

1. **Edit** script in `main.mb` or `IMPORT` modules.

2. **Check** — `moonbasic --check main.mb` (or `moonbasic test` for test folder).

3. **Read** `file:line:column` — jump to caret in editor.

4. **Fix** typo / arity — compiler **did you mean** helps command names.

5. **Re-check** until clean output.

6. **Run** — `moonrun main.mb` for window, audio, GPU.

**Why order matters:** `--check` is fast and needs no GPU; `moonrun` validates the full game.

---

## What the compiler reports

| Feature | Example |
|---------|---------|
| File and line | `main.mb:42:5` |
| Unknown command | `Unknown command: ENTITY.SETPOSITON` |
| Suggestions | `Did you mean: ENTITY.SETPOSITION` |
| Bad arity / types | Semantic phase after parse |
| Import errors | Missing `IMPORT` path |

**Example output:**

```text
main.mb:42:5
ENTITY.SETPOSITON(player, 0, 1, 5)
      ^^^^^^^^^^
Unknown command: ENTITY.SETPOSITON
Did you mean: ENTITY.SETPOSITION
```

---

## Common error types

| Symptom | Likely cause | Fix |
|---------|--------------|-----|
| Unknown command | Typo or wrong namespace | `HELP`, registry, did-you-mean |
| Wrong arg count | Extra/missing parameter | [COMMAND_REGISTRY.md](../COMMAND_REGISTRY.md) arity |
| Import failed | Path wrong | Forward slashes; file next to project |
| Undefined variable | Name typo | Case-insensitive but must be declared |
| Deprecated API | Old checklist name removed | `--strict-deprecated` lists replacements |

---

## Strict mode

```bash
moonbasic --check --strict-deprecated main.mb
```

**Why:** Surfaces APIs you should migrate before they break in a future release.

---

## Full example session

```bash
moonbasic new ErrorDemo
cd ErrorDemo
moonbasic --check main.mb
moonrun main.mb
```

In CI or VS Code task, run `--check` on every commit. See [DEBUG-AND-TESTING.md](DEBUG-AND-TESTING.md) for `moonbasic test`.

---

## Common mistakes

| Mistake | Fix |
|---------|-----|
| Only `moonrun` after big edit | Run `--check` first |
| Ignore column caret | Error is at `^` under typo |
| Wrong namespace | `SETPOS` on entity → `ENTITY.SETPOS` |
| Case confusion | Commands are case-insensitive — typo still errors |
| Runtime crash = compiler | Logic bugs need `DEBUG.LOG` |

Runtime crashes and call stacks: [ERROR_MESSAGES.md](../../ERROR_MESSAGES.md).

---

## See also

- [DEBUG-AND-TESTING.md](DEBUG-AND-TESTING.md) — tests, watches, FPS graph
- [PROJECT-WORKFLOW.md](PROJECT-WORKFLOW.md) — `moonbasic test` in projects
- [COMMAND_REGISTRY.md](../COMMAND_REGISTRY.md) — every overload
- [11-TOOLING.md](../11-TOOLING.md) — HELP, MODULE, CLI overview
