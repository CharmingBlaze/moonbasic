package pkgmgr

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// UploadGitHubReleaseAsset uploads assetPath to an existing GitHub release tag.
// Requires MOONBASIC_PUBLISH_REPO=owner/repo and MOONBASIC_PUBLISH_TAG=tag and GITHUB_TOKEN (or GH_TOKEN).
func UploadGitHubReleaseAsset(assetPath string) error {
	repo := strings.TrimSpace(os.Getenv("MOONBASIC_PUBLISH_REPO"))
	tag := strings.TrimSpace(os.Getenv("MOONBASIC_PUBLISH_TAG"))
	token := strings.TrimSpace(os.Getenv("GITHUB_TOKEN"))
	if token == "" {
		token = strings.TrimSpace(os.Getenv("GH_TOKEN"))
	}
	if repo == "" || tag == "" || token == "" {
		return fmt.Errorf("publish: set MOONBASIC_PUBLISH_REPO, MOONBASIC_PUBLISH_TAG, and GITHUB_TOKEN")
	}
	parts := strings.SplitN(repo, "/", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return fmt.Errorf("publish: MOONBASIC_PUBLISH_REPO must be owner/repo")
	}
	owner, rname := parts[0], parts[1]

	client := &http.Client{Timeout: 120 * time.Second}
	api := fmt.Sprintf("https://api.github.com/repos/%s/%s/releases/tags/%s", owner, rname, tag)
	req, err := http.NewRequest(http.MethodGet, api, nil)
	if err != nil {
		return err
	}
	req.Header.Set("Accept", "application/vnd.github+json")
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("github release %s: %s: %s", tag, resp.Status, truncate(string(body), 200))
	}
	var rel struct {
		ID int64 `json:"id"`
	}
	if err := json.Unmarshal(body, &rel); err != nil {
		return fmt.Errorf("github: decode release: %w", err)
	}

	data, err := os.ReadFile(assetPath)
	if err != nil {
		return err
	}
	base := filepath.Base(assetPath)
	uploadURL := fmt.Sprintf("https://uploads.github.com/repos/%s/%s/releases/%d/assets?name=%s", owner, rname, rel.ID, base)
	req2, err := http.NewRequest(http.MethodPost, uploadURL, bytes.NewReader(data))
	if err != nil {
		return err
	}
	req2.Header.Set("Accept", "application/vnd.github+json")
	req2.Header.Set("Authorization", "Bearer "+token)
	req2.Header.Set("Content-Type", "application/octet-stream")
	resp2, err := client.Do(req2)
	if err != nil {
		return err
	}
	defer resp2.Body.Close()
	body2, _ := io.ReadAll(resp2.Body)
	if resp2.StatusCode != http.StatusCreated {
		return fmt.Errorf("github upload: %s: %s", resp2.Status, truncate(string(body2), 200))
	}
	fmt.Printf("uploaded release asset %s\n", base)
	return nil
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}
