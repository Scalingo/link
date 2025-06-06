package migrations

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/pkg/errors"
	etcdv3 "go.etcd.io/etcd/client/v3"

	"github.com/Scalingo/link/v3/models"
)

// v0IP represents an virtual IP as stored in the v0 data version. The main difference as compared to v1 is the LeaseID.
type v0IP struct {
	ID      string `json:"id"`
	IP      string `json:"ip"`
	LeaseID int64  `json:"lease_id,omitempty"`
}

func (v0IP v0IP) convertToV1() models.Endpoint {
	return models.Endpoint{
		ID: v0IP.ID,
		IP: v0IP.IP,
	}
}

type v0EtcdStorage struct {
	etcdClient     etcdv3.KV
	leaseManagerID etcdv3.LeaseID
}

func newV0EtcdStorage(etcdClient etcdv3.KV, leaseManagerID etcdv3.LeaseID) v0EtcdStorage {
	return v0EtcdStorage{
		etcdClient:     etcdClient,
		leaseManagerID: leaseManagerID,
	}
}

// getIPs gets in etcd the list of virtual IPs on the given host. It parses the result in the v0 data version.
func (e v0EtcdStorage) getIPs(ctx context.Context, hostname string) ([]v0IP, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	resp, err := e.etcdClient.Get(ctx, fmt.Sprintf("%s/hosts/%s", models.EtcdLinkDirectory, hostname), etcdv3.WithPrefix())
	if err != nil {
		return nil, errors.Wrap(err, "fail to get list of IPs from etcd")
	}

	ips := make([]v0IP, 0)
	for _, kv := range resp.Kvs {
		var ip v0IP
		err := json.Unmarshal(kv.Value, &ip)
		if err != nil {
			return nil, errors.Wrapf(err, "invalid json for %s", kv.Key)
		}
		ips = append(ips, ip)
	}
	return ips, nil
}

// isMaster checks if the given vIP in v0 data version is master.
func (e v0EtcdStorage) isMaster(ctx context.Context, ip v0IP) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	key := fmt.Sprintf("%s/default/%s", models.EtcdLinkDirectory, storableIP(ip.IP))
	resp, err := e.etcdClient.Get(ctx, key)
	if err != nil {
		return false, errors.Wrapf(err, "fail to get lock of the IP %s", ip.IP)
	}

	if len(resp.Kvs) != 1 {
		return false, fmt.Errorf("invalid etcd state (key '%s' not found!)", key)
	}

	return resp.Kvs[0].Lease == ip.LeaseID, nil
}

// putIP puts the v1 IP in the new etcd key.
func (e v0EtcdStorage) putIP(ctx context.Context, endpoint models.Endpoint, hostname string) error {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	key := fmt.Sprintf("%s/default/%s", models.EtcdLinkDirectory, endpoint.IP)
	_, err := e.etcdClient.Put(ctx, key, hostname, etcdv3.WithLease(e.leaseManagerID))
	if err != nil {
		return errors.Wrapf(err, "fail to put the IP '%s' in etcd", endpoint.IP)
	}
	return nil
}

func storableIP(i string) string {
	return strings.ReplaceAll(i, "/", "_")
}
