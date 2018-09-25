package nsqprobe

import (
	"crypto/tls"
	"strconv"

	"github.com/Scalingo/go-philae/httpprobe"
)

type NSQProbe struct {
	http httpprobe.HTTPProbe
}

func NewNSQProbe(name, host string, port int, tlsConfig *tls.Config) NSQProbe {
	scheme := "http"
	if tlsConfig != nil {
		scheme = "https"
	}
	return NSQProbe{
		http: httpprobe.NewHTTPProbe(
			name, scheme+"://"+host+":"+strconv.Itoa(port)+"/ping",
			httpprobe.HTTPOptions{TLSConfig: tlsConfig},
		),
	}
}

func (p NSQProbe) Name() string {
	return p.http.Name()
}

func (p NSQProbe) Check() error {
	return p.http.Check()
}
