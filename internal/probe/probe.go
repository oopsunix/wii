package probe

import (
	"context"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/oopsunix/wii/internal/model"
	"github.com/oopsunix/wii/internal/provider"
)

// ProbeVersions probes version strings for all candidates using a bounded worker pool.
func ProbeVersions(ctx context.Context, candidates []model.Tool, batchSize int) {
	total := len(candidates)
	if total == 0 {
		return
	}

	var wg sync.WaitGroup
	sem := make(chan struct{}, batchSize)
	var done int
	var mu sync.Mutex

	for i := range candidates {
		// Check cache first (populated by providers)
		if ver, ok := provider.CacheGet(candidates[i].Name); ok {
			candidates[i].Version = ver
			mu.Lock()
			done++
			mu.Unlock()
			continue
		}

		wg.Add(1)
		sem <- struct{}{}

		go func(idx int) {
			defer wg.Done()
			defer func() { <-sem }()

			ver := probeCommand(ctx, candidates[idx].Name)
			candidates[idx].Version = ver
			provider.CacheSet(candidates[idx].Name, ver)

			mu.Lock()
			done++
			if isTTY() {
				pct := done * 100 / total
				fmt.Fprintf(os.Stderr, "\r[%d/%d] %d%%  ", done, total, pct)
			}
			mu.Unlock()
		}(i)
	}

	wg.Wait()

	if isTTY() {
		fmt.Fprintf(os.Stderr, "\r\033[K")
	}
}

// probeCommand tries to get a version string for the named command.
// Uses --version and -V flags first, then falls back to "version" subcommand.
func probeCommand(ctx context.Context, name string) string {
	accelEnv := getAccelEnv(name)

	for _, flag := range []string{"--version", "-V", "version"} {
		var output string
		var err error

		if accelEnv != "" {
			output, err = provider.RunCommandWithEnv(ctx, 2*time.Second, accelEnv, name, flag)
		} else {
			output, err = provider.RunCommand(ctx, 2*time.Second, name, flag)
		}

		// Even if err != nil, we still check output (some commands exit non-zero but print version)
		if len(output) > 0 {
			ver := ExtractVersion(output)
			if ver != "?" {
				return ver
			}
		}

		// If command timed out, don't try more flags
		if err != nil && ctx.Err() != nil {
			break
		}
	}

	return "?"
}

func getAccelEnv(name string) string {
	if name == "brew" {
		return "HOMEBREW_NO_AUTO_UPDATE=1"
	}
	return ""
}

func isTTY() bool {
	fi, err := os.Stderr.Stat()
	if err != nil {
		return false
	}
	return fi.Mode()&os.ModeCharDevice != 0
}

// DeduplicateByFamily removes family duplicates: same family name = keep first only.
func DeduplicateByFamily(candidates []model.Tool) []model.Tool {
	seen := make(map[string]bool)
	var result []model.Tool

	for _, t := range candidates {
		family := extractFamily(t.Name)
		if family == "" {
			result = append(result, t)
			continue
		}

		if seen[family] {
			continue
		}
		seen[family] = true
		result = append(result, t)
	}
	return result
}

// specialFamilyNames maps base names to tool families for deduplication.
var specialFamilyNames = map[string]string{
	"uv":   "uv",
	"uvw":  "uv",
	"uvx":  "uv",
	"pip":  "pip",
	"pip3": "pip",
	// npm ecosystem
	"npm":     "npm",
	"npx":     "npm",
	"pnpm":    "pnpm",
	"pnpx":    "pnpm",
	"yarn":    "yarn",
	"yarnpkg": "yarn",
}

func extractFamily(name string) string {
	// Strip common extensions for matching
	base := strings.TrimSuffix(name, ".exe")
	base = strings.TrimSuffix(base, ".cmd")

	// Check special family names
	if family, ok := specialFamilyNames[base]; ok {
		return family
	}

	// Strip last -suffix
	if idx := strings.LastIndex(base, "-"); idx > 0 {
		return base[:idx]
	}
	// Strip from first digit
	for i, c := range base {
		if c >= '0' && c <= '9' {
			if i > 0 {
				return base[:i]
			}
			return ""
		}
	}
	return ""
}
