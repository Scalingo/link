package redisprobe

import (
	"net/url"

	errgo "gopkg.in/errgo.v1"
	redis "gopkg.in/redis.v4"
)

type RedisProbe struct {
	name     string
	password string
	host     string
}

func NewRedisProbeFromURL(name, serviceUrl string) RedisProbe {
	url, _ := url.Parse(serviceUrl)
	password := ""

	if url.User != nil {
		password, _ = url.User.Password()
	}

	return NewRedisProbe(name, url.Host, password)
}

func NewRedisProbe(name, host, password string) RedisProbe {
	return RedisProbe{
		name:     name,
		password: password,
		host:     host,
	}
}

func (p RedisProbe) Name() string {
	return p.name
}

func (p RedisProbe) Check() error {
	client := redis.NewClient(&redis.Options{
		Addr:     p.host,
		Password: p.password,
		DB:       0,
	})
	defer client.Close()

	_, err := client.Ping().Result()

	if err != nil {
		return errgo.Notef(err, "Unable to contact host")
	}

	return nil
}
