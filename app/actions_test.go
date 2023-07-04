package app

import (
	"testing"

	"github.com/piotr-ku/yaml-runner-go/system"
	"github.com/stretchr/testify/assert"
)

const defaultShell = "/bin/bash"

func TestExecuteActions(t *testing.T) {
	tests := []struct {
		name    string
		actions []Action
		facts   Facts
		stdout  string
		stderr  string
	}{
		{
			name:    "No actions and no facts",
			actions: []Action{},
			facts:   Facts{},
			stdout:  empty,
			stderr:  empty,
		},
		{
			name:    "Facts with no actions",
			actions: []Action{},
			facts: Facts{
				"TEST1": Fact{Name: "TEST1", Command: "echo test1",
					Shell: defaultShell, Result: system.Command{
						Rc: 0, Stdout: "test1",
					},
				},
			},
			stdout: empty,
			stderr: empty,
		},
		{
			name: "Action without rules",
			actions: []Action{
				{
					Command: "echo action 1",
					Rules:   []string{},
					Shell:   defaultShell,
				},
			},
			facts: Facts{
				"TEST2": Fact{Name: "TEST2", Command: "echo test2",
					Shell: defaultShell, Result: system.Command{
						Rc: 0, Stdout: "test2",
					},
				},
			},
			stdout: "^time=[^ ]+ level=DEBUG msg=\"action executed\" " +
				"command=\"echo action 1\" " +
				"dir=[^ ]+ rc=0 stdout=\"action 1\" stderr=\"\" " +
				"error=<nil>\n$",
			stderr: empty,
		},
		{
			name: "Action with a rule returned 0 code",
			actions: []Action{
				{
					Command: "echo action 2",
					Rules: []string{
						"echo rule 1",
					},
					Shell: defaultShell,
				},
			},
			facts: Facts{
				"TEST3": Fact{Name: "TEST3", Command: "echo test3",
					Shell: defaultShell, Result: system.Command{
						Rc: 0, Stdout: "test3",
					},
				},
			},
			stdout: "^time=[^ ]+ level=DEBUG msg=\"rule checked\" " +
				"command=\"echo rule 1\" " +
				"dir=[^ ]+ rc=0 stdout=\"rule 1\" stderr=\"\" error=<nil>\n" +
				"time=[^ ]+ level=DEBUG msg=\"action executed\" " +
				"command=\"echo action 2\" " +
				"dir=[^ ]+ rc=0 stdout=\"action 2\" stderr=\"\" error=<nil>\n$",
			stderr: empty,
		},
		{
			name: "Action with a rule returned non 0 code",
			actions: []Action{
				{
					Command: "echo action 3",
					Rules: []string{
						"echo rule 2; exit 1;",
					},
					Shell: defaultShell,
				},
			},
			facts: Facts{
				"TEST4": Fact{Name: "TEST4", Command: "echo test4",
					Shell: defaultShell, Result: system.Command{
						Rc: 0, Stdout: "test4",
					},
				},
			},
			stdout: "^time=[^ ]+ level=DEBUG msg=\"rule checked\" " +
				"command=\"echo rule 2; exit 1;\" " +
				"dir=[^ ]+ rc=1 stdout=\"rule 2\" stderr=\"\" " +
				"error=\"exit status 1\"\n$",
			stderr: empty,
		},
		{
			name: "Action with two rules, first return zero, " +
				"second non-zero return code",
			actions: []Action{
				{
					Command: "echo action 4",
					Rules: []string{
						"echo rule 1;",
						"echo rule 2; exit 1;",
					},
					Shell: defaultShell,
				},
			},
			facts: Facts{},
			stdout: "^time=[^ ]+ level=DEBUG msg=\"rule checked\" " +
				"command=\"echo rule 1;\" " +
				"dir=[^ ]+ rc=0 stdout=\"rule 1\" stderr=\"\" error=<nil>\n" +
				"time=[^ ]+ level=DEBUG msg=\"rule checked\" " +
				"command=\"echo rule 2; exit 1;\" " +
				"dir=[^ ]+ rc=1 stdout=\"rule 2\" stderr=\"\" " +
				"error=\"exit status 1\"\n$",
			stderr: empty,
		},
		{
			name: "Action with two rules, first return non-zero, " +
				"second zero return code",
			actions: []Action{
				{
					Command: "echo action 5",
					Rules: []string{
						"echo rule 1; exit 1;",
						"echo rule 2;",
					},
					Shell: defaultShell,
				},
			},
			facts: Facts{},
			stdout: "^time=[^ ]+ level=DEBUG " +
				"msg=\"rule checked\" " +
				"command=\"echo rule 1; exit 1;\" " +
				"dir=[^ ]+ rc=1 stdout=\"rule 1\" stderr=\"\" " +
				"error=\"exit status 1\"\n$",
			stderr: empty,
		},
		{
			name: "Action returned non-zero code",
			actions: []Action{
				{
					Command: "echo action 6; exit 1",
					Shell:   defaultShell,
				},
			},
			stdout: empty,
			stderr: "^time=[^ ]+ level=ERROR msg=\"action executed\" " +
				"command=\"echo action 6; exit 1\" " +
				"dir=[^ ]+ rc=1 stdout=\"action 6\" " +
				"stderr=\"\" error=\"exit status 1\"\n$",
		},
		{
			name: "Action returned zero code but not empty stderr",
			actions: []Action{
				{
					Command: "echo action 7 1>&2",
					Shell:   defaultShell,
				},
			},
			stdout: "^time=[^ ]+ level=WARN msg=\"action executed\" " +
				"command=\"echo action 7 1>&2\" " +
				"dir=[^ ]+ rc=0 stdout=\"\" stderr=\"action 7\" " +
				"error=<nil>\n$",
			stderr: empty,
		},
	}

	for _, test := range tests {
		// Set log settings and clear buffers
		system.LogInit(system.LogConfig{
			File:  "testing_buffer",
			Level: "debug",
			Quiet: false,
			JSON:  false,
		})

		executeActions(test.actions, test.facts)
		assert.Regexp(t, test.stdout, system.GetTestingStdout())
		assert.Regexp(t, test.stderr, system.GetTestingStderr())
	}
}
