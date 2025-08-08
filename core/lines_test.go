package core

import (
	"reflect"
	"testing"
)

func TestToLineMap(t *testing.T) {
	content := "line 1\nline 2\nline 3"
	expected := map[string]string{
		"1": "line 1",
		"2": "line 2",
		"3": "line 3",
	}
	result := ToLineMap(content)
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
	result := FromLineMap(lineMap)
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
	result := FromLineMap(lineMap)
	if result != expected {
		t.Errorf("FromLineMap() with invalid key = %q, want %q", result, expected)
	}
}

func TestApplyPatch(t *testing.T) {
	originalContent := "line 1\nline 2\nline 3"
	patchJSON := `{"2": "new line 2"}`
	expected := "line 1\nnew line 2\nline 3"
	result, err := ApplyPatch(originalContent, patchJSON)
	if err != nil {
		t.Fatalf("ApplyPatch() error = %v", err)
	}
	if result != expected {
		t.Errorf("ApplyPatch() = %q, want %q", result, expected)
	}
}

func TestApplyPatch_InvalidJSON(t *testing.T) {
	originalContent := "line 1\nline 2\nline 3"
	patchJSON := `{"2": "new line 2"`
	_, err := ApplyPatch(originalContent, patchJSON)
	if err == nil {
		t.Error("ApplyPatch() with invalid JSON should have returned an error")
	}
}
