//go:build !windows && (darwin || freebsd || linux)

package raylibpurego

import "github.com/ebitengine/purego"

func loadDynamicLibrary(path string) (uintptr, error) {
	return purego.Dlopen(path, purego.RTLD_NOW|purego.RTLD_GLOBAL)
}
