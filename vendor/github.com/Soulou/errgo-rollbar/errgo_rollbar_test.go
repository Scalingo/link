package errgorollbar

import (
	"errors"
	"os"
	"testing"

	"github.com/rollbar/rollbar-go"
	"gopkg.in/errgo.v1"
)

func init() {
	rollbar.SetEnvironment("test")
	rollbar.SetToken(os.Getenv("TOKEN"))
}

func a() error {
	return errgo.New("error")
}

func b() error {
	return errgo.Mask(a())
}

func c() error {
	return errgo.Mask(b(), errgo.Any)
}

func note() error {
	return errgo.Notef(c(), "error of %s", "c")
}

func TestBuildStack(t *testing.T) {
	werr := Wrap(note().Error(), note(), 0)
	stack := werr.Stack()
	// Stack: a() -> b() -> c() -> note() -> TestBuildStack() -> testing.TestRunner() -> runtime.main()
	if len(stack) != 7 {
		t.Errorf("len(stack) = %d != 7", len(stack))
	}

	werr = Wrap(c().Error(), c(), 0)
	stack = werr.Stack()
	if len(stack) != 6 {
		t.Errorf("len(stack) = %d != 6", len(stack))
	}

	werr = Wrap("test empty", nil, 0)
	stack = werr.Stack()
	if len(stack) != 3 {
		t.Errorf("len(stack) = %d != 3", len(stack))
	}

	werr = Wrap("error", errors.New("error"), 0)
	stack = werr.Stack()
	if len(stack) != 3 {
		t.Errorf("len(stack) = %d != 3", len(stack))
	}

	err := Wrap(c().Error(), c(), 0)
	rollbar.Error(rollbar.ERR, err)
	rollbar.Wait()
}
