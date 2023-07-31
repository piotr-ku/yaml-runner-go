package system

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const stderrLevel string = "error"

// TestLogTargets is a test function for testing the logTargets function.
//
// It tests the logTargets function by providing different test cases and
// comparing the expected output with the actual output. The test cases
// include different levels, file paths, and quiet flags. The function
// initializes the logging configuration, calls the logTargets function,
// and then compares the output with the expected output using the assert
// package. The function ensures that the number of targets returned by
// logTargets matches the expected number of targets. If the number of
// targets is not zero, the function asserts that the returned targets
// match the expected targets.
func TestLogTargets(t *testing.T) {
	// temporary test file
	const testLogFile string = "/tmp/test.log"
	const errorFormat string = "level: %s, file: %s, quiet: %t, " +
		"got: %+v, expected: %+v"

	// tests table
	var tests = []struct {
		description string
		level       string
		file        string
		quiet       bool
		want        []string
	}{
		{level: "debug", file: "", quiet: false,
			want: []string{"stdout"}},
		{level: "info", file: "", quiet: false,
			want: []string{"stdout"}},
		{level: "warn", file: "", quiet: false,
			want: []string{"stdout"}},
		{level: "error", file: "", quiet: false,
			want: []string{"stderr"}},
		{level: "incorrect", file: "", quiet: false,
			want: []string{"stdout"}},
		{level: "debug", file: "", quiet: true,
			want: []string{}},
		{level: "info", file: "", quiet: true,
			want: []string{}},
		{level: "warn", file: "", quiet: true,
			want: []string{}},
		{level: "error", file: "", quiet: true,
			want: []string{}},
		{level: "incorrect", file: "", quiet: true,
			want: []string{}},
		{level: "debug", file: testLogFile, quiet: false,
			want: []string{"stdout", "file"}},
		{level: "info", file: testLogFile, quiet: false,
			want: []string{"stdout", "file"}},
		{level: "warn", file: testLogFile, quiet: false,
			want: []string{"stdout", "file"}},
		{level: "error", file: testLogFile, quiet: false,
			want: []string{"stderr", "file"}},
		{level: "incorrect", file: testLogFile, quiet: false,
			want: []string{"stdout", "file"}},
		{level: "debug", file: testLogFile, quiet: true,
			want: []string{"file"}},
		{level: "info", file: testLogFile, quiet: true,
			want: []string{"file"}},
		{level: "warn", file: testLogFile, quiet: true,
			want: []string{"file"}},
		{level: "error", file: testLogFile, quiet: true,
			want: []string{"file"}},
		{level: "incorrect", file: testLogFile, quiet: true,
			want: []string{"file"}},
	}

	for _, test := range tests {
		// logging init
		config := LogConfig{File: test.file, Quiet: test.quiet, JSON: false,
			Level: test.level}
		LogInit(config)
		// get targets
		targets := logTargets(test.level)
		// check number of targets
		if !(len(targets) == 0 && len(test.want) == 0) {
			assert.Equal(t, test.want, targets, errorFormat, test.level,
				test.file, test.quiet, targets, test.want)
		} else {
			assert.NotNil(t, len(targets))
			assert.NotNil(t, len(test.want))
		}
	}
}

// TestLogTextHandler is a test function for testing a log text handler.
//
// It iterates over the log levels and performs the following steps:
//   - Initializes log buffering with text format.
//   - Logs a test message with a specific set of fields.
//   - Formats the expected output based on the log level and the test message.
//   - Retrieves the actual output from the testing buffer.
//   - Checks if the actual output contains the expected output.
func TestLogTextHandler(t *testing.T) {
	for _, level := range []string{"debug", "info", "warn", "error"} {
		// log buffering, text format
		LogInit(LogConfig{File: "testing_buffer", Quiet: false, JSON: false})

		// log
		Log(level, "logging test", "field1", "preface-flinch-suspense")

		// expected output
		format := "level=%s msg=\"%s\" field1=%s"
		expected := fmt.Sprintf(format, strings.ToUpper(level),
			"logging test", "preface-flinch-suspense")

		// got
		got := testingStdout.String()
		if level == stderrLevel {
			got = testingStderr.String()
		}

		// content test
		assert.Contains(t, got, expected)
	}
}

