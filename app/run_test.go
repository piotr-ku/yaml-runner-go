package app

import (
	"testing"

	"github.com/piotr-ku/yaml-runner-go/system"
	"github.com/stretchr/testify/assert"
)

const emptyConfigHash = 0xe8b4543d

// TestRunEmptyConfig tests the Run function with an empty configuration.
//
// The function initializes a Config struct with specific values for
// the Daemon, Logging, Facts, Actions, and Hash fields. It then asserts
// that the result of calling the Run function with a testing
// configuration file and an empty Config struct matches the expected value.
func TestRunEmptyConfig(t *testing.T) {
	expect := Config{
		Daemon: Daemon{
			Interval: "5s",
		},
		Logging: system.LogConfig{
			File:  "./yaml-runner-go.log",
			Level: "error",
			Quiet: true,
			JSON:  true,
		},
		Facts: []Fact{
			{
				Name:    "shellTest",
				Command: "echo $0",
				Shell:   "/bin/zsh",
				Result: system.Command{
					Command:     "",
					Environment: map[string]string(nil),
					Directory:   "",
					Timeout:     0,
					Shell:       "",
					Stdout:      "",
					Stderr:      "",
					Rc:          0, Error: error(nil),
				},
			},
			{
				Name: "apacheIsRunning",
				Command: "curl --connect-timeout 1 -s http://localhost:80/; " +
					"echo $?;",
				Shell: "",
				Result: system.Command{
					Command:     "",
					Environment: map[string]string(nil),
					Directory:   "",
					Timeout:     0,
					Shell:       "",
					Stdout:      "",
					Stderr:      "",
					Rc:          0,
					Error:       error(nil),
				},
			},
			{
				Name:    "loadAverage1",
				Command: "uptime | awk '{ print $9; }' | cut -d\\. -f1",
				Shell:   "",
				Result: system.Command{
					Command:     "",
					Environment: map[string]string(nil),
					Directory:   "",
					Timeout:     0,
					Shell:       "",
					Stdout:      "",
					Stderr:      "",
					Rc:          0,
					Error:       error(nil),
				},
			},
		},
		Actions: []Action{
			{
				Command: "echo $0",
				Rules:   []string(nil),
				Shell:   "/bin/zsh",
			},
			{
				Command: "echo \"Stopping apache\"",
				Rules: []string{
					"[[ ${loadAverage1} -gt 15 ]]",
					"[[ ${apacheIsRunning} -eq 0 ]]",
				},
				Shell: "",
			},
			{
				Command: "echo \"Starting apache\"",
				Rules: []string{
					"[[ ${loadAverage1} -lt 15 ]]",
					"[[ ${apacheIsRunning} -ne 0 ]]",
				},
				Shell: "",
			},
		},
		Hash: emptyConfigHash,
	}

	assert.Equal(t, expect, Run(testingConfigFile, Config{}))
}
