package ip

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Scalingo/go-utils/logger"
	"github.com/Scalingo/link/models"
	"github.com/coreos/etcd/clientv3"
)

func (m *manager) lockRoutine(ctx context.Context, c chan string) {
	log := logger.Get(ctx).WithField("process", "etcd_watcher")
	etcdKey := fmt.Sprintf("%s/default/%s", models.ETCD_LINK_DIRECTORY, strings.Replace(m.ip, "/", "_", -1))

	grant, err := m.etcd.Grant(ctx, 5)
	if err != nil {
		log.WithError(err).Error("Fail to get lease")
		panic("ETCD ERROR")
	}

	for {
		_, err := m.etcd.Txn(ctx).
			If(clientv3.Compare(clientv3.CreateRevision(etcdKey), "=", 0)).
			Then(clientv3.OpPut(etcdKey, "locked", clientv3.WithLease(grant.ID))).
			Commit()
		if err != nil {
			c <- FaultEvent
		} else {
			resp, err := m.etcd.Get(ctx, etcdKey)
			if err != nil {
				c <- FaultEvent
			} else {
				if len(resp.Kvs) != 1 {
					// DAFUK :/
					c <- FaultEvent
				} else {
					if resp.Kvs[0].Lease == int64(grant.ID) {
						c <- ElectedEvent
					} else {
						c <- DemotedEvent
					}
				}
			}
		}

		time.Sleep(3 * time.Second)

		_, err = m.etcd.KeepAliveOnce(ctx, grant.ID)
		if err != nil {
			log.WithError(err).Error("Fail to send keep alive")
		}

		m.stopMutex.RLock()
		stopping := m.stopping
		m.stopMutex.RUnlock()

		if stopping {
			// We should stop! Wait for another 10 seconds (the time needed dor someone else to get the IP), demote our self and stop
			time.Sleep(10 * time.Second)
			c <- DemotedEvent
			close(c)
			return
		}
	}
}
