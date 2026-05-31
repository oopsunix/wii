package model

// Tool represents a discovered CLI tool on the system.
type Tool struct {
	Name    string // Command name as found in PATH
	Version string // Version string; "?" means not yet probed
	Path    string // Full filesystem path to the executable
	Source  string // Source label (e.g. Homebrew, Cargo, User Local)
}

// Section groups tools by their installation source.
type Section struct {
	Label string // Source label displayed in the table header
	Tools []Tool
}

// ScanResult holds the output of a PATH scan.
type ScanResult struct {
	Candidates []Tool
	Total      int
}

// DevEnv represents a detected development environment.
type DevEnv struct {
	Name      string // Display name (e.g. "Python", "Java")
	Version   string // Version string, empty if not installed
	Path      string // Full path to the executable
	Installed bool
}

// Config holds runtime configuration parsed from environment variables.
type Config struct {
	Format   string // Output format: table, json, csv
	Sort     string // Sort key: name, path, source
	ScanOnly bool   // Exit after scan, print candidate count
	Batch    int    // Parallel probe batch size (default 64)
	FullPath bool   // Show absolute paths instead of ~abbreviated
	NoColor  bool   // Disable ANSI colors
	ASCII    bool   // Use ASCII box-drawing instead of Unicode
}
