package ip

import (
	"context"
	"time"

	"github.com/Scalingo/go-utils/logger"
	"github.com/Scalingo/link/locker"
	"github.com/pkg/errors"
)

var (
	// ErrIsNotMaster is an error sent by Failover when the node is not currently master
	ErrIsNotMaster = errors.New("This node is not master of this IP")

	// ErrNoOtherHosts is an error sent by Failover when there is no other node to fail over.
	ErrNoOtherHosts = errors.New("No other nodes are listening for this IP")
)

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

// Failover trigger a failover on the current IP
// This function will refuse to trigger a failover if the node is not master or if there are no other nodes.
// To trigger the failover, we will Unlock the IP (remove the lock key) and update the Link between the Host and the IP
// Updating the link will notify watchers on this IP and other hosts will try to get the IP.
func (m *manager) Failover(ctx context.Context) error {
	isMaster, err := m.locker.IsMaster(ctx)
	if err != nil {
		return errors.Wrap(err, "fail to check if the node is master")
	}
	if !isMaster {
		return ErrIsNotMaster
	}
	hosts, err := m.storage.IPHosts(ctx, m.IP())
	if err != nil {
		return errors.Wrap(err, "fail to list other nodes listening for this IP")
	}
	if len(hosts) <= 1 {
		return ErrNoOtherHosts
	}

	// Unlock the IP
	err = m.locker.Unlock(ctx)
	if err != nil {
		return errors.Wrap(err, "fail to unlock current IP")
	}

	// LinkIP will update the IP Link that will trigger every other Watchers on this IP
	err = m.storage.LinkIP(ctx, m.IP())
	if err != nil {
		return errors.Wrap(err, "fail to update the IP Link")
	}
	return nil
}
