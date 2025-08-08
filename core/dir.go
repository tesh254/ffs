package core

import (
	"os"
)

// CreateDir creates a directory at the specified path, along with any necessary parents.
// If the directory already exists, CreateDir does nothing and returns nil.
func CreateDir(path string) error {
	return os.MkdirAll(path, os.ModePerm)
}

// DeleteDir removes a directory at the specified path, along with any children it contains.
// If the path does not exist, DeleteDir does nothing and returns nil.
func DeleteDir(path string) error {
	return os.RemoveAll(path)
}
