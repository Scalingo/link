package githubprobe

import (
	"encoding/json"
	"io"

	errgo "gopkg.in/errgo.v1"

	"github.com/Scalingo/go-philae/httpprobe"
)

type GithubChecker struct{}

type GithubProbe struct {
	http httpprobe.HTTPProbe
}

func (_ GithubChecker) Check(body io.Reader) error {
	var result GithubStatusResponse

	err := json.NewDecoder(body).Decode(&result)
	if err != nil {
		return errgo.Notef(err, "Invalid json")
	}

	if result.Status == "major" {
		return errgo.Newf("Github is probably down")
	}
	return nil
}

func NewGithubProbe(name string) GithubProbe {
	return GithubProbe{
		http: httpprobe.NewHTTPProbe(name, "https://status.github.com/api/status.json", httpprobe.HTTPOptions{
			Checker: GithubChecker{},
		}),
	}
}

func (p GithubProbe) Name() string {
	return p.http.Name()
}

func (p GithubProbe) Check() error {
	return p.http.Check()
}
