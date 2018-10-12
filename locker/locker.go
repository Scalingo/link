package locker

import "context"

type Locker interface {
	Refresh(context.Context) error
	Unlock(context.Context) error
	IsMaster(context.Context) (bool, error)
	Stop(context.Context) error
}
