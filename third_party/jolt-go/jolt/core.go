package jolt

// #include "wrapper/core.h"
import "C"
import (
	"fmt"
	"math"
)

// Init initializes Jolt Physics (call once at startup)
func Init() error {
	result := C.JoltInit()
	if result == 0 {
		return fmt.Errorf("failed to initialize Jolt")
	}
	return nil
}

// Shutdown cleans up Jolt resources (call once at exit)
func Shutdown() {
	C.JoltShutdown()
}

// DegreesToRadians converts degrees to radians
func DegreesToRadians(degrees float32) float32 {
	return degrees * math.Pi / 180.0
}
