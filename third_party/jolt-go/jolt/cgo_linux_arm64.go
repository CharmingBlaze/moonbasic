//go:build linux && arm64

package jolt

/*
#cgo LDFLAGS: -L${SRCDIR}/lib/linux_arm64 -ljolt_wrapper -lJolt -lstdc++ -lm -lpthread
*/
import "C"
