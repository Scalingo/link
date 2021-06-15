package ip

import (
	"context"
	"time"

	"github.com/Scalingo/go-utils/logger"
	"github.com/Scalingo/link/locker"
	"github.com/pkg/errors"
)

/* Stop process:
  1. Fetch information about other hosts registered on this IP and our state on the IP
  2. Begin the shutdown process (/!\ THIS SHOULD BE PROTECTED BY THE STOP MUTEX OR IT COULD LEAD TO INVALID STATE)
	2.1. Set the stop flag to true, since we are in the lock it will not impact
	the other process yet but it will ensure that if the process stop
	unexpectedly, all the other goroutine will stop and we will remove the IP by
	letting it decay.
	2.2. If we are the owner of the lock on the IP remove it
	2.3. Remove us from the list of potential hosts for this IP (UnlinkIP). This will trigger other hosts to try to get the IP
	2.4. Set the state machine to a state where it will remove the IP.
	2.5. Remove the IP from the interface.
*/

func (m *manager) Stop(ctx context.Context) error {
	log := logger.Get(ctx).WithField("task", "stop")
	ctx = logger.ToCtx(ctx, log)

	log.Info("Start stop preflight checks")
	hosts, err := m.storage.GetIPHosts(ctx, m.IP())
	if err != nil {
		return errors.Wrap(err, "fail to get new hosts")
	}

	isMaster, err := m.locker.IsMaster(ctx)
	if err != nil {
		if err == locker.ErrInvalidEtcdState { // If the key does not exist!
			isMaster = false
		} else {
			return errors.Wrap(err, "fail to know if we are master")
		}
	}

	m.stopMutex.Lock()
	defer m.stopMutex.Unlock()
	log.Info("Start stop process")

	m.stopped = true

	log.Info("Stop the locker")
	err = m.locker.Stop(ctx)
	if err != nil {
		return errors.Wrap(err, "fail to stop locker")
	}

	log.Info("Stop the watcher")
	err = m.watcher.Stop(ctx)
	if err != nil {
		log.WithError(err).Error("Fail to stop watcher")
	}

	log.Info("Unlink IP from the host")
	err = m.storage.UnlinkIPFromCurrentHost(ctx, m.ip)
	if err != nil {
		return errors.Wrap(err, "fail to unlink IP")
	}

	if isMaster && len(hosts) > 1 {
		log.Info("We were not alone, wait for other hosts to remove the IP")
		err := m.waitForReallocation(ctx)
		if err != nil {
			log.WithError(err).Error("Fail to reallocate IP, continuing shutdown")
		}
	}

	log.Info("Demoting ourself")
	if m.Status() != FAILING {
		// We cannot use SendEvent here because the isStopped method is blocked by the stopMutex
		m.eventChan <- DemotedEvent
	}

	// We can stop the FSM we do not need it anymore
	close(m.eventChan)

	log.Info("Stop process ended!")
	return nil
}

// ipCheckLoop will try to get the VIP at regular intervals
// This loop will try to get the IP if the current primary crashed.
func (m *manager) ipCheckLoop(ctx context.Context) {
	interval := time.Duration(m.IP().KeepaliveInterval) * time.Second
	if interval == 0 {
		interval = m.config.KeepAliveInterval
	}
	for {
		if m.isStopped() {
			return
		}

		m.tryToGetIP(ctx)

		time.Sleep(interval)
	}
}

func (m *manager) tryToGetIP(ctx context.Context) {
	if m.Status() == FAILING {
		return
	}

	log := logger.Get(ctx)
	// Refresh will try to set the lock on Etcd, if it fails, this mean that
	// there was an issue while trying to communicate with the Etcd cluser.
	err := m.locker.Refresh(ctx)
	if err != nil {
		// We do not want to send a fault event on every connection error. We will
		// wait to have multiple connection failure before sending an event to the
		// state machine.
		// This is done because the Fault event will make this instance PRIMARY
		// even if there is another PRIAMRY node in the cluster. This could lead to
		// connection RESET and client connection errors.
		m.keepaliveRetry++
		log.WithError(err).Info("Fail to refresh lock (retry)")
		if m.keepaliveRetry > m.config.KeepAliveRetry {
			// The connection with etcd is definitely lost, send a Fault event.
			log.WithError(err).Error("Fail to refresh lock")
			m.sendEvent(FaultEvent)
		}
		return
	}
	if m.keepaliveRetry > 0 {
		// We restored the connection with etcd, reset the fault counter.
		log.Infof("Lock refreshed after %v retries", m.keepaliveRetry)
		m.keepaliveRetry = 0
	}

	isMaster, err := m.locker.IsMaster(ctx)
	if err != nil {
		log.WithError(err).Error("Fail to check lock")
		// Fault event should only be handled in the Refresh operation where the
		// retry loop is present
		return
	}

	if isMaster {
		log.Debug("we are master, sending elected event")
		m.sendEvent(ElectedEvent)
	} else {
		log.Debug("we are not master, sending demoted event")
		m.sendEvent(DemotedEvent)
	}
}

// onTopologyChange is the called by the watcher. This will be called in one of 3 cases:
// 1. There is a node that joined the pool of nodes interested by this IP
// 2. There is a node that leaved the pool of nodes interested by this IP
// 3. The current master node is trying to initiate a failover
func (m *manager) onTopologyChange(ctx context.Context) {
	log := logger.Get(ctx)
	if m.isStopped() {
		return
	}

	// If we are already master we do not want to failover because:
	// 1. We already are master and our lock is still valid, no reason to refresh it.
	// 2. Another host has left the pool but since we are master, he was standby, no actions to take there.
	// 3. If there was a failover, it was initiated by the master node (us) and we do not want to take back the lock.
	if m.Status() == ACTIVATED {
		log.Info("Network topology changed but we are master, no actions needed")
		return
	}

	log.Info("Network topology changed, trying to get the IP")
	m.tryToGetIP(ctx)
}

func (m *manager) isStopped() bool {
	m.stopMutex.RLock()
	defer m.stopMutex.RUnlock()
	return m.stopped
}
