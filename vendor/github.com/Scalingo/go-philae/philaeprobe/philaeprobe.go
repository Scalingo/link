package philaeprobe

import (
	"encoding/json"
	"errors"
	"io"

	errgo "gopkg.in/errgo.v1"

	"github.com/Scalingo/go-philae/httpprobe"
	"github.com/Scalingo/go-philae/prober"
)

type PhilaeProbe struct {
	http httpprobe.HTTPProbe
}

type PhilaeChecker struct{}

func (_ PhilaeChecker) Check(body io.Reader) error {
	var result prober.Result

	err := json.NewDecoder(body).Decode(&result)
	if err != nil {
		return errgo.Notef(err, "Invalid json")
	}

	if !result.Healthy {
		reason := ""
		for _, probe := range result.Probes {
			if !probe.Healthy {
				reason += probe.Name + " is down (" + probe.Comment + "),"
			}
		}
		return errors.New(reason)
	}

	return nil
}

func NewPhilaeProbe(name, endpoint string, dialTimeout, responseTimeout int) PhilaeProbe {
	return PhilaeProbe{
		http: httpprobe.NewHTTPProbe(name, endpoint, httpprobe.HTTPOptions{
			Checker:         PhilaeChecker{},
			DialTimeout:     dialTimeout,
			ResponseTimeout: responseTimeout,
		}),
	}
}

func (p PhilaeProbe) Name() string {
	return p.http.Name()
}

func (p PhilaeProbe) Check() error {
	return p.http.Check()
}
