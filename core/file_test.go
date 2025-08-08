package core

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestReadFile(t *testing.T) {
	// Create a temporary file with some content
	content := []byte("hello world")
	tmpfile, err := ioutil.TempFile("", "test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name()) // clean up

	if _, err := tmpfile.Write(content); err != nil {
		t.Fatal(err)
	}
	if err := tmpfile.Close(); err != nil {
		t.Fatal(err)
	}

	// Read the file using the function
	readContent, err := ReadFile(tmpfile.Name())
	if err != nil {
		t.Fatalf("ReadFile failed: %v", err)
	}

	// Check if the content is the same
	if string(readContent) != string(content) {
		t.Errorf("ReadFile content mismatch: got %q, want %q", readContent, content)
	}
}

func TestWriteFile(t *testing.T) {
	// Create a temporary file path
	tmpfile, err := ioutil.TempFile("", "test")
	if err != nil {
		t.Fatal(err)
	}
	path := tmpfile.Name()
	tmpfile.Close()
	os.Remove(path) // Remove the temp file created by TempFile

	// Write to the file using the function
	content := []byte("hello again")
	if err := WriteFile(path, content); err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}
	defer os.Remove(path) // clean up

	// Read the file to check if the content is correct
	readContent, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatalf("Failed to read back file: %v", err)
	}

	if string(readContent) != string(content) {
		t.Errorf("WriteFile content mismatch: got %q, want %q", readContent, content)
	}
}

func TestDeleteFile(t *testing.T) {
	// Create a temporary file
	tmpfile, err := ioutil.TempFile("", "test")
	if err != nil {
		t.Fatal(err)
	}
	path := tmpfile.Name()
	tmpfile.Close()

	// Delete the file using the function
	if err := DeleteFile(path); err != nil {
		t.Fatalf("DeleteFile failed: %v", err)
	}

	// Check if the file still exists
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Errorf("DeleteFile failed: file %q still exists", path)
	}
}
