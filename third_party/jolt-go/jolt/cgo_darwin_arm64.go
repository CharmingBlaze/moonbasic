//go:build darwin && arm64

package jolt

/*
#cgo LDFLAGS: -L${SRCDIR}/lib/darwin_arm64 -ljolt_wrapper -lJolt -lc++
*/
import "C"
