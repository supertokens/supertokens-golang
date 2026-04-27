/* Copyright (c) 2025, VRAI Labs and/or its affiliates. All rights reserved.
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

package webauthn

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/supertokens/supertokens-golang/recipe/webauthn/webauthnmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
	"github.com/supertokens/supertokens-golang/test/unittesting"
)

func TestDefaultConfigSetsValuesCorrectlyForWebauthnRecipe(t *testing.T) {
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
			Init(nil),
		},
	}

	err := supertokens.Init(configValue)
	if err != nil {
		t.Error(err.Error())
	}

	recipeInstance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		t.Error(err.Error())
	}

	req := httptest.NewRequest(http.MethodGet, "https://api.supertokens.io", nil)
	relyingPartyId, err := recipeInstance.Config.GetRelyingPartyId("public", req, nil)
	if err != nil {
		t.Error(err.Error())
	}
	relyingPartyName, err := recipeInstance.Config.GetRelyingPartyName("public", req, nil)
	if err != nil {
		t.Error(err.Error())
	}
	origin, err := recipeInstance.Config.GetOrigin("public", req, nil)
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, "api.supertokens.io", relyingPartyId)
	assert.Equal(t, "SuperTokens", relyingPartyName)
	assert.Equal(t, "https://supertokens.io", origin)
	assert.Nil(t, recipeInstance.Config.ValidateEmailAddress("user@example.com", "public", nil))
	assert.Equal(t, "Email is not valid", *recipeInstance.Config.ValidateEmailAddress("invalid-email", "public", nil))
}

func TestCustomConfigOverridesAreAppliedForWebauthnRecipe(t *testing.T) {
	BeforeEach()
	connectionURI := unittesting.StartUpST("localhost", "8080")
	defer AfterEach()

	customValidateError := "custom invalid email"
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
			Init(&webauthnmodels.TypeInput{
				GetRelyingPartyId: func(tenantId string, req *http.Request, userContext supertokens.UserContext) (string, error) {
					return "custom-rp-id", nil
				},
				GetRelyingPartyName: func(tenantId string, userContext supertokens.UserContext) (string, error) {
					return "Custom RP Name", nil
				},
				GetOrigin: func(tenantId string, req *http.Request, userContext supertokens.UserContext) (string, error) {
					return "https://custom.example.com", nil
				},
				ValidateEmailAddress: func(email string, tenantId string, userContext supertokens.UserContext) *string {
					if email == "blocked@example.com" {
						return &customValidateError
					}
					return nil
				},
			}),
		},
	}

	err := supertokens.Init(configValue)
	if err != nil {
		t.Error(err.Error())
	}

	recipeInstance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		t.Error(err.Error())
	}

	req := httptest.NewRequest(http.MethodGet, "https://api.supertokens.io", nil)
	relyingPartyId, err := recipeInstance.Config.GetRelyingPartyId("public", req, nil)
	if err != nil {
		t.Error(err.Error())
	}
	relyingPartyName, err := recipeInstance.Config.GetRelyingPartyName("public", req, nil)
	if err != nil {
		t.Error(err.Error())
	}
	origin, err := recipeInstance.Config.GetOrigin("public", req, nil)
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, "custom-rp-id", relyingPartyId)
	assert.Equal(t, "Custom RP Name", relyingPartyName)
	assert.Equal(t, "https://custom.example.com", origin)
	assert.Nil(t, recipeInstance.Config.ValidateEmailAddress("allowed@example.com", "public", nil))
	assert.Equal(t, "custom invalid email", *recipeInstance.Config.ValidateEmailAddress("blocked@example.com", "public", nil))
}

func TestGetRecoverAccountLinkReturnsExpectedValue(t *testing.T) {
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
			Init(nil),
		},
	}

	err := supertokens.Init(configValue)
	if err != nil {
		t.Error(err.Error())
	}

	recipeInstance, err := GetRecipeInstanceOrThrowError()
	if err != nil {
		t.Error(err.Error())
	}

	req := httptest.NewRequest(http.MethodGet, "https://api.supertokens.io", nil)
	link, err := GetRecoverAccountLink(recipeInstance.RecipeModule.GetAppInfo(), "token-123", "tenant1", req, nil)
	if err != nil {
		t.Error(err.Error())
	}

	assert.Equal(t, "https://supertokens.io/auth/webauthn/recover?token=token-123&tenantId=tenant1", link)
}
