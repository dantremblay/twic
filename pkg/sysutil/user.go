package sysutil

import "os"

// IsRoot returns true if the current process is running as root (UID 0).
func IsRoot() bool {
	return os.Getuid() == 0
}

// FileExists returns true if the path exists.
func FileExists(path string) bool {
	_, err := os.Lstat(path)
	return err == nil
}
