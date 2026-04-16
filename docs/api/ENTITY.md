# Entity Commands

Commands for creating and managing game entities — the primary building block of any moonBASIC 3D game. Entities combine a 3D position, rotation, scale, model mesh, collision shape, physics body, and rendering state into a single handle.

## Core Concepts

- **Entity** — A game object with position, rotation, scale, model, color, alpha, collision type, and optional physics body.
- **Collision types** — Integer tags (1–255) assigned to entities. Collisions are checked between type pairs (e.g., "player" type 1 vs "enemy" type 2).
- **Entity.Update(dt)** — The per-frame entity tick. Moves entities by their velocity, applies gravity, checks collisions, and syncs with physics bodies.
- **Entity.DrawAll()** — Renders all visible entities in a single call. Must be inside a `Camera.Begin` / `Camera.End` block.
- Entities are **heap handles** and must be freed when removed from the game.

---

## Creation

### `Entity.Create()` / `Entity.Make()`

Creates an empty entity at the origin with no model.

**Returns:** `handle`

```basic
ent = Entity.Create()
```

---

### `Entity.CreateBox(w, h, d, r, g, b)` / `Entity.CreateCube(w, h, d, r, g, b)`

Creates an entity with a built-in colored box mesh.

- `w`, `h`, `d` (float) — Width, height, depth.
- `r`, `g`, `b` (int) — Color (0–255).

**Returns:** `handle`

```basic
wall = Entity.CreateBox(10, 3, 0.5, 128, 128, 128)
```

---

## Position & Transform

### `Entity.SetPos(entityHandle, x, y, z)` / `ent.pos(x, y, z)`

Sets the entity's world position.

- `x`, `y`, `z` (float) — World coordinates.

```basic
Entity.SetPos(player, 0, 1, 0)
; or
player.pos(0, 1, 0)
```

---

### `Entity.GetPos(entityHandle)` / `Entity.GetPosition(entityHandle)`

Returns the entity's position. Access X, Y, Z via return values.

---

### `Entity.GetXZ(entityHandle)`

Returns only the X and Z components (useful for top-down distance checks).

---

### `Entity.Rotate(entityHandle, pitch, yaw, roll)` / `Entity.Turn(entityHandle, pitch, yaw, roll)`

Sets the entity's rotation in Euler degrees.

- `pitch` (float) — X-axis rotation.
- `yaw` (float) — Y-axis rotation.
- `roll` (float) — Z-axis rotation.

---

### `Entity.GetRot(entityHandle)`

Returns the entity's rotation.

---

### `Entity.Scale(entityHandle, sx, sy, sz)` / `ent.scale(sx, sy, sz)`

Sets the entity's scale.

---

### `Entity.GetScale(entityHandle)`

Returns the entity's scale.

---

### `Entity.Move(entityHandle, forward, right, up)`

Moves the entity relative to its orientation. Positive `forward` moves in the direction the entity faces.

- `forward` (float) — Forward/backward.
- `right` (float) — Left/right strafe.
- `up` (float) — Up/down.

**How it works:** Computes a direction vector from the entity's current yaw and applies the movement. This is the primary way to move entities in gameplay code.

```basic
; Move entity forward at 5 units/sec
Entity.Move(player, 5 * dt, 0, 0)
```

---

### `Entity.Push(entityHandle, vx, vy, vz)`

Pushes the entity by an absolute world-space vector (ignores orientation).

---

### `Entity.Translate(entityHandle, dx, dy, dz)`

Translates the entity by a delta in world space.

---

### `Entity.Jump(entityHandle, force)`

Applies an upward velocity to the entity (for simple platformer physics).

- `force` (float) — Jump velocity.

---

## Extended Creation

### `Entity.CreateSphere(segments, r, g, b)`

Creates an entity with a built-in sphere mesh.

- `segments` (int) — Sphere tessellation (8–32 typical).
- `r`, `g`, `b` (int) — Color.

**Returns:** `handle`

---

### `Entity.CreateCylinder(segments, r, g, b)`

Creates an entity with a built-in cylinder mesh.

---

### `Entity.CreatePlane(divs, r, g, b)`

