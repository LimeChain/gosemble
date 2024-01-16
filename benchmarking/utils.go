package benchmarking

import (
	"os"
	"path/filepath"
	"runtime"
)

func fullPath(runtimePath string) (string, bool) {
	_, b, _, _ := runtime.Caller(0)
	basepath := filepath.Dir(b)
	fullPath := filepath.Join(basepath, "../", runtimePath)

	_, err := os.Stat(fullPath)
	isValidPath := err == nil || !os.IsNotExist(err)

	return fullPath, isValidPath
}
