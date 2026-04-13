package webhook

import (
	"context"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/Scalingo/go-utils/errors/v2"
	"github.com/Scalingo/link/v3/models"
	"github.com/Scalingo/link/v3/plugin"
)

const Name = "webhook"

type Factory struct {
	httpClient       *http.Client
	encryptedStorage models.EncryptedStorage
}

type PluginConfig struct {
	URL     string            `json:"url"`
	Headers map[string]string `json:"headers,omitempty"`
}

type StorablePluginConfig struct {
	URL     string                              `json:"url"`
	Headers map[string]models.EncryptedDataLink `json:"headers,omitempty"`
}

func Register(ctx context.Context, registry plugin.Registry, encryptedStorage models.EncryptedStorage) error {
	registry.Register(ctx, Name, Factory{
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

	return &Plugin{
		endpoint:   endpoint,
		cfg:        cfg,
		httpClient: httpClient,
	}, nil
}

func (f Factory) Validate(ctx context.Context, endpoint models.Endpoint) error {
	_, err := parseConfig(ctx, endpoint)
	if err != nil {
		return errors.Wrap(ctx, err, "invalid plugin config")
	}

	return nil
}

func (f Factory) Mutate(ctx context.Context, endpoint models.Endpoint) (json.RawMessage, error) {
	cfg, err := parseConfig(ctx, endpoint)
	if err != nil {
		return nil, errors.Wrap(ctx, err, "parse config")
	}

	if f.encryptedStorage == nil {
		return nil, errors.New(ctx, "encrypted storage is required")
	}

	storable := StorablePluginConfig{
		URL: cfg.URL,
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

	raw, err := json.Marshal(storable)
	if err != nil {
		return nil, errors.Wrap(ctx, err, "marshal storable config")
	}

	return raw, nil
}

func parseConfig(ctx context.Context, endpoint models.Endpoint) (PluginConfig, error) {
	cfg := PluginConfig{}
	err := json.Unmarshal(endpoint.PluginConfig, &cfg)
	if err != nil {
		return cfg, errors.Wrap(ctx, err, "unmarshal plugin config")
	}

	err = validateURL(ctx, cfg.URL)
	if err != nil {
		return cfg, err
	}

	if cfg.Headers == nil {
		cfg.Headers = make(map[string]string)
	}

	return cfg, nil
}

func (f Factory) parseStorableConfig(ctx context.Context, endpoint models.Endpoint) (PluginConfig, error) {
	storable := StorablePluginConfig{}
	err := json.Unmarshal(endpoint.PluginConfig, &storable)
	if err != nil {
		return PluginConfig{}, errors.Wrap(ctx, err, "unmarshal plugin config")
	}

	err = validateURL(ctx, storable.URL)
	if err != nil {
		return PluginConfig{}, err
	}

	cfg := PluginConfig{
		URL:     storable.URL,
		Headers: make(map[string]string),
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

func validateURL(ctx context.Context, configURL string) error {
	builder := errors.NewValidationErrorsBuilder()
	if configURL == "" {
		builder.Set("url", "URL is required")
	} else {
		parsedURL, err := url.ParseRequestURI(configURL)
		if err != nil {
			builder.Set("url", "URL should be valid")
		} else {
			scheme := strings.ToLower(parsedURL.Scheme)
			if scheme != "http" && scheme != "https" {
				builder.Set("url", "URL scheme must be http or https")
			}
		}
	}

	if verr := builder.Build(); verr != nil {
		return errors.Wrap(ctx, verr, "validate config")
	}

	return nil
}
