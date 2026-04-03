package pkgmgr

import (
	"os"
	"path/filepath"
	goruntime "runtime"
)

// CacheDir returns the root directory for installed packages.
func CacheDir() (string, error) {
	if d := os.Getenv("MOONBASIC_CACHE"); d != "" {
		return filepath.Clean(d), nil
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	switch goruntime.GOOS {
	case "windows":
		if la := os.Getenv("LOCALAPPDATA"); la != "" {
			return filepath.Join(la, "moonbasic", "packages"), nil
		}
		return filepath.Join(home, "AppData", "Local", "moonbasic", "packages"), nil
	default:
		return filepath.Join(home, ".local", "share", "moonbasic", "packages"), nil
	}
}

// InstallPath is the directory for one installed version.
func InstallPath(name, version string) (string, error) {
	root, err := CacheDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(root, name, version), nil
}

// AllInstalledRootDirs returns every version directory under the package cache (for INCLUDE search).
func AllInstalledRootDirs() ([]string, error) {
	root, err := CacheDir()
	if err != nil {
		return nil, err
	}
	fi, err := os.Stat(root)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	if !fi.IsDir() {
		return nil, nil
	}
	var out []string
	pkgs, err := os.ReadDir(root)
	if err != nil {
		return nil, err
	}
	for _, p := range pkgs {
		if !p.IsDir() {
			continue
		}
		pdir := filepath.Join(root, p.Name())
		vers, err := os.ReadDir(pdir)
		if err != nil {
			continue
		}
		for _, v := range vers {
			if !v.IsDir() {
				continue
			}
			inst := filepath.Join(pdir, v.Name())
			out = append(out, inst)
		}
	}
	return out, nil
}
