package app

import (
	"os"
)

// Instance represents a running instance of the application.
type Instance struct {
	// File system for the currently open data "file" (directory).
	DataFileSystem os.Root
}
