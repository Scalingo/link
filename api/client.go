package api

import (
	"context"
)

type Client interface {
	ListEndpoints(ctx context.Context) ([]Endpoint, error)
	GetEndpoint(ctx context.Context, id string) (Endpoint, error)
	GetEndpointHosts(ctx context.Context, id string) ([]Host, error)
	AddEndpoint(ctx context.Context, params AddEndpointParams) (Endpoint, error)
	UpdateEndpoint(ctx context.Context, id string, params UpdateEndpointParams) (Endpoint, error)
	RemoveEndpoint(ctx context.Context, id string) error
	Failover(ctx context.Context, id string) error
	RotateEncryptionKey(ctx context.Context) error
	Version(ctx context.Context) (string, error)
}
