package pkgmgr

import (
	"os"
	"path/filepath"
	"testing"
)

func TestInstallFromDir(t *testing.T) {
	dir := filepath.Join("..", "..", "testdata", "packages", "demo_extra")
	if _, err := os.Stat(dir); err != nil {
		t.Skip("demo_extra package not present")
	}
	t.Setenv("MOONBASIC_CACHE", t.TempDir())
	if err := Install(dir, nil); err != nil {
		t.Fatal(err)
	}
	list, err := AllInstalledRootDirs()
	if err != nil {
		t.Fatal(err)
	}
	if len(list) != 1 {
		t.Fatalf("got %d roots", len(list))
	}
}
