package etcd

import (
	"fmt"

	"github.com/coreos/etcd/clientv3"
)

// ClientFromEnv generates a etcd client (API v3) from the environment
// Look at ConfigFromEnv to get details about the environment variables used
func ClientFromEnv() (*clientv3.Client, error) {
	config, err := ConfigFromEnv()
	if err != nil {
		return nil, fmt.Errorf("fail to create etcd v3 config: %v", err)
	}

	client, err := clientv3.New(config)
	if err != nil {
		return nil, fmt.Errorf("fail to create etcdv3 client: %v", err)
	}
	return client, nil
}