// TestLogJSONHandler is a test function for testing a log JSON handler.
//
// This function iterates through the log levels ["debug", "info", "warn",
// "error"] and checks if the logging output is in the expected JSON format.
// It initializes the log buffer in JSON format, logs a message with
// the specified log level, and compares the actual output with
// the expected output.
func TestLogJSONHandler(t *testing.T) {
	for _, level := range []string{"debug", "info", "warn", "error"} {
		// output JSON format
		type Output struct {
			Level   string `json:"level"`
			Message string `json:"msg"`
			Field   string `json:"field1"`
		}

		// log buffering, JSON format
		LogInit(LogConfig{File: "testing_buffer", Quiet: false, JSON: true})

		// log
		Log(level, "TestLogJSONHandler", "field1", "zen-snagged-travel")

		// expected
		expected := Output{
			Level:   strings.ToUpper(level),
			Message: "TestLogJSONHandler",
			Field:   "zen-snagged-travel",
		}

		// choose stdout/stderr buffer
		output := testingStdout.Bytes()
		if level == stderrLevel {
			output = testingStderr.Bytes()
		}

		// decode JSON
		var got Output
		err := json.Unmarshal(output, &got)
		assert.Nil(t, err)

		// Compare the actual and expected output
		assert.Equal(t, expected, got)
	}
}

// TestNewLogBuilder is a test function for testing the NewLogBuilder function.
//
// It creates a new LogBuilder instance, sets the level to "INFO", and
// verifies that the instance is created correctly.
func TestNewLogBuilder(t *testing.T) {
	message := "Test message"
	builder := NewLogBuilder(message)

	// Verify that a new LogBuilder instance is created
	if !assert.NotNil(t, builder) {
		return
	}

	// Verify that the level is set to "INFO" by default
	assert.Equal(t, "INFO", builder.level)

	// Verify that the message is set correctly
	assert.Equal(t, message, builder.message)

	// Verify that the params slice is empty
	assert.Equal(t, 0, len(builder.params))
}

// TestLevel verifies that the Level method sets the level correctly
// in the LogBuilder instance.
func TestLevel(t *testing.T) {
	message := "Test message"
	level := "DEBUG"
	builder := NewLogBuilder(message)

	// Call the Level method
	builder.Level(level)

	// Verify that the level is set correctly
	assert.Equal(t, level, builder.level)
}

// TestSet verifies that the Set method adds parameters correctly to
// the LogBuilder instance.
func TestSet(t *testing.T) {
	message := "Test message"
	params := []interface{}{"param1", "foothill-constrict-employee"}
	builder := NewLogBuilder(message)

	// Call the Set method to add parameters
	builder.Set(params...)

	// Verify that the params slice is updated correctly
	assert.Equal(t, len(params), len(builder.params))

	// Verify that each parameter is added correctly
	for i, param := range params {
		assert.Equal(t, param, builder.params[i],
			"Unexpected param at index %d. Expected: %v, Got: %v",
			i, param, builder.params[i])
	}
}

// TestSave verifies that the Set method correctly run Log function
func TestSave(t *testing.T) {
	for _, level := range []string{"debug", "info", "warn", "error"} {
		// log buffering, text format
		LogInit(LogConfig{File: "testing_buffer", Quiet: false, JSON: false})

		// Log
		message := "audio-yen-suing"
		params := []interface{}{"field1", "gratified-print-unsmooth"}
		builder := NewLogBuilder(message)

		// Call the Set method to add parameters
		builder.Level(level)
		builder.Set(params...)
		builder.Save()

		// expected output
		expected := fmt.Sprintf("level=%s msg=%s field1=%s",
			strings.ToUpper(level), message, "gratified-print-unsmooth")

		// got
		got := testingStdout.String()

		if level == stderrLevel {
			got = testingStderr.String()
		}

		// content test
		assert.Contains(t, got, expected)
	}
}

// TestLogIncorrectLevel verifies that Log function will panic when incorrect
// log level is passed.
func TestLogIncorrectLevel(t *testing.T) {
	LogInit(LogConfig{File: "testing_buffer", Quiet: false, JSON: false})
	assert.Panics(t, func() { Log("incorrect_level", "msg") })
}

func TestLogInitInvalidFilePath(t *testing.T) {
	var rc int
	const codeIOError = 64
	MockOsExit = func(code int) {
		rc = code
	}
	defer func() {
		MockOsExit = os.Exit
	}()

	LogInit(LogConfig{File: "/not/existing/file", Quiet: true, JSON: false})
	assert.Equal(t, codeIOError, rc)
}
