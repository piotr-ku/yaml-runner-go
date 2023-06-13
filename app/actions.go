package app

import "github.com/piotr-ku/yaml-runner-go/system"

// executeActions executes a list of actions based on the provided facts.
func executeActions(actions []Action, facts Facts) {
	for _, action := range actions {
		// check action rules
		if checkActionRules(action, facts) {
			c := system.NewCommand(action.Command)
			c.Environment = facts
			c.Execute()
			logActionExecuted(action, &c)
		}
	}
}

// checkActionRules checks the rules of an action against the provided facts.
// It returns true if all rules pass, otherwise false.
func checkActionRules(action Action, facts Facts) bool {
	for _, rule := range action.Rules {
		c := system.NewCommand(rule)
		c.Environment = facts
		c.Execute()
		logRuleChecked(rule, &c)
		if c.Rc != 0 {
			return false
		}
	}
	return true
}

// logRuleChecked logs the result of a rule check.
func logRuleChecked(rule string, c *system.Command) {
	l := system.NewLogBuilder("rule checked")
	l.Level("debug")
	l.Set("command", rule)
	l.Set("dir", c.Directory)
	l.Set("rc", c.Rc)
	l.Set("stdout", c.Stdout)
	l.Set("stderr", c.Stderr)
	l.Set("error", c.Error)
	l.Save()
}

// logActionExecuted logs the execution of an action.
func logActionExecuted(action Action, c *system.Command) {
	var level string
	switch {
	case c.Error != nil:
		level = "error"
	case c.Stderr != "":
		level = "warn"
	default:
		level = "debug"
	}

	l := system.NewLogBuilder("action executed")
	l.Level(level)
	l.Set("command", action.Command)
	l.Set("dir", c.Directory)
	l.Set("rc", c.Rc)
	l.Set("stdout", c.Stdout)
	l.Set("stderr", c.Stderr)
	l.Set("error", c.Error)
	l.Save()
}
