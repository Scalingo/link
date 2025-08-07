package models

import "context"

// Storage engine needed for the LinK persistent memory
type Storage interface {
	GetEndpoints(ctx context.Context) (Endpoints, error) // GetEndpoints configured for this host
	AddEndpoint(ctx context.Context, endpoint Endpoint) (Endpoint, error)
	UpdateEndpoint(ctx context.Context, endpoint Endpoint) error
	RemoveEndpoint(ctx context.Context, id string) error

	GetCurrentHost(ctx context.Context) (Host, error) // Get host configuration for the current host
	SaveHost(ctx context.Context, host Host) error    // Save host modifications

	LinkEndpointWithCurrentHost(ctx context.Context, key string) error   // Link an Endpoint to the current host
	UnlinkEndpointFromCurrentHost(ctx context.Context, key string) error // Unlink an Endpoint from the current host
	GetEndpointHosts(ctx context.Context, key string) ([]string, error)  // List all hosts linked to the Endpoint

	GetEncryptedData(ctx context.Context, endpointID string, encryptedDataId string) (EncryptedData, error)
	UpsertEncryptedData(ctx context.Context, endpointID string, data EncryptedData) (EncryptedDataLink, error)
	RemoveEncryptedDataForEndpoint(ctx context.Context, endpointID string) error
	ListEncryptedDataForHost(ctx context.Context) ([]EncryptedData, error)
}
