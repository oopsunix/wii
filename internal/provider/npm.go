//go:build windows || darwin || linux

package provider

import (
	"context"
	"encoding/json"
	"strings"
	"time"
)

type npmProvider struct{}

func init() {
	Register(&npmProvider{})
}

func (p *npmProvider) Name() string { return "npm" }

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

func (p *npmProvider) parseTextOutput(out string) error {
	for _, line := range strings.Split(out, "\n") {
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
