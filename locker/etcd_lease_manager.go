package locker

import (
	"context"
	"sync"
	"time"

	"github.com/gofrs/uuid/v5"
	"go.etcd.io/etcd/api/v3/v3rpc/rpctypes"
	etcdv3 "go.etcd.io/etcd/client/v3"

	"github.com/Scalingo/go-utils/errors/v2"
	"github.com/Scalingo/go-utils/logger"
	"github.com/Scalingo/link/v2/config"
	"github.com/Scalingo/link/v2/models"
)

const DataVersion = 3

// ErrCallbackNotFound is launched when a user tries to delete a callback that does not exist
var ErrCallbackNotFound = errors.New(context.Background(), "lease callback not found")

// ErrGetLeaseTimeout is launched when a user calls GetLease and we fail to provide one in time
var ErrGetLeaseTimeout = errors.New(context.Background(), "timeout while trying to get lease")

// LeaseChangedCallback is a callback called by the lease manager when the leaseID has changed so that all managers could try to regenerate their keys
type LeaseChangedCallback func(ctx context.Context, oldLeaseID, newLeaseID etcdv3.LeaseID)

// EtcdLeaseManager let you get the current server lease for the server
type EtcdLeaseManager interface {
	Start(ctx context.Context) error
	Stop(ctx context.Context)
	GetLease(ctx context.Context) (etcdv3.LeaseID, error)                                      // This will get the current lease or wait for one to be generated
	MarkLeaseAsDirty(ctx context.Context, leaseID etcdv3.LeaseID) error                        // This is sent by clients if they think that there might be an issue with the Lease
	SubscribeToLeaseChange(ctx context.Context, callback LeaseChangedCallback) (string, error) // Subscribe to lease changes. This function returns an ID that should be used to unsubscribe.
	UnsubscribeToLeaseChange(ctx context.Context, id string) error                             // Unsubscribe from the lease changes
}

type etcdLeaseManager struct {
	stopper            chan bool
	config             config.Config
	leases             etcdv3.Lease
	kv                 etcdv3.KV
	storage            models.Storage
	leaseID            etcdv3.LeaseID
	callbacks          map[string]LeaseChangedCallback
	lastRefreshedAt    time.Time
	leaseDirtyNotifier chan etcdv3.LeaseID
	leaseErrorNotifier chan bool
	forceLeaseRefresh  bool
	callbackLock       *sync.RWMutex
	leaseLock          *sync.RWMutex
}

// NewEtcdLeaseManager returns a default manager that implements the EtcdLeaseManager interface
func NewEtcdLeaseManager(_ context.Context, config config.Config, storage models.Storage, etcd *etcdv3.Client) EtcdLeaseManager {
	return &etcdLeaseManager{
		stopper:            make(chan bool, 1),
		leaseDirtyNotifier: make(chan etcdv3.LeaseID, 1),
		leaseErrorNotifier: make(chan bool, 1),
		leases:             etcd,
		kv:                 etcd,
		config:             config,
		storage:            storage,
		forceLeaseRefresh:  false,
		callbacks:          make(map[string]LeaseChangedCallback),
		callbackLock:       &sync.RWMutex{},
		leaseLock:          &sync.RWMutex{},
	}
}

func (m *etcdLeaseManager) GetLease(ctx context.Context) (etcdv3.LeaseID, error) {
	log := logger.Get(ctx)
	// If the lease has been generated, send it
	m.leaseLock.RLock()
	leaseID := m.leaseID
	m.leaseLock.RUnlock()
	if leaseID != 0 {
		log.Debug("Lease has already been generated")
		return leaseID, nil
	}

	log.Debug("Generating a new lease")
	// If the lease has not been generated yet (or is dirty)
	// Prepare the return channel
	leaseChan := make(chan etcdv3.LeaseID, 1)

	// Use our own subscribe mechanism to know when the new lease has been generated
	id, err := m.SubscribeToLeaseChange(ctx, func(_ context.Context, _, leaseID etcdv3.LeaseID) {
		leaseChan <- leaseID
	})
	if err != nil {
		return etcdv3.NoLease, errors.Wrap(ctx, err, "fail to subscribe to leaseID changes")
	}
	defer func() {
		err := m.UnsubscribeToLeaseChange(ctx, id) // Do not forget to clean it
		if err != nil && err != ErrCallbackNotFound {
			log.WithError(err).Error("fail to unsubscribe from lease changes")
		}
	}()

	// The timer is a safeguard to prevent a race condition between the initial lease generation and the GetLease calls. By waiting max twice the KeepAliveInterval, we are safe.
	timer := time.NewTimer(2 * m.config.KeepAliveInterval)
	select {
	case <-timer.C:
		// If the command timed out
		return etcdv3.NoLease, ErrGetLeaseTimeout
	case leaseID = <-leaseChan:
		log.Debug("Lease has been generated in time")
	}
	return leaseID, nil // We got the lease in time \o/
}

