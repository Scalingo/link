package models

import "context"

// Storage engine needed for the LinK persistent memory
type Storage interface {
	GetEndpoints(context.Context) (Endpoints, error) // GetEndpoints configured for this host
	AddEndpoint(context.Context, Endpoint) (Endpoint, error)
	UpdateEndpoint(ctx context.Context, ip Endpoint) error
	RemoveEndpoint(context.Context, string) error

	GetCurrentHost(context.Context) (Host, error) // Get host configuration for the current host
	SaveHost(context.Context, Host) error         // Save host modifications

	LinkEndpointWithCurrentHost(context.Context, Endpoint) error   // Link an Endpoint to the current host
	UnlinkEndpointFromCurrentHost(context.Context, Endpoint) error // Unlink an Endpoint from the current host
	GetEndpointHosts(context.Context, Endpoint) ([]string, error)  // List all hosts linked to the Endpoint
}
