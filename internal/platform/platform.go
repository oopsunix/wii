package platform

import "regexp"

// Color represents terminal color attributes.
type Color struct {
	Bold bool
	ANSI string // ANSI escape code, e.g. "\033[36m"
}

// Platform defines the interface for OS-specific behavior.
// Each platform (darwin, linux, windows, bsd) provides its own implementation.
type Platform interface {
	// SystemDirs returns regex patterns matching directories to skip during PATH scan.
	SystemDirs() []*regexp.Regexp

	// FamilySkip returns a regex matching architecture-variant filenames to skip.
	// Returns nil if no skip pattern is needed.
	FamilySkip() *regexp.Regexp

	// GUISkip returns a regex matching GUI-only executables to skip.
	// Returns nil if no skip pattern is needed.
	GUISkip() *regexp.Regexp

	// SectionLabel maps a directory path to its source label (e.g. "Homebrew", "Cargo").
	SectionLabel(dir string) string

	// SectionColor returns the color for a given source label.
	SectionColor(label string) Color
}

// IsSystemDir checks whether a directory matches any of the system dir patterns.
func IsSystemDir(dir string, patterns []*regexp.Regexp) bool {
	for _, re := range patterns {
		if re.MatchString(dir) {
			return true
		}
	}
	return false
}
