package tcpprobe

import (
	"net"
	"time"

	errgo "gopkg.in/errgo.v1"
)

type TCPProbe struct {
	name     string
	endpoint string
	options  TCPOptions
}

type TCPOptions struct {
	Timeout time.Duration
}

func NewTCPProbe(name, endpoint string, opts TCPOptions) TCPProbe {
	return TCPProbe{
		name:     name,
		endpoint: endpoint,
		options:  opts,
	}
}

func (p TCPProbe) Name() string {
	return p.name
}

func (p TCPProbe) Check() error {
	dialer := &net.Dialer{}
	timeout := p.options.Timeout
	if timeout == 0 {
		timeout = 5 * time.Second
	}
	dialer.Timeout = timeout

	a, err := dialer.Dial("tcp", p.endpoint)
	if err != nil {
		return errgo.Notef(err, "fail to contact endpoint")
	}

	a.Close()
	return nil
}
