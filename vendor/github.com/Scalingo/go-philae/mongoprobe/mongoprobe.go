package mongoprobe

import (
	"crypto/tls"
	"errors"
	"net"
	"net/url"
	"strings"

	errgo "gopkg.in/errgo.v1"
	mgo "gopkg.in/mgo.v2"
)

type MongoProbe struct {
	name string
	url  string
}

func NewMongoProbe(name, url string) MongoProbe {
	return MongoProbe{
		name: name,
		url:  url,
	}
}

func (p MongoProbe) Name() string {
	return p.name
}

func (p MongoProbe) Check() error {
	session, err := p.buildSession(p.url)
	if err != nil {
		return errgo.Notef(err, "Unable to contact server")
	}
	defer session.Close()

	err = session.Ping()
	if err != nil {
		return errgo.Notef(err, "Unable to send query")
	}
	return nil
}

func (p MongoProbe) buildSession(rawURL string) (*mgo.Session, error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, errors.New("not a valid URL")
	}

	withTLS := false
	if u.Query().Get("ssl") == "true" {
		withTLS = true
		rawURL = strings.Replace(rawURL, "?ssl=true", "", 1)
		rawURL = strings.Replace(rawURL, "&ssl=true", "", 1)
	}
	info, err := mgo.ParseURL(rawURL)
	if err != nil {
		return nil, err
	}
	if withTLS {
		tlsConfig := &tls.Config{}
		tlsConfig.InsecureSkipVerify = true
		info.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
			conn, err := tls.Dial("tcp", addr.String(), tlsConfig)
			return conn, err
		}
	}

	s, err := mgo.DialWithInfo(info)
	if err != nil {
		return nil, err
	}
	return s, nil
}
