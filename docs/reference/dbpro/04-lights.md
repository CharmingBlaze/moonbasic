# DBPro — Lights

moonBASIC: **`Light.Make()`**, **`Light.Set*()`** — [LIGHT.md](../LIGHT.md).

| DBPro | moonBASIC | Notes |
|-------|-----------|--------|
| **MAKE LIGHT (light, type)** | ✓ **`Light.Make()`** | Returns handle. |
| **DELETE LIGHT (light)** | ✓ **`Light.Free()`** | |
| **POSITION LIGHT (light, x, y, z)** | ✓ **`Light.SetPos()`** | |
| **ROTATE LIGHT (light, x, y, z)** | ✓ **`Light.SetDir()`** | |
| **SET LIGHT RANGE (light, dist)** | ✓ **`Light.SetRange()`** | |
| **SET LIGHT COLOR (light, r, g, b)** | ✓ **`Light.SetColor()`** | |
| **SET AMBIENT LIGHT (r, g, b)** | ✓ **`Render.SetAmbient()`** | |
| **SET LIGHT TYPE** | ≈ **`Light.Make()`** + parameters | |
| **SET LIGHT TO OBJECT** | ≈ follow entity in loop | Use entity position + **`Light.SetPos()`**. |
