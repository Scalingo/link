package plugin

import (
	"context"

	"github.com/Scalingo/go-utils/errors/v2"
	"github.com/Scalingo/link/v2/models"
)

var ErrPluginNotFound = errors.New(context.Background(), "plugin not found")

type Factory interface {
	Validate(ctx context.Context, endpoint models.Endpoint) error
	Create(ctx context.Context, endpoint models.Endpoint) (Plugin, error)
}

type Registry interface {
	Register(ctx context.Context, pluginName string, factory Factory)
	Create(ctx context.Context, endpoint models.Endpoint) (Plugin, error)
	Validate(ctx context.Context, endpoint models.Endpoint) error
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
