package api

type ErrNotFound struct {
	message string
}

func (e ErrNotFound) Error() string {
	return e.message
}
