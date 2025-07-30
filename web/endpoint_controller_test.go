package web

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Scalingo/go-utils/errors/v2"
	"github.com/Scalingo/link/v3/api"
	"github.com/Scalingo/link/v3/models"
	"github.com/Scalingo/link/v3/scheduler"
	"github.com/Scalingo/link/v3/scheduler/schedulermock"
)

func TestEndpointCreator_Update(t *testing.T) {
	ctx := context.Background()
	linkIPId := "vip-11111111-1111-1111-1111-111111111111"
	tests := map[string]struct {
		body               string
		expectScheduler    func(*schedulermock.MockScheduler)
		expectedStatusCode int
		expectedBody       string
		expectedError      string
		linkIPId           string
	}{
		"With an invalid ID": {
			linkIPId:           "invalid-id",
			expectedError:      "Invalid endpoint ID",
			expectedStatusCode: http.StatusBadRequest,
		},
		"With an unknown IP": {
			expectScheduler: func(m *schedulermock.MockScheduler) {
				m.EXPECT().GetEndpoint(gomock.Any(), linkIPId).Return(nil)
			},
			expectedError:      "Endpoint not found",
			expectedStatusCode: http.StatusNotFound,
		},
		"With an invalid body": {
			body: "INVALID",
			expectScheduler: func(m *schedulermock.MockScheduler) {
				m.EXPECT().GetEndpoint(gomock.Any(), linkIPId).Return(&scheduler.EndpointWithStatus{})
			},
			expectedError: "invalid JSON",
		},
		"With a validation error": {
			body: `{"healthchecks": [{"port": 0}]}`,
			expectScheduler: func(m *schedulermock.MockScheduler) {
				m.EXPECT().GetEndpoint(gomock.Any(), linkIPId).Return(&scheduler.EndpointWithStatus{})
			},
			expectedError: "validate health checks",
		},
		"if it fails to update the IP": {
			body: `{"healthchecks": [{"type": "TCP", "host": "a.dev", "port": 12345}]}`,
			expectScheduler: func(m *schedulermock.MockScheduler) {
				m.EXPECT().GetEndpoint(gomock.Any(), linkIPId).Return(&scheduler.EndpointWithStatus{
					Endpoint: models.Endpoint{ID: linkIPId},
					Status:   api.Activated,
				})
				m.EXPECT().UpdateEndpoint(gomock.Any(), models.Endpoint{
					ID: linkIPId,
					Checks: []models.HealthCheck{
						{
							Type: "TCP",
							Host: "a.dev",
							Port: 12345,
						},
					},
				}).Return(errors.New(ctx, "err update IP"))
			},
			expectedError: "err update IP",
		},
		"When everything works fine": {
			body: `{"healthchecks": [{"type": "TCP", "host": "a.dev", "port": 12345}]}`,
			expectScheduler: func(m *schedulermock.MockScheduler) {
				m.EXPECT().GetEndpoint(gomock.Any(), linkIPId).Return(&scheduler.EndpointWithStatus{
					Endpoint: models.Endpoint{ID: linkIPId},
					Status:   api.Activated,
				})
				m.EXPECT().UpdateEndpoint(gomock.Any(), models.Endpoint{
					ID: linkIPId,
					Checks: []models.HealthCheck{
						{
							Type: "TCP",
							Host: "a.dev",
							Port: 12345,
						},
					},
				})
			},
			expectedBody:       fmt.Sprintf(`{"id":"%s","status":"ACTIVATED","checks":[{"type":"TCP","host":"a.dev","port":12345}],"healthcheck_interval":0,"plugin":"arp"}`+"\n", linkIPId),
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

			if test.linkIPId == "" {
				test.linkIPId = linkIPId
			}

			err := EndpointController{
				scheduler: scheduler,
			}.Update(res, req, map[string]string{"id": test.linkIPId})
			if test.expectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), test.expectedError)
			} else {
				require.NoError(t, err)
			}

			if test.expectedBody != "" {
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
