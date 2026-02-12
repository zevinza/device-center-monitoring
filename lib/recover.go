package lib

import (
	"errors"
	"fmt"
	"log"
	"runtime"
)

// Recover errors
func Recover() {
	if err := recover(); nil != err {
		PrintStackTrace(err)
	}
}

// StackTrace return error stack
func StackTrace(e interface{}) error {
	stack := make([]byte, 1024*8)
	stack = stack[:runtime.Stack(stack, false)]
	stackTraces := fmt.Sprintf("panic: %v\n%s\n", e, stack)

	return errors.New(stackTraces)
}

// PrintStackTrace print stack trace
func PrintStackTrace(e interface{}) {
	log.Println(StackTrace(e))
}
