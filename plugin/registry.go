package plugin

import (
	"context"
	"encoding/json"

	"github.com/Scalingo/go-utils/errors/v2"
	"github.com/Scalingo/link/v3/models"
)

var ErrPluginNotFound = errors.New(context.Background(), "plugin not found")

// Every plugin must implement a Factory that will be used to create the plugin instance.
// This Factory is also responsible for handling validation and mutation of the plugin config.

// Factory is the strict minimum interface used by
type Factory interface {
	Validate(ctx context.Context, endpoint models.Endpoint) error
	Create(ctx context.Context, endpoint models.Endpoint) (Plugin, error)
}

// MutableFactory is an interface that plugins can used if they need to change
// the plugin configuration before storing it.
type MutableFactory interface {
	Validate(ctx context.Context, endpoint models.Endpoint) error
	Create(ctx context.Context, endpoint models.Endpoint) (Plugin, error)
	// Mutate will let a plugin customize what's being stored in the database.
	// It takes the endpoint as input and returns a json.RawMessage representing the plugin config as it will be stored in the database.
	// This is useful for plugins that need to encrypt some fields before storing them.
	Mutate(ctx context.Context, endpoint models.Endpoint) (json.RawMessage, error)
}

type Registry interface {
	Register(ctx context.Context, pluginName string, factory Factory)
	Create(ctx context.Context, endpoint models.Endpoint) (Plugin, error)
	Validate(ctx context.Context, endpoint models.Endpoint) error
	Mutate(ctx context.Context, endpoint models.Endpoint) (json.RawMessage, error)
}

type registry struct {
	plugins map[string]Factory
}

func NewRegistry() Registry {
	return &registry{
		plugins: make(map[string]Factory),
	}
}

func (r *registry) Register(_ context.Context, pluginName string, factory Factory) {
	r.plugins[pluginName] = factory
}

func (r *registry) Create(ctx context.Context, endpoint models.Endpoint) (Plugin, error) {
	pluginName := endpoint.Plugin
	if pluginName == "" {
		pluginName = "arp"
	}

	factory, ok := r.plugins[pluginName]
	if !ok {
		return nil, ErrPluginNotFound
	}
	p, err := factory.Create(ctx, endpoint)
	if err != nil {
		return nil, errors.Wrap(ctx, err, "create plugin")
	}

	return p, nil
}

func (r *registry) Validate(ctx context.Context, endpoint models.Endpoint) error {
	factory, ok := r.plugins[endpoint.Plugin]
	if !ok {
		return ErrPluginNotFound
	}
	err := factory.Validate(ctx, endpoint)
	if err != nil {
		return errors.Wrap(ctx, err, "validate plugin")
	}
	return nil
}

func (r *registry) Mutate(ctx context.Context, endpoint models.Endpoint) (json.RawMessage, error) {
	factory, ok := r.plugins[endpoint.Plugin]
	if !ok {
		return nil, ErrPluginNotFound
	}
	mutableFactory, ok := factory.(MutableFactory)
	if !ok {
		return endpoint.PluginConfig, nil
	}

	res, err := mutableFactory.Mutate(ctx, endpoint)
	if err != nil {
		return nil, errors.Wrap(ctx, err, "mutate plugin")
	}
	return res, nil
}
