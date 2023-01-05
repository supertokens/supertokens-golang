package tperrors

type ClientTypeNotFoundError struct {
	Msg string
}

func (e ClientTypeNotFoundError) Error() string {
	return e.Msg
}
