package etcd

import (
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"

	"go.etcd.io/etcd/v3/clientv3"
	"go.etcd.io/etcd/v3/pkg/transport"
)

// ConfigFromEnv generates a etcd clientv3 config from the environment using the following variables:
// * ETCD_HOSTS: The different endpoints of etcd members
// * ETCD_TLS_CERT: Path to the TLS X.509 certificate
// * ETCD_TLS_KEY: Path to the private key authenticating the certificate
// * ETCD_CACERT: Path to the CA cert signing the etcd member certifcates
func ConfigFromEnv() (res clientv3.Config, _ error) {
	endpoints := []string{}
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

	config := clientv3.Config{
		Endpoints:   endpoints,
		DialTimeout: 5 * time.Second,
	}

	if tls {
		tlsInfo := transport.TLSInfo{
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
