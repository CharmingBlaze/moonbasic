package pkgmgr

import (
	"archive/zip"
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Install installs a package from a local directory, zip file, http(s) URL, or registry name.
func Install(target string, reg Registry) error {
	target = strings.TrimSpace(target)
	if target == "" {
		return fmt.Errorf("install: missing target")
	}

	// Local directory
	if fi, err := os.Stat(target); err == nil && fi.IsDir() {
		return installFromDir(target)
	}
	// Local zip
	if strings.HasSuffix(strings.ToLower(target), ".zip") {
		if fi, err := os.Stat(target); err == nil && !fi.IsDir() {
			return installFromZipFile(target)
		}
	}
	// URL
	if strings.HasPrefix(target, "http://") || strings.HasPrefix(target, "https://") {
		return installFromURL(target, "")
	}
	// Registry name
	if reg == nil {
		reg = DefaultRegistryFromEnv()
	}
	entry, err := reg.Lookup(target)
	if err != nil {
		return err
	}
	return installFromURL(entry.URL, entry.SHA256)
}

func installFromDir(dir string) error {
	manPath := filepath.Join(dir, "manifest.json")
	m, err := ParseManifestFile(manPath)
	if err != nil {
		return fmt.Errorf("install: %w", err)
	}
	mbcPath := filepath.Join(dir, m.EntryMBC)
	if _, err := os.Stat(mbcPath); err != nil {
		return fmt.Errorf("install: entry_mbc %q: %w", m.EntryMBC, err)
	}
	if err := verifyMBCSha(m, mbcPath); err != nil {
		return err
	}
	return copyInstalled(m, func(dst string) error {
		return copyDirContents(dir, dst)
	})
}

func installFromZipFile(zipPath string) error {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return err
	}
	defer r.Close()
	return installFromZipReader(&r.Reader)
}

func installFromZipReader(z *zip.Reader) error {
	var manifestData []byte
	files := map[string]*zip.File{}
	for _, f := range z.File {
		name := filepath.ToSlash(f.Name)
		// skip junk dirs
		base := filepath.Base(name)
		files[base] = f
		if base == "manifest.json" {
			rc, err := f.Open()
			if err != nil {
				return err
			}
			manifestData, err = io.ReadAll(rc)
			rc.Close()
			if err != nil {
				return err
			}
		}
	}
	if len(manifestData) == 0 {
		return fmt.Errorf("install: zip missing manifest.json")
	}
	m, err := ParseManifest(manifestData)
	if err != nil {
		return err
	}
	if _, ok := files[m.EntryMBC]; !ok {
		return fmt.Errorf("install: zip missing %q", m.EntryMBC)
	}
	tmpDir, err := os.MkdirTemp("", "moonbasic-pkg-*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpDir)
	for _, f := range z.File {
		base := filepath.Base(filepath.ToSlash(f.Name))
		if base == "" || strings.Contains(base, "..") {
			continue
		}
		out := filepath.Join(tmpDir, base)
		if err := extractZipFile(f, out); err != nil {
			return err
		}
	}
	mbcPath := filepath.Join(tmpDir, m.EntryMBC)
	if err := verifyMBCSha(m, mbcPath); err != nil {
		return err
	}
	return copyInstalled(m, func(dst string) error {
		return copyDirContents(tmpDir, dst)
	})
}

func installFromURL(url, wantSHA string) error {
	client := &http.Client{Timeout: 120 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return fmt.Errorf("install: download: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("install: download HTTP %s", resp.Status)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if wantSHA != "" {
		sum := sha256.Sum256(body)
		if hex.EncodeToString(sum[:]) != strings.ToLower(wantSHA) {
			return fmt.Errorf("install: sha256 mismatch")
		}
	}
	z, err := zip.NewReader(bytes.NewReader(body), int64(len(body)))
	if err != nil {
		return fmt.Errorf("install: expected zip archive: %w", err)
	}
	return installFromZipReader(z)
}

func extractZipFile(f *zip.File, dest string) error {
	rc, err := f.Open()
	if err != nil {
		return err
	}
	defer rc.Close()
	if err := os.MkdirAll(filepath.Dir(dest), 0755); err != nil {
		return err
	}
	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, rc)
	return err
}

func verifyMBCSha(m *Manifest, mbcPath string) error {
	if m.SHA256MBC == "" {
		return nil
	}
	data, err := os.ReadFile(mbcPath)
	if err != nil {
		return err
	}
	sum := sha256.Sum256(data)
	got := hex.EncodeToString(sum[:])
	if !strings.EqualFold(got, m.SHA256MBC) {
		return fmt.Errorf("install: sha256_mbc mismatch (got %s)", got)
	}
	return nil
}

func copyInstalled(m *Manifest, fill func(dst string) error) error {
	dst, err := InstallPath(m.Name, m.Version)
	if err != nil {
		return err
	}
	if err := os.RemoveAll(dst); err != nil {
		return err
	}
	if err := os.MkdirAll(dst, 0755); err != nil {
		return err
	}
	if err := fill(dst); err != nil {
		return err
	}
	fmt.Printf("installed %s@%s -> %s\n", m.Name, m.Version, dst)
	return nil
}

func copyDirContents(src, dst string) error {
	entries, err := os.ReadDir(src)
	if err != nil {
		return err
	}
	for _, e := range entries {
		name := e.Name()
		if name == "." || name == ".." {
			continue
		}
		from := filepath.Join(src, name)
		to := filepath.Join(dst, name)
		if e.IsDir() {
			if err := os.MkdirAll(to, 0755); err != nil {
				return err
			}
			if err := copyDirContents(from, to); err != nil {
				return err
			}
			continue
		}
		if err := copyFile(from, to); err != nil {
			return err
		}
	}
	return nil
}

func copyFile(from, to string) error {
	b, err := os.ReadFile(from)
	if err != nil {
		return err
	}
	return os.WriteFile(to, b, 0644)
}
