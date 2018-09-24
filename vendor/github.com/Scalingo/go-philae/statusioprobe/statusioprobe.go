package statusioprobe

import (
	"encoding/json"
	"io"

	errgo "gopkg.in/errgo.v1"

	"github.com/Scalingo/go-philae/httpprobe"
)

type StatusIOProbe struct {
	http httpprobe.HTTPProbe
}

type StatusIOChecker struct{}

func NewStatusIOProbe(name string, id string) StatusIOProbe {
	return StatusIOProbe{
		http: httpprobe.NewHTTPProbe(name, "https://api.status.io/1.0/status/"+id, httpprobe.HTTPOptions{
			Checker: StatusIOChecker{},
		}),
	}
}

func (p StatusIOProbe) Name() string {
	return p.http.Name()
}

func (p StatusIOProbe) Check() error {
	return p.http.Check()
}

func (_ StatusIOChecker) Check(body io.Reader) error {
	var result StatusIOResponse

	err := json.NewDecoder(body).Decode(&result)
	if err != nil {
		return errgo.Notef(err, "Invalid json")
	}

	if result.Result.Overall.StatusCode >= 400 {
		return errgo.Newf("Incident reported!")
	}

	return nil
}
