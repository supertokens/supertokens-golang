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

package thirdpartyemailpassword

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/recipe/thirdpartyemailpassword/tpepmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
	"github.com/supertokens/supertokens-golang/test/unittesting"
)

func TestDefaultConfigForThirdPartyEmailPasswordRecipe(t *testing.T) {
	configValue := supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: "http://localhost:8080",
		},
		AppInfo: supertokens.AppInfo{
			APIDomain:     "api.supertokens.io",
			AppName:       "SuperTokens",
			WebsiteDomain: "supertokens.io",
		},
		RecipeList: []supertokens.Recipe{
			Init(
				&tpepmodels.TypeInput{
					Providers: []tpmodels.TypeProvider{},
				},
			),
		},
	}

	BeforeEach()
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()
	err := supertokens.Init(configValue)

	if err != nil {
		t.Error(err.Error())
	}

	thirdpartyemailpassword, err := getRecipeInstanceOrThrowError()
	if err != nil {
		t.Error(err.Error())
	}

	assert.Nil(t, thirdpartyemailpassword.thirdPartyRecipe)

	emailpassword := thirdpartyemailpassword.emailPasswordRecipe
	signUpFeature := emailpassword.Config.SignUpFeature
	assert.Equal(t, 2, len(signUpFeature.FormFields))
	for _, formField := range signUpFeature.FormFields {
		assert.False(t, formField.Optional)
		assert.NotNil(t, formField.Validate)
	}

	signInFeature := emailpassword.Config.SignInFeature
	assert.Equal(t, 2, len(signInFeature.FormFields))
	for _, formField := range signInFeature.FormFields {
		assert.False(t, formField.Optional)
		assert.NotNil(t, formField.Validate)
	}

	resetPasswordUsingTokenFeature := emailpassword.Config.ResetPasswordUsingTokenFeature
	assert.Equal(t, 1, len(resetPasswordUsingTokenFeature.FormFieldsForGenerateTokenForm))
	assert.Equal(t, "email", resetPasswordUsingTokenFeature.FormFieldsForGenerateTokenForm[0].ID)
	assert.Equal(t, 1, len(resetPasswordUsingTokenFeature.FormFieldsForPasswordResetForm))
	assert.Equal(t, "password", resetPasswordUsingTokenFeature.FormFieldsForPasswordResetForm[0].ID)
}

func TestDefaultConfigForThirdPartyEmailPasswordRecipeWithProvider(t *testing.T) {
	configValue := supertokens.TypeInput{
		Supertokens: &supertokens.ConnectionInfo{
			ConnectionURI: "http://localhost:8080",
		},
		AppInfo: supertokens.AppInfo{
			APIDomain:     "api.supertokens.io",
			AppName:       "SuperTokens",
			WebsiteDomain: "supertokens.io",
		},
		RecipeList: []supertokens.Recipe{
			Init(
				&tpepmodels.TypeInput{
					Providers: []tpmodels.TypeProvider{
						customProvider2,
					},
				},
			),
		},
	}

	BeforeEach()
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()
	err := supertokens.Init(configValue)

	if err != nil {
		t.Error(err.Error())
	}

	thirdpartyemailpassword, err := getRecipeInstanceOrThrowError()
	if err != nil {
		t.Error(err.Error())
	}

	assert.NotNil(t, thirdpartyemailpassword.thirdPartyRecipe)

	emailpassword := thirdpartyemailpassword.emailPasswordRecipe
	signUpFeature := emailpassword.Config.SignUpFeature
	assert.Equal(t, 2, len(signUpFeature.FormFields))
	for _, formField := range signUpFeature.FormFields {
		assert.False(t, formField.Optional)
		assert.NotNil(t, formField.Validate)
	}

	signInFeature := emailpassword.Config.SignInFeature
	assert.Equal(t, 2, len(signInFeature.FormFields))
	for _, formField := range signInFeature.FormFields {
		assert.False(t, formField.Optional)
		assert.NotNil(t, formField.Validate)
	}

	resetPasswordUsingTokenFeature := emailpassword.Config.ResetPasswordUsingTokenFeature
	assert.Equal(t, 1, len(resetPasswordUsingTokenFeature.FormFieldsForGenerateTokenForm))
	assert.Equal(t, "email", resetPasswordUsingTokenFeature.FormFieldsForGenerateTokenForm[0].ID)
	assert.Equal(t, 1, len(resetPasswordUsingTokenFeature.FormFieldsForPasswordResetForm))
	assert.Equal(t, "password", resetPasswordUsingTokenFeature.FormFieldsForPasswordResetForm[0].ID)
}
