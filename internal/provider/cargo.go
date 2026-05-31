package provider

import (
	"context"
	"regexp"
	"strings"
	"time"
)

type cargoProvider struct{}

var cargoLineRe = regexp.MustCompile(`^(\S+)\s+v([^:]+):`)

func init() {
	Register(&cargoProvider{})
}

func (p *cargoProvider) Name() string { return "cargo" }

func (p *cargoProvider) Available() bool {
	_, err := lookPath("cargo")
	return err == nil
}

func (p *cargoProvider) Fetch(ctx context.Context) error {
	out, err := RunCommand(ctx, 3*time.Second, "cargo", "install", "--list")
	if err != nil {
		return err
	}

	for _, line := range strings.Split(out, "\n") {
		if m := cargoLineRe.FindStringSubmatch(line); len(m) > 2 {
			name := strings.TrimSpace(m[1])
			ver := strings.TrimSpace(m[2])
			if name != "" && ver != "" {
				CacheSet(name, ver)
			}
		}
	}
	return nil
}
