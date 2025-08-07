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
	"github.com/Scalingo/link/v3/config"
)

const (
	EtcdLinkDirectory = "/link"

	etcdTimeout = 5 * time.Second

	encryptedDataIDPrefix = "sec-"
	endpointIDPrefix      = "vip-"
)

// Used keys:
// /link/default/IP => Locks
// /link/hosts/HOSTNAME/IP => IP config
// /link/config/HOSTNAME => Hostname config
// /link/ips/IP/HOSTNAME => Link between IP and host
// /link/secrets/hosts/HOSTNAME/ENDPOINT_ID/ENCRYPTED_DATA_ID => Encrypted data for an endpoint (used by Plugins)

var (
	ErrEndpointAlreadyPresent = errors.New("endpoint already present")
	ErrHostNotFound           = errors.New("host not found")
	ErrEncryptedDataNotFound  = errors.New("encrypted data not found")
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
		return nil, errors.Wrap(err, "get etcd client")
	}
	defer closer.Close()

	ctx, cancel := context.WithTimeout(ctx, etcdTimeout)
	defer cancel()

	resp, err := client.Get(ctx, fmt.Sprintf("%s/hosts/%s", EtcdLinkDirectory, e.hostname), etcdv3.WithPrefix())
	if err != nil {
		return nil, errors.Wrap(err, "get list of endpoints from etcd")
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
		return endpoint, errors.Wrap(err, "get etcd client")
	}
	defer closer.Close()

	id, err := uuid.NewV4()
	if err != nil {
		return endpoint, errors.Wrap(err, "generate ID")
	}

	endpoint.ID = endpointIDPrefix + id.String()

	value, err := json.Marshal(endpoint)
	if err != nil {
		return endpoint, errors.Wrap(err, "marshal endpoint")
	}

	etcdCtx, cancel := context.WithTimeout(ctx, etcdTimeout)
	defer cancel()
	_, err = client.Put(etcdCtx, e.keyFor(endpoint), string(value))
	if err != nil {
		return endpoint, errors.Wrapf(err, "save endpoint")
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
		return errors.Wrap(err, "open client")
	}
	defer closer.Close()

	value, err := json.Marshal(endpoint)
	if err != nil {
		return errors.Wrap(err, "marshal endpoint")
	}

	etcdKey := e.keyFor(endpoint)
	log.WithFields(logrus.Fields{
		"etcd_key": etcdKey,
		"value":    string(value),
	}).Debug("Update the endpoint in etcd storage")

	etcdCtx, cancel := context.WithTimeout(ctx, etcdTimeout)
	defer cancel()
	_, err = client.Put(etcdCtx, etcdKey, string(value))
	if err != nil {
		return errors.Wrap(err, "update the endpoint in etcd storage")
	}
	return nil
}

func (e EtcdStorage) RemoveEndpoint(ctx context.Context, id string) error {
	client, closer, err := e.newEtcdClient()
	if err != nil {
		return errors.Wrap(err, "get etcd client")
	}
	defer closer.Close()

	ctx, cancel := context.WithTimeout(ctx, etcdTimeout)
	defer cancel()

	_, err = client.Delete(ctx, fmt.Sprintf("%s/hosts/%s/%s", EtcdLinkDirectory, e.hostname, id))
	if err != nil {
		return errors.Wrap(err, "delete endpoint")
	}
	return nil
}

func (e EtcdStorage) GetCurrentHost(ctx context.Context) (Host, error) {
	host, err := e.getHost(ctx, e.hostname)
	if err != nil {
		return host, errors.Wrap(err, "get current host")
	}

	return host, nil
}

func (e EtcdStorage) getHost(ctx context.Context, hostname string) (Host, error) {
	var host Host
	client, clientClose, err := e.newEtcdClient()
	if err != nil {
		return host, errors.Wrap(err, "get etcd client")
	}
	defer clientClose.Close()

	ctx, cancel := context.WithTimeout(ctx, etcdTimeout)
	defer cancel()

	resp, err := client.Get(ctx, e.keyForHost(hostname))
	if err != nil {
		return host, errors.Wrap(err, "get host from etcd")
	}

	if len(resp.Kvs) == 0 {
		return host, ErrHostNotFound
	}

	err = json.Unmarshal(resp.Kvs[0].Value, &host)
	if err != nil {
		return host, errors.Wrap(err, "decode host config")
	}

	return host, nil
}

