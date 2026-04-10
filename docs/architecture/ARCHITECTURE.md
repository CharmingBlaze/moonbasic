# MoonBASIC Architecture: The Three-Pillar System

MoonBASIC employs a "Zero-DLL" architecture designed to produce static, high-performance runtime targets capable of seamless scripting operations across desktop domains. To maintain execution speed exceeding conventional C interpreters, the architecture separates responsibilities into three foundational pillars avoiding runtime garbage collection pauses on mathematical operations mapping straight to WebAssembly abstractions.

## Pillar 1: The Host (Go Language VM)
The interpreter loops construct the AST execution blocks into compiled AST arrays run inherently inside the native Go runtime. The Host controls:
- **Asset Management**: Disk I/O logic, VFS packers, shader caches, and file decoders (primarily through Pure-Go targets such as `qmuntal/gltf`).
- **Window Presentation**: Leveraging specific `-tags static` Raylib integrations built exclusively inside static Go contexts rendering strictly on the main OS thread.
- **VM Logistics**: Register allocations, stack framing bounds logic, handling errors, array allocation bounds arrays (The Hard Guard mechanism).

## Pillar 2: The Guest (WebAssembly Jolt Physics)
Physics calculations, linear math interpolations, constraint solving, and kinematic sweeping require continuous CPU cache utilization bypassing Go's GC entirely to prevent frame staggering.
- Driven unconditionally by `wazero` (pure-Go WASM runtime).
- Encapsulates `jolt.wasm`, evaluating all rigid-body matrix computations natively on its own linear heap memory model.

## Pillar 3: The Bridge (Zero-Copy Shared Memory)
The secret to avoiding cross-boundary performance penalties between Pillar 1 (Host) and Pillar 2 (Guest) is the Shared Buffer View. 
Instead of making API calls (e.g., `Entity.GetPosition()`) crossing Wasm-To-Go translation layers, we linearly map the WASM Guest Memory block immediately into Go `[]float32` slice bounds referencing arrays via direct pointer mapping.

Entity transforms dynamically adjust sequentially at predetermined array strides evaluating logic internally from:
`BaseAddress + (EntityID * 16 bytes)` producing $O(1)$ immediate CPU time updates shared instantaneously!
