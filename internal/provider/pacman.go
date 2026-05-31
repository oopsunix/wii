//go:build linux || freebsd || openbsd || netbsd

package provider

import (
	"context"
	"strings"
	"time"
)

type pacmanProvider struct{}

func init() {
	Register(&pacmanProvider{})
}

func (p *pacmanProvider) Name() string { return "pacman" }

func (p *pacmanProvider) Available() bool {
	_, err := lookPath("pacman")
	return err == nil
}

func (p *pacmanProvider) Fetch(ctx context.Context) error {
	out, err := RunCommand(ctx, 3*time.Second, "pacman", "-Q")
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
		if name != "" && ver != "" {
			CacheSet(name, ver)
		}
	}
	return nil
}
