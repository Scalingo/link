package secrets

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Scalingo/link/v3/api"
	"github.com/Scalingo/link/v3/config"
	"github.com/Scalingo/link/v3/models"
	outscalepublicip "github.com/Scalingo/link/v3/plugin/outscale_public_ip"
	"github.com/Scalingo/link/v3/tests/integration/utils"
)

func TestKeyRotation(t *testing.T) {
	utils.CleanupEtcdData(t)
	t.Cleanup(func() {
		utils.CleanupEtcdData(t)
	})

	oldEncryptionKey := "a-very-long-encryption-key-1234567890abcdef-1234567890abcdef"

	binPath := utils.BuildLinkBinary(t)
	link := utils.StartLinK(t, binPath, utils.WithEnv("SECRET_STORAGE_ENCRYPTION_KEY", oldEncryptionKey))

	client := api.NewHTTPClient(api.WithURL(link.URL()))
	_, err := client.AddEndpoint(t.Context(), api.AddEndpointParams{
		Plugin: outscalepublicip.Name,
		PluginConfig: outscalepublicip.PluginConfig{
			AccessKey:  "test-access-key",
			SecretKey:  "test-secret-key",
			Region:     "dev",
			PublicIPID: "test-public-ip-id",
			NICID:      "test-nic-id",
		},
	})
	require.NoError(t, err)
	endpoints, err := client.ListEndpoints(t.Context())
	require.NoError(t, err)
	require.Len(t, endpoints, 1)

	link.Stop(t)

	newEncryptionKey := "another-very-long-encryption-key-1234567890abcdef-1234567890abcdef"
	link = utils.StartLinK(t, binPath,
		utils.WithEnv("SECRET_STORAGE_ENCRYPTION_KEY", newEncryptionKey),
		utils.WithEnv("SECRET_STORAGE_ALTERNATE_KEYS", oldEncryptionKey),
	)
	client = api.NewHTTPClient(api.WithURL(link.URL()))

	err = client.RotateEncryptionKey(t.Context())
	require.NoError(t, err)

	endpoints, err = client.ListEndpoints(t.Context())
	require.NoError(t, err)
	require.Len(t, endpoints, 1)

	link.Stop(t)

	// Restart with the new key only
	link = utils.StartLinK(t, binPath, utils.WithEnv("SECRET_STORAGE_ENCRYPTION_KEY", newEncryptionKey))
	client = api.NewHTTPClient(api.WithURL(link.URL()))

	endpoints, err = client.ListEndpoints(t.Context())
	require.NoError(t, err)
	require.Len(t, endpoints, 1)

	config := config.Config{
		Hostname:                   "test-host",
		SecretStorageEncryptionKey: newEncryptionKey,
	}
	storage := models.NewEtcdStorage(config)

	storedEndpoints, err := storage.GetEndpoints(t.Context())
	require.NoError(t, err)
	require.Len(t, storedEndpoints, 1)

	var pluginConfig outscalepublicip.StorablePluginConfig
	err = json.Unmarshal(storedEndpoints[0].PluginConfig, &pluginConfig)
	require.NoError(t, err)

	var accessKey, secretKey string
	secretStorage, err := models.NewEncryptedStorage(t.Context(), config, storage)
	require.NoError(t, err)
	err = secretStorage.Decrypt(t.Context(), pluginConfig.AccessKey, &accessKey)
	require.NoError(t, err)
	err = secretStorage.Decrypt(t.Context(), pluginConfig.SecretKey, &secretKey)
	require.NoError(t, err)
	assert.Equal(t, "test-access-key", accessKey)
	assert.Equal(t, "test-secret-key", secretKey)

	config.SecretStorageEncryptionKey = oldEncryptionKey
	secretStorage, err = models.NewEncryptedStorage(t.Context(), config, storage)
	require.NoError(t, err)
	err = secretStorage.Decrypt(t.Context(), pluginConfig.AccessKey, &accessKey)
	require.Error(t, err)
	err = secretStorage.Decrypt(t.Context(), pluginConfig.SecretKey, &secretKey)
	require.Error(t, err)
}
