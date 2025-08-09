package core

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCreateDir_AlreadyExists(t *testing.T) {
	// Create a temporary directory path
	dir, err := os.MkdirTemp("", "test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	// Create a directory using the function
	if err := CreateDir(dir); err != nil {
		t.Fatalf("CreateDir failed: %v", err)
	}

	// Check if the directory was created
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		t.Errorf("CreateDir failed: directory %q was not created", dir)
	}
}

func TestCreateDir(t *testing.T) {
	// Create a temporary directory path
	dir, err := os.MkdirTemp("", "test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir)

	// Create a directory using the function
	newDirPath := dir + "/newdir"
	if err := CreateDir(newDirPath); err != nil {
		t.Fatalf("CreateDir failed: %v", err)
	}

	// Check if the directory was created
	if _, err := os.Stat(newDirPath); os.IsNotExist(err) {
		t.Errorf("CreateDir failed: directory %q was not created", newDirPath)
	}
}

func TestDeleteDir_NonExistent(t *testing.T) {
	err := DeleteDir("non-existent-dir")
	if err != nil {
		t.Errorf("DeleteDir with non-existent directory should not have returned an error, got %v", err)
	}
}

func TestDeleteDir(t *testing.T) {
	// Create a temporary directory
	dir, err := os.MkdirTemp("", "test")
	if err != nil {
		t.Fatal(err)
	}

	// Delete the directory using the function
	if err := DeleteDir(dir); err != nil {
		t.Fatalf("DeleteDir failed: %v", err)
	}

	// Check if the directory still exists
	if _, err := os.Stat(dir); !os.IsNotExist(err) {
		t.Errorf("DeleteDir failed: directory %q still exists", dir)
	}
}

func TestBuildDirectoryTree(t *testing.T) {
	// Create a temporary directory for testing
	tmpDir, err := os.MkdirTemp("", "testTree")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create test files and directories
	//- tmpDir/
	//  - text_file.txt
	//  - binary_file.bin
	//  - emptyDir/
	//  - subDir/
	//    - another_file.txt
	if err = os.WriteFile(tmpDir+"/text_file.txt", []byte("hello world"), 0644); err != nil {
		t.Fatalf("Failed to write file: %v", err)
	}
	if err = os.WriteFile(tmpDir+"/binary_file.bin", []byte{0x00, 0x01, 0x02, 0x03}, 0644); err != nil {
		t.Fatalf("Failed to write file: %v", err)
	}
	if err = os.Mkdir(tmpDir+"/emptyDir", 0755); err != nil {
		t.Fatalf("Failed to create dir: %v", err)
	}
	if err = os.Mkdir(tmpDir+"/subDir", 0755); err != nil {
		t.Fatalf("Failed to create dir: %v", err)
	}
	if err = os.WriteFile(tmpDir+"/subDir/file2.txt", []byte("world"), 0644); err != nil {
		t.Fatalf("Failed to write file: %v", err)
	}

	// Build the directory tree
	tree, err := buildDirectoryTree(tmpDir, nil, nil)
	if err != nil {
		t.Fatalf("buildDirectoryTree failed: %v", err)
	}

	// Assertions
	if tree.Name != filepath.Base(tmpDir) {
		t.Errorf("Expected root name %s, got %s", filepath.Base(tmpDir), tree.Name)
	}
	if len(tree.Children) != 4 {
		t.Errorf("Expected 4 children, got %d", len(tree.Children))
	}

	// Find and check text_file.txt
	var textFile DirectoryTree
	for _, child := range tree.Children {
		if child.Name == "text_file.txt" {
			textFile = child
			break
		}
	}
	if textFile.Name == "" {
		t.Fatal("text_file.txt not found in directory tree")
	}
	if textFile.IsBinary {
		t.Errorf("Expected text_file.txt to not be binary")
	}

	// Find and check binary_file.bin
	var binaryFile DirectoryTree
	for _, child := range tree.Children {
		if child.Name == "binary_file.bin" {
			binaryFile = child
			break
		}
	}
	if binaryFile.Name == "" {
		t.Fatal("binary_file.bin not found in directory tree")
	}
	if !binaryFile.IsBinary {
		t.Errorf("Expected binary_file.bin to be binary")
	}
}
