package include

import (
	"path/filepath"
)

// Resolve resolves an INCLUDE path relative to the host file's directory, or returns a clean absolute path.
func Resolve(hostFile, includePath string) (string, error) {
	if filepath.IsAbs(includePath) {
		return filepath.Clean(includePath), nil
	}
	dir := filepath.Dir(hostFile)
	if dir == "" || dir == "." {
		dir = "."
	}
	return filepath.Join(dir, filepath.FromSlash(includePath)), nil
}
