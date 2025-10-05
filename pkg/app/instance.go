package app

import (
	"os"
)

type Instance struct {
	// File system for the currently open data "file" (directory).
	DataFileSystem os.Root
}
