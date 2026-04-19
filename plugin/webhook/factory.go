package webhook

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/kelseyhightower/envconfig"

	"github.com/Scalingo/go-utils/errors/v2"
	"github.com/Scalingo/link/v3/api"
	"github.com/Scalingo/link/v3/models"
	"github.com/Scalingo/link/v3/plugin"
)

const Name = api.PluginWebhook

type Config struct {
	RefreshEvery time.Duration `envconfig:"WEBHOOK_REFRESH_INTERVAL" default:"5m"`
}

type Factory struct {
	config           Config
	httpClient       *http.Client
	encryptedStorage models.EncryptedStorage
}

type PluginConfig = api.WebhookPluginConfig

type StorablePluginConfig struct {
	URL        string                              `json:"url"`
	Headers    map[string]models.EncryptedDataLink `json:"headers,omitempty"`
	ResourceID string                              `json:"resource_id"`
}

func Register(ctx context.Context, registry plugin.Registry, encryptedStorage models.EncryptedStorage) error {
	var config Config
	err := envconfig.Process("", &config)
	if err != nil {
		return errors.Wrap(ctx, err, "parse environment")
	}

	registry.Register(ctx, Name, Factory{
		config:           config,
		httpClient:       &http.Client{Timeout: 5 * time.Second},
		encryptedStorage: encryptedStorage,
	})
	return nil
}

func (f Factory) Create(ctx context.Context, endpoint models.Endpoint) (plugin.Plugin, error) {
	cfg, err := f.parseStorableConfig(ctx, endpoint)
	if err != nil {
		return nil, errors.Wrap(ctx, err, "parse config")
	}

	httpClient := f.httpClient
	if httpClient == nil {
		httpClient = &http.Client{Timeout: 5 * time.Second}
	}

	refreshEvery := f.config.RefreshEvery
	if refreshEvery <= 0 {
		refreshEvery = 5 * time.Minute
	}

	return &Plugin{
		endpoint:     endpoint,
		cfg:          cfg,
		httpClient:   httpClient,
		refreshEvery: refreshEvery,
	}, nil
}

func (f Factory) Validate(_ context.Context, endpoint models.Endpoint) error {
	validations := errors.NewValidationErrorsBuilder()
	var cfg PluginConfig
	err := json.Unmarshal(endpoint.PluginConfig, &cfg)
	if err != nil {
		validations.Set("plugin_config", "invalid JSON: "+err.Error())
		return validations.Build()
	}

	if cfg.ResourceID == "" {
		validations.Set("plugin_config.resource_id", "missing resource ID")
	}

	if cfg.URL == "" {
		validations.Set("plugin_config.url", "missing URL")
	} else {
		parsedURL, err := url.ParseRequestURI(cfg.URL)
		if err != nil {
			validations.Set("plugin_config.url", "invalid URL")
		} else {
			scheme := strings.ToLower(parsedURL.Scheme)
			if scheme != "http" && scheme != "https" {
				validations.Set("plugin_config.url", "invalid URL scheme")
			}
		}
	}

	validationErr := validations.Build()
	if validationErr != nil {
		return validationErr
	}

	return nil
}

func (f Factory) Mutate(ctx context.Context, endpoint models.Endpoint) (json.RawMessage, error) {
	if f.encryptedStorage == nil {
		return nil, errors.New(ctx, "encrypted storage is required")
	}

	var cfg PluginConfig
	err := json.Unmarshal(endpoint.PluginConfig, &cfg)
	if err != nil {
		return nil, errors.Wrap(ctx, err, "unmarshal plugin config")
	}

	storable := StorablePluginConfig{
		URL:        cfg.URL,
		ResourceID: cfg.ResourceID,
	}

	if len(cfg.Headers) > 0 {
		storable.Headers = make(map[string]models.EncryptedDataLink, len(cfg.Headers))
	}

	for name, value := range cfg.Headers {
		encryptedHeader, err := f.encryptedStorage.Encrypt(ctx, endpoint.ID, value)
		if err != nil {
			return nil, errors.Wrap(ctx, err, "encrypt header "+name)
		}
		storable.Headers[name] = encryptedHeader
	}

	raw, _ := json.Marshal(storable)

	return raw, nil
}

func (f Factory) parseStorableConfig(ctx context.Context, endpoint models.Endpoint) (PluginConfig, error) {
	storable := StorablePluginConfig{}
	err := json.Unmarshal(endpoint.PluginConfig, &storable)
	if err != nil {
		return PluginConfig{}, errors.Wrap(ctx, err, "unmarshal plugin config")
	}

	cfg := PluginConfig{
		URL:        storable.URL,
		Headers:    make(map[string]string),
		ResourceID: storable.ResourceID,
	}

	if len(storable.Headers) == 0 {
		return cfg, nil
	}

	if f.encryptedStorage == nil {
		return PluginConfig{}, errors.New(ctx, "encrypted storage is required")
	}

	for name, encryptedHeader := range storable.Headers {
		var value string
		err := f.encryptedStorage.Decrypt(ctx, encryptedHeader, &value)
		if err != nil {
			return PluginConfig{}, errors.Wrap(ctx, err, "decrypt header "+name)
		}
		cfg.Headers[name] = value
	}

	return cfg, nil
}
