package main

import (
	"encoding/json"
	"fmt"
	"os"
	"regexp"
)

type Command struct {
	Key          string   `json:"key"`
	Args         []string `json:"args,omitempty"`
	Returns      string   `json:"returns,omitempty"`
	Namespace    string   `json:"namespace,omitempty"`
	Description  string   `json:"description,omitempty"`
	Example      string   `json:"example,omitempty"`
	DeprecatedOf string   `json:"deprecated_of,omitempty"`
}

type Root struct {
	Commands []map[string]interface{} `json:"commands"`
}

func main() {
	path := "compiler/builtinmanifest/commands.json"
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("Error reading: %v\n", err)
		return
	}

	var root Root
	if err := json.Unmarshal(data, &root); err != nil {
		fmt.Printf("Error parsing: %v\n", err)
		return
	}

	re := regexp.MustCompile(`([a-zA-Z0-9]+)(\$|#|%|\?)`)

	for _, cmd := range root.Commands {
		// Clean description
		if desc, ok := cmd["description"].(string); ok {
			cmd["description"] = re.ReplaceAllString(desc, "$1")
		}
		// Clean example
		if ex, ok := cmd["example"].(string); ok {
			cmd["example"] = re.ReplaceAllString(ex, "$1")
		}
		// Clean args list (if any are strings with suffixes)
		if args, ok := cmd["args"].([]interface{}); ok {
			for i, arg := range args {
				if s, ok := arg.(string); ok {
					args[i] = re.ReplaceAllString(s, "$1")
				}
			}
		}
	}

	out, err := json.MarshalIndent(root, "", "  ")
	if err != nil {
		fmt.Printf("Error marshaling: %v\n", err)
		return
	}

	if err := os.WriteFile(path, out, 0644); err != nil {
		fmt.Printf("Error writing: %v\n", err)
		return
	}
	fmt.Println("Successfully cleaned commands.json")
}
