//go:build darwin || linux || freebsd || openbsd || netbsd

package provider

import (
	"context"
	"strings"
	"time"
)

type brewProvider struct{}

func init() {
	Register(&brewProvider{})
}

func (p *brewProvider) Name() string { return "brew" }

func (p *brewProvider) Available() bool {
	_, err := lookPath("brew")
	return err == nil
}

func (p *brewProvider) Fetch(ctx context.Context) error {
	out, err := RunCommandWithEnv(ctx, 3*time.Second, "HOMEBREW_NO_AUTO_UPDATE=1", "brew", "list", "--versions")
	if err != nil {
		return err
	}

	for _, line := range strings.Split(out, "\n") {
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		name := fields[0]
		ver := fields[1]
		if name == "" || ver == "" {
			continue
		}
		CacheSet(name, ver)
		// Register unversioned alias for versioned formulas (e.g. python@3.11 -> python)
		if idx := strings.Index(name, "@"); idx > 0 {
			CacheSet(name[:idx], ver)
		}
	}
	return nil
}
