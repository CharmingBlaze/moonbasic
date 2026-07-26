moonBASIC IDE — unzip and code (Windows / Linux / macOS)
========================================================

Everything you need is in THIS folder. No Go, Node, VS Code, or hunting for DLLs.


START IN 30 SECONDS
-------------------

  1. Extract this archive somewhere permanent (Desktop, Documents, Projects…).

  2. Start the IDE:

       Windows:  double-click  START-IDE.bat
                 (or moonbasic-ide.exe)

       macOS:    double-click  START-IDE.command
                 If Gatekeeper blocks it: right-click → Open
                 (Terminal alternative: ./START-IDE.sh)

       Linux:    chmod +x START-IDE.sh moonbasic-ide moonbasic moonrun
                 ./START-IDE.sh

  3. Status bar should say "Toolchain ready".

  4. File → Open Samples Folder → open hello.mb → press F5.

That's it. moonbasic and moonrun are already beside the IDE.



WINDOWS NOTES
-------------

  Keep every file from this zip together (including any lib*.dll next to moonrun.exe).
  Do not mix moonrun from an older release with a newer IDE folder.

WHAT'S IN THIS FOLDER
---------------------

  START-IDE.* / START-IDE.command   Launch the editor
  moonbasic-ide (.exe / bare / .app) Desktop IDE (docs also built into the app)
  moonbasic (.exe)                  Compiler + LSP (--check, .mb → .mbc)
  moonrun (.exe)                    Game runtime (F5)
  docs/                             Full documentation (Begin Here, guides, API)
  samples/                          hello.mb, spin_cube.mb — try these first
  ADD-TO-PATH.*                     Optional: use moonbasic/moonrun in any terminal
  README-IDE-RELEASE.txt            This file


KEYBOARD
--------

  F5                Run game (moonrun)
  Shift+F5          Stop game
  Ctrl/⌘+Shift+C    Check syntax
  Ctrl/⌘+Shift+B    Compile to .mbc
  Alt+H             Help at cursor

  macOS: use ⌘ (Command) instead of Ctrl for editor shortcuts.


TERMINAL (OPTIONAL)
-------------------

  The IDE does not need PATH. For a shell in any folder:

       Windows:  double-click ADD-TO-PATH.bat  (then open a NEW terminal)
       Linux/macOS:  ./ADD-TO-PATH.sh

  Then:

       moonbasic new MyGame
       cd MyGame
       moonrun main.mb


IF SOMETHING'S WRONG
--------------------

  Status "No toolchain"
    → Keep moonbasic and moonrun in THIS same folder (don't move only the IDE).
    → Or: File → Settings → Toolchain → Browse…

  Status "moonbasic only — F5 needs moonrun"
    → Copy moonrun from this same release into the IDE folder.

  IDE won't start
    Windows: install WebView2 Evergreen Runtime (usually already on Win 10/11).
    Linux:   install WebKitGTK + GTK3 (OS packages, not moonBASIC sidecars):
               Arch / CachyOS:  sudo pacman -S webkit2gtk-4.1 gtk3 mesa wayland libxkbcommon
               Debian / Ubuntu: sudo apt install libwebkit2gtk-4.1-0 libgtk-3-0 libgl1
    macOS:   right-click START-IDE.command → Open (Gatekeeper).


MORE
----

  Docs inside the IDE: Documentation → BEGIN_HERE.md / GETTING_STARTED.md
  More examples: https://github.com/CharmingBlaze/moonbasic/tree/main/examples
  Releases: https://github.com/CharmingBlaze/moonbasic/releases
  Engine source: https://github.com/CharmingBlaze/moonbasic-compiler
