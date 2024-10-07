package util

import (
	"context"
	"errors"
	"os/exec"
	"syscall"
	"time"
)

func ExecCommand(script string, secDuration time.Duration) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), secDuration)
	defer cancel()
	cmd := exec.CommandContext(ctx, script)
	if err := cmd.Start(); err != nil {
		return -1, err
	}
	if err := cmd.Wait(); err != nil {
		if isTimeoutError(err) {
			err = errors.New(script + " timeout!! kill process")
			return -1, err
		}
	}
	exitCode := cmd.ProcessState.ExitCode()
	return exitCode, nil
}

// https://github.com/YoshikiShibata/oak/blob/master/java.go#L95
func isTimeoutError(err error) bool {
	exitErr, ok := err.(*exec.ExitError)
	if !ok {
		return false
	}

	status := exitErr.Sys().(syscall.WaitStatus)
	return status.Signaled()
}
