package config

import "fmt"

// Version information set at build time via ldflags.
var (
	// Version is the current version of the program.
	Version = "1.0.0"

	// Commit is the git commit hash.
	Commit = "unknown"

	// Date is the build date.
	Date = "unknown"
)

// BuildInfo returns the full version string in the format: version-commit(date).
func BuildInfo() string {
	return fmt.Sprintf("%s-%s(%s)", Version, Commit, Date)
}
