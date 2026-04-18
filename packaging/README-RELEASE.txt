moonBASIC — quick start (pre-built binaries)
==============================================

GitHub Releases also ship a smaller **compiler-only** download (no moonrun, CGO off for the
compiler — no raylib.dll next to moonbasic.exe). See dist/README.md in the repo.

WHAT'S IN THIS FOLDER
---------------------
  moonbasic (or moonbasic.exe)  — Compiler: turn .mb source into .mbc bytecode, --check, --lsp
                                 (uses the full builtin catalog — same command names as the engine API)
  moonrun   (or moonrun.exe)     — Full game runtime: compile and run .mb / .mbc (graphics, physics, audio)

  For “all commands” at RUN TIME (playing/running a game), you need moonrun in this folder.
  For “all commands” at CHECK/COMPILE time only, moonbasic alone is enough.

  Windows (full runtime zip): the two executables and this README are enough — **libgcc**,
  **libstdc++**, and **winpthread** are linked into the `.exe` files. Raylib is compiled from
  sources (no `raylib.dll`). You should not need MinGW companion DLLs beside the binaries.

FIRST STEPS
-----------
  1. Extract this zip/tar anywhere you like (Desktop, Projects, etc.).

  2. Open a terminal in that folder:
       Windows: Shift+right-click the folder → "Open in Terminal", or cmd/PowerShell and cd to the folder.
       Linux:   cd /path/to/extracted/folder

  3. Check that it works:
       Windows:   moonbasic.exe --version
       Linux:     chmod +x moonbasic moonrun   (first time only)
                    ./moonbasic --version

  4. Compile a program to bytecode:
       moonbasic path\to\yourgame.mb
     This writes yourgame.mbc next to the source.

  5. Run a game (needs full runtime):
       moonrun path\to\yourgame.mb
     or:  moonrun yourgame.mbc
     moonrun compiles .mb inside the same program — you do NOT need Go, GCC, or moonbasic
     on PATH to play; extract the zip and run.

TIPS
----
  • Put the binaries on your PATH if you want to run them from any directory.
  • For editor support, run:  moonbasic --lsp  (stdio language server)
  • More help: https://github.com/CharmingBlaze/moonbasic/blob/main/docs/GETTING_STARTED.md

Linux: if the app fails to start, ensure GPU drivers and a normal desktop OpenGL stack
       are installed (run-time libs, not compiler -dev packages). See docs/BUILDING.md
       only if you build from source.

Windows: run from a normal folder. If Windows reports a missing **non-system** DLL, you may
          have a partial copy, a mixed install, or an antivirus quarantine — re-extract the
          **entire** full-runtime zip from the same release.

Windows: "Entry Point Not Found" / nanosleep64 — usually from copying only `moonrun.exe`
          without the matching `moonbasic.exe` from the **same** zip, or from PATH picking up
          an older MinGW `libwinpthread` / `libgcc` DLL. Re-extract the full zip; do not drop
          stray MinGW DLLs next to the exes unless you know you need them.
