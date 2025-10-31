package util

import (
	"context"
	"errors"
	"os"
	"os/exec"
	"syscall"
)

// RunAttached starts the command attached to the current stdio and waits for it.
func RunAttached(cmd *exec.Cmd) error {
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		return err
	}
	if err := cmd.Wait(); err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			if status, ok := exitErr.Sys().(syscall.WaitStatus); ok {
				if status.Signaled() && status.Signal() == syscall.SIGINT {
					return context.Canceled
				}
			}
		}
		return err
	}
	return nil
}
