package models

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/Scalingo/go-utils/etcd"
	"github.com/Scalingo/go-utils/logger"
	"github.com/Scalingo/link/v2/config"
	"github.com/gofrs/uuid"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	clientv3 "go.etcd.io/etcd/client/v3"
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

type etcdStorage struct {
	hostname string
}

func NewEtcdStorage(config config.Config) etcdStorage {
	return etcdStorage{
		hostname: config.Hostname,
	}
}

func (e etcdStorage) GetIPs(ctx context.Context) ([]IP, error) {
	client, closer, err := e.newEtcdClient()
	if err != nil {
		return nil, errors.Wrap(err, "fail to get etcd client")
	}
	defer closer.Close()

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resp, err := client.Get(ctx, fmt.Sprintf("%s/hosts/%s", EtcdLinkDirectory, e.hostname), clientv3.WithPrefix())
	if err != nil {
		return nil, errors.Wrap(err, "fail to get list of IPs from etcd")
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

	client, closer, err := e.newEtcdClient()
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

	etcdCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	_, err = client.Put(etcdCtx, e.keyFor(ip), string(value))
	if err != nil {
		return ip, errors.Wrapf(err, "fail to save IP")
	}

	return ip, nil
}

func (e etcdStorage) UpdateIP(ctx context.Context, ip IP) error {
	log := logger.Get(ctx)
	if ip.ID == "" {
		return fmt.Errorf("invalid IP ID: %s", ip.IP)
	}

	client, closer, err := e.newEtcdClient()
	if err != nil {
		return errors.Wrap(err, "fail to open client")
	}
	defer closer.Close()
	ip.Status = ""

	value, err := json.Marshal(ip)
	if err != nil {
		return errors.Wrap(err, "fail to marshal IP")
	}

	etcdKey := e.keyFor(ip)
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

func (e etcdStorage) RemoveIP(ctx context.Context, id string) error {
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

func (e etcdStorage) GetCurrentHost(ctx context.Context) (Host, error) {
	host, err := e.getHost(ctx, e.hostname)
	if err != nil {
		return host, errors.Wrap(err, "fail to get current host")
	}

	return host, nil
}

func (e etcdStorage) getHost(ctx context.Context, hostname string) (Host, error) {
	var host Host
	client, close, err := e.newEtcdClient()
	if err != nil {
		return host, errors.Wrap(err, "fail to get etcd client")
	}
	defer close.Close()

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

func (e etcdStorage) SaveHost(ctx context.Context, host Host) error {
	client, closer, err := e.newEtcdClient()
	if err != nil {
		return errors.Wrap(err, "fail to get etcd client")
	}
	defer closer.Close()

	value, err := json.Marshal(host)
	if err != nil {
		return errors.Wrap(err, "fail to marshal host")
	}

	etcdCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	_, err = client.Put(etcdCtx, e.keyForHost(e.hostname), string(value))
	if err != nil {
		return errors.Wrapf(err, "fail to save host")
	}

	return nil
}

func (e etcdStorage) LinkIPWithCurrentHost(ctx context.Context, ip IP) error {
	key := fmt.Sprintf("%s/ips/%s/%s", EtcdLinkDirectory, ip.StorableIP(), e.hostname)
	client, closer, err := e.newEtcdClient()
	if err != nil {
		return errors.Wrap(err, "fail to get etcd client")
	}
	defer closer.Close()

	payload, err := json.Marshal(IPLink{
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

func (e etcdStorage) UnlinkIPFromCurrentHost(ctx context.Context, ip IP) error {
	key := fmt.Sprintf("%s/ips/%s/%s", EtcdLinkDirectory, ip.StorableIP(), e.hostname)
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

func (e etcdStorage) GetIPHosts(ctx context.Context, ip IP) ([]string, error) {
	key := fmt.Sprintf("%s/ips/%s", EtcdLinkDirectory, ip.StorableIP())

	client, closer, err := e.newEtcdClient()
	if err != nil {
		return nil, errors.Wrap(err, "fail to get etcd client")
	}
	defer closer.Close()

	etcdCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resp, err := client.Get(etcdCtx, key, clientv3.WithPrefix())
	if err != nil {
		return nil, errors.Wrap(err, "fail to list ip links")
	}

	results := make([]string, resp.Count)
	for i, kv := range resp.Kvs {
		results[i] = string(kv.Key)
	}
	return results, nil
}

func (e etcdStorage) newEtcdClient() (clientv3.KV, io.Closer, error) {
	c, err := etcd.ClientFromEnv()
	if err != nil {
		return nil, nil, errors.Wrapf(err, "fail to get etcd client from config")
	}

	return clientv3.KV(c), c, nil
}

func (e etcdStorage) keyFor(ip IP) string {
	return fmt.Sprintf("%s/hosts/%s/%s", EtcdLinkDirectory, e.hostname, ip.ID)
}

func (e etcdStorage) keyForHost(hostname string) string {
	return fmt.Sprintf("%s/config/%s", EtcdLinkDirectory, hostname)
}
