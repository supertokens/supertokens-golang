package session

const (
	UnauthorizedError       = "UNAUTHORISED"
	TryRefreshTokenError    = "TRY_REFRESH_TOKEN"
	TokenTheftDetectedError = "TOKEN_THEFT_DETECTED"
)

type SessionError struct{}
