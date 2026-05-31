//go:build !windows

package devenv

import "syscall"

func sysProcAttr() *syscall.SysProcAttr {
	return nil
}
