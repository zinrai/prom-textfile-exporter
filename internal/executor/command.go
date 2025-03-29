package executor

import (
	"context"
	"fmt"
	"os/exec"
	"syscall"
	"time"
)

type ExecuteCommandResult struct {
	Output     string
	ExitCode   int
	Successful bool
	Error      error
}

// executes a command and returns its output, exit code, and error
func ExecuteCommand(commandStr string, timeoutSec int) (string, int, error) {
	result := ExecuteCommandWithResult(commandStr, timeoutSec)
	return result.Output, result.ExitCode, result.Error
}

// executes a command and returns a comprehensive result
func ExecuteCommandWithResult(commandStr string, timeoutSec int) ExecuteCommandResult {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutSec)*time.Second)
	defer cancel()

	if len(commandStr) == 0 {
		return ExecuteCommandResult{
			Output:     "",
			ExitCode:   1,
			Successful: false,
			Error:      fmt.Errorf("empty command"),
		}
	}

	cmd := exec.CommandContext(ctx, "sh", "-c", commandStr)

	// Set process group IDs to ensure termination, including child processes
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}

	output, err := cmd.CombinedOutput()

	result := ExecuteCommandResult{
		Output:     string(output),
		ExitCode:   0,
		Successful: (err == nil),
		Error:      err,
	}

	if err != nil {
		// In case of timeout
		if ctx.Err() == context.DeadlineExceeded {
			// Terminate the entire process group
			if cmd.Process != nil {
				syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
			}
			result.ExitCode = 124 // Using 124 as timeout exit code (like timeout command)
			result.Error = fmt.Errorf("command timed out after %d seconds: %w", timeoutSec, ctx.Err())
		} else if exitErr, ok := err.(*exec.ExitError); ok {
			// Normal command execution error ( non-zero exit code )
			if status, ok := exitErr.Sys().(syscall.WaitStatus); ok {
				result.ExitCode = status.ExitStatus()
			}
		} else {
			// Command execution failure (e.g., command not found)
			result.ExitCode = 127 // Command not found or similar
		}
	}

	return result
}
