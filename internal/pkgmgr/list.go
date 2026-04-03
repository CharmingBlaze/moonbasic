package pkgmgr

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// ListInstalled prints installed packages to w (name, version, path).
func ListInstalled(w io.Writer) error {
	root, err := CacheDir()
	if err != nil {
		return err
	}
	fi, err := os.Stat(root)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Fprintf(w, "(no packages installed; cache %s)\n", root)
			return nil
		}
		return err
	}
	if !fi.IsDir() {
		return fmt.Errorf("cache path is not a directory: %s", root)
	}
	names, err := os.ReadDir(root)
	if err != nil {
		return err
	}
	if len(names) == 0 {
		fmt.Fprintf(w, "(empty cache %s)\n", root)
		return nil
	}
	for _, n := range names {
		if !n.IsDir() {
			continue
		}
		pkgDir := filepath.Join(root, n.Name())
		vers, err := os.ReadDir(pkgDir)
		if err != nil {
			continue
		}
		for _, v := range vers {
			if !v.IsDir() {
				continue
			}
			inst := filepath.Join(pkgDir, v.Name())
			manPath := filepath.Join(inst, "manifest.json")
			desc := ""
			if m, err := ParseManifestFile(manPath); err == nil {
				desc = m.Description
			}
			if desc != "" {
				fmt.Fprintf(w, "%s@%s\t%s\t%s\n", n.Name(), v.Name(), inst, desc)
			} else {
				fmt.Fprintf(w, "%s@%s\t%s\n", n.Name(), v.Name(), inst)
			}
		}
	}
	return nil
}
