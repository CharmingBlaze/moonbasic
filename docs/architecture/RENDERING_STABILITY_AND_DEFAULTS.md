# Rendering Stability and Default States

This document records the "Black Screen" post-mortem and establishes guidelines for maintaining rendering stability in moonBASIC.

## Post-Mortem: The Black Screen Issue (April 2026)

### Symptoms
Scripts running in graphical mode (especially 3D) would occasionally show a completely black viewport, despite clearing with non-black colors and drawing geometry.

### Root Causes

1.  **Near-Zero Ambient Light Defaults**: The default 3D PBR ambient light was set to `0.06`. In scenes without explicit lights, this resulted in nearly unlit (black) geometry that was difficult to distinguish from the background.
2.  **Global 2D Darkness Overlay**: The `mblight2d` module registered a per-frame draw hook that applied a full-viewport black rectangle to simulate darkness when 2D ambient light was low. Because this hook was registered globally, it affected 3D programs even if they didn't use 2D lighting features.
3.  **GPU Swap-Chain Latency**: On some Windows systems (especially with integrated Intel GPUs), the first few frames after `WINDOW.OPEN` are not correctly displayed orPresented. Without a "warmup" period of blank frames, the initial render state could appear corrupted or black.
4.  **Purego Disparity**: The Raylib "sidecar" (purego) path lacked the warmup frames present in the CGO path, leading to platform-specific startup failures.

### The Fix

-   **Ambient Boost**: Increased default 3D ambient light from `0.06` to `0.12`.
-   **Hook Hardening**: The `mblight2d` overlay now strictly checks for the presence of 2D lights (`len(lights) > 0`) before drawing any masking rectangle.
-   **Uniform Warmup**: Added a consistent 2-frame warmup sequence (plus input drain) to `WINDOW.OPEN` in both CGO and PureGo implementations.
-   **Clamped Alpha**: Added safety clamping and early returns to the overlay logic to prevent zero-alpha or negative-alpha edge cases.

---

## Guidelines for Future Development

To prevent similar "unintended darkness" states in the future, adhere to the following rules:

### 1. Visibility by Default
-   All default engine states must allow the user to see *something*. 
-   Avoid "cinematic" defaults (like heavy fog or pitch blackness) as the initial state. Let the developer opt-in to these.

### 2. Guard Global Per-Frame Hooks
-   Modules that register global frame hooks (like `mblight2d` or `mbtransition`) **MUST** include an inexpensive early-return check.
-   A hook should only perform drawing operations if its specific feature set is actually being used by the active script.

### 3. Cross-Driver Portability
-   Any "warmup" or "guard" logic added to the CGO path to fix hardware-specific glitches **MUST** also be implemented in the PureGo sidecar path (`purego_minimal.go`).
-   Test graphical changes on both dedicated GPUs and integrated graphics (e.g. Intel Iris Xe) where possible.

### 4. Visibility Checklist for New Features
- [ ] Does this feature draw over the entire viewport?
- [ ] If yes, is there a clear "off" state when not in use?
- [ ] Are the default properties bright enough to be seen in a dark room?
- [ ] Does it gracefully handle cases where the user has not created any light sources?
