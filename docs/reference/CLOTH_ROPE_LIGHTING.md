# Cloth-style effects, rope, and lighting (beginner bridge)

This page ties together **projectiles**, **2D rope/bridge chains**, **3D lights + ambient**, and what we **do not** simulate as full cloth/soft-body yet.

---

## Projectiles

### 3D — `Entity.Shoot()`

`Entity.Shoot(prefab, speed, lifetime [, shape])` clones the prefab, clears duplicated physics, aligns to the shooter, builds a **Jolt** body with **CCD** (continuous collision), sets velocity along **pitch/yaw**, and uses **`Entity.DestroyAfter()`**. See [PROJECTILES.md](PROJECTILES.md).

### 2D — `Body2D.Shoot()`

`Body2D.Shoot(shooter, speed, lifetime [, radius])` spawns a **Box2D** dynamic circle with **`SetBullet(true)`** for tunneling resistance, velocity along the shooter **angle**, and an internal **`Body2D.Free()`** schedule.

**Filtering:** neither path yet applies automatic **“ignore shooter”** physics filters end-to-end; use gameplay rules, thicker geometry, or layers when Jolt/Box2D filters are fully wired.

---

## Rope / bridge (2D)

`Physics2D.CreateRope(x1, y1, x2, y2, segments, mode)` returns a **handle array** of **`Body2D`** instances: small **static** anchors plus **dynamic** links joined with **distance joints**.

- **`mode = "bridge"`** — static anchors at **both** ends (good for platforms).
- **`mode = "rope"`** — pinned at the **start** only; the chain dangles toward `(x2, y2)`.

Free **every** body handle from the array when you are done (joints are destroyed with the bodies). Details: [PHYSICS2D.md](PHYSICS2D.md).

---

## Lighting

### Ambient

- **`Render.SetAmbient(r, g, b [, a])`** — hemispheric-style **ambient tint** for the shared **PBR** path.
- **`World.SetAmbient(r, g, b [, a])`** — **alias**, same overloads (reads like “global lighting” in tutorials).

### Directional sun + extras

Constructors such as **`Light.Make("directional")`**, **`Light.Make("point")`**, **`Light.Make("spot")`** allocate **CPU light objects**. The **stock PBR fragment shader** in this repo shades with **one directional** vector plus **ambient**; see [LIGHT.md](LIGHT.md) and [CAMERA_LIGHT_RENDER.md](CAMERA_LIGHT_RENDER.md) for how sun + shadow interact.

### Parenting a light to an entity

- **`Light.SetParent(light, id)`** — for **point** and **spot** lights, **world position** follows the entity each frame. **Spot direction** is not automatically aligned to the camera; update direction or targets in script if you want a flashlight cone to track view yaw.

There is **no** separate GPU “8 lights” array in the default PBR shader here—extra handles are for API completeness and future/custom shading—not a hard cap error like classic multi-light demos.

---

## 3D “cape / tail” (faux cloth)

Full **cloth** or **Jolt soft bodies** wired to skinned meshes are **not** exposed as a one-liner yet. The usual **indie pattern** is:

1. A short **chain of rigid bodies** (spheres/capsules) constrained behind the character.
2. A **ribbon** or **strip** drawn through those positions each frame (`Draw3D.Line()`, triangle strip, or a thin mesh).

When/if a dedicated **`Entity.AttachCape()`**-style command lands, it would wrap those two steps—not a full FEM cloth solve.

---

## See also

- [GAMEHELPERS.md](GAMEHELPERS.md) — overview table  
- [PHYSICS2D.md](PHYSICS2D.md), [PHYSICS3D.md](PHYSICS3D.md)  
- [LIGHT.md](LIGHT.md), [PROJECTILES.md](PROJECTILES.md)
