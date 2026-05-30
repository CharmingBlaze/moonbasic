# INCLUDE Commands

Compile-time merging of **MoonBASIC source files** (`.mb`) into a single program. Use this to split **game logic** across modules; **documentation** in **`.md`** files is not compiled (see [**`STYLE_GUIDE.md`**](../../STYLE_GUIDE.md) for reference layout).

Page shape follows [**`DOC_STYLE_GUIDE.md`**](../DOC_STYLE_GUIDE.md) (**WAVE pattern**).

## Core Workflow

1. Put shared **`FUNCTION`** definitions, **`TYPE`**, and top-level data in library files (e.g. `lib/ui.mb`).
2. From your entry file (e.g. `main.mb`), write **`INCLUDE "relative/path.mb"`** where you want that file’s **top-level statements** to appear.
3. Set **`MOONBASIC_PATH`** (optional) so includes can resolve from standard roots when not found next to the including file (see [**`PACKAGES.md`**](../PACKAGES.md)).
4. Compile or **`--check`** the **entry** `.mb` only — the compiler expands all **`INCLUDE`** directives before semantic analysis and codegen.

**Important:** The compiler only reads **`.mb`** sources. Reference prose in **`docs/**/*.md`** is for humans; it is **not** an **`INCLUDE`** target.

---

### `INCLUDE "path.mb"` 

Merges another **`.mb`** file into the current compilation unit.

- **Arguments / form**:
  - **`path.mb`**: Quoted path. The keyword **`INCLUDE`** is case-insensitive (`include` and **`INCLUDE`** are equivalent).
- **Path resolution**:
  1. If the path is **absolute**, it is used as-is (cleaned).
  2. Otherwise it is resolved **relative to the directory of the file that contains this `INCLUDE`**.
  3. If the file is not found, the compiler searches extra roots from **`MOONBASIC_PATH`** and **installed package** directories (see pipeline **`SyncPackageIncludeRoots`** / **`TryOpenInclude`**).
- **Merge semantics** (after parse, before codegen):
  - **Types** and **functions** collected from included files are **prepended** to the host program’s type and function lists (see **`compiler/include/expand.go`**).
  - **Top-level statements** from the included file are **inserted in place of** the **`INCLUDE`** line.
  - **`INCLUDE`** inside a function body is expanded the same way (nested statement lists).
- **Include-once:** The same resolved **absolute** path is expanded **at most once** per compilation. Later **`INCLUDE`** of the same file is skipped (no duplicate top-level runs).
- **Cycles:** A circular chain (e.g. A includes B includes A) is a **compile-time error** with a clear message.
- **Globals and names:** After expansion, the program is a **single** AST. Top-level **`DIM`**, assignments, and **`FUNCTION`** names follow the same visibility rules as one large file—symbols defined in an included module are visible to the host and to other includes (subject to declaration order for forward references as usual). Use **`LOCAL`** inside **`FUNCTION`** for temporaries that must not collide with globals ([**`LANGUAGE.md`**](../LANGUAGE.md)).
- **Runtime:** There is **no per-frame cost** for **`INCLUDE`**; merging happens at **compile time** only.

---

## Full Example

**`main.mb`** (entry point):

```basic
INCLUDE "include_example_lib.mb"

; Calls into the included file's FUNCTION; uses shared top-level name
PRINT "Count is " + STR(NextCount())
```

**`include_example_lib.mb`** (beside `main.mb` or on an include root):

```basic
; Shared global (top-level)
count = 0

FUNCTION NextCount()
    count = count + 1
    RETURN count
ENDFUNCTION
```

Compile or type-check **`main.mb`** only: **`go run . --check main.mb`**. The merged program behaves as if both files were one source.

---

## See also

- [**`LANGUAGE.md`**](../LANGUAGE.md) — functions, **`LOCAL`**, control flow
- [**`PACKAGES.md`**](../PACKAGES.md) — **`MOONBASIC_PATH`**, shipping optional **`.mb`** with packages
- [**`STYLE_GUIDE.md`**](../../STYLE_GUIDE.md) — naming and project conventions
- [**`ARCHITECTURE.md`**](../../ARCHITECTURE.md) — compiler pipeline and **`INCLUDE`** stage
- [**`DOC_STYLE_GUIDE.md`**](../DOC_STYLE_GUIDE.md) — WAVE reference layout for other pages
