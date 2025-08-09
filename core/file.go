package core

import (
	"bytes"
	"os"
)

// readFile reads the content of a file at the given path and returns it as a byte slice.
// It returns an error if the file cannot be read.
func readFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

// writeFile writes data to a file at the given path, creating the file if it doesn't exist.
// To ensure data integrity, it performs an atomic write by first writing to a temporary file
// and then renaming it to the final destination.
func writeFile(path string, data []byte) error {
	tempFile, err := os.CreateTemp("", "ffs-")
	if err != nil {
		return err
	}
	defer os.Remove(tempFile.Name())

	if _, err := tempFile.Write(data); err != nil {
		return err
	}

	if err := tempFile.Close(); err != nil {
		return err
	}

	return os.Rename(tempFile.Name(), path)
}

// deleteFile removes the file at the given path.
// It returns an error if the file cannot be removed.
func deleteFile(path string) error {
	return os.Remove(path)
}

// isBinary checks if a file is likely binary by reading its first 1024 bytes
// and checking for the presence of null bytes.
func isBinary(path string) bool {
	file, err := os.Open(path)
	if err != nil {
		return false // Or handle error appropriately
	}
	defer file.Close()

	buffer := make([]byte, 1024)
	n, err := file.Read(buffer)
	if err != nil && err.Error() != "EOF" {
		return false // Or handle error
	}

	return bytes.Contains(buffer[:n], []byte{0})
}
