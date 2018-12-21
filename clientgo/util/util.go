package util

import (
	"fmt"
	"runtime"
)

const crashMessage = `panic: %v

%s

Oh noes! process crashed!

Please submit the stack trace and any relevant information to:
administrator`

func MakePanicTrace(err interface{}) string {
	stackBuf := make([]byte, 4096)
	n := runtime.Stack(stackBuf, false)
	return fmt.Sprintf(crashMessage, err, stackBuf[:n])
}
