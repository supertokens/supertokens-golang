package errors

type FieldError struct {
	Msg     string
	Payload []ErrorPayload
}

type ErrorPayload struct {
	ID    string `json:"id"`
	Error string `json:"error"`
}

func (err FieldError) Error() string {
	return err.Msg
}
