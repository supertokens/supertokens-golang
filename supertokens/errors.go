package supertokens

const BadInputErrorStr = "BAD_INPUT_ERROR"

// BadInputError used for non specific exceptions
type BadInputError struct {
	Msg string
}

func (err BadInputError) Error() string {
	return err.Msg
}
