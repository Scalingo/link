package integration

import (
	"encoding/json"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Scalingo/link/v3/api"
	"github.com/Scalingo/link/v3/config"
	"github.com/Scalingo/link/v3/models"
	outscalepublicip "github.com/Scalingo/link/v3/plugin/outscale_public_ip"
	"github.com/Scalingo/link/v3/tests/integration/utils"
)

func Test_SecretStorageUpgrade(t *testing.T) {
	// When working on 3.1.0, we had to change the secret storage format.
	// This test ensures that the upgrade from 3.0.1 to the current version works correctly.

	t.Run("Upgrading from 3.0.1 to current version", func(t *testing.T) {
		utils.CleanupEtcdData(t)
		t.Cleanup(func() {
			utils.CleanupEtcdData(t)
		})

		encryptionKey := "a-very-long-encryption-key-1234567890abcdef-1234567890abcdef"
		cfg := config.Config{
			Hostname:                   "test-host",
			SecretStorageEncryptionKey: encryptionKey,
		}
		storage := models.NewEtcdStorage(cfg)

		// Start a LinK using the old version (3.0.1)
		oldBinPath := utils.DownloadLinKVersion(t, "3.0.1")
		oldLink := utils.StartLinK(t, oldBinPath,
			utils.WithEnv("SECRET_STORAGE_ENCRYPTION_KEY", encryptionKey),
		)

		// Create an endpoint with the old version
		// This will use the old encryption format.
		linkClient := api.NewHTTPClient(api.WithURL(oldLink.URL()))
		expectedEndpoint, err := linkClient.AddEndpoint(t.Context(), api.AddEndpointParams{
			Plugin: outscalepublicip.Name,
			PluginConfig: outscalepublicip.PluginConfig{
				AccessKey: "TESTACCESSKEY",
				SecretKey: "TESTSECRETKEY",
				Region:    "test-region",

				PublicIPID: "ip-abc123",
				NICID:      "nic-abc123",
			},
		})
		require.NoError(t, err)

		// Check that the endpoint is stored with the old encryption format
		storedEndpoints, err := storage.GetEndpoints(t.Context())
		require.NoError(t, err)
		endpointIndex := slices.IndexFunc(storedEndpoints, func(endpoint models.Endpoint) bool {
			return expectedEndpoint.ID == endpoint.ID
		})
		require.NotEqual(t, -1, endpointIndex)

		var pluginConfig outscalepublicip.StorablePluginConfig
		err = json.Unmarshal(storedEndpoints[endpointIndex].PluginConfig, &pluginConfig)
		require.NoError(t, err)

		assert.NotEmpty(t, pluginConfig.AccessKey.Data)
		assert.NotEmpty(t, pluginConfig.SecretKey.Data)
		assert.Equal(t, "aes-cfb", pluginConfig.AccessKey.Type)
		assert.Equal(t, "aes-cfb", pluginConfig.SecretKey.Type)

		oldLink.Stop(t)

		// Start a LinK using the current version and the same encryption key

		newBinPath := utils.BuildLinKBinary(t)
		newLink := utils.StartLinK(t, newBinPath, utils.WithEnv("SECRET_STORAGE_ENCRYPTION_KEY", encryptionKey))

		linkClient = api.NewHTTPClient(api.WithURL(newLink.URL()))
		endpoints, err := linkClient.ListEndpoints(t.Context())
		require.NoError(t, err)
		endpointIndex = slices.IndexFunc(endpoints, func(endpoint api.Endpoint) bool {
			return expectedEndpoint.ID == endpoint.ID
		})
		require.NotEqual(t, -1, endpointIndex)
		assert.Equal(t, outscalepublicip.Name, endpoints[endpointIndex].Plugin)

		encryptedStorage, err := models.NewEncryptedStorage(t.Context(), cfg, storage)
		require.NoError(t, err)

		storedEndpoints, err = storage.GetEndpoints(t.Context())
		require.NoError(t, err)
		endpointIndex = slices.IndexFunc(storedEndpoints, func(endpoint models.Endpoint) bool {
			return expectedEndpoint.ID == endpoint.ID
		})
		require.NotEqual(t, -1, endpointIndex)

		pluginConfig = outscalepublicip.StorablePluginConfig{}

		err = json.Unmarshal(storedEndpoints[endpointIndex].PluginConfig, &pluginConfig)
		require.NoError(t, err)

		assert.Empty(t, pluginConfig.AccessKey.Data)
		assert.Empty(t, pluginConfig.SecretKey.Data)
		assert.Empty(t, pluginConfig.AccessKey.Type)
		assert.Empty(t, pluginConfig.SecretKey.Type)

		assert.NotEmpty(t, pluginConfig.AccessKey.ID)
		assert.NotEmpty(t, pluginConfig.SecretKey.ID)

		var accessKey, secretKey string
		err = encryptedStorage.Decrypt(t.Context(), pluginConfig.AccessKey, &accessKey)
		require.NoError(t, err)
		err = encryptedStorage.Decrypt(t.Context(), pluginConfig.SecretKey, &secretKey)
		require.NoError(t, err)
		assert.Equal(t, "TESTACCESSKEY", accessKey)
		assert.Equal(t, "TESTSECRETKEY", secretKey)
	})
}
