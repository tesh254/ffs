package core

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

// EditInstruction represents a single edit operation
type EditInstruction struct {
	Action     string `json:"action"`      // "replace" or "insert"
	LineNumber int    `json:"line_number"` // 1-based line number
	NewContent string `json:"new_content"` // Content to insert or replace
}

// FileEditRequest represents the edit instructions
type FileEditRequest struct {
	FilePath string            `json:"file_path"`
	Edits    []EditInstruction `json:"edits"`
}

// readFileLines reads a file and returns its content as a slice of strings.
func readFileLines(filePath string) ([]string, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %v", filePath, err)
	}
	return strings.Split(string(content), "\n"), nil
}

// printDiff prints the diff between two sets of lines with line numbers
func printDiff(topLines, original, updated, bottomLines []string, startLine int, highlight bool) {
	// Define ANSI color codes
	const (
		lightBlue = "\033[94m"      // Bright blue for line numbers
		gray      = "\033[90m"      // Gray for context lines
		white     = "\033[97m"      // White text for original and updated content
		redText   = "\033[31m"      // Red text for removed lines
		greenText = "\033[92m"      // Bright green text for added lines
		redBg     = "\033[41m"      // Red background for highlight
		greenBg   = "\033[48;5;34m" // Darker green background for highlight
		reset     = "\033[0m"       // Reset formatting
	)

	// Select color mode based on highlight flag
	removedPrefix := redText
	addedPrefix := greenText
	removedBg := ""
	addedBg := ""
	if highlight {
		removedPrefix = ""
		addedPrefix = ""
		removedBg = redBg
		addedBg = greenBg
	}

	// Print top context lines
	for i, line := range topLines {
		fmt.Printf("%s%3d| %s%s%s\n", lightBlue, startLine+i, gray, line, reset)
	}
	// Print removed lines (original) in red with white text
	for i, line := range original {
		fmt.Printf("%s%3d| %s%s%s- %s%s\n", lightBlue, startLine+len(topLines)+i, removedPrefix, removedBg, white, line, reset)
	}
	// Print added lines (updated) in green with white text
	for i, line := range updated {
		// Use the edit line number as the base for inserts, incrementing for each added line
		baseLine := startLine + len(topLines)
		if len(original) == 0 {
			baseLine = startLine + len(topLines)
		}
		fmt.Printf("%s%3d| %s%s%s+ %s%s\n", lightBlue, baseLine+i, addedPrefix, addedBg, white, line, reset)
	}
	// Print bottom context lines
	for i, line := range bottomLines {
		fmt.Printf("%s%3d| %s%s%s\n", lightBlue, startLine+len(topLines)+len(original)+i, gray, line, reset)
	}
}

// generateDiff returns the original and updated lines with context for diff display
func generateDiff(edit EditInstruction, lines []string) (topLines, original, updated, bottomLines []string) {
	if len(lines) == 1 {
		return lines, nil, strings.Split(edit.NewContent, "\n"), nil
	}

	contextLines := 2 // Number of context lines before and after
	start := max(1, edit.LineNumber-contextLines)
	end := min(len(lines), edit.LineNumber+contextLines)

	topLines = lines[start-1 : edit.LineNumber-1]
	if edit.Action == "replace" {
		original = lines[edit.LineNumber-1 : edit.LineNumber]
	} else {
		original = nil
	}
	updated = strings.Split(edit.NewContent, "\n")
	bottomLines = lines[edit.LineNumber-1 : end]

	return topLines, original, updated, bottomLines
}

// min returns the smaller of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// max returns the larger of two integers
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// validateEdits validates the line numbers in the edit requests against the original file lines
func validateEdits(request FileEditRequest, lines []string) error {
	for _, edit := range request.Edits {
		if edit.LineNumber < 1 || edit.LineNumber > len(lines)+1 {
			return fmt.Errorf("invalid line number %d for file %s (file has %d lines)", edit.LineNumber, request.FilePath, len(lines))
		}
	}
	return nil
}

