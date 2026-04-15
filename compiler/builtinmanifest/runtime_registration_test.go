package builtinmanifest

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// findRepoRoot walks upward from the test working directory until go.mod is found.
func findRepoRoot(t *testing.T) string {
	t.Helper()
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatal("go.mod not found from test working directory")
		}
		dir = parent
	}
}

// loadRuntimeSourceHaystack concatenates all runtime/**/*.go files for substring audits.
func loadRuntimeSourceHaystack(t *testing.T) string {
	t.Helper()
	root := findRepoRoot(t)
	rtDir := filepath.Join(root, "runtime")

	var b strings.Builder
	err := filepath.WalkDir(rtDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		if filepath.Ext(path) != ".go" {
			return nil
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		b.Write(data)
		b.WriteByte('\n')
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return b.String()
}

// assertManifestKeyLiteralInRuntime checks that key appears as a Register/stub string, or (for PARTICLE3D.*)
// that the mirrored PARTICLE.* literal exists (regDual).
func assertManifestKeyLiteralInRuntime(t *testing.T, haystack, key string) {
	t.Helper()
	if strings.HasPrefix(key, "PARTICLE3D.") {
		sibling := "PARTICLE." + strings.TrimPrefix(key, "PARTICLE3D.")
		sneedle := `"` + sibling + `"`
		if !strings.Contains(haystack, sneedle) {
			t.Errorf("manifest %s: expected dual sibling literal %s under runtime/ (regDual)", key, sneedle)
		}
		return
	}
	needle := `"` + key + `"`
	if !strings.Contains(haystack, needle) {
		t.Errorf("manifest %s: expected literal %s somewhere under runtime/", key, needle)
	}
}

// TestManifestCreateKeysAppearInRuntimeSources ensures every manifest command whose key
// ends with ".CREATE" appears as a literal "NS.METHOD" string in runtime sources
// (Register(...), regDual(...), stub name lists, etc.). Catches manifest-only CREATE rows
// with no runtime alias, including split build-tag files.
func TestManifestCreateKeysAppearInRuntimeSources(t *testing.T) {
	tbl := Default()
	haystack := loadRuntimeSourceHaystack(t)

	for key := range tbl.Commands {
		if !strings.HasSuffix(key, ".CREATE") {
			continue
		}
		assertManifestKeyLiteralInRuntime(t, haystack, key)
	}
}

// TestManifestMakeKeysAppearInRuntimeSources is the same audit for deprecated *.MAKE siblings
// (must still be registered at runtime until removal).
func TestManifestMakeKeysAppearInRuntimeSources(t *testing.T) {
	tbl := Default()
	haystack := loadRuntimeSourceHaystack(t)

	for key := range tbl.Commands {
		if !strings.HasSuffix(key, ".MAKE") {
			continue
		}
		assertManifestKeyLiteralInRuntime(t, haystack, key)
	}
}
