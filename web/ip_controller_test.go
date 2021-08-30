package web

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Scalingo/link/v2/api"
	"github.com/Scalingo/link/v2/models"
	"github.com/Scalingo/link/v2/scheduler"
	"github.com/Scalingo/link/v2/scheduler/schedulermock"
)

func TestIPController_Create(t *testing.T) {
	examples := []struct {
		Name               string
		Input              string
		SchedulerMock      func(*schedulermock.MockScheduler)
		ExpectedStatusCode int
		ExpectedBody       string
		ExpectedError      string
	}{
		{
			Name:               "With an invalid body",
			Input:              "INVALID",
			ExpectedStatusCode: http.StatusBadRequest,
			ExpectedBody:       `{"error": "invalid json"}`,
		}, {
			Name:               "With an invalid CIDR",
			Input:              `{"ip": "INVALID!!!"}`,
			ExpectedStatusCode: http.StatusBadRequest,
			ExpectedBody:       `{"error": "invalid IP"}`,
		}, {
			Name:               "With a port of 0 for the health check",
			Input:              `{"ip": "10.0.0.1/32", "checks": [{"port": 0}]}`,
			ExpectedStatusCode: http.StatusBadRequest,
			ExpectedBody:       `{"error": "health check port cannot be negative"}`,
		}, {
			Name:               "With a port of 65536 for the health check",
			Input:              `{"ip": "10.0.0.1/32", "checks": [{"port": 65536}]}`,
			ExpectedStatusCode: http.StatusBadRequest,
			ExpectedBody:       `{"error": "health check port cannot be greater than 65535"}`,
		}, {
			Name:  "With an IP that already has been assigned",
			Input: `{"ip": "10.0.0.1/32"}`,
			SchedulerMock: func(mock *schedulermock.MockScheduler) {
				mock.EXPECT().Start(gomock.Any(), gomock.Any()).Return(models.IP{}, scheduler.ErrIPAlreadyAssigned)
			},
			ExpectedStatusCode: http.StatusBadRequest,
			ExpectedBody:       `{"error": "IP already assigned"}`,
		}, {
			Name:  "When the scheduler fails",
			Input: `{"ip": "10.0.0.1/32"}`,
			SchedulerMock: func(mock *schedulermock.MockScheduler) {
				mock.EXPECT().Start(gomock.Any(), gomock.Any()).Return(models.IP{}, errors.New("SchedFail !"))
			},
			ExpectedError: "SchedFail !",
		}, {
			Name:  "When everything works fine",
			Input: `{"ip": "10.0.0.1/32", "no_network": true}`,
			SchedulerMock: func(mock *schedulermock.MockScheduler) {
				mock.EXPECT().Start(gomock.Any(), gomock.Any()).Return(models.IP{
					ID:        "test",
					IP:        "10.0.0.1/32",
					NoNetwork: true,
				}, nil)
			},
			ExpectedBody:       `{"id":"test","ip":"10.0.0.1/32","healthcheck_interval":0,"no_network":true}` + "\n",
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

func TestIPController_Patch(t *testing.T) {
	linkIPId := "my-id"
	tests := map[string]struct {
		body               string
		expectScheduler    func(*schedulermock.MockScheduler)
		expectedStatusCode int
		expectedBody       string
		expectedError      string
	}{
		"With an unknown IP": {
			expectScheduler: func(m *schedulermock.MockScheduler) {
				m.EXPECT().GetIP(gomock.Any(), linkIPId).Return(nil)
			},
			expectedStatusCode: http.StatusNotFound,
			expectedBody:       `{"resource": "IP", "error": "not found"}`,
		},
		"With an invalid body": {
			body: "INVALID",
			expectScheduler: func(m *schedulermock.MockScheduler) {
				m.EXPECT().GetIP(gomock.Any(), linkIPId).Return(&api.IP{})
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       `{"error": "invalid json to patch the IP"}`,
		},
		"With a port of 0 for the health check": {
			body: `{"healthchecks": [{"port": 0}]}`,
			expectScheduler: func(m *schedulermock.MockScheduler) {
				m.EXPECT().GetIP(gomock.Any(), linkIPId).Return(&api.IP{})
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       `{"error": "health check port cannot be negative"}`,
		},
		"With a port of 65536 for the health check": {
			body: `{"healthchecks": [{"port": 65536}]}`,
			expectScheduler: func(m *schedulermock.MockScheduler) {
				m.EXPECT().GetIP(gomock.Any(), linkIPId).Return(&api.IP{})
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       `{"error": "health check port cannot be greater than 65535"}`,
		},
		"if it fails to update the IP": {
			body: `{"healthchecks": [{"port": 12345}]}`,
			expectScheduler: func(m *schedulermock.MockScheduler) {
				m.EXPECT().GetIP(gomock.Any(), linkIPId).Return(&api.IP{
					IP:     models.IP{ID: linkIPId},
					Status: api.Activated,
				})
				m.EXPECT().UpdateIP(gomock.Any(), models.IP{
					ID: linkIPId,
					Checks: []models.Healthcheck{
						{Port: 12345},
					},
				}).Return(errors.New("err update IP"))
			},
			expectedError: "err update IP",
		},
		"When everything works fine": {
			body: `{"healthchecks": [{"port": 12345}]}`,
			expectScheduler: func(m *schedulermock.MockScheduler) {
				m.EXPECT().GetIP(gomock.Any(), linkIPId).Return(&api.IP{
					IP:     models.IP{ID: linkIPId},
					Status: api.Activated,
				})
				m.EXPECT().UpdateIP(gomock.Any(), models.IP{
					ID: linkIPId,
					Checks: []models.Healthcheck{
						{Port: 12345},
					},
				})
			},
			expectedBody:       fmt.Sprintf(`{"id":"%s","ip":"","checks":[{"type":"","host":"","port":12345}],"healthcheck_interval":0,"no_network":false}`+"\n", linkIPId),
			expectedStatusCode: http.StatusOK,
		},
		"When a user tries to update the NoNetwork flag": {
			body: `{"no_network": true}`,
			expectScheduler: func(m *schedulermock.MockScheduler) {
				m.EXPECT().GetIP(gomock.Any(), linkIPId).Return(&api.IP{
					IP: models.IP{ID: linkIPId, NoNetwork: false},
				})
				m.EXPECT().UpdateIP(gomock.Any(), models.IP{
					ID:        linkIPId,
					NoNetwork: true,
				})
			},
			expectedBody:       fmt.Sprintf(`{"id":"%s","ip":"","healthcheck_interval":0,"no_network":true}`+"\n", linkIPId),
			expectedStatusCode: http.StatusOK,
		},
	}

	for msg, test := range tests {
		t.Run(msg, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			scheduler := schedulermock.NewMockScheduler(ctrl)

			if test.expectScheduler != nil {
				test.expectScheduler(scheduler)
			}

			req := httptest.NewRequest("POST", "/ips", bytes.NewBufferString(test.body))
			res := httptest.NewRecorder()

			err := ipController{
				scheduler: scheduler,
			}.Patch(res, req, map[string]string{"id": linkIPId})
			if len(test.expectedError) > 0 {
				require.Error(t, err)
				assert.Contains(t, err.Error(), test.expectedError)
			} else {
				require.NoError(t, err)
			}

			if len(test.expectedBody) > 0 {
				body, err := ioutil.ReadAll(res.Body)
				require.NoError(t, err)
				assert.Equal(t, test.expectedBody, string(body))
			}

			if test.expectedStatusCode != 0 {
				assert.Equal(t, test.expectedStatusCode, res.Code)
			}
		})
	}
}
