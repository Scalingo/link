package plugin

import (
	"context"
)

// Plugin is an interface used by endpoint plugins to manage the actions needed to activate and deactivate an endpoint.
// A plugin is responsible for handling a single Endpoint.
type Plugin interface {
	// Activate is called when the endpoint needs to be activated on the current host.
	// This method is called when the state machine transition to the ACTIVATED state.
	Activate(ctx context.Context) error

	// TODO (leo): DeActivate
	// Disable is called when the endpoint needs to be disabled on the current host.
	// This method is called when the state machine transition from the ACTIVATED state to any other state.
	Disable(ctx context.Context) error

	// Ensure is called at regular interval when the endpoint is in the ACTIVATED state.
	Ensure(ctx context.Context) error

	// LockKey returns a string representing the key that represent the Endpoint.
	// This is the key that will be used for the primary election, all endpoints with the same key will be part of the same election.
	// Note: There's no prefix per plugin, if they same key is used by multiple plugins, they will all be part of the same election.
	ElectionKey(ctx context.Context) string
}
