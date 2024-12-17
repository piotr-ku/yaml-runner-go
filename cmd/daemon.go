package cmd

import (
	"time"

	"github.com/piotr-ku/yaml-runner-go/app"
	"github.com/piotr-ku/yaml-runner-go/system"
	"github.com/spf13/cobra"
)

// daemonCmd represents the daemon command
var daemonCmd = &cobra.Command{
	Use:   "daemon",
	Short: "Run actions periodically in the background",
	Run: func(_ *cobra.Command, _ []string) {
		// Minimal logging level
		level := "info"
		if DebugMode {
			level = "debug"
		}
		overwrite := app.Config{
			// Default daemon settings
			Daemon: app.Daemon{
				Interval: DaemonInterval,
			},
			// Default logging settings
			Logging: system.LogConfig{
				File:  LogFile,
				Quiet: QuietMode,
				JSON:  LogJSON,
				Level: level,
			},
		}

		for {
			// Save start time
			startTime := time.Now()
			// Run application and save configuration
			config := app.Run(ConfigFile, overwrite)
			minInterval, _ := time.ParseDuration(config.Daemon.Interval)
			// Calculate how long we should wait for the next run
			stopTime := time.Now()
			runDuration := stopTime.Sub(startTime)
			// Sleep if run duration is less than minimal interval
			if runDuration < time.Duration(minInterval) {
				diff := minInterval.Milliseconds() - runDuration.Milliseconds()
				wait := time.Duration(diff) * time.Millisecond
				// Log
				system.Log("debug", "sleeping", "ms", wait.Milliseconds())
				// Wait
				time.Sleep(wait)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(daemonCmd)
}
