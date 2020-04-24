package ip

import (
	"context"
	"time"

	"github.com/Scalingo/go-utils/logger"
	"github.com/Scalingo/link/locker"
	"github.com/pkg/errors"
)

var (
	// ErrReallocationTimedOut is an error returned by waitForReallocation is the reallocation did not happen in less than KeepAliveInterval
	ErrReallocationTimedOut = errors.New("Reallocation timed out")
)

/* Stop process:
  1. Fetch information about other hosts registered on this IP and our state on the IP
  3. Begin the shutdown process (/!\ THIS SHOULD BE PROTECTED BY THE STOP MUTEX OR IT COULD LEAD TO INVALID STATE)
	3.1. Set the stop flag to true, since we are in the lock it will not impact the other process yet but it will ensure that if the process stop unexpectedly, all the other goroutine will stop and we will remove the IP by letting it decay.
	3.2. If we are the owner of the lock on the IP remove it
	3.3. Remove us from the list of potential hosts for this IP (UnlinkIP). This will trigger other hosts to try to get the IP
	3.4. Set the state machine to a state where it will remove the IP.
	3.5. Remove the IP from the interface.
*/

func (m *manager) Stop(ctx context.Context) error {
	log := logger.Get(ctx).WithField("task", "stop")
	ctx = logger.ToCtx(ctx, log)

	log.Info("Start stop preflight checks")
	hosts, err := m.storage.IPHosts(ctx, m.IP())
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
	err = m.storage.UnlinkIP(ctx, m.ip)
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
	if m.stateMachine.Current() != FAILING {
		// We cannot use SendEvent here because the isStopping method is blocked by the stopMutex
		m.eventChan <- DemotedEvent
	}

	// We can stop the FSM we do not need it anymore
	close(m.eventChan)

	log.Info("Remove the IP from our interface")
	err = m.networkInterface.RemoveIP(m.ip.IP)
	if err != nil {
		return errors.Wrap(err, "fail to remove IP from interface")
	}

	log.Info("Stop process ended!")
	return nil
}

// ipCheckLoop checks every KeepaliveInterval if we should try to get the IP
// if we should try to get the IP, it will launch the tryToGetIP method that will do the heavy lifting.
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
	if m.stateMachine.Current() == FAILING {
		return
	}

	log := logger.Get(ctx)
	err := m.locker.Refresh(ctx)
	if err != nil {
		m.keepaliveRetry++
		log.WithError(err).Info("Fail to refresh lock (retry)")
		if m.keepaliveRetry > m.config.KeepAliveRetry {
			log.WithError(err).Error("Fail to refresh lock")
			m.sendEvent(FaultEvent)
		}
		return
	}
	if m.keepaliveRetry > 0 {
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

func (m *manager) onTopologyChange(ctx context.Context) {
	log := logger.Get(ctx)
	if m.isStopped() {
		return
	}

	log.Info("Network topology changed, trying to get the IP")
	m.tryToGetIP(ctx)
}

func (m *manager) waitForReallocation(ctx context.Context) error {
	log := logger.Get(ctx)
	startTime := time.Now()
	for {
		time.Sleep(100 * time.Millisecond)
		isMaster, err := m.locker.IsMaster(ctx)
		if err != nil && err == locker.ErrInvalidEtcdState { // The key does not exist so nobody took the lease yet
			continue
		}
		if err != nil {
			log.WithError(err).Error("Fail to check if we are master, retrying...")
		}

		if !isMaster {
			return nil // Someone else took the lease
		}

		if time.Now().Sub(startTime) > m.config.KeepAliveInterval {
			return ErrReallocationTimedOut
		}
	}
}
