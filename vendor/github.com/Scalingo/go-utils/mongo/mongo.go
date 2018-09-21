package mongo

import (
	"context"
	"crypto/tls"
	"errors"
	"net"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/Scalingo/go-utils/logger"
	"github.com/sirupsen/logrus"
	"gopkg.in/mgo.v2"
)

var (
	DefaultDatabaseName string
	sessionOnce         = sync.Once{}
	_session            *mgo.Session
)

func Session(log logrus.FieldLogger) *mgo.Session {
	sessionOnce.Do(func() {
		log = log.WithField("process", "mongo-init")
		err := errors.New("")
		for err != nil {
			_session, err = buildSession(logger.ToCtx(context.Background(), log))
			if err != nil {
				log.WithField("err", err).WithField("action", "wait 10sec").Info("init mongo: fail to create session")
				time.Sleep(10 * time.Second)
			}
		}
	})
	return _session
}

func buildSession(ctx context.Context) (*mgo.Session, error) {
	log := logger.Get(ctx)
	rawURL := os.Getenv("MONGO_URL")
	if rawURL == "" {
		rawURL = "mongodb://localhost:27017/" + DefaultDatabaseName
	}

	u, err := url.Parse(rawURL)
	if err != nil {
		return nil, errors.New("not a valid MONGO_URL")
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

	log.WithField("mongodb_host", u.Host).Info("init mongo")
	s, err := mgo.DialWithInfo(info)
	if err != nil {
		return nil, err
	}
	return s, nil
}
