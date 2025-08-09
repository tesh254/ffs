package core

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type errorMarshaler struct{}

func (e errorMarshaler) MarshalJSON() ([]byte, error) {
	return nil, os.ErrInvalid
}

func TestWorkingDirectoryTree(t *testing.T) {
	// Create a temporary directory and change into it
	tmpDir, err := os.MkdirTemp("", "testWdTree")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current working directory: %v", err)
	}
	if err = os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}
	defer os.Chdir(originalWd)

	// Create test files and directories
	if err = os.WriteFile("file1.txt", []byte("file1"), 0644); err != nil {
		t.Fatalf("Failed to write file: %v", err)
	}
	if err = os.Mkdir("subdir", 0755); err != nil {
		t.Fatalf("Failed to create dir: %v", err)
	}
	if err = os.WriteFile("subdir/file2.txt", []byte("file2"), 0644); err != nil {
		t.Fatalf("Failed to write file: %v", err)
	}

	// Test with no include/exclude
	tree, err := WorkingDirectoryTree(nil, nil)
	if err != nil {
		t.Fatalf("WorkingDirectoryTree failed: %v", err)
	}
	if tree.Name != filepath.Base(tmpDir) {
		t.Errorf("Expected root name %s, got %s", filepath.Base(tmpDir), tree.Name)
	}
	if len(tree.Children) != 2 {
		t.Errorf("Expected 2 children, got %d", len(tree.Children))
	}

	// Test with include
	tree, err = WorkingDirectoryTree([]string{"*.txt"}, nil)
	if err != nil {
		t.Fatalf("WorkingDirectoryTree with include failed: %v", err)
	}
	if len(tree.Children) != 2 {
		t.Errorf("Expected 2 children with include, got %d", len(tree.Children))
	}

	// Test with exclude
	tree, err = WorkingDirectoryTree(nil, []string{"subdir"})
	if err != nil {
		t.Fatalf("WorkingDirectoryTree with exclude failed: %v", err)
	}
	if len(tree.Children) != 1 {
		t.Errorf("Expected 1 child with exclude, got %d", len(tree.Children))
	}
}

func TestPrintDirectoryTree(t *testing.T) {
	tree := DirectoryTree{
		Name: "root",
		Children: []DirectoryTree{
			{Name: "file1.txt", IsFile: true},
			{Name: "subdir", Children: []DirectoryTree{
				{Name: "file2.txt", IsFile: true},
			}},
		},
	}

	// Capture stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	PrintDirectoryTree(tree, false)

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	expected := "├── file1.txt\n└── subdir\n    └── file2.txt\n"
	if !strings.Contains(output, expected) {
		t.Errorf("PrintDirectoryTree output mismatch:\ngot\n%q\nwant\n%q", output, expected)
	}
}

func TestPrintDirectoryTree_JSON(t *testing.T) {
	tree := DirectoryTree{
		Path: "root",
		Name: "root",
		Children: []DirectoryTree{
			{Path: "root/file1.txt", Name: "file1.txt", IsFile: true},
		},
	}

	// Capture stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	PrintDirectoryTree(tree, true)

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	if !strings.Contains(output, `"path": "root/file1.txt"`) {
		t.Errorf("PrintDirectoryTree JSON output mismatch: got %q", output)
	}

	// Test marshalling of invalid UTF-8
	tree = DirectoryTree{
		Name: string([]byte{0xff}),
	}
	r, w, _ = os.Pipe()
	os.Stdout = w

	PrintDirectoryTree(tree, true)

	w.Close()
	os.Stdout = old

	buf.Reset()
	io.Copy(&buf, r)
	output = buf.String()

	if !strings.Contains(output, `\ufffd`) {
		t.Errorf("Expected invalid UTF-8 to be escaped in JSON output, got %s", output)
	}

}

func TestGetTreeMinifiedJSON(t *testing.T) {
	tree := DirectoryTree{
		Path:   "root",
		Name:   "root",
		IsFile: false,
		Children: []DirectoryTree{
			{
				Path:   "root/file.txt",
				Name:   "file.txt",
				IsFile: true,
				Size:   123,
			},
		},
		Size: 123,
	}

	jsonString, err := GetTreeMinifiedJSON(tree)
	if err != nil {
		t.Fatalf("GetTreeMinifiedJSON failed: %v", err)
	}

	expected := `{"path":"root","name":"root","is_file":false,"children":[{"path":"root/file.txt","name":"file.txt","is_file":true,"size":123}],"size":123}`
	if jsonString != expected {
		t.Errorf("GetTreeMinifiedJSON mismatch:\ngot\n%s\nwant\n%s", jsonString, expected)
	}

	// Test marshalling of invalid UTF-8
	tree = DirectoryTree{
		Path: "root",
		Name: string([]byte{0xff}),
	}
	jsonString, err = GetTreeMinifiedJSON(tree)
	if err != nil {
		t.Fatalf("GetTreeMinifiedJSON failed for invalid UTF-8: %v", err)
	}
	if !strings.Contains(jsonString, `\ufffd`) {
		t.Errorf("Expected invalid UTF-8 to be escaped, got %s", jsonString)
	}

}

func TestWorkingDirectoryTree_Error(t *testing.T) {
	originalWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("Failed to get current working directory: %v", err)
	}
	defer os.Chdir(originalWd)

	tmpDir, err := os.MkdirTemp("", "testWdTreeError")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	if err = os.Chdir(tmpDir); err != nil {
		t.Fatalf("Failed to change directory: %v", err)
	}

	if err = os.Chmod(tmpDir, 0000); err != nil {
		t.Fatalf("Failed to chmod temp dir: %v", err)
	}

	_, err = WorkingDirectoryTree(nil, nil)
	if err == nil {
		t.Error("WorkingDirectoryTree should have failed with an unreadable directory")
	}

	if err := os.Chmod(tmpDir, 0755); err != nil {
		t.Fatalf("Failed to chmod temp dir back: %v", err)
	}
}
