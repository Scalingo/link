package migrations

import (
	"context"
	"io"

	scalingoerrors "github.com/Scalingo/go-utils/errors"
	"github.com/Scalingo/go-utils/etcd"
	"github.com/Scalingo/go-utils/logger"
	"github.com/Scalingo/link/locker"
	"github.com/Scalingo/link/models"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"go.etcd.io/etcd/v3/clientv3"
)

type V0toV1 struct {
	hostname     string
	leaseManager locker.EtcdLeaseManager
	storage      models.Storage
}

func NewV0toV1Migration(hostname string, leaseManager locker.EtcdLeaseManager, storage models.Storage) V0toV1 {
	return V0toV1{
		hostname:     hostname,
		leaseManager: leaseManager,
		storage:      storage,
	}
}

func (m V0toV1) NeedsMigration(ctx context.Context) (bool, error) {
	log := logger.Get(ctx)

	host, err := m.storage.GetCurrentHost(ctx)
	if err != nil && scalingoerrors.RootCause(err) != models.ErrHostNotFound {
		return false, errors.Wrap(err, "fail to get current host to check if it needs data migration from v0 to v1")
	}

	if scalingoerrors.RootCause(err) == models.ErrHostNotFound {
		return true, nil
	}

	if host.DataVersion >= 1 {
		log.Info("Current host does not need data migration from v0 to v1")
		return false, nil
	}

	log.Info("Current host needs data migration from v0 to v1")
	return true, nil
}

func (m V0toV1) Migrate(ctx context.Context) error {
	log := logger.Get(ctx)
	log.Info("Migrate data from v0 to v1")

	leaseManagerID, err := m.leaseManager.GetLease(ctx)
	if err != nil {
		return errors.Wrap(err, "fail to get lease manager ID")
	}

	etcdClient, closer, err := newEtcdClient()
	if err != nil {
		return errors.Wrap(err, "fail to create etcd client to migrate data from v0 to v1")
	}
	defer closer.Close()

	v0Storage := newV0EtcdStorage(etcdClient, leaseManagerID)
	ips, err := v0Storage.getIPs(ctx, m.hostname)
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
		err = v0Storage.putIP(ctx, ip.convertToV1(), m.hostname)
		if err != nil {
			return errors.Wrap(err, "fail to update the IP data during the migration from v0 to v1")
		}
	}

	host, err := m.storage.GetCurrentHost(ctx)
	if err != nil {
		return errors.Wrap(err, "fail to get current host to update its data version")
	}
	host.DataVersion = locker.DataVersion

	err = m.storage.SaveHost(ctx, host)
	if err != nil {
		return errors.Wrap(err, "fail to save host at the end of the v0 to v1 migration")
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
