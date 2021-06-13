package session

const (
	RefreshAPIPath             = "/session/refresh"
	SignoutAPIPath             = "/signout"

	AntiCSRF_VIA_TOKEN         = "VIA_TOKEN"
	AntiCSRF_VIA_CUSTOM_HEADER = "VIA_CUSTOM_HEADER"
	AntiCSRF_NONE              = "NONE"
	
	CookieSameSite_NONE        = "none"
	CookieSameSite_LAX         = "lax"
	CookieSameSite_STRICT      = "strict"
)
