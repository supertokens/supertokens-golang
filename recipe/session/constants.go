package session

const (
	refreshAPIPath = "/session/refresh"
	signoutAPIPath = "/signout"

	antiCSRF_VIA_TOKEN         = "VIA_TOKEN"
	antiCSRF_VIA_CUSTOM_HEADER = "VIA_CUSTOM_HEADER"
	antiCSRF_NONE              = "NONE"

	cookieSameSite_NONE   = "none"
	cookieSameSite_LAX    = "lax"
	cookieSameSite_STRICT = "strict"
)
