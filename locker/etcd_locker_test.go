package locker

import (
	"context"
	"testing"
	"time"

	"github.com/Scalingo/link/config"
	"github.com/Scalingo/link/etcdmock"
	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/mvcc/mvccpb"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRefresh(t *testing.T) {
	key := "/test"
	examples := []struct {
		Name            string
		InitialLeaseID  int
		ExpectedLeaseID int64
		ExpectedKV      func(*gomock.Controller, *etcdmock.MockKV)
		ExpectedLease   func(*etcdmock.MockLease)
		ExpectedError   string
	}{
		{
			Name:           "When we cannot get the lease",
			InitialLeaseID: 0,
			ExpectedLease: func(mock *etcdmock.MockLease) {
				mock.EXPECT().Grant(gomock.Any(), int64(6)).Return(nil, errors.New("NOP"))
			},
			ExpectedError: "NOP",
		}, {
			Name:           "When the transaction fails",
			InitialLeaseID: 0,
			ExpectedLease: func(mock *etcdmock.MockLease) {
				mock.EXPECT().Grant(gomock.Any(), int64(6)).Return(&clientv3.LeaseGrantResponse{
					ID: 12,
				}, nil)
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
				mock.EXPECT().Grant(gomock.Any(), int64(6)).Return(&clientv3.LeaseGrantResponse{
					ID: 12,
				}, nil)

				mock.EXPECT().KeepAliveOnce(gomock.Any(), clientv3.LeaseID(12)).Return(nil, errors.New("NOP"))
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
			Name:           "When everything succeed",
			InitialLeaseID: 0,
			ExpectedLease: func(mock *etcdmock.MockLease) {
				mock.EXPECT().Grant(gomock.Any(), int64(6)).Return(&clientv3.LeaseGrantResponse{
					ID: 12,
				}, nil)

				mock.EXPECT().KeepAliveOnce(gomock.Any(), clientv3.LeaseID(12)).Return(nil, nil)
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

			if example.ExpectedLease != nil {
				example.ExpectedLease(leaseMock)
			}

			if example.ExpectedKV != nil {
				example.ExpectedKV(ctrl, kvMock)
			}

			locker := &etcdLocker{
				kvEtcd:    kvMock,
				leaseEtcd: leaseMock,
				key:       key,
				config: config.Config{
					KeepAliveInterval: 3 * time.Second,
				},
				leaseID: clientv3.LeaseID(example.InitialLeaseID),
			}

			err := locker.Refresh(ctx)
			if len(example.ExpectedError) != 0 {
				require.Error(t, err)
				assert.Contains(t, err.Error(), example.ExpectedError)
			} else {
				assert.NoError(t, err)
			}

			assert.Equal(t, locker.leaseID, clientv3.LeaseID(example.ExpectedLeaseID))
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
			Name:          "When there is an issue with ETCD",
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
