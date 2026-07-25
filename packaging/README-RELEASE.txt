moonBASIC — quick start (pre-built binaries)

==============================================



GitHub Releases also ship a smaller **compiler-only** download (no moonrun, CGO off for the

compiler — no raylib.dll next to moonbasic.exe). See dist/README.md in the repo.



WHAT'S IN THIS FOLDER

---------------------

  moonbasic (or moonbasic.exe)  — Compiler: turn .mb source into .mbc bytecode, --check, --lsp,

                                 moonbasic new (scaffold a project)

                                 (uses the full builtin catalog — same command names as the engine API)

  moonrun   (or moonrun.exe)     — Full game runtime: compile and run .mb / .mbc (graphics, physics, audio)



  For “all commands” at RUN TIME (playing/running a game), you need moonrun in this folder.

  For “all commands” at CHECK/COMPILE time only, moonbasic alone is enough.



  Windows (full runtime zip): the two executables and this README are enough — **libgcc**,

  **libstdc++**, and **winpthread** are linked into the `.exe` files. Raylib, Jolt, and ENet

  are compiled/linked into the binaries (no `raylib.dll`, no `libenet`). You should not need

  MinGW companion DLLs beside the binaries.



  Linux / macOS: same idea — extract and run `./moonrun`. No `libenet` / Homebrew `enet`

  package is required for networking. Linux `moonrun` also embeds `libstdc++` / `libgcc`

  (Jolt) so you should not need matching g++ runtime packages from the build host.



FIRST STEPS

-----------

  1. Extract this zip/tar anywhere you like (Desktop, Projects, etc.).



  2. Open a terminal in that folder:

       Windows: Shift+right-click the folder → "Open in Terminal", or cmd/PowerShell and cd to the folder.

       Linux:   cd /path/to/extracted/folder

       macOS:   cd /path/to/extracted/folder && chmod +x moonbasic moonrun



  3. Check that it works:

       Windows:   moonrun.exe --version

       Linux/macOS: ./moonrun --version



  4. Start a new game (optional):

       moonbasic new MyGame

       cd MyGame

       moonrun main.mb



  5. Or run an existing script:

       moonrun path\to\yourgame.mb

     moonrun compiles .mb inside the same program — you do NOT need Go or GCC on the player machine.



  6. Lint without running (optional):

       moonbasic --check path\to\yourgame.mb



  7. Compile to bytecode only (optional):

       moonbasic path\to\yourgame.mb   → writes yourgame.mbc next to the source



TIPS

----

  • Language reference (syntax, $"..." strings, ENUM, multi-return): docs/LANGUAGE.md on GitHub.

  • Example projects (tilemap, gamepad, platformer): examples/ folder in the source repo.

  • Visual Studio Code / Cursor (easiest):

      Windows: double-click  INSTALL-VSCODE.bat  in this folder

      Linux/macOS:  ./INSTALL-VSCODE.sh

      Or run:  moonbasic install-vscode

    This installs the extension and sets moonbasic.languageServerPath / moonbasic.moonrunPath

    automatically. The .vsix is included in this zip (moonbasic-*-vscode.vsix).

    Then open any .mb file — Ctrl+F5 run, Ctrl+Shift+C check, Alt+H help.

    Run and Debug -> "Debug moonBASIC" for breakpoints (needs moonrun in this folder).

  • For editor support (any client), run:  moonbasic --lsp  (stdio language server)

  • Porting from BlitzBASIC? See docs/reference/MIGRATION.md for commands not in this release.

  • More help: https://github.com/CharmingBlaze/moonbasic/blob/main/docs/GETTING_STARTED.md

  • Downloads / IDE: https://github.com/CharmingBlaze/moonbasic/releases/latest



Linux: if the app fails to start, ensure a normal desktop + GPU stack (runtime libs,

       not compiler -dev packages). Typical packages:



         Arch / CachyOS:  mesa  wayland  libxkbcommon  libglvnd

         Debian / Ubuntu: libgl1  libwayland-client0  libxkbcommon0



       You do **not** need `libenet`, `raylib`, or `jolt` packages — those are linked

       into `moonrun`. See docs/BUILDING.md only if you build from source.



  Older releases (before static ENet): if `ldd moonrun` shows `libenet.so.7 => not found`,

       either upgrade to a newer full-runtime zip, or on Arch/CachyOS install the `enet`

       package (`sudo pacman -S enet`).



Windows: run from a normal folder. If Windows reports a missing **non-system** DLL, you may

          have a partial copy, a mixed install, or an antivirus quarantine — re-extract the

          **entire** full-runtime zip from the same release.



Windows: "Entry Point Not Found" / nanosleep64 — usually from copying only `moonrun.exe`

          without the matching `moonbasic.exe` from the **same** zip, or from PATH picking up

          an older MinGW `libwinpthread` / `libgcc` DLL. Re-extract the full zip; do not drop

          stray MinGW DLLs next to the exes unless you know you need them.

