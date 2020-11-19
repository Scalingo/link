package logrus_rollbar

import (
	"fmt"
	"runtime"
	"strconv"

	"github.com/pkg/errors"
	"github.com/rollbar/rollbar-go"
)

// Rollbar package expect such error:
// type CauseStacker interface {
//   error
//   Cause() error
//   Stack() Stack
// }

var (
	_ rollbar.CauseStacker = wrappedError{}
)

type wrappedError struct {
	err error
	msg string
}

func (err wrappedError) Error() string {
	return err.msg
}

type causer interface {
	Cause() error
}

type stackTracer interface {
	StackTrace() errors.StackTrace
}

func Wrap(msg string, err error) wrappedError {
	return wrappedError{msg: msg, err: err}
}

func (err wrappedError) Cause() error {
	return err.err
}

func (werr wrappedError) Stack() []runtime.Frame {
	stack := []runtime.Frame{}
	err := werr.err

	// We're going to the deepest call
	for {
		c, ok := err.(causer)
		if !ok {
			break
		}
		err = c.Cause()
	}

	// Return an empty stack
	tracer, ok := err.(stackTracer)
	if !ok {
		return stack
	}

	errorsStack := tracer.StackTrace()
	for i := len(errorsStack) - 1; i >= 0; i-- {
		f := errorsStack[i]
		line, _ := strconv.Atoi(fmt.Sprintf("%d", f))
		frame := runtime.Frame{
			File:     fmt.Sprintf("%+s", f),
			Line:     line,
			Function: fmt.Sprintf("%n", f),
		}
		stack = append([]runtime.Frame{frame}, stack...)
	}

	return stack
}
