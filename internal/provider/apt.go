//go:build linux || freebsd || openbsd || netbsd

package provider

import (
	"context"
	"strings"
	"time"
)

type aptProvider struct{}

func init() {
	Register(&aptProvider{})
}

func (p *aptProvider) Name() string { return "apt" }

func (p *aptProvider) Available() bool {
	_, err := lookPath("dpkg-query")
	return err == nil
}

func (p *aptProvider) Fetch(ctx context.Context) error {
	out, err := RunCommand(ctx, 3*time.Second, "dpkg-query", "-W", "-f", "${Package}\t${Version}\n")
	if err != nil {
		return err
	}

	for _, line := range strings.Split(out, "\n") {
		parts := strings.SplitN(line, "\t", 2)
		if len(parts) < 2 {
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
