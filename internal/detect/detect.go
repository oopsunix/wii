package detect

import "runtime"

// OSType represents the detected operating system category.
type OSType int

const (
	OSMacOS OSType = iota
	OSLinux
	OSWindows
	OSBSD
	OSUnknown
)

// String returns a human-readable name for the OS type.
func (o OSType) String() string {
	switch o {
	case OSMacOS:
		return "macos"
	case OSLinux:
		return "linux"
	case OSWindows:
		return "windows"
	case OSBSD:
		return "bsd"
	default:
		return "unknown"
	}
}

// OS returns the detected operating system type using runtime.GOOS.
func OS() OSType {
	switch runtime.GOOS {
	case "darwin":
		return OSMacOS
	case "linux":
		return OSLinux
	case "windows":
		return OSWindows
	case "freebsd", "openbsd", "netbsd", "dragonfly":
		return OSBSD
	default:
		return OSUnknown
	}
}
