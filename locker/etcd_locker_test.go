package locker

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/Scalingo/link/v2/config"
	"github.com/Scalingo/link/v2/etcdmock"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.etcd.io/etcd/v3/clientv3"
	"go.etcd.io/etcd/v3/mvcc/mvccpb"
)

func TestRefresh(t *testing.T) {
	key := "/test"
	examples := []struct {
		Name                      string
		LeaseSubscriberID         string
		ExpectedLeaseSubscriberID string
		ExpectedKV                func(*gomock.Controller, *etcdmock.MockKV)
		ExpectedLeaseManager      func(*MockEtcdLeaseManager)
		ExpectedError             string
	}{
		{
			Name:                      "When we: did not subscribe to change yet and the transaction succeed",
			LeaseSubscriberID:         "",
			ExpectedLeaseSubscriberID: "id-1",
			ExpectedLeaseManager: func(mock *MockEtcdLeaseManager) {
				mock.EXPECT().SubscribeToLeaseChange(gomock.Any(), gomock.Any()).Return("id-1", nil)
				mock.EXPECT().GetLease(gomock.Any()).Return(clientv3.LeaseID(12), nil)
			},
			ExpectedKV: func(ctrl *gomock.Controller, mock *etcdmock.MockKV) {
				txnMock := etcdmock.NewMockTxn(ctrl)
				txnMock.EXPECT().If(clientv3.Compare(clientv3.CreateRevision(key), "=", 0)).Return(txnMock)
				txnMock.EXPECT().Then(clientv3.OpPut(key, "hostname", clientv3.WithLease(12))).Return(txnMock)
				txnMock.EXPECT().Commit().Return(nil, nil)
				mock.EXPECT().Txn(gomock.Any()).Return(txnMock)
			},
		},
		{
			Name:                      "When the transaction fails",
			LeaseSubscriberID:         "id-1",
			ExpectedLeaseSubscriberID: "id-1",
			ExpectedLeaseManager: func(mock *MockEtcdLeaseManager) {
				mock.EXPECT().GetLease(gomock.Any()).Return(clientv3.LeaseID(12), nil)
				mock.EXPECT().MarkLeaseAsDirty(gomock.Any(), clientv3.LeaseID(12)).Return(nil)
			},
			ExpectedKV: func(ctrl *gomock.Controller, mock *etcdmock.MockKV) {
				txnMock := etcdmock.NewMockTxn(ctrl)
				txnMock.EXPECT().If(clientv3.Compare(clientv3.CreateRevision(key), "=", 0)).Return(txnMock)
				txnMock.EXPECT().Then(clientv3.OpPut(key, "hostname", clientv3.WithLease(12))).Return(txnMock)
				txnMock.EXPECT().Commit().Return(nil, errors.New("NOP"))
				mock.EXPECT().Txn(gomock.Any()).Return(txnMock)
			},
			ExpectedError: "NOP",
		},
	}

	for _, example := range examples {
		t.Run(example.Name, func(t *testing.T) {
			ctx := context.Background()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			kvMock := etcdmock.NewMockKV(ctrl)
			leaseManagerMock := NewMockEtcdLeaseManager(ctrl)

			if example.ExpectedKV != nil {
				example.ExpectedKV(ctrl, kvMock)
			}

			if example.ExpectedLeaseManager != nil {
				example.ExpectedLeaseManager(leaseManagerMock)
			}

			locker := &etcdLocker{
				kvEtcd:            kvMock,
				key:               key,
				leaseManager:      leaseManagerMock,
				leaseSubscriberID: example.LeaseSubscriberID,
				config: config.Config{
					KeepAliveInterval: 3 * time.Second,
					Hostname:          "hostname",
				},
				lock: &sync.Mutex{},
			}

			err := locker.Refresh(ctx)
			if len(example.ExpectedError) != 0 {
				require.Error(t, err)
				assert.Contains(t, err.Error(), example.ExpectedError)
			} else {
				assert.NoError(t, err)
			}
			assert.Equal(t, example.ExpectedLeaseSubscriberID, locker.leaseSubscriberID)
		})
	}
}

