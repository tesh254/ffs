package agent

import (
	"github.com/tesh254/ffs/core"
)

// Suggestion represents a proposed change to a file, typically from an LLM agent.
// It targets specific lines to be updated rather than replacing the entire file content.
type Suggestion struct {
	FilePath    string // The path to the file to be modified.
	LineChanges string // A JSON string representing a map of line numbers to their new content.
}

// ApplySuggestion takes a suggestion with line-specific changes and applies them to the file.
// It reads the original file, applies the line changes, and writes the updated content back.
// This approach is more precise and efficient than overwriting the entire file.
func ApplySuggestion(suggestion Suggestion) (string, error) {
	// Read the original file content
	originalContent, err := core.ReadFile(suggestion.FilePath)
	if err != nil {
		return "", err
	}

	// Apply the patch to the original content
	newContent, err := core.ApplyPatch(string(originalContent), suggestion.LineChanges)
	if err != nil {
		return "", err
	}

	// Write the new content back to the file
	if err := core.WriteFile(suggestion.FilePath, []byte(newContent)); err != nil {
		return "", err
	}

	return newContent, nil
}
