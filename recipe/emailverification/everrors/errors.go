package everrors

const (
	UnknownUserIdErrorStr = "UNKNOWN_USER_ID"
)

type UnknownUserIdError struct {
	Msg string
}

func (err UnknownUserIdError) Error() string {
	return err.Msg
}
