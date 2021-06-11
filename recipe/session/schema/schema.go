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
	CookieSecure             *bool
	CookieSameSite           *CookieSameSiteStr
	SessionExpiredStatusCode *int
	CookieDomain             *string
	AntiCsrf                 *AntiCsrfStr
	Override                 *struct {
		Functions func(originalImplementation RecipeImplementation) RecipeImplementation
		APIs      func(originalImplementation APIImplementation) APIImplementation
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
		Functions func(originalImplementation RecipeImplementation) RecipeImplementation
		APIs      func(originalImplementation APIImplementation) APIImplementation
	}
}

type SessionContainerInterface struct {
	RevokeSession     func()
	GetSessionData    func() interface{}
	UpdateSessionData func(newSessionData interface{}) interface{}
	GetUserId         func() string
	GetJWTPayload     func() interface{}
	GetHandle         func() string
	GetAccessToken    func() string
	UpdateJWTPayload  func(newJWTPayload interface{})
}

type RecipeImplementation struct {
	CreateNewSession            func(res http.ResponseWriter, userID string, jwtPayload interface{}, sessionData interface{}) SessionContainerInterface
	GetSession                  func(req *http.Request, res http.ResponseWriter, options *VerifySessionOptions) *SessionContainerInterface
	RefreshSession              func(req *http.Request, res http.ResponseWriter) SessionContainerInterface
	RevokeAllSessionsForUser    func(userID string) []string
	GetAllSessionHandlesForUser func(userID string) []string
	RevokeSession               func(sessionHandle string) bool
	RevokeMultipleSessions      func(sessionHandles []string) []string
	GetSessionData              func(sessionHandle string) interface{}
	UpdateSessionData           func(sessionHandle string, newSessionData interface{})
	GetJWTPayload               func(sessionHandle string) interface{}
	UpdateJWTPayload            func(sessionHandle string, newJWTPayload interface{})
}

type VerifySessionOptions struct {
	AntiCsrfCheck   *bool
	SessionRequired *bool
}

type APIOptions struct {
	RecipeImplementation RecipeImplementation
	Config               TypeNormalisedInput
	RecipeID             string
	Req                  *http.Request
	Res                  http.ResponseWriter
}

type APIImplementation struct {
	RefreshPOST   func(options APIOptions)
	SignOutPOST   func(options APIOptions) map[string]string
	VerifySession func(verifySessionOptions *VerifySessionOptions, options APIOptions)
}
