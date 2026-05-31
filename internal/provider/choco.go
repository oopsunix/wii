//go:build windows

package provider

import (
	"context"
	"strings"
	"time"
)

type chocoProvider struct{}

func init() {
	Register(&chocoProvider{})
}

func (p *chocoProvider) Name() string { return "choco" }

func (p *chocoProvider) Available() bool {
	_, err := lookPath("choco")
	return err == nil
}

func (p *chocoProvider) Fetch(ctx context.Context) error {
	out, err := RunCommand(ctx, 5*time.Second, "choco", "list", "--local-only", "--limit-output")
	if err != nil {
		return err
	}

	for _, line := range strings.Split(out, "\n") {
		if strings.HasPrefix(line, "Chocolatey") {
			continue
		}
		parts := strings.SplitN(line, "|", 2)
		if len(parts) != 2 {
			continue
		}
		name := strings.TrimSpace(parts[0])
		ver := strings.TrimSpace(parts[1])
		if name != "" && ver != "" {
			CacheSet(name, ver)
		}
	}
	return nil
}
