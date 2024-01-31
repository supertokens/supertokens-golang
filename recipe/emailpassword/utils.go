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

package emailpassword

import (
	"encoding/json"
	"reflect"
	"regexp"

	"github.com/supertokens/supertokens-golang/ingredients/emaildelivery"
	"github.com/supertokens/supertokens-golang/recipe/emailpassword/emaildelivery/backwardCompatibilityService"
	"github.com/supertokens/supertokens-golang/recipe/emailpassword/epmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func validateAndNormaliseUserInput(recipeInstance *Recipe, appInfo supertokens.NormalisedAppinfo, config *epmodels.TypeInput) epmodels.TypeNormalisedInput {

	typeNormalisedInput := makeTypeNormalisedInput(recipeInstance)

	if config != nil && config.SignUpFeature != nil {
		typeNormalisedInput.SignUpFeature = validateAndNormaliseSignupConfig(config.SignUpFeature)
		typeNormalisedInput.ResetPasswordUsingTokenFeature = validateAndNormaliseResetPasswordUsingTokenConfig(typeNormalisedInput.SignUpFeature)
	}

	// we must call this after validateAndNormaliseSignupConfig
	typeNormalisedInput.SignInFeature = validateAndNormaliseSignInConfig(typeNormalisedInput.SignUpFeature)

	typeNormalisedInput.ResetPasswordUsingTokenFeature = validateAndNormaliseResetPasswordUsingTokenConfig(typeNormalisedInput.SignUpFeature)

	typeNormalisedInput.GetEmailDeliveryConfig = func(recipeImpl epmodels.RecipeInterface) emaildelivery.TypeInputWithService {
		sendPasswordResetEmail := DefaultCreateAndSendCustomPasswordResetEmail(appInfo)

		emailService := backwardCompatibilityService.MakeBackwardCompatibilityService(recipeImpl, appInfo, sendPasswordResetEmail)
		if config != nil && config.EmailDelivery != nil && config.EmailDelivery.Service != nil {
			emailService = *config.EmailDelivery.Service
		}
		result := emaildelivery.TypeInputWithService{
			Service: emailService,
		}
		if config != nil && config.EmailDelivery != nil && config.EmailDelivery.Override != nil {
			result.Override = config.EmailDelivery.Override
		}

		return result
	}

	if config != nil && config.Override != nil {
		if config.Override.Functions != nil {
			typeNormalisedInput.Override.Functions = config.Override.Functions
		}
		if config.Override.APIs != nil {
			typeNormalisedInput.Override.APIs = config.Override.APIs
		}
	}

	return typeNormalisedInput
}

func makeTypeNormalisedInput(recipeInstance *Recipe) epmodels.TypeNormalisedInput {
	signUpConfig := validateAndNormaliseSignupConfig(nil)
	return epmodels.TypeNormalisedInput{
		SignUpFeature:                  signUpConfig,
		SignInFeature:                  validateAndNormaliseSignInConfig(signUpConfig),
		ResetPasswordUsingTokenFeature: validateAndNormaliseResetPasswordUsingTokenConfig(signUpConfig),
		Override: epmodels.OverrideStruct{
			Functions: func(originalImplementation epmodels.RecipeInterface) epmodels.RecipeInterface {
				return originalImplementation
			},
			APIs: func(originalImplementation epmodels.APIInterface) epmodels.APIInterface {
				return originalImplementation
			},
		},
	}
}

func validateAndNormaliseResetPasswordUsingTokenConfig(signUpConfig epmodels.TypeNormalisedInputSignUp) epmodels.TypeNormalisedInputResetPasswordUsingTokenFeature {
	normalisedInputResetPasswordUsingTokenFeature := epmodels.TypeNormalisedInputResetPasswordUsingTokenFeature{
		FormFieldsForGenerateTokenForm: nil,
		FormFieldsForPasswordResetForm: nil,
	}

	if len(signUpConfig.FormFields) > 0 {
		var (
			formFieldsForPasswordResetForm []epmodels.NormalisedFormField
			formFieldsForGenerateTokenForm []epmodels.NormalisedFormField
		)
		for _, FormField := range signUpConfig.FormFields {
			if FormField.ID == "password" {
				formFieldsForPasswordResetForm = append(formFieldsForPasswordResetForm, FormField)
			}
			if FormField.ID == "email" {
				formFieldsForGenerateTokenForm = append(formFieldsForGenerateTokenForm, FormField)
			}
		}
		normalisedInputResetPasswordUsingTokenFeature.FormFieldsForGenerateTokenForm = formFieldsForGenerateTokenForm
		normalisedInputResetPasswordUsingTokenFeature.FormFieldsForPasswordResetForm = formFieldsForPasswordResetForm
	}

	return normalisedInputResetPasswordUsingTokenFeature
}
func validateAndNormaliseSignInConfig(signUpConfig epmodels.TypeNormalisedInputSignUp) epmodels.TypeNormalisedInputSignIn {
	return epmodels.TypeNormalisedInputSignIn{
		FormFields: normaliseSignInFormFields(signUpConfig.FormFields),
	}
}

