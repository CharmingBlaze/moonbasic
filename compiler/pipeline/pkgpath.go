package pipeline

import (
	"os"
	"path/filepath"
	"strings"

	"moonbasic/compiler/include"
	"moonbasic/internal/pkgmgr"
)

// SyncPackageIncludeRoots updates INCLUDE search paths from MOONBASIC_PATH and installed packages.
func SyncPackageIncludeRoots() {
	var roots []string
	if p := os.Getenv("MOONBASIC_PATH"); p != "" {
		for _, s := range filepath.SplitList(p) {
			s = strings.TrimSpace(s)
			if s == "" {
				continue
			}
			if abs, err := filepath.Abs(s); err == nil {
				roots = append(roots, abs)
			}
		}
	}
	if extra, err := pkgmgr.AllInstalledRootDirs(); err == nil {
		roots = append(roots, extra...)
	}
	include.SetPackageRoots(roots)
}
