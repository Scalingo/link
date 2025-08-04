package upgrades

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	v2api "github.com/Scalingo/link/v2/api"
	"github.com/Scalingo/link/v3/api"
	"github.com/Scalingo/link/v3/tests/integration/utils"
)

func Test_UpdateFrom2ToCurrent(t *testing.T) {
	t.Run("Upgrading from the latest v2 to the current v3", func(t *testing.T) {
		utils.CleanupEtcdData(t)
		t.Cleanup(func() {
			utils.CleanupEtcdData(t)
		})

		// Start a LinK v2
		oldBinPath := utils.DownloadLinKVersion(t, "2.0.7")
		oldLink := utils.StartLinK(t, oldBinPath)

		// Create an IP in v2
		v2client := v2api.NewHTTPClient(v2api.WithURL(oldLink.URL()))
		ips, err := v2client.ListIPs(t.Context())
		require.NoError(t, err)
		assert.Empty(t, ips)

		_, err = v2client.AddIP(t.Context(), "10.20.0.1/32", v2api.AddIPParams{})
		require.NoError(t, err)
		ips, err = v2client.ListIPs(t.Context())
		require.NoError(t, err)
		assert.Len(t, ips, 1)
		assert.Equal(t, "10.20.0.1/32", ips[0].IP.IP)

		oldLink.Stop(t)

		// Start a LinK v3
		newBinPath := utils.BuildLinkBinary(t)
		newLink := utils.StartLinK(t, newBinPath, utils.WithEnv("SECRET_STORAGE_ENCRYPTION_KEY", "a-very-long-encryption-key-1234567890abcdef-1234567890abcdef"))

		// Check that the IP is still there
		client := api.NewHTTPClient(api.WithURL(newLink.URL()))
		endpoints, err := client.ListEndpoints(t.Context())
		require.NoError(t, err)
		assert.Len(t, endpoints, 1)

		assert.Equal(t, "10.20.0.1_32", endpoints[0].ElectionKey)
	})
}
