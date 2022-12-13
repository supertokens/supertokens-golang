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
	Cookie TokenTransferMethod = "cookie"
	Header TokenTransferMethod = "header"
	Any    TokenTransferMethod = "any"
)

type HandshakeInfo struct {
	rawJwtSigningPublicKeyList     []KeyInfo
	AntiCsrf                       string
	AccessTokenBlacklistingEnabled bool
	AccessTokenValidity            uint64
	RefreshTokenValidity           uint64
}

func (h *HandshakeInfo) GetJwtSigningPublicKeyList() []KeyInfo {
	result := []KeyInfo{}
	for _, key := range h.rawJwtSigningPublicKeyList {
		if key.ExpiryTime > getCurrTimeInMS() {
			result = append(result, key)
		}
	}
	return result
}

func (h *HandshakeInfo) SetJwtSigningPublicKeyList(updatedList []KeyInfo) {
	h.rawJwtSigningPublicKeyList = updatedList
}

func getCurrTimeInMS() uint64 {
	return uint64(time.Now().UnixNano() / 1000000)
}

type KeyInfo struct {
	PublicKey  string
	ExpiryTime uint64
	CreatedAt  uint64
}

type CreateOrRefreshAPIResponse struct {
	Session        SessionStruct                   `json:"session"`
	AccessToken    CreateOrRefreshAPIResponseToken `json:"accessToken"`
	RefreshToken   CreateOrRefreshAPIResponseToken `json:"refreshToken"`
	IDRefreshToken CreateOrRefreshAPIResponseToken `json:"idRefreshToken"`
	AntiCsrfToken  *string                         `json:"antiCsrfToken"`
}

type SessionStruct struct {
	Handle                string                 `json:"handle"`
	UserID                string                 `json:"userId"`
	UserDataInAccessToken map[string]interface{} `json:"userDataInJWT"`
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

type RegenerateAccessTokenResponse struct {
	Status      string                          `json:"status"`
	Session     SessionStruct                   `json:"session"`
	AccessToken CreateOrRefreshAPIResponseToken `json:"accessToken"`
}

type TypeInput struct {
	CookieSecure             *bool
	CookieSameSite           *string
	SessionExpiredStatusCode *int
	InvalidClaimStatusCode   *int
	CookieDomain             *string
	AntiCsrf                 *string
	Override                 *OverrideStruct
	ErrorHandlers            *ErrorHandlers
	Jwt                      *JWTInputConfig
	GetTokenTransferMethod   func(req *http.Request, forCreateNewSession bool, userContext supertokens.UserContext) string
}

type JWTInputConfig struct {
	Issuer                           *string
	Enable                           bool
	PropertyNameInAccessTokenPayload *string
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
	RefreshTokenPath         supertokens.NormalisedURLPath
	CookieDomain             *string
	CookieSameSite           string
	CookieSecure             bool
	SessionExpiredStatusCode int
	InvalidClaimStatusCode   int
	AntiCsrf                 string
	Override                 OverrideStruct
	ErrorHandlers            NormalisedErrorHandlers
	Jwt                      JWTNormalisedConfig
	GetTokenTransferMethod   func(req *http.Request, forCreateNewSession bool, userContext supertokens.UserContext) string
}

type JWTNormalisedConfig struct {
	Issuer                           *string
	Enable                           bool
	PropertyNameInAccessTokenPayload string
}

type VerifySessionOptions struct {
	AntiCsrfCheck                 *bool
	SessionRequired               *bool
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

type TypeSessionContainer struct {
	RevokeSession            func() error
	GetSessionData           func() (map[string]interface{}, error)
	UpdateSessionData        func(newSessionData map[string]interface{}) error
	GetUserID                func() string
	GetAccessTokenPayload    func() map[string]interface{}
	GetHandle                func() string
	GetAccessToken           func() string
	UpdateAccessTokenPayload func(newAccessTokenPayload map[string]interface{}) error // Deprecated: use MergeIntoAccessTokenPayload instead
	GetTimeCreated           func() (uint64, error)
	GetExpiry                func() (uint64, error)

	RevokeSessionWithContext            func(userContext supertokens.UserContext) error
	GetSessionDataWithContext           func(userContext supertokens.UserContext) (map[string]interface{}, error)
	UpdateSessionDataWithContext        func(newSessionData map[string]interface{}, userContext supertokens.UserContext) error
	GetUserIDWithContext                func(userContext supertokens.UserContext) string
	GetAccessTokenPayloadWithContext    func(userContext supertokens.UserContext) map[string]interface{}
	GetHandleWithContext                func(userContext supertokens.UserContext) string
	GetAccessTokenWithContext           func(userContext supertokens.UserContext) string
	UpdateAccessTokenPayloadWithContext func(newAccessTokenPayload map[string]interface{}, userContext supertokens.UserContext) error // Deprecated: use MergeIntoAccessTokenPayloadWithContext instead
	GetTimeCreatedWithContext           func(userContext supertokens.UserContext) (uint64, error)
	GetExpiryWithContext                func(userContext supertokens.UserContext) (uint64, error)

	MergeIntoAccessTokenPayloadWithContext func(accessTokenPayloadUpdate map[string]interface{}, userContext supertokens.UserContext) error

	AssertClaimsWithContext     func(claimValidators []claims.SessionClaimValidator, userContext supertokens.UserContext) error
	FetchAndSetClaimWithContext func(claim *claims.TypeSessionClaim, userContext supertokens.UserContext) error
	SetClaimValueWithContext    func(claim *claims.TypeSessionClaim, value interface{}, userContext supertokens.UserContext) error
	GetClaimValueWithContext    func(claim *claims.TypeSessionClaim, userContext supertokens.UserContext) interface{}
	RemoveClaimWithContext      func(claim *claims.TypeSessionClaim, userContext supertokens.UserContext) error

	MergeIntoAccessTokenPayload func(accessTokenPayloadUpdate map[string]interface{}) error

	AssertClaims     func(claimValidators []claims.SessionClaimValidator) error
	FetchAndSetClaim func(claim *claims.TypeSessionClaim) error
	SetClaimValue    func(claim *claims.TypeSessionClaim, value interface{}) error
	GetClaimValue    func(claim *claims.TypeSessionClaim) interface{}
	RemoveClaim      func(claim *claims.TypeSessionClaim) error
}

type SessionContainer = *TypeSessionContainer

type SessionInformation struct {
	SessionHandle      string
	UserId             string
	SessionData        map[string]interface{}
	Expiry             uint64
	AccessTokenPayload map[string]interface{}
	TimeCreated        uint64
}

const SessionContext int = iota