Creates an entity with a flat plane mesh.

---

### `Entity.CreateMesh(meshHandle)`

Creates an entity from a raw mesh handle.

---

### `Entity.Load(filePath)` / `Entity.LoadMesh(filePath)` / `LoadMesh(filePath)`

Creates an entity by loading a 3D model file (`.glb`, `.gltf`, `.obj`, `.iqm`, `.vox`).

- `filePath` (string) — Path to model file.

**Returns:** `handle`

```basic
hero = Entity.Load("assets/hero.glb")
Entity.SetPos(hero, 0, 0, 0)
```

---

### `Entity.LoadAnimatedMesh(filePath)`

Loads a model with skeletal animation data.

**Returns:** `handle`

---

## Appearance

### `Entity.Color(entityHandle, r, g, b)` / `ent.col(r, g, b)`

Sets the entity's tint color.

- `r`, `g`, `b` (int) — Color (0–255).

---

### `Entity.GetColor(entityHandle)`

Returns the entity's color.

---

### `Entity.Alpha(entityHandle, alpha)` / `EntityAlpha(entityHandle, alpha)`

Sets the entity's alpha transparency. 0.0 = invisible, 1.0 = opaque.

- `alpha` (float) — Transparency value.

---

### `Entity.GetAlpha(entityHandle)`

Returns the entity's alpha value.

---

### `Entity.Shininess(entityHandle, shininess)` / `EntityShininess(entityHandle, shininess)`

Sets the specular shininess of the entity's material.

- `shininess` (float) — Shininess power (higher = tighter highlight).

---

### `Entity.Texture(entityHandle, textureHandle)` / `EntityTexture(entityHandle, textureHandle)`

Applies a texture to the entity's mesh.

- `textureHandle` (handle) — Texture loaded with `Texture.Load`.

```basic
brickTex = Texture.Load("assets/brick.png")
Entity.Texture(wall, brickTex)
```

---

### `Entity.FX(entityHandle, flags)`

Sets rendering effect flags on the entity (e.g., fullbright, vertex colors, flat shading).

- `flags` (int) — Bitmask of FX flags (Blitz3D compatible).

---

### `Entity.Blend(entityHandle, mode)` / `EntityBlend(entityHandle, mode)`

Sets the blend mode for the entity. 1 = Alpha, 2 = Multiply, 3 = Additive.

---

### `Entity.Order(entityHandle, order)`

Sets the render order for the entity. Lower values render first. Used for transparency sorting.

---

### `Entity.SetCullMode(entityHandle, mode)`

Sets the face culling mode for the entity's mesh rendering.

- `mode` (int) — 0 = back-face culling (default), 1 = front-face, 2 = none.

---

## Collision

### `Entity.Type(entityHandle, typeID)` / `EntityType(entityHandle, typeID)`

Assigns a collision type ID to an entity.

- `typeID` (int) — Collision type (1–255).

**How it works:** Entities only collide with other entities when `Collisions(typeA, typeB, method)` has been called to define interaction.

```basic
Entity.Type(player, 1)
Entity.Type(enemy, 2)
Entity.Type(bullet, 3)
```

---

### `Entity.Radius(entityHandle, radius)` / `EntityRadius(entityHandle, radius)`

Sets the entity's collision sphere radius.

- `radius` (float) — Collision radius.

```basic
Entity.Radius(player, 0.5)
```

---

### `Entity.Box(entityHandle, x, y, z, w, h, d)`

Sets an axis-aligned bounding box for collision instead of a sphere.

---

### `Collisions(typeA, typeB, method)`

Enables collision checking between two entity types.

- `typeA`, `typeB` (int) — Entity type IDs.
- `method` (int) — Collision method (1 = sphere-sphere, 2 = sphere-box, 3 = box-box).

```basic
Collisions(1, 2, 1)   ; Player (1) vs Enemy (2), sphere-sphere
Collisions(1, 3, 2)   ; Player (1) vs Walls (3), sphere-box
```

---

### `EntityCollided(entityHandle, typeID)`

Returns `TRUE` if the entity collided with any entity of the given type this frame.

