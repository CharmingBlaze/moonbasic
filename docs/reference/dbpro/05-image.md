# DBPro — Image (CPU)

moonBASIC: **`IMAGE.*`** (Raylib **Image** on CPU) — [IMAGE.md](../IMAGE.md).

| DBPro | moonBASIC | Notes |
|-------|-----------|--------|
| **LOAD IMAGE (file$, img)** | ✓ **`IMAGE.LOAD`**, **`IMAGE.LOADRAW`** | Returns **handle**. |
| **DELETE IMAGE** | ✓ **`IMAGE.FREE`** | |
| **SAVE IMAGE** | ✓ **`IMAGE.EXPORT`** | |
| **PASTE IMAGE** | ≈ **`IMAGE.DRAWIMAGE`**, blit to another image | |
| **GET IMAGE** | ≈ **`IMAGE` crop / copy** helpers | See manifest. |
| **SET IMAGE COLORKEY** / **TRANSPARENCY** | ≈ **`IMAGE.COLORREPLACE`**, **`ALPHACLEAR`**, etc. | |
| **IMAGE WIDTH** / **HEIGHT** | ✓ **`IMAGE.WIDTH`**, **`IMAGE.HEIGHT`** | |
