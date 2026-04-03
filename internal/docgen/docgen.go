package docgen

import (
	"os"
	"path/filepath"
)

// Generate reads commands JSON from jsonPath and writes a static site to outDir.
func Generate(jsonPath, outDir string) error {
	cmds, err := LoadCommandsJSON(jsonPath)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(outDir, 0755); err != nil {
		return err
	}
	if err := WriteStyle(outDir); err != nil {
		return err
	}
	return WriteSite(outDir, cmds)
}

// DefaultCommandsPath returns repo-relative path to embedded manifest source.
func DefaultCommandsPath() string {
	return filepath.Join("compiler", "builtinmanifest", "commands.json")
}
