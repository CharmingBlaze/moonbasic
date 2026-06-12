# Entity system — objects in your game world

> **Entities** are moonBASIC’s main way to represent things in the world: the player, enemies, crates, cameras-as-objects, and empty pivots for hierarchy.

**Namespaces:** `ENTITY` · `ENT` (shortcuts) · **Status:** Shipped · **Platform:** full runtime

**All `ENTITY` commands:** [COMMAND_REGISTRY.md#core-window-time](COMMAND_REGISTRY.md#core-window-time) (ENTITY section) · **Deep reference:** [reference/ENTITY.md](../../reference/ENTITY.md)

**Case:** `entity.setpos` = `ENTITY.SETPOS`

---

## Table of contents

- [At a glance](#at-a-glance)
- [When to use entities](#when-to-use-entities)
- [Choose the right tool](#choose-the-right-tool)
- [Core workflow](#core-workflow)
- [Creating entities](#creating-entities)
- [Transform and motion](#transform-and-motion)
- [Hierarchy](#hierarchy)
- [Visibility, names, and tags](#visibility-names-and-tags)
- [Drawing and updates](#drawing-and-updates)
- [Physics and animation hooks](#physics-and-animation-hooks)
- [Full example](#full-example)
- [Common mistakes](#common-mistakes)
- [Memory notes](#memory-notes)
- [See also](#see-also)

---

## At a glance

| Idea | Detail |
|------|--------|
| **You get** | Integer/handle IDs for game objects with position, rotation, scale, parent/child, tags |
| **You need first** | Window + loop ([00-START.md](../00-START.md)) |
| **Typical games** | 3D adventures, shooters, anything with “things” in a scene |
| **Not for** | Pure HUD text — use `DRAW.TEXT` / `GUI.*` ([UI-MENUS.md](UI-MENUS.md)) |

**Why entities exist:** The engine can **sort**, **parent**, **animate**, **sync physics**, and **draw** collections of objects if they share one model. Without entities you would track positions yourself and call low-level draw for each mesh.

---

## When to use entities

**Use when:**

- You have multiple objects that move, rotate, or get destroyed.
- You need parent/child attachments (weapon on hand, wheel on car).
- You want `SCENE.DRAW()` or `ENTITY.DRAWALL()` to render the world.
- You will attach physics or animation to objects.

**Skip when:**

- Single full-screen 2D demo with only `DRAW.RECTANGLE` — entities optional.
- One-off debug line — use `DEBUG.DRAWLINE`.

---

## Choose the right tool

| I want to… | Use | Not |
|------------|-----|-----|
| Place a box in 3D | `ENTITY.CREATECUBE` | Raw `MESH.CUBE` + manual matrix (unless custom mesh) |
| Load a character model | `ENTITY.LOAD` / `MODEL.LOAD` + `SETMODEL` | Drawing model handle without entity |
| Move every frame | `ENTITY.MOVE` / handle `.move` + `APP.DELTA()` | Teleport only when physics should integrate |
| “What did I click?” | `PICK.*` on entities ([COLLISION-3D.md](COLLISION-3D.md)) | `ENTITY` alone |
| 2D sprite actor | `SPRITE.*` ([07-2D-WORLD.md](../07-2D-WORLD.md)) | `ENTITY` for simple 2D unless you need 3D pipeline |

---

## Core workflow

1. **Create** — `ENTITY.CREATE`, `CREATECUBE`, `LOAD`, …  
   **Why:** Registers the object in the scene graph and allocates render/physics slots.

2. **Transform** — `SETPOS`, `SETROT`, `SETSCALE` (or `.pos()`, `.turn()`).  
   **Why:** Defines where the object appears before first draw.

3. **Optional setup** — parent, tags, material, physics (`ADDPHYSICS`).  
   **Why:** Gameplay queries (`HASTAG`) and simulation need metadata early.

4. **Each frame — simulate** — `ENTITY.UPDATE(dt)` or physics step + sync.  
   **Why:** Animation clips and internal lerps advance with time.

5. **Each frame — draw** — inside `RENDER.BEGIN`: `SCENE.DRAW()` or `ENTITY.DRAWALL()`.  
   **Why:** Drawing outside the 3D pass uses wrong matrices.

6. **Destroy** — `ENTITY.FREE` when removed permanently.  
   **Why:** Frees GPU/physics links; leaks in long sessions otherwise.

---

## Creating entities

| Command | Why use it |
|---------|------------|
| `ENTITY.CREATE(name)` | Empty pivot — logic anchor, camera rig, spawn point |
| `ENTITY.CREATECUBE(w,h,d)` | Fast placeholder or crate |
| `ENTITY.CREATESPHERE(r)` | Balls, orbs, simple props |
| `ENTITY.CREATEPLANE(w,h)` | Floor card, quad |
| `ENTITY.LOAD(path)` | GLB/OBJ character or prop |
| `ENTITY.COPY(ent)` | Duplicate with same setup |

**Example:**

```basic
; Hero from file; crate as primitive
hero = ENTITY.LOAD("assets/hero.glb")
crate = ENTITY.CREATECUBE(1, 1, 1)
ENTITY.SETPOS(crate, 3, 0, 0)
```

---

## Transform and motion

| Command | When | Why |
|---------|------|-----|
| `ENTITY.SETPOS(ent, x,y,z)` | Spawn, teleport, snap to ground | Sets absolute world (or local) position |
| `ENTITY.MOVE(ent, dx,dy,dz)` | Every frame with `APP.DELTA()` | Delta movement keeps speed consistent across FPS |
| `ENTITY.SETROT` / `TURN` | Aim, spin | Rotation in degrees |
| `ENTITY.SETSCALE` | Resize prop | Non-uniform scale for art |
| `ENTITY.LOOKAT` / `POINTAT` | AI, turrets | Aim without manual trig |
| `ENTITY.X/Y/Z(ent)` | Read position | Cheaper than unpacking full struct |

**Example — spin with delta time:**

```basic
WHILE NOT APP.SHOULDCLOSE()
    cube.turn(0, 90 * APP.DELTA(), 0)   ; 90 degrees per second on Y
    ; ... draw ...
WEND
```

**Handle chaining** (same commands): `hero.pos(0, 1, 0).turn(0, 45, 0)`

---

## Hierarchy

| Command | Why |
|---------|-----|
| `ENTITY.PARENT(child, parent)` | Child moves with parent (hand → body) |
| `ENTITY.PARENTCLEAR(child)` | Detach without destroying |
| `ENTITY.GETCHILD(parent, i)` | Iterate attached parts |

Parent transforms combine: child **local** offset rotates with parent **world** matrix.

```basic
hand = ENTITY.CREATECUBE(0.2, 0.2, 0.5)
ENTITY.PARENT(hand, hero)
ENTITY.SETPOS(hand, 0.5, 1.2, 0)   ; relative to hero
```

---

## Visibility, names, and tags

| Command | Why |
|---------|-----|
| `ENTITY.SHOW` / `HIDE` | Cutscenes, respawn without destroy |
| `ENTITY.SETNAME` | Debug, save files, editor labels |
| `ENTITY.SETTAG` / `HASTAG` | “enemy”, “pickup” — gameplay filters |

```basic
ENTITY.SETTAG(goblin, "enemy")
IF ENTITY.HASTAG(goblin, "enemy") THEN score = score + 10
```

---

## Drawing and updates

| Command | When to call | Why |
|---------|--------------|-----|
| `ENTITY.UPDATE(dt)` | Once per frame before draw | Advances animation / internal state |
| `ENTITY.DRAWALL()` | Inside `RENDER.BEGIN` … `END` | Renders sorted scene |
| `ENTITY.DRAW(ent)` | Single object highlight | Debug or special pass |
| `SCENE.DRAW()` | Active scene batch | Works with `SCENE.REGISTER` workflow |

**Order matters:**

```basic
RENDER.CLEAR(20, 22, 30)
RENDER.BEGIN(cam)
ENTITY.UPDATE(APP.DELTA())
SCENE.DRAW()          ; or ENTITY.DRAWALL()
RENDER.END()
RENDER.FRAME()
```

---

## Physics and animation hooks

| Command | Why |
|---------|-----|
| `ENTITY.ADDPHYSICS(ent)` | Jolt body tied to entity transform |
| `ENTITY.SETSTATIC(ent, true)` | Floors, walls — no motion |
| `ENTITY.PLAYANIM(ent, "Run")` | Skeletal clip on loaded model |
| `ENTITY.STOPANIM(ent)` | Cut animation |

3D physics detail: [COLLISION-3D.md](COLLISION-3D.md). Animation FSM: [reference/ANIM.md](../../reference/ANIM.md).

---

## Full example

```basic
APP.OPEN(960, 540, "Entity demo")
APP.SETFPS(60)

cam = CAMERA.CREATE()
CAMERA.SETACTIVE(cam)
CAMERA.SETPOS(cam, 0, 3, -8)
CAMERA.LOOKAT(cam, 0, 0, 0)

player = ENTITY.CREATECUBE(1, 1, 1)
ENTITY.SETTAG(player, "player")
ENTITY.SETPOS(player, 0, 0, 0)

sword = ENTITY.CREATECUBE(0.1, 0.1, 0.8)
ENTITY.PARENT(sword, player)
ENTITY.SETPOS(sword, 0.6, 0.5, 0)

WHILE NOT APP.SHOULDCLOSE()
    IF INPUT.KEYDOWN(KEY_D) THEN player.move(3 * APP.DELTA(), 0, 0)
    player.turn(0, 45 * APP.DELTA(), 0)

    ENTITY.UPDATE(APP.DELTA())
    RENDER.CLEAR(25, 28, 35)
    RENDER.BEGIN(cam)
    ENTITY.DRAWALL()
    RENDER.END()
    RENDER.FRAME()
WEND

ENTITY.FREE(sword)
ENTITY.FREE(player)
CAMERA.FREE(cam)
APP.CLOSE()
```

Run: `moonrun yourfile.mb` · Check: `moonbasic --check yourfile.mb`

---

## Common mistakes

| Mistake | Why it breaks | Fix |
|---------|---------------|-----|
| Draw without `RENDER.BEGIN` | Wrong camera matrices | Bracket 3D with `RENDER.BEGIN(cam)` |
| Forget `RENDER.FRAME()` | Frozen window | Call every loop iteration |
| Move without `DELTA()` | Speed depends on FPS | Multiply by `APP.DELTA()` |
| Never `ENTITY.FREE` | Handle leaks on respawn | Free when removing permanently |
| Parent after world `SETPOS` confusion | Child jumps | Set parent then local offset |
| `UPDATE` after draw | One-frame lag | `UPDATE` before `DRAWALL` |

---

## Memory notes

- `ENTITY.FREE(id)` when object is gone for good.
- `FREE.ALL` / `ERASE ALL` on major scene transitions — see [MEMORY.md](../../MEMORY.md).
- `ENTITY.COPY` duplicates state; still free copies when done.

---

## See also

- [01-CORE.md](../01-CORE.md) — overview
- [COLLISION-3D.md](COLLISION-3D.md) — physics on entities
- [reference/BEGINNER_FULL_STACK.md](../../reference/BEGINNER_FULL_STACK.md) — `NAVTO`, `SETHEALTH`, juice
