//go:build !windows && !(darwin || freebsd || linux)

package raylibpurego

import (
	"fmt"
	"runtime"
)

func loadDynamicLibrary(path string) (uintptr, error) {
	_ = path
	return 0, fmt.Errorf("raylibpurego: dynamic load not supported on GOOS=%s", runtime.GOOS)
}
