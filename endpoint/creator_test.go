package endpoint

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/golang/mock/gomock"

	"github.com/Scalingo/link/v2/api"
	"github.com/Scalingo/link/v2/models"
	"github.com/Scalingo/link/v2/models/modelsmock"
	"github.com/Scalingo/link/v2/plugin/pluginmock"
	"github.com/Scalingo/link/v2/scheduler/schedulermock"
)

func Test_Creator_CreateEndpoint(t *testing.T) {
	specs := []struct {
		Name          string
		Params        CreateEndpointParams
		Storage       func(t *testing.T, mock *modelsmock.MockStorage)
		Registry      func(mock *pluginmock.MockRegistry)
		Scheduler     func(mock *schedulermock.MockScheduler)
		ExpectedError string
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
			ExpectedError: "Health check type is not supported",
		}, {
			Name: "Plugin validation failed",
			Registry: func(mock *pluginmock.MockRegistry) {
				mock.EXPECT().Validate(gomock.Any(), gomock.Any()).Return(errors.New("plugin validation error"))
			},
			ExpectedError: "plugin validation error",
		}, {
			Name: "Storage error",
			Registry: func(mock *pluginmock.MockRegistry) {
				mock.EXPECT().Validate(gomock.Any(), gomock.Any()).Return(nil)
				mock.EXPECT().Mutate(gomock.Any(), gomock.Any()).Return(nil, nil)
			},
			Storage: func(_ *testing.T, mock *modelsmock.MockStorage) {
				mock.EXPECT().AddEndpoint(gomock.Any(), gomock.Any()).Return(models.Endpoint{}, errors.New("storage error"))
			},
			ExpectedError: "storage error",
		}, {
			Name: "Scheduler Start error",
			Registry: func(mock *pluginmock.MockRegistry) {
				mock.EXPECT().Validate(gomock.Any(), gomock.Any()).Return(nil)
				mock.EXPECT().Mutate(gomock.Any(), gomock.Any()).Return(nil, nil)
			},
			Storage: func(_ *testing.T, mock *modelsmock.MockStorage) {
				mock.EXPECT().AddEndpoint(gomock.Any(), gomock.Any()).Return(models.Endpoint{ID: "test-id"}, nil)
				mock.EXPECT().RemoveEndpoint(gomock.Any(), "test-id").Return(nil)
			},
			Scheduler: func(mock *schedulermock.MockScheduler) {
				mock.EXPECT().Start(gomock.Any(), gomock.Any()).Return(models.Endpoint{}, errors.New("scheduler error"))
			},
			ExpectedError: "scheduler error",
		}, {
			Name: "Successful creation",
			Registry: func(mock *pluginmock.MockRegistry) {
				mock.EXPECT().Validate(gomock.Any(), gomock.Any()).Return(nil)
				mock.EXPECT().Mutate(gomock.Any(), gomock.Any()).Return(nil, nil)
			},
			Storage: func(_ *testing.T, mock *modelsmock.MockStorage) {
				mock.EXPECT().AddEndpoint(gomock.Any(), gomock.Any()).Return(models.Endpoint{ID: "test-id"}, nil)
			},
			Scheduler: func(mock *schedulermock.MockScheduler) {
				mock.EXPECT().Start(gomock.Any(), gomock.Any()).Return(models.Endpoint{ID: "test-id"}, nil)
			},
			ExpectedError: "",
		}, {
			Name: "The plugin needs to perform a mutation",
			Registry: func(mock *pluginmock.MockRegistry) {
				mock.EXPECT().Validate(gomock.Any(), gomock.Any()).Return(nil)
				mock.EXPECT().Mutate(gomock.Any(), gomock.Any()).Return([]byte(`{"key":"value"}`), nil)
			},
			Storage: func(t *testing.T, mock *modelsmock.MockStorage) {
				mock.EXPECT().AddEndpoint(gomock.Any(), gomock.Any()).Return(models.Endpoint{ID: "test-id"}, nil).Do(func(_ context.Context, endpoint models.Endpoint) {
					assert.JSONEq(t, `{"key":"value"}`, string(endpoint.PluginConfig))
				})
			},
			Scheduler: func(mock *schedulermock.MockScheduler) {
				mock.EXPECT().Start(gomock.Any(), gomock.Any()).Return(models.Endpoint{ID: "test-id"}, nil)
			},
			ExpectedError: "",
		},
	}

	for _, spec := range specs {
		t.Run(spec.Name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			ctx := context.Background()

			// Mock dependencies
			mockStorage := modelsmock.NewMockStorage(ctrl)
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

			// Call the function under test
			creator := creator{
				storage:   mockStorage,
				registry:  mockRegistry,
				scheduler: mockScheduler,
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
