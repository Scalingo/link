package migrations

import (
	"context"
	"io"

	"github.com/Scalingo/go-utils/etcd"
	"github.com/Scalingo/go-utils/logger"
	"github.com/Scalingo/link/locker"
	"github.com/Scalingo/link/models"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"go.etcd.io/etcd/v3/clientv3"
)

type V0toV1 struct {
	// TODO Do we need the lease manager
	leaseManager locker.EtcdLeaseManager
	storage      models.Storage
}

func NewV0toV1Migration(leaseManager locker.EtcdLeaseManager, storage models.Storage) V0toV1 {
	return V0toV1{
		leaseManager: leaseManager,
		storage:      storage,
	}
}

func (m V0toV1) NeedsMigration(ctx context.Context) bool {
	log := logger.Get(ctx)
	log.Info("Data migration from v0 to v1 is needed")
	// TODO
	return false
}

func (m V0toV1) Migrate(ctx context.Context) error {
	log := logger.Get(ctx)
	log.Info("Migrate data from v0 to v1")

	host, err := m.storage.GetCurrentHost(ctx)
	if err != nil {
		return errors.Wrap(err, "fail to get the current host information to migrate data from v0 to v1")
	}

	etcdClient, closer, err := newEtcdClient()
	if err != nil {
		return errors.Wrap(err, "fail to create etcd client to migrate data from v0 to v1")
	}
	defer closer.Close()

	v0Storage := newV0EtcdStorage(etcdClient, clientv3.LeaseID(host.LeaseID))
	ips, err := v0Storage.getIPs(ctx, host.Hostname)
	if err != nil {
		return errors.Wrap(err, "fail to get the list of v0 IPs")
	}

	for _, ip := range ips {
		log := log.WithFields(logrus.Fields{
			"id":       ip.ID,
			"ip":       ip.IP,
			"lease_id": ip.LeaseID,
		})

		isMaster, err := v0Storage.isMaster(ctx, ip)
		if err != nil {
			return errors.Wrap(err, "fail to get the master status of the IP")
		}
		if !isMaster {
			log.Info("Host is not master of this IP, nothing to do")
			continue
		}

		log.Info("Host is master of this IP, migrate the data")
		err = v0Storage.putIP(ctx, ip.convertToV1(), host.Hostname)
		if err != nil {
			return errors.Wrap(err, "fail to update the IP data during the migration from v0 to v1")
		}
	}

	log.Info("End of the data migration from v0 to v1")
	return nil
}

func newEtcdClient() (clientv3.KV, io.Closer, error) {
	c, err := etcd.ClientFromEnv()
	if err != nil {
		return nil, nil, errors.Wrapf(err, "fail to get etcd client from config")
	}

	return clientv3.KV(c), c, nil
}