func (m *etcdLeaseManager) MarkLeaseAsDirty(_ context.Context, leaseID etcdv3.LeaseID) error {
	m.leaseDirtyNotifier <- leaseID
	return nil
}

func (m *etcdLeaseManager) SubscribeToLeaseChange(ctx context.Context, callback LeaseChangedCallback) (string, error) {
	if callback == nil {
		return "", errors.New(ctx, "callback cannot be nil")
	}

	log := logger.Get(ctx)
	log.Debug("Subscribe to lease change")
	uuid, err := uuid.NewV4()
	if err != nil {
		return "", errors.Wrap(ctx, err, "fail to generate UUID")
	}

	m.callbackLock.Lock()
	defer m.callbackLock.Unlock()
	id := uuid.String()
	m.callbacks[id] = callback
	return id, nil
}

func (m *etcdLeaseManager) UnsubscribeToLeaseChange(_ context.Context, id string) error {
	m.callbackLock.Lock()
	defer m.callbackLock.Unlock()

	_, ok := m.callbacks[id]
	if !ok {
		return ErrCallbackNotFound
	}
	delete(m.callbacks, id)

	return nil
}

func (m *etcdLeaseManager) Start(ctx context.Context) error {
	log := logger.Get(ctx).WithField("source", "etcd-lease-manager")

	// Step 1: Fetch the lease associated to the host
	// We do this to keep lease between restart. If we do not do that, we
	// might trigger unwanted failover since old lease will expire and we do
	// not have any guarantee that we are the one that will take those locks.
	log.Info("Getting old leaseID")
	host, err := m.storage.GetCurrentHost(ctx)
	if err != nil && !errors.Is(err, models.ErrHostNotFound) {
		return errors.Wrap(ctx, err, "fail to find host config")
	}
	if host.LeaseID != 0 {
		log.Infof("Starting with LeaseID=%v", host.LeaseID)
		m.leaseID = etcdv3.LeaseID(host.LeaseID)
		m.lastRefreshedAt = time.Now()
	} else {
		log.Info("LeaseID not found, starting with LeaseID=0")
		m.forceLeaseRefresh = true
		m.leaseDirtyNotifier <- m.leaseID
	}

	_, err = m.SubscribeToLeaseChange(ctx, m.storeLeaseChange)
	if err != nil {
		return errors.Wrap(ctx, err, "fail to subscribe to lease changes")
	}

	// Step 2: Start the lease refresher. This codes has to be running constantly to keep the lease
	// alive. If this loop stop (or fails for a long time), the other nodes will try to get the lock
	// and we will loose our endpoints.
	go func() {
		log := log.WithField("source", "etcd-lease-manager-refresh")
		ctx := logger.ToCtx(ctx, log)
		log.Info("Starting lease refresher")
		ticker := time.NewTicker(m.config.KeepAliveInterval)
		for {
			// Should we skip the refresher in this loop iteration ?
			// This is used if a client marked a lease as dirty even if it wasn't
			runRefresher := true

			select {
			case <-ticker.C:
			case leaseID := <-m.leaseDirtyNotifier:
				runRefresher = !m.isLeaseDirty(ctx, leaseID)
				log.Debugf("A lease is dirty. Is the current one dirty? %v", runRefresher)
			case <-m.leaseErrorNotifier:
				log.Debug("Notified of an error in the refresh process, retry immediately")
			case <-m.stopper:
				log.Info("Stopping lease refresher")
				ticker.Stop()
				return
			}

			if runRefresher {
				err := m.refresh(ctx)
				if err != nil {
					log.WithError(err).Error("fail to refresh lease")
				}
			}
		}
	}()
	return nil
}

