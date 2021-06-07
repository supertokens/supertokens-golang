package schema

import (
	"net/http"

	"github.com/supertokens/supertokens-golang/supertokens"
)

type AntiCsrfStr string
type CookieSameSiteStr string

const (
	ViaToken        AntiCsrfStr = "VIA_TOKEN"
	ViaCustomHeader AntiCsrfStr = "VIA_CUSTOM_HEADER"
	NoneAntiCsrf    AntiCsrfStr = "NONE"

	Strict     CookieSameSiteStr = "strict"
	Lax        CookieSameSiteStr = "lax"
	NoneCookie CookieSameSiteStr = "NONE"
)

type RecipeImplementation struct {
	Querier       supertokens.Querier
	Config        TypeNormalisedInput
	HandshakeInfo HandshakeInfo
}

type APIImplementation struct{}

type HandshakeInfo struct {
	JWTSigningPublicKey            string
	AntiCsrf                       AntiCsrfStr
	AccessTokenBlacklistingEnabled bool
	JWTSigningPublicKeyExpiryTime  int
	AccessTokenValidity            int
	RefreshTokenValidity           int
}

type CreateOrRefreshAPIResponse struct {
	Session struct {
		Handle        string
		UserID        string
		UserDataInJWT interface{}
	}
	AccessToken struct {
		Token       string
		Expiry      int
		CreatedTime int
	}
	RefreshToken struct {
		Token       string
		Expiry      int
		CreatedTime int
	}
	IDRefreshToken struct {
		Token       string
		Expiry      int
		CreatedTime int
	}
	AntiCsrfToken *string
}

type TypeInput struct {
	CookieSecure             bool
	CookieSameSite           CookieSameSiteStr
	SessionExpiredStatusCode int
	CookieDomain             string
	AntiCsrf                 AntiCsrfStr
	Override                 struct {
		functions func(originalImplementation RecipeImplementation) RecipeInterface
		apis      func(originalImplementation APIImplementation) APIInterface
	}
}

type TypeNormalisedInput struct {
	RefreshTokenPath         supertokens.NormalisedURLPath
	CookieDomain             *string
	CookieSameSite           CookieSameSiteStr
	CookieSecure             bool
	SessionExpiredStatusCode int
	AntiCsrf                 AntiCsrfStr
	Override                 struct {
		Functions func(originalImplementation RecipeImplementation) RecipeInterface
		Apis      func(originalImplementation APIImplementation) APIInterface
	}
}

type SessionContainerInterface interface {
	RevokeSession()
	GetSessionData() interface{}
	UpdateSessionData(newSessionData interface{}) interface{}
	GetUserId() string
	GetJWTPayload() interface{}
	GetHandle() string
	GetAccessToken() string
	UpdateJWTPayload(newJWTPayload interface{})
}

type RecipeInterface interface {
	CreateNewSession(
		res http.ResponseWriter,
		userID string,
		jwtPayload interface{},
		sessionData interface{},
	) SessionContainerInterface
	GetSession(
		req *http.Request,
		res http.ResponseWriter,
		options VerifySessionOptions,
	) *SessionContainerInterface
	RefreshSession(req *http.Request, res http.ResponseWriter) SessionContainerInterface
	RevokeAllSessionsForUser(userID string) []string
	GetAllSessionHandlesForUser(userID string) []string
	RevokeSession(sessionHandle string) bool
	RevokeMultipleSessions(sessionHandles []string) []string
	GetSessionData(sessionHandle string) interface{}
	UpdateSessionData(sessionHandle string, newSessionData interface{})
	GetJWTPayload(sessionHandle string) interface{}
	UpdateJWTPayload(sessionHandle string, newJWTPayload interface{})
}

type VerifySessionOptions struct {
	AntiCsrfCheck   bool
	SessionRequired bool
}

type APIOptions struct {
	RecipeImplementation RecipeInterface
	Config               TypeNormalisedInput
	RecipeID             string
	Req                  *http.Request
	Res                  http.ResponseWriter
}

type APIInterface interface {
	RefreshPOST(options APIOptions)
	SignOutPOST(options APIOptions) map[string]string
	VerifySession(verifySessionOptions *VerifySessionOptions, options APIOptions)
}
