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
				URL:        "https://example.com/webhook",
				ResourceID: "resource-1",
				Headers: map[string]string{
					"Authorization": "Bearer token",
				},
			},
		},
		{
			name:          "invalid json",
			cfg:           `{"url":`,
			expectedError: "invalid JSON",
		},
		{
			name: "missing URL",
			cfg: PluginConfig{
				ResourceID: "resource-1",
				Headers:    map[string]string{"X-Test": "x"},
			},
			expectedError: "missing URL",
		},
		{
			name: "missing resource id",
			cfg: PluginConfig{
				URL: "https://example.com/webhook",
			},
			expectedError: "missing resource ID",
		},
		{
			name: "invalid URL",
			cfg: PluginConfig{
				URL:        "::::",
				ResourceID: "resource-1",
			},
			expectedError: "invalid URL",
		},
		{
			name: "invalid scheme",
			cfg: PluginConfig{
				URL:        "ftp://example.com",
				ResourceID: "resource-1",
			},
			expectedError: "invalid URL scheme",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var raw []byte
			switch v := tt.cfg.(type) {
			case string:
				raw = []byte(v)
			default:
				raw, _ = json.Marshal(tt.cfg)
			}

			endpoint := models.Endpoint{PluginConfig: raw}
			err := Factory{}.Validate(t.Context(), endpoint)

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
	raw, _ := json.Marshal(StorablePluginConfig{URL: "https://example.com/hook", ResourceID: "resource-123"})

	endpoint := models.Endpoint{
		ID:           "vip-test-id",
		Plugin:       Name,
		PluginConfig: raw,
	}

	p, err := Factory{}.Create(t.Context(), endpoint)
	require.NoError(t, err)
	require.NotNil(t, p)

	assert.Equal(t, "webhook/resource-123", p.ElectionKey(t.Context()))
}

func TestFactoryCreateWithEncryptedHeaders(t *testing.T) {
	ctx := t.Context()
	ctrl := gomock.NewController(t)
	mockStorage := models.NewMockEncryptedStorage(ctrl)

	raw, _ := json.Marshal(StorablePluginConfig{
		URL:        "https://example.com/hook",
		ResourceID: "resource-123",
		Headers: map[string]models.EncryptedDataLink{
			"Authorization": {ID: "auth-id", EndpointID: "vip-test-id"},
			"X-App":         {ID: "x-app-id", EndpointID: "vip-test-id"},
		},
	})

	endpoint := models.Endpoint{
		ID:           "vip-test-id",
		Plugin:       Name,
		PluginConfig: raw,
	}

	mockStorage.EXPECT().
		Decrypt(ctx, models.EncryptedDataLink{ID: "auth-id", EndpointID: "vip-test-id"}, gomock.Any()).
		DoAndReturn(func(_ context.Context, _ models.EncryptedDataLink, out any) error {
			outStr, ok := out.(*string)
			if !ok {
				return assert.AnError
			}
			*outStr = "Bearer token"
			return nil
		})
	mockStorage.EXPECT().
		Decrypt(ctx, models.EncryptedDataLink{ID: "x-app-id", EndpointID: "vip-test-id"}, gomock.Any()).
		DoAndReturn(func(_ context.Context, _ models.EncryptedDataLink, out any) error {
			outStr, ok := out.(*string)
			if !ok {
				return assert.AnError
			}
			*outStr = "link"
			return nil
		})

	p, err := Factory{encryptedStorage: mockStorage}.Create(ctx, endpoint)
	require.NoError(t, err)
	require.NotNil(t, p)
	assert.Equal(t, "webhook/resource-123", p.ElectionKey(ctx))
}

func TestFactoryMutate(t *testing.T) {
	ctx := t.Context()
	ctrl := gomock.NewController(t)
	mockStorage := models.NewMockEncryptedStorage(ctrl)

	raw, _ := json.Marshal(PluginConfig{
		URL:        "https://example.com/hook",
		ResourceID: "resource-123",
		Headers: map[string]string{
			"Authorization": "Bearer token",
			"X-App":         "link",
		},
	})

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
	assert.Equal(t, "resource-123", stored.ResourceID)
	assert.Equal(t, "auth-id", stored.Headers["Authorization"].ID)
	assert.Equal(t, "x-app-id", stored.Headers["X-App"].ID)
	assert.Equal(t, "vip-test-id", stored.Headers["Authorization"].EndpointID)
	assert.Equal(t, "vip-test-id", stored.Headers["X-App"].EndpointID)
}
