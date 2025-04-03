package locker

import "context"

// Locker is a structure that let distribute locks across the entire networks
type Locker interface {
	Refresh(ctx context.Context) error          // Refresh must be called to refresh our TTL or to try to get the lock
	Unlock(ctx context.Context) error           // Unlock remove the lock and mark the lock accessible for re-election
	IsMaster(ctx context.Context) (bool, error) // IsMaster return true if the lock is allocated to our node
	Stop(ctx context.Context) error             // Stop the locker (cleanup)
}
