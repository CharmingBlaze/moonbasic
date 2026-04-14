//go:build windows && amd64 && cgo

/*
WINDOWS CGO LDFLAGS
The linker flags below are specifically ordered to ensure proper resolution of Jolt 
and standard library symbols on Windows. 

LIBRARIES:
- libJolt_wrapper.a: C-wrapper bridge exposing CharacterVirtual and World APIs.
- libJolt.a: The core Jolt Physics engine (compiled with SIMD).

BUILD STEPS:
Static libraries must be pre-compiled using the PowerShell script in 
third_party/jolt-go/scripts/build-libs-windows.ps1 before building the Go runtime.
*/

package jolt

/*
#cgo LDFLAGS: -L${SRCDIR}/lib/windows_amd64 -ljolt_wrapper -lJolt -lstdc++
*/
import "C"
