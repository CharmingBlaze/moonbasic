# Changelog

This document tracks the recent development history of moonBASIC.

---

## Version 0.1 (April 2026)

### April 4, 2026

-   **Math**: `CLAMP`, `LERP`, and `WRAP` commands now use the formulas from Raylib 5.5 for better consistency and performance. Trigonometric functions still use the standard Go math library.
-   **File I/O**: Clarified the behavior of file writing commands. `FILE.WRITE` and its alias `WRITEFILE` write raw data, while `FILE.WRITELN` and `WRITEFILELN` append a newline character.
-   **Control Flow**: Implemented a full suite of `DO...LOOP` structures (`DO WHILE`, `DO UNTIL`, `DO...LOOP WHILE`, `DO...LOOP UNTIL`). Added `EXIT` and `CONTINUE` statements for all loop types (`FOR`, `WHILE`, `REPEAT`, `DO`) and `EXIT FUNCTION` for early returns from functions.
-   **Parser**: Fixed a bug where `NEXT` in a `FOR` loop could incorrectly consume the first part of the next statement if it was on a new line. The optional variable after `NEXT` is now only considered if it's on the same line.
-   **Codebase**: Refactored loop-related parsing logic into a dedicated `parser_stmts_loop.go` file for better organization.
-   **Internal**: Removed obsolete `strmod.*` command registrations. The core `runtime` now handles all string-related built-in commands.
