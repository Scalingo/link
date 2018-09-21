package ip

import (
	"context"
	"fmt"
	"time"

	"github.com/Scalingo/go-etcd-lock/lock"
	"github.com/Scalingo/go-utils/logger"
	"github.com/Scalingo/link/models"
)

func (m *manager) lockRoutine(ctx context.Context, c chan string) {
	log := logger.Get(ctx).WithField("process", "etcd_watcher")
	etcdKey := fmt.Sprintf("%s/locks/default/%s", models.ETCD_LINK_DIRECTORY, m.ip)
	for {
		waitTime := 3 * time.Second

		l, err := m.locker.Acquire(etcdKey, 2)
		if _, ok := err.(*lock.Error); ok {
			// The key is already locked
			c <- DemotedEvent
		} else if err != nil {
			log.WithError(err).Error("ETCD Lock failed")
			c <- FaultEvent
		} else {
			// The lock is our, get the key!
			c <- ElectedEvent
			waitTime = 2 * time.Second
		}

		time.Sleep(waitTime)
		if err == nil {
			l.Release()
		}

		m.stopMutex.RLock()
		stopping := m.stopping
		m.stopMutex.RUnlock()

		if stopping {
			// We should stop! Wait for another 3 seconds (the time needed dor someone else to get the IP), demote our self and stop
			time.Sleep(3 * time.Second)
			c <- DemotedEvent
			close(c)
			return
		}
	}
}
