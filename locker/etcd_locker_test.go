package locker

import (
	"context"
	"testing"
	"time"

	"github.com/Scalingo/link/config"
	"github.com/Scalingo/link/etcdmock"
	"github.com/Scalingo/link/models"
	"github.com/Scalingo/link/models/modelsmock"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/etcdserver/api/v3rpc/rpctypes"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRefresh(t *testing.T) {
	key := "/test"
	examples := []struct {
		Name             string
		InitialLeaseID   int
		LastLeaseRefresh time.Time
		ExpectedLeaseID  int64
		ExpectedKV       func(*gomock.Controller, *etcdmock.MockKV)
		ExpectedLease    func(*etcdmock.MockLease)
		ExpectedStorage  func(*modelsmock.MockStorage)
		ExpectedError    string
	}{
		{
			Name:           "When we cannot get the lease",
			InitialLeaseID: 0,
			ExpectedLease: func(mock *etcdmock.MockLease) {
				mock.EXPECT().Grant(gomock.Any(), int64(15)).Return(nil, errors.New("NOP"))
			},
			ExpectedError: "NOP",
		}, {
			Name:           "When the transaction fails and no lease were configured",
			InitialLeaseID: 0,
			ExpectedLease: func(mock *etcdmock.MockLease) {
				mock.EXPECT().Grant(gomock.Any(), int64(15)).Return(&clientv3.LeaseGrantResponse{
					ID: 12,
				}, nil)
			},
			ExpectedStorage: func(mock *modelsmock.MockStorage) {
				mock.EXPECT().UpdateIP(gomock.Any(), gomock.Any()).Return(nil).Do(func(ctx context.Context, ip models.IP) {
					assert.Equal(t, ip.LeaseID, int64(12))
				})
			},
			ExpectedKV: func(ctrl *gomock.Controller, mock *etcdmock.MockKV) {
				txnMock := etcdmock.NewMockTxn(ctrl)
				txnMock.EXPECT().If(clientv3.Compare(clientv3.CreateRevision(key), "=", 0)).Return(txnMock)
				txnMock.EXPECT().Then(clientv3.OpPut(key, "locked", clientv3.WithLease(12))).Return(txnMock)
				txnMock.EXPECT().Commit().Return(nil, errors.New("NOP"))

				mock.EXPECT().Txn(gomock.Any()).Return(txnMock)
			},
			ExpectedError:   "NOP",
			ExpectedLeaseID: 0,
		}, {
			Name:             "When the transaction fails and the lease was not expired",
			InitialLeaseID:   0,
			LastLeaseRefresh: time.Now(),
			ExpectedLease: func(mock *etcdmock.MockLease) {
				mock.EXPECT().Grant(gomock.Any(), int64(15)).Return(&clientv3.LeaseGrantResponse{
					ID: 12,
				}, nil)
			},
			ExpectedStorage: func(mock *modelsmock.MockStorage) {
				mock.EXPECT().UpdateIP(gomock.Any(), gomock.Any()).Return(nil).Do(func(ctx context.Context, ip models.IP) {
					assert.Equal(t, ip.LeaseID, int64(12))
				})
			},
			ExpectedKV: func(ctrl *gomock.Controller, mock *etcdmock.MockKV) {
				txnMock := etcdmock.NewMockTxn(ctrl)
				txnMock.EXPECT().If(clientv3.Compare(clientv3.CreateRevision(key), "=", 0)).Return(txnMock)
				txnMock.EXPECT().Then(clientv3.OpPut(key, "locked", clientv3.WithLease(12))).Return(txnMock)
				txnMock.EXPECT().Commit().Return(nil, errors.New("NOP"))

				mock.EXPECT().Txn(gomock.Any()).Return(txnMock)
			},
			ExpectedError:   "NOP",
			ExpectedLeaseID: 12,
		}, {
			Name:             "When the transaction fails and the lease was expired",
			InitialLeaseID:   0,
			LastLeaseRefresh: time.Now().Add(-1 * time.Hour),
			ExpectedLease: func(mock *etcdmock.MockLease) {
				mock.EXPECT().Grant(gomock.Any(), int64(15)).Return(&clientv3.LeaseGrantResponse{
					ID: 12,
				}, nil)
			},
			ExpectedStorage: func(mock *modelsmock.MockStorage) {
				mock.EXPECT().UpdateIP(gomock.Any(), gomock.Any()).Return(nil).Do(func(ctx context.Context, ip models.IP) {
					assert.Equal(t, ip.LeaseID, int64(12))
				})
			},
			ExpectedKV: func(ctrl *gomock.Controller, mock *etcdmock.MockKV) {
				txnMock := etcdmock.NewMockTxn(ctrl)
				txnMock.EXPECT().If(clientv3.Compare(clientv3.CreateRevision(key), "=", 0)).Return(txnMock)
				txnMock.EXPECT().Then(clientv3.OpPut(key, "locked", clientv3.WithLease(12))).Return(txnMock)
				txnMock.EXPECT().Commit().Return(nil, errors.New("NOP"))

				mock.EXPECT().Txn(gomock.Any()).Return(txnMock)
			},
			ExpectedError:   "NOP",
			ExpectedLeaseID: 0,
		}, {
			Name:           "When keepalive fail",
			InitialLeaseID: 0,
			ExpectedLease: func(mock *etcdmock.MockLease) {
				mock.EXPECT().Grant(gomock.Any(), int64(15)).Return(&clientv3.LeaseGrantResponse{
					ID: 12,
				}, nil)

				mock.EXPECT().KeepAliveOnce(gomock.Any(), clientv3.LeaseID(12)).Return(nil, errors.New("NOP"))
			},
			ExpectedStorage: func(mock *modelsmock.MockStorage) {
				mock.EXPECT().UpdateIP(gomock.Any(), gomock.Any()).Return(nil).Do(func(ctx context.Context, ip models.IP) {
					assert.Equal(t, ip.LeaseID, int64(12))
				})
			},
			ExpectedKV: func(ctrl *gomock.Controller, mock *etcdmock.MockKV) {
				txnMock := etcdmock.NewMockTxn(ctrl)
				txnMock.EXPECT().If(clientv3.Compare(clientv3.CreateRevision(key), "=", 0)).Return(txnMock)
				txnMock.EXPECT().Then(clientv3.OpPut(key, "locked", clientv3.WithLease(12))).Return(txnMock)
				txnMock.EXPECT().Commit().Return(nil, nil)

				mock.EXPECT().Txn(gomock.Any()).Return(txnMock)
			},
			ExpectedLeaseID: 0,
		}, {
			Name:           "When keepalive fail because lease is not found",
			InitialLeaseID: 123,
			ExpectedLease: func(mock *etcdmock.MockLease) {
				mock.EXPECT().KeepAliveOnce(gomock.Any(), clientv3.LeaseID(123)).Return(nil, rpctypes.ErrLeaseNotFound)
			},
			ExpectedKV: func(ctrl *gomock.Controller, mock *etcdmock.MockKV) {
				txnMock := etcdmock.NewMockTxn(ctrl)
				txnMock.EXPECT().If(clientv3.Compare(clientv3.CreateRevision(key), "=", 0)).Return(txnMock)
				txnMock.EXPECT().Then(clientv3.OpPut(key, "locked", clientv3.WithLease(123))).Return(txnMock)
				txnMock.EXPECT().Commit().Return(nil, nil)

				mock.EXPECT().Txn(gomock.Any()).Return(txnMock)
			},
			ExpectedLeaseID: 0,
		}, {
			Name:           "When everything succeed",
			InitialLeaseID: 0,
			ExpectedLease: func(mock *etcdmock.MockLease) {
				mock.EXPECT().Grant(gomock.Any(), int64(15)).Return(&clientv3.LeaseGrantResponse{
					ID: 12,
				}, nil)

				mock.EXPECT().KeepAliveOnce(gomock.Any(), clientv3.LeaseID(12)).Return(nil, nil)
			},
			ExpectedStorage: func(mock *modelsmock.MockStorage) {
				mock.EXPECT().UpdateIP(gomock.Any(), gomock.Any()).Return(nil).Do(func(ctx context.Context, ip models.IP) {
					assert.Equal(t, ip.LeaseID, int64(12))
				})
			},
			ExpectedKV: func(ctrl *gomock.Controller, mock *etcdmock.MockKV) {
				txnMock := etcdmock.NewMockTxn(ctrl)
				txnMock.EXPECT().If(clientv3.Compare(clientv3.CreateRevision(key), "=", 0)).Return(txnMock)
				txnMock.EXPECT().Then(clientv3.OpPut(key, "locked", clientv3.WithLease(12))).Return(txnMock)
				txnMock.EXPECT().Commit().Return(nil, nil)

				mock.EXPECT().Txn(gomock.Any()).Return(txnMock)
			},
			ExpectedLeaseID: 12,
		},
	}

	for _, example := range examples {
		t.Run(example.Name, func(t *testing.T) {
			ctx := context.Background()
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			kvMock := etcdmock.NewMockKV(ctrl)
			leaseMock := etcdmock.NewMockLease(ctrl)
			storageMock := modelsmock.NewMockStorage(ctrl)

			if example.ExpectedLease != nil {
				example.ExpectedLease(leaseMock)
			}

			if example.ExpectedKV != nil {
				example.ExpectedKV(ctrl, kvMock)
			}

			if example.ExpectedStorage != nil {
				example.ExpectedStorage(storageMock)
			}

			locker := &etcdLocker{
				kvEtcd:    kvMock,
				leaseEtcd: leaseMock,
				key:       key,
				config: config.Config{
					KeepAliveInterval: 3 * time.Second,
				},
				leaseID:          clientv3.LeaseID(example.InitialLeaseID),
				lastLeaseRefresh: example.LastLeaseRefresh,
				storage:          storageMock,
			}

			err := locker.Refresh(ctx)
			if len(example.ExpectedError) != 0 {
				require.Error(t, err)
				assert.Contains(t, err.Error(), example.ExpectedError)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, locker.leaseID, clientv3.LeaseID(example.ExpectedLeaseID))
			// Let the coroutine time to finish
			time.Sleep(200 * time.Millisecond)
		})
	}
}

func Test_IsMaster(t *testing.T) {

	examples := []struct {
		Name           string
		OurLeaseID     clientv3.LeaseID
		CurrentLeaseID clientv3.LeaseID
		EtcdError      error
		Expected       bool
		ExpectedError  string
	}{
		{
			Name:          "When there is an issue with etcd",
			EtcdError:     errors.New("NOP"),
			ExpectedError: "NOP",
		}, {
			Name:           "when we are not master",
			OurLeaseID:     10,
			CurrentLeaseID: 11,
			Expected:       false,
		}, {
			Name:           "when we are master",
			OurLeaseID:     10,
			CurrentLeaseID: 10,
			Expected:       true,
		},
	}

	for _, example := range examples {
		t.Run(example.Name, func(t *testing.T) {
			key := "/test"
			ctx := context.Background()

			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			etcdMock := etcdmock.NewMockKV(ctrl)
			etcdMock.EXPECT().Get(gomock.Any(), key).Return(&clientv3.GetResponse{
				Kvs: []*mvccpb.KeyValue{
					{Lease: int64(example.CurrentLeaseID)},
				},
			}, example.EtcdError)

			locker := &etcdLocker{
				kvEtcd:  etcdMock,
				leaseID: example.OurLeaseID,
				key:     key,
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
