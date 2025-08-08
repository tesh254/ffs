package core

import (
	"io/ioutil"
	"os"
	"testing"
)

func TestCreateDir(t *testing.T) {
	// Create a temporary directory path
	dir, err := ioutil.TempDir("", "test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(dir) // clean up

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

func TestDeleteDir(t *testing.T) {
	// Create a temporary directory
	dir, err := ioutil.TempDir("", "test")
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
