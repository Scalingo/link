package prober

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Scalingo/go-philae/sampleprobe"
	. "github.com/smartystreets/goconvey/convey"
)

func TestProber(t *testing.T) {
	ctx := context.Background()
	Convey("With healthy probes", t, func() {
		p := NewProber()
		p.AddProbe(sampleprobe.NewSampleProbe("a", true))
		p.AddProbe(sampleprobe.NewSampleProbe("b", true))

		res := p.Check(ctx)

		So(res.Healthy, ShouldBeTrue)
		So(len(res.Probes), ShouldEqual, 2)
		So(validateProbe(res.Probes, "a", true), ShouldBeNil)
		So(validateProbe(res.Probes, "b", true), ShouldBeNil)
	})

	Convey("With unhealthy probes", t, func() {
		p := NewProber()
		p.AddProbe(sampleprobe.NewSampleProbe("a", false))
		p.AddProbe(sampleprobe.NewSampleProbe("b", false))

		res := p.Check(ctx)

		So(res.Healthy, ShouldBeFalse)
		So(len(res.Probes), ShouldEqual, 2)
		So(validateProbe(res.Probes, "a", false), ShouldBeNil)
		So(validateProbe(res.Probes, "b", false), ShouldBeNil)
	})

	Convey("With a healthy probe and a unhealthy probe", t, func() {
		p := NewProber()
		p.AddProbe(sampleprobe.NewSampleProbe("a", true))
		p.AddProbe(sampleprobe.NewSampleProbe("b", false))

		res := p.Check(ctx)

		So(res.Healthy, ShouldBeFalse)
		So(len(res.Probes), ShouldEqual, 2)
		So(validateProbe(res.Probes, "a", true), ShouldBeNil)
		So(validateProbe(res.Probes, "b", false), ShouldBeNil)
	})

	Convey("With a probe that times out", t, func() {
		p := NewProber()
		p.AddProbe(sampleprobe.NewTimedSampleProbe("test", true, 4*time.Second))
		p.AddProbe(sampleprobe.NewTimedSampleProbe("test", true, 4*time.Second))
		p.AddProbe(sampleprobe.NewTimedSampleProbe("test", true, 4*time.Second))
		p.AddProbe(sampleprobe.NewTimedSampleProbe("test", true, 4*time.Second))
		start := time.Now()
		res := p.Check(ctx)
		duration := time.Now().Sub(start)

		So(duration, ShouldBeLessThan, 3*time.Second)

		So(res.Healthy, ShouldBeFalse)
		So(len(res.Probes), ShouldEqual, 4)
		for _, p := range res.Probes {
			So(p.Comment, ShouldEqual, "Probe timeout")
			So(p.Healthy, ShouldBeFalse)
		}
	})
}

func validateProbe(probes []*ProbeResult, name string, healthy bool) error {
	for _, probe := range probes {
		if probe.Name == name {
			So(probe.Healthy, ShouldEqual, healthy)
			return nil
		}
	}

	return errors.New("Unable to find node " + name)
}
