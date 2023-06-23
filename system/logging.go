package system

import (
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"os"

	"golang.org/x/exp/slog"
)

// LogConfig represents the configuration options for logging.
type LogConfig struct {
	// The file path where log entries will be written.
	File string `file,validate:"filepath"`
	// The minimal log level to be logged.
	Level string `minimal_level,validate:"oneof=debug info warn error"`
	// Whether to suppress console output of log entries.
	Quiet bool
	// Whether to format log entries in JSON format.
	JSON bool
}

var loggers map[string]*slog.Logger
var testingStdout bytes.Buffer
var testingStderr bytes.Buffer

// LogInit initializes the logging system based on the provided configuration.
// It sets up loggers for writing to stdout/stderr or file, and sets the minimum
// logging level. If the configuration specifies "testing_buffer" as the file,
// it redirects logging output to a testing buffer.The loggers are stored in
// the loggers map for later use.
func LogInit(config LogConfig) {
	// stdout/stderr
	var stdout io.Writer = os.Stdout
	var stderr io.Writer = os.Stderr
	// log file permission
	const logFilePermission fs.FileMode = 0600

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

	// We will collect loggers in the temporary variable.
	_loggers := map[string]*slog.Logger{}

	// Initialize file logger if the file path is specified and
	// is not "testing_buffer".
	if config.File != "" && config.File != "testing_buffer" {
		f, err := os.OpenFile(config.File, os.O_RDWR|os.O_CREATE|os.O_APPEND,
			logFilePermission)
		// notest
		if err != nil {
			FatalError("IOError", err.Error())
		}
		_loggers["file"] = logHandler(f, options, config)
	}

	// Initialize stdout logger if Quiet flag is not set.
	if !config.Quiet {
		_loggers["stdout"] = logHandler(stdout, options, config)
		_loggers["stderr"] = logHandler(stderr, options, config)
	}

	// Set the loggers variable to the collected loggers.
	loggers = _loggers
}

// logHandler creates a logger with the specified output, options,
// and JSON format flag.
func logHandler(output io.Writer, options *slog.HandlerOptions,
	config LogConfig) *slog.Logger {
	switch config.JSON {
	case true:
		return slog.New(slog.NewJSONHandler(output, options))
	default:
		return slog.New(slog.NewTextHandler(output, options))
	}
}

// Log saves a log message with the specified level and parameters
// to the configured log targets.
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
			FatalError("LogError",
				fmt.Sprintf("last log has incorrect level: %s", level))
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
		_, handlerEnabled := loggers[handler]
		if handlerEnabled {
			targets = append(targets, handler)
		}
	}

	return targets
}

// LogBuilder represents a builder for creating log entries.
type LogBuilder struct {
	// The log level of the entry.
	level string
	// The log message.
	message string
	// Optional parameters to be included in the log message.
	params []interface{}
}

// NewLogBuilder creates a new LogBuilder instance with the specified
// log message.
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

// Save builds the log parameters and invokes the Log function
// to save the log message.
func (b *LogBuilder) Save() {
	Log(b.level, b.message, b.params...)
}