func normaliseSignInFormFields(formFields []epmodels.NormalisedFormField) []epmodels.NormalisedFormField {
	normalisedFormFields := make([]epmodels.NormalisedFormField, 0)
	if len(formFields) > 0 {
		for _, formField := range formFields {
			if formField.ID != "password" && formField.ID != "email" {
				continue
			}
			var validate func(value interface{}, tenantId string) *string
			if formField.ID == "password" {
				validate = defaultValidator
			} else if formField.ID == "email" {
				validate = formField.Validate
			}
			normalisedFormFields = append(normalisedFormFields, epmodels.NormalisedFormField{
				ID:       formField.ID,
				Validate: validate,
				Optional: false,
			})
		}
	}
	return normalisedFormFields
}

func validateAndNormaliseSignupConfig(config *epmodels.TypeInputSignUp) epmodels.TypeNormalisedInputSignUp {
	if config == nil {
		return epmodels.TypeNormalisedInputSignUp{
			FormFields: NormaliseSignUpFormFields(nil),
		}
	}
	return epmodels.TypeNormalisedInputSignUp{
		FormFields: NormaliseSignUpFormFields(config.FormFields),
	}
}

func NormaliseSignUpFormFields(formFields []epmodels.TypeInputFormField) []epmodels.NormalisedFormField {
	var (
		normalisedFormFields     []epmodels.NormalisedFormField
		formFieldPasswordIDCount = 0
		formFieldEmailIDCount    = 0
	)

	if len(formFields) > 0 {
		for _, formField := range formFields {
			var (
				validate func(value interface{}, tenantId string) *string
				optional bool = false
			)
			if formField.ID == "password" {
				formFieldPasswordIDCount++
				validate = defaultPasswordValidator
				if formField.Validate != nil {
					validate = formField.Validate
				}
			} else if formField.ID == "email" {
				formFieldEmailIDCount++
				validate = defaultEmailValidator
				if formField.Validate != nil {
					validate = formField.Validate
				}
			} else {
				validate = defaultValidator
				if formField.Validate != nil {
					validate = formField.Validate
				}
				if formField.Optional != nil {
					optional = *formField.Optional
				}
			}
			normalisedFormFields = append(normalisedFormFields, epmodels.NormalisedFormField{
				ID:       formField.ID,
				Validate: validate,
				Optional: optional,
			})
		}
	}
	if formFieldPasswordIDCount == 0 {
		normalisedFormFields = append(normalisedFormFields, epmodels.NormalisedFormField{
			ID:       "password",
			Validate: defaultPasswordValidator,
			Optional: false,
		})
	}
	if formFieldEmailIDCount == 0 {
		normalisedFormFields = append(normalisedFormFields, epmodels.NormalisedFormField{
			ID:       "email",
			Validate: defaultEmailValidator,
			Optional: false,
		})
	}
	return normalisedFormFields
}

func defaultValidator(_ interface{}, tenantId string) *string {
	return nil
}

func defaultPasswordValidator(value interface{}, tenantId string) *string {
	// length >= 8 && < 100
	// must have a number and a character

	if reflect.TypeOf(value).Kind() != reflect.String {
		msg := "Development bug: Please make sure the password field yields a string"
		return &msg
	}
	if len(value.(string)) < 8 {
		msg := "Password must contain at least 8 characters, including a number"
		return &msg
	}
	if len(value.(string)) >= 100 {
		msg := "Password's length must be lesser than 100 characters"
		return &msg
	}
	alphaCheck, err := regexp.Match(`^.*[A-Za-z]+.*$`, []byte(value.(string)))
	if err != nil || !alphaCheck {
		msg := "Password must contain at least one alphabet"
		return &msg
	}
	numCheck, err := regexp.Match(`^.*[0-9]+.*$`, []byte(value.(string)))
	if err != nil || !numCheck {
		msg := "Password must contain at least one number"
		return &msg
	}
	return nil
}

func defaultEmailValidator(value interface{}, tenantId string) *string {
	if reflect.TypeOf(value).Kind() != reflect.String {
		msg := "Development bug: Please make sure the email field yields a string"
		return &msg
	}
	emailCheck, err := regexp.Match(`^(([^<>()\[\]\\.,;:\s@"]+(\.[^<>()\[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$`, []byte(value.(string)))
	if err != nil || !emailCheck {
		msg := "Email is invalid"
		return &msg
	}
	return nil
}

func parseUser(value interface{}) (*epmodels.User, error) {
	respJSON, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}
	var supertokensUser supertokens.User
	err = json.Unmarshal(respJSON, &supertokensUser)

	if err != nil {
		return nil, err
	}

	epUser := epmodels.User{
		ID:         supertokensUser.ID,
		Email:      supertokensUser.Emails[0],
		TimeJoined: supertokensUser.TimeJoined,
		TenantIds:  supertokensUser.TenantIDs,
	}

	return &epUser, nil
}
