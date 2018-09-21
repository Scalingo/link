package nsqproducer

import (
	"strconv"

	"github.com/Scalingo/go-utils/env"
	nsq "github.com/nsqio/go-nsq"
	"github.com/pkg/errors"
)

func FromEnv() (*NsqProducer, error) {
	E := env.InitMapFromEnv(map[string]string{
		"NSQD_TLS":        "false",
		"NSQD_TLS_CACERT": "",
		"NSQD_TLS_CERT":   "",
		"NSQD_TLS_KEY":    "",

		"NSQLOOKUPD_URLS": "localhost:4161",

		"NSQD_HOST": "localhost",
		"NSQD_PORT": "4150",

		"NSQ_MAX_IN_FLIGHT": "20",
	})

	nsqConfig, err := NsqConfigFromEnv(E)
	if err != nil {
		return nil, err
	}

	return New(ProducerOpts{
		Host: E["NSQD_TLS"],
		Port: E["NSQD_PORT"],

		NsqConfig: nsqConfig,
	})
}

func NsqConfigFromEnv(E map[string]string) (*nsq.Config, error) {
	nsqConfig := nsq.NewConfig()
	if E["NSQD_TLS"] == "true" {
		nsqConfig.Set("tls_v1", true)
		nsqConfig.Set("tls_root_ca_file", E["NSQD_TLS_CACERT"])
		nsqConfig.Set("tls_cert", E["NSQD_TLS_CERT"])
		nsqConfig.Set("tls_key", E["NSQD_TLS_KEY"])
	}

	maxInFlight, err := strconv.Atoi(E["NSQ_MAX_IN_FLIGHT"])
	if err != nil {
		return nil, errors.Wrapf(err, "invalid max in flight: %s", E["NSQ_MAX_IN_FLIGHT"])
	}

	nsqConfig.Set("max_in_flight", maxInFlight)

	return nsqConfig, nil
}
