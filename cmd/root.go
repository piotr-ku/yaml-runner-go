package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	ConfigFile     string
	LogFile        string
	LogJSON        bool
	QuietMode      bool
	DebugMode      bool
	DaemonInterval string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use: "yaml-runner-go",
	Short: `An application that executes commands based on rules specified
in a YAML file.`,
	Long: `The application executes commands based on the rules defined
in a YAML file. It can be run once or as a daemon to execute
commands at specific intervals.`,
	Args: cobra.NoArgs,
}

// Execute adds all child commands to the root command
// and sets flags appropriately. This is called by main.main().
// It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err) // nolint:revive
		os.Exit(1)       // nolint:revive
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&ConfigFile, "config", "./config.yaml",
		"configuration file in yaml format")
	rootCmd.PersistentFlags().StringVar(&DaemonInterval, "interval", "",
		"set daemon interval")
	rootCmd.PersistentFlags().StringVar(&LogFile, "log", "",
		"enable logging to the file")
	rootCmd.PersistentFlags().BoolVar(&LogJSON, "json", false,
		"enable JSON formatting for the output")
	rootCmd.PersistentFlags().BoolVar(&QuietMode, "quiet", false,
		"enable quiet mode")
	rootCmd.PersistentFlags().BoolVar(&DebugMode, "debug", false,
		"enable debug logging")
}
