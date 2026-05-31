package provider

import (
	"context"
	"os"
	"os/exec"
	"time"
)

// lookPath wraps exec.LookPath for testability.
var lookPath = exec.LookPath

// RunCommand runs a command with a timeout and returns its combined stdout+stderr.
func RunCommand(ctx context.Context, timeout time.Duration, name string, args ...string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Env = os.Environ()
	configureCmd(cmd)

	out, err := cmd.CombinedOutput()
	if err != nil && len(out) == 0 {
		return "", err
	}
	return string(out), nil
}

// RunCommandWithEnv runs a command with extra environment variables and returns combined stdout+stderr.
func RunCommandWithEnv(ctx context.Context, timeout time.Duration, env string, name string, args ...string) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Env = append(os.Environ(), env)
	configureCmd(cmd)

	out, err := cmd.CombinedOutput()
	if err != nil && len(out) == 0 {
		return "", err
	}
	return string(out), nil
}
