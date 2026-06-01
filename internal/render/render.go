package render

import (
	"os"
	"runtime"
	"strings"

	"wii/internal/model"
)

// Renderer defines the interface for output rendering.
type Renderer interface {
	Render(devEnvs []model.DevEnv, sections []model.Section)
}

// New returns the appropriate Renderer for the given format.
func New(cfg *model.Config) Renderer {
	switch cfg.Format {
	case "json":
		return &jsonRenderer{cfg: cfg}
	case "csv":
		return &csvRenderer{cfg: cfg}
	default:
		return &compactRenderer{cfg: cfg}
	}
}

// GroupBySection groups candidates by their source label, preserving insertion order.
func GroupBySection(candidates []model.Tool) []model.Section {
	labelOrder := []string{}
	labelSet := make(map[string]bool)
	labelTools := make(map[string][]model.Tool)

	for _, t := range candidates {
		if !labelSet[t.Source] {
			labelSet[t.Source] = true
			labelOrder = append(labelOrder, t.Source)
		}
		labelTools[t.Source] = append(labelTools[t.Source], t)
	}

	sections := make([]model.Section, 0, len(labelOrder))
	for _, label := range labelOrder {
		sections = append(sections, model.Section{
			Label: label,
			Tools: labelTools[label],
		})
	}
	return sections
}

// sectionPriority defines the display order of sections.
var sectionPriority = map[string]int{
	"User Local":   0,
	"System Local": 1,
}

// SortSections sorts sections and their tools based on the sort key.
func SortSections(sections []model.Section, sortKey string) {
	switch sortKey {
	case "name":
		for i := range sections {
			sortByName(sections[i].Tools)
		}
	case "path":
		for i := range sections {
			sortByPath(sections[i].Tools)
		}
	case "source":
		sortByLabel(sections)
		for i := range sections {
			sortByName(sections[i].Tools)
		}
	default:
		// Default order: User Local → System Local → package managers (alphabetical)
		for i := range sections {
			sortByName(sections[i].Tools)
		}
		sortByPriority(sections)
	}
}

func sortByPriority(sections []model.Section) {
	for i := 1; i < len(sections); i++ {
		for j := i; j > 0; j-- {
			if sectionOrder(sections[j]) < sectionOrder(sections[j-1]) {
				sections[j], sections[j-1] = sections[j-1], sections[j]
			}
		}
	}
}

func sectionOrder(s model.Section) int {
	if p, ok := sectionPriority[s.Label]; ok {
		return p
	}
	return 100 // package managers sort after User/System Local
}

func sortByName(tools []model.Tool) {
	for i := 1; i < len(tools); i++ {
		for j := i; j > 0 && tools[j].Name < tools[j-1].Name; j-- {
			tools[j], tools[j-1] = tools[j-1], tools[j]
		}
	}
}

func sortByPath(tools []model.Tool) {
	for i := 1; i < len(tools); i++ {
		for j := i; j > 0 && tools[j].Path < tools[j-1].Path; j-- {
			tools[j], tools[j-1] = tools[j-1], tools[j]
		}
	}
}

func sortByLabel(sections []model.Section) {
	for i := 1; i < len(sections); i++ {
		for j := i; j > 0 && sections[j].Label < sections[j-1].Label; j-- {
			sections[j], sections[j-1] = sections[j-1], sections[j]
		}
	}
}

// displayPath returns the path to display, abbreviating $HOME to ~ unless fullPath is set.
// On Windows, always return the full path (no ~ abbreviation).
func displayPath(path string, fullPath bool) string {
	if fullPath || runtime.GOOS == "windows" {
		return path
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return path
	}
	return strings.Replace(path, home, "~", 1)
}