- `typeID` (int) — Type to check against.

**Returns:** `bool`

```basic
IF EntityCollided(player, 2) THEN
    PRINT "Hit an enemy!"
ENDIF
```

---

### `EntityHitsType(entityHandle, typeID)`

Returns `TRUE` if the entity hit an entity of the given type (edge-triggered).

---

### `CountCollisions(entityHandle)`

Returns the number of collisions the entity had this frame.

**Returns:** `int`

---

### `GetCollisionEntity(entityHandle, index)`

Returns the handle of the entity involved in collision at the given index.

**Returns:** `handle`

---

### `ResetEntity(entityHandle)`

Resets the entity's collision state. Call after teleporting an entity to prevent ghost collisions.

---

## Visibility

### `Entity.Visible(entityHandle, visible)` / `Entity.SetVisible(entityHandle, visible)` / `EntityVisible(entityHandle, visible)`

Sets whether an entity is visible.

- `visible` (bool) — `TRUE` to show, `FALSE` to hide.

---

### `HideEntity(entityHandle)` / `ShowEntity(entityHandle)`

Easy Mode aliases for hiding/showing entities.

---

## Parenting

### `Entity.Parent(childHandle, parentHandle)`

Makes an entity a child of another. The child's transform becomes relative to the parent.

```basic
; Attach weapon to player hand
Entity.Parent(sword, player)
Entity.SetPos(sword, 0.5, 0, 0)   ; Offset relative to player
```

---

### `Entity.ParentClear(entityHandle)` / `Entity.Unparent(entityHandle)`

Removes the entity from its parent (detaches to world space).

---

### `Entity.CountChildren(entityHandle)`

Returns the number of child entities.

**Returns:** `int`

---

### `Entity.GetChild(entityHandle, index)`

Returns a child entity by index.

**Returns:** `handle`

---

### `Entity.FindChild(entityHandle, name)`

Finds a child entity by name (from the model hierarchy).

**Returns:** `handle`

---

## Spatial Queries

### `Entity.EntityX(entityHandle)` / `EntityX(entityHandle)`

Returns the entity's X position. **Returns:** `float`

### `Entity.EntityY(entityHandle)` / `EntityY(entityHandle)`

Returns the entity's Y position. **Returns:** `float`

### `Entity.EntityZ(entityHandle)` / `EntityZ(entityHandle)`

Returns the entity's Z position. **Returns:** `float`

### `Entity.EntityPitch(entityHandle)` / `EntityPitch(entityHandle)`

Returns the entity's pitch (X rotation) in degrees. **Returns:** `float`

### `Entity.EntityYaw(entityHandle)` / `EntityYaw(entityHandle)`

Returns the entity's yaw (Y rotation) in degrees. **Returns:** `float`

### `Entity.EntityRoll(entityHandle)` / `EntityRoll(entityHandle)`

Returns the entity's roll (Z rotation) in degrees. **Returns:** `float`

---

### `Entity.Distance(entityA, entityB)`

Returns the 3D distance between two entities.

**Returns:** `float`

```basic
dist = Entity.Distance(player, enemy)
IF dist < 2.0 THEN
    PRINT "Too close!"
ENDIF
```

---

### `Entity.InView(entityHandle)`

Returns `TRUE` if the entity is within the camera's view frustum.

**Returns:** `bool`

---

### `Entity.TFormPoint(x, y, z, srcEntity, dstEntity)`

Transforms a point from one entity's local space to another's (or to world space if dstEntity is 0).

---

### `Entity.TFormVector(vx, vy, vz, srcEntity, dstEntity)` / `TFormVector(...)`

Transforms a direction vector between entity coordinate spaces.

---

### `Entity.DeltaX(entityHandle)` / `Entity.DeltaY(entityHandle)` / `Entity.DeltaZ(entityHandle)`

Returns the per-frame position delta (how far the entity moved this frame) on each axis.

**Returns:** `float`

---

### `Entity.MatrixElement(entityHandle, row, col)`

Returns a specific element from the entity's 4×4 transform matrix.

**Returns:** `float`

---

