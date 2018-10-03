package logrus_rollbar

import (
	"github.com/pkg/errors"
)

func errorsFoo() error {
	return errors.Wrapf(errorsBar(), "errors from Bar")
}

func errorsBar() error {
	return errors.New("a new error")
}
