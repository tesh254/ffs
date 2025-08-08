package diff

import (
	"testing"
)

func TestGenerateAndApplyPatch(t *testing.T) {
	text1 := "The quick brown fox jumps over the lazy dog."
	text2 := "The quick yellow fox jumps over the lazy cat."

	// Generate the patch
	patch := GenerateDiff(text1, text2)

	// Apply the patch
	patchedText, err := ApplyPatch(text1, patch)
	if err != nil {
		t.Fatalf("ApplyPatch failed: %v", err)
	}

	// Check if the patched text is the same as the original text2
	if patchedText != text2 {
		t.Errorf("Patched text mismatch: got %q, want %q", patchedText, text2)
	}
}

func TestApplyInvalidPatch(t *testing.T) {
	text1 := "Some text."
	invalidPatch := "invalid patch data"

	// Try to apply the invalid patch
	_, err := ApplyPatch(text1, invalidPatch)
	if err == nil {
		t.Errorf("ApplyPatch should have failed for invalid patch data, but it didn't")
	}
}
