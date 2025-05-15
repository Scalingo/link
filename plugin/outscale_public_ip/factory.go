package outscalepublicip

import (
	"context"
	"encoding/json"
	"time"

	osc "github.com/outscale/osc-sdk-go/v2"

	"github.com/kelseyhightower/envconfig"

	"github.com/Scalingo/go-utils/errors/v2"
	"github.com/Scalingo/link/v2/models"
	"github.com/Scalingo/link/v2/plugin"
)

const Name = "outscale_public_ip"

type Config struct {
	RefreshEvery time.Duration `envconfig:"OUTSCALE_PUBLIC_IP_REFRESH_INTERVAL" default:"1m"`
}

func Register(ctx context.Context, registry plugin.Registry, encryptedStorage models.EncryptedStorage) error {
	var config Config
	err := envconfig.Process("", &config)
	if err != nil {
		return errors.Wrap(ctx, err, "parse environment")
	}

	registry.Register(ctx, Name, Factory{
		config:           config,
		encryptedStorage: encryptedStorage,
	})

	return nil
}

type Factory struct {
	config           Config
	encryptedStorage models.EncryptedStorage
}

func (f Factory) Create(ctx context.Context, endpoint models.Endpoint) (plugin.Plugin, error) {
	oscClient := osc.NewAPIClient(osc.NewConfiguration())

	var cfg StorablePluginConfig
	err := json.Unmarshal(endpoint.PluginConfig, &cfg)
	if err != nil {
		return nil, errors.Wrap(ctx, err, "unmarshal plugin config")
	}

	plugin := &Plugin{
		oscClient:    oscClient,
		refreshEvery: f.config.RefreshEvery,
		region:       cfg.Region,
		publicIPID:   cfg.PublicIPID,
		nicID:        cfg.NICID,
	}

	err = f.encryptedStorage.Decrypt(ctx, cfg.AccessKey, &plugin.accessKey)
	if err != nil {
		return nil, errors.Wrap(ctx, err, "decrypt access key")
	}
	err = f.encryptedStorage.Decrypt(ctx, cfg.SecretKey, &plugin.secretKey)
	if err != nil {
		return nil, errors.Wrap(ctx, err, "decrypt secret key")
	}

	return plugin, nil
}

type PluginConfig struct {
	AccessKey string `json:"access_key"`
	SecretKey string `json:"secret_key"`
	Region    string `json:"region"`

	PublicIPID string `json:"public_ip_id"`
	NICID      string `json:"nic_id"`
}

func (f Factory) Validate(_ context.Context, endpoint models.Endpoint) error {
	validations := errors.NewValidationErrorsBuilder()
	var req PluginConfig
	err := json.Unmarshal(endpoint.PluginConfig, &req)
	if err != nil {
		validations.Set("plugin_config", "invalid JSON: "+err.Error())
		return validations.Build()
	}

	if req.AccessKey == "" {
		validations.Set("plugin_config.access_key", "missing access key")
	}
	if req.SecretKey == "" {
		validations.Set("plugin_config.secret_key", "missing secret key")
	}
	if req.Region == "" {
		validations.Set("plugin_config.region", "missing region")
	}
	if req.PublicIPID == "" {
		validations.Set("plugin_config.public_ip_id", "missing public IP ID")
	}
	if req.NICID == "" {
		validations.Set("plugin_config.nic_id", "missing NIC ID")
	}

	validationErr := validations.Build()
	if validationErr != nil {
		return validationErr
	}

	return nil
}

func (f Factory) Mutate(ctx context.Context, endpoint models.Endpoint) (json.RawMessage, error) {
	var req PluginConfig

	err := json.Unmarshal(endpoint.PluginConfig, &req)
	if err != nil {
		return nil, errors.Wrap(ctx, err, "unmarshal plugin config")
	}

	cfg := StorablePluginConfig{
		Region:     req.Region,
		PublicIPID: req.PublicIPID,
		NICID:      req.NICID,
	}

	cfg.AccessKey, err = f.encryptedStorage.Encrypt(ctx, req.AccessKey)
	if err != nil {
		return nil, errors.Wrap(ctx, err, "encrypt access key")
	}
	cfg.SecretKey, err = f.encryptedStorage.Encrypt(ctx, req.SecretKey)
	if err != nil {
		return nil, errors.Wrap(ctx, err, "encrypt secret key")
	}

	res, _ := json.Marshal(cfg)

	return res, nil
}

type StorablePluginConfig struct {
	AccessKey models.EncryptedData `json:"access_key"`
	SecretKey models.EncryptedData `json:"secret_key"`
	Region    string               `json:"region"`

	PublicIPID string `json:"public_ip_id"`
	NICID      string `json:"nic_id"`
}
