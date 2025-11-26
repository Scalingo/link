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

	binPath := utils.BuildLinKBinary(t)
	linkProcess := utils.StartLinK(t, binPath,
		utils.WithEnv("SECRET_STORAGE_ENCRYPTION_KEY", oldEncryptionKey),
	)

	linkClient := api.NewHTTPClient(api.WithURL(linkProcess.URL()))

	_, err := linkClient.AddEndpoint(t.Context(), api.AddEndpointParams{
		Plugin: outscalepublicip.Name,
		PluginConfig: outscalepublicip.PluginConfig{
			AccessKey:  "TESTACCESSKEY",
			SecretKey:  "TESTSECRETKEY",
			Region:     "test-region",
			PublicIPID: "ip-abc123",
			NICID:      "nic-abc123",
		},
	})
	require.NoError(t, err)

	endpoints, err := linkClient.ListEndpoints(t.Context())
	require.NoError(t, err)
	require.Len(t, endpoints, 1)

	linkProcess.Stop(t)

	newEncryptionKey := "another-very-long-encryption-key-1234567890abcdef-1234567890abcdef"
	linkProcess = utils.StartLinK(t, binPath,
		utils.WithEnv("SECRET_STORAGE_ENCRYPTION_KEY", newEncryptionKey),
		utils.WithEnv("SECRET_STORAGE_ALTERNATE_KEYS", oldEncryptionKey),
	)
	linkClient = api.NewHTTPClient(api.WithURL(linkProcess.URL()))

	err = linkClient.RotateEncryptionKey(t.Context())
	require.NoError(t, err)

	endpoints, err = linkClient.ListEndpoints(t.Context())
	require.NoError(t, err)
	require.Len(t, endpoints, 1)

	linkProcess.Stop(t)

	// Restart with the new key only
	linkProcess = utils.StartLinK(t, binPath, utils.WithEnv("SECRET_STORAGE_ENCRYPTION_KEY", newEncryptionKey))
	linkClient = api.NewHTTPClient(api.WithURL(linkProcess.URL()))

	endpoints, err = linkClient.ListEndpoints(t.Context())
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
	assert.Equal(t, "TESTACCESSKEY", accessKey)
	assert.Equal(t, "TESTSECRETKEY", secretKey)

	config.SecretStorageEncryptionKey = oldEncryptionKey
	secretStorage, err = models.NewEncryptedStorage(t.Context(), config, storage)
	require.NoError(t, err)
	err = secretStorage.Decrypt(t.Context(), pluginConfig.AccessKey, &accessKey)
	require.Error(t, err)
	err = secretStorage.Decrypt(t.Context(), pluginConfig.SecretKey, &secretKey)
	require.Error(t, err)
}
