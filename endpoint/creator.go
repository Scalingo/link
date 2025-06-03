package endpoint

import (
	"context"
	"encoding/json"

	"github.com/Scalingo/go-utils/errors/v2"
	"github.com/Scalingo/go-utils/logger"
	"github.com/Scalingo/link/v3/api"
	"github.com/Scalingo/link/v3/models"
	"github.com/Scalingo/link/v3/plugin"
	"github.com/Scalingo/link/v3/scheduler"
)

type CreateEndpointParams struct {
	HealthCheckInterval int               `json:"healthcheck_interval"`
	Checks              []api.HealthCheck `json:"checks"`
	Plugin              string            `json:"plugin_name"`
	PluginConfig        json.RawMessage   `json:"plugin_config"`
}

type Creator interface {
	CreateEndpoint(ctx context.Context, params CreateEndpointParams) (models.Endpoint, error)
}

type creator struct {
	storage   models.Storage
	scheduler scheduler.Scheduler
	registry  plugin.Registry
}

func NewCreator(storage models.Storage, scheduler scheduler.Scheduler, registry plugin.Registry) Creator {
	return &creator{
		storage:   storage,
		scheduler: scheduler,
		registry:  registry,
	}
}

func (c *creator) CreateEndpoint(ctx context.Context, params CreateEndpointParams) (models.Endpoint, error) {
	log := logger.Get(ctx)
	checks := models.HealthChecksFromAPIType(params.Checks)

	log.Info("Validating Health checks")
	validationErr := checks.Validate(ctx)
	if validationErr != nil {
		return models.Endpoint{}, errors.Wrap(ctx, validationErr, "validate health checks")
	}

	endpoint := models.Endpoint{
		HealthCheckInterval: params.HealthCheckInterval,
		Checks:              checks,
		Plugin:              params.Plugin,
		PluginConfig:        params.PluginConfig,
	}

	log.Info("Validating plugin")

	err := c.registry.Validate(ctx, endpoint)
	if err != nil {
		return endpoint, errors.Wrap(ctx, err, "validate plugin")
	}

	pluginConfig, err := c.registry.Mutate(ctx, endpoint)
	if err != nil {
		return endpoint, errors.Wrap(ctx, err, "mutate plugin")
	}
	endpoint.PluginConfig = pluginConfig

	log.Info("Creating endpoint in database")

	endpoint, err = c.storage.AddEndpoint(ctx, endpoint)
	if err != nil {
		return endpoint, errors.Wrap(ctx, err, "create endpoint")
	}

	log.Info("Starting the endpoint scheduler")

	ctx, log = logger.WithStructToCtx(ctx, "endpoint", endpoint)
	schedulerCtx := logger.ToCtx(context.Background(), log)

	endpointID := endpoint.ID
	endpoint, err = c.scheduler.Start(schedulerCtx, endpoint) // nolint: contextcheck // We use a background context since this context will continue to live in the endpoint manager
	if err != nil {
		c.removeEndpointFromStorage(ctx, endpointID)
		return endpoint, errors.Wrap(ctx, err, "start scheduler")
	}

	return endpoint, nil
}

func (c *creator) removeEndpointFromStorage(ctx context.Context, endpointID string) {
	log := logger.Get(ctx)
	log.Info("Removing endpoint from database")

	err := c.storage.RemoveEndpoint(ctx, endpointID)
	if err != nil {
		log.WithError(err).Error("Failed to remove endpoint from database")
	}
}