func (e EtcdStorage) SaveHost(ctx context.Context, host Host) error {
	client, closer, err := e.newEtcdClient()
	if err != nil {
		return errors.Wrap(err, "get etcd client")
	}
	defer closer.Close()

	value, _ := json.Marshal(host)

	etcdCtx, cancel := context.WithTimeout(ctx, etcdTimeout)
	defer cancel()
	_, err = client.Put(etcdCtx, e.keyForHost(e.hostname), string(value))
	if err != nil {
		return errors.Wrapf(err, "save host")
	}

	return nil
}

func (e EtcdStorage) LinkEndpointWithCurrentHost(ctx context.Context, lockKey string) error {
	key := fmt.Sprintf("%s/ips/%s/%s", EtcdLinkDirectory, lockKey, e.hostname)
	client, closer, err := e.newEtcdClient()
	if err != nil {
		return errors.Wrap(err, "get etcd client")
	}
	defer closer.Close()

	payload, err := json.Marshal(EndpointLink{
		UpdatedAt: time.Now(),
	})
	if err != nil {
		return errors.Wrap(err, "encode endpoint Link")
	}

	etcdCtx, cancel := context.WithTimeout(ctx, etcdTimeout)
	defer cancel()
	_, err = client.Put(etcdCtx, key, string(payload))
	if err != nil {
		return errors.Wrap(err, "save endpoint link")
	}
	return nil
}

func (e EtcdStorage) UnlinkEndpointFromCurrentHost(ctx context.Context, lockKey string) error {
	key := fmt.Sprintf("%s/ips/%s/%s", EtcdLinkDirectory, lockKey, e.hostname)
	client, closer, err := e.newEtcdClient()
	if err != nil {
		return errors.Wrap(err, "get etcd client")
	}
	defer closer.Close()

	etcdCtx, cancel := context.WithTimeout(ctx, etcdTimeout)
	defer cancel()
	_, err = client.Delete(etcdCtx, key)
	if err != nil {
		return errors.Wrap(err, "delete endpoint link")
	}
	return nil
}

func (e EtcdStorage) GetEndpointHosts(ctx context.Context, lockKey string) ([]string, error) {
	key := fmt.Sprintf("%s/ips/%s", EtcdLinkDirectory, lockKey)

	client, closer, err := e.newEtcdClient()
	if err != nil {
		return nil, errors.Wrap(err, "get etcd client")
	}
	defer closer.Close()

	etcdCtx, cancel := context.WithTimeout(ctx, etcdTimeout)
	defer cancel()

	resp, err := client.Get(etcdCtx, key, etcdv3.WithPrefix())
	if err != nil {
		return nil, errors.Wrap(err, "list ip links")
	}

	results := make([]string, resp.Count)
	for i, kv := range resp.Kvs {
		results[i] = filepath.Base(string(kv.Key))
	}
	return results, nil
}

func (e EtcdStorage) GetEncryptedData(ctx context.Context, endpointID string, encryptedDataId string) (EncryptedData, error) {
	client, closer, err := e.newEtcdClient()
	if err != nil {
		return EncryptedData{}, errors.Wrap(err, "get etcd client")
	}
	defer closer.Close()
	etcdCtx, cancel := context.WithTimeout(ctx, etcdTimeout)
	defer cancel()

	resp, err := client.Get(etcdCtx, e.keyForEncryptedData(endpointID, encryptedDataId))
	if err != nil {
		return EncryptedData{}, errors.Wrap(err, "get encrypted data")
	}
	if len(resp.Kvs) == 0 {
		return EncryptedData{}, ErrEncryptedDataNotFound
	}
	var encryptedData EncryptedData
	err = json.Unmarshal(resp.Kvs[0].Value, &encryptedData)
	if err != nil {
		return EncryptedData{}, errors.Wrap(err, "decode encrypted data")
	}
	return encryptedData, nil
}

