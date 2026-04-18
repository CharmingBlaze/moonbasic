# MoonBASIC: Your First Hour

Welcome to the future of rapid game development. MoonBASIC is designed to be the fastest path from a "blank screen" to a "playable game," combining the readable soul of BASIC with a cutting-edge 3D engine and professional physics.

---

## 1. The MoonBASIC Philosophy

MoonBASIC isn't your grandfather's BASIC. It's a high-performance, object-oriented engine designed for the modern era.

> [!IMPORTANT]
> **Why MoonBASIC?**
> - **Zero Boilerplate:** No complex project setup. One `.mb` file is a complete game.
> - **Fluent API:** Write code that reads like a sentence using **Method Chaining**.
> - **Native Power:** Built on **Raylib** (Graphics) and **Jolt** (Physics) for AAA-grade stability.
> - **Beginner First:** The engine handles the "boring" parts (memory management, windowing) so you can focus on **fun**.

---

## 2. Language Primer (The Essentials)

If you can write a sentence, you can write MoonBASIC.

### Variables & Data
No keywords required. Just name it and assign it.
```basic
name    = "Hero"      ; Strings use double quotes
hp      = 100         ; Integers for whole numbers
speed   = 5.5         ; Floats for precision
isAlive = TRUE        ; Booleans (TRUE/FALSE)
```

### Clean Logic
Everything in MoonBASIC is designed for readability.
```basic
IF hp < 20 THEN
    PRINT "Danger! Health is low."
ENDIF

WHILE isAlive
    ; The game loop runs here
WEND
```

---

## 3. The Modern Touch: Fluent Chaining

This is where MoonBASIC outshines legacy engines. Most commands return the object they modified, allowing you to "chain" configurations together.

> [!TIP]
> **Fluent vs. Procedural**
> Instead of repeating `ENTITY.SET...` over and over, you can do it in one line. It's faster to write and easier to read.

```basic
; ❌ The "Old" Clunky Way
hero = ENTITY.LOAD("hero.glb")
ENTITY.SETPOS(hero, 0, 10, 0)
ENTITY.SETCOLOR(hero, 255, 0, 0)
ENTITY.SETALPHA(hero, 0.5)

; ✅ The Modern MoonBASIC Way
hero = ENTITY.LOAD("hero.glb").SETPOS(0, 10, 0).SETCOLOR(255, 0, 0).SETALPHA(0.5)
```

---

## 4. Build a Game in 15 Lines

Let's build a functional 2D prototype with frame-rate independent movement.

```basic
WINDOW.OPEN(1280, 720, "MoonBASIC Move")
WINDOW.SETFPS(60)

x = 640 : y = 360  ; Starting coordinates

WHILE NOT WINDOW.SHOULDCLOSE()
    dt = TIME.DELTA() ; Get the time since the last frame
    
    ; Rapid movement using Input Axis
    dx = INPUT.AXIS(KEY_A, KEY_D)
    dy = INPUT.AXIS(KEY_W, KEY_S)
    
    x = x + dx * 400 * dt
    y = y + dy * 400 * dt
    
    RENDER.CLEAR(30, 40, 50)
    DRAW.CIRCLE(x, y, 40, 255, 200, 0, 255)
    RENDER.FRAME()
WEND

WINDOW.CLOSE()
```

---

## 5. Transitioning to 3D

Going 3D in MoonBASIC is exactly like 2D, just with an extra dimension.

```basic
WINDOW.OPEN(1280, 720, "3D World")
cam = CAMERA.CREATE().SETPOS(0, 5, 10).SETTARGET(0, 0, 0)
cube = ENTITY.CREATECUBE(2.0).SETCOLOR(100, 200, 255, 255)

WHILE NOT WINDOW.SHOULDCLOSE()
    dt = TIME.DELTA()
    
    ; Rotate the cube over time
    cube.SETROT(0, TIME.GET() * 50, 0)
    
    RENDER.CLEAR(15, 15, 25)
    RENDER.BEGIN3D(cam)
        ENTITY.DRAWALL()
        DRAW3D.GRID(50, 1.0)
    RENDER.END3D()
    RENDER.FRAME()
WEND
```

---

## 6. Prototyping with "Gameplay Helpers"

MoonBASIC includes "Full Stack" helpers to solve common game problems instantly.

- **`ENTITY.NAVTO`**: Pathfind to a destination.
- **`ENTITY.WITHINRADIUS`**: Check for proximity triggers.
- **`CAMERA.ORBIT`**: Professional 3rd-person camera in one line.
- **`ENTITY.DAMAGE`**: Built-in health and combat logic.

> [!NOTE]
> Check out the **[Beginner Full Stack Guide](reference/BEGINNER_FULL_STACK.md)** to see these in action.

---

## Your Journey Continues

You've mastered the basics. Now, explore the specialized systems that make MoonBASIC a powerhouse:

1. **[Animation 3D](reference/ANIMATION_3D.md)**: Bring your characters to life with skeletal animation.
2. **[Physics Engine](reference/PHYSICS3D.md)**: Add weight and impact using the Jolt solver.
3. **[Modern Rendering](reference/RENDER.md)**: Unlock PBR, Bloom, and SSAO for AAA visuals.
4. **[Multiplayer](reference/NETWORK.md)**: Build online games with the ENet networking stack.

**Welcome to the MoonBASIC community. Now go build something amazing!**