### `PointEntity(entityHandle, targetHandle)`

Rotates an entity to face another entity (look-at).

```basic
PointEntity(turret, player)
```

---

## Velocity & Forces

### `Entity.Velocity(entityHandle, vx, vy, vz)`

Sets the entity's velocity vector directly.

- `vx`, `vy`, `vz` (float) — Velocity components.

---

### `Entity.AddForce(entityHandle, fx, fy, fz)`

Adds a force to the entity (accumulated, applied during `Entity.Update`).

---

### `Entity.Slide(entityHandle, enabled)`

Enables or disables sliding collision response (entity slides along walls instead of stopping).

- `enabled` (bool) — `TRUE` for slide behavior.

---

## Collision Detail Queries

### `Entity.CollisionX(entityHandle, index)` / `Entity.CollisionY(...)` / `Entity.CollisionZ(...)`

Returns the world-space position of a collision contact point.

- `index` (int) — Collision index (from `CountCollisions`).

**Returns:** `float`

---

### `Entity.CollisionNX(entityHandle, index)` / `Entity.CollisionNY(...)` / `Entity.CollisionNZ(...)`

Returns the collision surface normal vector at a contact point.

**Returns:** `float`

```basic
; Detect wall vs floor collision by checking normal Y
ny = Entity.CollisionNY(player, 0)
IF ny > 0.7 THEN
    grounded = TRUE   ; Standing on ground
ELSE
    grounded = FALSE   ; Hitting a wall
ENDIF
```

---

## Picking

### `Entity.Pick(entityHandle, range)`

Casts a ray from the entity's position along its forward direction and returns the first entity hit.

- `range` (float) — Maximum pick distance.

**Returns:** `handle` — Hit entity, or 0 if nothing hit.

---

### `Entity.PickMode(entityHandle, mode)`

Sets the picking mode for an entity (whether it can be picked, and how).

- `mode` (int) — 0 = not pickable, 1 = sphere, 2 = polygon.

---

## Sprite Entities

### `Entity.LoadSprite(filePath)` / `LoadSprite(filePath)`

Loads a textured billboard sprite as an entity. The sprite always faces the camera.

**Returns:** `handle`

---

### `Entity.CreateSprite()`

Creates an empty sprite entity.

**Returns:** `handle`

---

### `ScaleSprite(entityHandle, sx, sy)`

Sets the sprite entity's scale.

---

### `SpriteViewMode(entityHandle, mode)` / `Entity.SpriteViewMode(entityHandle, mode)`

Sets how the sprite faces the camera: 1 = fixed, 2 = free, 3 = upright.

---

## Brushes (Blitz3D Compatibility)

Brushes are material objects for painting entity surfaces with colors and textures (Blitz3D pattern).

### `LoadBrush(filePath)`

Loads a brush from a texture file.

**Returns:** `handle`

### `FreeBrush(brushHandle)`

Frees a brush.

### `BrushColor(brushHandle, r, g, b)`

Sets the brush color.

### `BrushAlpha(brushHandle, alpha)`

Sets the brush alpha.

### `BrushBlend(brushHandle, mode)`

Sets the brush blend mode.

### `GetEntityBrush(entityHandle)`

Returns the brush applied to an entity.

**Returns:** `handle`

### `PaintSurface(surfaceHandle, brushHandle)`

Applies a brush to a model surface.

### `GetSurfaceBrush(surfaceHandle)`

Returns the brush on a surface.

**Returns:** `handle`

---

## Per-Frame Update

### `Entity.Update(dt)`

Updates all entities: applies velocity, gravity, collision detection and response. Call once per frame.

- `dt` (float) — Delta time.

**How it works:**
1. For each entity with a velocity, moves it by `velocity * dt`.
2. Applies gravity if set.
3. Runs the collision detection pass between all type pairs registered with `Collisions()`.
4. Resolves collisions by sliding entities apart.
5. Syncs positions with Jolt physics bodies (if `Entity.AddPhysics` was called).

```basic
WHILE NOT Window.ShouldClose()
    dt = Time.Delta()
    Entity.Update(dt)
    ; ... render ...
WEND
```

