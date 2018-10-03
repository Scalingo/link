package api

import (
	"context"

	"github.com/Scalingo/link/models"
)

type Client interface {
	ListIPs(ctx context.Context) ([]IP, error)
	AddIP(ctx context.Context, ip string, checks ...models.Healthcheck) (IP, error)
	RemoveIP(ctx context.Context, id string) error
}
