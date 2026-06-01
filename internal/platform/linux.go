//go:build linux

package platform

import (
	"regexp"
)

type linuxPlatform struct{}

var linuxSystemDirs = []*regexp.Regexp{
	regexp.MustCompile(`^(/bin|/sbin|/usr/bin|/usr/sbin|/lib/systemd/)`),
}

var linuxFamilySkip = regexp.MustCompile(`.*-(x86_64|aarch64|i686|armv7l|armhf)$`)
// GUISkipRE matches GUI-only executables that must never be probed (probing would launch the GUI).
var GUISkipRE = regexp.MustCompile(`^(xdg-open|notify-send|zenity|kdialog|gvfs-open|gvfs-mount|gvfs-set-attribute|gvfs-copy|gvfs-move|gvfs-rm|gvfs-mkdir|gvfs-monitor-dir|gvfs-monitor-file|gvfs-ls|gvfs-info|gvfs-cat|gvfs-tree|gvfs-save|gnome-open|kde-open|exo-open|gvfsd|gvfsd-metadata|gnome-terminal|konsole|xterm|gucharmap|gnome-calculator|baobab|eog|evince|gedit|gnome-text-editor|nautilus|totem|yelp|systemctl|journalctl)$`)

func (p *linuxPlatform) SystemDirs() []*regexp.Regexp {
	return linuxSystemDirs
}

func (p *linuxPlatform) FamilySkip() *regexp.Regexp {
	return linuxFamilySkip
}

func (p *linuxPlatform) GUISkip() *regexp.Regexp {
	return GUISkipRE
}

// New returns the Linux platform implementation.
func New() Platform {
	return &linuxPlatform{}
}

func (p *linuxPlatform) SectionColor(label string) Color {
	switch label {
	case "Homebrew":
		return Color{Bold: true, ANSI: "\033[33m"} // Yellow
	case "Snap":
		return Color{Bold: true, ANSI: "\033[35m"} // Magenta
	case "Cargo":
		return Color{Bold: false, ANSI: "\033[33m"} // Yellow
	case "Go":
		return Color{Bold: false, ANSI: "\033[36m"} // Cyan
	case "nvm":
		return Color{Bold: false, ANSI: "\033[32m"} // Green
	case "pyenv":
		return Color{Bold: false, ANSI: "\033[34m"} // Blue
	case "Deno":
		return Color{Bold: false, ANSI: "\033[36m"} // Cyan
	case "Nix":
		return Color{Bold: true, ANSI: "\033[34m"} // Blue
	case "User Local":
		return Color{Bold: false, ANSI: "\033[36m"} // Cyan
	case "npm Global":
		return Color{Bold: false, ANSI: "\033[31m"} // Red
	case "System Local":
		return Color{Bold: false, ANSI: "\033[32m"} // Green
	default:
		return Color{Bold: false, ANSI: "\033[0m"} // Default
	}
}

func (p *linuxPlatform) SectionLabel(dir string) string {
	return "Other"
}
