package system

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"golang.org/x/exp/slog"
)

type LogConfig struct {
	File  string `file,validate:"filepath"`
	Level string `minimal_level,validate:"oneof=debug info warn error"`
	Quiet bool
	Json  bool
}

var loggers map[string]*slog.Logger
var testingStdout bytes.Buffer
var testingStderr bytes.Buffer

// LogInit initializes the logging system based on the provided configuration.
// It sets up loggers for writing to stdout/stderr or file, and sets the minimum logging level.
// If the configuration specifies "testing_buffer" as the file, it redirects logging output to a testing buffer.
// The loggers are stored in the loggers map for later use.
func LogInit(config LogConfig) {
	// stdout/stderr
	var stdout io.Writer = os.Stdout
	var stderr io.Writer = os.Stderr

	// buffer for testing
	if config.File == "testing_buffer" {
		testingStdout.Reset()
		testingStderr.Reset()
		stdout = &testingStdout
		stderr = &testingStderr
	}

	// set minimum logging level
	var minimumLevel = new(slog.LevelVar)
	switch config.Level {
	default:
		minimumLevel.Set(slog.LevelDebug)
	case "info":
		minimumLevel.Set(slog.LevelInfo)
	case "warn":
		minimumLevel.Set(slog.LevelWarn)
	case "error":
		minimumLevel.Set(slog.LevelError)
	}

	// default options
	options := &slog.HandlerOptions{Level: minimumLevel}

	// handler creates a logger with the specified output, options, and JSON format flag.
	handler := func(output io.Writer, options *slog.HandlerOptions, json bool) *slog.Logger {
		if json {
			return slog.New(slog.NewJSONHandler(output, options))
		} else {
			return slog.New(slog.NewTextHandler(output, options))
		}
	}

	// We will collect loggers in the temporary variable.
	_loggers := map[string]*slog.Logger{}

	// Initialize file logger if the file path is specified and is not "testing_buffer".
	if config.File != "" && config.File != "testing_buffer" {
		f, err := os.OpenFile(config.File, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0600)
		// notest
		if err != nil {
			FatalError("IOError", err.Error())
		}
		_loggers["file"] = handler(f, options, config.Json)
	}

	// Initialize stdout logger if Quiet flag is not set.
	if !config.Quiet {
		_loggers["stdout"] = handler(stdout, options, config.Json)
		_loggers["stderr"] = handler(stderr, options, config.Json)
	}

	// Set the loggers variable to the collected loggers.
	loggers = _loggers
}

// Log saves a log message with the specified level and parameters to the configured log targets.
func Log(level string, message string, params ...interface{}) {
	for _, handler := range logTargets(level) {
		switch level {
		case "debug":
			loggers[handler].Debug(message, params...)
		case "info":
			loggers[handler].Info(message, params...)
		case "warn":
			loggers[handler].Warn(message, params...)
		case "error":
			loggers[handler].Error(message, params...)
		// notest
		default:
			loggers[handler].Warn(message, params...)
			FatalError("LogError", fmt.Sprintf("last log has incorrect level: %s", level))
		}
	}
}

// logTargets returns a list of log targets based on the specified level.
func logTargets(level string) []string {
	var targets []string
	var output string

	switch level {
	case "error":
		output = "stderr"
	default:
		output = "stdout"
	}

	for _, handler := range []string{output, "file"} {
		_, handler_enabled := loggers[handler]
		if handler_enabled {
			targets = append(targets, handler)
		}
	}

	return targets
}

type LogBuilder struct {
	level   string
	message string
	params  []interface{}
}

// NewLogBuilder creates a new LogBuilder instance with the specified log message.
func NewLogBuilder(message string) *LogBuilder {
	return &LogBuilder{
		level:   "INFO",
		message: message,
		params:  []interface{}{},
	}
}

// Level sets the logging level for the LogBuilder.
func (b *LogBuilder) Level(level string) *LogBuilder {
	b.level = level
	return b
}

// Set adds parameters to the LogBuilder.
func (b *LogBuilder) Set(params ...interface{}) *LogBuilder {
	b.params = append(b.params, params...)
	return b
}

// Save builds the log parameters and invokes the Log function to save the log message.
func (b *LogBuilder) Save() {
	Log(b.level, b.message, b.params...)
}
