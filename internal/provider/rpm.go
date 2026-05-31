//go:build linux || freebsd || openbsd || netbsd

package provider

import (
	"context"
	"strings"
	"time"
)

type rpmProvider struct{}

func init() {
	Register(&rpmProvider{})
}

func (p *rpmProvider) Name() string { return "rpm" }

func (p *rpmProvider) Available() bool {
	_, err := lookPath("rpm")
	return err == nil
}

func (p *rpmProvider) Fetch(ctx context.Context) error {
	out, err := RunCommand(ctx, 5*time.Second, "rpm", "-qa", "--queryformat", "%{NAME} %{VERSION}\n")
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
		if name == "gpg-pubkey" {
			continue
		}
		if name != "" && ver != "" {
			CacheSet(name, ver)
		}
	}
	return nil
}
