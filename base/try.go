package base

import (
	"runtime"
	//"runtime/debug"
)

const StackStringSize = 16384

// Attempts to run `attempt` and recovers from any panics, returning the panic object or nil if success.
func Try(attempt func()) (panicked interface{}, stackTrace string) {
	defer func() {
		if o := recover(); o != nil {
			panicked = o

			// Get stack trace:
			stkBytes := make([]byte, StackStringSize, StackStringSize)
			n := runtime.Stack(stkBytes, false)
			stackTrace = string(stkBytes[:n])

			return
		}
	}()

	attempt()

	panicked = nil
	stackTrace = ""
	return
}

func PanicIf(err error) {
	if err != nil {
		panic(err)
	}
}
