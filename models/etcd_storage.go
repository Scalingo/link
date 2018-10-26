package models

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/Scalingo/go-utils/etcd"
	"github.com/Scalingo/link/config"
	"github.com/coreos/etcd/clientv3"
	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

const (
	ETCD_LINK_DIRECTORY = "/link"
)

var (
	ErrIPAlreadyPresent = errors.New("IP already present")
)

type etcdStorage struct {
	hostname string
}

func NewEtcdStorage(config config.Config) etcdStorage {
	return etcdStorage{
		hostname: config.Hostname,
	}
}

func (e etcdStorage) GetIPs(ctx context.Context) ([]IP, error) {
	client, closer, err := e.NewEtcdClient()
	if err != nil {
		return nil, errors.Wrap(err, "fail to get etcd client")
	}
	defer closer.Close()

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resp, err := client.Get(ctx, fmt.Sprintf("%s/hosts/%s", ETCD_LINK_DIRECTORY, e.hostname), clientv3.WithPrefix())
	if err != nil {
		return nil, err
	}

	ips := make([]IP, 0)
	for _, kv := range resp.Kvs {
		var ip IP
		err := json.Unmarshal(kv.Value, &ip)
		if err != nil {
			return nil, errors.Wrapf(err, "invalid json for %s", kv.Key)
		}
		ips = append(ips, ip)
	}
	return ips, nil
}

func (e etcdStorage) AddIP(ctx context.Context, ip IP) (IP, error) {
	ips, err := e.GetIPs(ctx)
	if err != nil {
		return ip, errors.Wrap(err, "fail to contact etcd")
	}

	for _, i := range ips {
		if i.IP == ip.IP {
			return i, ErrIPAlreadyPresent
		}
	}

	client, closer, err := e.NewEtcdClient()
	if err != nil {
		return ip, errors.Wrap(err, "fail to get etcd client")
	}
	defer closer.Close()

	id, err := uuid.NewV4()
	if err != nil {
		return ip, errors.Wrap(err, "fail to generate ID")
	}

	ip.ID = "vip-" + id.String()

	// We do not want to store the status in database. This will always be STANDBY.
	// So we set it to "" and let the omitempty do its job.
	ip.Status = ""
	value, err := json.Marshal(ip)
	if err != nil {
		return ip, errors.Wrap(err, "fail to marshal IP")
	}
	key := fmt.Sprintf("%s/hosts/%s/%s", ETCD_LINK_DIRECTORY, e.hostname, ip.ID)

	etcdCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	_, err = client.Put(etcdCtx, key, string(value))
	if err != nil {
		return ip, errors.Wrapf(err, "fail to save IP")
	}

	return ip, nil
}

func (e etcdStorage) RemoveIP(ctx context.Context, id string) error {
	client, closer, err := e.NewEtcdClient()
	if err != nil {
		return errors.Wrap(err, "fail to get etcd client")
	}
	defer closer.Close()

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	_, err = client.Delete(ctx, fmt.Sprintf("%s/hosts/%s/%s", ETCD_LINK_DIRECTORY, e.hostname, id))
	if err != nil {
		return errors.Wrap(err, "fail to delete IP")
	}
	return nil
}

func (e etcdStorage) NewEtcdClient() (clientv3.KV, io.Closer, error) {
	c, err := etcd.ClientFromEnv()
	if err != nil {
		return nil, nil, errors.Wrapf(err, "fail to get etcd client from config")
	}

	return clientv3.KV(c), c, nil
}
