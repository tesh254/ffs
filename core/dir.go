package core

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type DirectoryTree struct {
	Path     string          `json:"path"`
	Name     string          `json:"name"`
	IsFile   bool            `json:"is_file"`
	IsBinary bool            `json:"is_binary,omitempty"`
	Children []DirectoryTree `json:"children,omitempty"`
	Size     int64           `json:"size"`
}

// createDir creates a directory at the specified path, along with any necessary parents.
// If the directory already exists, createDir does nothing and returns nil.
func createDir(path string) error {
	return os.MkdirAll(path, os.ModePerm)
}

// deleteDir removes a directory at the specified path, along with any children it contains.
// If the path does not exist, deleteDir does nothing and returns nil.
func deleteDir(path string) error {
	return os.RemoveAll(path)
}

// buildDirectoryTree recursively builds a DirectoryTree from a given path.
func buildDirectoryTree(path string, include, exclude []string) (DirectoryTree, error) {
	info, err := os.Stat(path)
	if err != nil {
		return DirectoryTree{}, err
	}

	// Check exclude patterns against the base name.
	for _, pattern := range exclude {
		if matched, _ := filepath.Match(pattern, info.Name()); matched {
			return DirectoryTree{}, nil // Excluded
		}
	}

	if !info.IsDir() {
		// It's a file. Check include patterns.
		if len(include) > 0 {
			included := false
			for _, pattern := range include {
				if matched, _ := filepath.Match(pattern, info.Name()); matched {
					included = true
					break
				}
			}
			if !included {
				return DirectoryTree{}, nil
			}
		}
		// Included file.
		return DirectoryTree{
				Path:     path,
				Name:     info.Name(),
				IsFile:   true,
				IsBinary: IsBinary(path),
				Size:     info.Size(),
			},
			nil
	}

	// It's a directory. Recurse.
	entries, err := os.ReadDir(path)
	if err != nil {
		return DirectoryTree{}, err
	}

	var children []DirectoryTree
	var size int64
	for _, entry := range entries {
		childPath := filepath.Join(path, entry.Name())
		child, err := buildDirectoryTree(childPath, include, exclude)
		if err != nil {
			// Log error and continue
			fmt.Printf("error processing %s: %v\n", childPath, err)
			continue
		}
		if child.Path != "" { // If not excluded/filtered
			children = append(children, child)
			size += child.Size
		}
	}

	// If it's a directory, it's only included if it has children after filtering,
	// unless there are no include patterns (in which case empty dirs are fine).
	if len(children) == 0 && len(include) > 0 {
		return DirectoryTree{}, nil
	}

	return DirectoryTree{
			Path:     path,
			Name:     info.Name(),
			IsFile:   false,
			Children: children,
			Size:     size,
		},
		nil
}

// workingDirectoryTree constructs and prints a directory tree of the current working directory.
func workingDirectoryTree(include, exclude []string) (DirectoryTree, error) {
	dir, err := os.Getwd()
	if err != nil {
		return DirectoryTree{}, fmt.Errorf("could not get working directory: %w", err)
	}

	tree, err := buildDirectoryTree(dir, include, exclude)
	if err != nil {
		return DirectoryTree{}, fmt.Errorf("could not build directory tree for %q: %w", dir, err)
	}

	return tree, nil
}

// printTree prints the directory tree in a formatted manner.
func printTree(tree DirectoryTree) {
	var printRecursive func([]DirectoryTree, string)
	printRecursive = func(children []DirectoryTree, prefix string) {
		for i, child := range children {
			connector := "├── "
			newPrefix := prefix + "│   "
			if i == len(children)-1 {
				connector = "└── "
				newPrefix = prefix + "    "
			}

			fmt.Printf("%s%s%s\n", prefix, connector, child.Name)

			if len(child.Children) > 0 {
				printRecursive(child.Children, newPrefix)
			}
		}
	}
	printRecursive(tree.Children, "")
}

func getTreeMinifiedJSON(tree DirectoryTree) (string, error) {
	jsonBytes, err := json.Marshal(tree)
	if err != nil {
		return "", fmt.Errorf("could not marshal directory tree to JSON: %w", err)
	}
	return string(jsonBytes), nil
}

func printTreeInJSON(tree DirectoryTree) {
	jsonBytes, err := json.MarshalIndent(tree, "", "  ")
	if err != nil {
		fmt.Printf("could not marshal directory tree to JSON: %v", err)
		return
	}
	fmt.Println(string(jsonBytes))
}
