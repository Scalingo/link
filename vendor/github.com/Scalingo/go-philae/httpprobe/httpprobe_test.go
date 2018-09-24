package httpprobe

import (
	"errors"
	"testing"

	httpmock "gopkg.in/jarcoal/httpmock.v1"

	. "github.com/smartystreets/goconvey/convey"
)

func TestHttpProbe(t *testing.T) {
	Convey("With a unaivalable server", t, func() {
		p := NewHTTPProbe("http", "http://localhost:6666", HTTPOptions{testing: true})
		err := p.Check()
		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldStartWith, "Unable to send request")
	})

	Convey("With a server responding 5XX", t, func() {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()
		httpmock.RegisterResponder("GET", "http://scalingo.com/",
			httpmock.NewStringResponder(500, `Error`))

		Convey("With a normal node, it should refuse this node", func() {
			p := NewHTTPProbe("http", "http://scalingo.com/", HTTPOptions{testing: true})

			err := p.Check()
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldStartWith, "Invalid return code")
		})

		Convey("With a node accepting 500 responses, it should accept this node", func() {
			p := NewHTTPProbe("http", "http://scalingo.com/", HTTPOptions{
				testing:            true,
				ExpectedStatusCode: 500,
			})
			err := p.Check()
			So(err, ShouldBeNil)
		})

		Convey("With a node accepting 400 responses, it should not accept this node", func() {
			p := NewHTTPProbe("http", "http://scalingo.com/", HTTPOptions{
				testing:            true,
				ExpectedStatusCode: 400,
			})
			err := p.Check()
			So(err, ShouldNotBeNil)
			So(err.Error(), ShouldEqual, "Unexpected status code: 500 (expected: 400)")
		})

	})
	Convey("With a server responding 2XX", t, func() {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()
		httpmock.RegisterResponder("GET", "http://scalingo.com/",
			httpmock.NewStringResponder(200, `Error`))

		Convey("With a vanilla probe", func() {
			p := NewHTTPProbe("http", "http://scalingo.com/", HTTPOptions{testing: true})
			So(p.Check(), ShouldBeNil)
		})
		Convey("With a probe with a custom checker", func() {
			Convey("With a success checker", func() {
				p := NewHTTPProbe("http", "http://scalingo.com/", HTTPOptions{
					testing: true,
					Checker: newTestChecker(nil),
				})
				So(p.Check(), ShouldBeNil)
			})

			Convey("With a failing checker", func() {
				p := NewHTTPProbe("http", "http://scalingo.com/", HTTPOptions{
					testing: true,
					Checker: newTestChecker(errors.New("test error")),
				})
				err := p.Check()

				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "test error")
			})
		})
	})
}
