# Strings, color tuples, and allocation habits

moonBASIC **interns** dynamic string results on the VM heap (similar in spirit to many scripting runtimes). You do not manage that by hand for normal `PRINT` / HUD text, but **tight gameplay loops** (physics, thousands of entities) should still avoid **building new strings every frame**.

This page maps a practical split: **numeric / boolean string APIs** vs **string-producing APIs**, aligned with how teams budget work in engines.

## Core Workflow

- **Hot path (numeric result):** use `LEN`, `INSTR`, `VAL`, `CONTAINS`, `STARTSWITH` — no string allocation.
- **Cold path (display/logging):** use `STR`, `INTERP`, `FORMAT`, string `+` — fine for HUD updates and one-off prints.
- **Color tuples:** `COLOR.TOHSV(col)` returns a 3-float tuple; use `COLOR.TOHSVX/Y/Z` in physics loops.

---

## 1–10 — Prefer in hot paths (no new string *result*)

These builtins return **`int`**, **`float`**, or **`bool`** — you are not assembling a new text body as the primary result (contrast with `STR`, `+` on strings, `INTERP`, etc.):

| # | Commands | Typical use |
|---|----------|-------------|
| 1 | `LEN(text)` | Length |
| 2 | `INSTR(hay, needle, [start])` | Search index (1-based) |
| 3 | `ASC(char)` | First code unit as integer |
| 4 | `CONTAINS(hay, needle)` | Boolean membership |
| 5 | `STARTSWITH(text, prefix)` | Boolean prefix |
| 6 | `ENDSWITH(text, suffix)` | Boolean suffix |
| 7 | `VAL(text)` | Parse to float |
| 8 | `COUNT(text, substr)` | Count occurrences — returns **int** |
| 9 | `ISALPHA` / `ISNUMERIC` / `ISALPHANUM` | Character-class tests — **bool** |
| 10 | `COLOR.TOHSVX` / `TOHSVY` / `TOHSVZ` | One HSV component from a color handle — **float** |

Use these for **per-entity** or **physics-tick** logic when the question is *search / parse / measure*, not *display*.

---

## 11–12 — Fine for UI / debug; avoid in inner physics loops

These **produce new string or tuple handles** intended for **HUD, labels, logs, tooling**:

| # | Commands | Notes |
|---|----------|--------|
| 11 | `INTERP` / `STRING.INTERP` / `STRING.INTERP$` | Fill `"{0}"` … `"{9}"` placeholders (1–10 arg overloads). `STRING.INTERP$` is the string-returning alias. See [STRING.md](STRING.md). |
| 12 | `COLOR.TOHSV(color)` | Returns **`(h, s, v)`** tuple (three floats) in one call — pairs with `COLOR.FROMHSV` / `COLOR.HSV`; see [COLOR.md](COLOR.md). |

Also treat as “cold path”: `STR`, string `+`, `FORMAT`, `LEFT` / `MID` / `REPLACE`, `SPLIT` / `JOIN`, `COLOR.TOHEX`, etc.

---

## Full Example

String handling with hot-path and cold-path split.

```basic
; --- COLD path: build display strings once per event ---
score   = 0
scoreStr = "Score: 0"

IF KEYPRESSED(KEY_SPACE) THEN
    score    = score + 10
    scoreStr = "Score: " + STR(score)   ; string concat fine here (event, not every frame)
END IF

; --- HOT path: numeric only ---
hp = 100.0
IF hp < 0.01 THEN PRINT "dead"         ; no string built in physics tick

; --- INTERP for formatted output ---
msg = INTERP("Player {0} scored {1} points!", "Alice", STR(score))
PRINT msg
```

---

## See also

- [STRING.md](STRING.md) — full string API including **`INTERP`**
- [COLOR.md](COLOR.md) — **`COLOR.FROMHSV`** and **`COLOR.TOHSV`**
- [MEMORY.md](../MEMORY.md) — handles and `FREE.ALL`
