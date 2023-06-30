package system

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestGetTestingStdout verifies that GetTestingStdout returns
// the expected value.
func TestGetTestingStdout(t *testing.T) {
	expected := "test stdout"

	// Write the expected value to the testingStdout buffer
	testingStdout.Reset()
	length, err := testingStdout.WriteString(expected)

	// Assert that the number of bytes written is equal to the length
	// of the expected value
	assert.Equal(t, len(expected), length)

	// Assert that no error occurred during the write operation
	assert.Nil(t, err)

	// Assert that the value returned by GetTestingStdout is equal
	// to the expected value
	assert.Equal(t, expected, GetTestingStdout())
}

// TestGetTestingStderr verifies that GetTestingStderr returns
// the expected value.
func TestGetTestingStderr(t *testing.T) {
	expected := "test stderr"

	// Write the expected value to the testingStderr buffer
	testingStderr.Reset()
	length, err := testingStderr.WriteString(expected)

	// Assert that the number of bytes written is equal to the length
	// of the expected value
	assert.Equal(t, len(expected), length)

	// Assert that no error occurred during the write operation
	assert.Nil(t, err)

	// Assert that the value returned by GetTestingStderr is equal
	// to the expected value
	assert.Equal(t, expected, GetTestingStderr())
}
