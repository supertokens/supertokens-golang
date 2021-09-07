package errors

type FieldError struct {
	Msg     string
	Payload []ErrorPayload
}

type ErrorPayload struct {
	ID    string
	Error string
}

func (err FieldError) Error() string {
	return err.Msg
}
