package ffs

import "github.com/tesh254/ffs/core"

// ffs is the default implementation of the FileSystem interface.
// It uses the core package to interact with the underlying file system.
type ffs struct{}

// New returns a new instance of the default FileSystem implementation.
func New() FileSystem {
	return &ffs{}
}

// File returns a new File instance for the given path.
func (f *ffs) File(path string) File {
	return &file{path: path}
}

// Dir returns a new Dir instance for the given path.
func (f *ffs) Dir(path string) Dir {
	return &dir{path: path}
}

// file is the default implementation of the File interface.
type file struct {
	path string
}

// Read reads the content of the file.
func (f *file) Read() ([]byte, error) {
	return core.ReadFile(f.path)
}

// Write writes data to the file.
func (f *file) Write(data []byte) error {
	return core.WriteFile(f.path, data)
}

// Delete deletes the file.
func (f *file) Delete() error {
	return core.DeleteFile(f.path)
}

// Path returns the path of the file.
func (f *file) Path() string {
	return f.path
}

// dir is the default implementation of the Dir interface.
type dir struct {
	path string
}

// Create creates the directory.
func (d *dir) Create() error {
	return core.CreateDir(d.path)
}

// Delete deletes the directory.
func (d *dir) Delete() error {
	return core.DeleteDir(d.path)
}

// Path returns the path of the directory.
func (d *dir) Path() string {
	return d.path
}
