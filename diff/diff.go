package diff

import (
	"github.com/sergi/go-diff/diffmatchpatch"
)

// GenerateDiff creates a patch string between two texts.
// The patch is generated using the diff-match-patch library and is returned in a compact format.
func GenerateDiff(text1, text2 string) string {
	dmp := diffmatchpatch.New()
	diffs := dmp.DiffMain(text1, text2, false)
	return dmp.DiffToDelta(diffs)
}

// ApplyPatch applies a patch to a text and returns the patched text.
// It returns an error if the patch is invalid or cannot be applied.
func ApplyPatch(text, patch string) (string, error) {
	dmp := diffmatchpatch.New()
	diffs, err := dmp.DiffFromDelta(text, patch)
	if err != nil {
		return "", err
	}
	patchedText := dmp.DiffText2(diffs)
	return patchedText, nil
}
