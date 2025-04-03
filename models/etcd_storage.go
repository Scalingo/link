package models

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"path/filepath"
	"time"

	"github.com/gofrs/uuid/v5"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	etcdv3 "go.etcd.io/etcd/client/v3"

	"github.com/Scalingo/go-utils/etcd"
	"github.com/Scalingo/go-utils/logger"
	"github.com/Scalingo/link/v2/config"
)

const (
	EtcdLinkDirectory = "/link"
)

// Used keys:
// /link/default/IP => Locks
// /link/hosts/HOSTNAME/IP => IP config
// /link/config/HOSTNAME => Hostname config
// /link/ips/IP/HOSTNAME => Link between IP and host

var (
	ErrIPAlreadyPresent = errors.New("IP already present")
	ErrHostNotFound     = errors.New("host not found")
)

type EtcdStorage struct {
	hostname string
}

func NewEtcdStorage(config config.Config) EtcdStorage {
	return EtcdStorage{
		hostname: config.Hostname,
	}
}

func (e EtcdStorage) GetEndpoints(ctx context.Context) (Endpoints, error) {
	client, closer, err := e.newEtcdClient()
	if err != nil {
		return nil, errors.Wrap(err, "fail to get etcd client")
	}
	defer closer.Close()

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resp, err := client.Get(ctx, fmt.Sprintf("%s/hosts/%s", EtcdLinkDirectory, e.hostname), etcdv3.WithPrefix())
	if err != nil {
		return nil, errors.Wrap(err, "fail to get list of IPs from etcd")
	}

	endpoints := make(Endpoints, 0)
	for _, kv := range resp.Kvs {
		var endpoint Endpoint
		err := json.Unmarshal(kv.Value, &endpoint)
		if err != nil {
			return nil, errors.Wrapf(err, "invalid json for %s", kv.Key)
		}
		endpoints = append(endpoints, endpoint)
	}
	return endpoints, nil
}

func (e EtcdStorage) AddEndpoint(ctx context.Context, endpoint Endpoint) (Endpoint, error) {
	client, closer, err := e.newEtcdClient()
	if err != nil {
		return endpoint, errors.Wrap(err, "fail to get etcd client")
	}
	defer closer.Close()

	id, err := uuid.NewV4()
	if err != nil {
		return endpoint, errors.Wrap(err, "fail to generate ID")
	}

	endpoint.ID = "vip-" + id.String()

	value, err := json.Marshal(endpoint)
	if err != nil {
		return endpoint, errors.Wrap(err, "fail to marshal IP")
	}

	etcdCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	_, err = client.Put(etcdCtx, e.keyFor(endpoint), string(value))
	if err != nil {
		return endpoint, errors.Wrapf(err, "fail to save IP")
	}

	return endpoint, nil
}

func (e EtcdStorage) UpdateEndpoint(ctx context.Context, endpoint Endpoint) error {
	log := logger.Get(ctx)
	if endpoint.ID == "" {
		return fmt.Errorf("invalid endpoint ID: %s", endpoint.ID)
	}

	client, closer, err := e.newEtcdClient()
	if err != nil {
		return errors.Wrap(err, "fail to open client")
	}
	defer closer.Close()

	value, err := json.Marshal(endpoint)
	if err != nil {
		return errors.Wrap(err, "fail to marshal endpoint")
	}

	etcdKey := e.keyFor(endpoint)
	log.WithFields(logrus.Fields{
		"etcd_key": etcdKey,
		"value":    string(value),
	}).Debug("Update the IP in etcd storage")

	etcdCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	_, err = client.Put(etcdCtx, etcdKey, string(value))
	if err != nil {
		return errors.Wrap(err, "fail to update the IP in etcd storage")
	}
	return nil
}

func (e EtcdStorage) RemoveEndpoint(ctx context.Context, id string) error {
	client, closer, err := e.newEtcdClient()
	if err != nil {
		return errors.Wrap(err, "fail to get etcd client")
	}
	defer closer.Close()

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err = client.Delete(ctx, fmt.Sprintf("%s/hosts/%s/%s", EtcdLinkDirectory, e.hostname, id))
	if err != nil {
		return errors.Wrap(err, "fail to delete IP")
	}
	return nil
}

