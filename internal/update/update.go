package update

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/oopsunix/wii/internal/config"
)

const repo = "oopsunix/wii"

// Result holds the outcome of an update attempt.
type Result struct {
	OK       bool
	Version  string
	Err      error
	Manual   bool // true if user needs to download manually
}

// CheckAndUpdate checks for a new release and attempts to replace the running binary.
// Returns empty Result if no update is available. Times out after 30 seconds.
func CheckAndUpdate() Result {
	ch := make(chan Result, 1)
	go func() {
		ch <- doUpdate()
	}()
	select {
	case r := <-ch:
		return r
	case <-time.After(10 * time.Second):
		return Result{Err: fmt.Errorf("update check timed out")}
	}
}

func doUpdate() Result {
	r := checkAndDownload()
	if r.Err != nil {
		return r
	}
	if !r.OK {
		return r
	}

	if err := replaceBinary(); err != nil {
		r.Err = err
		r.Manual = true
	}
	return r
}

func checkAndDownload() Result {
	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get("https://api.github.com/repos/" + repo + "/releases/latest")
	if err != nil {
		return Result{Err: err}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return Result{Err: fmt.Errorf("GitHub API returned %d", resp.StatusCode)}
	}

	var release struct {
		TagName string `json:"tag_name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return Result{Err: err}
	}

	latest := strings.TrimPrefix(release.TagName, "v")
	current := config.Version
	if !isNewer(current, latest) {
		return Result{OK: false}
	}

	// Download new binary
	ext := ""
	if runtime.GOOS == "windows" {
		ext = ".exe"
	}
	assetName := fmt.Sprintf("wii_%s_%s%s", runtime.GOOS, runtime.GOARCH, ext)
	downloadURL := fmt.Sprintf("https://github.com/%s/releases/download/%s/%s", repo, release.TagName, assetName)

	exePath, err := os.Executable()
	if err != nil {
		return Result{Version: latest, Err: err}
	}
	exePath, err = resolveSymlink(exePath)
	if err != nil {
		return Result{Version: latest, Err: err}
	}
	tmpPath := exePath + ".new"

	dlClient := &http.Client{Timeout: 60 * time.Second}
	if err := downloadFile(dlClient, downloadURL, tmpPath); err != nil {
		os.Remove(tmpPath)
		return Result{Version: latest, Err: err}
	}

	return Result{OK: true, Version: latest}
}

func downloadFile(client *http.Client, url, dest string) error {
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download returned %d", resp.StatusCode)
	}

	f, err := os.OpenFile(dest, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o755)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, resp.Body)
	return err
}

func replaceBinary() error {
	exePath, err := os.Executable()
	if err != nil {
		return err
	}
	exePath, err = resolveSymlink(exePath)
	if err != nil {
		return err
	}
	tmpPath := exePath + ".new"
	oldPath := exePath + ".old"

	// Clean up leftover files from previous attempts
	os.Remove(oldPath)

	// Try direct rename (atomic on Unix, may fail on Windows due to file lock)
	if err := os.Rename(tmpPath, exePath); err == nil {
		return nil
	}

	// Fallback: rename running binary out of the way, then move new one in
	if err := os.Rename(exePath, oldPath); err != nil {
		os.Remove(tmpPath)
		return fmt.Errorf("cannot replace binary: %w", err)
	}
	if err := os.Rename(tmpPath, exePath); err != nil {
		// Restore original
		os.Rename(oldPath, exePath)
		return fmt.Errorf("cannot move new binary: %w", err)
	}

	return nil
}

func isNewer(current, latest string) bool {
	c := parseVersion(current)
	l := parseVersion(latest)
	for i := range 3 {
		if l[i] > c[i] {
			return true
		}
		if l[i] < c[i] {
			return false
		}
	}
	return false
}

func parseVersion(v string) [3]int {
	var parts [3]int
	for i, s := range strings.SplitN(v, ".", 3) {
		fmt.Sscanf(s, "%d", &parts[i])
	}
	return parts
}

func resolveSymlink(path string) (string, error) {
	fi, err := os.Lstat(path)
	if err != nil {
		return path, err
	}
	if fi.Mode()&os.ModeSymlink == 0 {
		return path, nil
	}
	return os.Readlink(path)
}
