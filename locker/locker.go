package locker

import "context"

// Locker is a structure that let distribute locks accross the entire networks
type Locker interface {
	Refresh(context.Context) error          // Refresh must be called to refresh our TTL or to try to get the lock
	Unlock(context.Context) error           // Unlock remove the lock and mark the key accessible for re-election
	IsMaster(context.Context) (bool, error) // IsMaster return true if the lock is allocated to our node
	Stop(context.Context) error             // Stop the locker (cleanup)
}
