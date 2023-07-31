package system

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestFatalError tests the FatalError function.
//
// It verifies the behavior of the FatalError function by setting up mocked
// dependencies and running various test cases. It checks the return code,
// logs, and the error message printed to stderr.
func TestFatalError(t *testing.T) {
	var rc int
	MockOsExit = func(code int) {
		rc = code
	}
	defer func() {
		MockOsExit = os.Exit
	}()

	const codeIOError = 64
	const codeParseError = 65
	const codeValidationError = 66
	const codeOSError = 67

	tests := []struct {
		name     string
		error    string
		expected int
	}{
		{
			name:     "",
			error:    "",
			expected: 1,
		},
		{
			name:     "UnexistingError",
			error:    "This error does not exists.",
			expected: 1,
		},
		{
			name:     "OK",
			error:    "",
			expected: 0,
		},
		{
			name:     "Unknown",
			error:    "Unknown error",
			expected: 1,
		},
		{
			name:     "IOError",
			error:    "IOError error",
			expected: codeIOError,
		},
		{
			name:     "ParseError",
			error:    "ParseError error",
			expected: codeParseError,
		},
		{
			name:     "ValidationError",
			error:    "ValidationError error",
			expected: codeValidationError,
		},
		{
			name:     "OSError",
			error:    "OSError error",
			expected: codeOSError,
		},
	}

	for _, test := range tests {
		// Set log settings and clear buffers
		LogInit(LogConfig{
			File:  "testing_buffer",
			Level: "info",
			Quiet: false,
			JSON:  false,
		})

		// Call the tested function
		FatalError(test.name, test.error)

		// Test return code
		assert.Equal(t, test.expected, rc)

		// Test logs
		assert.Equal(t, "", GetTestingStdout())
		assert.Regexp(t, fmt.Sprintf(" level=ERROR "+
			"msg=\"FATAL ERROR: %s %s\" ", test.name, test.error),
			GetTestingStderr())
	}
}