func (e EtcdStorage) GetCurrentHost(ctx context.Context) (Host, error) {
	host, err := e.getHost(ctx, e.hostname)
	if err != nil {
		return host, errors.Wrap(err, "fail to get current host")
	}

	return host, nil
}

func (e EtcdStorage) getHost(ctx context.Context, hostname string) (Host, error) {
	var host Host
	client, clientClose, err := e.newEtcdClient()
	if err != nil {
		return host, errors.Wrap(err, "fail to get etcd client")
	}
	defer clientClose.Close()

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resp, err := client.Get(ctx, e.keyForHost(hostname))
	if err != nil {
		return host, errors.Wrap(err, "fail to get host from etcd")
	}

	if len(resp.Kvs) == 0 {
		return host, ErrHostNotFound
	}

	err = json.Unmarshal(resp.Kvs[0].Value, &host)
	if err != nil {
		return host, errors.Wrap(err, "fail to decode host config")
	}

	return host, nil
}

func (e EtcdStorage) SaveHost(ctx context.Context, host Host) error {
	client, closer, err := e.newEtcdClient()
	if err != nil {
		return errors.Wrap(err, "fail to get etcd client")
	}
	defer closer.Close()

	value, _ := json.Marshal(host)

	etcdCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	_, err = client.Put(etcdCtx, e.keyForHost(e.hostname), string(value))
	if err != nil {
		return errors.Wrapf(err, "fail to save host")
	}

	return nil
}

func (e EtcdStorage) LinkEndpointWithCurrentHost(ctx context.Context, lockKey string) error {
	key := fmt.Sprintf("%s/ips/%s/%s", EtcdLinkDirectory, lockKey, e.hostname)
	client, closer, err := e.newEtcdClient()
	if err != nil {
		return errors.Wrap(err, "fail to get etcd client")
	}
	defer closer.Close()

	payload, err := json.Marshal(EndpointLink{
		UpdatedAt: time.Now(),
	})
	if err != nil {
		return errors.Wrap(err, "fail to encode IP Link")
	}

	etcdCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	_, err = client.Put(etcdCtx, key, string(payload))
	if err != nil {
		return errors.Wrap(err, "fail to save ip link")
	}
	return nil
}

func (e EtcdStorage) UnlinkEndpointFromCurrentHost(ctx context.Context, lockKey string) error {
	key := fmt.Sprintf("%s/ips/%s/%s", EtcdLinkDirectory, lockKey, e.hostname)
	client, closer, err := e.newEtcdClient()
	if err != nil {
		return errors.Wrap(err, "fail to get etcd client")
	}
	defer closer.Close()

	etcdCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	_, err = client.Delete(etcdCtx, key)
	if err != nil {
		return errors.Wrap(err, "fail to delete IP link")
	}
	return nil
}

func (e EtcdStorage) GetEndpointHosts(ctx context.Context, lockKey string) ([]string, error) {
	key := fmt.Sprintf("%s/ips/%s", EtcdLinkDirectory, lockKey)

	client, closer, err := e.newEtcdClient()
	if err != nil {
		return nil, errors.Wrap(err, "fail to get etcd client")
	}
	defer closer.Close()

	etcdCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resp, err := client.Get(etcdCtx, key, etcdv3.WithPrefix())
	if err != nil {
		return nil, errors.Wrap(err, "fail to list ip links")
	}

	results := make([]string, resp.Count)
	for i, kv := range resp.Kvs {
		results[i] = filepath.Base(string(kv.Key))
	}
	return results, nil
}

func (e EtcdStorage) newEtcdClient() (etcdv3.KV, io.Closer, error) {
	c, err := etcd.ClientFromEnv()
	if err != nil {
		return nil, nil, errors.Wrapf(err, "fail to get etcd client from config")
	}

	return etcdv3.KV(c), c, nil
}

func (e EtcdStorage) keyFor(endpoint Endpoint) string {
	return fmt.Sprintf("%s/hosts/%s/%s", EtcdLinkDirectory, e.hostname, endpoint.ID)
}

func (e EtcdStorage) keyForHost(hostname string) string {
	return fmt.Sprintf("%s/config/%s", EtcdLinkDirectory, hostname)
}
