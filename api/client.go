package api

import (
	"context"
)

type Client interface {
	ListEndpoints(ctx context.Context) ([]Endpoint, error)
	GetEndpoint(ctx context.Context, id string) (Endpoint, error)
	AddEndpoint(ctx context.Context, params AddEndpointParams) (Endpoint, error)
	UpdateEndpoint(ctx context.Context, id string, params UpdateEndpointParams) (Endpoint, error)
	RemoveEndpoint(ctx context.Context, id string) error
	Failover(ctx context.Context, id string) error
	Version(ctx context.Context) (string, error)
}
