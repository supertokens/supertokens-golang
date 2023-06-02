/* Copyright (c) 2021, VRAI Labs and/or its affiliates. All rights reserved.
 *
 * This software is licensed under the Apache License, Version 2.0 (the
 * "License") as published by the Apache Software Foundation.
 *
 * You may not use this file except in compliance with the License. You may
 * obtain a copy of the License at http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
 * WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
 * License for the specific language governing permissions and limitations
 * under the License.
 */

package sessmodels

import (
	"github.com/MicahParks/keyfunc"
	"net/http"
	"time"

	"github.com/supertokens/supertokens-golang/recipe/openid/openidmodels"
	"github.com/supertokens/supertokens-golang/recipe/session/claims"
	"github.com/supertokens/supertokens-golang/supertokens"
)

type TokenType string

const (
	AccessToken  TokenType = "access"
	RefreshToken TokenType = "refresh"
)

// When adding a new token transfer method, it's also necessary to update the related constant (availableTokenTransferMethods) in `session/constants.go`
type TokenTransferMethod string

const (
	CookieTransferMethod TokenTransferMethod = "cookie"
	HeaderTransferMethod TokenTransferMethod = "header"
	AnyTransferMethod    TokenTransferMethod = "any"
)

type GetJWKSResult struct {
	JWKS        *keyfunc.JWKS
	Error       error
	LastFetched int64
}

type GetJWKSFunctionObject struct {
	Fn   GetJWKSFunction
	Path string
}

type GetJWKSFunction = func(string) GetJWKSResult

func getCurrTimeInMS() uint64 {
	return uint64(time.Now().UnixNano() / 1000000)
}

type KeyInfo struct {
	PublicKey  string
	ExpiryTime uint64
	CreatedAt  uint64
}

type CreateOrRefreshAPIResponse struct {
	Session       SessionStruct                   `json:"session"`
	AccessToken   CreateOrRefreshAPIResponseToken `json:"accessToken"`
	RefreshToken  CreateOrRefreshAPIResponseToken `json:"refreshToken"`
	AntiCsrfToken *string                         `json:"antiCsrfToken"`
}

type SessionStruct struct {
	Handle                string                 `json:"handle"`
	UserID                string                 `json:"userId"`
	UserDataInAccessToken map[string]interface{} `json:"userDataInJWT"`
	ExpiryTime            uint64                 `json:"expiryTime"`
}

type CreateOrRefreshAPIResponseToken struct {
	Token       string `json:"token"`
	Expiry      uint64 `json:"expiry"`
	CreatedTime uint64 `json:"createdTime"`
}

type GetSessionResponse struct {
	Status      string
	Session     SessionStruct                   `json:"session"`
	AccessToken CreateOrRefreshAPIResponseToken `json:"accessToken"`
}

type RegenerateAccessTokenResponse struct {
	Status      string                          `json:"status"`
	Session     SessionStruct                   `json:"session"`
	AccessToken CreateOrRefreshAPIResponseToken `json:"accessToken"`
}

type TypeInput struct {
	CookieSecure                                 *bool
	CookieSameSite                               *string
	SessionExpiredStatusCode                     *int
	InvalidClaimStatusCode                       *int
	CookieDomain                                 *string
	AntiCsrf                                     *string
	Override                                     *OverrideStruct
	ErrorHandlers                                *ErrorHandlers
	GetTokenTransferMethod                       func(req *http.Request, forCreateNewSession bool, userContext supertokens.UserContext) TokenTransferMethod
	ExposeAccessTokenToFrontendInCookieBasedAuth bool
	UseDynamicAccessTokenSigningKey              *bool
}

type OverrideStruct struct {
	Functions     func(originalImplementation RecipeInterface) RecipeInterface
	APIs          func(originalImplementation APIInterface) APIInterface
	OpenIdFeature *openidmodels.OverrideStruct
}

type ErrorHandlers struct {
	OnUnauthorised       func(message string, req *http.Request, res http.ResponseWriter) error
	OnTokenTheftDetected func(sessionHandle string, userID string, req *http.Request, res http.ResponseWriter) error
	OnInvalidClaim       func(validationErrors []claims.ClaimValidationError, req *http.Request, res http.ResponseWriter) error
}

type TypeNormalisedInput struct {
	RefreshTokenPath                             supertokens.NormalisedURLPath
	CookieDomain                                 *string
	CookieSameSite                               string
	CookieSecure                                 bool
	SessionExpiredStatusCode                     int
	InvalidClaimStatusCode                       int
	AntiCsrf                                     string
	Override                                     OverrideStruct
	ErrorHandlers                                NormalisedErrorHandlers
	GetTokenTransferMethod                       func(req *http.Request, forCreateNewSession bool, userContext supertokens.UserContext) TokenTransferMethod
	ExposeAccessTokenToFrontendInCookieBasedAuth bool
	UseDynamicAccessTokenSigningKey              bool
}

