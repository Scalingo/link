package models

import "context"

type Storage interface {
	GetIPs(context.Context) ([]IP, error)
	AddIP(context.Context, IP) (IP, error)
	UpdateIP(ctx context.Context, ip IP) error
	RemoveIP(context.Context, string) error

	GetHost(context.Context) (Host, error)
	SaveHost(context.Context, Host) error
}
