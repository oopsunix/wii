//go:build windows

package provider

import (
	"context"
	"regexp"
	"strings"
	"time"
)

type wingetProvider struct{}

var wingetLineRe = regexp.MustCompile(`^(.+?)\s{2,}\S+\s+(\d\S*)`)

func init() {
	Register(&wingetProvider{})
}

func (p *wingetProvider) Name() string { return "winget" }

func (p *wingetProvider) Available() bool {
	// winget may not be in Git Bash PATH; check via cmd.exe
	if _, err := lookPath("winget"); err == nil {
		return true
	}
	out, err := RunCommand(context.Background(), 2*time.Second, "cmd.exe", "/c", "where", "winget")
	return err == nil && strings.TrimSpace(out) != ""
}

func (p *wingetProvider) Fetch(ctx context.Context) error {
	out, err := RunCommand(ctx, 8*time.Second, "cmd.exe", "/c", "winget", "list")
	if err != nil {
		// Fallback to direct winget call
		out, err = RunCommand(ctx, 8*time.Second, "winget", "list")
		if err != nil {
			return err
		}
	}

	for _, line := range strings.Split(out, "\n") {
		if m := wingetLineRe.FindStringSubmatch(line); len(m) > 2 {
			name := strings.TrimSpace(m[1])
			ver := strings.TrimSpace(m[2])
			if name != "" && ver != "" {
				CacheSet(name, ver)
			}
		}
	}
	return nil
}
