//go:build cgo

package mbgui

import "embed"

// Official raygui binary styles from https://github.com/raysan5/raygui/tree/master/styles
// (raylib license / zlib — see raygui_styles/README.md).
//
//go:embed raygui_styles/*.rgs
var embeddedRayguiStyles embed.FS
