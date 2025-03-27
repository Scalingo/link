package web

import (
	"bytes"
	"fmt"
	"io"
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
				mock.EXPECT().Start(gomock.Any(), gomock.Any()).Return(models.Endpoint{}, scheduler.ErrIPAlreadyAssigned)
			},
			ExpectedStatusCode: http.StatusBadRequest,
			ExpectedBody:       `{"error": "IP already assigned"}`,
		}, {
			Name:  "When the scheduler fails",
			Input: `{"ip": "10.0.0.1/32"}`,
			SchedulerMock: func(mock *schedulermock.MockScheduler) {
				mock.EXPECT().Start(gomock.Any(), gomock.Any()).Return(models.Endpoint{}, errors.New("SchedFail !"))
			},
			ExpectedError: "SchedFail !",
		}, {
			Name:  "When everything works fine",
			Input: `{"ip": "10.0.0.1/32"}`,
			SchedulerMock: func(mock *schedulermock.MockScheduler) {
				mock.EXPECT().Start(gomock.Any(), gomock.Any()).Return(models.Endpoint{
					ID: "test",
					IP: "10.0.0.1/32",
				}, nil)
			},
			ExpectedBody:       `{"id":"test","ip":"10.0.0.1/32","healthcheck_interval":0}` + "\n",
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
				body, err := io.ReadAll(resp.Body)
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
				m.EXPECT().GetEndpoint(gomock.Any(), linkIPId).Return(nil)
			},
			expectedStatusCode: http.StatusNotFound,
			expectedBody:       `{"resource": "IP", "error": "not found"}`,
		},
		"With an invalid body": {
			body: "INVALID",
			expectScheduler: func(m *schedulermock.MockScheduler) {
				m.EXPECT().GetEndpoint(gomock.Any(), linkIPId).Return(&scheduler.EndpointWithStatus{})
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       `{"error": "invalid json"}`,
		},
		"With a port of 0 for the health check": {
			body: `{"healthchecks": [{"port": 0}]}`,
			expectScheduler: func(m *schedulermock.MockScheduler) {
				m.EXPECT().GetEndpoint(gomock.Any(), linkIPId).Return(&scheduler.EndpointWithStatus{})
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       `{"error": "health check port cannot be negative"}`,
		},
		"With a port of 65536 for the health check": {
			body: `{"healthchecks": [{"port": 65536}]}`,
			expectScheduler: func(m *schedulermock.MockScheduler) {
				m.EXPECT().GetEndpoint(gomock.Any(), linkIPId).Return(&scheduler.EndpointWithStatus{})
			},
			expectedStatusCode: http.StatusBadRequest,
			expectedBody:       `{"error": "health check port cannot be greater than 65535"}`,
		},
		"if it fails to update the IP": {
			body: `{"healthchecks": [{"port": 12345}]}`,
			expectScheduler: func(m *schedulermock.MockScheduler) {
				m.EXPECT().GetEndpoint(gomock.Any(), linkIPId).Return(&scheduler.EndpointWithStatus{
					Endpoint: models.Endpoint{ID: linkIPId},
					Status:   api.Activated,
				})
				m.EXPECT().UpdateEndpoint(gomock.Any(), models.Endpoint{
					ID: linkIPId,
					Checks: []models.HealthCheck{
						{Port: 12345},
					},
				}).Return(errors.New("err update IP"))
			},
			expectedError: "err update IP",
		},
		"When everything works fine": {
			body: `{"healthchecks": [{"port": 12345}]}`,
			expectScheduler: func(m *schedulermock.MockScheduler) {
				m.EXPECT().GetEndpoint(gomock.Any(), linkIPId).Return(&scheduler.EndpointWithStatus{
					Endpoint: models.Endpoint{ID: linkIPId},
					Status:   api.Activated,
				})
				m.EXPECT().UpdateEndpoint(gomock.Any(), models.Endpoint{
					ID: linkIPId,
					Checks: []models.HealthCheck{
						{Port: 12345},
					},
				})
			},
			expectedBody:       fmt.Sprintf(`{"id":"%s","ip":"","status":"ACTIVATED","checks":[{"type":"","host":"","port":12345}],"healthcheck_interval":0}`+"\n", linkIPId),
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
				body, err := io.ReadAll(res.Body)
				require.NoError(t, err)
				assert.Equal(t, test.expectedBody, string(body))
			}

			if test.expectedStatusCode != 0 {
				assert.Equal(t, test.expectedStatusCode, res.Code)
			}
		})
	}
}
