package ffs

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/tesh254/ffs/core"
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

func TestApplyPatch(t *testing.T) {
	// Create a temporary file with some content
	tmpfile, err := os.CreateTemp("", "test.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	content := "line1\nline2\nline3"
	if _, err = tmpfile.Write([]byte(content)); err != nil {
		t.Fatal(err)
	}
	tmpfile.Close()

	// Test replacing a line
	replaceRequest := core.FileEditRequest{
		FilePath: tmpfile.Name(),
		Edits: []core.EditInstruction{
			{
				Action:     "replace",
				LineNumber: 2,
				NewContent: "this is a replaced line",
			},
		},
	}

	err = core.ApplyPatch(replaceRequest, true, false, false)
	if err != nil {
		t.Errorf("ApplyPatch failed: %v", err)
	}

	// Read the file and check the content
	newContent, err := os.ReadFile(tmpfile.Name())
	if err != nil {
		t.Fatal(err)
	}
	expectedAfterReplace := "line1\nthis is a replaced line\nline3"
	if string(newContent) != expectedAfterReplace {
		t.Errorf("ApplyPatch did not replace the line correctly. Got %q, expected %q", string(newContent), expectedAfterReplace)
	}

	// Test adding a line
	addRequest := core.FileEditRequest{
		FilePath: tmpfile.Name(),
		Edits: []core.EditInstruction{
			{
				Action:     "insert",
				LineNumber: 2,
				NewContent: "this is an added line",
			},
		},
	}
	err = core.ApplyPatch(addRequest, true, false, false)
	if err != nil {
		t.Errorf("ApplyPatch failed: %v", err)
	}

	// Read the file and check the content
	newContent, err = os.ReadFile(tmpfile.Name())
	if err != nil {
		t.Fatal(err)
	}
	expectedAfterAdd := "line1\nthis is an added line\nthis is a replaced line\nline3"
	if string(newContent) != expectedAfterAdd {
		t.Errorf("ApplyPatch did not add the line correctly. Got %q, expected %q", string(newContent), expectedAfterAdd)
	}
}