---

## Rendering

### `Entity.DrawAll()`

Draws every visible entity. Must be called inside a camera block.

**How it works:** Iterates all live entities, skips hidden ones, applies their transform (position × rotation × scale), sets their color/alpha, and draws their mesh. Entities with physics bodies have their transform synced from the physics world first.

```basic
Camera.Begin(cam)
    Entity.DrawAll()
Camera.End(cam)
```

---

### `Entity.Draw(entityHandle)`

Draws a single entity. Useful when you need custom draw order.

---

## Entity Freeing

### `Entity.FreeEntities()`

Frees **all** entities at once. Use for scene transitions.

---

## Physics Integration

### `Entity.AddPhysics(entityHandle, shapeType, mass)` / `Entity.Physics(entityHandle, ...)`

Adds a Jolt physics body to an entity. The entity's position and rotation will be synced to/from the physics simulation each frame.

- `entityHandle` (handle) — Entity.
- `shapeType` (string/int) — Shape: `"box"`, `"sphere"`, `"capsule"`.
- `mass` (float) — Mass in kg. 0 = static (immovable).

**How it works:** Creates a Jolt rigid body with the given shape, places it at the entity's current position, and links it. During `Entity.Update`, the entity's transform is read from the physics body (for dynamic) or written to it (for kinematic).

```basic
; Dynamic physics crate
crate = Entity.CreateBox(1, 1, 1, 180, 120, 60)
Entity.SetPos(crate, 0, 5, 0)
Entity.AddPhysics(crate, "box", 10)

; Static floor
floor = Entity.CreateBox(20, 0.5, 20, 100, 100, 100)
Entity.SetPos(floor, 0, -0.25, 0)
Entity.AddPhysics(floor, "box", 0)
```

---

### `Entity.SetBounce(entityHandle, restitution)`

Sets how bouncy the entity's physics body is. 0 = no bounce, 1 = perfect bounce.

**How it works:** Updates both the entity's internal bounce value and the Jolt body's restitution coefficient via `SetRestitutionToIndex`.

---

### `Entity.SetFriction(entityHandle, friction)`

Sets the friction of the entity's physics body.

**How it works:** Updates both the internal friction value and the Jolt body friction via `SetFrictionToIndex`.

---

### `Entity.SetGravity(entityHandle, gravity)` / `ent.SetGravity(g)`

Sets per-entity gravity override.

---

## Terrain Interaction

### `Entity.ClampToTerrain(entityHandle)`

Snaps an entity's Y position to the terrain height at its current XZ position.

---

### `Terrain.SnapY(entityHandle)` / `Terrain.Place(entityHandle)`

Places an entity on the terrain surface.

---

## Easy Mode Shortcuts

