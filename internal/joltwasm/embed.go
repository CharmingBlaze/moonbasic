package joltwasm

import (
	_ "embed"
)

//go:embed jolt.wasm
var EmbeddedWASM []byte

// GetJoltBinary returns the byte footprint of the Jolt physics engine compiled
// to WebAssembly. This embeds Jolt natively directly inside the Go runtime
// eliminating the need for trailing binaries or dynamic library linkages on load.
func GetJoltBinary() []byte {
	return EmbeddedWASM
}
