package utils

import (
	"os"
	"path/filepath"
	"strings"
)

func IsGoRun() bool {
	execPath, err := os.Executable()
	if err != nil {
		return false
	}
	return strings.HasPrefix(execPath, filepath.Join(os.TempDir(), "go-build"))
}
