package errors

const BadInputErrorStr = "BAD_INPUT_ERROR"

// BadInputError used for non specific exceptions
type BadInputError struct {
	Msg         string
	ActualError error
}

func MakeBadInputError(msg string) BadInputError {
	// TODO
	return BadInputError{Msg: msg}
}

func (err BadInputError) Error() string {
	return err.Msg
}
