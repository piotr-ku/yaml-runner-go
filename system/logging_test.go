package system

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"
)

const unexpectedResult string = "Unexpected %s. Expected: %v, Got: %v"
const stderrLevel string = "error"

// TestLogTargets contains an unit test for logTargets() function
func TestLogTargets(test *testing.T) {
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

	for _, t := range tests {
		// logging init
		config := LogConfig{File: t.file, Quiet: t.quiet, JSON: false,
			Level: t.level}
		LogInit(config)
		// get targets
		targets := logTargets(t.level)
		// check number of targets
		if !reflect.DeepEqual(targets, t.want) && !(len(targets) == 0 &&
			len(t.want) == 0) {
			test.Errorf(errorFormat, t.level, t.file, t.quiet, targets, t.want)
		}
	}
}

// TestLogTextHandler verifies that Log function generate
// a proper log string when using a text handler
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
		if !strings.Contains(got, expected) {
			t.Errorf(unexpectedResult, "log content", expected, got)
		}
	}
}

// TestLogTextHandler verifies that Log function generate
// a proper log string when using a log handler

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
		if err != nil {
			t.Fatalf("Incorrect JSON log: %s", err)
		}

		// Compare the actual and expected output
		if !reflect.DeepEqual(got, expected) {
			t.Errorf("Decoded JSON does not match the expected output.\n"+
				"Expected: %+v\nGot:   %+v", expected, got)
		}
	}
}

// TestNewLogBuilder verifies that a new LogBuilder instance
// is created correctly.
func TestNewLogBuilder(t *testing.T) {
	message := "Test message"
	builder := NewLogBuilder(message)

	// Verify that a new LogBuilder instance is created
	if builder == nil {
		t.Error("NewLogBuilder should not return nil")
		return
	}

	// Verify that the level is set to "INFO" by default
	if builder.level != "INFO" {
		t.Errorf("Unexpected default level. Expected: %s, Got: %s", "INFO",
			builder.level)
	}

	// Verify that the message is set correctly
	if builder.message != message {
		t.Errorf("Unexpected message. Expected: %s, Got: %s", message,
			builder.message)
	}

	// Verify that the params slice is empty
	if len(builder.params) != 0 {
		t.Error("Unexpected params slice length. Expected: 0, Got:",
			len(builder.params))
	}
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
	if builder.level != level {
		t.Errorf("Unexpected level. Expected: %s, Got: %s", level,
			builder.level)
	}
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
	if len(builder.params) != len(params) {
		t.Errorf("Unexpected number of params. Expected: %d, Got: %d",
			len(params), len(builder.params))
	}

	// Verify that each parameter is added correctly
	for i, param := range params {
		if builder.params[i] != param {
			t.Errorf("Unexpected param at index %d. Expected: %v, Got: %v",
				i, param, builder.params[i])
		}
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
		if !strings.Contains(got, expected) {
			t.Errorf("Unexpected log content. Expected: <%s>, Got: <%s>",
				expected, got)
		}
	}
}
