# DBPro — Lights

moonBASIC: **`LIGHT.MAKE`**, **`LIGHT.SET*`** — [LIGHT.md](../LIGHT.md).

| DBPro | moonBASIC | Notes |
|-------|-----------|--------|
| **MAKE LIGHT (light)** | ✓ **`LIGHT.MAKE`** | Handle-based. |
| **DELETE LIGHT** | ✓ **`LIGHT.FREE`** | |
| **POSITION LIGHT** | ✓ **`LIGHT.SETPOSITION`** / **`SETPOS`** | |
| **POINT LIGHT** | ✓ **`LIGHT.SETTARGET`**, **`SETDIR`** | Type-dependent (directional vs spot). |
| **SET LIGHT RANGE** | ✓ **`LIGHT.SETRANGE`**, **`SETINTENSITY`** | |
| **SET LIGHT COLOR** | ✓ **`LIGHT.SETCOLOR`** | |
| **SET LIGHT TYPE** | ≈ **`LIGHT.MAKE`** + parameters | Types differ from DBPro enum. |
| **SET LIGHT TO OBJECT** | ≈ follow entity in loop | No single **attach-to-entity** name in all builds — use **entity position** + **`LIGHT.SETPOSITION`**. |
