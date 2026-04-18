# Getting Started with MoonBASIC

Welcome to MoonBASIC. Whether you are installing the engine for the first time or writing your first lines of code, this guide will get you up and running in minutes.

> [!TIP]
> **New to game development?**
> Start with **[MoonBASIC: Your First Hour](FIRST_HOUR.md)** for a friendly introduction to the language, modern **Method Chaining**, and rapid prototyping.

---

## 1. Installation

### Pre-built Binaries (Recommended)
Download the latest archive for your platform from **[GitHub Releases](https://github.com/CharmingBlaze/moonbasic/releases/latest)**.

| OS | Archive |
|----|---------|
| **Windows** | `moonbasic-v*-windows-amd64.zip` |
| **Linux** | `moonbasic-v*-linux-amd64.tar.gz` |

**Contents:**
- `moonbasic` / `moonbasic.exe`: The compiler (`.mb` → `.mbc`).
- `moonrun` / `moonrun.exe`: The game runtime (executes your code).

### Building from Source
If you prefer to build from source, you will need **Go 1.25.3+** and a **C toolchain** (MinGW-w64 for Windows, GCC for Linux).

> [!IMPORTANT]
> **CGO is Required**
> MoonBASIC relies on CGO for its high-performance graphics and physics engines. Ensure `CGO_ENABLED=1` is set in your environment.

**Windows Build:**
```bat
git clone https://github.com/CharmingBlaze/moonbasic
cd moonbasic
set CGO_ENABLED=1
go build -o moonbasic.exe .
go build -tags fullruntime -o moonrun.exe ./cmd/moonrun
```

---

## 2. Your First Program

Create a file named `hello.mb`:
```basic
PRINT "Hello, MoonBASIC!"
```

Run it using the runtime:
```bash
moonrun hello.mb
```

---

## 3. Opening a Window

MoonBASIC makes window management effortless. Create `display.mb`:

```basic
WINDOW.OPEN(1280, 720, "MoonBASIC Window")
WINDOW.SETFPS(60)

WHILE NOT WINDOW.SHOULDCLOSE()
    RENDER.CLEAR(30, 40, 60)
    DRAW.TEXT("Press ESC to exit", 540, 350, 20, 255, 255, 255, 255)
    RENDER.FRAME()
WEND

WINDOW.CLOSE()
```

---

## 4. Modern 3D with Method Chaining

MoonBASIC supports **Method Chaining** (Fluent API), allowing you to configure objects in a single, readable line.

```basic
WINDOW.OPEN(1280, 720, "3D Cube Demo")
cam = CAMERA.CREATE().SETPOS(0, 5, 10).SETTARGET(0, 0, 0)
cube = ENTITY.CREATECUBE(2.0).SETCOLOR(100, 200, 255, 255)

WHILE NOT WINDOW.SHOULDCLOSE()
    ; Update rotation using a fluent method
    cube.SETROT(0, TIME.GET() * 50, 0)

    RENDER.CLEAR(10, 10, 20)
    RENDER.BEGIN3D(cam)
        ENTITY.DRAWALL()
        DRAW3D.GRID(50, 1.0)
    RENDER.END3D()
    RENDER.FRAME()
WEND
```

---

## 5. Modern Blitz-Style (High Fidelity)

For advanced users, MoonBASIC provides a "High Fidelity" path with PBR materials, dynamic lighting, and SSAO.

```basic
WINDOW.OPEN(1920, 1080, "Project: High Fidelity")
cam = CAMERA.CREATE().SETPOS(0, 5, 10)
sun = LIGHT.CREATEDIRECTIONAL(0, -1, 0, 255, 255, 200, 2.0)

; Load a high-poly model with modern effects
car = ENTITY.LOADMESH("supercar.glb").SETPBR(0.9, 0.1)
RENDER.SETSSAO(TRUE)
RENDER.SETBLOOM(0.8)

WHILE NOT WINDOW.SHOULDCLOSE()
    CAMERA.FOLLOWENTITY(cam, car, 10.0, 3.0, 5.0)
    
    ENTITY.UPDATE(TIME.DELTA())

    RENDER.CLEAR(12, 14, 22)
    RENDER.BEGIN3D(cam)
        ENTITY.DRAWALL()
    RENDER.END3D()
    RENDER.FRAME()
WEND
```

---

## Next Steps

Explore the specialized documentation to master every aspect of the engine:

| Topic | Reference |
|-------|-----------|
| **Core Workflow** | [Programming Guide](PROGRAMMING.md) |
| **Language Syntax** | [Language Reference](LANGUAGE.md) |
| **3D Entities** | [Entity Reference](reference/ENTITY.md) |
| **Physics** | [Physics 3D Reference](reference/PHYSICS3D.md) |
| **Atmosphere** | [Camera & Render Reference](reference/CAMERA_LIGHT_RENDER.md) |
| **Gameplay Helpers** | [Beginner Full Stack](reference/BEGINNER_FULL_STACK.md) |

**Happy Coding!**
