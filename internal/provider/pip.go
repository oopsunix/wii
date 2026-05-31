//go:build windows || darwin || linux

package provider

import (
	"context"
	"encoding/json"
	"strings"
	"time"
)

type pipProvider struct{}

func init() {
	Register(&pipProvider{})
}

func (p *pipProvider) Name() string { return "Python Scripts" }

func (p *pipProvider) Available() bool {
	// Check pip, pip3, pip3.x
	for _, cmd := range []string{"pip", "pip3"} {
		if _, err := lookPath(cmd); err == nil {
			return true
		}
	}
	return false
}

func (p *pipProvider) Fetch(ctx context.Context) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// Try pip list --format=json first
	out, err := RunCommand(ctx, 5*time.Second, "pip", "list", "--format=json")
	if err != nil {
		// Fallback to pip3
		out, err = RunCommand(ctx, 5*time.Second, "pip3", "list", "--format=json")
		if err != nil {
			return p.fetchWithoutJSON()
		}
	}

	var packages []struct {
		Name    string `json:"name"`
		Version string `json:"version"`
	}
	if err := json.Unmarshal([]byte(out), &packages); err != nil {
		return p.parseTextOutput(out)
	}

	for _, pkg := range packages {
		if pkg.Name != "" && pkg.Version != "" {
			CacheSet(pkg.Name, pkg.Version)
		}
	}
	return nil
}

func (p *pipProvider) fetchWithoutJSON() error {
	out, err := RunCommand(context.Background(), 5*time.Second, "pip", "list")
	if err != nil {
		out, err = RunCommand(context.Background(), 5*time.Second, "pip3", "list")
		if err != nil {
			return err
		}
	}
	return p.parseTextOutput(out)
}

func (p *pipProvider) parseTextOutput(out string) error {
	lines := strings.Split(out, "\n")
	for i, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" || i == 0 || i == 1 {
			continue // Skip header lines
		}
		parts := strings.Fields(line)
		if len(parts) >= 2 {
			name := parts[0]
			ver := parts[1]
			if name != "" && ver != "" {
				CacheSet(name, ver)
			}
		}
	}
	return nil
}
