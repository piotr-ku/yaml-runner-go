package system

// This file contains variables and functions used exclusively for testing
// purposes. Exported functions are accessible across packages.

import "bytes"

var testingStdout bytes.Buffer
var testingStderr bytes.Buffer

// GetTestingStdout returns testing stdout buffer
func GetTestingStdout() string {
	return testingStdout.String()
}

// GetTestingStderr returns testing stderr buffer
func GetTestingStderr() string {
	return testingStderr.String()
}
