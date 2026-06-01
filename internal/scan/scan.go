package scan

import (
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/oopsunix/wii/internal/model"
	"github.com/oopsunix/wii/internal/platform"
)

// blocklistSuffixes are file suffixes to skip (not CLI tools).
var blocklistSuffixes = []string{"-config", ".py", ".sh", ".bat", ".cmd"}

// devEnvDirPatterns are path patterns that belong to known development environments.
// These directories are skipped because dev environments are already detected separately.
// Note: Python Scripts and Go bin directories are NOT skipped because they contain user-installed CLI tools.
var devEnvDirPatterns = []string{
	"/java/", "/jdk", "/jre",
	"/nodejs",
	"/dotnet",
	"/git/cmd", "/git/mingw64/bin",
	"/microsoft vs code/",
	"/windows kits/",      // Windows SDK / Performance Toolkit
	"/windows performance toolkit/",
}

// isDevEnvDir checks if a directory belongs to a known development environment.
// This handles more complex patterns that can't be expressed as simple substring matches.
func isDevEnvDir(dir string) bool {
	// Normalize: trim trailing slash, convert to lowercase
	lower := strings.TrimRight(strings.ToLower(dir), "/")

	for _, pattern := range devEnvDirPatterns {
		if strings.Contains(lower, pattern) {
			return true
		}
	}

	// Python runtime directory: contains /python/ but NOT /scripts
	if strings.Contains(lower, "/python/") && !strings.Contains(lower, "/scripts") {
		return true
	}

	return false
}

// pkgManagerDirPatterns maps path patterns to package manager names.
// Platform-specific patterns are added in initPkgManagers().
var pkgManagerDirPatterns = []struct {
	pattern string
	label   string
}{
	// Cross-platform
	{"/.cargo/", "Cargo"},
	{"/.local/bin", "User Local"},

	// Homebrew (macOS Apple Silicon, Intel, Linuxbrew)
	{"/homebrew/bin", "Homebrew"},
	{"/homebrew/homebrew/bin", "Homebrew"},

	// Windows
	{"/scoop/", "Scoop"},
	{"/appdata/roaming/npm", "npm Global"},
	{"/chocolatey/", "Chocolatey"},
	{"/choco/", "Chocolatey"},
	{"/windowsapps/", "winget"},  // Windows Store apps
}

// isPkgManagerDir checks if a directory is a package manager directory.
// This handles more complex patterns that can't be expressed as simple substring matches.
func isPkgManagerDir(dir string) (string, bool) {
	// Normalize: trim trailing slash, convert to lowercase
	lower := strings.TrimRight(strings.ToLower(dir), "/")

	// Check simple patterns first
	for _, pm := range pkgManagerDirPatterns {
		if strings.Contains(lower, pm.pattern) {
			return pm.label, true
		}
	}

	// Python Scripts directory: contains /python/ and ends with /scripts
	if strings.Contains(lower, "/python/") && strings.HasSuffix(lower, "/scripts") {
		return "Python Scripts", true
	}

	// Go tools directory: ends with /go/bin (user-installed Go tools)
	// Exclude the Go runtime itself (e.g., D:/Dev/Go/bin, C:/Go/bin)
	if strings.HasSuffix(lower, "/go/bin") {
		// Check if it's a user Go path (contains /users/ or /home/)
		if strings.Contains(lower, "/users/") || strings.Contains(lower, "/home/") {
			return "Go Tools", true
		}
	}

	return "", false
}

// pkgManagerWhitelist defines which tools to keep for each package manager.
// An empty map means keep all tools; a populated map acts as a filter.
var pkgManagerWhitelist = map[string]map[string]bool{
	"Cargo": {
		"cargo":         true,
		"rustc":         true,
		"rustup":        true,
		"rust-analyzer": true,
	},
	"Homebrew":       {}, // empty = keep all (brew-installed CLI tools)
	"npm Global":     {}, // empty = keep all (user-installed npm packages)
	"Python Scripts": {}, // empty = keep all (pip-installed CLI tools)
	"Go Tools":       {}, // empty = keep all (go install CLI tools)
	"winget":         {}, // empty = keep all
}

// SetWhitelist overrides the tool whitelist for a given package manager label.
func SetWhitelist(label string, names map[string]bool) {
	pkgManagerWhitelist[label] = names
}

