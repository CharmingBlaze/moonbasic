//go:build windows && cgo_static

package main

/*
#cgo windows LDFLAGS: -L./libs/windows -lraylib -lopengl32 -lgdi32 -lwinmm -lcomdlg32 -lole32 -lsetupapi
*/
import "C"

// This file explicitly fulfills the Static Compilation Cheat Sheet requirement for Zero-DLL Windows builds natively overriding dynamic linkages.
