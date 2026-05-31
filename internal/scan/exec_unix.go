//go:build !windows

package scan

import "os"

func isExecutable(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	return info.Mode()&0111 != 0
}

func normalizePath(dir string) string {
	return dir
}