// This method is called to refresh the current lease. If the currentLease is dirty or if it has not be generated: generate a new lease
func (m *etcdLeaseManager) refresh(ctx context.Context) error {
	log := logger.Get(ctx).WithField("source", "etcd-lease-manager")

	m.leaseLock.Lock()
	defer m.leaseLock.Unlock()

	// If the lease has not been generated yet (or if it is dirty)
	if m.leaseID == 0 || m.forceLeaseRefresh || m.hasLeaseExpired(ctx) {
		switch {
		case m.forceLeaseRefresh:
			log.Info("Lease is dirty, regenerating lease")
		case m.leaseID == 0:
			log.Info("LeaseID = 0, regenerating lease")
		default:
			log.Info("Lease has expired, regenerating lease")
		}

		oldLeaseID := m.leaseID
		leaseTime := int64(m.config.LeaseTime().Seconds())
		grant, err := m.leases.Grant(ctx, leaseTime)
		if err != nil {
			return errors.Wrap(ctx, err, "fail to regenerate lease")
		}
		log.WithField("lease_id", grant.ID).Info("New LeaseID generated")

		m.leaseID = grant.ID
		m.lastRefreshedAt = time.Now()
		m.forceLeaseRefresh = false
		m.notifyLeaseChanged(ctx, oldLeaseID, m.leaseID)
		return nil
	}

	// Here the lease is still valid, we just need to refresh it.
	_, err := m.leases.KeepAliveOnce(ctx, m.leaseID)
	if err != nil {
		if err, ok := err.(rpctypes.EtcdError); ok && rpctypes.Error(err) == rpctypes.ErrLeaseNotFound {
			log.WithError(err).Error("Keep alive failed: lease not found, regenerate lease")
			m.forceLeaseRefresh = true
			m.leaseErrorNotifier <- true
			return nil
		}
		return errors.Wrap(ctx, err, "keep alive failed but the lease might still be valid, continuing")
	}
	m.lastRefreshedAt = time.Now()

	return nil
}

func (m *etcdLeaseManager) Stop(_ context.Context) {
	m.stopper <- true
}

func (m *etcdLeaseManager) notifyLeaseChanged(ctx context.Context, oldLeaseID, newLeaseID etcdv3.LeaseID) {
	m.callbackLock.RLock()
	defer m.callbackLock.RUnlock()

	for _, callback := range m.callbacks {
		go callback(ctx, oldLeaseID, newLeaseID)
	}
}

func (m *etcdLeaseManager) hasLeaseExpired(_ context.Context) bool {
	return m.lastRefreshedAt.IsZero() || time.Now().After(m.lastRefreshedAt.Add(m.config.LeaseTime()))
}

func (m *etcdLeaseManager) isLeaseDirty(ctx context.Context, leaseID etcdv3.LeaseID) bool {
	log := logger.Get(ctx)
	m.leaseLock.RLock()
	defer m.leaseLock.RUnlock()

	if leaseID != m.leaseID {
		log.Infof("We got notified that there was an issue with lease %v but current lease is %v", leaseID, m.leaseID)
		// This lease is not the current one so there is nothing to do
		return false
	}
	// If the key has not expired, there's nothing to do. If there is an issue the refresher will pick that up and mange it by itself.
	if m.hasLeaseExpired(ctx) {
		log.Infof("We got notified that there was an issue with lease %v generated on %v and indeed it's expired. Resetting it.", leaseID, m.lastRefreshedAt)
		m.forceLeaseRefresh = true
		return true
	}

	log.Infof("We got notified that there was an issue with lease %v but we've found no issue", m.leaseID)
	return false
}

func (m *etcdLeaseManager) storeLeaseChange(ctx context.Context, _, leaseID etcdv3.LeaseID) {
	log := logger.Get(ctx)
	err := m.storage.SaveHost(ctx, models.Host{
		Hostname:    m.config.Hostname,
		LeaseID:     int64(leaseID),
		DataVersion: DataVersion,
	})
	if err != nil {
		log.WithError(err).Error("Fail to save new lease")
	}
}
