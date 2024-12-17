package cmd

import (
	"github.com/piotr-ku/yaml-runner-go/app"
	"github.com/piotr-ku/yaml-runner-go/system"
	"github.com/spf13/cobra"
)

// oneshotCmd represents the oneshot command
var oneshotCmd = &cobra.Command{
	Use:   "oneshot",
	Short: "Runs actions ones end exit",
	Run: func(_ *cobra.Command, _ []string) {
		// Minimal logging level
		level := "info"
		if DebugMode {
			level = "debug"
		}
		overwrite := app.Config{
			// Default logging settings
			Logging: system.LogConfig{
				File:  LogFile,
				Quiet: QuietMode,
				JSON:  LogJSON,
				Level: level,
			},
		}
		app.Run(ConfigFile, overwrite)
	},
}

func init() {
	rootCmd.AddCommand(oneshotCmd)
}
