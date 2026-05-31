//go:build freebsd || openbsd || netbsd || dragonfly

package platform

import (
	"regexp"
)

type bsdPlatform struct{}

var bsdSystemDirs = []*regexp.Regexp{
	regexp.MustCompile(`^(/bin|/sbin|/usr/bin|/usr/sbin|/rescue/)`),
}

var bsdGUISkip = regexp.MustCompile(`^(xdg-open|notify-send|zenity|xterm|xclock|xlogo|xeyes|xcalc|xedit|xman|xclipboard|startx|xinit|gvfs-open)$`)

func (p *bsdPlatform) SystemDirs() []*regexp.Regexp {
	return bsdSystemDirs
}

func (p *bsdPlatform) FamilySkip() *regexp.Regexp {
	return nil
}

func (p *bsdPlatform) GUISkip() *regexp.Regexp {
	return bsdGUISkip
}

// New returns the BSD platform implementation.
func New() Platform {
	return &bsdPlatform{}
}

func (p *bsdPlatform) SectionColor(label string) Color {
	switch label {
	case "User Local":
		return Color{Bold: false, ANSI: "\033[36m"} // Cyan
	case "System Local":
		return Color{Bold: false, ANSI: "\033[32m"} // Green
	default:
		return Color{Bold: false, ANSI: "\033[0m"} // Default
	}
}

func (p *bsdPlatform) SectionLabel(dir string) string {
	return "Other"
}
