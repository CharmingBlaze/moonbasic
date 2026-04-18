moonBASIC — quick start (pre-built binaries)
==============================================

GitHub Releases also ship a smaller **compiler-only** download (no moonrun, CGO off for the
compiler — no raylib.dll next to moonbasic.exe). See dist/README.md in the repo.

WHAT'S IN THIS FOLDER
---------------------
  moonbasic (or moonbasic.exe)  — Compiler: turn .mb source into .mbc bytecode, --check, --lsp
  moonrun   (or moonrun.exe)     — Full game runtime: compile and run .mb / .mbc (graphics, physics, audio)

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

TIPS
----
  • Put the binaries on your PATH if you want to run them from any directory.
  • For editor support, run:  moonbasic --lsp  (stdio language server)
  • More help: https://github.com/CharmingBlaze/moonbasic/blob/main/docs/GETTING_STARTED.md

Linux: if the app fails to start, install your distro's OpenGL / X11 / Wayland dev
       packages (see docs/BUILDING.md — same libraries the binary was linked against).

Windows: run from a normal folder; if Windows reports a missing DLL, install the
          latest "Microsoft Visual C++ Redistributable" for x64, or use MSYS2 MinGW
          runtimes if you built from source yourself.
