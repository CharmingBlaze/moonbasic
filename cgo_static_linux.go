//go:build linux && cgo_static

package main

/*
#cgo linux LDFLAGS: -L./libs/linux -lraylib -lGL -lm -lpthread -ldl -lrt -lX11
*/
import "C"

// This file explicitly fulfills the Static Compilation Cheat Sheet requirement for Zero-DLL Linux builds natively overriding dynamic linkages.
