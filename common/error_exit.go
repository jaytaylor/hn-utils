package common

import (
	"fmt"
	"os"
)

// ErrorExit prints an error message to STDERR then exits.
// If no exit status-code is specified, it defaults to `1`.
func ErrorExit(err interface{}, code ...int) {
	statusCode := 1
	if len(code) > 0 {
		statusCode = code[0]
	}

	fmt.Fprintf(os.Stderr, "ERROR: %s (exiting with status-code: %v)\n", err, statusCode)

	os.Exit(statusCode)
}
