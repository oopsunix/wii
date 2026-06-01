package render

import (
	"fmt"
	"os"
	"strings"

	"github.com/oopsunix/wii/internal/model"
)

// ANSI color codes.
const (
	colorReset  = "\033[0m"
	colorBold   = "\033[1m"
	colorDim    = "\033[2m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
)

type compactRenderer struct {
	cfg *model.Config
}

func (r *compactRenderer) Render(devEnvs []model.DevEnv, sections []model.Section) {
	noColor := r.cfg.NoColor || !isStdoutTTY()

	// Filter to only installed dev environments
	var installedEnvs []model.DevEnv
	for _, env := range devEnvs {
		if env.Installed && env.Version != "" {
			installedEnvs = append(installedEnvs, env)
		}
	}

	// Calculate column widths across all data
	nameWidth, verWidth := r.calcColumnWidths(installedEnvs, sections)

	// Dev environment summary
	if len(installedEnvs) > 0 {
		r.printDevEnvs(installedEnvs, noColor, nameWidth, verWidth)
		fmt.Println()
	}

	// Tool sections
	if len(sections) == 0 {
		if len(devEnvs) == 0 {
			r.printLine("No tools found.", noColor, "")
		}
		return
	}

	for i, sec := range sections {
		r.printSection(sec, noColor, nameWidth, verWidth)
		if i < len(sections)-1 {
			fmt.Println()
		}
	}
}

func (r *compactRenderer) calcColumnWidths(devEnvs []model.DevEnv, sections []model.Section) (int, int) {
	nameWidth := 10 // minimum width
	verWidth := 10  // minimum width

	// Check dev env names
	for _, env := range devEnvs {
		if len(env.Name) > nameWidth {
			nameWidth = len(env.Name)
		}
		ver := env.Version
		if !env.Installed || ver == "" {
			ver = "Not found"
		}
		if len(ver) > verWidth {
			verWidth = len(ver)
		}
	}

	// Check tool names and versions
	for _, sec := range sections {
		for _, t := range sec.Tools {
			name := strings.TrimSuffix(t.Name, ".exe")
			name = strings.TrimSuffix(name, ".cmd")
			if len(name) > nameWidth {
				nameWidth = len(name)
			}
			if len(t.Version) > verWidth {
				verWidth = len(t.Version)
			}
		}
	}

	return nameWidth + 2, verWidth + 2 // add padding
}

func (r *compactRenderer) printDevEnvs(devEnvs []model.DevEnv, noColor bool, nameWidth, verWidth int) {
	title := "Current development environment"
	r.printHeader(title, noColor)

	for i, env := range devEnvs {
		isLast := i == len(devEnvs)-1
		connector := "├─"
		if isLast {
			connector = "└─"
		}

		path := displayPath(env.Path, r.cfg.FullPath)

		if noColor {
			fmt.Printf(" %s %-*s %-*s %s\n", connector, nameWidth, env.Name, verWidth, env.Version, path)
		} else {
			fmt.Printf(" %s%s%s %s%-*s%s %s%-*s%s %s%s%s\n",
				colorDim, connector, colorReset,
				colorBold, nameWidth, env.Name, colorReset,
				colorGreen, verWidth, env.Version, colorReset,
				colorDim, path, colorReset,
			)
		}
	}
}

func (r *compactRenderer) printSection(sec model.Section, noColor bool, nameWidth, verWidth int) {
	r.printHeader(sec.Label, noColor)

	for i, t := range sec.Tools {
		isLast := i == len(sec.Tools)-1
		connector := "├─"
		if isLast {
			connector = "└─"
		}

		name := strings.TrimSuffix(t.Name, ".exe")
		name = strings.TrimSuffix(name, ".cmd")
		path := displayPath(t.Path, r.cfg.FullPath)

		if noColor {
			fmt.Printf(" %s %-*s %-*s %s\n", connector, nameWidth, name, verWidth, t.Version, path)
		} else {
			vColor := colorGreen
			if t.Version == "?" || t.Version == "" {
				vColor = colorYellow
			}
			fmt.Printf(" %s%s%s %s%-*s%s %s%-*s%s %s%s%s\n",
				colorDim, connector, colorReset,
				colorBold, nameWidth, name, colorReset,
				vColor, verWidth, t.Version, colorReset,
				colorDim, path, colorReset,
			)
		}
	}
}

func (r *compactRenderer) printHeader(title string, noColor bool) {
	if noColor {
		fmt.Printf("[*] %s\n", title)
	} else {
		fmt.Printf("%s[*]%s %s%s%s\n",
			colorBold, colorReset,
			colorBold, title, colorReset,
		)
	}
}

func (r *compactRenderer) printLine(text string, noColor bool, color string) {
	if noColor || color == "" {
		fmt.Printf("  %s\n", text)
	} else {
		fmt.Printf("  %s%s%s\n", color, text, colorReset)
	}
}

func isStdoutTTY() bool {
	fi, err := os.Stdout.Stat()
	if err != nil {
		return false
	}
	return fi.Mode()&os.ModeCharDevice != 0
}
