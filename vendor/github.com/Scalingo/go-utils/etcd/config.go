package etcd

import (
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"

	etcdtransport "go.etcd.io/etcd/client/pkg/v3/transport"
	etcdclient "go.etcd.io/etcd/client/v3"
)

// ConfigFromEnv generates a etcd clientv3 config from the environment using the following variables:
// * ETCD_HOSTS: The different endpoints of etcd members
// * ETCD_TLS_CERT: Path to the TLS X.509 certificate
// * ETCD_TLS_KEY: Path to the private key authenticating the certificate
// * ETCD_CACERT: Path to the CA cert signing the etcd member certifcates
func ConfigFromEnv() (res etcdclient.Config, _ error) {
	var endpoints []string
	etcdURLs := strings.Split(os.Getenv("ETCD_HOSTS"), ",")
	tls := false

	for _, u := range etcdURLs {
		etcdURL, err := url.Parse(strings.TrimSpace(u))
		if err != nil {
			return res, fmt.Errorf("invalid URL in ETCD_HOSTS %s: %v", u, err)
		}
		if etcdURL.Scheme == "https" {
			tls = true
		}
		endpoints = append(endpoints, etcdURL.Host)
	}

	config := etcdclient.Config{
		Endpoints:   endpoints,
		DialTimeout: 5 * time.Second,
	}

	if tls {
		tlsInfo := etcdtransport.TLSInfo{
			CertFile:      os.Getenv("ETCD_TLS_CERT"),
			KeyFile:       os.Getenv("ETCD_TLS_KEY"),
			TrustedCAFile: os.Getenv("ETCD_CACERT"),
		}
		tlsConfig, err := tlsInfo.ClientConfig()
		if err != nil {
			return res, fmt.Errorf("fail to create tls info config: %v", err)
		}
		config.TLS = tlsConfig
	}
	return config, nil
}