type JWTNormalisedConfig struct {
	Issuer                           *string
	Enable                           bool
	PropertyNameInAccessTokenPayload string
}

type VerifySessionOptions struct {
	AntiCsrfCheck                 *bool
	SessionRequired               *bool
	CheckDatabase                 *bool
	OverrideGlobalClaimValidators func(globalClaimValidators []claims.SessionClaimValidator, sessionContainer SessionContainer, userContext supertokens.UserContext) ([]claims.SessionClaimValidator, error)
}

type APIOptions struct {
	RecipeImplementation RecipeInterface
	Config               TypeNormalisedInput
	RecipeID             string
	Req                  *http.Request
	Res                  http.ResponseWriter
	OtherHandler         http.HandlerFunc

	ClaimValidatorsAddedByOtherRecipes []claims.SessionClaimValidator
}

type NormalisedErrorHandlers struct {
	OnUnauthorised       func(message string, req *http.Request, res http.ResponseWriter) error
	OnTryRefreshToken    func(message string, req *http.Request, res http.ResponseWriter) error
	OnTokenTheftDetected func(sessionHandle string, userID string, req *http.Request, res http.ResponseWriter) error
	OnInvalidClaim       func(validationErrors []claims.ClaimValidationError, req *http.Request, res http.ResponseWriter) error
}

type SessionTokens struct {
	AccessToken                   string
	RefreshToken                  *string
	AntiCsrfToken                 *string
	FrontToken                    string
	AccessAndFrontendTokenUpdated bool
}

type RequestResponseInfo struct {
	Res                 http.ResponseWriter
	Req                 *http.Request
	TokenTransferMethod TokenTransferMethod
}

type TypeSessionContainer struct {
	RevokeSession                  func() error
	GetSessionDataInDatabase       func() (map[string]interface{}, error)
	UpdateSessionDataInDatabase    func(newSessionData map[string]interface{}) error
	GetUserID                      func() string
	GetAccessTokenPayload          func() map[string]interface{}
	GetHandle                      func() string
	GetAllSessionTokensDangerously func() SessionTokens
	GetAccessToken                 func() string
	GetTimeCreated                 func() (uint64, error)
	GetExpiry                      func() (uint64, error)

	RevokeSessionWithContext               func(userContext supertokens.UserContext) error
	GetSessionDataInDatabaseWithContext    func(userContext supertokens.UserContext) (map[string]interface{}, error)
	UpdateSessionDataInDatabaseWithContext func(newSessionData map[string]interface{}, userContext supertokens.UserContext) error
	GetUserIDWithContext                   func(userContext supertokens.UserContext) string
	GetAccessTokenPayloadWithContext       func(userContext supertokens.UserContext) map[string]interface{}
	GetHandleWithContext                   func(userContext supertokens.UserContext) string
	GetAccessTokenWithContext              func(userContext supertokens.UserContext) string
	GetTimeCreatedWithContext              func(userContext supertokens.UserContext) (uint64, error)
	GetExpiryWithContext                   func(userContext supertokens.UserContext) (uint64, error)

	MergeIntoAccessTokenPayloadWithContext func(accessTokenPayloadUpdate map[string]interface{}, userContext supertokens.UserContext) error

	AssertClaimsWithContext     func(claimValidators []claims.SessionClaimValidator, userContext supertokens.UserContext) error
	FetchAndSetClaimWithContext func(claim *claims.TypeSessionClaim, userContext supertokens.UserContext) error
	SetClaimValueWithContext    func(claim *claims.TypeSessionClaim, value interface{}, userContext supertokens.UserContext) error
	GetClaimValueWithContext    func(claim *claims.TypeSessionClaim, userContext supertokens.UserContext) interface{}
	RemoveClaimWithContext      func(claim *claims.TypeSessionClaim, userContext supertokens.UserContext) error

	MergeIntoAccessTokenPayload func(accessTokenPayloadUpdate map[string]interface{}) error

	AssertClaims            func(claimValidators []claims.SessionClaimValidator) error
	FetchAndSetClaim        func(claim *claims.TypeSessionClaim) error
	SetClaimValue           func(claim *claims.TypeSessionClaim, value interface{}) error
	GetClaimValue           func(claim *claims.TypeSessionClaim) interface{}
	RemoveClaim             func(claim *claims.TypeSessionClaim) error
	AttachToRequestResponse func(info RequestResponseInfo) error
}

type SessionContainer = *TypeSessionContainer

type SessionInformation struct {
	SessionHandle                    string
	UserId                           string
	SessionDataInDatabase            map[string]interface{}
	Expiry                           uint64
	CustomClaimsInAccessTokenPayload map[string]interface{}
	TimeCreated                      uint64
}

type ParsedJWTInfo struct {
	RawTokenString string
	RawPayload     string
	Header         string
	Payload        map[string]interface{}
	Signature      string
	Version        int
	KID            *string
}

const SessionContext int = iota
