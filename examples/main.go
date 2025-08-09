package main

import (
	"fmt"
	"log"

	"github.com/tesh254/ffs/agent"
	"github.com/tesh254/ffs/core"
)

func main() {
	// --- Replacing Patch Example ---
	fmt.Println("--- Replacing Patch Example ---")
	replacingFilePath := "replacing_example.txt"
	originalContentReplacing := "line1\nline2\nline3"
	if err := core.WriteFile(replacingFilePath, []byte(originalContentReplacing)); err != nil {
		log.Fatalf("Failed to create initial file: %v", err)
	}
	defer core.DeleteFile(replacingFilePath)

	fmt.Println("Original file content:")
	fmt.Println(originalContentReplacing)
	fmt.Println("--------------------")

	suggestionReplacing := agent.Suggestion{
		FilePath:    replacingFilePath,
		LineChanges: `{"2": "this is a replaced line"}`,
		PatchType:   core.PatchTypeReplacing,
	}

	newContentReplacing, err := agent.ApplySuggestion(suggestionReplacing)
	if err != nil {
		log.Fatalf("Failed to apply suggestion: %v", err)
	}

	fmt.Println("New file content after applying suggestion:")
	fmt.Println(newContentReplacing)
	fmt.Println("--------------------")

	fmt.Println("Diff:")
	core.PrintDiff(originalContentReplacing, newContentReplacing)
	fmt.Println("--------------------")

	// --- Adding Patch Example ---
	fmt.Println("--- Adding Patch Example ---")
	addingFilePath := "adding_example.txt"
	originalContentAdding := "line1\nline3"
	if err := core.WriteFile(addingFilePath, []byte(originalContentAdding)); err != nil {
		log.Fatalf("Failed to create initial file: %v", err)
	}
	defer core.DeleteFile(addingFilePath)

	fmt.Println("Original file content:")
	fmt.Println(originalContentAdding)
	fmt.Println("--------------------")

	suggestionAdding := agent.Suggestion{
		FilePath:    addingFilePath,
		LineChanges: `{"2": "this is an added line"}`,
		PatchType:   core.PatchTypeAdding,
	}

	newContentAdding, err := agent.ApplySuggestion(suggestionAdding)
	if err != nil {
		log.Fatalf("Failed to apply suggestion: %v", err)
	}

	fmt.Println("New file content after applying suggestion:")
	fmt.Println(newContentAdding)
	fmt.Println("--------------------")

	fmt.Println("Diff:")
	core.PrintDiff(originalContentAdding, newContentAdding)
	fmt.Println("--------------------")

	// --- Tree Example ---
	fmt.Println("--- Directory Example ---")
	tree, err := core.WorkingDirectoryTree(nil, []string{".git", "node_modules", ".DS_Store"})
	if err != nil {
		log.Fatalf("Failed to get working directory tree: %v", err)
	}
	fmt.Println("Text Print Example:")
	core.PrintDirectoryTree(tree, false)
	fmt.Println("--------------------")
	fmt.Println("JSON Print Example:")
	core.PrintDirectoryTree(tree, true)
}
