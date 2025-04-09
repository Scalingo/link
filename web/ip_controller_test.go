package web

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Scalingo/link/v2/endpoint"
	"github.com/Scalingo/link/v2/endpoint/endpointmock"
	"github.com/Scalingo/link/v2/models"
)

func TestIPController_Create(t *testing.T) {
	examples := []struct {
		Name               string
		Input              string
		EndpointCreator    func(mock *endpointmock.MockCreator)
		ExpectedStatusCode int
		ExpectedBody       string
		ExpectedError      string
	}{
		{
			Name:  "When the creator fails",
			Input: `{"ip": "10.0.0.1/32"}`,
			EndpointCreator: func(mock *endpointmock.MockCreator) {
				mock.EXPECT().CreateEndpoint(gomock.Any(), endpoint.CreateEndpointParams{
					Plugin:       "arp",
					PluginConfig: []byte(`{"ip":"10.0.0.1/32"}`),
				}).Return(models.Endpoint{}, errors.New("creator fails"))
			},
			ExpectedError: "creator fails",
		}, {
			Name:  "When everything works fine",
			Input: `{"ip": "10.0.0.1/32"}`,
			EndpointCreator: func(mock *endpointmock.MockCreator) {
				mock.EXPECT().CreateEndpoint(gomock.Any(), endpoint.CreateEndpointParams{
					Plugin:       "arp",
					PluginConfig: []byte(`{"ip":"10.0.0.1/32"}`),
				}).Return(models.Endpoint{
					ID:     "test",
					Plugin: "arp",
				}, nil)
			},
			ExpectedBody:       `{"id":"test","healthcheck_interval":0,"plugin":"arp"}` + "\n",
			ExpectedStatusCode: http.StatusCreated,
		},
	}

	for _, example := range examples {
		t.Run(example.Name, func(t *testing.T) {
			ctrl := gomock.NewController(t)

			endpointCreator := endpointmock.NewMockCreator(ctrl)
			if example.EndpointCreator != nil {
				example.EndpointCreator(endpointCreator)
			}

			ipCtrl := IPController{
				endpointCreator: endpointCreator,
			}

			req := httptest.NewRequest("POST", "/ips", bytes.NewBufferString(example.Input))
			resp := httptest.NewRecorder()

			err := ipCtrl.Create(resp, req, nil)
			if example.ExpectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), example.ExpectedError)
			} else {
				require.NoError(t, err)
			}

			if example.ExpectedBody != "" {
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
