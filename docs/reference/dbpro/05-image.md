# DBPro — Image (CPU)

moonBASIC: **`IMAGE.*`** (Raylib **Image** on CPU) — [IMAGE.md](../IMAGE.md).

| DBPro | moonBASIC | Notes |
|-------|-----------|--------|
| **LOAD IMAGE** | ✓ **`Image.Load()`** | Returns **handle**. |
| **DELETE IMAGE** | ✓ **`Image.Free()`** | |
| **SAVE IMAGE** | ✓ **`Image.Export()`** | |
| **PASTE IMAGE** | ≈ **`Draw.Texture()`** | |
| **IMAGE WIDTH** | ✓ **`Image.Width()`** | |
| **IMAGE HEIGHT** | ✓ **`Image.Height()`** | |
| **ROTATE IMAGE** | ✓ **`Image.Rotate()`** | |
| **MIRROR IMAGE** | ✓ **`Image.FlipH()`** | |
| **FLIP IMAGE** | ✓ **`Image.FlipV()`** | |
| **GET IMAGE** | ≈ **`Image.Copy()`** / **`Image.Crop()`** | |
| **SET IMAGE COLORKEY** | ≈ **`Image.ColorReplace()`** / **`AlphaClear()`** | |
