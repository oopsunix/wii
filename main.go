package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"runtime"

	"github.com/oopsunix/wii/internal/config"
	"github.com/oopsunix/wii/internal/devenv"
	"github.com/oopsunix/wii/internal/model"
	"github.com/oopsunix/wii/internal/platform"
	"github.com/oopsunix/wii/internal/probe"
	"github.com/oopsunix/wii/internal/provider"
	"github.com/oopsunix/wii/internal/render"
	"github.com/oopsunix/wii/internal/scan"
	"github.com/oopsunix/wii/internal/update"
)

func main() {
	// Command line flags
	format := flag.String("f", "table", "Output format: table, json, csv")
	noColor := flag.Bool("nc", false, "Disable color output")
	concurrency := flag.Int("c", runtime.NumCPU(), "Number of concurrent workers")
	showVersion := flag.Bool("v", false, "Show version information")
	flag.Parse()

	if *showVersion {
		fmt.Println(config.Version)
		os.Exit(0)
	}

	fmt.Fprintln(os.Stderr, config.BuildInfo())

	// Async update check
	updateDone := make(chan update.Result, 1)
	go func() {
		updateDone <- update.CheckAndUpdate()
	}()
	defer func() {
		printUpdateResult(<-updateDone)
	}()

	// Build config
	cfg := &model.Config{
		Format:  *format,
		Batch:   *concurrency,
		NoColor: *noColor,
	}

	// Validate format
	switch cfg.Format {
	case "json", "csv", "table":
	default:
		fmt.Fprintf(os.Stderr, "Error: unsupported format %q (use table, json, csv)\n", cfg.Format)
		os.Exit(1)
	}

	// Validate concurrency
	if cfg.Batch < 1 {
		cfg.Batch = 1
	}
	if cfg.Batch > 128 {
		cfg.Batch = 128
	}

	plat := platform.New()
	scanner := scan.NewScanner(plat)

	// Populate package manager whitelists before scanning
	ctx := context.Background()
	allNames := provider.ResolveNames(ctx)
	for label, names := range allNames {
		scan.SetWhitelist(label, names)
	}

	result := scanner.ScanPath()

	// Detect development environments
	devEnvs := devenv.DetectAll()

	if result.Total == 0 && len(devEnvs) == 0 {
		renderer := render.New(cfg)
		renderer.Render(nil, nil)
		return
	}

	// Probe tool versions
	if result.Total > 0 {
		// Phase 1: bulk-query package managers for cached versions
		fmt.Fprintf(os.Stderr, "Concurrency: %d workers\n", cfg.Batch)
		fmt.Fprintf(os.Stderr, "Querying package managers...\n")
		provider.FetchAll(ctx)

		// Phase 1.5: replace PATH-scanned entries with provider-level entries
		// (e.g. Homebrew formulae instead of individual binaries)
		if entries := provider.ResolveEntries(ctx); len(entries) > 0 {
			for source := range entries {
				filtered := result.Candidates[:0]
				for _, c := range result.Candidates {
					if c.Source != source {
						filtered = append(filtered, c)
					}
				}
				result.Candidates = filtered
			}
			for _, es := range entries {
				result.Candidates = append(result.Candidates, es...)
			}
			result.Total = len(result.Candidates)
		}

		// Phase 2: probe remaining tools in parallel
		fmt.Fprintf(os.Stderr, "Probing versions...\n")
		probe.ProbeVersions(ctx, result.Candidates, cfg.Batch)

		// Phase 3: deduplicate family tools (uv/uvw/uvx -> uv)
		result.Candidates = probe.DeduplicateByFamily(result.Candidates)
	}

	sections := render.GroupBySection(result.Candidates)
	render.SortSections(sections, cfg.Sort)

	renderer := render.New(cfg)
	renderer.Render(devEnvs, sections)
}

func printUpdateResult(r update.Result) {
	if !r.OK {
		if r.Err != nil {
			fmt.Fprintf(os.Stderr, "Update check failed: %v\n", r.Err)
		}
		return
	}
	if r.Manual {
		fmt.Fprintf(os.Stderr, "New version %s available, please download manually: https://github.com/oopsunix/wii/releases\n", r.Version)
		return
	}
	fmt.Fprintf(os.Stderr, "Updated to %s, please restart to apply.\n", r.Version)
}
