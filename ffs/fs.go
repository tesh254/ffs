package ffs

// FileSystem provides an interface for file and directory operations.
// It abstracts the underlying file system, allowing for easier testing and extension.
type FileSystem interface {
	File(path string) File
	Dir(path string) Dir
}

// File provides an interface for file-specific operations.
type File interface {
	Read() ([]byte, error)
	Write(data []byte) error
	Delete() error
	Path() string
}

// Dir provides an interface for directory-specific operations.
type Dir interface {
	Create() error
	Delete() error
	Path() string
}
