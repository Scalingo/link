package etcd

import (
	"fmt"

	etcdclient "go.etcd.io/etcd/client/v3"
)

// ClientFromEnv generates a etcd client (API v3) from the environment
// Look at ConfigFromEnv to get details about the environment variables used
func ClientFromEnv() (*etcdclient.Client, error) {
	config, err := ConfigFromEnv()
	if err != nil {
		return nil, fmt.Errorf("fail to create etcd v3 config: %v", err)
	}

	newClient, err := etcdclient.New(config)
	if err != nil {
		return nil, fmt.Errorf("fail to create etcd v3 client: %v", err)
	}
	return newClient, nil
}
