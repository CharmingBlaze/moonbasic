# Lights

| Designed | Implementation | Memory / notes |
|----------|----------------|----------------|
| **CreateLight (type, parent)** | **`LIGHT.MAKE`** + **`LIGHT.SET*`** for kind | **Heap handle** — **`LIGHT.FREE`**. |
| **LightColor** | **`LIGHT.SETCOLOR`** | |
| **LightRange** | **`LIGHT.SETRANGE`**, **`SETINTENSITY`** | |
| **LightCone** | **`LIGHT.SETINNERCONE`**, **`SETOUTERCONE`** | [LIGHT.md](../LIGHT.md) |
| **LightPosition** | **`LIGHT.SETPOSITION`** | |
| **LightPointAt** | **`LIGHT.SETTARGET`**, **`SETDIR`** | |
