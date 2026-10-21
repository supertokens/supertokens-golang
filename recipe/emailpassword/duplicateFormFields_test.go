/* Copyright (c) 2026, VRAI Labs and/or its affiliates. All rights reserved.
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

package emailpassword

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/supertokens/supertokens-golang/recipe/emailpassword/api"
	"github.com/supertokens/supertokens-golang/recipe/emailpassword/epmodels"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/supertokens"
	"github.com/supertokens/supertokens-golang/test/unittesting"
)

func TestSignUpUsesTheFirstOfDuplicateFormFields(t *testing.T) {
	BeforeEach()
	connectionURI := unittesting.StartUpST("localhost", "8080")
	defer AfterEach()
	configValue := supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: connectionURI,
		},
		AppInfo: supertokens.AppInfo{
			APIDomain:     "api.supertokens.io",
			AppName:       "SuperTokens",
			WebsiteDomain: "supertokens.io",
		},
		RecipeList: []supertokens.Recipe{
			Init(&epmodels.TypeInput{
				SignUpFeature: &epmodels.TypeInputSignUp{
					FormFields: []epmodels.TypeInputFormField{
						{
							ID:       "name",
							Optional: func() *bool { b := true; return &b }(),
						},
					},
				},
			}),
			session.Init(nil),
		},
	}
	err := supertokens.Init(configValue)
	assert.NoError(t, err)
	mux := http.NewServeMux()
	testServer := httptest.NewServer(supertokens.Middleware(mux))
	defer testServer.Close()

	postBody, err := json.Marshal(map[string]interface{}{
		"formFields": []map[string]interface{}{
			{"id": "email", "value": "first@example.com"},
			{"id": "password", "value": "validpass123"},
			{"id": "email", "value": "second@example.com"},
		},
	})
	assert.NoError(t, err)

	resp, err := http.Post(testServer.URL+"/auth/signup", "application/json", bytes.NewBuffer(postBody))
	assert.NoError(t, err)
	dataInBytes, err := io.ReadAll(resp.Body)
	assert.NoError(t, err)
	resp.Body.Close()

	// The account is created with the same (first) email that was validated
	assert.Equal(t, 200, resp.StatusCode)
	var data map[string]interface{}
	_ = json.Unmarshal(dataInBytes, &data)
	assert.Equal(t, "OK", data["status"])

	user, err := GetUserByEmail("public", "first@example.com")
	assert.NoError(t, err)
	assert.NotNil(t, user)
	user, err = GetUserByEmail("public", "second@example.com")
	assert.NoError(t, err)
	assert.Nil(t, user)
}

func TestEmailPasswordAPIsUseFirstMatchingFormField(t *testing.T) {
	formFields := []epmodels.TypeFormField{
		{ID: "email", Value: "first@example.com"},
		{ID: "password", Value: "firstpass123"},
		{ID: "email", Value: "second@example.com"},
		{ID: "password", Value: "secondpass123"},
	}

	stopErr := errors.New("stop")
	var signUpEmail, signUpPassword, signInEmail, signInPassword, resetTokenEmail, resetPassword string

	signUp := func(email string, password string, tenantId string, userContext supertokens.UserContext) (epmodels.SignUpResponse, error) {
		signUpEmail, signUpPassword = email, password
		return epmodels.SignUpResponse{}, stopErr
	}
	signIn := func(email string, password string, tenantId string, userContext supertokens.UserContext) (epmodels.SignInResponse, error) {
		signInEmail, signInPassword = email, password
		return epmodels.SignInResponse{}, stopErr
	}
	getUserByEmail := func(email string, tenantId string, userContext supertokens.UserContext) (*epmodels.User, error) {
		resetTokenEmail = email
		return nil, stopErr
	}
	resetPasswordUsingToken := func(token string, newPassword string, tenantId string, userContext supertokens.UserContext) (epmodels.ResetPasswordUsingTokenResponse, error) {
		resetPassword = newPassword
		return epmodels.ResetPasswordUsingTokenResponse{}, stopErr
	}

	options := epmodels.APIOptions{
		RecipeImplementation: epmodels.RecipeInterface{
			SignUp:                  &signUp,
			SignIn:                  &signIn,
			GetUserByEmail:          &getUserByEmail,
			ResetPasswordUsingToken: &resetPasswordUsingToken,
		},
	}
	userContext := &map[string]interface{}{}
	apiImpl := api.MakeAPIImplementation()

	_, err := (*apiImpl.SignUpPOST)(formFields, "public", options, userContext)
	assert.Equal(t, stopErr, err)
	assert.Equal(t, "first@example.com", signUpEmail)
	assert.Equal(t, "firstpass123", signUpPassword)

	_, err = (*apiImpl.SignInPOST)(formFields, "public", options, userContext)
	assert.Equal(t, stopErr, err)
	assert.Equal(t, "first@example.com", signInEmail)
	assert.Equal(t, "firstpass123", signInPassword)

	_, err = (*apiImpl.GeneratePasswordResetTokenPOST)(formFields, "public", options, userContext)
	assert.Equal(t, stopErr, err)
	assert.Equal(t, "first@example.com", resetTokenEmail)

	_, err = (*apiImpl.PasswordResetPOST)(formFields, "token", "public", options, userContext)
	assert.Equal(t, stopErr, err)
	assert.Equal(t, "firstpass123", resetPassword)
}
