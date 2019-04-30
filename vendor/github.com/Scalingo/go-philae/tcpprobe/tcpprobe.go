package tcpprobe

import (
	"context"
	"fmt"
	"net"
	"time"

	errgo "gopkg.in/errgo.v1"
)

type Resolver interface {
	LookupIPAddr(context.Context, string) ([]net.IPAddr, error)
}

type Dialer interface {
	DialContext(context.Context, string, string) (net.Conn, error)
}

type TCPProbe struct {
	name     string
	endpoint string
	options  TCPOptions
}

type TCPOptions struct {
	Timeout  time.Duration
	Resolver Resolver
	Dialer   Dialer
}

func NewTCPProbe(name, endpoint string, opts TCPOptions) TCPProbe {
	if opts.Resolver == nil {
		opts.Resolver = net.DefaultResolver
	}
	if opts.Dialer == nil {
		opts.Dialer = &net.Dialer{}
	}
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
	timeout := p.options.Timeout
	if timeout == 0 {
		timeout = 5 * time.Second
	}

	host, port, err := net.SplitHostPort(p.endpoint)
	if err != nil {
		return errgo.Notef(err, "invalid endpoint %v, should be host:port", p.endpoint)
	}

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	addrs, err := p.options.Resolver.LookupIPAddr(ctx, host)
	if err != nil {
		return errgo.Notef(err, "DNS resolution failed for %v", host)
	}

	endpoint := addrs[0].String() + ":" + port
	if len(addrs[0].IP) == net.IPv6len {
		endpoint = fmt.Sprintf("[%s]:%s", addrs[0].String(), port)
	}

	a, err := p.options.Dialer.DialContext(ctx, "tcp", endpoint)
	if err != nil {
		return errgo.Notef(err, "fail to open TCP connection")
	}

	err = a.Close()
	if err != nil {
		return errgo.Notef(err, "fail to close connection")
	}
	return nil
}
