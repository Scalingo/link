package errgorollbar

import (
	"runtime"
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
	_ rollbar.Stacker = wrappedError{}
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

func (werr wrappedError) Stack() []runtime.Frame {
	stack := []runtime.Frame{}
	err := werr.err
	for err != nil {
		if errgoErr, ok := err.(*errgo.Err); !ok {
			break
		} else {
			frame := runtime.Frame{
				File: errgoErr.File,
				Line: errgoErr.Line,
			}
			stack = append([]runtime.Frame{frame}, stack...)
			err = errgoErr.Underlying()
		}
	}

	// The result is the concatenation of the stack given by the current
	// execution flow and the stack determined by the pkg/errors error
	// Plus we ignore all the intermediate frames from rollbar lib which
	// is adding something around 4 levels of depth.
	rawStack := getCallersFrames(2 + werr.skip)
	execStack := []runtime.Frame{}
	for _, frame := range rawStack {
		if !strings.Contains(frame.File, "rollbar/rollbar-go") {
			execStack = append(execStack, frame)
		}
	}
	return append(stack, execStack...)
}

func getCallersFrames(skip int) []runtime.Frame {
	pc := make([]uintptr, 100)
	runtime.Callers(1+skip, pc)
	fr := runtime.CallersFrames(pc)

	return framesToSlice(fr)
}

// framesToSlice extracts all the runtime.Frame from runtime.Frames.
func framesToSlice(fr *runtime.Frames) []runtime.Frame {
	frames := make([]runtime.Frame, 0)

	for frame, more := fr.Next(); frame != (runtime.Frame{}); frame, more = fr.Next() {
		frames = append(frames, frame)

		if !more {
			break
		}
	}

	return frames
}
