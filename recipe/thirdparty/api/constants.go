package api

// If Third Party login is used with one of the following development keys, then the dev authorization url and the redirect url will be used.
// When adding or changing client id's they should be in the following order: Google and Facebook
var DevOauthClientIds = map[string]bool{
	"1060725074195-kmeum4crr01uirfl2op9kd5acmi9jutn.apps.googleusercontent.com": true, // google
	"467101b197249757c71f": true, // github
}

const (
	DevOauthAuthorisationUrl = "https://supertokens.io/dev/oauth/redirect-to-provider"
	DevOauthRedirectUrl      = "https://supertokens.io/dev/oauth/redirect-to-app"
)
