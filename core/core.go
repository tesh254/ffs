package core

// ReadFile reads the content of a file at the given path.
func ReadFile(path string) ([]byte, error) {
	return readFile(path)
}

// GetFileLineMap reads and parses a file to line map.
// The line number being the key and the contents being the value.
func FileToLineMap(path string) (map[string]string, error) {
	content, err := readFile(path)
	if err != nil {
		return nil, err
	}
	return toLineMap(string(content)), nil
}

// LineMapToFile converts a line map back to a file content.
func LineMapToFile(lineMap map[string]string) string {
	return fromLineMap(lineMap)
}

// WriteFile writes data to a file at the given path.
func WriteFile(path string, data []byte) error {
	return writeFile(path, data)
}

// DeleteFile removes the file at the given path.
func DeleteFile(path string) error {
	return deleteFile(path)
}

// CreateDir creates a directory at the specified path.
func CreateDir(path string) error {
	return createDir(path)
}

// DeleteDir removes a directory at the specified path.
func DeleteDir(path string) error {
	return deleteDir(path)
}

// ApplyPatch takes the original content and a JSON string representing the patch,
// then returns the updated content.
func ApplyPatch(originalContent string, patchJSON string) (string, error) {
	return applyPatch(originalContent, patchJSON)
}

// WorkingDirectoryTree returns a tree of the current working directory
func WorkingDirectoryTree(include, exclude []string) (DirectoryTree, error) {
	tree, err := workingDirectoryTree(include, exclude)
	return tree, err
}

// PrintDirectoryTree prints the directory tree in a human-readable format.
func PrintDirectoryTree(tree DirectoryTree, inJson bool) {
	if !inJson {
		printTree(tree)
	} else {
		printTreeInJSON(tree)
	}
}

// GetTreeMinifiedJSON returns the directory tree as a minified JSON string.
func GetTreeMinifiedJSON(tree DirectoryTree) (string, error) {
	return getTreeMinifiedJSON(tree)
}

// BuildDirTree builds a tree based on path provided
func BuildDirTree(path string, include, exclude []string) (DirectoryTree, error) {
	tree, err := buildDirectoryTree(path, include, exclude)
	return tree, err
}
