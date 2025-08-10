package main

import (
	"fmt"
	"log"

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

	editRequestReplacing := core.FileEditRequest{
		FilePath: replacingFilePath,
		Edits: []core.EditInstruction{
			{
				Action:     "replace",
				LineNumber: 2,
				NewContent: "this is a replaced line",
			},
		},
	}

	if err := core.ApplyPatch(editRequestReplacing, true, true, true); err != nil {
		log.Fatalf("Failed to apply suggestion: %v", err)
	}

	newContentReplacing, err := core.ReadFile(replacingFilePath)
	if err != nil {
		log.Fatalf("Failed to read file after applying patch: %v", err)
	}

	fmt.Println("New file content after applying suggestion:")
	fmt.Println(string(newContentReplacing))
	fmt.Println("--------------------")

	// --- Adding Patch Example ---
	fmt.Println("--- Adding Patch Example ---")
	addingFilePath := "adding_example.txt"
	originalContentAdding := "line1\nline3"
	if err = core.WriteFile(addingFilePath, []byte(originalContentAdding)); err != nil {
		log.Fatalf("Failed to create initial file: %v", err)
	}
	defer core.DeleteFile(addingFilePath)

	fmt.Println("Original file content:")
	fmt.Println(originalContentAdding)
	fmt.Println("--------------------")

	editRequestAdding := core.FileEditRequest{
		FilePath: addingFilePath,
		Edits: []core.EditInstruction{
			{
				Action:     "insert",
				LineNumber: 2,
				NewContent: "this is an added line",
			},
		},
	}

	if err = core.ApplyPatch(editRequestAdding, true, true, true); err != nil {
		log.Fatalf("Failed to apply suggestion: %v", err)
	}

	newContentAdding, err := core.ReadFile(addingFilePath)
	if err != nil {
		log.Fatalf("Failed to read file after applying patch: %v", err)
	}

	fmt.Println("New file content after applying suggestion:")
	fmt.Println(string(newContentAdding))
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
