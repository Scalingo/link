package sampleprobe

import (
	"time"

	errgo "gopkg.in/errgo.v1"
)

// Used for tests only

type SampleProbe struct {
	name   string
	result bool
	time   time.Duration
}

func NewSampleProbe(name string, result bool) SampleProbe {
	return SampleProbe{
		name:   name,
		result: result,
		time:   1 * time.Millisecond,
	}
}

func NewTimedSampleProbe(name string, result bool, time time.Duration) SampleProbe {
	return SampleProbe{
		name:   name,
		result: result,
		time:   time,
	}
}

func (s SampleProbe) Name() string {
	return s.name
}

func (s SampleProbe) Check() error {
	time.Sleep(s.time)
	if s.result {
		return nil
	}
	return errgo.New("error")
}
