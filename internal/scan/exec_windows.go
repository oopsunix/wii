//go:build windows

package scan

import (
	"os"
	"path/filepath"
	"strings"
)

func isExecutable(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".exe", ".bat", ".cmd", ".com":
		return true
	case "":
		// No extension — could be a shell script in Git Bash or scoop shim
		info, err := os.Stat(path)
		if err != nil {
			return false
		}
		return info.Mode()&0111 != 0
	default:
		return false
	}
}

func normalizePath(dir string) string {
	// Convert backslashes to forward slashes for consistent matching
	return strings.ReplaceAll(dir, "\\", "/")
}
