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

// ApplyPatch takes the original content and a JSON string representing the patch,
// then returns the updated content.
func ApplyPatch(originalContent string, patchJSON string) (string, error) {
	return applyPatch(originalContent, patchJSON)
}
