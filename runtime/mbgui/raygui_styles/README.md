# Bundled raygui `.rgs` styles

Binary style files are copied from the official [raygui `styles/`](https://github.com/raysan5/raygui/tree/master/styles) tree (same names as upstream folders: `style_<name>.rgs` saved as `<name>.rgs`).

They are loaded at runtime via `GUI.THEMEAPPLY` (embedded in the moonBASIC binary when built with CGO). Upstream notes: raylib **5.5** and raygui **4.5**+; see the [raygui styles README](https://github.com/raysan5/raygui/blob/master/styles/README.md).

Licensing follows [raylib / raygui](https://github.com/raysan5/raygui) (zlib); font licenses per style live in the upstream style folders.
