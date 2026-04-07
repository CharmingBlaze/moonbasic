//go:build linux && amd64

package jolt

/*
#cgo LDFLAGS: -L${SRCDIR}/lib/linux_amd64 -ljolt_wrapper -lJolt -lstdc++ -lm -lpthread
*/
import "C"
