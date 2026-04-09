//go:build windows

package raylibpurego

import (
	"fmt"
	"syscall"
)

func loadDynamicLibrary(path string) (uintptr, error) {
	handle, err := syscall.LoadLibrary(path)
	if err != nil {
		return 0, fmt.Errorf("LoadLibrary %q: %w", path, err)
	}
	return uintptr(handle), nil
}
