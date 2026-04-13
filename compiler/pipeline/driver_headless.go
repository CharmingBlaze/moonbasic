//go:build !fullruntime

package pipeline

import (
	"moonbasic/hal"
	"moonbasic/drivers/video/null"
)

// DefaultDriver returns the Null driver for compiler-only / headless builds.
func DefaultDriver() hal.Driver {
	d := null.NewDriver()
	return hal.Driver{
		Video:  d,
		Input:  d,
		System: d,
	}
}
