package system

import (
	"fmt"
	"os"
	"runtime"
)

var returnCodes = map[string]int{
	"OK":              0,
	"Unknown":         1,
	"IOError":         64,
	"ParseError":      65,
	"ValidationError": 66,
	"OSError":         67,
}
var MockOsExit = os.Exit

// FatalError tries to write a log error and exist with the status code
func FatalError(name string, error string) {
	// Get runtime info
	pc, filename, line, _ := runtime.Caller(1)
	fn := runtime.FuncForPC(pc).Name()

	// Save logs
	Log("error", fmt.Sprintf("FATAL ERROR: %s %s", name, error), "file",
		filename, "line", line, "fn", fn)

	// Get return code number
	code, exists := returnCodes[name]
	if !exists {
		code = 1
	}

	// Exit
	MockOsExit(code)
}
