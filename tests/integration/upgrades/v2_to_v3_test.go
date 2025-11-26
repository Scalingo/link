package upgrades

import (
	"slices"
	"testing"

	"github.com/stretchr/testify/require"

	v2api "github.com/Scalingo/link/v2/api"
	"github.com/Scalingo/link/v3/api"
	"github.com/Scalingo/link/v3/tests/integration/utils"
)

func Test_UpdateFromV2ToV3(t *testing.T) {
	t.Run("Upgrading from the latest v2 to the v3 on the current commit", func(t *testing.T) {
		utils.CleanupEtcdData(t)
		t.Cleanup(func() {
			utils.CleanupEtcdData(t)
		})

		// Start a LinK v2
		oldBinPath := utils.DownloadLinKVersion(t, "2.0.7")
		linkV2Process := utils.StartLinK(t, oldBinPath)

		// Create an IP in v2
		v2client := v2api.NewHTTPClient(v2api.WithURL(linkV2Process.URL()))

		expectedIP, err := v2client.AddIP(t.Context(), "10.20.0.1/32", v2api.AddIPParams{})
		require.NoError(t, err)
		ips, err := v2client.ListIPs(t.Context())
		require.NoError(t, err)
		containsIP := slices.ContainsFunc(ips, func(ip v2api.IP) bool {
			return expectedIP.ID == ip.ID
		})
		require.True(t, containsIP)

		linkV2Process.Stop(t)

		// Start a LinK v3
		newBinPath := utils.BuildLinKBinary(t)
		linkV3Process := utils.StartLinK(t, newBinPath,
			utils.WithEnv("SECRET_STORAGE_ENCRYPTION_KEY", "a-very-long-encryption-key-1234567890abcdef-1234567890abcdef"),
		)

		// Check that the IP is still there
		v3client := api.NewHTTPClient(api.WithURL(linkV3Process.URL()))
		endpoints, err := v3client.ListEndpoints(t.Context())
		require.NoError(t, err)
		containsIP = slices.ContainsFunc(endpoints, func(endpoint api.Endpoint) bool {
			return expectedIP.StorableIP() == endpoint.ElectionKey
		})
		require.True(t, containsIP)
	})
}
