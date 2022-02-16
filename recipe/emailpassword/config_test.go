/*
 * Copyright (c) 2021, VRAI Labs and/or its affiliates. All rights reserved.
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
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/supertokens/supertokens-golang/recipe/emailpassword/epmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
	"github.com/supertokens/supertokens-golang/test/unittesting"
)

func TestDefaultConfigForEmailPasswordModule(t *testing.T) {
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
			Init(nil),
		},
	}
	BeforeEach()
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()
	err := supertokens.Init(configValue)
	if err != nil {
		t.Error(err.Error())
	}
	singletonEmailPasswordInstance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		t.Error(err.Error())
	}
	signupFeature := singletonEmailPasswordInstance.Config.SignUpFeature
	assert.Equal(t, len(signupFeature.FormFields), 2)
	for i := 0; i < len(signupFeature.FormFields); i++ {
		assert.Equal(t, signupFeature.FormFields[i].Optional, false)
		assert.NotNil(t, *signupFeature.FormFields[i].Validate(""))
	}

	singInFeature := singletonEmailPasswordInstance.Config.SignInFeature
	assert.Equal(t, len(singInFeature.FormFields), 2)
	for i := 0; i < len(singInFeature.FormFields); i++ {
		assert.Equal(t, singInFeature.FormFields[i].Optional, false)
	}

	resetPasswordUsingTokenFeature := singletonEmailPasswordInstance.Config.ResetPasswordUsingTokenFeature
	assert.Equal(t, len(resetPasswordUsingTokenFeature.FormFieldsForGenerateTokenForm), 1)
	assert.Equal(t, resetPasswordUsingTokenFeature.FormFieldsForGenerateTokenForm[0].ID, "email")
	assert.Equal(t, len(resetPasswordUsingTokenFeature.FormFieldsForPasswordResetForm), 1)
	assert.Equal(t, resetPasswordUsingTokenFeature.FormFieldsForPasswordResetForm[0].ID, "password")

}

func TestChangedConfigForEmailPasswordModule(t *testing.T) {
	customOptionalValue := false
	customReturnValueFromValidator := "test"
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
			Init(&epmodels.TypeInput{
				SignUpFeature: &epmodels.TypeInputSignUp{
					FormFields: []epmodels.TypeInputFormField{
						{
							ID:       "test",
							Optional: &customOptionalValue,
							Validate: func(value interface{}) *string {
								return &customReturnValueFromValidator
							},
						},
					},
				},
			}),
		},
	}
	BeforeEach()
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()
	err := supertokens.Init(configValue)
	if err != nil {
		t.Error(err.Error())
	}
	singletonEmailPasswordInstance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		t.Error(err.Error())
	}
	signupFeature := singletonEmailPasswordInstance.Config.SignUpFeature
	formFields := signupFeature.FormFields
	assert.Equal(t, len(formFields), 3)
	var testFormField epmodels.NormalisedFormField
	for _, formField := range formFields {
		if formField.ID == "test" {
			testFormField = formField
		}
	}
	assert.NotNil(t, testFormField)
	assert.Equal(t, false, testFormField.Optional)
	assert.Equal(t, "test", *testFormField.Validate(""))

}

func TestNoEmailPasswordValidatorsGivenShouldAddThem(t *testing.T) {
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
			Init(&epmodels.TypeInput{
				SignUpFeature: &epmodels.TypeInputSignUp{
					FormFields: []epmodels.TypeInputFormField{
						{
							ID: "email",
						},
						{
							ID: "password",
						},
					},
				},
			}),
		},
	}
	BeforeEach()
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()
	err := supertokens.Init(configValue)
	if err != nil {
		t.Error(err.Error())
	}
	singletonEmailPasswordInstance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		t.Error(err.Error())
	}
	signupFeature := singletonEmailPasswordInstance.Config.SignUpFeature
	formFields := signupFeature.FormFields
	assert.NotNil(t, *formFields[0].Validate(""))
	assert.NotNil(t, *formFields[1].Validate(""))

}

func TestToCheckTheDefaultEmailPasswordValidators(t *testing.T) {
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
			Init(nil),
		},
	}
	BeforeEach()
	unittesting.StartUpST("localhost", "8080")
	defer AfterEach()
	err := supertokens.Init(configValue)
	if err != nil {
		t.Error(err.Error())
	}
	singletonEmailPasswordInstance, err := getRecipeInstanceOrThrowError()
	if err != nil {
		t.Error(err.Error())
	}
	signupFeature := singletonEmailPasswordInstance.Config.SignUpFeature
	formFields := signupFeature.FormFields
	var emailFormField epmodels.NormalisedFormField
	var passwordFormField epmodels.NormalisedFormField
	for _, formFiled := range formFields {
		if formFiled.ID == "email" {
			emailFormField = formFiled
		} else {
			passwordFormField = formFiled
		}
	}
	assert.Equal(t, "Email is invalid", *emailFormField.Validate("aaaaa"))
	assert.Equal(t, "Email is invalid", *emailFormField.Validate("aaaaa@aaaaa"))
	assert.Equal(t, "Email is invalid", *emailFormField.Validate("random  User   @randomMail.com"))
	assert.Equal(t, "Email is invalid", *emailFormField.Validate("*@*"))
	assert.Nil(t, emailFormField.Validate("validemail@supertokens.io"))

	assert.Equal(t, "Password must contain at least 8 characters, including a number", *passwordFormField.Validate("aaaa"))
	assert.Equal(t, "Password must contain at least one number", *passwordFormField.Validate("aaaaaaaaa"))
	assert.Equal(t, "Password must contain at least one alphabet", *passwordFormField.Validate("1234*-56*789"))
	assert.Nil(t, passwordFormField.Validate("validPass123"))

}
