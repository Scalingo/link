package philaeprobe

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Scalingo/go-philae/prober"
	. "github.com/smartystreets/goconvey/convey"
)

func TestPhilaeProbe(t *testing.T) {
	Convey("With a unaivalable server", t, func() {
		p := NewPhilaeProbe("http", "http://localhost:6666", 0, 0)
		err := p.Check()
		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldStartWith, "Unable to send request")
	})

	Convey("With a server responding 5XX", t, func() {
		ts := launchTestServer(500, "Error")
		defer ts.Close()

		p := NewPhilaeProbe("http", ts.URL, 0, 0)
		err := p.Check()
		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldStartWith, "Invalid return code")
	})

	Convey("With a server responding 2XX but an invalid json", t, func() {
		ts := launchTestServer(200, "Salut salut")
		defer ts.Close()

		p := NewPhilaeProbe("http", ts.URL, 0, 0)
		err := p.Check()
		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldStartWith, "Invalid json")
	})

	Convey("With a server responding 2XX but an unhealthy probe", t, func() {
		result := &prober.Result{
			Healthy: false,
			Probes: []*prober.ProbeResult{
				&prober.ProbeResult{
					Name:    "node-1",
					Healthy: false,
					Comment: "pas bien",
				},
			},
		}
		ts := launchJSONTestServer(200, result)
		defer ts.Close()

		p := NewPhilaeProbe("http", ts.URL, 0, 0)
		err := p.Check()
		So(err, ShouldNotBeNil)
		So(err.Error(), ShouldEqual, "node-1 is down (pas bien),")
	})

	Convey("With a server responding 2XX and an healthy probe", t, func() {
		result := &prober.Result{
			Healthy: true,
			Probes:  []*prober.ProbeResult{},
		}
		ts := launchJSONTestServer(200, result)
		defer ts.Close()

		p := NewPhilaeProbe("http", ts.URL, 0, 0)
		So(p.Check(), ShouldBeNil)
	})
}

func launchJSONTestServer(statusCode int, response interface{}) *httptest.Server {
	a, err := json.Marshal(response)
	if err != nil {
		panic(err)
	}

	return launchTestServer(statusCode, string(a))
}

func launchTestServer(statusCode int, response string) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(statusCode)
		fmt.Fprintln(w, response)
	}))
}
