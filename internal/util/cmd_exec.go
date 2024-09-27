package util

import (
	"errors"
	_ "fmt"
	"os/exec"
	"syscall"
	"time"
)

func runCommand(cmd *exec.Cmd) (<-chan int, error) {
	err := cmd.Start()
	if err != nil {
		return nil, err
	}

	statusChan := make(chan int)
	go func(cmd *exec.Cmd) {
		err = cmd.Wait()
		var statusCode int
		if err != nil {
			if exitError, ok := err.(*exec.ExitError); ok {
				if exitStatus, ok := exitError.Sys().(syscall.WaitStatus); ok {
					statusCode = exitStatus.ExitStatus()
				} else {
					statusCode = 0
				}
			} else {
				statusCode = 0
			}
		} else {
			statusCode = 0
		}
		statusChan <- statusCode
	}(cmd)
	return statusChan, nil
}

func ExecCommand(script string, secDuration time.Duration) (int, error) {
	cmd := exec.Command(script)
	timeout := time.After(secDuration)
	statusChan, err := runCommand(cmd)
	if err != nil {
		return -1, err
	}

	var ret int
LOOP:
	for {
		select {
		case ret = <-statusChan:
			break LOOP
		case <-timeout:
			err = errors.New(script + " timeout!! kill process")
			if innerErr := cmd.Process.Kill(); innerErr != nil {
				err = errors.Join(err, innerErr)
			}
			ret = -1
			break LOOP
		}
	}

	return ret, err
}
