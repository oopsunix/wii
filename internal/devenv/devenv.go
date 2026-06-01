package devenv

import (
	"context"
	"os/exec"
	"regexp"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/oopsunix/wii/internal/model"
)

var versionRe = regexp.MustCompile(`([0-9]+\.[0-9]+(\.[0-9]+){0,2})`)

// detectors defines the dev environments to detect.
var detectors = []struct {
	name string
	cmd  string
	flag string
}{
	// Runtimes
	{"Python", "python", "--version"},
	{"Java", "java", "-version"},
	{"Go", "go", "version"},
	{"Rust", "rustc", "--version"},
	{"Node", "node", "--version"},
	{".NET", "dotnet", "--version"},
	{"PHP", "php", "--version"},
	{"Ruby", "ruby", "--version"},
	{"Deno", "deno", "--version"},
	{"Bun", "bun", "--version"},

	// Package managers
	{"pip", "pip", "--version"},
	{"npm", "npm", "--version"},
	{"pnpm", "pnpm", "--version"},
	{"yarn", "yarn", "--version"},
	{"Cargo", "cargo", "--version"},

	// Build tools
	{"Maven", "mvn", "--version"},
	{"Gradle", "gradle", "--version"},
	{"CMake", "cmake", "--version"},

	// Compilers
	{"GCC", "gcc", "--version"},
	{"Clang", "clang", "--version"},

	// Container
	{"Docker", "docker", "--version"},
	{"Podman", "podman", "--version"},
	{"Kubectl", "kubectl", "version"},

	// VCS
	{"Git", "git", "--version"},
}

// DetectAll probes all known dev environments in parallel.
func DetectAll() []model.DevEnv {
	results := make([]model.DevEnv, len(detectors))
	var wg sync.WaitGroup

	for i, d := range detectors {
		wg.Add(1)
		go func(idx int, name, cmd, flag string) {
			defer wg.Done()
			results[idx] = detectOne(name, cmd, flag)
		}(i, d.name, d.cmd, d.flag)
	}

	wg.Wait()

	// Sort alphabetically by name
	sort.Slice(results, func(i, j int) bool {
		return results[i].Name < results[j].Name
	})

	return results
}

func detectOne(name, cmd, flag string) model.DevEnv {
	// Find the executable path
	path, err := exec.LookPath(cmd)
	if err != nil {
		return model.DevEnv{Name: name, Version: "", Path: "", Installed: false}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	c := exec.CommandContext(ctx, cmd, flag)
	c.SysProcAttr = sysProcAttr()

	// CombinedOutput captures both stdout and stderr (java -version outputs to stderr)
	out, err := c.CombinedOutput()
	if err != nil && len(out) == 0 {
		return model.DevEnv{Name: name, Version: "", Path: path, Installed: false}
	}

	ver := extractVersion(string(out))
	return model.DevEnv{Name: name, Version: ver, Path: path, Installed: ver != ""}
}

func extractVersion(output string) string {
	// Try first 3 lines
	lines := strings.SplitN(output, "\n", 4)
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		if m := versionRe.FindString(line); m != "" {
			return m
		}
	}
	return ""
}
