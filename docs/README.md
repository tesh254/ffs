# FFS (For File's Sake) Documentation

Welcome to the documentation for the `ffs` library. This guide provides a comprehensive overview of the library's features and how to use them effectively.

## Getting Started

### Installation

To use `ffs` in your Go project, you can install it using `go get`:

```bash
go get github.com/tesh254/ffs
```

### Basic Usage

The `ffs` package provides a high-level, fluent API for file and directory operations. Here's a quick example of how to write to a file:

```go
package main

import (
	"fmt"
	"log"

	"github.com/tesh254/ffs"
)

func main() {
	fs := ffs.New()
	file := fs.File("example.txt")

	err := file.Write([]byte("Hello from ffs!"))
	if err != nil {
		log.Fatalf("Failed to write to file: %v", err)
	}

	fmt.Println("Successfully wrote to example.txt")
}
```

## High-Level API: The `ffs` Package

The `ffs` package is the recommended way to interact with the filesystem. It provides `File` and `Dir` objects that make I/O operations intuitive and easy to manage.

### Creating an `ffs` Instance

First, create a new `ffs` instance:

```go
import "github.com/tesh254/ffs"

fs := ffs.New()
```

### Working with Files

Use the `File()` method to get a `File` object.

```go
file := fs.File("path/to/your/file.txt")
```

**Reading a File:**

```go
content, err := file.Read()
if err != nil {
    // Handle error
}
fmt.Println(string(content))
```

**Writing a File:**

```go
err := file.Write([]byte("Hello, world!"))
if err != nil {
    // Handle error
}
```

## Low-Level API: The `core` Package

The `core` package provides a set of low-level functions for interacting with the filesystem. These functions are used by the high-level `ffs` package, but they can also be used directly when you need more control.

### Reading a File

```go
import "github.com/tesh254/ffs/core"

content, err := core.ReadFile("path/to/your/file.txt")
if err != nil {
    // Handle error
}
fmt.Println(string(content))
```


### Deleting a File

```go
import "github.com/tesh254/ffs/core"

err := core.DeleteFile("path/to/your/file.txt")
if err != nil {
    // Handle error
}
```

### Creating a Directory

```go
import "github.com/tesh254/ffs/core"

err := core.CreateDir("path/to/your/directory")
if err != nil {
    // Handle error
}
```

### Deleting a Directory

```go
import "github.com/tesh254/ffs/core"

err := core.DeleteDir("path/to/your/directory")
if err != nil {
    // Handle error
}
```

### Applying a Patch

The `ApplyPatch` function allows you to apply a set of line-based changes to a string. The patch is provided as a JSON string that maps line numbers to their new content.

```go
import "github.com/tesh254/ffs/core"

originalContent := "line 1\nline 2\nline 3"
patchJSON := `{"2": "new line 2"}`
newContent, err := core.ApplyPatch(originalContent, patchJSON)
if err != nil {
    // Handle error
}
fmt.Println(newContent)
// Output:
// line 1
// new line 2
// line 3
```

### Working with Directory Trees

The `core` package provides functions for generating and displaying directory trees.

**Generating a Directory Tree:**

The `WorkingDirectoryTree` function returns a `DirectoryTree` struct that represents the directory structure of the current working directory. You can provide optional `include` and `exclude` patterns to filter the results.

```go
import "github.com/tesh254/ffs/core"

// Get the full directory tree
tree, err := core.WorkingDirectoryTree(nil, nil)
if err != nil {
    // Handle error
}

// Get the directory tree including only .go files
goFilesTree, err := core.WorkingDirectoryTree([]string{"*.go"}, nil)
if err != nil {
    // Handle error
}
```

**Printing a Directory Tree:**

The `PrintDirectoryTree` function prints a `DirectoryTree` in a human-readable format. You can also print the tree in JSON format.

```go
import "github.com/tesh254/ffs/core"

tree, _ := core.WorkingDirectoryTree(nil, nil)

// Print in human-readable format
core.PrintDirectoryTree(tree, false)

// Print in JSON format
core.PrintDirectoryTree(tree, true)
```

**Getting a Minified JSON Representation:**

The `GetTreeMinifiedJSON` function returns a minified JSON string representation of a `DirectoryTree`.

```go
import "github.com/tesh254/ffs/core"

tree, _ := core.WorkingDirectoryTree(nil, nil)
jsonString, err := core.GetTreeMinifiedJSON(tree)
if err != nil {
    // Handle error
}
fmt.Println(jsonString)
```

**Deleting a File:**

```go
err := file.Delete()
if err != nil {
    // Handle error
}
```

### Working with Directories

Use the `Dir()` method to get a `Dir` object.

```go
dir := fs.Dir("path/to/your/directory")
```

**Creating a Directory:**

```go
err := dir.Create()
if err != nil {
    // Handle error
}
```

**Deleting a Directory:**

```go
err := dir.Delete()
if err != nil {
    // Handle error
}
```

## LLM Agent Integration

The `agent` package is designed to help AI agents apply changes to files in a precise and efficient manner. It's particularly useful for applying line-specific modifications suggested by LLMs.

### Applying a Suggestion

The `agent.ApplySuggestion` function takes a `Suggestion` object, which specifies the file to be changed and a JSON map of line numbers to their new content.

Here's an example of how to apply a suggestion to update the second line of a file:

```go
package main

import (
	"fmt"
	"log"

	"github.com/tesh254/ffs/agent"
	"github.com/tesh254/ffs/core"
)

func main() {
	// Setup: Create a file to modify
	filePath := "example.txt"
	originalContent := "Line 1: This is the original first line.\nLine 2: This is the original second line."
	if err := core.WriteFile(filePath, []byte(originalContent)); err != nil {
		log.Fatalf("Failed to create initial file: %v", err)
	}
	defer core.DeleteFile(filePath)

	// An agent suggests a change to the second line
	suggestion := agent.Suggestion{
		FilePath:    filePath,
		LineChanges: `{"2": "Line 2: This line has been updated by an agent."}`,
	}

	// Apply the suggestion
	newContent, err := agent.ApplySuggestion(suggestion)
	if err != nil {
		log.Fatalf("Failed to apply suggestion: %v", err)
	}

	fmt.Println("File content after applying suggestion:")
	fmt.Println(newContent)
}
```

### Writing a File

```go
import "github.com/tesh254/ffs/core"

err := core.WriteFile("path/to/your/file.txt", []byte("Hello, world!"))
if err != nil {
    // Handle error
}
```
