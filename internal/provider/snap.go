//go:build linux

package provider

import (
	"context"
	"strings"
	"time"
)

type snapProvider struct{}

func init() {
	Register(&snapProvider{})
}

func (p *snapProvider) Name() string { return "snap" }

func (p *snapProvider) Available() bool {
	_, err := lookPath("snap")
	return err == nil
}

func (p *snapProvider) Fetch(ctx context.Context) error {
	out, err := RunCommand(ctx, 3*time.Second, "snap", "list")
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
		if name == "Name" { // skip header
			continue
		}
		if name != "" && ver != "" {
			CacheSet(name, ver)
		}
	}
	return nil
}
