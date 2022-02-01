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
	"github.com/supertokens/supertokens-golang/supertokens"
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
	Status             string `json:"status"`
	GetSessionResponse GetSessionResponse
}

type TypeInput struct {
	CookieSecure             *bool
	CookieSameSite           *string
	SessionExpiredStatusCode *int
	CookieDomain             *string
	AntiCsrf                 *string
	Override                 *OverrideStruct
	ErrorHandlers            *ErrorHandlers
	Jwt                      *JWTInputConfig
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
}

type TypeNormalisedInput struct {
	RefreshTokenPath         supertokens.NormalisedURLPath
	CookieDomain             *string
	CookieSameSite           string
	CookieSecure             bool
	SessionExpiredStatusCode int
	AntiCsrf                 string
	Override                 OverrideStruct
	ErrorHandlers            NormalisedErrorHandlers
	Jwt                      JWTNormalisedConfig
}

type JWTNormalisedConfig struct {
	Issuer                           *string
	Enable                           bool
	PropertyNameInAccessTokenPayload string
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

type NormalisedErrorHandlers struct {
	OnUnauthorised       func(message string, req *http.Request, res http.ResponseWriter) error
	OnTryRefreshToken    func(message string, req *http.Request, res http.ResponseWriter) error
	OnTokenTheftDetected func(sessionHandle string, userID string, req *http.Request, res http.ResponseWriter) error
}

type SessionContainer struct {
	RevokeSession            func(userContext supertokens.UserContext) error
	GetSessionData           func(userContext supertokens.UserContext) (map[string]interface{}, error)
	UpdateSessionData        func(newSessionData map[string]interface{}, userContext supertokens.UserContext) error
	GetUserID                func() string
	GetAccessTokenPayload    func() map[string]interface{}
	GetHandle                func() string
	GetAccessToken           func() string
	UpdateAccessTokenPayload func(newAccessTokenPayload map[string]interface{}, userContext supertokens.UserContext) error
	GetTimeCreated           func(userContext supertokens.UserContext) (uint64, error)
	GetExpiry                func(userContext supertokens.UserContext) (uint64, error)
}

type SessionInformation struct {
	SessionHandle      string
	UserId             string
	SessionData        map[string]interface{}
	Expiry             uint64
	AccessTokenPayload map[string]interface{}
	TimeCreated        uint64
}

const SessionContext int = iota
