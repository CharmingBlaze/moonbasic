// Package lsp implements a minimal Language Server Protocol (stdio) for moonBASIC.
package lsp

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"moonbasic/compiler/builtinmanifest"
	moonerrors "moonbasic/compiler/errors"
	"moonbasic/compiler/pipeline"
)

// Serve runs the LSP server on stdin/stdout until shutdown or EOF.
func Serve() error {
	s := &server{
		docs:  make(map[string]string),
		table: builtinmanifest.Default(),
	}
	br := bufio.NewReader(os.Stdin)
	for {
		body, err := readFramedMessage(br)
		if err != nil {
			if errors.Is(err, io.EOF) {
				return nil
			}
			return err
		}
		var env struct {
			JSONRPC string           `json:"jsonrpc"`
			ID      *json.RawMessage `json:"id"`
			Method  string           `json:"method"`
			Params  json.RawMessage  `json:"params"`
		}
		if err := json.Unmarshal(body, &env); err != nil {
			continue
		}
		if env.Method == "" {
			continue
		}
		if env.ID == nil {
			s.handleNotify(env.Method, env.Params)
			continue
		}
		resp := s.handleRequest(env.Method, env.Params)
		out := map[string]any{
			"jsonrpc": "2.0",
			"id":      json.RawMessage(*env.ID),
		}
		if resp.err != nil {
			out["error"] = map[string]any{"code": -32603, "message": resp.err.Error()}
		} else {
			out["result"] = resp.result
		}
		raw, _ := json.Marshal(out)
		if err := writeFramedMessage(os.Stdout, raw); err != nil {
			return err
		}
	}
}

type server struct {
	docs  map[string]string // uri -> full text
	table *builtinmanifest.Table
}

type reqResult struct {
	result any
	err    error
}

func (s *server) handleNotify(method string, params json.RawMessage) {
	switch method {
	case "initialized":
		return
	case "textDocument/didOpen":
		var p struct {
			TextDocument struct {
				URI  string `json:"uri"`
				Text string `json:"text"`
			} `json:"textDocument"`
		}
		_ = json.Unmarshal(params, &p)
		if p.TextDocument.URI != "" {
			s.docs[p.TextDocument.URI] = p.TextDocument.Text
			s.publishDiagnostics(p.TextDocument.URI, p.TextDocument.Text)
		}
	case "textDocument/didChange":
		var p struct {
			TextDocument struct {
				URI string `json:"uri"`
			} `json:"textDocument"`
			ContentChanges []struct {
				Text string `json:"text"`
			} `json:"contentChanges"`
		}
		_ = json.Unmarshal(params, &p)
		if p.TextDocument.URI != "" && len(p.ContentChanges) > 0 {
			t := p.ContentChanges[len(p.ContentChanges)-1].Text
			s.docs[p.TextDocument.URI] = t
			s.publishDiagnostics(p.TextDocument.URI, t)
		}
	case "textDocument/didClose":
		var p struct {
			TextDocument struct {
				URI string `json:"uri"`
			} `json:"textDocument"`
		}
		_ = json.Unmarshal(params, &p)
		delete(s.docs, p.TextDocument.URI)
	case "shutdown":
		return
	case "exit":
		os.Exit(0)
	}
}

func (s *server) handleRequest(method string, params json.RawMessage) reqResult {
	switch method {
	case "initialize":
		return reqResult{result: map[string]any{
			"capabilities": map[string]any{
				"textDocumentSync": 1,
				"hoverProvider":    true,
				"completionProvider": map[string]any{
					"triggerCharacters": []string{"."},
				},
			},
			"serverInfo": map[string]string{"name": "moonbasic-lsp", "version": "0.1"},
		}}
	case "shutdown":
		return reqResult{result: nil}
	case "textDocument/hover":
		return reqResult{result: s.hover(params)}
	case "textDocument/completion":
		return reqResult{result: s.completion(params)}
	default:
		return reqResult{result: nil}
	}
}

func (s *server) publishDiagnostics(uri, text string) {
	path := filePathFromURI(uri)
	name := filepath.Base(path)
	if path == "" {
		name = "buffer.mb"
		path = name
	}
	err := pipeline.CheckSource(path, text)
	var diags []any
	var me *moonerrors.MoonError
	if err != nil && errors.As(err, &me) {
		diags = append(diags, map[string]any{
			"range": map[string]any{
				"start": map[string]uint32{"line": uint32(me.Line - 1), "character": uint32(me.Col - 1)},
				"end":   map[string]uint32{"line": uint32(me.Line - 1), "character": uint32(me.Col + 20)},
			},
			"severity": 1,
			"source":   "moonbasic",
			"message":  me.Message,
		})
	} else if err != nil {
		diags = append(diags, map[string]any{
			"range": map[string]any{
				"start": map[string]uint32{"line": 0, "character": 0},
				"end":   map[string]uint32{"line": 0, "character": 1},
			},
			"severity": 1,
			"source":   "moonbasic",
			"message":  err.Error(),
		})
	}
	notif := map[string]any{
		"jsonrpc": "2.0",
		"method":  "textDocument/publishDiagnostics",
		"params": map[string]any{
			"uri":         uri,
			"diagnostics": diags,
		},
	}
	raw, _ := json.Marshal(notif)
	_ = writeFramedMessage(os.Stdout, raw)
}

