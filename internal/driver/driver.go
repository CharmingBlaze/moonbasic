// Package driver describes how the engine selects a Raylib backend: native CGO (statically linked)
// versus a dynamically loaded sidecar (purego + raylib.dll/.so/.dylib next to the executable).
//
// Runtime selection is exposed as [GetDefaultDriver]. Windowing implementations still live under
// runtime/window and internal/raylibpurego; this package is the single place for probe logic.
package driver

import (
	"fmt"
	"os"
	"strings"

	"moonbasic/internal/raylibpurego"
)

// Kind is the resolved backend for the current process.
type Kind int

const (
	KindUnavailable Kind = iota
	KindNativeCGO
	KindPuregoDLL
)

func (k Kind) String() string {
	switch k {
	case KindNativeCGO:
		return "native_cgo"
	case KindPuregoDLL:
		return "purego_dll"
	default:
		return "unavailable"
	}
}

// Selection is the result of [GetDefaultDriver] (or [ProbeRaylibSharedObject] alone).
type Selection struct {
	Kind   Kind
	Detail string
}

func (s Selection) String() string {
	if s.Detail == "" {
		return s.Kind.String()
	}
	return fmt.Sprintf("%s: %s", s.Kind.String(), s.Detail)
}

// EnvDriver is the name of the environment variable that overrides automatic probing.
// Values: "auto" (default), "cgo", "purego" (or "dll").
const EnvDriver = "MOONBASIC_DRIVER"

// GetDefaultDriver picks a backend without requiring the user to pass build tags:
//   - If the binary was built with cgo and Raylib is linked, native CGO is preferred.
//   - Otherwise it tries to load the conventional sidecar shared library via purego.
//
// Override with MOONBASIC_DRIVER=cgo|purego|auto.
func GetDefaultDriver() Selection {
	switch strings.ToLower(strings.TrimSpace(os.Getenv(EnvDriver))) {
	case "cgo", "native":
		if nativeCGOLinked {
			return Selection{Kind: KindNativeCGO, Detail: "MOONBASIC_DRIVER=cgo; raylib linked at build (cgo)"}
		}
		return Selection{Kind: KindUnavailable, Detail: "MOONBASIC_DRIVER=cgo but this binary was built with CGO disabled (no linked raylib)"}
	case "purego", "dll", "sidecar":
		return probePurego()
	case "auto", "":
		if nativeCGOLinked {
			return Selection{Kind: KindNativeCGO, Detail: "default: native CGO (raylib linked at build)"}
		}
		return probePurego()
	default:
		return Selection{Kind: KindUnavailable, Detail: "unknown " + EnvDriver + " value (use auto, cgo, or purego)"}
	}
}

// ProbeRaylibSharedObject attempts to open the conventional Raylib DLL/SO next to the executable.
// On success the library stays loaded for the process lifetime (same as a real init path).
func ProbeRaylibSharedObject() Selection {
	return probePurego()
}

func probePurego() Selection {
	_, err := raylibpurego.Load("")
	if err != nil {
		return Selection{Kind: KindUnavailable, Detail: err.Error()}
	}
	return Selection{Kind: KindPuregoDLL, Detail: "loaded sidecar " + raylibpurego.LibPath()}
}
