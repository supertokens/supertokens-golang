package errors

type ForbiddenAccessError struct {
	Msg string
}

func (err ForbiddenAccessError) Error() string {
	return err.Msg
}
