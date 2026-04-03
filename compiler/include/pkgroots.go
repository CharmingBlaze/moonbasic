package include

import (
	"os"
	"path/filepath"
	"sync"
)

var (
	pkgMu    sync.RWMutex
	pkgRoots []string
)

// SetPackageRoots sets extra directories searched when INCLUDE cannot resolve from the host file.
// Pass nil or empty to clear.
func SetPackageRoots(roots []string) {
	pkgMu.Lock()
	defer pkgMu.Unlock()
	if len(roots) == 0 {
		pkgRoots = nil
		return
	}
	pkgRoots = append([]string(nil), roots...)
}

func packageRootsSnapshot() []string {
	pkgMu.RLock()
	defer pkgMu.RUnlock()
	return append([]string(nil), pkgRoots...)
}

// TryOpenInclude searches roots for includePath (as in the source) and returns path and contents if found.
func TryOpenInclude(includePath string) (absPath string, data []byte, err error) {
	for _, root := range packageRootsSnapshot() {
		cand := filepath.Join(root, filepath.FromSlash(includePath))
		cand = filepath.Clean(cand)
		if !filepath.IsAbs(cand) {
			var e error
			cand, e = filepath.Abs(cand)
			if e != nil {
				continue
			}
		}
		rootAbs, err := filepath.Abs(root)
		if err != nil {
			continue
		}
		if rel, err := filepath.Rel(rootAbs, cand); err != nil || rel == ".." || len(rel) >= 3 && rel[:3] == ".."+string(filepath.Separator) {
			continue
		}
		b, err := os.ReadFile(cand)
		if err == nil {
			return cand, b, nil
		}
	}
	return "", nil, os.ErrNotExist
}
