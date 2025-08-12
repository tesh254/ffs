package core

import (
	"os"
	"testing"
)

func TestReadFile_NonExistent(t *testing.T) {
	_, err := ReadFile("non-existent-file")
	if err == nil {
		t.Error("ReadFile with non-existent file should have returned an error")
	}
}

func TestReadFile(t *testing.T) {
	// Create a temporary file with some content
	content := []byte("hello world")
	tmpfile, err := os.CreateTemp("", "test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name()) // clean up

	if _, tmpFileError := tmpfile.Write(content); tmpFileError != nil {
		t.Fatal(tmpFileError)
	}
	if err = tmpfile.Close(); err != nil {
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

func TestWriteFile_DirNonExistent(t *testing.T) {
	err := WriteFile("non-existent-dir/file", []byte("hello"))
	if err == nil {
		t.Error("WriteFile to a non-existent directory should have returned an error")
	}
}

func TestWriteFile(t *testing.T) {
	// Create a temporary file path
	tmpfile, err := os.CreateTemp("", "test")
	if err != nil {
		t.Fatal(err)
	}
	path := tmpfile.Name()
	tmpfile.Close()
	os.Remove(path) // Remove the temp file created by TempFile

	// Write to the file using the function
	content := []byte("hello again")
	if err = WriteFile(path, content); err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}
	defer os.Remove(path) // clean up

	// Read the file to check if the content is correct
	readContent, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("Failed to read back file: %v", err)
	}

	if string(readContent) != string(content) {
		t.Errorf("WriteFile content mismatch: got %q, want %q", readContent, content)
	}

	// Test error case
	err = WriteFile("/non-existent-dir/file", []byte("hello"))
	if err == nil {
		t.Error("WriteFile to a non-existent directory should have returned an error")
	}
}

func TestDeleteFile_NonExistent(t *testing.T) {
	err := DeleteFile("non-existent-file")
	if err == nil {
		t.Error("DeleteFile with non-existent file should have returned an error")
	}
}

func TestDeleteFile(t *testing.T) {
	// Create a temporary file
	tmpfile, err := os.CreateTemp("", "test")
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

func TestIsBinary(t *testing.T) {
	// Create a temporary text file
	textFile, err := os.CreateTemp("", "testText")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(textFile.Name())
	if _, err = textFile.WriteString("this is a text file"); err != nil {
		t.Fatal(err)
	}
	textFile.Close()

	if IsBinary(textFile.Name()) {
		t.Errorf("IsBinary failed: expected %s to be a text file", textFile.Name())
	}

	// Create a temporary binary file
	binaryFile, err := os.CreateTemp("", "testBinary")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(binaryFile.Name())
	if _, err = binaryFile.Write([]byte{0, 1, 2, 3}); err != nil {
		t.Fatal(err)
	}
	binaryFile.Close()

	if !IsBinary(binaryFile.Name()) {
		t.Errorf("IsBinary failed: expected %s to be a binary file", binaryFile.Name())
	}

	// Test non-existent file
	if IsBinary("non-existent-file") {
		t.Error("IsBinary with non-existent file should have returned false")
	}

	// Test read error
	file, err := os.CreateTemp("", "test")
	if err != nil {
		t.Fatal(err)
	}
	name := file.Name()
	file.Close()
	// Re-opening the file and closing it will not cause a read error.
	// Instead, we make the file unreadable.
	if err := os.Chmod(name, 0200); err != nil {
		t.Fatal(err)
	}
	if IsBinary(name) {
		t.Error("IsBinary with read error should have returned false")
	}
	// Cleanup
	if err := os.Chmod(name, 0600); err != nil {
		t.Fatal(err)
	}
	os.Remove(name)
}
