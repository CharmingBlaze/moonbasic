package pkgmgr

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
)

var nameRe = regexp.MustCompile(`^[a-z0-9][a-z0-9_-]*$`)

// Manifest is the v1 package manifest (manifest.json).
type Manifest struct {
	Name        string            `json:"name"`
	Version     string            `json:"version"`
	Description string            `json:"description,omitempty"`
	Moonbasic   string            `json:"moonbasic,omitempty"`
	EntryMBC    string            `json:"entry_mbc"`
	Deps        map[string]string `json:"deps,omitempty"`
	SHA256MBC   string            `json:"sha256_mbc,omitempty"`
}

// ParseManifestFile reads and validates manifest.json from path.
func ParseManifestFile(path string) (*Manifest, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return ParseManifest(data)
}

// ParseManifest decodes JSON bytes into a validated Manifest.
func ParseManifest(data []byte) (*Manifest, error) {
	var m Manifest
	if err := json.Unmarshal(data, &m); err != nil {
		return nil, fmt.Errorf("manifest: invalid JSON: %w", err)
	}
	if err := m.Validate(); err != nil {
		return nil, err
	}
	return &m, nil
}

// Validate checks required fields and naming rules.
func (m *Manifest) Validate() error {
	if m.Name == "" {
		return fmt.Errorf("manifest: name is required")
	}
	if !nameRe.MatchString(m.Name) {
		return fmt.Errorf("manifest: name %q must match [a-z0-9][a-z0-9_-]*", m.Name)
	}
	if m.Version == "" {
		return fmt.Errorf("manifest: version is required")
	}
	if m.EntryMBC == "" {
		return fmt.Errorf("manifest: entry_mbc is required")
	}
	if filepath.Base(m.EntryMBC) != m.EntryMBC || m.EntryMBC == "." || m.EntryMBC == ".." {
		return fmt.Errorf("manifest: entry_mbc must be a plain filename")
	}
	if filepath.Ext(m.EntryMBC) != ".mbc" {
		return fmt.Errorf("manifest: entry_mbc must end with .mbc")
	}
	return nil
}
