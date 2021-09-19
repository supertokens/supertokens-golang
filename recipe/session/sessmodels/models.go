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
	Handle        string                 `json:"handle"`
	UserID        string                 `json:"userId"`
	UserDataInJWT map[string]interface{} `json:"userDataInJWT"`
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
	Override                 *OverrideStruct
	ErrorHandlers            *ErrorHandlers
}

type OverrideStruct struct {
	Functions func(originalImplementation RecipeInterface) RecipeInterface
	APIs      func(originalImplementation APIInterface) APIInterface
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
	RevokeSession     func() error
	GetSessionData    func() (map[string]interface{}, error)
	UpdateSessionData func(newSessionData map[string]interface{}) error
	GetUserID         func() string
	GetJWTPayload     func() map[string]interface{}
	GetHandle         func() string
	GetAccessToken    func() string
	UpdateJWTPayload  func(newJWTPayload map[string]interface{}) error
	GetTimeCreated    func() (uint64, error)
	GetExpiry         func() (uint64, error)
}

type SessionInformation struct {
	SessionHandle string
	UserId        string
	SessionData   map[string]interface{}
	Expiry        uint64
	JwtPayload    map[string]interface{}
	TimeCreated   uint64
}

const SessionContext int = iota
