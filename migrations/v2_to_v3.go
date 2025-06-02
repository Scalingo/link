package migrations

import (
	"context"
	"encoding/json"

	"github.com/Scalingo/go-utils/errors/v2"
	"github.com/Scalingo/go-utils/logger"
	"github.com/Scalingo/link/v3/models"
	"github.com/Scalingo/link/v3/plugin/arp"
)

type V2toV3 struct {
	hostname string
	storage  models.Storage
}

func NewV2toV3Migration(hostname string, storage models.Storage) V2toV3 {
	return V2toV3{
		hostname: hostname,
		storage:  storage,
	}
}

func (m V2toV3) Name() string {
	return "v2 to v3"
}

func (m V2toV3) NeedsMigration(ctx context.Context) (bool, error) {
	log := logger.Get(ctx)

	endpoints, err := m.storage.GetEndpoints(ctx)
	if err != nil {
		return false, errors.Wrap(ctx, err, "fail to get endpoints to check if it needs data migration from v2 to v3")
	}

	for _, endpoint := range endpoints {
		if endpoint.Plugin == "" {
			log.Info("Current host needs data migration from v2 to v3")
			return true, nil
		}
	}
	log.Info("Current host does not need data migration from v2 to v3")
	return false, nil
}

func (m V2toV3) Migrate(ctx context.Context) error {
	log := logger.Get(ctx)
	log.Info("Migrate data from v2 to v3")

	endpoints, err := m.storage.GetEndpoints(ctx)
	if err != nil {
		return errors.Wrap(ctx, err, "get endpoints")
	}

	for _, endpoint := range endpoints {
		if endpoint.Plugin != "" {
			continue
		}
		ctx, log := logger.WithStructToCtx(ctx, "endpoint", endpoint)
		log.Info("Migrating endpoint to v3")
		pluginConfig, _ := json.Marshal(arp.PluginConfig{
			IP: endpoint.IP,
		})

		endpoint.Plugin = "arp"
		endpoint.PluginConfig = pluginConfig

		err = m.storage.UpdateEndpoint(ctx, endpoint)
		if err != nil {
			return errors.Wrap(ctx, err, "update endpoint")
		}
	}

	host, err := m.storage.GetCurrentHost(ctx)
	if err != nil {
		return errors.Wrap(ctx, err, "get current host")
	}
	host.DataVersion = 3

	err = m.storage.SaveHost(ctx, host)
	if err != nil {
		return errors.Wrap(ctx, err, "save host")
	}

	log.Info("Data migration from v2 to v3 done")
	return nil
}
