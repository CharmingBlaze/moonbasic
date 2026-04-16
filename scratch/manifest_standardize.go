package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type Command struct {
	Key          string   `json:"key"`
	Args         []string `json:"args,omitempty"`
	Returns      string   `json:"returns,omitempty"`
	Description  string   `json:"description,omitempty"`
	Namespace    string   `json:"namespace,omitempty"`
	Phase        string   `json:"phase,omitempty"`
	DeprecatedOf string   `json:"deprecated_of,omitempty"`
	Example      string   `json:"example,omitempty"`
	Errors       []string `json:"errors,omitempty"`
}

type Manifest struct {
	Commands []Command `json:"commands"`
}

func main() {
	data, err := os.ReadFile("compiler/builtinmanifest/commands.json")
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return
	}

	var m Manifest
	if err := json.Unmarshal(data, &m); err != nil {
		fmt.Printf("Error unmarshaling JSON: %v\n", err)
		return
	}

	newCommands := make([]Command, 0, len(m.Commands))
	canonicalKeys := make(map[string]bool)
	for _, cmd := range m.Commands {
		canonicalKeys[cmd.Key] = true
	}

	for _, cmd := range m.Commands {
		key := cmd.Key
		changed := false
		newKey := key

		// 1. Rename MAKE to CREATE
		if strings.HasSuffix(key, ".MAKE") {
			newKey = strings.TrimSuffix(key, "MAKE") + "CREATE"
			changed = true
		} else if strings.Contains(key, ".MAKE") && !strings.HasSuffix(key, ".MAKE") {
			// e.g. MODEL.MAKECUBE -> MODEL.CREATECUBE
			newKey = strings.Replace(key, ".MAKE", ".CREATE", 1)
			changed = true
		}

		// 2. Rename SETPOSITION to SETPOS
		if strings.HasSuffix(key, ".SETPOSITION") {
			newKey = strings.TrimSuffix(key, "SETPOSITION") + "SETPOS"
			changed = true
		}

		if changed {
			// If the new canonical key doesn't exist, rename this one and add old as alias
			// Or if it DOES exist, this one becomes the alias.
			if !canonicalKeys[newKey] {
				fmt.Printf("Renaming %s to %s\n", key, newKey)
				
				// Create the new canonical command
				canonicalCmd := cmd
				canonicalCmd.Key = newKey
				
				// Create the legacy alias
				aliasCmd := cmd
				aliasCmd.DeprecatedOf = newKey
				aliasCmd.Description = fmt.Sprintf("DEPRECATED alias of %s. Use %s.", newKey, newKey)
				
				newCommands = append(newCommands, canonicalCmd)
				newCommands = append(newCommands, aliasCmd)
				
				// Mark as seen so we don't add it twice if there are multiple overloads
				canonicalKeys[newKey] = true
			} else {
				fmt.Printf("Key %s already has canonical %s. Marking as deprecated.\n", key, newKey)
				aliasCmd := cmd
				aliasCmd.DeprecatedOf = newKey
				if !strings.Contains(aliasCmd.Description, "DEPRECATED") {
					aliasCmd.Description = fmt.Sprintf("DEPRECATED alias of %s. Use %s. %s", newKey, newKey, aliasCmd.Description)
				}
				newCommands = append(newCommands, aliasCmd)
			}
		} else {
			newCommands = append(newCommands, cmd)
		}
	}

	m.Commands = newCommands
	
	// Sort or just write back
	out, err := json.MarshalIndent(m, "", "  ")
	if err != nil {
		fmt.Printf("Error marshaling JSON: %v\n", err)
		return
	}

	if err := os.WriteFile("compiler/builtinmanifest/commands.json", out, 0644); err != nil {
		fmt.Printf("Error writing file: %v\n", err)
		return
	}
	fmt.Println("Manifest migration complete.")
}
