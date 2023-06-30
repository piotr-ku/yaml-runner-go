package app

import (
	"github.com/piotr-ku/yaml-runner-go/system"
)

// Fact provides a data format for the facts defined
// in the configuration file.
type Fact struct {
	Name    string         `validate:"required"` // fact name
	Command string         `validate:"required"` // fact command
	Shell   string         // fact shell
	Result  system.Command // fact result
}

// LogFactGathered logs the details of a fact that has been gathered.
func (fact *Fact) logFactGathered(c system.Command) {
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

// Facts represents a map of fact names to their corresponding values.
type Facts map[string]Fact

func (facts Facts) toEnvironment() map[string]string {
	environment := make(map[string]string)

	for key, fact := range facts {
		if fact.Result.Stdout != "" && fact.Result.Rc == 0 {
			environment[key] = fact.Result.Stdout
		}
	}

	return environment
}

// gatherFacts collects facts by executing commands and saves the results
// in a temporary storage.
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
		_ = c.Execute()
		// log
		fact.logFactGathered(c)
		// add result
		fact.Result = c

		// save fact value to the temporary storage
		gatheredFacts[fact.Name] = fact
	}

	return gatheredFacts
}
