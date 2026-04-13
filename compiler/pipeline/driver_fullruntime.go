//go:build fullruntime

package pipeline

import (
	"moonbasic/hal"
	"moonbasic/drivers/video/raylib"
)

// DefaultDriver returns the Raylib driver for interactive builds.
func DefaultDriver() hal.Driver {
	d := raylib.NewDriver()
	return hal.Driver{
		Video:  d,
		Input:  d,
		System: d,
	}
}
