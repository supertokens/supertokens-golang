/*
 * Copyright (c) 2024, VRAI Labs and/or its affiliates. All rights reserved.
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

package session

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/supertokens/supertokens-golang/recipe/session/claims"
	sessionErrors "github.com/supertokens/supertokens-golang/recipe/session/errors"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
	"github.com/supertokens/supertokens-golang/test/unittesting"

	"github.com/stretchr/testify/assert"
)

func TestSessionErrorHandlerOverides(t *testing.T) {
	BeforeEach()

	customAntiCsrfVal := "VIA_TOKEN"
	configValue := supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: "http://localhost:8080",
		},
		AppInfo: supertokens.AppInfo{
			AppName:       "SuperTokens",
			WebsiteDomain: "supertokens.io",
			APIDomain:     "api.supertokens.io",
		},
		RecipeList: []supertokens.Recipe{
			Init(&sessmodels.TypeInput{
				AntiCsrf: &customAntiCsrfVal,
				ErrorHandlers: &sessmodels.ErrorHandlers{
					OnUnauthorised: func(message string, req *http.Request, res http.ResponseWriter) error {
						res.WriteHeader(401)
						res.Write([]byte("unauthorised from errorHandler"))
						return nil
					},
					OnTokenTheftDetected: func(sessionHandle, userID string, req *http.Request, res http.ResponseWriter) error {
						res.WriteHeader(403)
						res.Write([]byte("token theft detected from errorHandler"))
						return nil
					},
					OnTryRefreshToken: func(message string, req *http.Request, res http.ResponseWriter) error {
						res.WriteHeader(401)
						res.Write([]byte("try refresh token from errorHandler"))
						return nil
					},
					OnInvalidClaim: func(validationErrors []claims.ClaimValidationError, req *http.Request, res http.ResponseWriter) error {
						res.WriteHeader(403)
						res.Write([]byte("invalid claim from errorHandler"))
						return nil
					},
					OnClearDuplicateSessionCookies: func(message string, req *http.Request, res http.ResponseWriter) error {
						res.WriteHeader(200)
						res.Write([]byte("clear duplicate session cookies from errorHandler"))
						return nil
					},
				},
				GetTokenTransferMethod: func(req *http.Request, forCreateNewSession bool, userContext supertokens.UserContext) sessmodels.TokenTransferMethod {
					return sessmodels.CookieTransferMethod
				},
			}),
		},
	}

	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()
	err := supertokens.Init(configValue)
	if err != nil {
		t.Error(err.Error())
	}

	mux := http.NewServeMux()

	mux.HandleFunc("/test/unauthorized", func(rw http.ResponseWriter, r *http.Request) {
		supertokens.ErrorHandler(sessionErrors.UnauthorizedError{}, r, rw)
	})

	mux.HandleFunc("/test/try-refresh", func(rw http.ResponseWriter, r *http.Request) {
		supertokens.ErrorHandler(sessionErrors.TryRefreshTokenError{}, r, rw)
	})

	mux.HandleFunc("/test/token-theft", func(rw http.ResponseWriter, r *http.Request) {
		supertokens.ErrorHandler(sessionErrors.TokenTheftDetectedError{}, r, rw)
	})

	mux.HandleFunc("/test/claim-validation", func(rw http.ResponseWriter, r *http.Request) {
		supertokens.ErrorHandler(sessionErrors.InvalidClaimError{}, r, rw)
	})

	mux.HandleFunc("/test/clear-duplicate-session", func(rw http.ResponseWriter, r *http.Request) {
		supertokens.ErrorHandler(sessionErrors.ClearDuplicateSessionCookiesError{}, r, rw)
	})

	testServer := httptest.NewServer(supertokens.Middleware(mux))
	defer func() {
		testServer.Close()
	}()

	t.Run("should override session errorHandlers", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet, testServer.URL+"/test/unauthorized", nil)
		assert.NoError(t, err)

		res, err := http.DefaultClient.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, 401, res.StatusCode)

		content, err := io.ReadAll(res.Body)
		assert.NoError(t, err)
		assert.Equal(t, `unauthorised from errorHandler`, string(content))

		req, err = http.NewRequest(http.MethodGet, testServer.URL+"/test/try-refresh", nil)
		assert.NoError(t, err)

		res, err = http.DefaultClient.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, 401, res.StatusCode)

		content, err = io.ReadAll(res.Body)
		assert.NoError(t, err)
		assert.Equal(t, `try refresh token from errorHandler`, string(content))

		req, err = http.NewRequest(http.MethodGet, testServer.URL+"/test/token-theft", nil)
		assert.NoError(t, err)

		res, err = http.DefaultClient.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, 403, res.StatusCode)

		content, err = io.ReadAll(res.Body)
		assert.NoError(t, err)
		assert.Equal(t, `token theft detected from errorHandler`, string(content))

		req, err = http.NewRequest(http.MethodGet, testServer.URL+"/test/claim-validation", nil)
		assert.NoError(t, err)

		res, err = http.DefaultClient.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, 403, res.StatusCode)

		content, err = io.ReadAll(res.Body)
		assert.NoError(t, err)
		assert.Equal(t, `invalid claim from errorHandler`, string(content))

		req, err = http.NewRequest(http.MethodGet, testServer.URL+"/test/clear-duplicate-session", nil)
		assert.NoError(t, err)

		res, err = http.DefaultClient.Do(req)
		assert.NoError(t, err)
		assert.Equal(t, 200, res.StatusCode)

		content, err = io.ReadAll(res.Body)
		assert.NoError(t, err)
		assert.Equal(t, `clear duplicate session cookies from errorHandler`, string(content))
	})
}
