package config

// Version information set at build time via ldflags.
var (
	// Version is the current version of the program.
	Version = "1.0.0"

	// Commit is the git commit hash.
	Commit = "unknown"

	// Date is the build date.
	Date = "unknown"
)

// BuildInfo returns the full version string.
func BuildInfo() string {
	return Version
}
