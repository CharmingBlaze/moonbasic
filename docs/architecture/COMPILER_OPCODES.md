# Opcode & Boundary Technical Reference

MoonBASIC VM relies on an Opcode instruction registry driving logic across the interpreter arrays mapping accurately to the GPU interactions. Every instruction bridges through `vm/opcode/opcode.go`.

## Memory Interactions

Opcodes like `OpSyncPhysics` invoke mathematical translations reading specific registers directly from shared slices (`Float32View` BaseOffset addressing arrays) bound inherently through Wazero.

For all custom bindings involving the GPU:
- `OpUserFn` dispatches calls interacting with the Handle Engine.
- Handlers specifically implement Zero-Copy arrays accessing limits, extracting properties matching `Registry[int32]Asset` models avoiding garbage collections natively prior to triggering `Raylib.*` abstractions dynamically.

Refer to active opcode array boundaries when optimizing logic inside the script sequences evaluating loop contexts cleanly avoiding excessive CGO interactions.
