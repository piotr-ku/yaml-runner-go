package system

import (
	"testing"
)

// TestReturnCode contains an unit test for returnCode function
func TestReturnCode(t *testing.T) {
	var tests = []struct {
		description string
		input       string
		want        int
	}{
		{"OK should be 0", "OK", 0},
		{"Unknown should be 1", "Unknown", 1},
		{"NotExisted should be 1", "NotExisted", 1},
		{"IOError should be 64", "IOError", 64},
		{"ParseError should be 65", "ParseError", 65},
	}

	for _, tt := range tests {
		result := returnCode(tt.input)
		if result != tt.want {
			t.Errorf(tt.description, result, tt.want)
		}
	}
}
