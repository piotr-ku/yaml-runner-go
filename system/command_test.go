package system

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewCommand(t *testing.T) {
	command := "echo test"
	pwd, _ := os.Getwd()
	c := NewCommand(command)
	const timeout int = 5

	tests := []struct {
		Expected any
		Got      any
		Desc     string
	}{
		{Expected: command, Got: c.Command, Desc: "command"},
		{Expected: 0, Got: len(c.Environment), Desc: "environment length"},
		{Expected: pwd, Got: c.Directory, Desc: "directory"},
		{Expected: timeout, Got: c.Timeout, Desc: "timeout"},
		{Expected: "/bin/sh", Got: c.Shell, Desc: "shell"},
		{Expected: "", Got: c.Stdout, Desc: "stdout"},
		{Expected: "", Got: c.Stderr, Desc: "stderr"},
		{Expected: 0, Got: c.Rc, Desc: "return code"},
		{Expected: nil, Got: c.Error, Desc: "error"},
	}

	for _, test := range tests {
		assert.Equal(t, test.Expected, test.Got, test)
	}
}

func TestNewCommandGetCwdError(t *testing.T) {
	functionGetwd = func() (string, error) {
		return "/not/existing/directory", errors.New("os.Getwd error")
	}
	assert.Panics(t, func() { NewCommand("echo test") })
	functionGetwd = os.Getwd
}

func TestCommand(t *testing.T) {
	// commands
	var commands = []struct {
		Command string
		Rc      int
		Stdout  string
		Stderr  string
		Error   error
	}{
		{"echo test", 0, "test", "", nil},
		{"echo test 1>&2", 0, "", "test", nil},
		{"exit 1", 1, "", "", errors.New("exit status 1")},
	}

	for _, command := range commands {
		// run command
		cmd := NewCommand(command.Command)
		cmd.Shell = "/bin/bash"
		_ = cmd.Execute()

		var tests = []struct {
			Expected    any
			Got         any
			Description string
		}{
			{command.Rc, cmd.Rc, "return codes"},
			{command.Stdout, cmd.Stdout, "stdout"},
			{command.Stderr, cmd.Stderr, "stderr"},
			{command.Error == nil, cmd.Error == nil, "error"},
		}

		for _, test := range tests {
			assert.Equal(t, test.Expected, test.Got,
				"Unexpected %s for `%s`.", test.Description, cmd.Command)
		}
	}
}

func TestCommandEnviroment(t *testing.T) {
	command := "echo ${VAR1}"
	// run command
	cmd := NewCommand(command)
	cmd.Environment = map[string]string{"VAR1": "test"}
	_ = cmd.Execute()

	// Verify expected stdout
	assert.Equal(t, "test", cmd.Stdout)
}

func TestCommandWorkingDirectory(t *testing.T) {
	command := "pwd"
	// run command
	cmd := NewCommand(command)
	cmd.Directory = "/"
	_ = cmd.Execute()

	// Verify expected stdout
	assert.Equal(t, "/", cmd.Stdout)
}

func TestCommandShell(t *testing.T) {
	command := "echo $0"
	// run command
	cmd := NewCommand(command)
	cmd.Shell = "/bin/bash"
	_ = cmd.Execute()

	// Verify expected stdout
	assert.Equal(t, "/bin/bash", cmd.Stdout)
}
