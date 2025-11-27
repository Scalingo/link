package integration

import (
	"encoding/json"
	"slices"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	linkapi "github.com/Scalingo/link/v3/api"
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

	binPath := utils.BuildLinKBinary(t)

	oldEncryptionKey := "a-very-long-encryption-key-1234567890abcdef-1234567890abcdef"
	linkProcess := utils.StartLinK(t, binPath,
		utils.WithEnv("SECRET_STORAGE_ENCRYPTION_KEY", oldEncryptionKey),
	)

	// Initialize LinK with a endpoint
	linkClient := linkapi.NewHTTPClient(linkapi.WithURL(linkProcess.URL()))
	expectedEndpoint, err := linkClient.AddEndpoint(t.Context(), linkapi.AddEndpointParams{
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
	containsEndpoint := slices.ContainsFunc(endpoints, func(endpoint linkapi.Endpoint) bool {
		return expectedEndpoint.ID == endpoint.ID
	})
	require.True(t, containsEndpoint)

	// Stop LinK...
	linkProcess.Stop(t)
	// ... and start it again with a new encryption key to test the key rotation
	newEncryptionKey := "another-very-long-encryption-key-1234567890abcdef-1234567890abcdef"
	linkProcess = utils.StartLinK(t, binPath,
		utils.WithEnv("SECRET_STORAGE_ENCRYPTION_KEY", newEncryptionKey),
		utils.WithEnv("SECRET_STORAGE_ALTERNATE_KEYS", oldEncryptionKey),
	)
	linkClient = linkapi.NewHTTPClient(linkapi.WithURL(linkProcess.URL()))

	err = linkClient.RotateEncryptionKey(t.Context())
	require.NoError(t, err)

	endpoints, err = linkClient.ListEndpoints(t.Context())
	require.NoError(t, err)
	containsEndpoint = slices.ContainsFunc(endpoints, func(endpoint linkapi.Endpoint) bool {
		return expectedEndpoint.ID == endpoint.ID
	})
	require.True(t, containsEndpoint)

	// Stop Link again...
	linkProcess.Stop(t)
	// ... and restart it with the new encryption key only
	linkProcess = utils.StartLinK(t, binPath,
		utils.WithEnv("SECRET_STORAGE_ENCRYPTION_KEY", newEncryptionKey),
	)
	linkClient = linkapi.NewHTTPClient(linkapi.WithURL(linkProcess.URL()))

	endpoints, err = linkClient.ListEndpoints(t.Context())
	require.NoError(t, err)
	containsEndpoint = slices.ContainsFunc(endpoints, func(endpoint linkapi.Endpoint) bool {
		return expectedEndpoint.ID == endpoint.ID
	})
	require.True(t, containsEndpoint)

	config := config.Config{
		Hostname:                   "test-host",
		SecretStorageEncryptionKey: newEncryptionKey,
	}
	storage := models.NewEtcdStorage(config)

	// Check in the storage that the data are still accessible despite the key rotation
	storedEndpoints, err := storage.GetEndpoints(t.Context())
	require.NoError(t, err)
	endpointIndex := slices.IndexFunc(storedEndpoints, func(endpoint models.Endpoint) bool {
		return expectedEndpoint.ID == endpoint.ID
	})
	require.NotEqual(t, -1, endpointIndex)

	var pluginConfig outscalepublicip.StorablePluginConfig
	err = json.Unmarshal(storedEndpoints[endpointIndex].PluginConfig, &pluginConfig)
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
