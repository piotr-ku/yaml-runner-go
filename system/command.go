package system

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"
)

// Command represents a system command to be executed.
type Command struct {
	Command     string            // The command to be executed.
	Environment map[string]string // Environment variables for the command.
	Directory   string            // Working directory for the command.
	Timeout     int               // Timeout duration in seconds.
	Shell       string            // Shell used to execute the command.
	Stdout      string            // Standard output of the command.
	Stderr      string            // Standard error of the command.
	Rc          int               // Return code of the command.
	Error       error             // Error encountered during command execution.
}

var functionGetwd = os.Getwd

// NewCommand creates a new Command with default settings.
func NewCommand(command string) Command {
	pwd, err := functionGetwd()
	const timeout int = 5
	if err != nil {
		panic(err.Error())
	}
	return Command{
		Command:   command,
		Directory: pwd,
		Timeout:   timeout,
		Shell:     "/bin/sh",
	}
}

// Execute executes the command and captures its output.
func (c *Command) Execute() error {
	// Set command timeout
	ctx, cancel := context.WithTimeout(context.Background(),
		time.Duration(c.Timeout)*time.Second)
	defer cancel()

	// Set command with context
	cmd := exec.CommandContext(ctx, c.Shell, "-c", c.Command)

	// Set environment variables
	cmd.Env = os.Environ()
	for key, value := range c.Environment {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%v", key, value))
	}

	// Set working directory
	cmd.Dir = c.Directory

	// Capture stdout/stderr
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()

	// Save command stdout/stderr and return code
	c.Stdout = strings.Trim(stdout.String(), "\n")
	c.Stderr = strings.Trim(stderr.String(), "\n")
	c.Rc = cmd.ProcessState.ExitCode()
	c.Error = err

	return err
}
