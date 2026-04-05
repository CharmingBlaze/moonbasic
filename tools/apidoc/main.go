// Apidoc writes docs/API_CONSISTENCY.md from compiler/builtinmanifest/commands.json.
// Run from the module root: go run ./tools/apidoc
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type root struct {
	Commands []struct {
		Key       string   `json:"key"`
		Args      []string `json:"args"`
		Returns   string   `json:"returns,omitempty"`
		Namespace string   `json:"namespace,omitempty"`
	} `json:"commands"`
}

func main() {
	path := filepath.Join("compiler", "builtinmanifest", "commands.json")
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "read %s: %v (run from repository root)\n", path, err)
		os.Exit(1)
	}
	var r root
	if err := json.Unmarshal(data, &r); err != nil {
		fmt.Fprintf(os.Stderr, "parse: %v\n", err)
		os.Exit(1)
	}
	type sig struct {
		key, ret string
		args     []string
	}
	byNS := make(map[string][]sig)
	for _, c := range r.Commands {
		parts := strings.SplitN(c.Key, ".", 2)
		ns := parts[0]
		if c.Namespace != "" {
			ns = strings.ToUpper(c.Namespace)
		}
		byNS[ns] = append(byNS[ns], sig{key: c.Key, ret: c.Returns, args: c.Args})
	}
	nss := make([]string, 0, len(byNS))
	for ns := range byNS {
		nss = append(nss, ns)
	}
	sort.Strings(nss)

	var b strings.Builder
	b.WriteString("# moonBASIC API consistency\n\n")
	b.WriteString("This document is generated from `compiler/builtinmanifest/commands.json`.\n\n")
	b.WriteString("Refresh: `go run ./tools/apidoc` (from the repository root).\n\n")
	b.WriteString("## Related documentation\n\n")
	b.WriteString("- **[ERROR_MESSAGES.md](../ERROR_MESSAGES.md)** — compile-time vs runtime errors, did-you-mean, heap handle hints.\n")
	b.WriteString("- **[ROADMAP.md](../ROADMAP.md)** — phased engineering plan (polish → rendering → 2D → systems → …).\n\n")
	b.WriteString("## Naming conventions\n\n")
	b.WriteString("- **Registry / source form**: `NS.ACTION` in uppercase with a dot (e.g. `CAMERA.SETPOS`).\n")
	b.WriteString("- **Handle methods** (on a handle value): `cam.SetPos` dispatches to `CAMERA.SETPOS`. **`SetPosition`** is an alias for **`SetPos`** where both are registered (same handler).\n")
	b.WriteString("- **Spatial handles** (`Camera3D`, `Body3D`, `Model`, `Sprite`, `Light2D`): use **`SETPOS`** for position. Aliases **`SETPOSITION`** exist for **Camera**, **Model**, **Body3D**, **Sprite**, **Light2D** — same implementation as `SETPOS`.\n")
	b.WriteString("- **3D lights** (`LIGHT.*`): use **`LIGHT.SETDIR`** for the directional sun (normalized). **`LIGHT.SETPOS`** stores point/spot position; **`LIGHT.SETTARGET`** moves the shadow frustum look-at; **`RENDER.SETAMBIENT`** sets PBR ambient tint.\n")
	b.WriteString("- **`MODEL.SETPOS`**: sets the model root transform to a **translation matrix** (replaces prior rotation/scale on that matrix).\n")
	b.WriteString("- **Creation verbs**: `*.MAKE` for procedural handles; `*.LOAD` for assets (`SPRITE.LOAD`, `MODEL.LOAD`); materials use `MATERIAL.MAKEDEFAULT` / `MATERIAL.MAKEPBR`.\n\n")
	b.WriteString("## Default values (common `Make` paths)\n\n")
	b.WriteString("| Command | Defaults |\n")
	b.WriteString("|----------|----------|\n")
	b.WriteString("| `CAMERA.MAKE` | position (0, 2, 8), target (0, 0, 0), up (0, 1, 0), FOV 45°, perspective |\n")
	b.WriteString("| `LIGHT.MAKE` | kind `directional`, white, intensity 1.0, direction toward normalized (-1,-2,-1) |\n")
	b.WriteString("| `BODY3D.MAKE` | no args → **DYNAMIC** motion type |\n")
	b.WriteString("| `MATERIAL.MAKEDEFAULT` / `MAKEPBR` | see `runtime/mbmodel3d` (material modules) |\n\n")
	b.WriteString("## Debug watch overlay\n\n")
	b.WriteString("`DEBUG.WATCH(label$, value)` stores rows; `DEBUG.WATCHCLEAR` clears them. With **CGO** and Raylib, the window pipeline may draw a **top-left overlay** each frame (`runtime/mbdebug/overlay_cgo.go`) when **`DEBUG.ENABLE`** was called or the host enabled **`Registry.DebugMode`** (e.g. **`--debug`**). **`DEBUG.DISABLE`** clears the user override. Without CGO, watches are stored but not drawn.\n\n")
	b.WriteString("## Errors\n\n")
	b.WriteString("- **Compile-time**: unknown `NS.METHOD` → did-you-mean within namespace + manifest listing (see `compiler/semantic/cmdhint.go`).\n")
	b.WriteString("- **Runtime**: VM wraps native errors with **source file and line** when available (`vm/vm.go`). Unknown registry keys → `runtime.FormatUnknownRegistryCommand`.\n\n")
	b.WriteString("## Commands by namespace\n\n")
	for _, ns := range nss {
		entries := byNS[ns]
		sort.Slice(entries, func(i, j int) bool {
			if entries[i].key != entries[j].key {
				return entries[i].key < entries[j].key
			}
			return len(entries[i].args) < len(entries[j].args)
		})
		b.WriteString("### ")
		b.WriteString(ns)
		b.WriteString("\n\n")
		for _, e := range entries {
			b.WriteString("- **`")
			b.WriteString(e.key)
			b.WriteString("`**")
			if len(e.args) > 0 {
				b.WriteString(" - args: ")
				b.WriteString(strings.Join(e.args, ", "))
			} else {
				b.WriteString(" - args: (none)")
			}
			if e.ret != "" {
				b.WriteString(" -> returns ")
				b.WriteString(e.ret)
			}
			b.WriteString("\n")
		}
		b.WriteString("\n")
	}
	out := filepath.Join("docs", "API_CONSISTENCY.md")
	if err := os.WriteFile(out, []byte(b.String()), 0o644); err != nil {
		fmt.Fprintf(os.Stderr, "write %s: %v\n", out, err)
		os.Exit(1)
	}
	fmt.Fprintf(os.Stderr, "wrote %s\n", out)
}
