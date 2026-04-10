# Compiler specification (selected topics)

Focused notes for contributors and tooling. Authoritative language rules remain in [`compiler/errors/MoonBasic.md`](../compiler/errors/MoonBasic.md); pipeline API in [`ARCHITECTURE.md`](../ARCHITECTURE.md).

## Safety guarding: `ENTITY` spatial macros

The fast path for **`ENTITY.X` / `Y` / `Z` / `P` / `W` / `YAW` / `R`** lowers to **`OpEntityPropGet`** / **`OpEntityPropSet`** (see **`compiler/codegen`**, **`vm/vm_dispatch.go`**).

### Compile time

When the entity index is a **numeric literal** (integer or whole float), the **semantic analyzer** (so **`go run . --check`** catches it) and the **code generator** reject:

- **Negative** indices.
- Indices **≥ `runtime.MaxEntitySpatialIndex`** (exclusive cap, currently **2²⁴**).

Implementation: **`compiler/entityspatial/validate.go`** (shared rules); **`compiler/semantic/analyze.go`**; **`compiler/codegen/entity_macro_validate.go`**.

### Run time

The VM always validates the resolved **numeric** id (including after **`EntityRef`** handle decode) before touching shared SoA memory:

- **`id < 0`** or **`id ≥ MaxEntitySpatialIndex`** → clear **`ENTITY:`** runtime error (no silent slice access).
- If **`Registry.Spatial`** is used and the index is **in range of the SoA slice** but **`Registry.EntityIDActive(id)`** is false (set by **mbentity**), the VM errors instead of reading/writing a stale slot.

Non-literal indices (e.g. **`ENTITY.X(e)`**) are not constant-folded at compile time; they rely on the VM guards above.
