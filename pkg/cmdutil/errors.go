package cmdutil

type FlagError struct {
	err error
}

func (fe *FlagError) Error() string {
	return fe.err.Error()
}

func (fe *FlagError) Unwrap() error {
	return fe.err
}

func FlagErrorWrap(err error) *FlagError {
	return &FlagError{err: err}
}
