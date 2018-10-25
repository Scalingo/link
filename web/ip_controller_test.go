package web

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Scalingo/link/models"
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
			SchedulerMock: func(mock *schedulermock.MockScheduler) {
				mock.EXPECT().Start(gomock.Any(), gomock.Any()).Return(models.IP{}, scheduler.ErrNotStopping)
			},
			ExpectedStatusCode: http.StatusBadRequest,
			ExpectedBody:       `{"msg": "IP already assigned"}`,
		}, {
			Name:  "When the scheduler fails",
			Input: `{"ip": "10.0.0.1/32"}`,
			SchedulerMock: func(mock *schedulermock.MockScheduler) {
				mock.EXPECT().Start(gomock.Any(), gomock.Any()).Return(models.IP{}, errors.New("SchedFail !"))
			},
			ExpectedError: "SchedFail !",
		}, {
			Name:  "When everything works fine",
			Input: `{"ip": "10.0.0.1/32"}`,
			SchedulerMock: func(mock *schedulermock.MockScheduler) {
				mock.EXPECT().Start(gomock.Any(), gomock.Any()).Return(models.IP{
					ID: "test",
					IP: "10.0.0.1/32",
				}, nil)
			},
			ExpectedBody:       `{"id":"test","ip":"10.0.0.1/32"}` + "\n",
			ExpectedStatusCode: http.StatusCreated,
		},
	}

	for _, example := range examples {
		t.Run(example.Name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			schedulerMock := schedulermock.NewMockScheduler(ctrl)

			ipCtrl := ipController{
				scheduler: schedulerMock,
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
