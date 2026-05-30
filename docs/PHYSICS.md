# Getting Started with 3D Physics

MoonBASIC features a high-performance 3D physics engine powered by **Jolt Physics**. This guide will help you set up your first simulation and understand how physical bodies interact with your game world.

> [!NOTE]
> For technical details on the engine's internal dual-path architecture and build requirements, see [docs/architecture/PHYSICS_INTERNAL.md](architecture/PHYSICS_INTERNAL.md).

---

## 1. Setting up the World

Before you can create any physical objects, you must initialize the physics world. This is usually done once at the start of your game.

```basic
; Initialize the physics world
PHYSICS3D.START()

; Set the global gravity (X, Y, Z)
; 0, -10, 0 is typical earth-like gravity
PHYSICS3D.SETGRAVITY(0, -9.81, 0)
```

In your main game loop, you must call `PHYSICS3D.UPDATE()` to advance the simulation.

```basic
WHILE NOT WINDOW.SHOULDCLOSE()
    PHYSICS3D.UPDATE()
    ; ... your game logic ...
    RENDER.FRAME()
WEND
```

---

## 2. Bodies and Shapes

In MoonBASIC, physical objects are called **Bodies**. To create a body, you follow a "Build-and-Commit" pattern:

1. **Create a Definition**: Tell the engine if the body is `STATIC` (unmoving), `DYNAMIC` (responds to forces), or `KINEMATIC` (moved by script).
2. **Add Shapes**: Attach boxes, spheres, or capsules to the definition.
3. **Commit**: Insert the finished body into the world.

### Example: A Falling Ball
```basic
; 1. Create a dynamic definition
ballDef = BODY3D.CREATE("DYNAMIC")

; 2. Add a sphere shape (radius 1.0)
BODY3D.ADDSPHERE(ballDef, 1.0)

; 3. Commit to the world at height Y=10
ball = BODY3D.COMMIT(ballDef, 0, 10, 0)
```

---

## 3. Linking Physics to Visuals

By default, physical bodies are invisible. To see them, you should link them to a visual **Entity** (like a 3D Model or a Primitive).

The command `ENTITY.LINKPHYSBUFFER` automatically syncs the position and rotation of a physics body to your model every frame.

```basic
; Create a visual sphere model
ballModel = MODEL.CREATESPHERE(1.0)
ballEnt   = ENTITY.CREATE(ballModel)

; Link it to our physical body
ENTITY.LINKPHYSBUFFER(ballEnt, ball.bufferIndex())
```

---

## 4. Forces and Impulses

You can move dynamic bodies by applying forces or impulses:

- **ApplyForce**: Continuous pressure (like a rocket engine).
- **ApplyImpulse**: A sudden hit (like a bat hitting a ball).

```basic
; Jump! Apply an upward impulse
IF Input.KeyHit(KEY_SPACE) THEN ball.applyImpulse(0, 5, 0)
```

---

## 5. Collision Events

To detect when objects hit each other, use the `ONCOLLISION` callback.

```basic
; Register a callback when the ball hits the floor
PHYSICS3D.ONCOLLISION(ball, floor, "OnBallHitFloor")

SUB OnBallHitFloor(hA, hB)
    PRINT "The ball bounced!"
END SUB
```

---

## Next Steps

- **[PHYSICS3D Reference](reference/PHYSICS3D.md)** — Full list of all physics commands.
- **[Character Physics](reference/CHARACTER_PHYSICS.md)** — How to make players walk, jump, and climb.
- **[Advanced Physics](reference/PHYSICS_ADVANCED.md)** — Joints, vehicles, and complex constraints.
