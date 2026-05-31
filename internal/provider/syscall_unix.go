//go:build !windows

package provider

import "os/exec"

func configureCmd(cmd *exec.Cmd) {}
