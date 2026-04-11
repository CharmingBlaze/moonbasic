# Lights

| Designed | Implementation | Memory / notes |
|----------|----------------|----------------|
| **CreateLight (type, parent)** | **`Light.Make()`** | **Heap handle** — **`Light.Free()`**. |
| **LightColor(id, r, g, b)** | **`Light.SetColor()`** | |
| **LightRange(id, dist)** | **`Light.SetRange()`** | |
| **LightCone(id, inner, outer)** | **`Light.SetInnerCone()`**, **`Light.SetOuterCone()`** | [LIGHT.md](../LIGHT.md) |
| **LightPosition(id, x, y, z)** | **`Light.SetPos()`** | |
| **LightPointAt(id, x, y, z)** | **`Light.SetDir()`** | |
