# Lights

| Designed | Implementation | Memory / notes |
|----------|----------------|----------------|
| **CreateLight (type, parent)** | **`Light.Create()`** (Pascal) / deprecated **`Light.Make()`** | **Heap handle** — **`Light.Free()`**; registry **`LIGHT.CREATE`** (deprecated **`LIGHT.MAKE`**). |
| **LightColor(id, r, g, b)** | **`Light.SetColor()`** | Registry **`LIGHT.SETCOLOR`**. |
| **LightRange(id, dist)** | **`Light.SetRange()`** | Registry **`LIGHT.SETRANGE`**. |
| **LightCone(id, inner, outer)** | **`Light.SetInnerCone()`**, **`Light.SetOuterCone()`** | Registry **`LIGHT.SETINNERCONE`**, **`LIGHT.SETOUTERCONE`** — [LIGHT.md](../LIGHT.md) |
| **LightPosition(id, x, y, z)** | **`Light.SetPos()`** | Registry **`LIGHT.SETPOS`** (deprecated **`LIGHT.SETPOSITION`**). |
| **LightPointAt(id, x, y, z)** | **`Light.SetDir()`** | Registry **`LIGHT.SETDIR`**. |
