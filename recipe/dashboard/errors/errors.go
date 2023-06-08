package errors

type OperationNotAllowedError struct {
	Msg string
}

func (err OperationNotAllowedError) Error() string {
	if err.Msg == "" {
		return "You are not permitted to perform this operation"
	}

	return err.Msg
}