func (s *server) hover(params json.RawMessage) any {
	var p struct {
		TextDocument struct {
			URI string `json:"uri"`
		} `json:"textDocument"`
		Position struct {
			Line      int `json:"line"`
			Character int `json:"character"`
		} `json:"position"`
	}
	_ = json.Unmarshal(params, &p)
	text := s.docs[p.TextDocument.URI]
	if text == "" {
		return nil
	}
	lines := strings.Split(text, "\n")
	if p.Position.Line < 0 || p.Position.Line >= len(lines) {
		return nil
	}
	line := lines[p.Position.Line]
	key, ok := dottedCommandAt(line, p.Position.Character)
	if !ok {
		return nil
	}
	cmd, found := s.table.FirstOverload(key)
	if !found {
		if alt, ok2 := s.table.BestSimilarKey(key, 3); ok2 {
			return map[string]any{
				"contents": map[string]any{
					"kind":  "markdown",
					"value": fmt.Sprintf("Unknown command `%s`. Did you mean **`%s`**?", key, alt),
				},
			}
		}
		return nil
	}
	doc := formatCommandDoc(key, cmd)
	return map[string]any{
		"contents": map[string]any{
			"kind":  "markdown",
			"value": doc,
		},
	}
}

func formatCommandDoc(key string, c builtinmanifest.Command) string {
	var b strings.Builder
	fmt.Fprintf(&b, "### `%s`\n\n", key)
	if len(c.Args) > 0 {
		b.WriteString("**Arguments:** ")
		parts := make([]string, len(c.Args))
		for i, a := range c.Args {
			parts[i] = argKindName(a)
		}
		b.WriteString(strings.Join(parts, ", "))
		b.WriteString("\n\n")
	}
	if c.Returns != "" {
		fmt.Fprintf(&b, "**Returns:** `%s`\n\n", c.Returns)
	}
	if c.Phase != "" {
		fmt.Fprintf(&b, "**Phase:** `%s`\n\n", c.Phase)
	}
	if c.Pure {
		b.WriteString("**Pure:** yes\n\n")
	}
	if c.Stub != "" {
		fmt.Fprintf(&b, "> %s\n", c.Stub)
	}
	return b.String()
}

func argKindName(k builtinmanifest.ArgKind) string {
	switch k {
	case builtinmanifest.Int:
		return "int"
	case builtinmanifest.Float:
		return "float"
	case builtinmanifest.String:
		return "string"
	case builtinmanifest.Bool:
		return "bool"
	case builtinmanifest.Handle:
		return "handle"
	default:
		return "any"
	}
}

var dottedCmdRE = regexp.MustCompile(`\b([A-Z][A-Z0-9_#]*)\s*\.\s*([A-Z][A-Z0-9_#]*)`)

func dottedCommandAt(line string, col int) (string, bool) {
	if col < 0 {
		return "", false
	}
	for _, loc := range dottedCmdRE.FindAllStringSubmatchIndex(line, -1) {
		fullStart, fullEnd := loc[0], loc[1]
		if col >= fullStart && col <= fullEnd {
			ns := line[loc[2]:loc[3]]
			meth := line[loc[4]:loc[5]]
			return builtinmanifest.Key(ns, meth), true
		}
	}
	return "", false
}

func (s *server) completion(params json.RawMessage) any {
	var p struct {
		TextDocument struct {
			URI string `json:"uri"`
		} `json:"textDocument"`
		Position struct {
			Line      int `json:"line"`
			Character int `json:"character"`
		} `json:"position"`
	}
	_ = json.Unmarshal(params, &p)
	text := s.docs[p.TextDocument.URI]
	if text == "" {
		return map[string]any{"isIncomplete": false, "items": []any{}}
	}
	lines := strings.Split(text, "\n")
	if p.Position.Line < 0 || p.Position.Line >= len(lines) {
		return map[string]any{"isIncomplete": false, "items": []any{}}
	}
	line := lines[p.Position.Line]
	col := p.Position.Character
	if col > len(line) {
		col = len(line)
	}
	prefix := line[:col]
	dot := strings.LastIndex(prefix, ".")
	if dot < 0 {
		return map[string]any{"isIncomplete": false, "items": []any{}}
	}
	nsPart := strings.TrimSpace(prefix[:dot])
	if i := strings.LastIndexAny(nsPart, " \t(,:"); i >= 0 {
		nsPart = strings.TrimSpace(nsPart[i+1:])
	}
	nsPart = strings.ToUpper(nsPart)
	keys := s.table.KeysWithNamespacePrefix(nsPart)
	var items []any
	for _, k := range keys {
		suf := strings.TrimPrefix(k, nsPart+".")
		items = append(items, map[string]any{
			"label":      suf,
			"kind":       3,
			"insertText": suf,
		})
	}
	return map[string]any{"isIncomplete": false, "items": items}
}

func filePathFromURI(raw string) string {
	u, err := url.Parse(raw)
	if err != nil || u.Scheme != "file" {
		return ""
	}
	p := u.Path
	unescaped, err := url.PathUnescape(p)
	if err == nil && unescaped != "" {
		p = unescaped
	}
	if len(p) >= 3 && p[0] == '/' && p[2] == ':' {
		p = p[1:]
	}
	return filepath.FromSlash(p)
}
