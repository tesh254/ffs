package main

import (
	"fmt"
	"log"

	"github.com/tesh254/ffs/agent"
	"github.com/tesh254/ffs/core"
)

func main() {
	// Create a temporary file for the example
	filePath := "example.txt"
	originalContent := "This is a sample file for the ffs library.\nThis should be the second line for the ffs sample file"
	if err := core.WriteFile(filePath, []byte(originalContent)); err != nil {
		log.Fatalf("Failed to create initial file: %v", err)
	}
	defer core.DeleteFile(filePath)

	fmt.Println("Original file content:")
	fmt.Println(originalContent)
	fmt.Println("--------------------")

	// Simulate an LLM agent suggestion with line-specific changes
	suggestion := agent.Suggestion{
		FilePath:    filePath,
		LineChanges: `{"2": "This is the updated content of the sample file, modified by the LLM agent."}`,
	}

	// Apply the suggestion
	newContent, err := agent.ApplySuggestion(suggestion)
	if err != nil {
		log.Fatalf("Failed to apply suggestion: %v", err)
	}

	fmt.Println("New file content after applying suggestion:")
	fmt.Println(newContent)
	fmt.Println("--------------------")

	// Verify the file content on disk
	finalContent, err := core.ReadFile(filePath)
	if err != nil {
		log.Fatalf("Failed to read final file content: %v", err)
	}

	fmt.Println("Final content read from disk:")
	fmt.Println(string(finalContent))
	fmt.Println("--------------------")

	if string(finalContent) == newContent {
		fmt.Println("Successfully applied suggestion to the file.")
	} else {
		fmt.Println("Failed to apply suggestion to the file.")
	}
}
