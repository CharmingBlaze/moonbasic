//go:build cgo || (windows && !cgo)

package shaders

import (
	"embed"
)

//go:embed shd/*
var EmbeddedShaders embed.FS
