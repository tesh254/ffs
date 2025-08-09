package agent

import (
	"os"
	"testing"

	"github.com/tesh254/ffs/core"
)

func TestApplySuggestion(t *testing.T) {
	// Create a temporary file with some content
	originalContent := "This is the original content."
	tmpfile, err := os.CreateTemp("", "test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name()) // clean up

	if err = core.WriteFile(tmpfile.Name(), []byte(originalContent)); err != nil {
		t.Fatal(err)
	}

	// Create a suggestion
	newContent := "This is the new and improved content."
	suggestion := Suggestion{
		FilePath:    tmpfile.Name(),
		LineChanges: `{"1": "This is the new and improved content."}`,
	}

	// Apply the suggestion
	appliedContent, err := ApplySuggestion(suggestion)
	if err != nil {
		t.Fatalf("ApplySuggestion failed: %v", err)
	}

	// Check if the returned content is correct
	if appliedContent != newContent {
		t.Errorf("ApplySuggestion returned wrong content: got %q, want %q", appliedContent, newContent)
	}

	// Read the file to check if the content is correct
	fileContent, err := core.ReadFile(tmpfile.Name())
	if err != nil {
		t.Fatalf("Failed to read back file: %v", err)
	}

	if string(fileContent) != newContent {
		t.Errorf("File content mismatch after applying suggestion: got %q, want %q", fileContent, newContent)
	}
}

func TestApplySuggestion_ReadFileError(t *testing.T) {
	suggestion := Suggestion{
		FilePath:    "/non/existent/file",
		LineChanges: `{"1": "This is the new and improved content."}`,
	}

	_, err := ApplySuggestion(suggestion)
	if err == nil {
		t.Fatal("ApplySuggestion should have failed")
	}
}

func TestApplySuggestion_ApplyPatchError(t *testing.T) {
	// Create a temporary file with some content
	originalContent := "This is the original content."
	tmpfile, err := os.CreateTemp("", "test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name()) // clean up

	if err = core.WriteFile(tmpfile.Name(), []byte(originalContent)); err != nil {
		t.Fatal(err)
	}

	// Create a suggestion with invalid JSON
	suggestion := Suggestion{
		FilePath:    tmpfile.Name(),
		LineChanges: `invalid json`,
	}

	_, err = ApplySuggestion(suggestion)
	if err == nil {
		t.Fatal("ApplySuggestion should have failed")
	}
}

func TestApplySuggestion_WriteFileError(t *testing.T) {
	// Create a temporary directory.
	tmpDir, err := os.MkdirTemp("", "test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a temporary file with some content in the temporary directory.
	originalContent := "This is the original content."
	tmpfile, err := os.CreateTemp(tmpDir, "test")
	if err != nil {
		t.Fatal(err)
	}

	if err = core.WriteFile(tmpfile.Name(), []byte(originalContent)); err != nil {
		t.Fatal(err)
	}

	// Make the directory read-only.
	if err = os.Chmod(tmpDir, 0555); err != nil {
		t.Fatal(err)
	}
	defer os.Chmod(tmpDir, 0755) // clean up

	// Create a suggestion
	suggestion := Suggestion{
		FilePath:    tmpfile.Name(),
		LineChanges: `{"1": "This is the new and improved content."}`,
	}

	_, err = ApplySuggestion(suggestion)
	if err == nil {
		t.Fatal("ApplySuggestion should have failed")
	}
}
