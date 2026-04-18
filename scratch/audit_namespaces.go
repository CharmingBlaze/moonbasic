//go:build ignore

package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type Command struct {
	Key         string `json:"key"`
	Description string `json:"description"`
	DeprecatedOf string `json:"deprecated_of"`
}

type Manifest struct {
	Commands []Command `json:"commands"`
}

func main() {
	data, err := os.ReadFile("compiler/builtinmanifest/commands.json")
	if err != nil {
		panic(err)
	}

	var manifest Manifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		panic(err)
	}
	commands := manifest.Commands

	// Collect namespaces and per-namespace command counts (skip deprecated aliases)
	nsCounts := map[string]int{}
	nsCommands := map[string][]string{}
	globals := []string{}
	skipped := 0

	for _, cmd := range commands {
		key := cmd.Key
		if key == "" {
			continue
		}
		// Skip legacy suffixed aliases
		if strings.ContainsAny(key, "$#%?") {
			skipped++
			continue
		}
		// Skip deprecated aliases (they have deprecated_of set and the key differs)
		if cmd.DeprecatedOf != "" && cmd.DeprecatedOf != key {
			skipped++
			continue
		}

		if strings.Contains(key, ".") {
			ns := strings.SplitN(key, ".", 2)[0]
			nsCounts[ns]++
			nsCommands[ns] = append(nsCommands[ns], key)
		} else {
			globals = append(globals, key)
		}
	}

	// Sort namespaces
	nsList := make([]string, 0, len(nsCounts))
	for ns := range nsCounts {
		nsList = append(nsList, ns)
	}
	sort.Strings(nsList)

	// Check which namespaces have a reference doc
	refDir := "docs/reference"
	entries, _ := os.ReadDir(refDir)
	docFiles := map[string]bool{}
	for _, e := range entries {
		name := strings.TrimSuffix(strings.ToUpper(e.Name()), ".MD")
		docFiles[name] = true
	}

	fmt.Printf("=== MoonBASIC Command Namespace Audit ===\n\n")
	fmt.Printf("Total commands (non-deprecated, non-suffixed): %d\n", len(commands)-skipped)
	fmt.Printf("Total global builtins: %d\n", len(globals))
	fmt.Printf("Total namespaces: %d\n\n", len(nsList))

	fmt.Printf("%-22s %6s  %s\n", "NAMESPACE", "CMDS", "DOC FILE")
	fmt.Printf("%s\n", strings.Repeat("-", 60))
	
	missing := []string{}
	for _, ns := range nsList {
		hasDoc := docFiles[ns]
		status := "✓"
		if !hasDoc {
			status = "✗ MISSING"
			missing = append(missing, ns)
		}
		fmt.Printf("%-22s %6d  %s\n", ns, nsCounts[ns], status)
	}

	fmt.Printf("\n=== UNDOCUMENTED NAMESPACES (%d) ===\n", len(missing))
	for _, ns := range missing {
		fmt.Printf("  %s (%d commands)\n", ns, nsCounts[ns])
		sort.Strings(nsCommands[ns])
		for _, cmd := range nsCommands[ns] {
			fmt.Printf("    - %s\n", cmd)
		}
	}

	// Check existing docs for namespaces that don't exist in manifest
	fmt.Printf("\n=== DOC FILES WITHOUT MATCHING NAMESPACE ===\n")
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".md") {
			continue
		}
		ns := strings.TrimSuffix(strings.ToUpper(e.Name()), ".MD")
		if nsCounts[ns] == 0 {
			// it's a thematic/guide doc, not namespace-specific – check a few known ones
			skipList := map[string]bool{
				"API_CONVENTIONS": true, "BEGINNER_FULL_STACK": true, "BLITZ2025": true,
				"BLITZ3D": true, "BLITZ_COMMAND_INDEX": true, "BLITZ_ESSENTIAL_API": true,
				"CAMERA_LIGHT_RENDER": true, "CHARACTER_PHYSICS": true, "CLOTH_ROPE_LIGHTING": true,
				"CSV_DATABASE": true, "DRAW_WRAPPERS": true, "ENTITYREF": true,
				"GAME_ENGINE_PATTERNS": true, "GAME_MATH_HELPERS": true, "GAMEPLAY_HELPERS": true,
				"GAMEHELPERS": true, "LESS_MATH": true, "LEVEL_COLLISION": true,
				"MODERN_BLITZ_COMMANDS": true, "NAVMESH": true, "PHYSICS2D": true,
				"PHYSICS_ADVANCED": true, "PHYSICS3D": true, "PROCEDURAL": true,
				"QOL": true, "RAYCAST": true, "RAYLIB_EXTRAS": true,
				"SCATTER_PROP_SPAWNER": true, "SCENE_ENGINE_BRIEF": true,
				"STRING_HEAP": true, "TEXTURE_DRAW_WRAPPERS": true,
				"UNIVERSAL_HANDLE_METHODS": true, "VEC_QUAT": true, "WAVE": true,
			}
			if !skipList[ns] && nsCounts[ns] == 0 {
				fmt.Printf("  %s\n", filepath.Join(refDir, e.Name()))
			}
		}
	}

	fmt.Printf("\n=== GLOBAL BUILTINS (no namespace) ===\n")
	sort.Strings(globals)
	for _, g := range globals {
		fmt.Printf("  %s\n", g)
	}
}
