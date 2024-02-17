package common

import (
	"os"
)

// GetProjectRoot returns the root directory of the project, e.g. "/root/project/GoWebDAV/"
func GetProjectRoot() string {
	wd, _ := os.Getwd()
	// look for parents until we find a directory with a .git folder
	for {
		if _, err := os.Stat(wd + "/.git"); err == nil {
			break
		}
		wd = wd[:len(wd)-1]
	}

	return wd
}