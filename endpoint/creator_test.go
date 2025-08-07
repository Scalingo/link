package endpoint

import (
	"context"
	"errors"
	"testing"

	"github.com/!scalingo/link/v2/models/modelsmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/golang/mock/gomock"

	"github.com/Scalingo/link/v3/api"
	"github.com/Scalingo/link/v3/models"
	"github.com/Scalingo/link/v3/plugin/pluginmock"
	"github.com/Scalingo/link/v3/scheduler/schedulermock"
)

func Test_Creator_CreateEndpoint(t *testing.T) {
	specs := []struct {
		Name                 string
		Params               CreateEndpointParams
		Storage              func(t *testing.T, mock *models.MockStorage)
		Registry             func(mock *pluginmock.MockRegistry)
		Scheduler            func(mock *schedulermock.MockScheduler)
		ExpectedError        string
		maxNumberOfEndpoints int
	}{
		{
			Name: "Invalid health check",
			Params: CreateEndpointParams{
				Checks: []api.HealthCheck{
					{
						Type: "INVALID",
					},
				},
			},
			Scheduler: func(mock *schedulermock.MockScheduler) {
				mock.EXPECT().EndpointCount().Return(0)
			},
			ExpectedError: "Health check type is not supported",
		}, {
			Name: "Invalid health check interval",
			Scheduler: func(mock *schedulermock.MockScheduler) {
				mock.EXPECT().EndpointCount().Return(0)
			},
			Params: CreateEndpointParams{
				HealthCheckInterval: -1,
			},
			ExpectedError: "Health check interval must be greater than 0",
		}, {
			Name: "Health check interval too high",
			Scheduler: func(mock *schedulermock.MockScheduler) {
				mock.EXPECT().EndpointCount().Return(0)
			},
			Params: CreateEndpointParams{
				HealthCheckInterval: 4000,
			},
			ExpectedError: "Health check interval must be less than or equal to 3600 seconds",
		}, {
			Name: "Plugin validation failed",
			Registry: func(mock *pluginmock.MockRegistry) {
				mock.EXPECT().Validate(gomock.Any(), gomock.Any()).Return(errors.New("plugin validation error"))
			},
			Scheduler: func(mock *schedulermock.MockScheduler) {
				mock.EXPECT().EndpointCount().Return(0)
			},
			ExpectedError: "plugin validation error",
		}, {
			Name: "Storage error",
			Registry: func(mock *pluginmock.MockRegistry) {
				mock.EXPECT().Validate(gomock.Any(), gomock.Any()).Return(nil)
				mock.EXPECT().Mutate(gomock.Any(), gomock.Any()).Return(nil, nil)
			},
			Storage: func(_ *testing.T, mock *models.MockStorage) {
				mock.EXPECT().AddEndpoint(gomock.Any(), gomock.Any()).Return(models.Endpoint{}, errors.New("storage error"))
			},
			Scheduler: func(mock *schedulermock.MockScheduler) {
				mock.EXPECT().EndpointCount().Return(0)
			},
			ExpectedError: "storage error",
		}, {
			Name: "Scheduler Start error",
			Registry: func(mock *pluginmock.MockRegistry) {
				mock.EXPECT().Validate(gomock.Any(), gomock.Any()).Return(nil)
				mock.EXPECT().Mutate(gomock.Any(), gomock.Any()).Return(nil, nil)
			},
			Storage: func(_ *testing.T, mock *models.MockStorage) {
				mock.EXPECT().AddEndpoint(gomock.Any(), gomock.Any()).Return(models.Endpoint{ID: "test-id"}, nil)
				mock.EXPECT().RemoveEndpoint(gomock.Any(), "test-id").Return(nil)
			},
			Scheduler: func(mock *schedulermock.MockScheduler) {
				mock.EXPECT().Start(gomock.Any(), gomock.Any()).Return(models.Endpoint{}, errors.New("scheduler error"))
				mock.EXPECT().EndpointCount().Return(0)
			},
			ExpectedError: "scheduler error",
		}, {
			Name: "Successful creation",
			Registry: func(mock *pluginmock.MockRegistry) {
				mock.EXPECT().Validate(gomock.Any(), gomock.Any()).Return(nil)
				mock.EXPECT().Mutate(gomock.Any(), gomock.Any()).Return(nil, nil)
			},
			Storage: func(_ *testing.T, mock *models.MockStorage) {
				mock.EXPECT().AddEndpoint(gomock.Any(), gomock.Any()).Return(models.Endpoint{ID: "test-id"}, nil)
			},
			Scheduler: func(mock *schedulermock.MockScheduler) {
				mock.EXPECT().Start(gomock.Any(), gomock.Any()).Return(models.Endpoint{ID: "test-id"}, nil)
				mock.EXPECT().EndpointCount().Return(0)
			},
			ExpectedError: "",
		}, {
			Name: "The plugin needs to perform a mutation",
			Registry: func(mock *pluginmock.MockRegistry) {
				mock.EXPECT().Validate(gomock.Any(), gomock.Any()).Return(nil)
				mock.EXPECT().Mutate(gomock.Any(), gomock.Any()).Return([]byte(`{"key":"value"}`), nil)
			},
			Storage: func(t *testing.T, mock *models.MockStorage) {
				mock.EXPECT().AddEndpoint(gomock.Any(), gomock.Any()).Return(models.Endpoint{ID: "test-id"}, nil).Do(func(_ context.Context, endpoint models.Endpoint) {
					assert.JSONEq(t, `{"key":"value"}`, string(endpoint.PluginConfig))
				})
			},
			Scheduler: func(mock *schedulermock.MockScheduler) {
				mock.EXPECT().Start(gomock.Any(), gomock.Any()).Return(models.Endpoint{ID: "test-id"}, nil)
				mock.EXPECT().EndpointCount().Return(0)
			},
			ExpectedError: "",
		}, {
			Name: "Too many endpoints configured",
			Scheduler: func(mock *schedulermock.MockScheduler) {
				mock.EXPECT().EndpointCount().Return(1001)
			},
			ExpectedError: "Too many endpoints configured: 1001, max allowed: 1000",
		}, {
			Name: "Too many endpoints check disabled",
			Registry: func(mock *pluginmock.MockRegistry) {
				mock.EXPECT().Validate(gomock.Any(), gomock.Any()).Return(nil)
				mock.EXPECT().Mutate(gomock.Any(), gomock.Any()).Return(nil, nil)
			},
			Storage: func(_ *testing.T, mock *modelsmock.MockStorage) {
				mock.EXPECT().AddEndpoint(gomock.Any(), gomock.Any()).Return(models.Endpoint{ID: "test-id"}, nil)
			},
			Scheduler: func(mock *schedulermock.MockScheduler) {
				mock.EXPECT().Start(gomock.Any(), gomock.Any()).Return(models.Endpoint{ID: "test-id"}, nil)
				mock.EXPECT().EndpointCount().Return(1001)
			},
			maxNumberOfEndpoints: -1,
			ExpectedError:        "",
		},
	}

	for _, spec := range specs {
		t.Run(spec.Name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			ctx := context.Background()

			// Mock dependencies
			mockStorage := models.NewMockStorage(ctrl)
			mockRegistry := pluginmock.NewMockRegistry(ctrl)
			mockScheduler := schedulermock.NewMockScheduler(ctrl)

			// Apply the mock setups if provided
			if spec.Storage != nil {
				spec.Storage(t, mockStorage)
			}
			if spec.Registry != nil {
				spec.Registry(mockRegistry)
			}
			if spec.Scheduler != nil {
				spec.Scheduler(mockScheduler)
			}

			if spec.maxNumberOfEndpoints == 0 {
				spec.maxNumberOfEndpoints = 1000
			}

			// Call the function under test
			creator := creator{
				storage:              mockStorage,
				registry:             mockRegistry,
				scheduler:            mockScheduler,
				maxNumberOfEndpoints: spec.maxNumberOfEndpoints,
			}
			_, err := creator.CreateEndpoint(ctx, spec.Params)

			// Assert the expected error
			if spec.ExpectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), spec.ExpectedError)
			} else {
				require.NoError(t, err)
			}
		})
	}
}
