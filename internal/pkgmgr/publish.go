package pkgmgr

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// PublishPack validates a package directory and writes a zip to outPath (e.g. mylib-1.0.0.zip).
func PublishPack(dir, outPath string) error {
	dir = filepath.Clean(dir)
	man, err := ParseManifestFile(filepath.Join(dir, "manifest.json"))
	if err != nil {
		return fmt.Errorf("publish: %w", err)
	}
	if outPath == "" {
		outPath = fmt.Sprintf("%s-%s.zip", man.Name, man.Version)
	}
	if err := publishPack(dir, outPath); err != nil {
		return err
	}
	fmt.Printf("wrote %s\n", outPath)
	return nil
}

// PublishReader returns a zip as an in-memory reader for uploads (caller closes).
func PublishReader(dir string) (name string, r io.ReadCloser, size int64, err error) {
	dir = filepath.Clean(dir)
	man, err := ParseManifestFile(filepath.Join(dir, "manifest.json"))
	if err != nil {
		return "", nil, 0, err
	}
	mbc := filepath.Join(dir, man.EntryMBC)
	if _, err := os.Stat(mbc); err != nil {
		return "", nil, 0, err
	}
	tmp, err := os.CreateTemp("", "moonbasic-publish-*.zip")
	if err != nil {
		return "", nil, 0, err
	}
	path := tmp.Name()
	_ = tmp.Close()
	defer func() {
		if err != nil {
			os.Remove(path)
		}
	}()
	if err = PublishPackToPath(dir, path); err != nil {
		return "", nil, 0, err
	}
	fi, err := os.Stat(path)
	if err != nil {
		return "", nil, 0, err
	}
	f, err := os.Open(path)
	if err != nil {
		return "", nil, 0, err
	}
	zname := fmt.Sprintf("%s-%s.zip", man.Name, man.Version)
	return zname, &unlinkOnClose{f: f, path: path}, fi.Size(), nil
}

// PublishPackToPath like PublishPack but explicit output path without printing.
func PublishPackToPath(dir, outPath string) error {
	return publishPack(dir, outPath)
}

func publishPack(dir, outPath string) error {
	dir = filepath.Clean(dir)
	man, err := ParseManifestFile(filepath.Join(dir, "manifest.json"))
	if err != nil {
		return err
	}
	mbc := filepath.Join(dir, man.EntryMBC)
	if _, err := os.Stat(mbc); err != nil {
		return err
	}
	f, err := os.Create(outPath)
	if err != nil {
		return err
	}
	defer f.Close()
	zw := zip.NewWriter(f)
	add := func(rel string) error {
		data, err := os.ReadFile(filepath.Join(dir, rel))
		if err != nil {
			return err
		}
		w, err := zw.Create(filepath.ToSlash(rel))
		if err != nil {
			return err
		}
		_, err = w.Write(data)
		return err
	}
	if err := add("manifest.json"); err != nil {
		zw.Close()
		return err
	}
	if err := add(man.EntryMBC); err != nil {
		zw.Close()
		return err
	}
	entries, err := os.ReadDir(dir)
	if err != nil {
		zw.Close()
		return err
	}
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		name := e.Name()
		if name == "manifest.json" || name == man.EntryMBC {
			continue
		}
		if strings.EqualFold(filepath.Ext(name), ".mb") {
			if err := add(name); err != nil {
				zw.Close()
				return err
			}
		}
	}
	return zw.Close()
}

type unlinkOnClose struct {
	f    *os.File
	path string
}

func (u *unlinkOnClose) Read(p []byte) (int, error) { return u.f.Read(p) }
func (u *unlinkOnClose) Close() error {
	err := u.f.Close()
	_ = os.Remove(u.path)
	return err
}
