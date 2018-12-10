package api

import (
	"context"

	"github.com/Scalingo/link/models"
)

type Client interface {
	ListIPs(ctx context.Context) ([]IP, error)
	GetIP(ctx context.Context, id string) (IP, error)
	AddIP(ctx context.Context, ip string, checks ...models.Healthcheck) (IP, error)
	RemoveIP(ctx context.Context, id string) error
	TryGetLock(ctx context.Context, id string) error
	Version(ctx context.Context) string
}
