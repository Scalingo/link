package outscalepublicip

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	osc "github.com/outscale/osc-sdk-go/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Scalingo/link/v3/services/outscale/outscalemock"
)

const (
	testPublicIPID = "pip-123"
	testNicID      = "nic-456"
	testLinkID     = "link-789"
)

func newPlugin(mockClient *outscalemock.MockPublicIPClient) *Plugin {
	return &Plugin{
		oscClient:    mockClient,
		refreshEvery: time.Minute,
		publicIPID:   testPublicIPID,
		nicID:        testNicID,
	}
}

func TestPlugin_Activate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockClient := outscalemock.NewMockPublicIPClient(ctrl)
		plugin := newPlugin(mockClient)

		resp := osc.LinkPublicIpResponse{}
		resp.SetLinkPublicIpId(testLinkID)

		mockClient.EXPECT().
			LinkPublicIP(gomock.Any(), osc.LinkPublicIpRequest{
				PublicIpId:  osc.PtrString(testPublicIPID),
				NicId:       osc.PtrString(testNicID),
				AllowRelink: osc.PtrBool(true),
			}).
			Return(resp, nil)

		err := plugin.Activate(context.Background())
		require.NoError(t, err)
		assert.Equal(t, testLinkID, plugin.linkPublicIPID)
	})

	t.Run("error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockClient := outscalemock.NewMockPublicIPClient(ctrl)
		plugin := newPlugin(mockClient)

		mockClient.EXPECT().
			LinkPublicIP(gomock.Any(), gomock.Any()).
			Return(osc.LinkPublicIpResponse{}, errors.New("link error"))

		err := plugin.Activate(context.Background())
		require.Error(t, err)
	})
}

func TestPlugin_Deactivate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockClient := outscalemock.NewMockPublicIPClient(ctrl)
		plugin := newPlugin(mockClient)
		plugin.linkPublicIPID = testLinkID

		mockClient.EXPECT().
			UnlinkPublicIP(gomock.Any(), osc.UnlinkPublicIpRequest{
				LinkPublicIpId: osc.PtrString(testLinkID),
			}).
			Return(osc.UnlinkPublicIpResponse{}, nil)

		err := plugin.Deactivate(context.Background())
		require.NoError(t, err)
		assert.Empty(t, plugin.linkPublicIPID)
	})

	t.Run("skip if not linked", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockClient := outscalemock.NewMockPublicIPClient(ctrl)
		plugin := newPlugin(mockClient)
		plugin.linkPublicIPID = ""

		// Should not call UnlinkPublicIP
		err := plugin.Deactivate(context.Background())
		require.NoError(t, err)
	})

	t.Run("error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockClient := outscalemock.NewMockPublicIPClient(ctrl)
		plugin := newPlugin(mockClient)
		plugin.linkPublicIPID = testLinkID

		mockClient.EXPECT().
			UnlinkPublicIP(gomock.Any(), gomock.Any()).
			Return(osc.UnlinkPublicIpResponse{}, errors.New("unlink error"))

		err := plugin.Deactivate(context.Background())
		require.Error(t, err)
	})
}

func TestPlugin_Ensure(t *testing.T) {
	t.Run("already refreshed", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockClient := outscalemock.NewMockPublicIPClient(ctrl)
		plugin := newPlugin(mockClient)
		lastRefreshedAt := time.Now()
		plugin.lastRefreshedAt = lastRefreshedAt

		err := plugin.Ensure(context.Background())
		require.NoError(t, err)
		assert.Equal(t, lastRefreshedAt, plugin.lastRefreshedAt, "lastRefreshedAt should not be updated")
	})

	t.Run("read public ip error", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockClient := outscalemock.NewMockPublicIPClient(ctrl)
		plugin := newPlugin(mockClient)
		lastRefreshedAt := time.Now().Add(-2 * time.Minute)
		plugin.lastRefreshedAt = lastRefreshedAt

		mockClient.EXPECT().
			ReadPublicIP(gomock.Any(), testPublicIPID).
			Return(osc.PublicIp{}, errors.New("read error"))

		err := plugin.Ensure(context.Background())
		require.Error(t, err)
		assert.Equal(t, lastRefreshedAt, plugin.lastRefreshedAt, "lastRefreshedAt should not be updated")
	})

	t.Run("a public ip that is not linked on the correct NIC should trigger Activate", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockClient := outscalemock.NewMockPublicIPClient(ctrl)
		plugin := newPlugin(mockClient)
		plugin.lastRefreshedAt = time.Now().Add(-2 * time.Minute)

		publicIP := osc.PublicIp{}
		publicIP.SetNicId("other-nic")
		publicIP.SetLinkPublicIpId("other-link")

		mockClient.EXPECT().
			ReadPublicIP(gomock.Any(), testPublicIPID).
			Return(publicIP, nil)

		resp := osc.LinkPublicIpResponse{}
		resp.SetLinkPublicIpId(testLinkID)
		mockClient.EXPECT().
			LinkPublicIP(gomock.Any(), osc.LinkPublicIpRequest{
				PublicIpId:  osc.PtrString(testPublicIPID),
				NicId:       osc.PtrString(testNicID),
				AllowRelink: osc.PtrBool(true),
			}).
			Return(resp, nil)

		err := plugin.Ensure(context.Background())
		require.NoError(t, err)
		assert.Equal(t, testLinkID, plugin.linkPublicIPID)
		assert.Greater(t, plugin.lastRefreshedAt, time.Now().Add(-1*time.Minute), "lastRefreshedAt should be updated")
	})

	t.Run("if the link ID was updated but we're on the same NIC", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockClient := outscalemock.NewMockPublicIPClient(ctrl)
		plugin := newPlugin(mockClient)
		plugin.lastRefreshedAt = time.Now().Add(-2 * time.Minute)
		plugin.linkPublicIPID = "old-link"

		publicIP := osc.PublicIp{}
		publicIP.SetNicId(testNicID)
		publicIP.SetLinkPublicIpId(testLinkID)

		mockClient.EXPECT().
			ReadPublicIP(gomock.Any(), testPublicIPID).
			Return(publicIP, nil)

		err := plugin.Ensure(context.Background())
		require.NoError(t, err)
		assert.Equal(t, testLinkID, plugin.linkPublicIPID)
	})
}
