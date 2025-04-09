package ip

import (
	"context"
	"time"

	"github.com/pkg/errors"

	"github.com/Scalingo/go-utils/logger"
	"github.com/Scalingo/go-utils/retry"
)

var (
	// ErrIsNotMaster is an error sent by Failover when the node is not currently master
	ErrIsNotMaster = errors.New("this node is not master of this endpoint")

	// ErrNoOtherHosts is an error sent by Failover when there is no other node to fail over.
	ErrNoOtherHosts = errors.New("no other nodes are listening for this endpoint")

	// ErrReallocationTimedOut is an error returned by waitForReallocation if the reallocation did not happen in less than KeepAliveInterval
	ErrReallocationTimedOut = errors.New("reallocation timed out")
)

func (m *EndpointManager) waitForReallocation(ctx context.Context) error {
	log := logger.Get(ctx).WithField("process", "wait-for-reallocation")

	retryer := retry.New(retry.WithMaxDuration(m.config.KeepAliveInterval),
		retry.WithWaitDuration(100*time.Millisecond),
		retry.WithoutMaxAttempts())

	err := retryer.Do(ctx, func(ctx context.Context) error {
		isMaster, err := m.locker.IsMaster(ctx)
		if err != nil {
			log.WithError(err).Debug("Fail to check if we are still master")
			return err
		}

		if isMaster {
			log.Debug("We are still master")
			return errors.New("still master")
		}
		return nil
	})

	if err != nil {
		if retryErr, ok := err.(retry.RetryError); ok {
			if retryErr.Scope == retry.MaxDurationScope {
				return ErrReallocationTimedOut
			}
		}

		return errors.Wrap(err, "fail to wait for endpoint reallocation")
	}
	return nil
}

// Failover forces a change of the master. It can only be run on the current master instance for an endpoint.
// If there is another node available for this endpoint, it steps down as a master and ensure that another node becomes master.
// This function refuses to trigger a failover if the node is not master or if there are no other nodes.
// To trigger the failover, we unlock the endpoint (remove the lock key) and update the link between the host and the endpoint.
// Updating the link notifies watchers on this endpoint and other hosts will try to get the endpoint.
func (m *EndpointManager) Failover(ctx context.Context) error {
	isMaster, err := m.locker.IsMaster(ctx)
	if err != nil {
		return errors.Wrap(err, "fail to check if the node is master")
	}
	if !isMaster {
		return ErrIsNotMaster
	}
	hosts, err := m.storage.GetEndpointHosts(ctx, m.plugin.ElectionKey(ctx))
	if err != nil {
		return errors.Wrap(err, "fail to list other nodes listening for this endpoint")
	}
	if len(hosts) <= 1 {
		return ErrNoOtherHosts
	}

	// Unlock the endpoint
	err = m.locker.Unlock(ctx)
	if err != nil {
		return errors.Wrap(err, "fail to unlock current endpoint")
	}

	// Update the endpoint link that will trigger every other watchers on this endpoint
	err = m.storage.LinkEndpointWithCurrentHost(ctx, m.plugin.ElectionKey(ctx))
	if err != nil {
		return errors.Wrap(err, "fail to update the endpoint Link")
	}
	return nil
}
