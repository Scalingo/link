package errgorollbar

import (
	"strings"

	"github.com/rollbar/rollbar-go"
	"gopkg.in/errgo.v1"
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
	err  error
	msg  string
	skip int
}

func (err wrappedError) Error() string {
	return err.msg
}

func Wrap(msg string, err error, skip int) wrappedError {
	return wrappedError{msg: msg, err: err, skip: skip}
}

func (err wrappedError) Cause() error {
	return err.err
}

func (werr wrappedError) Stack() rollbar.Stack {
	stack := rollbar.Stack{}
	err := werr.err
	for err != nil {
		if errgoErr, ok := err.(*errgo.Err); !ok {
			break
		} else {
			frame := rollbar.Frame{
				Filename: errgoErr.File,
				Line:     errgoErr.Line,
			}
			stack = append([]rollbar.Frame{frame}, stack...)
			err = errgoErr.Underlying()
		}
	}

	// The result is the concatenation of the stack given by the current
	// execution flow and the stack determined by the pkg/errors error
	// Plus we ignore all the intermediate frames from rollbar lib which
	// is adding something around 4 levels of depth.
	rawStack := rollbar.BuildStack(werr.skip + 2)
	execStack := rollbar.Stack{}
	for _, frame := range rawStack {
		if !strings.Contains(frame.Filename, "rollbar/rollbar-go") {
			execStack = append(execStack, frame)
		}
	}
	return append(stack, execStack...)
}
