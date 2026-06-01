//go:build windows || darwin || linux

package provider

import (
	"context"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"wii/internal/model"
)

type npmProvider struct{}

func init() {
	Register(&npmProvider{})
}

func (p *npmProvider) Name() string { return "npm Global" }

func (p *npmProvider) Available() bool {
	_, err := lookPath("npm")
	return err == nil
}

func (p *npmProvider) Fetch(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	out, err := RunCommand(ctx, 5*time.Second, "npm", "list", "-g", "--depth=0", "--json")
	if err != nil {
		return err
	}

	var result struct {
		Dependencies map[string]struct {
			Version string `json:"version"`
		} `json:"dependencies"`
	}
	if err := json.Unmarshal([]byte(out), &result); err != nil {
		// Fallback: try parsing non-JSON output
		return p.parseTextOutput(out)
	}

	for name, dep := range result.Dependencies {
		if name != "" && dep.Version != "" {
			CacheSet(name, dep.Version)
			// Also cache common aliases
			if name == "npm" {
				CacheSet("npx", dep.Version)
			}
		}
	}
	return nil
}

// FetchNames returns the set of globally installed npm package names.
func (p *npmProvider) FetchNames(ctx context.Context) map[string]bool {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	names := make(map[string]bool)
	out, err := RunCommand(ctx, 5*time.Second, "npm", "list", "-g", "--depth=0", "--json")
	if err != nil {
		return names
	}

	var result struct {
		Dependencies map[string]struct {
			Version string `json:"version"`
		} `json:"dependencies"`
	}
	if err := json.Unmarshal([]byte(out), &result); err != nil {
		return names
	}

	for name, dep := range result.Dependencies {
		if name != "" && dep.Version != "" {
			names[name] = true
		}
	}
	return names
}

// FetchEntries returns one entry per globally installed npm package,
// using the package name (e.g. @anthropic-ai/claude-code) as the tool name.
func (p *npmProvider) FetchEntries(ctx context.Context) []model.Tool {
	prefix := npmPrefix()
	if prefix == "" {
		return nil
	}

	binDir := filepath.Join(prefix, "bin")

	// npm global node_modules location differs by platform:
	// macOS/Linux: {prefix}/lib/node_modules, Windows: {prefix}/node_modules
	modulesDir := filepath.Join(prefix, "lib", "node_modules")
	if runtime.GOOS == "windows" {
		modulesDir = filepath.Join(prefix, "node_modules")
	}
	entries, err := os.ReadDir(modulesDir)
	if err != nil {
		return nil
	}

	var tools []model.Tool
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		if strings.HasPrefix(entry.Name(), "@") {
			// Scoped packages: @scope/pkg
			scopeDir := filepath.Join(modulesDir, entry.Name())
			scopeEntries, err := os.ReadDir(scopeDir)
			if err != nil {
				continue
			}
			for _, se := range scopeEntries {
				if se.IsDir() && !strings.HasPrefix(se.Name(), ".") {
					pkgName := entry.Name() + "/" + se.Name()
					tools = appendNpmEntry(tools, pkgName, filepath.Join(scopeDir, se.Name()), binDir)
				}
			}
			continue
		}
		tools = appendNpmEntry(tools, entry.Name(), filepath.Join(modulesDir, entry.Name()), binDir)
	}
	return tools
}

func appendNpmEntry(tools []model.Tool, pkgName, pkgDir, binDir string) []model.Tool {
	ver, binName := readPkgMeta(pkgDir, pkgName)
	if ver == "" {
		ver = "?"
	}
	tools = append(tools, model.Tool{
		Name:    pkgName,
		Version: ver,
		Path:    filepath.Join(binDir, binName),
		Source:  "npm Global",
	})
	return tools
}

// readPkgMeta reads the version and first binary name from a package's package.json.
func readPkgMeta(pkgDir, pkgName string) (version, binName string) {
	data, err := os.ReadFile(filepath.Join(pkgDir, "package.json"))
	if err != nil {
		return "", ""
	}
	var pkg struct {
		Version string      `json:"version"`
		Bin     any `json:"bin"`
	}
	if json.Unmarshal(data, &pkg) != nil {
		return "", ""
	}
	binName = resolveBinName(pkg.Bin, pkgName)
	return pkg.Version, binName
}

// resolveBinName extracts the binary name from the bin field.
// Prefers the key matching the package name; falls back to the first key.
func resolveBinName(bin any, pkgName string) string {
	base := filepath.Base(pkgName)
	switch v := bin.(type) {
	case string:
		return base
	case map[string]any:
		if _, ok := v[base]; ok {
			return base
		}
		for name := range v {
			return name
		}
	}
	return base
}

// npmPrefix returns the npm global prefix directory (e.g. /opt/homebrew).
func npmPrefix() string {
	cmd := exec.Command("npm", "config", "get", "prefix")
	out, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

func (p *npmProvider) parseTextOutput(out string) error {
	for line := range strings.SplitSeq(out, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "{") || strings.HasPrefix(line, "}") {
			continue
		}
		// Format: /path/to/lib
		// ├── package@version
		// └── package@version
		if idx := strings.LastIndex(line, "@"); idx > 0 {
			name := line[strings.LastIndex(line, " ")+1 : idx]
			ver := line[idx+1:]
			name = strings.TrimSpace(name)
			if name != "" && ver != "" {
				CacheSet(name, ver)
			}
		}
	}
	return nil
}
