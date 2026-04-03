package pkgmgr

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

// Index lists remote packages (registry v1 index.json).
type Index struct {
	Packages map[string]IndexEntry `json:"packages"`
}

// IndexEntry describes one downloadable package.
type IndexEntry struct {
	Version string `json:"version"`
	URL     string `json:"url"`
	SHA256  string `json:"sha256,omitempty"`
}

// Registry resolves package names to download URLs.
type Registry interface {
	Lookup(name string) (IndexEntry, error)
}

// HTTTPRegistry fetches an index JSON from a base URL (MOONBASIC_REGISTRY).
type HTTTPRegistry struct {
	IndexURL string
	Client   *http.Client
}

// Lookup implements Registry.
func (r *HTTTPRegistry) Lookup(name string) (IndexEntry, error) {
	if r.IndexURL == "" {
		return IndexEntry{}, fmt.Errorf("registry: MOONBASIC_REGISTRY is not set (or use a file/URL path with install)")
	}
	client := r.Client
	if client == nil {
		client = &http.Client{Timeout: 30 * time.Second}
	}
	resp, err := client.Get(r.IndexURL)
	if err != nil {
		return IndexEntry{}, fmt.Errorf("registry: fetch index: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return IndexEntry{}, fmt.Errorf("registry: index HTTP %s", resp.Status)
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return IndexEntry{}, err
	}
	var idx Index
	if err := json.Unmarshal(data, &idx); err != nil {
		return IndexEntry{}, fmt.Errorf("registry: index JSON: %w", err)
	}
	e, ok := idx.Packages[name]
	if !ok {
		return IndexEntry{}, fmt.Errorf("registry: unknown package %q", name)
	}
	if e.URL == "" {
		return IndexEntry{}, fmt.Errorf("registry: package %q has no url", name)
	}
	return e, nil
}

// FileRegistry reads a local index JSON path.
type FileRegistry struct {
	Path string
}

// Lookup implements Registry.
func (r *FileRegistry) Lookup(name string) (IndexEntry, error) {
	data, err := os.ReadFile(r.Path)
	if err != nil {
		return IndexEntry{}, err
	}
	var idx Index
	if err := json.Unmarshal(data, &idx); err != nil {
		return IndexEntry{}, fmt.Errorf("registry: index JSON: %w", err)
	}
	e, ok := idx.Packages[name]
	if !ok {
		return IndexEntry{}, fmt.Errorf("registry: unknown package %q", name)
	}
	return e, nil
}

// DefaultRegistryFromEnv returns a registry from MOONBASIC_REGISTRY (URL or file path).
func DefaultRegistryFromEnv() Registry {
	raw := strings.TrimSpace(os.Getenv("MOONBASIC_REGISTRY"))
	if raw == "" {
		return &HTTTPRegistry{IndexURL: ""}
	}
	if strings.HasPrefix(raw, "http://") || strings.HasPrefix(raw, "https://") {
		return &HTTTPRegistry{IndexURL: raw}
	}
	return &FileRegistry{Path: raw}
}
