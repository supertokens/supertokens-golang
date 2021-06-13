package session

const (
	accessTokenCookieKey  = "sAccessToken"
	refreshTokenCookieKey = "sRefreshToken"

	// there are two of them because one is used by the server to check if the user is logged in and the other is checked by the frontend to see if the user is logged in.
	idRefreshTokenCookieKey = "sIdRefreshToken"
	idRefreshTokenHeaderKey = "id-refresh-token"

	antiCsrfHeaderKey = "anti-csrf"
	ridHeaderKey = "rid"

	frontTokenHeaderKey = "front-token"
)
