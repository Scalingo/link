package locker

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/Scalingo/link/config"
	"github.com/Scalingo/link/etcdmock"
	"github.com/coreos/etcd/etcdserver/api/v3rpc/rpctypes"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.etcd.io/etcd/clientv3"
)

func Test_refresh(t *testing.T) {
	examples := []struct {
		Name                      string
		InitialLeaseID            int64
		InitialLastRefreshedAt    time.Time
		MockLease                 func(*etcdmock.MockLease)
		ExpectedLastRefreshedAt   time.Time
		ExpectedLeaseID           int64
		ExpectedError             string
		ShouldNotifyLeaseError    bool
		ShouldCallLeaseSubscriber bool
		LeaseSubscriberOld        int64
		LeaseSubscriberNew        int64
	}{
		{
			Name:           "When the lease has not been generated and we fail to to create a lease",
			InitialLeaseID: 0,
			MockLease: func(mock *etcdmock.MockLease) {
				mock.EXPECT().Grant(gomock.Any(), int64(15)).Return(nil, errors.New("nop"))
			},
			ExpectedLeaseID: 0,
			ExpectedError:   "fail to regenerate lease",
		},
		{
			Name: "When the lease has not been generated",
			MockLease: func(mock *etcdmock.MockLease) {
				mock.EXPECT().Grant(gomock.Any(), int64(15)).Return(&clientv3.LeaseGrantResponse{
					ID: clientv3.LeaseID(1234),
				}, nil)
			},
			ExpectedLeaseID:           1234,
			ShouldCallLeaseSubscriber: true,
			LeaseSubscriberNew:        1234,
			LeaseSubscriberOld:        0,
		},
		{
			Name:                   "When the lease has expired",
			InitialLeaseID:         1234,
			InitialLastRefreshedAt: time.Now().Add(-1 * time.Hour),
			MockLease: func(mock *etcdmock.MockLease) {
				mock.EXPECT().Grant(gomock.Any(), int64(15)).Return(&clientv3.LeaseGrantResponse{
					ID: clientv3.LeaseID(1235),
				}, nil)
			},
			ExpectedLeaseID:           1235,
			ShouldCallLeaseSubscriber: true,
			LeaseSubscriberNew:        1235,
			LeaseSubscriberOld:        1234,
		},
		{
			Name:                   "When the lease has been generated but Etcd refuse the KeepAlive (not found)",
			InitialLeaseID:         1234,
			InitialLastRefreshedAt: time.Now(),
			MockLease: func(mock *etcdmock.MockLease) {
				mock.EXPECT().KeepAliveOnce(gomock.Any(), clientv3.LeaseID(1234)).Return(nil, rpctypes.ErrLeaseNotFound)
			},
			ExpectedLeaseID:        0,
			ShouldNotifyLeaseError: true,
		},
		{
			Name:                   "When the lease bas been generated but Etcd refuse the KeepAlive",
			InitialLeaseID:         1234,
			InitialLastRefreshedAt: time.Now(),
			MockLease: func(mock *etcdmock.MockLease) {
				mock.EXPECT().KeepAliveOnce(gomock.Any(), clientv3.LeaseID(1234)).Return(nil, errors.New("nop"))
			},
			ExpectedLeaseID: 1234,
			ExpectedError:   "nop",
		},
		{
			Name:                   "When the lease has been generated and Etcd accept the keepalive",
			InitialLeaseID:         1234,
			InitialLastRefreshedAt: time.Now(),
			MockLease: func(mock *etcdmock.MockLease) {
				mock.EXPECT().KeepAliveOnce(gomock.Any(), clientv3.LeaseID(1234)).Return(nil, nil)
			},
			ExpectedLeaseID: 1234,
		},
	}
	for _, example := range examples {
		t.Run(example.Name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			ctx := context.Background()
			leaseMock := etcdmock.NewMockLease(ctrl)

			if example.MockLease != nil {
				example.MockLease(leaseMock)
			}

			leaseManager := &etcdLeaseManager{
				leases: leaseMock,
				config: config.Config{
					KeepAliveInterval: 3 * time.Second,
				},
				leaseID:            clientv3.LeaseID(example.InitialLeaseID),
				lastRefreshedAt:    example.InitialLastRefreshedAt,
				leaseLock:          &sync.RWMutex{},
				callbackLock:       &sync.RWMutex{},
				leaseErrorNotifier: make(chan bool, 10),
				callbacks:          make(map[string]LeaseChangedCallback),
			}

			subscriberCalled := false
			var oldLeaseID, newLeaseID clientv3.LeaseID

			leaseManager.SubscribeToLeaseChange(ctx, func(_ context.Context, o, n clientv3.LeaseID) {
				subscriberCalled = true
				newLeaseID = n
				oldLeaseID = o
			})
			err := leaseManager.refresh(ctx)
			if example.ExpectedError == "" {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				assert.Contains(t, err.Error(), example.ExpectedError)
			}

			assert.Equal(t, clientv3.LeaseID(example.ExpectedLeaseID), leaseManager.leaseID)
			if example.ShouldNotifyLeaseError {
				assert.Len(t, leaseManager.leaseErrorNotifier, 1)
			} else {
				assert.Len(t, leaseManager.leaseErrorNotifier, 0)
			}

			// Wait for the corountine to be called
			time.Sleep(100 * time.Millisecond)
			assert.Equal(t, example.ShouldCallLeaseSubscriber, subscriberCalled)
			assert.Equal(t, clientv3.LeaseID(example.LeaseSubscriberNew), newLeaseID)
			assert.Equal(t, clientv3.LeaseID(example.LeaseSubscriberOld), oldLeaseID)
		})
	}
}
