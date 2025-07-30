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

func TestFactory_Validate(t *testing.T) {
	specs := []struct {
		Name          string
		Config        PluginConfig
		ExpectedError string
	}{
		{
			Name: "with a missing access key",
			Config: PluginConfig{
				AccessKey:  "",
				SecretKey:  "ABC1234",
				Region:     "eu-west-2",
				PublicIPID: "pip-123",
				NICID:      "nic-456",
			},
			ExpectedError: "missing access key",
		}, {
			Name: "with an invalid access key format",
			Config: PluginConfig{
				AccessKey:  "invalid-key",
				SecretKey:  "ABC1234",
				Region:     "eu-west-2",
				PublicIPID: "pip-123",
				NICID:      "nic-456",
			},
			ExpectedError: "invalid access key format",
		}, {
			Name: "with a missing secret key",
			Config: PluginConfig{
				AccessKey:  "ABC1234",
				SecretKey:  "",
				Region:     "eu-west-2",
				PublicIPID: "pip-123",
				NICID:      "nic-456",
			},
			ExpectedError: "missing secret key",
		}, {
			Name: "with an invalid secret key format",
			Config: PluginConfig{
				AccessKey:  "ABC1234",
				SecretKey:  "invalid-secret",
				Region:     "eu-west-2",
				PublicIPID: "pip-123",
				NICID:      "nic-456",
			},
			ExpectedError: "invalid secret key format",
		}, {
			Name: "with a missing region",
			Config: PluginConfig{
				AccessKey:  "ABC1234",
				SecretKey:  "ABC1234",
				Region:     "",
				PublicIPID: "pip-123",
				NICID:      "nic-456",
			},
			ExpectedError: "missing region",
		}, {
			Name: "with an invalid region",
			Config: PluginConfig{
				AccessKey:  "ABC1234",
				SecretKey:  "ABC1234",
				Region:     "invalid-region",
				PublicIPID: "pip-123",
				NICID:      "nic-456",
			},
			ExpectedError: "invalid region",
		}, {
			Name: "with a missing public IP ID",
			Config: PluginConfig{
				AccessKey:  "ABC1234",
				SecretKey:  "ABC1234",
				Region:     "eu-west-2",
				PublicIPID: "",
				NICID:      "nic-456",
			},
			ExpectedError: "missing public IP ID",
		}, {
			Name: "with an invalid public IP ID format",
			Config: PluginConfig{
				AccessKey:  "ABC1234",
				SecretKey:  "ABC1234",
				Region:     "eu-west-2",
				PublicIPID: "invalid-pip-id",
				NICID:      "nic-456",
			},
			ExpectedError: "invalid public IP ID format",
		}, {
			Name: "with a missing NIC ID",
			Config: PluginConfig{
				AccessKey:  "ABC1234",
				SecretKey:  "ABC1234",
				Region:     "eu-west-2",
				PublicIPID: "pip-123",
				NICID:      "",
			},
			ExpectedError: "missing NIC ID",
		}, {
			Name: "with an invalid NIC ID format",
			Config: PluginConfig{
				AccessKey:  "ABC1234",
				SecretKey:  "ABC1234",
				Region:     "eu-west-2",
				PublicIPID: "pip-123",
				NICID:      "invalid-nic-id",
			},
			ExpectedError: "invalid NIC ID format",
		}, {
			Name: "with a valid configuration",
			Config: PluginConfig{
				AccessKey:  "ABC1234",
				SecretKey:  "ABC1234",
				Region:     "eu-west-2",
				PublicIPID: "pip-123",
				NICID:      "nic-456",
			},
			ExpectedError: "",
		},
	}

	for _, spec := range specs {
		t.Run(spec.Name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockStorage := modelsmock.NewMockEncryptedStorage(ctrl)

			factory := Factory{
				encryptedStorage: mockStorage,
			}

			rawConfig, err := json.Marshal(spec.Config)
			require.NoError(t, err)

			endpoint := models.Endpoint{PluginConfig: rawConfig}

			err = factory.Validate(context.Background(), endpoint)
			if spec.ExpectedError != "" {
				require.Error(t, err)
				assert.Contains(t, err.Error(), spec.ExpectedError)
			} else {
				assert.NoError(t, err)
			}
		})
	}

}

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
