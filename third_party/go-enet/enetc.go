package enet

// #cgo !windows pkg-config: libenet
// #cgo windows CFLAGS: -Ienet/include/
// #cgo windows LDFLAGS: -Lenet/ -lenet -lws2_32 -lwinmm
// #include <enet/enet.h>
import "C"
import "fmt"

// Initialize enet
func Initialize() {
	C.enet_initialize()
}

// Deinitialize enet
func Deinitialize() {
	C.enet_deinitialize()
}

// LinkedVersion returns the linked version of enet currently being used.
// Returns MAJOR.MINOR.PATCH as a string.
func LinkedVersion() string {
	var version = uint32(C.enet_linked_version())
	major := uint8(version >> 16)
	minor := uint8(version >> 8)
	patch := uint8(version)
	return fmt.Sprintf("%d.%d.%d", major, minor, patch)
}
