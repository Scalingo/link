package web

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Scalingo/link/models"
	"github.com/Scalingo/link/models/modelsmock"
	"github.com/Scalingo/link/network/networkmock"
	"github.com/Scalingo/link/scheduler"
	"github.com/Scalingo/link/scheduler/schedulermock"
	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestIPController_Create(t *testing.T) {
	examples := []struct {
		Name               string
		Input              string
		InterfaceMock      func(mock *networkmock.MockNetworkInterface)
		StorageMock        func(mock *modelsmock.MockStorage)
		SchedulerMock      func(mock *schedulermock.MockScheduler)
		ExpectedStatusCode int
		ExpectedBody       string
		ExpectedError      string
	}{
		{
			Name:               "With an invalid body",
			Input:              "INVALID",
			ExpectedStatusCode: http.StatusBadRequest,
			ExpectedBody:       `{"msg": "invalid json"}`,
		}, {
			Name:               "With an invalid CIDR",
			Input:              `{"ip": "INVALID!!!"}`,
			ExpectedStatusCode: http.StatusBadRequest,
			ExpectedBody:       `{"msg": "invalid IP"}`,
		}, {
			Name:  "With an IP that already has been assigned",
			Input: `{"ip": "10.0.0.1/32"}`,
			InterfaceMock: func(mock *networkmock.MockNetworkInterface) {
				mock.EXPECT().HasIP("10.0.0.1/32").Return(true, nil)
			},
			StorageMock: func(mock *modelsmock.MockStorage) {
				mock.EXPECT().AddIP(gomock.Any(), gomock.Any()).Return(models.IP{}, models.ErrIPAlreadyPresent)
			},
			SchedulerMock: func(mock *schedulermock.MockScheduler) {
				mock.EXPECT().CancelStopping(gomock.Any(), gomock.Any()).Return(scheduler.ErrNotStopping)
			},
			ExpectedStatusCode: http.StatusBadRequest,
			ExpectedBody:       `{"msg": "IP already assigned"}`,
		}, {
			Name:  "When the storage fails",
			Input: `{"ip": "10.0.0.1/32"}`,
			InterfaceMock: func(mock *networkmock.MockNetworkInterface) {
				mock.EXPECT().HasIP("10.0.0.1/32").Return(false, nil)
			},
			StorageMock: func(mock *modelsmock.MockStorage) {
				mock.EXPECT().AddIP(gomock.Any(), gomock.Any()).Return(models.IP{}, errors.New("FAIL !"))
			},
			ExpectedError: "FAIL !",
		}, {
			Name:  "When the scheduler fails",
			Input: `{"ip": "10.0.0.1/32"}`,
			InterfaceMock: func(mock *networkmock.MockNetworkInterface) {
				mock.EXPECT().HasIP("10.0.0.1/32").Return(false, nil)
			},
			StorageMock: func(mock *modelsmock.MockStorage) {
				mock.EXPECT().AddIP(gomock.Any(), gomock.Any()).Return(models.IP{IP: "10.0.0.1/32"}, nil)
			},
			SchedulerMock: func(mock *schedulermock.MockScheduler) {
				mock.EXPECT().Start(gomock.Any(), gomock.Any()).Return(errors.New("SchedFail !"))
			},
			ExpectedError: "SchedFail !",
		}, {
			Name:  "When everything works fine",
			Input: `{"ip": "10.0.0.1/32"}`,
			InterfaceMock: func(mock *networkmock.MockNetworkInterface) {
				mock.EXPECT().HasIP("10.0.0.1/32").Return(false, nil)
			},
			StorageMock: func(mock *modelsmock.MockStorage) {
				mock.EXPECT().AddIP(gomock.Any(), gomock.Any()).Return(models.IP{IP: "10.0.0.1/32", ID: "test"}, nil)
			},
			SchedulerMock: func(mock *schedulermock.MockScheduler) {
				mock.EXPECT().Start(gomock.Any(), gomock.Any()).Return(nil)
			},
			ExpectedBody:       `{"id":"test","ip":"10.0.0.1/32"}` + "\n",
			ExpectedStatusCode: http.StatusCreated,
		},
	}

	for _, example := range examples {
		t.Run(example.Name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			storageMock := modelsmock.NewMockStorage(ctrl)
			schedulerMock := schedulermock.NewMockScheduler(ctrl)
			netMock := networkmock.NewMockNetworkInterface(ctrl)

			ipCtrl := ipController{
				storage:      storageMock,
				scheduler:    schedulerMock,
				netInterface: netMock,
			}

			if example.InterfaceMock != nil {
				example.InterfaceMock(netMock)
			}

			if example.StorageMock != nil {
				example.StorageMock(storageMock)
			}

			if example.SchedulerMock != nil {
				example.SchedulerMock(schedulerMock)
			}

			req := httptest.NewRequest("POST", "/ips", bytes.NewBufferString(example.Input))
			resp := httptest.NewRecorder()

			err := ipCtrl.Create(resp, req, nil)
			if len(example.ExpectedError) > 0 {
				require.Error(t, err)
				assert.Contains(t, err.Error(), example.ExpectedError)
			} else {
				require.NoError(t, err)
			}

			if len(example.ExpectedBody) > 0 {
				body, err := ioutil.ReadAll(resp.Body)
				require.NoError(t, err)
				assert.Equal(t, example.ExpectedBody, string(body))
			}

			if example.ExpectedStatusCode != 0 {
				assert.Equal(t, example.ExpectedStatusCode, resp.Code)
			}
		})
	}
}