func Test_IsMaster(t *testing.T) {

	examples := []struct {
		Name                 string
		CurrentLeaseID       clientv3.LeaseID
		MockEtcdLeaseManager func(mock *MockEtcdLeaseManager)
		EtcdError            error
		Expected             bool
		ExpectedError        string
	}{
		{
			Name:          "When there is an issue with etcd",
			EtcdError:     errors.New("NOP"),
			ExpectedError: "NOP",
		}, {
			Name: "when we are not master",
			MockEtcdLeaseManager: func(mock *MockEtcdLeaseManager) {
				mock.EXPECT().GetLease(gomock.Any()).Return(clientv3.LeaseID(10), nil)
			},
			CurrentLeaseID: 11,
			Expected:       false,
		}, {
			Name: "when we are master",
			MockEtcdLeaseManager: func(mock *MockEtcdLeaseManager) {
				mock.EXPECT().GetLease(gomock.Any()).Return(clientv3.LeaseID(10), nil)
			},
			CurrentLeaseID: 10,
			Expected:       true,
		}, {
			Name: "when there was an error while getting our leaseID",
			MockEtcdLeaseManager: func(mock *MockEtcdLeaseManager) {
				mock.EXPECT().GetLease(gomock.Any()).Return(clientv3.LeaseID(0), errors.New("Nop"))
			},
			ExpectedError: "Nop",
		},
	}

	for _, example := range examples {
		t.Run(example.Name, func(t *testing.T) {
			key := "/test"
			ctx := context.Background()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			etcdLeaseManager := NewMockEtcdLeaseManager(ctrl)
			if example.MockEtcdLeaseManager != nil {
				example.MockEtcdLeaseManager(etcdLeaseManager)
			}

			etcdMock := etcdmock.NewMockKV(ctrl)
			etcdMock.EXPECT().Get(gomock.Any(), key).Return(&clientv3.GetResponse{
				Kvs: []*mvccpb.KeyValue{
					{Lease: int64(example.CurrentLeaseID)},
				},
			}, example.EtcdError)

			locker := &etcdLocker{
				kvEtcd:       etcdMock,
				key:          key,
				leaseManager: etcdLeaseManager,
			}

			value, err := locker.IsMaster(ctx)
			if len(example.ExpectedError) > 0 {
				require.Error(t, err)
				assert.Contains(t, err.Error(), example.ExpectedError)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, example.Expected, value)
			}
		})
	}
}

func Test_Unlock(t *testing.T) {
	key := "/test"

	examples := []struct {
		Name               string
		CurrentLeaseID     clientv3.LeaseID
		ExpectLeaseManager func(mock *MockEtcdLeaseManager)
		ExpectKV           func(*etcdmock.MockKV)
		ExpectedError      string
	}{
		{
			Name: "when we are not master",
			ExpectLeaseManager: func(mock *MockEtcdLeaseManager) {
				mock.EXPECT().GetLease(gomock.Any()).Return(clientv3.LeaseID(10), nil)
			},
			CurrentLeaseID: 11,
			ExpectedError:  ErrNotMaster.Error(),
		}, {
			Name: "when we are master and etcd is not sending any error",
			ExpectLeaseManager: func(mock *MockEtcdLeaseManager) {
				mock.EXPECT().GetLease(gomock.Any()).Return(clientv3.LeaseID(10), nil)
			},
			CurrentLeaseID: 10,
			ExpectKV: func(m *etcdmock.MockKV) {
				m.EXPECT().Delete(gomock.Any(), key).Return(nil, nil)
			},
		}, {
			Name: "when we are master and etcd is sending an error",
			ExpectLeaseManager: func(mock *MockEtcdLeaseManager) {
				mock.EXPECT().GetLease(gomock.Any()).Return(clientv3.LeaseID(10), nil)
			},
			CurrentLeaseID: 10,
			ExpectKV: func(m *etcdmock.MockKV) {
				m.EXPECT().Delete(gomock.Any(), key).Return(nil, errors.New("HAHA NOPE!"))
			},
			ExpectedError: "HAHA NOPE!",
		},
	}

	for _, example := range examples {
		t.Run(example.Name, func(t *testing.T) {
			ctx := context.Background()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			etcdLeaseManager := NewMockEtcdLeaseManager(ctrl)
			if example.ExpectLeaseManager != nil {
				example.ExpectLeaseManager(etcdLeaseManager)
			}

			etcdMock := etcdmock.NewMockKV(ctrl)
			etcdMock.EXPECT().Get(gomock.Any(), key).Return(&clientv3.GetResponse{
				Kvs: []*mvccpb.KeyValue{
					{Lease: int64(example.CurrentLeaseID)},
				},
			}, nil)

			if example.ExpectKV != nil {
				example.ExpectKV(etcdMock)
			}

			locker := &etcdLocker{
				kvEtcd:       etcdMock,
				key:          key,
				leaseManager: etcdLeaseManager,
			}

			err := locker.Unlock(ctx)
			if len(example.ExpectedError) > 0 {
				require.Error(t, err)
				assert.Contains(t, err.Error(), example.ExpectedError)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
