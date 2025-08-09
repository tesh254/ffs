package core

import (
	"os"
)

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
