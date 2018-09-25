package httpprobe

import "io"

type HTTPChecker interface {
	Check(io.Reader) error
}

type testChecker struct {
	err error
}

func newTestChecker(err error) testChecker {
	return testChecker{
		err: err,
	}
}

func (t testChecker) Check(_ io.Reader) error {
	return t.err
}
