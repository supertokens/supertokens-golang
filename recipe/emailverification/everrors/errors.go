package everrors

const (
	UnknownUserIdErrorStr = "UNKNOWN_USER_ID"
)

// TryRefreshTokenError used for when the refresh API needs to be called
type UnknownUserIdError struct {
	Msg string
}

func (err UnknownUserIdError) Error() string {
	return err.Msg
}
