// Package engineassets documents where the engine ships resources via //go:embed so a
// distribution binary does not depend on a repo-relative /shaders or /scripts tree.
//
// Bundled today:
//   - runtime/shaders: standard GLSL (embed in package shaders)
//   - runtime/mbgui: official raygui .rgs themes (embedded_raygui_styles.go)
//   - compiler/builtinmanifest: commands.json (manifest_json.go)
//
// Game content (.mb, glTF, audio paths) remains on disk unless you package a payload bundle.
package engineassets

// AuditDone marks that the engine-owned embed locations above have been inventoried for single-file shipping.
const AuditDone = true
