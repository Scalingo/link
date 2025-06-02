package outscalepublicip

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Scalingo/link/v3/models"
	"github.com/Scalingo/link/v3/models/modelsmock"
)

func TestFactory_Mutate_Success(t *testing.T) {
	ctx := context.Background()
	ctrl := gomock.NewController(t)

	// Given a plugin config with sensitive data
	req := PluginConfig{
		AccessKey:  "my-access",
		SecretKey:  "my-secret",
		Region:     "eu-west-2",
		PublicIPID: "pip-123",
		NICID:      "nic-456",
	}
	raw, _ := json.Marshal(req)
	endpoint := models.Endpoint{PluginConfig: raw}

	mockStorage := modelsmock.NewMockEncryptedStorage(ctrl)
	mockStorage.EXPECT().Encrypt(ctx, "my-access").Return(models.EncryptedData{
		Type: models.EncryptedDataTypeAESCFB,
		Data: "enc_my-access",
	}, nil)
	mockStorage.EXPECT().Encrypt(ctx, "my-secret").Return(models.EncryptedData{
		Type: models.EncryptedDataTypeAESCFB,
		Data: "enc_my-secret",
	}, nil)

	f := Factory{encryptedStorage: mockStorage}

	// When we mutate the plugin config
	res, err := f.Mutate(ctx, endpoint)
	require.NoError(t, err)

	// It should encrypt the sensitive data and keep the rest
	var stored StorablePluginConfig
	err = json.Unmarshal(res, &stored)
	require.NoError(t, err)
	assert.Equal(t, models.EncryptedDataTypeAESCFB, stored.AccessKey.Type)
	assert.Equal(t, models.EncryptedDataTypeAESCFB, stored.SecretKey.Type)
	assert.Equal(t, "enc_my-access", stored.AccessKey.Data)
	assert.Equal(t, "enc_my-secret", stored.SecretKey.Data)
	assert.Equal(t, req.Region, stored.Region)
	assert.Equal(t, req.PublicIPID, stored.PublicIPID)
	assert.Equal(t, req.NICID, stored.NICID)
}
