//go:build !windows

package scan

import (
	"os"
	"strings"
)

// isUserPath checks if a directory is under the user's home directory.
func isUserPath(dir string) bool {
	home, err := os.UserHomeDir()
	if err != nil {
		return false
	}
	return strings.HasPrefix(dir, home)
}
