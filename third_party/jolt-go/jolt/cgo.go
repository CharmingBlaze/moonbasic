/*
Package jolt provides Go bindings for the Jolt Physics engine.

Architecture: Go code → CGO → C wrapper → Jolt C++ library
Uses opaque pointers to hide C++ types from Go.

Pre-built binaries are provided for darwin/arm64 and linux/amd64.
Platform-specific LDFLAGS are in cgo_<os>_<arch>.go files.
*/
package jolt

/*
#cgo CXXFLAGS: -std=c++17 -IJoltPhysics -DNDEBUG -DJPH_DISABLE_CUSTOM_ALLOCATOR -DJPH_PROFILE_ENABLED -DJPH_DEBUG_RENDERER -DJPH_OBJECT_STREAM
*/
import "C"
