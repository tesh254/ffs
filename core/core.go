package core

// ReadFile reads the content of a file at the given path.
func ReadFile(path string) ([]byte, error) {
	return readFile(path)
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

// ApplyPatch applies a patch to a file.
func ApplyPatch(request FileEditRequest, verbose, prompt, highlight bool) error {
	return editFileWorkflow(request, verbose, prompt, highlight)
}

// SearchFiles performs a concurrent search for a query in a given path.
func SearchFiles(rootPath, query string, options SearchOptions) ([]SearchResult, error) {
	return search(rootPath, query, options)
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

// ReadFileLines reads the lines of a file at the given path.
func ReadFileLines(path string) ([]string, error) {
	return readFileLines(path)
}

// ValidateEdits validates the edits against the file content.
func ValidateEdits(request FileEditRequest, lines []string) error {
	return validateEdits(request, lines)
}

// SortEdits sorts the edits by line number.
func SortEdits(edits []EditInstruction) []EditInstruction {
	return sortEdits(edits)
}

// GenerateDiff generates the diff for the edit.
func GenerateDiff(edit EditInstruction, lines []string) (topLines, original, updated, bottomLines []string) {
	return generateDiff(edit, lines)
}

// PrintDiff prints the diff for the edit.
func PrintDiff(topLines, original, updated, bottomLines []string, edit EditInstruction) {
	printDiff(topLines, original, updated, bottomLines, max(1, edit.LineNumber), true)
}

// PromptUser prompts the user for confirmation.
func PromptUser(question string) bool {
	return promptUser(question + " (y/n): ")
}

// ApplyEdits applies the edits to the file content.
func ApplyEdits(lines []string, edits []EditInstruction) ([]string, error) {
	return applyEdits(lines, edits)
}

// WriteFileLines writes the lines to a file at the given path.
func WriteFileLines(path string, lines []string) error {
	return writeFileLines(path, lines)
}
