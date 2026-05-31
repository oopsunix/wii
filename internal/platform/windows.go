//go:build windows

package platform

import (
	"regexp"
)

type windowsPlatform struct{}

// Paths are matched after normalization to forward slashes.
// Git for Windows has three main PATH entries: cmd/, mingw64/bin/, usr/bin/.
// cmd/ contains git.exe and other user-facing tools; the other two are system-level.
var windowsSystemDirs = []*regexp.Regexp{
	regexp.MustCompile(`(?i)^[a-z]:/windows(/|$)`),
	regexp.MustCompile(`(?i)^/c/Windows|^/d/Windows|^/proc`),
	regexp.MustCompile(`(?i)[/\\]Git[/\\](mingw64|mingw32|usr)[/\\]`),
	regexp.MustCompile(`(?i):[/\\]Program Files[/\\]Git[/\\](mingw64|mingw32|usr)[/\\]`),
	regexp.MustCompile(`(?i):[/\\]Program Files \(x86\)[/\\]Git[/\\](mingw64|mingw32|usr)[/\\]`),
}

// Skip architecture-variant siblings (e.g. tool-x86_64.exe, tool-arm64.exe)
// but NOT .exe files in general — on Windows, .exe ARE the executables.
var windowsFamilySkip = regexp.MustCompile(`.*-(x86_64|i686|aarch64|arm64|armv7l)\.exe$`)

// Skip GUI-only executables that would open windows when probed.
var windowsGUISkip = regexp.MustCompile(`(?i)^(notepad|calc|mspaint|write|wordpad|charmap|sndrec32|sndvol32|magnify|osk|dvdplay|wmplayer|explorer|control|taskmgr|regedit|mmc|cmd|powershell|msedge|chrome|firefox|iexplore|MicrosoftEdge|hh|winhlp32|` +
	// Git GUI tools
	`git-gui|gitk|` +
	// Java GUI tools
	`appletviewer|javaws|javaw|jmc|jvisualvm|jconsole|jhat|policytool|` +
	// Windows Performance Toolkit
	`WPRUI|wpa|wpaexporter|wpr|xbootmgr|xbootmgrsleep|xperf|` +
	// Bandizip GUI (Arkview.x64.exe etc.)
	`Bandizip|Updater|Arkview.*|` +
	// Nmap GUI
	`zenmap|` +
	// Windows Apps
	`XboxPcAppAdminServer|XboxPcAppCE|microsoftstore|olk|store|wt|` +
	// Visual Studio
	`devenv|` +
	// Python GUI variant
	`pythonw|` +
	// Uninstall/Installer/Setup/Updater programs
	`Uninstall|Installer|Setup|Updater|` +
	// Other GUI tools
	`bash|sh|mintty|winpty-agent)\.exe$`)

func (p *windowsPlatform) SystemDirs() []*regexp.Regexp {
	return windowsSystemDirs
}

func (p *windowsPlatform) FamilySkip() *regexp.Regexp {
	return windowsFamilySkip
}

func (p *windowsPlatform) GUISkip() *regexp.Regexp {
	return windowsGUISkip
}

// New returns the Windows platform implementation.
func New() Platform {
	return &windowsPlatform{}
}

func (p *windowsPlatform) SectionColor(label string) Color {
	switch label {
	case "Scoop":
		return Color{Bold: true, ANSI: "\033[36m"} // Cyan
	case "Chocolatey":
		return Color{Bold: true, ANSI: "\033[35m"} // Magenta
	case "npm Global":
		return Color{Bold: false, ANSI: "\033[31m"} // Red
	case "Cargo":
		return Color{Bold: false, ANSI: "\033[33m"} // Yellow
	case "Go":
		return Color{Bold: false, ANSI: "\033[36m"} // Cyan
	case "AppData":
		return Color{Bold: false, ANSI: "\033[34m"} // Blue
	case "MinGW":
		return Color{Bold: false, ANSI: "\033[32m"} // Green
	default:
		return Color{Bold: false, ANSI: "\033[0m"} // Default
	}
}

// SectionLabel is no longer used for detailed labeling.
// The scan package now handles all label logic (User/Global/package manager).
func (p *windowsPlatform) SectionLabel(dir string) string {
	return "Other"
}
