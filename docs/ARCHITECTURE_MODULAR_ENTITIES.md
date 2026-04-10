# Modular Entity Architecture Guide

## Overview

The `ent` struct in MoonBASIC has been skeletonized to improve memory efficiency and maintainability. Previously, every entity carried the memory overhead of all engine features (animation, collision clusters, physics, tweens), even if it was a simple static prop. 

We now use a **Modular Extension** pattern where specialized data is stored in a separate `entExt` struct, lazily allocated only when the entity requires those features.

## The `ent` Skeleton

The core `ent` struct (in `runtime/mbentity/entity_ent_cgo.go`) now only contains:
- **Identity**: `id`, `kind`.
- **Primary Transform**: `pos`, `pitch`, `yaw`, `roll`, `scale`.
- **Core Raylib Refs**: `rlModel`, `hasRLModel`, `rlTex`, `rlMat`.
- **Visibility**: `hidden`, `alpha`.

## The `entExt` Extension

All other fields live in `entExt`. This includes:
- **Collision bookkeeping**: `hits`, `hitPos`, `hitN`, `collType`, `slide`.
- **Advanced Rendering**: `isSprite`, `spriteMode`, `outlineThickness`, `outlineColor`.
- **State Management**: `ghostMode`, `ghostTimer`, `shadowCast`.
- **Animation Data**: `modelAnims`, `animLen`, `animSpeed`.
- **QoL Tweens**: `tweenFading`, `tweenTurning`, `pulseSpeed`.

## Usage Pattern: `getExt()`

To access extension fields, use the `getExt()` method. This ensures that the extension struct is initialized if it doesn't exist.

### ✅ Correct Pattern (Lazy Initialization)
```go
// Instead of:
// e.collided = true

// Use:
ext := e.getExt()
ext.collided = true
```

### ⚠️ Read-Only / Non-Allocating Access
If you only need to check if a feature is active without forcing an allocation:
```go
if e.ext != nil && e.ext.collided {
    // ...
}
```

## Migration Checklist

When updating legacy code to the modular system:
1. Identify fields that have moved to `entExt`.
2. Wrap accesses in `e.getExt()` for setters.
3. For heavy loops (like drawing or collision loops), call `getExt()` once at the top of the function if multiple fields are accessed.
4. Ensure `nil` safety if you are not using `getExt()` (direct `e.ext` usage).

## Adding New Features

Follow these steps to add a new entity feature:
1. Add the field to the `entExt` struct in `entity_ent_cgo.go`.
2. DO NOT add it to the core `ent` unless every single entity in the engine (thousands) will use it every frame.
3. Update the `ENTITY.*` command handlers to use the extension via `getExt()`.