func (e EtcdStorage) UpsertEncryptedData(ctx context.Context, endpointID string, encryptedData EncryptedData) (EncryptedDataLink, error) {
	client, closer, err := e.newEtcdClient()
	if err != nil {
		return EncryptedDataLink{}, errors.Wrap(err, "get etcd client")
	}
	defer closer.Close()

	if encryptedData.ID == "" {
		id, err := uuid.NewV4()
		if err != nil {
			return EncryptedDataLink{}, errors.Wrap(err, "generate ID")
		}

		encryptedData.ID = encryptedDataIDPrefix + id.String()
	}
	encryptedData.EndpointID = endpointID

	value, err := json.Marshal(encryptedData)
	if err != nil {
		return EncryptedDataLink{}, errors.Wrap(err, "marshal encrypted data")
	}
	etcdCtx, cancel := context.WithTimeout(ctx, etcdTimeout)
	defer cancel()
	_, err = client.Put(etcdCtx, e.keyForEncryptedData(endpointID, encryptedData.ID), string(value))
	if err != nil {
		return EncryptedDataLink{}, errors.Wrap(err, "save encrypted data")
	}
	return EncryptedDataLink{
		ID:         encryptedData.ID,
		EndpointID: endpointID,
	}, nil
}

func (e EtcdStorage) RemoveEncryptedDataForEndpoint(ctx context.Context, endpointID string) error {
	etcdCtx, cancel := context.WithTimeout(ctx, etcdTimeout)
	defer cancel()

	client, closer, err := e.newEtcdClient()
	if err != nil {
		return errors.Wrap(err, "get etcd client")
	}
	defer closer.Close()

	_, err = client.Delete(etcdCtx, fmt.Sprintf("%s/secrets/hosts/%s/%s", EtcdLinkDirectory, e.hostname, endpointID), etcdv3.WithPrefix())
	if err != nil {
		return errors.Wrap(err, "delete encrypted data")
	}
	return nil
}

func (e EtcdStorage) ListEncryptedDataForHost(ctx context.Context) ([]EncryptedData, error) {
	etcdCtx, cancel := context.WithTimeout(ctx, etcdTimeout)
	defer cancel()

	client, closer, err := e.newEtcdClient()
	if err != nil {
		return nil, errors.Wrap(err, "get etcd client")
	}
	defer closer.Close()

	resp, err := client.Get(etcdCtx, fmt.Sprintf("%s/secrets/hosts/%s", EtcdLinkDirectory, e.hostname), etcdv3.WithPrefix())
	if err != nil {
		return nil, errors.Wrap(err, "list encrypted data")
	}

	results := make([]EncryptedData, 0, resp.Count)
	for _, kv := range resp.Kvs {
		var encryptedData EncryptedData
		err = json.Unmarshal(kv.Value, &encryptedData)
		if err != nil {
			return nil, errors.Wrapf(err, "invalid json for %s", kv.Key)
		}

		results = append(results, encryptedData)
	}
	return results, nil
}

func (e EtcdStorage) newEtcdClient() (etcdv3.KV, io.Closer, error) {
	c, err := etcd.ClientFromEnv()
	if err != nil {
		return nil, nil, errors.Wrapf(err, "get etcd client from config")
	}

	return etcdv3.KV(c), c, nil
}

func (e EtcdStorage) keyFor(endpoint Endpoint) string {
	return fmt.Sprintf("%s/hosts/%s/%s", EtcdLinkDirectory, e.hostname, endpoint.ID)
}

func (e EtcdStorage) keyForHost(hostname string) string {
	return fmt.Sprintf("%s/config/%s", EtcdLinkDirectory, hostname)
}

func (e EtcdStorage) keyForEncryptedData(endpointID, encryptedDataID string) string {
	return fmt.Sprintf("%s/secrets/hosts/%s/%s/%s", EtcdLinkDirectory, e.hostname, endpointID, encryptedDataID)
}
