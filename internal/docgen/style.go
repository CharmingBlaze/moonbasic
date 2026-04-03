package docgen

import (
	"os"
	"path/filepath"
)

// WriteStyle writes minimal CSS.
func WriteStyle(outDir string) error {
	css := `body{font-family:system-ui,sans-serif;max-width:52rem;margin:1rem auto;padding:0 1rem;line-height:1.45}
code,pre{font-family:ui-monospace,monospace}
#q{width:100%;max-width:28rem;padding:.4rem;font-size:1rem}
#results{column-count:2;list-style:none;padding:0}
#results li{margin:.15rem 0}
nav{margin-bottom:1rem}
.sig{font-size:1.1rem}
.stub{color:#666}
pre{background:#f4f4f4;padding:.75rem;overflow:auto}`
	return os.WriteFile(filepath.Join(outDir, "style.css"), []byte(css), 0644)
}
