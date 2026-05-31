//go:build windows

package provider

import (
	"context"
	"regexp"
	"strings"
	"time"
)

type scoopProvider struct{}

var scoopLineRe = regexp.MustCompile(`^(\S+)\s+(\d\S*)`)

func init() {
	Register(&scoopProvider{})
}

func (p *scoopProvider) Name() string { return "scoop" }

func (p *scoopProvider) Available() bool {
	if _, err := lookPath("scoop"); err == nil {
		return true
	}
	out, err := RunCommand(context.Background(), 2*time.Second, "powershell.exe", "-NoProfile", "-Command", "Get-Command scoop -ErrorAction SilentlyContinue | Select-Object -ExpandProperty Source")
	return err == nil && strings.TrimSpace(out) != ""
}

func (p *scoopProvider) Fetch(ctx context.Context) error {
	out, err := RunCommand(ctx, 5*time.Second, "scoop", "list")
	if err != nil {
		out, err = RunCommand(ctx, 5*time.Second, "powershell.exe", "-NoProfile", "-Command", "scoop list")
		if err != nil {
			return err
		}
	}

	for _, line := range strings.Split(out, "\n") {
		if m := scoopLineRe.FindStringSubmatch(line); len(m) > 2 {
			name := strings.TrimSpace(m[1])
			ver := strings.TrimSpace(m[2])
			if name != "" && ver != "" {
				CacheSet(name, ver)
			}
		}
	}
	return nil
}
