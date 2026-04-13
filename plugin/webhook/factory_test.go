package webhook

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Scalingo/link/v3/models"
)

func TestFactoryValidate(t *testing.T) {
	tests := []struct {
		name          string
		cfg           any
		expectedError string
	}{
		{
			name: "valid config",
			cfg: PluginConfig{
				URL: "https://example.com/webhook",
				Headers: map[string]string{
					"Authorization": "Bearer token",
				},
			},
		},
		{
			name:          "invalid json",
			cfg:           `{"url":`,
			expectedError: "unmarshal plugin config",
		},
		{
			name: "missing URL",
			cfg: PluginConfig{
				Headers: map[string]string{"X-Test": "x"},
			},
			expectedError: "URL is required",
		},
		{
			name: "invalid URL",
			cfg: PluginConfig{
				URL: "::::",
			},
			expectedError: "URL should be valid",
		},
		{
			name: "invalid scheme",
			cfg: PluginConfig{
				URL: "ftp://example.com",
			},
			expectedError: "URL scheme must be http or https",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var raw []byte
			switch v := tt.cfg.(type) {
			case string:
				raw = []byte(v)
			default:
				var err error
				raw, err = json.Marshal(tt.cfg)
				require.NoError(t, err)
			}

			endpoint := models.Endpoint{PluginConfig: raw}
			err := Factory{}.Validate(context.Background(), endpoint)

			if tt.expectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.expectedError)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestFactoryCreate(t *testing.T) {
	raw, err := json.Marshal(StorablePluginConfig{URL: "https://example.com/hook"})
	require.NoError(t, err)

	endpoint := models.Endpoint{
		ID:           "vip-test-id",
		Plugin:       Name,
		PluginConfig: raw,
	}

	p, err := Factory{}.Create(context.Background(), endpoint)
	require.NoError(t, err)
	require.NotNil(t, p)

	assert.Equal(t, "vip-test-id", p.ElectionKey(context.Background()))
}

func TestFactoryCreateWithEncryptedHeaders(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	mockStorage := models.NewMockEncryptedStorage(ctrl)

	raw, err := json.Marshal(StorablePluginConfig{
		URL: "https://example.com/hook",
		Headers: map[string]models.EncryptedDataLink{
			"Authorization": {ID: "auth-id", EndpointID: "vip-test-id"},
			"X-App":         {ID: "x-app-id", EndpointID: "vip-test-id"},
		},
	})
	require.NoError(t, err)

	endpoint := models.Endpoint{
		ID:           "vip-test-id",
		Plugin:       Name,
		PluginConfig: raw,
	}

	mockStorage.EXPECT().
		Decrypt(ctx, models.EncryptedDataLink{ID: "auth-id", EndpointID: "vip-test-id"}, gomock.Any()).
		DoAndReturn(func(_ context.Context, _ models.EncryptedDataLink, out any) error {
			*(out.(*string)) = "Bearer token"
			return nil
		})
	mockStorage.EXPECT().
		Decrypt(ctx, models.EncryptedDataLink{ID: "x-app-id", EndpointID: "vip-test-id"}, gomock.Any()).
		DoAndReturn(func(_ context.Context, _ models.EncryptedDataLink, out any) error {
			*(out.(*string)) = "link"
			return nil
		})

	p, err := Factory{encryptedStorage: mockStorage}.Create(ctx, endpoint)
	require.NoError(t, err)
	require.NotNil(t, p)
}

func TestFactoryMutate(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)
	mockStorage := models.NewMockEncryptedStorage(ctrl)

	raw, err := json.Marshal(PluginConfig{
		URL: "https://example.com/hook",
		Headers: map[string]string{
			"Authorization": "Bearer token",
			"X-App":         "link",
		},
	})
	require.NoError(t, err)

	endpoint := models.Endpoint{
		ID:           "vip-test-id",
		Plugin:       Name,
		PluginConfig: raw,
	}

	mockStorage.EXPECT().
		Encrypt(ctx, "vip-test-id", "Bearer token").
		Return(models.EncryptedDataLink{ID: "auth-id", EndpointID: "vip-test-id"}, nil)
	mockStorage.EXPECT().
		Encrypt(ctx, "vip-test-id", "link").
		Return(models.EncryptedDataLink{ID: "x-app-id", EndpointID: "vip-test-id"}, nil)

	mutated, err := Factory{encryptedStorage: mockStorage}.Mutate(ctx, endpoint)
	require.NoError(t, err)

	var stored StorablePluginConfig
	err = json.Unmarshal(mutated, &stored)
	require.NoError(t, err)
	assert.Equal(t, "https://example.com/hook", stored.URL)
	assert.Equal(t, "auth-id", stored.Headers["Authorization"].ID)
	assert.Equal(t, "x-app-id", stored.Headers["X-App"].ID)
	assert.Equal(t, "vip-test-id", stored.Headers["Authorization"].EndpointID)
	assert.Equal(t, "vip-test-id", stored.Headers["X-App"].EndpointID)
}
