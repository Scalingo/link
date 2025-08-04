package api

import (
	"context"
)

type Client interface {
	ListIPs(context.Context) ([]IP, error)
	GetIP(ctx context.Context, id string) (IP, error)
	AddIP(ctx context.Context, ip string, params AddIPParams) (IP, error)
	UpdateIP(ctx context.Context, id string, params UpdateIPParams) (IP, error)
	RemoveIP(ctx context.Context, id string) error
	Failover(ctx context.Context, id string) error
	Version(context.Context) (string, error)
}
