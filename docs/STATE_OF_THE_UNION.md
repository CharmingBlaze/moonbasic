# MoonBASIC: State of the Union

This document summarizes the high-level goals and current status of the MoonBASIC engine for all contributors.

## Current Focus: Mario 64 Parity

We are currently in **Phase 2: High-Fidelity Refinement**. Our primary objective is to make the "Easy Mode" demo (`examples/mario64/main_easymode.mb`) represent the absolute best practices of the engine.

### Completed Milestones
- [x] **Jolt Integration**: Full high-performance rigid body and KCC support on Linux and Windows (CGO).
- [x] **Visual Snap Band**: Eliminated solver jitter for a "Nintendo-smooth" feel.
- [x] **Easy Mode API**: `Character.Create`, `Character.Update`, and `hero.Jump` abstractions.
- [x] **Decoupled Raylib & HAL**: Isolated engine core from GPU drivers; enabled headless compiler unit tests.
- [x] **Static Linking**: Supported standalone "Zero-DLL" builds using CGO + Zig CC.
- [x] **Cross-Platform Parity**: Achieved native physics stability on Windows via prebuilt static libraries.

### Active Priorities
- **Method Standardization**: Migrating all showcase scripts to use handle methods (`hero.SetGravity`) rather than namespaced calls (`CHARACTERREF.SETGRAVITY`). [IN PROGRESS]
- **Driver Parity**: Ensuring the `Null` driver implements all new rendering features for robust headless testing.
- **Advanced Jolt Integration**: Wiring `PHYSICS3D.PROCESSCOLLISIONS` to high-level entity listeners.

## The "Same Path" Philosophy

All MoonBASIC development follows the **Same Path** rule:

1. **One Script runs Everywhere**: A `.mb` file written using "Easy Mode" syntax on Windows must behave identically on Linux without code changes.
2. **Same physics on desktop**: **Windows and Linux** fullruntime builds with **CGO + Jolt** share one native path. Headless, mobile, or **`!cgo`** builds use stubs so the same `.mb` compiles; KCC gameplay requires the toolchain + Jolt libs documented in [BUILDING.md](BUILDING.md) and [JOLT_WINDOWS_PARITY.md](JOLT_WINDOWS_PARITY.md).
3. **Handle Consistency**: Handles (Entities, Characters, Models) are the primary currency of the API. Prefer handle methods over global namespaced commands in all user-facing examples.

## Near-Term Roadmap
- [ ] **Batch 7 Atmosphere**: Volumetric lighting and skybox blending.
- [ ] **Advanced KCC Features**: Crouching, swimming, and wall-jumps integrated into `Character.Create`.
- [ ] **IDE Tooling**: Improved VS Code IntelliSense for handle methods via `gopls` stubs.
