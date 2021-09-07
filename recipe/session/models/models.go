package models

import (
	"net/http"

	"github.com/supertokens/supertokens-golang/supertokens"
)

type HandshakeInfo struct {
	JWTSigningPublicKey            string
	AntiCsrf                       string
	AccessTokenBlacklistingEnabled bool
	JWTSigningPublicKeyExpiryTime  uint64
	AccessTokenValidity            uint64
	RefreshTokenValidity           uint64
	SigningKeyLastUpdated          uint64
}

type CreateOrRefreshAPIResponse struct {
	Session        SessionStruct                   `json:"session"`
	AccessToken    CreateOrRefreshAPIResponseToken `json:"accessToken"`
	RefreshToken   CreateOrRefreshAPIResponseToken `json:"refreshToken"`
	IDRefreshToken CreateOrRefreshAPIResponseToken `json:"idRefreshToken"`
	AntiCsrfToken  *string                         `json:"antiCsrfToken"`
}

type SessionStruct struct {
	Handle        string      `json:"handle"`
	UserID        string      `json:"userId"`
	UserDataInJWT interface{} `json:"userDataInJWT"`
}

type CreateOrRefreshAPIResponseToken struct {
	Token       string `json:"token"`
	Expiry      uint64 `json:"expiry"`
	CreatedTime uint64 `json:"createdTime"`
}

type GetSessionResponse struct {
	Session     SessionStruct                   `json:"session"`
	AccessToken CreateOrRefreshAPIResponseToken `json:"accessToken"`
}

type TypeInput struct {
	CookieSecure             *bool
	CookieSameSite           *string
	SessionExpiredStatusCode *int
	CookieDomain             *string
	AntiCsrf                 *string
	Override                 *struct {
		Functions func(originalImplementation RecipeInterface) RecipeInterface
		APIs      func(originalImplementation APIInterface) APIInterface
	}
	ErrorHandlers *ErrorHandlers
}

type ErrorHandlers struct {
	OnUnauthorised       func(message string, req *http.Request, res http.ResponseWriter) error
	OnTokenTheftDetected func(sessionHandle string, userID string, req *http.Request, res http.ResponseWriter) error
}

type TypeNormalisedInput struct {
	RefreshTokenPath         supertokens.NormalisedURLPath
	CookieDomain             *string
	CookieSameSite           string
	CookieSecure             bool
	SessionExpiredStatusCode int
	AntiCsrf                 string
	Override                 struct {
		Functions func(originalImplementation RecipeInterface) RecipeInterface
		APIs      func(originalImplementation APIInterface) APIInterface
	}
	ErrorHandlers NormalisedErrorHandlers
}

type VerifySessionOptions struct {
	AntiCsrfCheck   *bool
	SessionRequired *bool
}

type APIOptions struct {
	RecipeImplementation RecipeInterface
	Config               TypeNormalisedInput
	RecipeID             string
	Req                  *http.Request
	Res                  http.ResponseWriter
	OtherHandler         http.HandlerFunc
}

type SessionRecipe struct {
	RecipeModule supertokens.RecipeModule
	Config       TypeNormalisedInput
	RecipeImpl   RecipeInterface
	APIImpl      APIInterface
}

type NormalisedErrorHandlers struct {
	OnUnauthorised       func(message string, req *http.Request, res http.ResponseWriter) error
	OnTryRefreshToken    func(message string, req *http.Request, res http.ResponseWriter) error
	OnTokenTheftDetected func(sessionHandle string, userID string, req *http.Request, res http.ResponseWriter) error
}

type SessionContainer struct {
	RevokeSession     func() error
	GetSessionData    func() (interface{}, error)
	UpdateSessionData func(newSessionData interface{}) error
	GetUserID         func() string
	GetJWTPayload     func() interface{}
	GetHandle         func() string
	GetAccessToken    func() string
	UpdateJWTPayload  func(newJWTPayload interface{}) error
	GetTimeCreated    func() (uint64, error)
	GetExpiry         func() (uint64, error)
}

type SessionInformation struct {
	SessionHandle string
	UserId        string
	SessionData   interface{}
	Expiry        uint64
	JwtPayload    interface{}
	TimeCreated   uint64
}

const SessionContext int = iota
