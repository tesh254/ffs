package core

import (
	"io/ioutil"
	"os"
)

// ReadFile reads the content of a file at the given path and returns it as a byte slice.
// It returns an error if the file cannot be read.
func ReadFile(path string) ([]byte, error) {
	return ioutil.ReadFile(path)
}

// WriteFile writes data to a file at the given path, creating the file if it doesn't exist.
// To ensure data integrity, it performs an atomic write by first writing to a temporary file
// and then renaming it to the final destination.
func WriteFile(path string, data []byte) error {
	tempFile, err := ioutil.TempFile("", "ffs-")
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

// DeleteFile removes the file at the given path.
// It returns an error if the file cannot be removed.
func DeleteFile(path string) error {
	return os.Remove(path)
}
