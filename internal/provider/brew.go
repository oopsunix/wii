//go:build darwin || linux || freebsd || openbsd || netbsd

package provider

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"github.com/oopsunix/wii/internal/model"
)

// brewPrefixCandidates lists common Homebrew prefix paths in priority order.
var brewPrefixCandidates = []string{
	"/opt/homebrew",                  // macOS Apple Silicon
	"/usr/local",                     // macOS Intel
	"/home/linuxbrew/.linuxbrew",     // Linuxbrew
}

type brewProvider struct {
	prefix string // resolved Homebrew prefix (e.g. /opt/homebrew)
}

func init() {
	Register(&brewProvider{})
}

func (p *brewProvider) Name() string { return "Homebrew" }

func (p *brewProvider) Available() bool {
	for _, prefix := range brewPrefixCandidates {
		if info, err := os.Stat(filepath.Join(prefix, "Cellar")); err == nil && info.IsDir() {
			p.prefix = prefix
			return true
		}
	}
	return false
}

// readCellar reads installed formula names and versions from the Cellar directory.
// Layout: Cellar/{name}/{version}/
// Strategy: use directory name as version; fall back to INSTALL_RECEIPT.json
// when multiple version directories exist or directory name is unusable.
func readCellar(prefix string) map[string]string {
	versions := make(map[string]string)
	cellar := filepath.Join(prefix, "Cellar")
	entries, err := os.ReadDir(cellar)
	if err != nil {
		return versions
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		name := entry.Name()
		verDirs, err := os.ReadDir(filepath.Join(cellar, name))
		if err != nil || len(verDirs) == 0 {
			continue
		}

		// Single version directory: use directory name directly
		if len(verDirs) == 1 {
			ver := verDirs[0].Name()
			if ver != "" {
				versions[name] = ver
			}
			continue
		}

		// Multiple version directories: fall back to INSTALL_RECEIPT.json
		ver := readFormulaReceipt(filepath.Join(cellar, name, verDirs[0].Name()))
		if ver != "" {
			versions[name] = ver
		}
	}
	return versions
}

// readFormulaReceipt extracts the stable version from a formula's INSTALL_RECEIPT.json.
func readFormulaReceipt(dir string) string {
	data, err := os.ReadFile(filepath.Join(dir, "INSTALL_RECEIPT.json"))
	if err != nil {
		return ""
	}
	var receipt struct {
		Source struct {
			Versions struct {
				Stable string `json:"stable"`
			} `json:"versions"`
		} `json:"source"`
	}
	if err := json.Unmarshal(data, &receipt); err != nil {
		return ""
	}
	return receipt.Source.Versions.Stable
}

// readCaskroom reads installed cask names and versions from the Caskroom directory.
// Layout: Caskroom/{name}/{version}/ (ignore .metadata)
// Strategy: use directory name as version; fall back to INSTALL_RECEIPT.json
// in .metadata when no valid version directory is found.
func readCaskroom(prefix string) map[string]string {
	versions := make(map[string]string)
	caskroom := filepath.Join(prefix, "Caskroom")
	entries, err := os.ReadDir(caskroom)
	if err != nil {
		return versions
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		name := entry.Name()
		verDirs, err := os.ReadDir(filepath.Join(caskroom, name))
		if err != nil {
			continue
		}
		// Find version directories (skip dotfiles like .metadata)
		var verDir string
		for _, vd := range verDirs {
			if vd.IsDir() && !strings.HasPrefix(vd.Name(), ".") {
				verDir = vd.Name()
				break
			}
		}
		if verDir != "" {
			versions[name] = verDir
			continue
		}

		// No version directory found, fall back to INSTALL_RECEIPT.json
		ver := readCaskReceipt(filepath.Join(caskroom, name))
		if ver != "" {
			versions[name] = ver
		}
	}
	return versions
}

// readCaskReceipt extracts the version from a cask's .metadata/INSTALL_RECEIPT.json.
func readCaskReceipt(caskDir string) string {
	data, err := os.ReadFile(filepath.Join(caskDir, ".metadata", "INSTALL_RECEIPT.json"))
	if err != nil {
		return ""
	}
	var receipt struct {
		Source struct {
			Version string `json:"version"`
		} `json:"source"`
	}
	if err := json.Unmarshal(data, &receipt); err != nil {
		return ""
	}
	return receipt.Source.Version
}

func (p *brewProvider) Fetch(ctx context.Context) error {
	if p.prefix == "" {
		return nil
	}

	// Read formula versions from Cellar
	for name, ver := range readCellar(p.prefix) {
		CacheSet(name, ver)
		// Register unversioned alias for versioned formulas (e.g. python@3.11 -> python)
		if idx := strings.Index(name, "@"); idx > 0 {
			CacheSet(name[:idx], ver)
		}
	}

	// Read cask versions from Caskroom
	for name, ver := range readCaskroom(p.prefix) {
		CacheSet(name, ver)
	}

	return nil
}

// FetchNames returns the set of installed formula and cask names.
func (p *brewProvider) FetchNames(ctx context.Context) map[string]bool {
	names := make(map[string]bool)
	if p.prefix == "" {
		return names
	}

	cellar := filepath.Join(p.prefix, "Cellar")
	if entries, err := os.ReadDir(cellar); err == nil {
		for _, e := range entries {
			if e.IsDir() {
				names[e.Name()] = true
			}
		}
	}

	caskroom := filepath.Join(p.prefix, "Caskroom")
	if entries, err := os.ReadDir(caskroom); err == nil {
		for _, e := range entries {
			if e.IsDir() {
				names[e.Name()] = true
			}
		}
	}

	return names
}

// FetchEntries returns one entry per installed Homebrew formula/cask.
func (p *brewProvider) FetchEntries(ctx context.Context) []model.Tool {
	var entries []model.Tool
	if p.prefix == "" {
		return entries
	}

	binDir := filepath.Join(p.prefix, "bin")

	addEntry := func(name, ver string) {
		if ver == "" {
			ver = "?"
		}
		entries = append(entries, model.Tool{
			Name:    name,
			Version: ver,
			Path:    filepath.Join(binDir, name),
			Source:  "Homebrew",
		})
	}

	for name, ver := range readCellar(p.prefix) {
		addEntry(name, ver)
	}

	for name, ver := range readCaskroom(p.prefix) {
		addEntry(name, ver)
	}

	return entries
}
