package system

import (
	"errors"
	"os"
	"testing"
)

const unexpectedValue string = "Unexpected %s for `%s`. Expected: %v, Got: %v"

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
		if test.Expected != test.Got {
			t.Errorf("Unexpected %v. Expected: %v, Got: %v",
				test.Desc, test.Expected, test.Got)
		}
	}
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
			{command.Rc, cmd.Rc, "return code"},
			{command.Stdout, cmd.Stdout, "stdout"},
			{command.Stderr, cmd.Stderr, "stderr"},
			{command.Error == nil, cmd.Error == nil, "error"},
		}

		for _, test := range tests {
			if test.Expected != test.Got {
				t.Errorf(unexpectedValue,
					test.Description, cmd.Command, test.Expected, test.Got)
			}
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
	if cmd.Stdout != "test" {
		t.Errorf(unexpectedValue,
			"environment variable", command, "test", cmd.Stdout)
	}
}

func TestCommandWorkingDirectory(t *testing.T) {
	command := "pwd"
	expected := "/"
	// run command
	cmd := NewCommand(command)
	cmd.Directory = "/"
	_ = cmd.Execute()

	// Verify expected stdout
	if cmd.Stdout != expected {
		t.Errorf(unexpectedValue, "working directory", command, expected,
			cmd.Stdout)
	}
}

func TestCommandShell(t *testing.T) {
	command := "echo $0"
	expected := "/bin/bash"
	// run command
	cmd := NewCommand(command)
	cmd.Shell = "/bin/bash"
	_ = cmd.Execute()

	// Verify expected stdout
	if cmd.Stdout != expected {
		t.Errorf(unexpectedValue, "shell", command, expected, cmd.Stdout)
	}
}
