package system

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"
)

// TestLogTargets contains an unit test for logTargets() function
func TestLogTargets(test *testing.T) {
	// temporary test file
	TestLogFile := "/tmp/test.log"

	// tests table
	var tests = []struct {
		description string
		level       string
		file        string
		quiet       bool
		want        []string
	}{
		{level: "debug", file: "", quiet: false, want: []string{"stdout"}},
		{level: "info", file: "", quiet: false, want: []string{"stdout"}},
		{level: "warn", file: "", quiet: false, want: []string{"stdout"}},
		{level: "error", file: "", quiet: false, want: []string{"stderr"}},
		{level: "incorrect", file: "", quiet: false, want: []string{"stdout"}},
		{level: "debug", file: "", quiet: true, want: []string{}},
		{level: "info", file: "", quiet: true, want: []string{}},
		{level: "warn", file: "", quiet: true, want: []string{}},
		{level: "error", file: "", quiet: true, want: []string{}},
		{level: "incorrect", file: "", quiet: true, want: []string{}},
		{level: "debug", file: TestLogFile, quiet: false, want: []string{"stdout", "file"}},
		{level: "info", file: TestLogFile, quiet: false, want: []string{"stdout", "file"}},
		{level: "warn", file: TestLogFile, quiet: false, want: []string{"stdout", "file"}},
		{level: "error", file: TestLogFile, quiet: false, want: []string{"stderr", "file"}},
		{level: "incorrect", file: TestLogFile, quiet: false, want: []string{"stdout", "file"}},
		{level: "debug", file: TestLogFile, quiet: true, want: []string{"file"}},
		{level: "info", file: TestLogFile, quiet: true, want: []string{"file"}},
		{level: "warn", file: TestLogFile, quiet: true, want: []string{"file"}},
		{level: "error", file: TestLogFile, quiet: true, want: []string{"file"}},
		{level: "incorrect", file: TestLogFile, quiet: true, want: []string{"file"}},
	}

	for _, t := range tests {
		// logging init
		LogInit(LogConfig{File: t.file, Quiet: t.quiet, Json: false, Level: t.level})
		// get targets
		targets := logTargets(t.level)
		// check number of targets
		if !reflect.DeepEqual(targets, t.want) && !(len(targets) == 0 && len(t.want) == 0) {
			test.Errorf("level: %s, file: %s, quiet: %t, got: %+v, expected: %+v", t.level, t.file, t.quiet, targets, t.want)
		}
	}
}

// TestLogTextHandler verifies that Log function generate a proper log string when using a text handler
func TestLogTextHandler(t *testing.T) {
	for _, level := range []string{"debug", "info", "warn", "error"} {
		// log buffering, text format
		LogInit(LogConfig{File: "testing_buffer", Quiet: false, Json: false})

		// log
		Log(level, "logging test", "field1", 41)

		// expected output
		expected := fmt.Sprintf("level=%s msg=\"%s\" field1=%d", strings.ToUpper(level), "logging test", 41)

		// got
		got := testingStdout.String()
		if level == "error" {
			got = testingStderr.String()
		}

		// content test
		if !strings.Contains(got, expected) {
			t.Errorf("Unexpected log content. Expected: <%s>, Got: <%s>", expected, got)
		}
	}
}

// TestLogTextHandler verifies that Log function generate a proper log string when using a log handler

func TestLogJSONHandler(t *testing.T) {
	for _, level := range []string{"debug", "info", "warn", "error"} {
		// output JSON format
		type Output struct {
			Level   string `json:"level"`
			Message string `json:"msg"`
			Field   int    `json:"field1"`
		}

		// log buffering, JSON format
		LogInit(LogConfig{File: "testing_buffer", Quiet: false, Json: true})

		// log
		Log(level, "logging test", "field1", 41)

		// expected
		expected := Output{
			Level:   strings.ToUpper(level),
			Message: "logging test",
			Field:   41,
		}

		// choose stdout/stderr buffer
		output := testingStdout.Bytes()
		if level == "error" {
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
			t.Errorf("Decoded JSON does not match the expected output.\nExpected: %+v\nGot:   %+v", expected, got)
		}
	}
}

// TestNewLogBuilder verifies that a new LogBuilder instance is created correctly.
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
		t.Errorf("Unexpected default level. Expected: %s, Got: %s", "INFO", builder.level)
	}

	// Verify that the message is set correctly
	if builder.message != message {
		t.Errorf("Unexpected message. Expected: %s, Got: %s", message, builder.message)
	}

	// Verify that the params slice is empty
	if len(builder.params) != 0 {
		t.Error("Unexpected params slice length. Expected: 0, Got:", len(builder.params))
	}
}

// TestLevel verifies that the Level method sets the level correctly in the LogBuilder instance.
func TestLevel(t *testing.T) {
	message := "Test message"
	level := "DEBUG"
	builder := NewLogBuilder(message)

	// Call the Level method
	builder.Level(level)

	// Verify that the level is set correctly
	if builder.level != level {
		t.Errorf("Unexpected level. Expected: %s, Got: %s", level, builder.level)
	}
}

// TestSet verifies that the Set method adds parameters correctly to the LogBuilder instance.
func TestSet(t *testing.T) {
	message := "Test message"
	params := []interface{}{"param1", 42}
	builder := NewLogBuilder(message)

	// Call the Set method to add parameters
	builder.Set(params...)

	// Verify that the params slice is updated correctly
	if len(builder.params) != len(params) {
		t.Errorf("Unexpected number of params. Expected: %d, Got: %d", len(params), len(builder.params))
	}

	// Verify that each parameter is added correctly
	for i, param := range params {
		if builder.params[i] != param {
			t.Errorf("Unexpected param at index %d. Expected: %v, Got: %v", i, param, builder.params[i])
		}
	}
}

// TestSave verifies that the Set method correctly run Log function
func TestSave(t *testing.T) {
	for _, level := range []string{"debug", "info", "warn", "error"} {
		// log buffering, text format
		LogInit(LogConfig{File: "testing_buffer", Quiet: false, Json: false})

		// Log
		message := "Test message"
		params := []interface{}{"field1", 42}
		builder := NewLogBuilder(message)

		// Call the Set method to add parameters
		builder.Level(level)
		builder.Set(params...)
		builder.Save()

		// expected output
		expected := fmt.Sprintf("level=%s msg=\"%s\" field1=%d", strings.ToUpper(level), message, 42)

		// got
		got := testingStdout.String()

		if level == "error" {
			got = testingStderr.String()
		}

		// content test
		if !strings.Contains(got, expected) {
			t.Errorf("Unexpected log content. Expected: <%s>, Got: <%s>", expected, got)
		}
	}
}
