moonBASIC — quick start (pre-built binaries)

==============================================

Download archives from: https://github.com/CharmingBlaze/moonbasic/releases/latest

Documentation and examples: https://github.com/CharmingBlaze/moonbasic



WHAT'S IN THIS FOLDER

---------------------

  moonbasic (or moonbasic.exe)  — Compiler: turn .mb source into .mbc bytecode, --check, --lsp,

                                 moonbasic new (scaffold a project)

  moonrun   (or moonrun.exe)     — Full game runtime: compile and run .mb / .mbc (graphics, physics, audio)



  For RUN TIME commands (playing a game), you need moonrun in this folder.

  For CHECK/COMPILE/LSP only, moonbasic alone is enough (compiler-only zip).



FIRST STEPS

-----------

  1. Extract this zip/tar anywhere you like.



  2. Open a terminal in that folder.



  3. Check that it works:

       Windows:   moonrun.exe --version

       Linux/macOS: ./moonrun --version



  4. Start a new game:

       moonbasic new MyGame

       cd MyGame

       moonrun main.mb



  5. Or run an existing script:

       moonrun path\to\yourgame.mb



  6. Lint without running:

       moonbasic --check path\to\yourgame.mb



  7. Compile to bytecode only:

       moonbasic path\to\yourgame.mb   → writes yourgame.mbc next to the source



TIPS

----

  • Language reference: docs/LANGUAGE.md at github.com/CharmingBlaze/moonbasic

  • Examples: examples/ folder in that repository

  • VS Code: download moonbasic-<tag>-vscode.vsix from the same Releases page

  • Full guide: docs/GETTING_STARTED.md

  • Engine source (contributors only): github.com/CharmingBlaze/moonbasic-compiler



Windows: re-extract the full zip if executables fail to start; keep both .exe files from the same release.



Linux: ensure GPU drivers and a normal desktop OpenGL stack are installed.
