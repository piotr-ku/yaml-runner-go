package app

import (
	"testing"

	"github.com/piotr-ku/yaml-runner-go/system"
	"github.com/stretchr/testify/assert"
)

// empty represent a regular expression for empty string.
var empty = "^$"

// tests represent the test to perform
type test struct {
	name        string
	facts       []Fact
	expected    []string
	stdout      string
	stderr      string
	environment map[string]string
}

var tests = []test{
	{
		name: "Single fact with successful command execution",
		facts: []Fact{
			{
				Name:    "TEST1",
				Command: "echo test1",
				Shell:   "/bin/bash",
			},
		},
		expected:    []string{"test1"},
		stdout:      empty,
		stderr:      empty,
		environment: map[string]string{"TEST1": "test1"},
	},
	{
		name: "Single fact with command execution writing to stderr",
		facts: []Fact{
			{
				Name:    "TEST2",
				Command: "echo test2 1>&2",
				Shell:   "/bin/bash",
			},
		},
		expected: []string{""},
		stdout: "level=WARN msg=\"fact gathered\" name=TEST2 " +
			"command=\"echo test2 1>&2\" " +
			"dir=[^ ]+ rc=0 stdout=\"\" stderr=test2 error=<nil>",
		stderr:      empty,
		environment: map[string]string{},
	},
	{
		name: "Single fact with command execution resulting in an error",
		facts: []Fact{
			{
				Name:    "TEST3",
				Command: "echo test3 1>&2; exit 1;",
				Shell:   "/bin/bash",
			},
		},
		expected: []string{""},
		stdout:   empty,
		stderr: "level=ERROR msg=\"fact gathered\" name=TEST3 " +
			"command=\"echo test3 1>&2; exit 1;\" " +
			"dir=[^ ]+ rc=1 stdout=\"\" stderr=test3 " +
			"error=\"exit status 1\"",
		environment: map[string]string{},
	},
}

func TestGatherFacts(t *testing.T) {
	for _, test := range tests {
		// Set log settings and clear buffers
		system.LogInit(system.LogConfig{
			File:  "testing_buffer",
			Level: "info",
			Quiet: false,
			JSON:  false,
		})

		// Gather facts
		facts := gatherFacts(test.facts)

		// Test stdout
		assert.Equal(t, test.expected[0],
			facts[test.facts[0].Name].Result.Stdout)

		// Test logs
		assert.Regexp(t, test.stderr, system.GetTestingStderr())
		assert.Regexp(t, test.stdout, system.GetTestingStdout())

		// Test environment
		assert.Equal(t, test.environment, facts.toEnvironment())
	}
}
