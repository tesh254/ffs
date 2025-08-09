package core

import (
	"bytes"
	"io"
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestToLineMap(t *testing.T) {
	content := "line 1\nline 2\nline 3"
	expected := map[string]string{
		"1": "line 1",
		"2": "line 2",
		"3": "line 3",
	}
	result := toLineMap(content)
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("ToLineMap() = %v, want %v", result, expected)
	}
}

func TestFromLineMap(t *testing.T) {
	lineMap := map[string]string{
		"1": "line 1",
		"3": "line 3",
		"2": "line 2",
	}
	expected := "line 1\nline 2\nline 3"
	result := fromLineMap(lineMap)
	if result != expected {
		t.Errorf("FromLineMap() = %q, want %q", result, expected)
	}
}

func TestFromLineMap_InvalidKey(t *testing.T) {
	lineMap := map[string]string{
		"1":   "line 1",
		"foo": "bar",
		"2":   "line 2",
	}
	expected := "line 1\nline 2"
	result := fromLineMap(lineMap)
	if result != expected {
		t.Errorf("FromLineMap() with invalid key = %q, want %q", result, expected)
	}
}

func TestApplyPatch(t *testing.T) {
	// Test replacing patch
	originalContent := "line1\nline2\nline3"
	patchJSON := `{"2": "new line2"}`
	expectedContent := "line1\nnew line2\nline3"
	newContent, err := applyPatch(originalContent, patchJSON, PatchTypeReplacing)
	if err != nil {
		t.Fatalf("applyPatch failed: %v", err)
	}
	if newContent != expectedContent {
		t.Errorf("applyPatch replacing failed: got %q, want %q", newContent, expectedContent)
	}

	// Test adding patch
	originalContent = "line1\nline3"
	patchJSON = `{"2": "line2"}`
	expectedContent = "line1\nline2\nline3"
	newContent, err = applyPatch(originalContent, patchJSON, PatchTypeAdding)
	if err != nil {
		t.Fatalf("applyPatch failed: %v", err)
	}
	if newContent != expectedContent {
		t.Errorf("applyPatch adding failed: got %q, want %q", newContent, expectedContent)
	}
}

func TestApplyPatch_InvalidJSON(t *testing.T) {
	originalContent := "line 1\nline 2\nline 3"
	patchJSON := `{"2": "new line 2"`
	_, err := applyPatch(originalContent, patchJSON, PatchTypeReplacing)
	if err == nil {
		t.Error("ApplyPatch() with invalid JSON should have returned an error")
	}
}

func TestApplyAddingPatch_MultiLine(t *testing.T) {
	testCases := []struct {
		name     string
		original string
		patch    string
		expected string
	}{
		{
			name:     "insert at beginning",
			original: "line3\nline4",
			patch:    `{"1": "line1", "2": "line2"}`,
			expected: "line1\nline2\nline3\nline4",
		},
		{
			name:     "insert in middle",
			original: "line1\nline4",
			patch:    `{"2": "line2", "3": "line3"}`,
			expected: "line1\nline2\nline3\nline4",
		},
		{
			name:     "insert at end",
			original: "line1\nline2",
			patch:    `{"3": "line3", "4": "line4"}`,
			expected: "line1\nline2\nline3\nline4",
		},
		{
			name:     "insert with gaps",
			original: "line1\nline4",
			patch:    `{"2": "line2", "4": "line3"}`,
			expected: "line1\nline2\nline4\nline3",
		},
		{
			name:     "insert beyond end",
			original: "line1",
			patch:    `{"3": "line3"}`,
			expected: "line1\n\nline3",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			newContent, err := applyPatch(tc.original, tc.patch, PatchTypeAdding)
			if err != nil {
				t.Fatalf("applyPatch failed: %v", err)
			}
			if newContent != tc.expected {
				t.Errorf("applyPatch adding failed: got %q, want %q", newContent, tc.expected)
			}
		})
	}
}

func TestPrintDiff(t *testing.T) {
	original := "line1\nline2\nline3"
	new := "line1\nnew line2\nline3"

	// Capture stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	printDiff(original, new)

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	t.Logf("diff output: %q", output)
	if !strings.Contains(output, "\x1b[31m- line2\n\x1b[0m") {
		t.Errorf("PrintDiff output mismatch: expected red color for removed line")
	}
	if !strings.Contains(output, "\x1b[32m+ new line2\n\x1b[0m") {
		t.Errorf("PrintDiff output mismatch: expected green color for added line")
	}
}
