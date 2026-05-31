package probe

import "regexp"

var (
	// semverPrimary matches versions like 1.2, 1.2.3, 1.2.3.4
	semverPrimary = regexp.MustCompile(`([0-9]+\.[0-9]+(\.[0-9]+){0,2})`)
	// semverFallback matches "Version 1.2.3" style output
	semverFallback = regexp.MustCompile(`[Vv]ersion\s+([0-9]+(\.[0-9]+)+)`)
)

// ExtractVersion extracts a semver-like version string from command output.
// Tries multiple lines (up to 3) and two regex patterns.
// Returns "?" if no version can be extracted.
func ExtractVersion(output string) string {
	if output == "" {
		return "?"
	}

	lines := splitLines(output, 3)
	for _, line := range lines {
		if line == "" {
			continue
		}
		if m := semverPrimary.FindString(line); m != "" {
			return m
		}
		if m := semverFallback.FindStringSubmatch(line); len(m) > 1 {
			return m[1]
		}
	}

	return "?"
}

// splitLines returns up to maxLines non-empty lines from s.
func splitLines(s string, maxLines int) []string {
	var lines []string
	start := 0
	for i := 0; i < len(s) && len(lines) < maxLines; i++ {
		if s[i] == '\n' {
			line := s[start:i]
			lines = append(lines, line)
			start = i + 1
		}
	}
	if len(lines) < maxLines && start < len(s) {
		lines = append(lines, s[start:])
	}
	return lines
}
