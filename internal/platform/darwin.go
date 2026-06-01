//go:build darwin

package platform

import (
	"regexp"
)

type darwinPlatform struct{}

var darwinSystemDirs = []*regexp.Regexp{
	regexp.MustCompile(`^(/bin|/sbin|/usr/bin|/usr/sbin|/System/|/Library/Apple/)`),
}

var darwinFamilySkip = regexp.MustCompile(`.*-(intel64|arm64)$`)
// GUISkipRE matches GUI-only executables that must never be probed (probing would launch the GUI).
var GUISkipRE = regexp.MustCompile(`^(open|say|pbcopy|pbpaste|screencapture|sips|osascript|pinentry-mac|iina|orb|orbctl)$`)

func (p *darwinPlatform) SystemDirs() []*regexp.Regexp {
	return darwinSystemDirs
}

func (p *darwinPlatform) FamilySkip() *regexp.Regexp {
	return darwinFamilySkip
}

func (p *darwinPlatform) GUISkip() *regexp.Regexp {
	return GUISkipRE
}

// New returns the macOS platform implementation.
func New() Platform {
	return &darwinPlatform{}
}

func (p *darwinPlatform) SectionColor(label string) Color {
	switch label {
	case "Homebrew":
		return Color{Bold: true, ANSI: "\033[33m"} // Yellow
	case "User Local":
		return Color{Bold: false, ANSI: "\033[36m"} // Cyan
	case "npm Global":
		return Color{Bold: false, ANSI: "\033[31m"} // Red
	case "Python Framework":
		return Color{Bold: false, ANSI: "\033[34m"} // Blue
	case "System Local":
		return Color{Bold: false, ANSI: "\033[32m"} // Green
	default:
		return Color{Bold: false, ANSI: "\033[0m"} // Default
	}
}

func (p *darwinPlatform) SectionLabel(dir string) string {
	return "Other"
}
