package models

import "context"

// Storage engine needed for the LinK persistent memory
type Storage interface {
	GetIPs(context.Context) ([]IP, error) // GetIPs configured for this host
	AddIP(context.Context, IP) (IP, error)
	UpdateIP(ctx context.Context, ip IP) error
	RemoveIP(context.Context, string) error

	GetCurrentHost(context.Context) (Host, error) // Get host configuration for the current host
	SaveHost(context.Context, Host) error         // Save host modifications

	LinkIPWithCurrentHost(context.Context, IP) error   // Link an IP to the current host
	UnlinkIPFromCurrentHost(context.Context, IP) error // Unlink an IP from the current host
	GetIPHosts(context.Context, IP) ([]string, error)  // List all hosts linked to the IP
}
