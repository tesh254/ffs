# FFS (For File's Sake) Documentation

Welcome to the documentation for the `ffs` library. This guide provides a comprehensive overview of the library's features and how to use them effectively.

## Table of Contents

- [Getting Started](#getting-started)
  - [Installation](#installation)
  - [Basic Usage](#basic-usage)
- [High-Level API: The `ffs` Package](#high-level-api-the-ffs-package)
  - [Creating an `ffs` Instance](#creating-an-ffs-instance)
  - [Working with Files](#working-with-files)
    - [Reading a File](#reading-a-file)
    - [Writing a File](#writing-a-file)
    - [Deleting a File](#deleting-a-file)
  - [Working with Directories](#working-with-directories)
    - [Creating a Directory](#creating-a-directory)
    - [Deleting a Directory](#deleting-a-directory)
    - [Getting a Directory Tree](#getting-a-directory-tree)
- [Low-Level API: The `core` Package](#low-level-api-the-core-package)
  - [File Operations](#file-operations)
    - [ReadFile](#readfile)
    - [WriteFile](#writefile)
    - [DeleteFile](#deletefile)
  - [Directory Operations](#directory-operations)
    - [CreateDir](#createdir)
    - [DeleteDir](#deletedir)
  - [Patching](#patching)
    - [ApplyPatch](#applypatch)
  - [Directory Trees](#directory-trees)
    - [WorkingDirectoryTree](#workingdirectorytree)
    - [PrintDirectoryTree](#printdirectorytree)
    - [GetTreeMinifiedJSON](#gettreeminifiedjson)
    - [BuildDirTree](#builddirtree)
- [LLM Agent Integration](#llm-agent-integration)
  - [Applying a Suggestion](#applying-a-suggestion)

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

**Getting a Directory Tree:**

```go
// Get the directory tree including only .go files
tree, err := dir.Tree([]string{"*.go"}, nil)
if err != nil {
    // Handle error
}

// Print in human-readable format
core.PrintDirectoryTree(tree, false)
```

## Low-Level API: The `core` Package

The `core` package provides a set of low-level functions for interacting with the filesystem. These functions are used by the high-level `ffs` package, but they can also be used directly when you need more control.

### File Operations

#### ReadFile

The `ReadFile` function reads the entire content of a file at the specified path.

```go
import "github.com/tesh254/ffs/core"

content, err := core.ReadFile("path/to/your/file.txt")
if err != nil {
    // Handle error
}
fmt.Println(string(content))
```

#### WriteFile

The `WriteFile` function writes a slice of bytes to a file at the specified path. If the file doesn't exist, it will be created. If it exists, its contents will be overwritten.

```go
import "github.com/tesh254/ffs/core"

err := core.WriteFile("path/to/your/file.txt", []byte("Hello, world!"))
if err != nil {
    // Handle error
}
```

#### DeleteFile

The `DeleteFile` function removes a file from the filesystem.

```go
import "github.com/tesh254/ffs/core"

err := core.DeleteFile("path/to/your/file.txt")
if err != nil {
    // Handle error
}
```

### Directory Operations

#### CreateDir

The `CreateDir` function creates a new directory at the specified path.

```go
import "github.com/tesh254/ffs/core"

err := core.CreateDir("path/to/your/directory")
if err != nil {
    // Handle error
}
```

#### DeleteDir

The `DeleteDir` function removes a directory from the filesystem. The directory must be empty for it to be deleted.

```go
import "github.com/tesh254/ffs/core"

err := core.DeleteDir("path/to/your/directory")
if err != nil {
    // Handle error
}
```

### Patching

#### ApplyPatch

The `ApplyPatch` function allows you to apply a set of line-based changes to a file's content. The patch is provided as a `FileEditRequest` object. This is particularly useful for applying changes suggested by LLMs.

```go
import "github.com/tesh254/ffs/core"

// Create a file to modify
core.WriteFile("example.txt", []byte("line 1\nline 2\nline 3"))

// An agent suggests a change to the second line
request := core.FileEditRequest{
    Path: "example.txt",
    Edits: []core.FileEdit{
        {
            Line:    2,
            Content: "new line 2",
        },
    },
}

// Apply the suggestion
err := core.ApplyPatch(request, false, false, false)
if err != nil {
    // Handle error
}
```

### Directory Trees

#### WorkingDirectoryTree

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

#### PrintDirectoryTree

The `PrintDirectoryTree` function prints a `DirectoryTree` in a human-readable format. You can also print the tree in JSON format.

```go
import "github.com/tesh254/ffs/core"

tree, _ := core.WorkingDirectoryTree(nil, nil)

// Print in human-readable format
core.PrintDirectoryTree(tree, false)

// Print in JSON format
core.PrintDirectoryTree(tree, true)
```

#### GetTreeMinifiedJSON

The `GetTreeMinifiedJSON` function returns a minified JSON string representation of a `DirectoryTree`. This is useful for sending the directory structure to an LLM.

```go
import "github.com/tesh254/ffs/core"

tree, _ := core.WorkingDirectoryTree(nil, nil)
jsonString, err := core.GetTreeMinifiedJSON(tree)
if err != nil {
    // Handle error
}
fmt.Println(jsonString)
```

#### BuildDirTree

The `BuildDirTree` function returns a `DirectoryTree` struct that represents the directory structure of the given path. You can provide optional `include` and `exclude` patterns to filter the results.

```go
import "github.com/tesh254/ffs/core"

// Get the full directory tree
tree, err := core.BuildDirTree("path/to/dir", nil, nil)
if err != nil {
    // Handle error
}
```
