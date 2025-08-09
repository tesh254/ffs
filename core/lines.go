package core

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/sergi/go-diff/diffmatchpatch"
)

// toLineMap converts the content of a file into a map where keys are line numbers
// (as strings) and values are the content of each line.
func toLineMap(content string) map[string]string {
	lines := strings.Split(content, "\n")
	lineMap := make(map[string]string, len(lines))
	for i, line := range lines {
		lineMap[strconv.Itoa(i+1)] = line
	}
	return lineMap
}

// fromLineMap converts a map of lines back into a single string.
// The lines are sorted by line number before being joined together.
func fromLineMap(lineMap map[string]string) string {
	lines := make([]struct {
		num int
		str string
	}, 0, len(lineMap))

	for k, v := range lineMap {
		num, err := strconv.Atoi(k)
		if err != nil {
			continue // Or handle error appropriately
		}
		lines = append(lines, struct {
			num int
			str string
		}{num, v})
	}

	sort.Slice(lines, func(i, j int) bool {
		return lines[i].num < lines[j].num
	})

	var builder strings.Builder
	for i, line := range lines {
		if i > 0 {
			builder.WriteString("\n")
		}
		builder.WriteString(line.str)
	}
	return builder.String()
}

// applyPatch takes the original content and a JSON string representing the patch,
// then returns the updated content.
func applyPatch(originalContent string, patchJSON string, patchType PatchType) (string, error) {
	var patch map[string]string
	if err := json.Unmarshal([]byte(patchJSON), &patch); err != nil {
		return "", err
	}

	if patchType == PatchTypeAdding {
		return applyAddingPatch(originalContent, patch)
	}

	return applyReplacingPatch(originalContent, patch)
}

func applyReplacingPatch(originalContent string, patch map[string]string) (string, error) {
	lineMap := toLineMap(originalContent)
	for k, v := range patch {
		lineMap[k] = v
	}

	return fromLineMap(lineMap), nil
}

// printDiff prints a colored diff of the changes between two strings.
func printDiff(original, new string) {
	dmp := diffmatchpatch.New()
	a, b, c := dmp.DiffLinesToChars(original, new)
	diffs := dmp.DiffMain(a, b, false)
	result := dmp.DiffCharsToLines(diffs, c)

	for _, diff := range result {
		text := diff.Text
		switch diff.Type {
		case diffmatchpatch.DiffInsert:
			fmt.Printf("\x1b[32m+ %s\x1b[0m", text) // Green
		case diffmatchpatch.DiffDelete:
			fmt.Printf("\x1b[31m- %s\x1b[0m", text) // Red
		case diffmatchpatch.DiffEqual:
			fmt.Print(text)
		}
	}
	fmt.Println()
}

func applyAddingPatch(originalContent string, patch map[string]string) (string, error) {
	lines := strings.Split(originalContent, "\n")
	var patchItems []struct {
		line    int
		content string
	}
	for k, v := range patch {
		lineNum, err := strconv.Atoi(k)
		if err != nil {
			return "", err
		}
		patchItems = append(patchItems, struct {
			line    int
			content string
		}{lineNum, v})
	}
	sort.Slice(patchItems, func(i, j int) bool {
		return patchItems[i].line < patchItems[j].line
	})

	var newLines []string
	origIdx := 0
	patchIdx := 0
	lineCounter := 1

	for origIdx < len(lines) || patchIdx < len(patchItems) {
		if patchIdx < len(patchItems) && patchItems[patchIdx].line == lineCounter {
			newLines = append(newLines, patchItems[patchIdx].content)
			patchIdx++
		} else if origIdx < len(lines) {
			newLines = append(newLines, lines[origIdx])
			origIdx++
		} else {
			newLines = append(newLines, "")
		}
		lineCounter++
	}

	return strings.Join(newLines, "\n"), nil
}
