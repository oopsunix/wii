//go:build linux

package provider

import (
	"context"
	"strings"
	"time"
)

type flatpakProvider struct{}

func init() {
	Register(&flatpakProvider{})
}

func (p *flatpakProvider) Name() string { return "flatpak" }

func (p *flatpakProvider) Available() bool {
	_, err := lookPath("flatpak")
	return err == nil
}

func (p *flatpakProvider) Fetch(ctx context.Context) error {
	out, err := RunCommand(ctx, 3*time.Second, "flatpak", "list", "--columns=application,version")
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
