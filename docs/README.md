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
