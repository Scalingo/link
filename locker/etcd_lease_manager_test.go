package locker

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.etcd.io/etcd/api/v3/v3rpc/rpctypes"
	etcdv3 "go.etcd.io/etcd/client/v3"

	"github.com/Scalingo/link/v3/config"
	"github.com/Scalingo/link/v3/etcdmock"
)

type updadeResults struct {
	called     bool
	oldLeaseID etcdv3.LeaseID
	newLeaseID etcdv3.LeaseID
}

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
		ShouldForceLeaseRefresh   bool
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
				mock.EXPECT().Grant(gomock.Any(), int64(15)).Return(&etcdv3.LeaseGrantResponse{
					ID: etcdv3.LeaseID(1234),
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
				mock.EXPECT().Grant(gomock.Any(), int64(15)).Return(&etcdv3.LeaseGrantResponse{
					ID: etcdv3.LeaseID(1235),
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
				mock.EXPECT().KeepAliveOnce(gomock.Any(), etcdv3.LeaseID(1234)).Return(nil, rpctypes.ErrLeaseNotFound)
			},
			ExpectedLeaseID:         1234,
			ShouldNotifyLeaseError:  true,
			ShouldForceLeaseRefresh: true,
		},
		{
			Name:                   "When the lease bas been generated but Etcd refuse the KeepAlive",
			InitialLeaseID:         1234,
			InitialLastRefreshedAt: time.Now(),
			MockLease: func(mock *etcdmock.MockLease) {
				mock.EXPECT().KeepAliveOnce(gomock.Any(), etcdv3.LeaseID(1234)).Return(nil, errors.New("nop"))
			},
			ExpectedLeaseID: 1234,
			ExpectedError:   "nop",
		},
		{
			Name:                   "When the lease has been generated and Etcd accept the keepalive",
			InitialLeaseID:         1234,
			InitialLastRefreshedAt: time.Now(),
			MockLease: func(mock *etcdmock.MockLease) {
				mock.EXPECT().KeepAliveOnce(gomock.Any(), etcdv3.LeaseID(1234)).Return(nil, nil)
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
				leaseID:            etcdv3.LeaseID(example.InitialLeaseID),
				lastRefreshedAt:    example.InitialLastRefreshedAt,
				leaseLock:          &sync.RWMutex{},
				callbackLock:       &sync.RWMutex{},
				leaseErrorNotifier: make(chan bool, 10),
				callbacks:          make(map[string]LeaseChangedCallback),
			}

			leaseChangedChan := make(chan updadeResults)

			go func() {
				time.Sleep(100 * time.Millisecond)
				// If the LeaseChanged callback has not been sent, send default
				select {
				case leaseChangedChan <- updadeResults{
					called: false,
				}:
				default:
				}
			}()

			_, err := leaseManager.SubscribeToLeaseChange(ctx, func(_ context.Context, o, n etcdv3.LeaseID) {
				leaseChangedChan <- updadeResults{
					called:     true,
					newLeaseID: n,
					oldLeaseID: o,
				}
			})
			require.NoError(t, err)

			err = leaseManager.refresh(ctx)
			if example.ExpectedError == "" {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				assert.Contains(t, err.Error(), example.ExpectedError)
			}

			assert.Equal(t, etcdv3.LeaseID(example.ExpectedLeaseID), leaseManager.leaseID)
			if example.ShouldNotifyLeaseError {
				assert.Len(t, leaseManager.leaseErrorNotifier, 1)
			} else {
				assert.Empty(t, leaseManager.leaseErrorNotifier)
			}

			// Wait for the corountine to be called
			result := <-leaseChangedChan
			assert.Equal(t, example.ShouldCallLeaseSubscriber, result.called)
			assert.Equal(t, etcdv3.LeaseID(example.LeaseSubscriberNew), result.newLeaseID)
			assert.Equal(t, etcdv3.LeaseID(example.LeaseSubscriberOld), result.oldLeaseID)
			assert.Equal(t, example.ShouldForceLeaseRefresh, leaseManager.forceLeaseRefresh)
		})
	}
}
