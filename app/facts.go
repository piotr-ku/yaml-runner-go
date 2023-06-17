package app

import (
	"github.com/piotr-ku/yaml-runner-go/system"
)

// Facts represents a map of fact names to their corresponding values.
type Facts map[string]string

// gatherFacts collects facts by executing commands and saves the results in a temporary storage.
func gatherFacts(facts []Fact) Facts {
	// temporary storage
	gatheredFacts := Facts{}

	for _, fact := range facts {
		// create command
		c := system.NewCommand(fact.Command)
		// set shell
		if fact.Shell != "" {
			c.Shell = fact.Shell
		}
		// execute command
		c.Execute()
		// log
		LogFactGathered(fact, c)

		// save fact value to the temporary storage
		gatheredFacts[fact.Name] = c.Stdout
	}

	// log gathered facts
	system.Log("debug", "facts", "facts", gatheredFacts)

	return gatheredFacts
}

// LogFactGathered logs the details of a fact that has been gathered.
func LogFactGathered(fact Fact, c system.Command) {
	// determine log level based on command execution result
	var level string
	switch {
	case c.Error != nil:
		level = "error"
	case c.Stderr != "":
		level = "warn"
	default:
		level = "debug"
	}

	// build and save log entry
	l := system.NewLogBuilder("fact gathered")
	l.Level(level)
	l.Set("name", fact.Name)
	l.Set("command", fact.Command)
	l.Set("dir", c.Directory)
	l.Set("rc", c.Rc)
	l.Set("stdout", c.Stdout)
	l.Set("stderr", c.Stderr)
	l.Set("error", c.Error)
	l.Save()
}
