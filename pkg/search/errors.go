package search

// ErrNotFound represents a fail to identify a entity by its name or id
type ErrNotFound struct {
	EntityName string
	Reference  string
}

func (e ErrNotFound) Error() string {
	return "No " + e.EntityName + " with id or name containing '" +
		e.Reference + "' was found"
}
