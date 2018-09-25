package philaehandler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Scalingo/go-philae/prober"
	"github.com/Scalingo/go-philae/sampleprobe"
	. "github.com/smartystreets/goconvey/convey"
)

func TestPhilaeHandler(t *testing.T) {
	TestWorkingEndpoint := func(resp *http.Response) {
		result := &prober.Result{}
		json.NewDecoder(resp.Body).Decode(&result)
		So(result.Healthy, ShouldBeTrue)
		So(len(result.Probes), ShouldEqual, 1)
		So(result.Probes[0].Name, ShouldEqual, "test")
		So(result.Probes[0].Healthy, ShouldBeTrue)
	}

	Convey("With an healthy prober", t, func() {
		probe := prober.NewProber()
		probe.AddProbe(sampleprobe.NewSampleProbe("test", true))
		req := httptest.NewRequest("GET", "http://example.foo/_health", nil)

		Convey("the ServeHTTP method should render a healthy node", func() {
			w := httptest.NewRecorder()
			handler := NewHandler(probe)
			handler.ServeHTTP(w, req)

			resp := w.Result()

			So(resp.StatusCode, ShouldEqual, 200)
			TestWorkingEndpoint(resp)
		})

		Convey("the PhilaeRouter should route to the correct route", func() {
			w := httptest.NewRecorder()
			handler := NewPhilaeRouter(http.NotFoundHandler(), probe)

			handler.ServeHTTP(w, req)

			resp := w.Result()

			So(resp.StatusCode, ShouldEqual, 200)
			TestWorkingEndpoint(resp)
		})
		Convey("the PhilarRouter should route all the other routes to the other router", func() {
			w := httptest.NewRecorder()
			handler := NewPhilaeRouter(http.RedirectHandler("http://scalingo.com", 301), probe)
			req := httptest.NewRequest("GET", "http://example.foo/salut", nil)

			handler.ServeHTTP(w, req)

			resp := w.Result()
			So(resp.StatusCode, ShouldEqual, 301)
		})
	})
}