// sortEdits sorts a slice of EditInstruction by line number in ascending order.
// It returns a new slice, leaving the original unchanged.
func sortEdits(edits []EditInstruction) []EditInstruction {
	sortedEdits := make([]EditInstruction, len(edits))
	copy(sortedEdits, edits)
	sort.Slice(sortedEdits, func(i, j int) bool {
		return sortedEdits[i].LineNumber < sortedEdits[j].LineNumber
	})
	return sortedEdits
}

// promptUser asks for user confirmation and returns true if confirmed, false otherwise.
func promptUser(message string) bool {
	fmt.Print(message)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	response := strings.TrimSpace(strings.ToLower(scanner.Text()))
	return response == "y" || response == "yes"
}

// applyEdits applies a slice of sorted edits to a slice of lines and returns the updated lines.
func applyEdits(lines []string, edits []EditInstruction) ([]string, error) {
	updatedLines := make([]string, 0, len(lines)+len(edits))
	editIndex := 0
	lineIndex := 0

	for lineIndex < len(lines) {
		currentLineNumber := lineIndex + 1

		if editIndex < len(edits) && edits[editIndex].LineNumber == currentLineNumber {
			edit := edits[editIndex]
			newLines := strings.Split(edit.NewContent, "\n")

			switch edit.Action {
			case "insert":
				updatedLines = append(updatedLines, newLines...)
				updatedLines = append(updatedLines, lines[lineIndex])
			case "replace":
				updatedLines = append(updatedLines, newLines...)
			default:
				return nil, fmt.Errorf("invalid action %s", edit.Action)
			}

			editIndex++
			lineIndex++
		} else {
			updatedLines = append(updatedLines, lines[lineIndex])
			lineIndex++
		}
	}

	// Handle edits that are for lines beyond the original file length (e.g., append)
	for editIndex < len(edits) {
		edit := edits[editIndex]
		if edit.Action == "insert" {
			newLines := strings.Split(edit.NewContent, "\n")
			updatedLines = append(updatedLines, newLines...)
		}
		editIndex++
	}
	return updatedLines, nil
}

// writeFileLines writes a slice of strings to a file, with each string as a new line.
func writeFileLines(filePath string, lines []string) error {
	newContent := strings.Join(lines, "\n")
	if err := os.WriteFile(filePath, []byte(newContent), 0644); err != nil {
		return fmt.Errorf("failed to write file %s: %v", filePath, err)
	}
	return nil
}

// editFileWorkflow orchestrates the file editing process.
func editFileWorkflow(request FileEditRequest, verbose, prompt, highlight bool) error {
	// Read file
	lines, err := readFileLines(request.FilePath)
	if err != nil {
		return err
	}

	// Validate edits
	if err = validateEdits(request, lines); err != nil {
		return err
	}

	// Sort edits
	edits := sortEdits(request.Edits)

	// Show diffs
	if verbose {
		fmt.Printf("Proposed changes for %s:\n", request.FilePath)
		for _, edit := range edits {
			topLines, original, updated, bottomLines := generateDiff(edit, lines)
			fmt.Printf("\nEdit at line %d (%s):\n", edit.LineNumber, edit.Action)
			printDiff(topLines, original, updated, bottomLines, max(1, edit.LineNumber-2), highlight)
		}
	}

	// Prompt user
	if prompt {
		if !promptUser("\nApply these changes? (y/n): ") {
			return fmt.Errorf("user aborted the file edit operation")
		}
	}

	// Apply edits
	updatedLines, err := applyEdits(lines, edits)
	if err != nil {
		return fmt.Errorf("failed to apply edits to %s: %v", request.FilePath, err)
	}

	// Write file
	if err := writeFileLines(request.FilePath, updatedLines); err != nil {
		return err
	}

	if verbose {
		fmt.Printf("Successfully updated file %s\n", request.FilePath)
	}
	return nil
}
