package errors

import (
	"reflect"
)

const (
	UnauthorizedErrorStr       = "UNAUTHORISED"
	TryRefreshTokenErrorStr    = "TRY_REFRESH_TOKEN"
	TokenTheftDetectedErrorStr = "TOKEN_THEFT_DETECTED"
)

type SessionError struct{}

// TryRefreshTokenError used for when the refresh API needs to be called
type TryRefreshTokenError struct {
	Msg  string
	Type string
}

func MakeTryRefreshTokenError(msg string) TryRefreshTokenError {
	return TryRefreshTokenError{
		Msg:  msg,
		Type: TryRefreshTokenErrorStr,
	}
}

func (err TryRefreshTokenError) Error() string {
	return err.Msg
}

// TokenTheftDetectedError used for when token theft has happened for a session
type TokenTheftDetectedError struct {
	Msg     string
	Type    string
	Payload TokenTheftDetectedErrorPayload
}

type TokenTheftDetectedErrorPayload struct {
	SessionHandle string `json:"sessionHandle"`
	UserID        string `json:"userId"`
}

func MakeTokenTheftDetectedError(sessionHandle, userID, msg string) TokenTheftDetectedError {
	return TokenTheftDetectedError{
		Msg:  msg,
		Type: TokenTheftDetectedErrorStr,
		Payload: TokenTheftDetectedErrorPayload{
			SessionHandle: sessionHandle,
			UserID:        userID,
		},
	}
}

func (err TokenTheftDetectedError) Error() string {
	return err.Msg
}

// UnauthorizedError used for when the user has been logged out
type UnauthorizedError struct {
	Msg  string
	Type string
}

func MakeUnauthorizedError(msg string) UnauthorizedError {
	return UnauthorizedError{
		Msg:  msg,
		Type: UnauthorizedErrorStr,
	}
}

func (err UnauthorizedError) Error() string {
	return err.Msg
}

// IsTokenTheftDetectedError returns true if error is a TokenTheftDetectedError
func IsTokenTheftDetectedError(err error) bool {
	return reflect.TypeOf(err) == reflect.TypeOf(TokenTheftDetectedError{})
}

// IsUnauthorizedError returns true if error is a UnauthorizedError
func IsUnauthorizedError(err error) bool {
	return reflect.TypeOf(err) == reflect.TypeOf(UnauthorizedError{})
}

// IsTryRefreshTokenError returns true if error is a TryRefreshTokenError
func IsTryRefreshTokenError(err error) bool {
	return reflect.TypeOf(err) == reflect.TypeOf(TryRefreshTokenError{})
}
