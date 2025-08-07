package migrations

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"sync"

	"github.com/Scalingo/go-utils/errors/v3"
	"github.com/Scalingo/go-utils/logger"
	"github.com/Scalingo/link/v3/config"
	"github.com/Scalingo/link/v3/models"
	outscalepublicip "github.com/Scalingo/link/v3/plugin/outscale_public_ip"
)

type endpointWithConfig struct {
	endpoint models.Endpoint
	config   outscalepublicip.StorablePluginConfig
}

type V3toV4 struct {
	hostname         string
	storage          models.Storage
	encryptionKey    []byte
	encryptedStorage models.EncryptedStorage

	endpointsFetch     *sync.Once
	endpointsToMigrate []endpointWithConfig
}

func NewV3toV4Migration(ctx context.Context, cfg config.Config, storage models.Storage) (*V3toV4, error) {
	encryptedStorage, err := models.NewEncryptedStorage(ctx, cfg, storage)
	if err != nil {
		return nil, errors.Wrap(ctx, err, "init encrypted storage")
	}
	if len(cfg.SecretStorageEncryptionKey) < 32 {
		return nil, errors.New(ctx, "encryption key must be at least 32 bytes long")
	}

	key := sha256.Sum256([]byte(cfg.SecretStorageEncryptionKey))

	return &V3toV4{
		hostname:           cfg.Hostname,
		storage:            storage,
		encryptedStorage:   encryptedStorage,
		encryptionKey:      key[:],
		endpointsToMigrate: []endpointWithConfig{},
		endpointsFetch:     &sync.Once{},
	}, nil
}

func (m *V3toV4) Name() string {
	return "v3 to v4"
}

func (m *V3toV4) NeedsMigration(ctx context.Context) (bool, error) {
	log := logger.Get(ctx)

	endpoints, err := m.getEndpointsToMigrate(ctx)
	if err != nil {
		return false, errors.Wrap(ctx, err, "get endpoints to check if it needs data migration from v3 to v4")
	}

	if len(endpoints) == 0 {
		log.Info("No endpoints need data migration from v3 to v4")
		return false, nil
	}

	log.Infof("Found %d endpoints that need data migration from v3 to v4", len(endpoints))
	return true, nil
}

func (m *V3toV4) Migrate(ctx context.Context) error {
	log := logger.Get(ctx)
	log.Info("Migrate data from v3 to v4")

	endpoints, err := m.getEndpointsToMigrate(ctx)
	if err != nil {
		return errors.Wrap(ctx, err, "get endpoints to migrate from v3 to v4")
	}

	for _, toMigrate := range endpoints {
		endpoint := toMigrate.endpoint
		pluginConfig := toMigrate.config

		ctx, log := logger.WithStructToCtx(ctx, "endpoint", endpoint)
		log.Info("Migrating endpoint to v4")

		if pluginConfig.AccessKey.ID == "" {
			log.Info("Migrating access key")
			var accessKey string
			err := m.encryptedStorage.Decrypt(ctx, pluginConfig.AccessKey, &accessKey)
			if err != nil {
				return errors.Wrap(ctx, err, "decrypt old access key")
			}
			pluginConfig.AccessKey, err = m.encryptedStorage.Encrypt(ctx, endpoint.ID, accessKey)
			if err != nil {
				return errors.Wrap(ctx, err, "encrypt access key")
			}
		}

		if pluginConfig.SecretKey.ID == "" {
			log.Info("Migrating secret key")
			var secretKey string

			err := m.encryptedStorage.Decrypt(ctx, pluginConfig.SecretKey, &secretKey)
			if err != nil {
				return errors.Wrap(ctx, err, "decrypt old secret key")
			}

			pluginConfig.SecretKey, err = m.encryptedStorage.Encrypt(ctx, endpoint.ID, secretKey)
			if err != nil {
				return errors.Wrap(ctx, err, "encrypt secret key")
			}
		}

		endpoint.PluginConfig, err = json.Marshal(pluginConfig)
		if err != nil {
			return errors.Wrap(ctx, err, "marshal new plugin config")
		}

		err = m.storage.UpdateEndpoint(ctx, endpoint)
		if err != nil {
			return errors.Wrap(ctx, err, "update endpoint")
		}
	}
	return nil
}

func (m *V3toV4) getEndpointsToMigrate(ctx context.Context) ([]endpointWithConfig, error) {
	var fetchErr error

	m.endpointsFetch.Do(func() {
		m.endpointsToMigrate = make([]endpointWithConfig, 0)

		endpoints, err := m.storage.GetEndpoints(ctx)
		if err != nil {
			fetchErr = errors.Wrap(ctx, err, "get endpoints to migrate")
			return
		}

		for _, endpoint := range endpoints {
			ctx, _ = logger.WithStructToCtx(ctx, "endpoint", endpoint)

			if endpoint.Plugin != outscalepublicip.Name {
				continue
			}

			var pluginConfig outscalepublicip.StorablePluginConfig
			err := json.Unmarshal(endpoint.PluginConfig, &pluginConfig)
			if err != nil {
				fetchErr = errors.Wrap(ctx, err, "unmarshal plugin config")
				return
			}

			if pluginConfig.AccessKey.ID != "" && pluginConfig.SecretKey.ID != "" {
				continue
			}

			m.endpointsToMigrate = append(m.endpointsToMigrate, endpointWithConfig{
				endpoint: endpoint,
				config:   pluginConfig,
			})
		}
	})

	if fetchErr != nil {
		return nil, errors.Wrap(ctx, fetchErr, "fetch endpoints to migrate")
	}

	return m.endpointsToMigrate, nil
}
