package ip

import (
	"context"
	"time"

	"github.com/pkg/errors"

	"github.com/Scalingo/go-utils/logger"
	"github.com/Scalingo/link/v2/locker"
)

/* Stop process:
  1. Fetch information about other hosts registered on this IP and our state on the IP
  2. Begin the shutdown process (/!\ THIS SHOULD BE PROTECTED BY THE STOP MUTEX OR IT COULD LEAD TO INVALID STATE)
	2.1. Set the stop flag to true, since we are in the lock it will not impact
	the other processes yet. But it will ensure that if the process stops
	unexpectedly, all the other goroutines will stop and we will remove the IP by
	letting it decay.
	2.2. If we are the owner of the lock on the IP remove it
	2.3. Remove us from the list of potential hosts for this IP (UnlinkIPFromCurrentHost). This will trigger other hosts to try to get the IP
	2.4. Send a demoted event so that the state machine is in a state to remote the IP
*/

func (m *EndpointManager) Stop(ctx context.Context) error {
	log := logger.Get(ctx).WithField("process", "stop")
	ctx = logger.ToCtx(ctx, log)

	log.Info("Stops the IP manager")
	hosts, err := m.storage.GetEndpointHosts(ctx, m.plugin.ElectionKey(ctx))
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
	log.Info("Start the stop process")

	m.stopped = true

	log.Info("Stop the locker")
	err = m.locker.Stop(ctx)
	if err != nil {
		return errors.Wrap(err, "fail to stop locker")
	}

	log.Info("Stop the watcher")
	m.watcher.Stop(ctx)

	log.Info("Unlink IP from the host")
	err = m.storage.UnlinkEndpointFromCurrentHost(ctx, m.plugin.ElectionKey(ctx))
	if err != nil {
		return errors.Wrap(err, "fail to unlink IP")
	}

	if isMaster && len(hosts) > 1 {
		log.Info("We were not alone, wait for an other host to get the IP")
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

// ipCheckLoop tries to get the VIP at regular intervals. It is useful to get the IP if the current primary crashed.
func (m *EndpointManager) ipCheckLoop(ctx context.Context) {
	for {
		if m.isStopped() {
			return
		}

		m.tryToGetIP(ctx)

		time.Sleep(m.config.KeepAliveInterval)
	}
}

func (m *EndpointManager) tryToGetIP(ctx context.Context) {
	if m.Status() == FAILING {
		return
	}

	log := logger.Get(ctx)
	// Refresh will try to set the lock on etcd. If it fails, this means that
	// there was an issue while trying to communicate with the etcd cluster.
	err := m.locker.Refresh(ctx)
	if err != nil {
		// We do not want to send a fault event on every connection error. We
		// wait for multiple connection failures before sending an event to the
		// state machine.
		// This is done because the Fault event will make this instance master
		// even if there is another master node in the cluster. This could lead to
		// TCP connection RESET and client connection errors.
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
		log.Debug("We are master, sending elected event")
		m.sendEvent(ElectedEvent)
	} else {
		log.Debug("We are not master, sending demoted event")
		m.sendEvent(DemotedEvent)
	}
}

// onTopologyChange is called by the watcher. It is called in one of the 3 following cases:
// 1. There is a node that joined the pool of nodes interested by this IP
// 2. There is a node that left the pool of nodes interested by this IP
// 3. The current master node is trying to initiate a failover
func (m *EndpointManager) onTopologyChange(ctx context.Context) {
	log := logger.Get(ctx)
	if m.isStopped() {
		return
	}

	// If we are already master we do not want to failover because in the 3 cases possible for this method to be called
	// 1. We already are master and our lock is still valid, no reason to refresh it.
	// 2. Another host left the pool but since we are master, the one which left was standby and no action is needed.
	// 3. If there was a failover, it was initiated by the master node (us) and we do not want to take back the lock.
	if m.Status() == ACTIVATED {
		log.Info("Network topology changed but we are master, no actions needed")
		return
	}

	log.Info("Network topology changed, trying to get the IP")
	m.tryToGetIP(ctx)
}

func (m *EndpointManager) isStopped() bool {
	m.stopMutex.RLock()
	defer m.stopMutex.RUnlock()
	return m.stopped
}