| Shortcut | Maps To |
|----------|---------|
| `CreateEntity()` | `Entity.Create()` |
| `Cube(w, h, d, r, g, b)` | `Entity.CreateBox(w, h, d, r, g, b)` |
| `Sphere(seg, r, g, b)` | `Entity.CreateSphere(seg, r, g, b)` |
| `LoadMesh(path)` | `Entity.Load(path)` |
| `PositionEntity(e, x, y, z)` | `Entity.SetPos(e, x, y, z)` |
| `RotateEntity(e, p, y, r)` | `Entity.RotateEntity(e, p, y, r)` |
| `ScaleEntity(e, sx, sy, sz)` | `Entity.Scale(e, sx, sy, sz)` |
| `EntityColor(e, r, g, b)` | `Entity.Color(e, r, g, b)` |
| `EntityAlpha(e, a)` | `Entity.Alpha(e, a)` |
| `EntityShininess(e, s)` | `Entity.Shininess(e, s)` |
| `EntityBlend(e, m)` | `Entity.Blend(e, m)` |
| `EntityTexture(e, tex)` | `Entity.Texture(e, tex)` |
| `MoveEntity(e, f, r, u)` | `Entity.Move(e, f, r, u)` |
| `TurnEntity(e, p, y, r)` | `Entity.Turn(e, p, y, r)` |
| `PointEntity(e, target)` | Look-at target entity |
| `HideEntity(e)` / `ShowEntity(e)` | `Entity.Visible(e, FALSE/TRUE)` |
| `FreeEntity(e)` | `Entity.Free(e)` |
| `FreeEntities()` | `Entity.FreeEntities()` |
| `EntityX(e)` / `EntityY(e)` / `EntityZ(e)` | Position components |
| `EntityPitch(e)` / `EntityYaw(e)` / `EntityRoll(e)` | Rotation components |
| `EntityType(e, t)` | `Entity.Type(e, t)` |
| `EntityRadius(e, r)` | `Entity.Radius(e, r)` |
| `Collisions(a, b, m)` | Enable collision pair |
| `EntityCollided(e, t)` | Check collision |
| `CountCollisions(e)` | Count collisions |
| `GetCollisionEntity(e, i)` | Get collision entity |
| `ResetEntity(e)` | Reset collision state |
| `LoadSprite(path)` | `Entity.LoadSprite(path)` |
| `ScaleSprite(e, sx, sy)` | Scale sprite entity |
| `SpriteViewMode(e, m)` | Sprite billboard mode |
| `LoadBrush(path)` / `FreeBrush(b)` | Brush load/free |
| `BrushColor(b, r, g, b)` / `BrushAlpha(b, a)` | Brush properties |
| `PaintSurface(s, b)` | Apply brush to surface |
| `TFormVector(...)` | Transform vector between spaces |

---

## Full Example

A 3D scene with physics entities, collisions, and rendering.

```basic
Window.Open(1280, 720, "Entity Demo")
Window.SetFPS(60)

cam = Camera.Create()
cam.pos(0, 10, 20)
cam.look(0, 2, 0)
cam.fov(60)

; Set up physics
World.Gravity(0, -9.81, 0)

; Create floor
floor = Entity.CreateBox(20, 0.5, 20, 80, 120, 80)
Entity.SetPos(floor, 0, -0.25, 0)
Entity.AddPhysics(floor, "box", 0)

; Create falling crates
FOR i = 0 TO 9
    crate = Entity.CreateBox(1, 1, 1, 180 + i * 7, 120, 60)
    Entity.SetPos(crate, (i - 5) * 1.5, 5 + i * 2, 0)
    Entity.AddPhysics(crate, "box", 5)
    Entity.SetBounce(crate, 0.3)
NEXT

; Create player sphere
player = Entity.Create()
Entity.SetPos(player, 0, 1, 5)
Entity.Radius(player, 0.5)
Entity.Type(player, 1)

; Enable collision between crates and player
Entity.Type(crate, 2)
Collisions(1, 2, 1)

WHILE NOT Window.ShouldClose()
    dt = Time.Delta()

    ; Move player
    IF Input.KeyDown(KEY_A) THEN Entity.Move(player, 0, -5 * dt, 0)
    IF Input.KeyDown(KEY_D) THEN Entity.Move(player, 0, 5 * dt, 0)
    IF Input.KeyDown(KEY_W) THEN Entity.Move(player, 5 * dt, 0, 0)
    IF Input.KeyDown(KEY_S) THEN Entity.Move(player, -5 * dt, 0, 0)

    ; Update all entities (physics + collisions)
    Entity.Update(dt)

    ; Check collisions
    IF EntityCollided(player, 2) THEN
        PRINT "Touched a crate!"
    ENDIF

    ; Render
    Render.Clear(30, 30, 50)
    Camera.Begin(cam)
        Draw.Grid(20, 1.0)
        Entity.DrawAll()
    Camera.End(cam)

    Draw.Text("WASD = Move", 10, 10, 18, 255, 255, 255, 255)
    Render.Frame()
WEND

Entity.FreeEntities()
Camera.Free(cam)
Window.Close()
```

---

## See Also

- [WORLD](WORLD.md) — Gravity, world update, fog
- [PHYSICS](PHYSICS.md) — Low-level Jolt physics commands
- [CAMERA](CAMERA.md) — Camera for rendering entities
- [MODEL](MODEL.md) — Loading 3D models for entities
- [PLAYER](PLAYER.md) — Player character controller
