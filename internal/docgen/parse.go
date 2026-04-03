package docgen

import (
	"encoding/json"
	"os"
)

// CommandDoc is one manifest entry plus optional prose fields.
type CommandDoc struct {
	Key         string   `json:"key"`
	Args        []string `json:"args"`
	Returns     string   `json:"returns,omitempty"`
	Pure        bool     `json:"pure,omitempty"`
	Phase       string   `json:"phase,omitempty"`
	Namespace   string   `json:"namespace,omitempty"`
	Stub        string   `json:"stub,omitempty"`
	Description string   `json:"description,omitempty"`
	Example     string   `json:"example,omitempty"`
	Errors      []string `json:"errors,omitempty"`
}

type jsonRoot struct {
	Commands []CommandDoc `json:"commands"`
}

// LoadCommandsJSON reads commands.json (compiler manifest shape + optional doc fields).
func LoadCommandsJSON(path string) ([]CommandDoc, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var root jsonRoot
	if err := json.Unmarshal(data, &root); err != nil {
		return nil, err
	}
	return root.Commands, nil
}
