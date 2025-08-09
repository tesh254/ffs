package ffs

import (
	"os"
	"path/filepath"
	"testing"
)

func TestFFS(t *testing.T) {
	// Create a new temporary directory for testing.
	tmpDir, err := os.MkdirTemp("", "ffs-test")
	if err != nil {
		t.Fatalf("failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create a new ffs instance.
	fs := New()

	// Test creating a new directory.
	dirPath := filepath.Join(tmpDir, "test-dir")
	d := fs.Dir(dirPath)
	if err = d.Create(); err != nil {
		t.Fatalf("failed to create directory: %v", err)
	}

	// Test getting the path of the directory.
	if d.Path() != dirPath {
		t.Errorf("unexpected directory path: got %q, want %q", d.Path(), dirPath)
	}

	// Test creating a new file.
	filePath := filepath.Join(dirPath, "test-file.txt")
	f := fs.File(filePath)

	// Test getting the path of the file.
	if f.Path() != filePath {
		t.Errorf("unexpected file path: got %q, want %q", f.Path(), filePath)
	}

	// Test writing to the file.
	data := []byte("hello, world")
	if err = f.Write(data); err != nil {
		t.Fatalf("failed to write to file: %v", err)
	}

	// Test reading from the file.
	readData, err := f.Read()
	if err != nil {
		t.Fatalf("failed to read from file: %v", err)
	}
	if string(readData) != string(data) {
		t.Errorf("unexpected file content: got %q, want %q", string(readData), string(data))
	}

	// Test building the directory tree.
	tree, err := d.Tree(nil, nil)
	if err != nil {
		t.Fatalf("failed to build directory tree: %v", err)
	}
	if len(tree.Children) != 1 {
		t.Errorf("unexpected number of children in tree: got %d, want 1", len(tree.Children))
	}
	if tree.Children[0].Name != "test-file.txt" {
		t.Errorf("unexpected file in tree: got %q, want %q", tree.Children[0].Name, "test-file.txt")
	}

	// Test deleting the file.
	if err := f.Delete(); err != nil {
		t.Fatalf("failed to delete file: %v", err)
	}

	// Test deleting the directory.
	if err := d.Delete(); err != nil {
		t.Fatalf("failed to delete directory: %v", err)
	}
}
