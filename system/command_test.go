package system

import (
	"errors"
	"os"
	"testing"
)

func TestNewCommand(t *testing.T) {
	command := "echo test"
	pwd, _ := os.Getwd()
	c := NewCommand(command)

	// Verify that the command is set correctly
	if c.Command != command {
		t.Errorf("Unexpected message. Expected: %s, Got: %s", command, c.Command)
	}

	// Verify that the environment is not set
	if len(c.Environment) != 0 {
		t.Errorf("Unexpected environment length. Expected: %d, Got: %d", 0, len(c.Environment))
	}

	// Verify that the directory is set to "os.Getwd()" by default
	if c.Directory != pwd {
		t.Errorf("Unexpected directory. Expected: %s, Got: %s", pwd, c.Directory)
	}

	// Verify that the timeout is set to 5 by default
	if c.Timeout != 5 {
		t.Errorf("Unexpected timeout. Expected: %d, Got: %d", 5, c.Timeout)
	}

	// Verify that the shell is set to "/bin/sh" by default
	if c.Shell != "/bin/sh" {
		t.Errorf("Unexpected shell. Expected: %s, Got: %s", "/bin/sh", c.Shell)
	}

	// Verify that the stdout is empty by default
	if c.Stdout != "" {
		t.Errorf("Unexpected stdout. Expected: %s, Got: %s", "", c.Stdout)
	}

	// Verify that the stderr is empty by default
	if c.Stderr != "" {
		t.Errorf("Unexpected stderr. Expected: %s, Got: %s", "", c.Stderr)
	}

	// Verify that the return code is 0 by default
	if c.Rc != 0 {
		t.Errorf("Unexpected return code. Expected: %d, Got: %d", 0, c.Rc)
	}

	// Verify that the error is nil by default
	if c.Error != nil {
		t.Errorf("Unexpected error value. Expected: %v, Got: %v", nil, c.Error)
	}
}

func TestCommand(t *testing.T) {
	// tests table
	var tests = []struct {
		Command string
		Rc      int
		Stdout  string
		Stderr  string
		Error   error
	}{
		{Command: "echo test", Rc: 0, Stdout: "test", Stderr: "", Error: nil},
		{Command: "echo test 1>&2", Rc: 0, Stdout: "", Stderr: "test", Error: nil},
		{Command: "exit 1", Rc: 1, Stdout: "", Stderr: "", Error: errors.New("exit status 1")},
	}

	for _, test := range tests {
		// run command
		cmd := NewCommand(test.Command)
		cmd.Shell = "/bin/bash"
		cmd.Execute()

		// Verify expected return code
		if cmd.Rc != test.Rc {
			t.Errorf("Unexpected return code for `%s`. Expected: %d, Got: %d", test.Command, test.Rc, cmd.Rc)
		}

		// Verify expected stdout
		if cmd.Stdout != test.Stdout {
			t.Errorf("Unexpected stdout for `%s`. Expected: %s, Got: %s", test.Command, test.Stdout, cmd.Stdout)
		}

		// Verify expected stderr
		if cmd.Stderr != test.Stderr {
			t.Errorf("Unexpected stderr for `%s`. Expected: %s, Got: %s", test.Command, test.Stderr, cmd.Stderr)
		}

		// Verify expected error value
		if cmd.Error != nil && cmd.Error.Error() != test.Error.Error() {
			t.Errorf("Unexpected error value for `%s`. Expected: %v, Got: %v", test.Command, test.Error.Error(), cmd.Error.Error())
		}
	}
}

func TestCommandEnviroment(t *testing.T) {
	command := "echo ${VAR1}"
	// run command
	cmd := NewCommand(command)
	cmd.Environment = map[string]string{"VAR1": "test"}
	cmd.Execute()

	// Verify expected stdout
	if cmd.Stdout != "test" {
		t.Errorf("Unexpected stdout for `%s`. Expected: %s, Got: %s", command, "test", cmd.Stdout)
	}
}

func TestCommandWorkingDirectory(t *testing.T) {
	command := "pwd"
	expected := "/"
	// run command
	cmd := NewCommand(command)
	cmd.Directory = "/"
	cmd.Execute()

	// Verify expected stdout
	if cmd.Stdout != expected {
		t.Errorf("Unexpected stdout for `%s`. Expected: %s, Got: %s", command, expected, cmd.Stdout)
	}
}

func TestCommandShell(t *testing.T) {
	command := "echo ${SHELL}"
	expected := "/bin/bash"
	// run command
	cmd := NewCommand(command)
	cmd.Shell = "/bin/bash"
	cmd.Execute()

	// Verify expected stdout
	if cmd.Stdout != expected {
		t.Errorf("Unexpected stdout for `%s`. Expected: %s, Got: %s", command, expected, cmd.Stdout)
	}
}