// Scanner scans PATH directories for CLI tool candidates.
type Scanner struct {
	platform platform.Platform
}

func NewScanner(p platform.Platform) *Scanner {
	return &Scanner{platform: p}
}

func (s *Scanner) ScanPath() model.ScanResult {
	pathEnv := os.Getenv("PATH")
	sep := ":"
	if runtime.GOOS == "windows" {
		sep = ";"
	}

	dirs := strings.Split(pathEnv, sep)
	seenDirs := make(map[string]bool)
	seenNames := make(map[string]bool)
	systemDirs := s.platform.SystemDirs()

	var candidates []model.Tool

	for _, dir := range dirs {
		dir = normalizePath(dir)

		if dir == "" || seenDirs[dir] {
			continue
		}
		if platform.IsSystemDir(dir, systemDirs) {
			continue
		}

		// Determine label first - package manager dirs are never skipped
		label := resolveLabel(dir)
		isPkgManager := label == "Cargo" || label == "Scoop" || label == "npm Global" ||
			label == "Chocolatey" || label == "Python Scripts" || label == "Go Tools" ||
			label == "winget" || label == "NuGet" || label == "Homebrew"

		// Skip dev environment dirs only if NOT a package manager
		if !isPkgManager && isDevEnvDir(dir) {
			continue
		}

		info, err := os.Stat(dir)
		if err != nil || !info.IsDir() {
			continue
		}
		seenDirs[dir] = true

		entries, err := os.ReadDir(dir)
		if err != nil {
			continue
		}

		// Collect valid candidates in this directory
		var dirCandidates []struct {
			name    string
			fullDir string
		}
		dirName := filepath.Base(dir)

		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}
			name := entry.Name()
			if seenNames[name] {
				continue
			}
			// For non-pkg-manager dirs, skip scripts/configs
			if !isPkgManager && isSkippable(name) {
				continue
			}
			if !isExecutable(filepath.Join(dir, name)) {
				continue
			}
			if skip := s.platform.FamilySkip(); skip != nil && skip.MatchString(name) {
				continue
			}
			if skip := s.platform.GUISkip(); skip != nil && skip.MatchString(name) {
				continue
			}
			dirCandidates = append(dirCandidates, struct {
				name    string
				fullDir string
			}{name, dir})
		}

		if len(dirCandidates) == 0 {
			continue
		}

		if isPkgManager {
			// Package manager: keep whitelisted tools (empty whitelist = keep all)
			whitelist, hasWhitelist := pkgManagerWhitelist[label]
			for _, dc := range dirCandidates {
				if hasWhitelist && len(whitelist) > 0 && !whitelist[dc.name] {
					continue
				}
				seenNames[dc.name] = true
				candidates = append(candidates, model.Tool{
					Name:    dc.name,
					Version: "?",
					Path:    filepath.Join(dc.fullDir, dc.name),
					Source:  label,
				})
			}
		} else {
			// Non-package-manager: keep ONE tool per directory
			// Priority: tool matching directory name > first found
			chosen := dirCandidates[0]
			lowerDirName := strings.ToLower(dirName)
			for _, dc := range dirCandidates {
				base := strings.ToLower(strings.TrimSuffix(dc.name, ".exe"))
				if base == lowerDirName {
					chosen = dc
					break
				}
			}
			seenNames[chosen.name] = true
			candidates = append(candidates, model.Tool{
				Name:    chosen.name,
				Version: "?",
				Path:    filepath.Join(chosen.fullDir, chosen.name),
				Source:  label,
			})
		}
	}

	return model.ScanResult{
		Candidates: candidates,
		Total:      len(candidates),
	}
}

// isSkippable checks if a file should be skipped (script, config, etc.).
func isSkippable(name string) bool {
	lower := strings.ToLower(name)
	for _, suffix := range blocklistSuffixes {
		if strings.HasSuffix(lower, suffix) {
			return true
		}
	}
	return false
}

func resolveLabel(dir string) string {
	// Check package manager patterns (including complex patterns)
	if label, ok := isPkgManagerDir(dir); ok {
		return label
	}
	if isUserPath(dir) {
		return "User Local"
	}
	return "System Local"
}
