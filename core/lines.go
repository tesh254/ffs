package core

import (
	"encoding/json"
	"sort"
	"strconv"
	"strings"
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
func applyPatch(originalContent string, patchJSON string) (string, error) {
	var patch map[string]string
	if err := json.Unmarshal([]byte(patchJSON), &patch); err != nil {
		return "", err
	}

	lineMap := toLineMap(originalContent)
	for k, v := range patch {
		lineMap[k] = v
	}

	return fromLineMap(lineMap), nil
}
