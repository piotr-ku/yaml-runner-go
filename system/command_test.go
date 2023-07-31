package system

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestNewCommand is a test function that tests the NewCommand function.
//
// It initializes a command string, gets the current working directory,
// creates a new command, sets a timeout, and performs various tests.
// The function uses the assert.Equal function to compare the expected
// values with the actual values returned by the NewCommand function.
// It asserts that the command, environment length, directory, timeout,
// shell, stdout, stderr, return code, and error values are as expected.
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

// TestNewCommandGetCwdError is a unit test for the NewCommand function
// when there is an error in getting the current working directory.
//
// It mocks the functionGetwd variable to return an error and
// a non-existing directory path. Then it asserts that calling NewCommand
// with a specific command string will cause a panic. Finally, it restores
// the original value of the functionGetwd variable. This test is used
// to ensure that NewCommand handles the error in getting
// the current working directory correctly.
func TestNewCommandGetCwdError(t *testing.T) {
	functionGetwd = func() (string, error) {
		return "/not/existing/directory", errors.New("os.Getwd error")
	}
	assert.Panics(t, func() { NewCommand("echo test") })
	functionGetwd = os.Getwd
}

// TestCommand is a test function that validates the behavior of
// the Command function.
//
// It tests a set of predefined commands and verifies the return codes, stdout,
// stderr, and error values of each command. The function uses the `NewCommand`
// function to create a new command instance and sets the shell to `/bin/bash`.
// It then executes the command and compares the expected values with
// the actual values. The `assert.Equal` function is used to check that each
// expected value matches the corresponding actual value.
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

// TestCommandEnviroment tests the command environment.
//
// It sets up a command with the given environment, executes the command, and
// verifies the expected stdout.
func TestCommandEnviroment(t *testing.T) {
	command := "echo ${VAR1}"
	// run command
	cmd := NewCommand(command)
	cmd.Environment = map[string]string{"VAR1": "test"}
	_ = cmd.Execute()

	// Verify expected stdout
	assert.Equal(t, "test", cmd.Stdout)
}

// TestCommandWorkingDirectory tests the command working directory.
//
// It sets up a command with the given working directory, executes
// the command, and verifies the expected stdout.
func TestCommandWorkingDirectory(t *testing.T) {
	command := "pwd"
	// run command
	cmd := NewCommand(command)
	cmd.Directory = "/"
	_ = cmd.Execute()

	// Verify expected stdout
	assert.Equal(t, "/", cmd.Stdout)
}

// TestCommandShell tests the command shell.
//
// It sets up a command with the given shell, executes the command, and
// verifies the expected stdout.
func TestCommandShell(t *testing.T) {
	command := "echo $0"
	// run command
	cmd := NewCommand(command)
	cmd.Shell = "/bin/bash"
	_ = cmd.Execute()

	// Verify expected stdout
	assert.Equal(t, "/bin/bash", cmd.Stdout)
}
