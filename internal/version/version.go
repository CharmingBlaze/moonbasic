// Package version holds the moonBASIC build identity for CLI binaries.
//
// Release builds should set Version at link time, for example:
//
//	go build -ldflags="-X moonbasic/internal/version.Version=v1.2.18" ./cmd/moonbasic
package version

// Version is overridden with -ldflags=-X moonbasic/internal/version.Version=... for tagged releases.
var Version = "devel"
