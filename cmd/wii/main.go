package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"runtime"

	"wii/internal/config"
	"wii/internal/devenv"
	"wii/internal/model"
	"wii/internal/platform"
	"wii/internal/probe"
	"wii/internal/provider"
	"wii/internal/render"
	"wii/internal/scan"
)

func main() {
	// Command line flags
	format := flag.String("format", "table", "Output format: table, json, csv")
	noColor := flag.Bool("no-color", false, "Disable color output")
	concurrency := flag.Int("c", runtime.NumCPU(), "Number of concurrent workers")
	showVersion := flag.Bool("version", false, "Show version information")
	flag.Parse()

	// Handle version flag
	if *showVersion {
		fmt.Printf("wii %s\n", config.BuildInfo())
		os.Exit(0)
	}

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
		ctx := context.Background()

		// Phase 1: bulk-query package managers for cached versions
		fmt.Fprintf(os.Stderr, "Querying package managers...\n")
		provider.FetchAll(ctx)

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
