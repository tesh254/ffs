package agent

import (
	"ffs/core"
	"ffs/diff"
)

// Suggestion represents a proposed change to a file, typically from an LLM agent.
type Suggestion struct {
	FilePath    string // The path to the file to be modified.
	NewContent  string // The proposed new content for the file.
}

// ApplySuggestion takes a suggestion, computes a patch, and applies it to the specified file.
// It reads the original file, generates a diff against the suggested new content,
// applies the resulting patch, and writes the updated content back to the file.
// It returns the final content of the file after the patch has been applied.
func ApplySuggestion(suggestion Suggestion) (string, error) {
	// Read the original file content
	originalContent, err := core.ReadFile(suggestion.FilePath)
	if err != nil {
		return "", err
	}

	// Generate a diff between the original content and the new content
	patch := diff.GenerateDiff(string(originalContent), suggestion.NewContent)

	// Apply the patch
	newContent, err := diff.ApplyPatch(string(originalContent), patch)
	if err != nil {
		return "", err
	}

	// Write the new content to the file
	if err := core.WriteFile(suggestion.FilePath, []byte(newContent)); err != nil {
		return "", err
	}

	return newContent, nil
}
