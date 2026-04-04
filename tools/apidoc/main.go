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
	b.WriteString("## Naming conventions\n\n")
	b.WriteString("- **Registry / source form**: `NS.ACTION` in uppercase with a dot (e.g. `CAMERA.SETPOS`).\n")
	b.WriteString("- **Handle methods** (on a handle value): `cam.SetPos` dispatches to `CAMERA.SETPOS`. **`SetPosition`** is an alias for **`SetPos`** on spatial types.\n")
	b.WriteString("- **Spatial handles** (`Camera3D`, `Body3D`, `CharController`, `Model`, `Sprite`): use **`SetPos`** / `SETPOS` (sprite: x,y; others: x,y,z).\n")
	b.WriteString("- **`MODEL.SETPOS`**: sets the model root transform to a **translation matrix** (replaces prior rotation/scale on that matrix).\n")
	b.WriteString("- **`LIGHT.MAKE`**: zero arguments -> directional, white, intensity 1.0 (config handle for future lighting).\n")
	b.WriteString("- **`BODY3D.MAKE`**: zero arguments -> **DYNAMIC** motion builder.\n")
	b.WriteString("- **Errors**: include **file and line** when available; unknown commands use **did-you-mean** against the live registry.\n\n")
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
