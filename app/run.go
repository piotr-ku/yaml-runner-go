package app

import (
	"os"
	"strings"

	"github.com/piotr-ku/yaml-runner-go/system"
)

var applicationStarted bool
var configurationHash uint32

// Run executes all the actions defined in the configuration file.
// It loads the configuration from the specified file and merges it with the provided merge configuration.
// It initializes logging and gathers facts before executing the actions.
//
// Parameters:
// - configFile: The path to the configuration file.
// - merge: The merge configuration to combine with the loaded configuration.
//
// notest
func Run(configFile string, configArgs Config) Config {
	// Default settings
	config := Config{
		// Default daemon settings
		Daemon: Daemon{
			Interval: "2s",
		},
		// Default logging settings
		Logging: system.LogConfig{
			File:  "",
			Quiet: false,
			Json:  false,
			Level: "info",
		},
	}

	// Load configuration file
	contentFile := LoadConfigFile(configFile)
	config.Merge(contentFile)

	// Load configuration from arguments
	config.Merge(configArgs)

	// Calculate configuration hash
	config.CalculateHash()

	// Initialize logging
	system.LogInit(system.LogConfig{
		File:  config.Logging.File,
		Quiet: config.Logging.Quiet,
		Json:  config.Logging.Json,
		Level: config.Logging.Level,
	})

	// Log application startup
	if !applicationStarted {
		system.Log("info", "starting", "args", strings.Join(os.Args[1:], " "))
		applicationStarted = true
	}

	// Check if we should reload configuration
	if config.Hash != configurationHash {
		// Update configuration hash
		configurationHash = config.Hash

		// Log configuration changes
		system.Log("debug", "configuration hash", "hash", configurationHash)
		system.Log("info", "configuration loaded", "file", configFile, "facts", len(config.Facts), "actions", len(config.Actions))
		system.Log("debug", "configuration dump", "config", config)
	}

	// Gather facts
	facts := gatherFacts(config.Facts)

	// Execute actions
	executeActions(config.Actions, facts)

	// Return configuration
	return config
}
