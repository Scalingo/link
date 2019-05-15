package locker

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Scalingo/go-utils/logger"
	"github.com/Scalingo/link/config"
	"github.com/Scalingo/link/models"
	"github.com/coreos/etcd/clientv3"
	"github.com/pkg/errors"
)

type etcdLocker struct {
	kvEtcd           clientv3.KV
	leaseEtcd        clientv3.Lease
	leaseID          clientv3.LeaseID
	key              string
	config           config.Config
	lastLeaseRefresh time.Time
}

func NewEtcdLocker(config config.Config, etcd *clientv3.Client, ip string) *etcdLocker {
	key := fmt.Sprintf("%s/default/%s", models.ETCD_LINK_DIRECTORY, strings.Replace(ip, "/", "_", -1))
	return &etcdLocker{
		kvEtcd:    etcd,
		leaseEtcd: etcd,
		key:       key,
		leaseID:   0,
		config:    config,
	}
}

func (l *etcdLocker) Refresh(ctx context.Context) error {
	log := logger.Get(ctx)

	if l.leaseID == 0 {
		grant, err := l.leaseEtcd.Grant(ctx, int64(l.config.LeaseTime().Seconds()))
		if err != nil {
			return errors.Wrap(err, "fail to generate grant")
		}

		l.leaseID = grant.ID
	}

	// The goal of this transaction is to create the key with our leaseID only if this key does not exist
	// We use a transaction to make sure that concurrent tries wont interfere with each others.

	_, err := l.kvEtcd.Txn(ctx).
		// If the key does not exists (createRevision == 0)
		If(clientv3.Compare(clientv3.CreateRevision(l.key), "=", 0)).
		// Create it with our leaseID
		Then(clientv3.OpPut(l.key, "locked", clientv3.WithLease(l.leaseID))).
		Commit()
	if err != nil {
		if l.leaseExpired() {
			// We got an error, this can be because our leaseID is not valid anymore: Reset it
			oldLeaseID := l.leaseID
			l.leaseID = 0
			return errors.Wrapf(err, "fail to refresh lock: probably expired (leaseID = %v)", oldLeaseID)
		} else {
			// We got an error. This is probably not relaterd to an expired lease. Do not reset it
			return errors.Wrapf(err, "fail to refresh lock")
		}
	}

	_, err = l.leaseEtcd.KeepAliveOnce(ctx, l.leaseID)
	if err != nil {
		if l.leaseExpired() {
			l.leaseID = 0
			log.WithError(err).Error("Keep alive failed: Regenerate lease")
		} else {
			// We got an error while sending keepalive
			log.WithError(err).Error("Keep alive failed")
		}
	}

	l.lastLeaseRefresh = time.Now()
	return nil
}

func (l *etcdLocker) leaseExpired() bool {
	return l.lastLeaseRefresh.IsZero() || time.Now().After(l.lastLeaseRefresh.Add(l.config.LeaseTime()))
}

func (l *etcdLocker) Unlock(ctx context.Context) error {
	_, err := l.kvEtcd.Delete(ctx, l.key)
	if err != nil {
		return errors.Wrap(err, "fail to unlock key")
	}
	l.leaseID = 0
	return nil
}

func (l *etcdLocker) IsMaster(ctx context.Context) (bool, error) {
	resp, err := l.kvEtcd.Get(ctx, l.key)
	if err != nil {
		return false, errors.Wrap(err, "fail to get lock")
	}

	if len(resp.Kvs) != 1 {
		// DAFUK :/
		return false, errors.New("invalid etcd state (key not found!)")
	}

	return resp.Kvs[0].Lease == int64(l.leaseID), nil
}

func (l *etcdLocker) Stop(ctx context.Context) error {
	// Reset the lease and let the old lease die.
	// Setting the leaseID to 0 will ensure that the next time `Refresh` is
	// called, we will work with a new lease.
	l.leaseID = 0
	return nil
}
