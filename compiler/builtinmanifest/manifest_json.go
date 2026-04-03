package builtinmanifest

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"strings"
)

//go:embed commands.json
var embeddedCommandsJSON []byte

type jsonRoot struct {
	Commands []jsonCommand `json:"commands"`
}

type jsonCommand struct {
	Key       string   `json:"key"`
	Args      []string `json:"args"`
	Returns   string   `json:"returns,omitempty"`
	Pure      bool     `json:"pure,omitempty"`
	Phase     string   `json:"phase,omitempty"`
	Namespace string   `json:"namespace,omitempty"`
	Stub      string   `json:"stub,omitempty"`
}

func parseArgKind(s string) (ArgKind, error) {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "any":
		return Any, nil
	case "int":
		return Int, nil
	case "float":
		return Float, nil
	case "string":
		return String, nil
	case "bool":
		return Bool, nil
	case "handle":
		return Handle, nil
	default:
		return Any, fmt.Errorf("unknown arg kind %q", s)
	}
}

// ParseJSON builds a manifest table from the standard JSON schema.
func ParseJSON(data []byte) (*Table, error) {
	var root jsonRoot
	if err := json.Unmarshal(data, &root); err != nil {
		return nil, err
	}
	t := &Table{Commands: make(map[string][]Command)}
	for _, jc := range root.Commands {
		key := NormalizeCommand(jc.Key)
		if key == "" {
			return nil, fmt.Errorf("empty command key in manifest")
		}
		args := make([]ArgKind, len(jc.Args))
		for i, a := range jc.Args {
			k, err := parseArgKind(a)
			if err != nil {
				return nil, fmt.Errorf("%s: %w", key, err)
			}
			args[i] = k
		}
		cmd := Command{
			Key:       key,
			Args:      args,
			Returns:   jc.Returns,
			Pure:      jc.Pure,
			Phase:     jc.Phase,
			Namespace: jc.Namespace,
			Stub:      jc.Stub,
		}
		t.Commands[key] = append(t.Commands[key], cmd)
	}
	return t, nil
}

func mustDefaultTable() *Table {
	t, err := ParseJSON(embeddedCommandsJSON)
	if err != nil {
		panic("builtinmanifest: embedded JSON: " + err.Error())
	}
	return t
}
