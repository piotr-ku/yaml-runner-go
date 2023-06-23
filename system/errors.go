package system

import (
	"fmt"
	"os"
	"runtime"
)

// FatalError tries to write a log error and exist with the status code
// notest
func FatalError(name string, error string) {
	// Get runtime info
	pc, filename, line, _ := runtime.Caller(1)
	fn := runtime.FuncForPC(pc).Name()

	// Save logs
	Log("error", fmt.Sprintf("FATAL ERROR: %s %s", name, error), "file",
		filename, "line", line, "fn", fn)

	// Exit
	os.Exit(returnCode(name)) // nolint:revive
}

// returnCode returns a code number to return from defined errors
// If the error is not defined, it returns 1
func returnCode(name string) int {
	// return codes constans
	const codeIOError int = 64
	const codeParseError int = 65
	const codeValidationError int = 66
	const codeLogLevelError int = 67
	const codeOSError int = 68
	// codes defines possible application return codes
	codes := map[string]int{
		"OK":              0,
		"Unknown":         1,
		"IOError":         codeIOError,
		"ParseError":      codeParseError,
		"ValidationError": codeValidationError,
		"LogLevelError":   codeLogLevelError,
		"OSError":         codeOSError,
	}

	// Get return code number
	code, exists := codes[name]
	if !exists {
		code = 1
	}

	return code
}
