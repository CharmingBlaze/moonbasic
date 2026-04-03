package docgen

import (
	"encoding/json"
	"fmt"
	"html"
	"os"
	"path/filepath"
	"strings"
)

// WriteSite writes index.html, per-command pages, search-index.json, and app.js into outDir.
func WriteSite(outDir string, cmds []CommandDoc) error {
	if err := os.MkdirAll(filepath.Join(outDir, "commands"), 0755); err != nil {
		return err
	}
	if err := writeIndex(outDir); err != nil {
		return err
	}
	if err := writeAppJS(outDir); err != nil {
		return err
	}
	var search []searchEntry
	for _, c := range cmds {
		fn := fileNameForKey(c.Key)
		if err := writeCommandPage(outDir, fn, c); err != nil {
			return err
		}
		body := buildSearchBody(c)
		search = append(search, searchEntry{
			Title: c.Key,
			URL:   "commands/" + fn + ".html",
			Body:  body,
		})
	}
	idxPath := filepath.Join(outDir, "search-index.json")
	f, err := os.Create(idxPath)
	if err != nil {
		return err
	}
	defer f.Close()
	enc := json.NewEncoder(f)
	enc.SetIndent("", "  ")
	return enc.Encode(search)
}

type searchEntry struct {
	Title string `json:"title"`
	URL   string `json:"url"`
	Body  string `json:"body"`
}

func fileNameForKey(key string) string {
	s := strings.ReplaceAll(key, ".", "_")
	s = strings.ReplaceAll(s, "$", "str")
	return s
}

func buildSearchBody(c CommandDoc) string {
	var b strings.Builder
	b.WriteString(c.Key)
	b.WriteByte(' ')
	b.WriteString(strings.Join(c.Args, " "))
	b.WriteByte(' ')
	b.WriteString(c.Description)
	b.WriteByte(' ')
	b.WriteString(c.Example)
	return b.String()
}

func writeCommandPage(outDir, fn string, c CommandDoc) error {
	path := filepath.Join(outDir, "commands", fn+".html")
	sig := signatureHTML(c)
	var desc, ex, errSec strings.Builder
	if c.Description != "" {
		desc.WriteString("<p>" + html.EscapeString(c.Description) + "</p>")
	} else {
		desc.WriteString("<p><em>No description yet.</em></p>")
	}
	if c.Example != "" {
		ex.WriteString("<h2>Example</h2><pre><code>")
		ex.WriteString(html.EscapeString(c.Example))
		ex.WriteString("</code></pre>")
	}
	if len(c.Errors) > 0 {
		errSec.WriteString("<h2>Error codes</h2><ul>")
		for _, e := range c.Errors {
			errSec.WriteString("<li>" + html.EscapeString(e) + "</li>")
		}
		errSec.WriteString("</ul>")
	} else {
		errSec.WriteString("<h2>Error codes</h2><p>See <a href=\"../#errors\">error categories</a> and <code>compiler/errors/MoonBasic.md</code> in the repo.</p>")
	}
	stub := ""
	if c.Stub != "" {
		stub = "<p class=\"stub\"><strong>Stub:</strong> " + html.EscapeString(c.Stub) + "</p>"
	}
	page := fmt.Sprintf(`<!DOCTYPE html>
<html lang="en"><head><meta charset="utf-8"><title>%s — moonBASIC</title>
<link rel="stylesheet" href="../style.css"></head><body>
<nav><a href="../index.html">Index</a></nav>
<h1><code>%s</code></h1>
<p class="sig">%s</p>
%s
%s
%s
%s
</body></html>`, html.EscapeString(c.Key), html.EscapeString(c.Key), sig, stub, desc.String(), ex.String(), errSec.String())
	return os.WriteFile(path, []byte(page), 0644)
}

func signatureHTML(c CommandDoc) string {
	var b strings.Builder
	b.WriteString(html.EscapeString(c.Key))
	b.WriteByte('(')
	for i, a := range c.Args {
		if i > 0 {
			b.WriteString(", ")
		}
		b.WriteString(html.EscapeString(a))
	}
	b.WriteByte(')')
	if c.Returns != "" {
		b.WriteString(" → ")
		b.WriteString(html.EscapeString(c.Returns))
	}
	return b.String()
}

func writeIndex(outDir string) error {
	html := `<!DOCTYPE html>
<html lang="en"><head><meta charset="utf-8"><title>moonBASIC command reference</title>
<link rel="stylesheet" href="style.css"></head><body>
<h1>moonBASIC command reference</h1>
<p>Generated from <code>commands.json</code>. Use the search box to filter commands.</p>
<input id="q" type="search" placeholder="Search commands..." autocomplete="off">
<ul id="results"></ul>
<section id="errors"><h2>Error categories</h2>
<ul>
<li><strong>Lexer</strong> — invalid characters, unterminated strings</li>
<li><strong>Parse</strong> — syntax structure</li>
<li><strong>Type</strong> — argument kinds, unknown builtins</li>
<li><strong>Runtime</strong> — native engine failures</li>
<li><strong>CodeGen</strong> — internal compiler IR issues</li>
</ul>
<p>Formatted messages use <code>[moonBASIC] Category in file line col:</code> (see <code>compiler/errors</code>).</p>
</section>
<script src="app.js" defer></script>
</body></html>`
	return os.WriteFile(filepath.Join(outDir, "index.html"), []byte(html), 0644)
}

func writeAppJS(outDir string) error {
	js := `fetch('search-index.json').then(r=>r.json()).then(idx=>{
  const ul=document.getElementById('results');
  const inp=document.getElementById('q');
  function render(filter){
    ul.innerHTML='';
    const q=(filter||'').toLowerCase();
    idx.forEach(item=>{
      if(!q || item.title.toLowerCase().includes(q) || item.body.toLowerCase().includes(q)){
        const li=document.createElement('li');
        const a=document.createElement('a');
        a.href=item.url; a.textContent=item.title;
        li.appendChild(a);
        ul.appendChild(li);
      }
    });
  }
  inp.addEventListener('input',()=>render(inp.value));
  render('');
});`
	return os.WriteFile(filepath.Join(outDir, "app.js"), []byte(js), 0644)
}
