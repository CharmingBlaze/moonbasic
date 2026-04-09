// Package raylibpurego loads Raylib as a shared library without CGO, using ebitengine/purego.
// Use [Load] / [LoadFrom], then [RegisterGame] for a minimal window/input/draw set.
// See cmd/puregohello for a CGO_ENABLED=0 smoke test and docs/architecture/ZERO_CGO_RAYLIB.md.
package raylibpurego
