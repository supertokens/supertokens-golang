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

package api

import (
	"github.com/supertokens/supertokens-golang/recipe/emailpassword/epmodels"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func MakeAPIImplementation() epmodels.APIInterface {
	emailExistsGET := func(email string, options epmodels.APIOptions, userContext supertokens.UserContext) (epmodels.EmailExistsGETResponse, error) {
		user, err := (*options.RecipeImplementation.GetUserByEmail)(email, userContext)
		if err != nil {
			return epmodels.EmailExistsGETResponse{}, err
		}
		return epmodels.EmailExistsGETResponse{
			OK: &struct{ Exists bool }{Exists: user != nil},
		}, nil
	}

	generatePasswordResetTokenPOST := func(formFields []epmodels.TypeFormField, options epmodels.APIOptions, userContext supertokens.UserContext) (epmodels.GeneratePasswordResetTokenPOSTResponse, error) {
		var email string
		for _, formField := range formFields {
			if formField.ID == "email" {
				email = formField.Value
			}
		}

		user, err := (*options.RecipeImplementation.GetUserByEmail)(email, userContext)
		if err != nil {
			return epmodels.GeneratePasswordResetTokenPOSTResponse{}, err
		}

		if user == nil {
			return epmodels.GeneratePasswordResetTokenPOSTResponse{
				OK: &struct{}{},
			}, nil
		}

		response, err := (*options.RecipeImplementation.CreateResetPasswordToken)(user.ID, userContext)
		if err != nil {
			return epmodels.GeneratePasswordResetTokenPOSTResponse{}, err
		}
		if response.UnknownUserIdError != nil {
			return epmodels.GeneratePasswordResetTokenPOSTResponse{
				OK: &struct{}{},
			}, nil
		}

		passwordResetLink, err := options.Config.ResetPasswordUsingTokenFeature.GetResetPasswordURL(*user, userContext)

		if err != nil {
			return epmodels.GeneratePasswordResetTokenPOSTResponse{}, err
		}

		passwordResetLink = passwordResetLink + "?token=" + response.OK.Token + "&rid=" + options.RecipeID

		options.Config.ResetPasswordUsingTokenFeature.CreateAndSendCustomEmail(*user, passwordResetLink, userContext)

		return epmodels.GeneratePasswordResetTokenPOSTResponse{
			OK: &struct{}{},
		}, nil
	}

	passwordResetPOST := func(formFields []epmodels.TypeFormField, token string, options epmodels.APIOptions, userContext supertokens.UserContext) (epmodels.ResetPasswordUsingTokenResponse, error) {
		var newPassword string
		for _, formField := range formFields {
			if formField.ID == "password" {
				newPassword = formField.Value
			}
		}

		response, err := (*options.RecipeImplementation.ResetPasswordUsingToken)(token, newPassword, userContext)
		if err != nil {
			return epmodels.ResetPasswordUsingTokenResponse{}, err
		}

		return response, nil
	}

	signInPOST := func(formFields []epmodels.TypeFormField, options epmodels.APIOptions, userContext supertokens.UserContext) (epmodels.SignInResponse, error) {
		var email string
		for _, formField := range formFields {
			if formField.ID == "email" {
				email = formField.Value
			}
		}
		var password string
		for _, formField := range formFields {
			if formField.ID == "password" {
				password = formField.Value
			}
		}

		response, err := (*options.RecipeImplementation.SignIn)(email, password, userContext)
		if err != nil {
			return epmodels.SignInResponse{}, err
		}
		if response.WrongCredentialsError != nil {
			return response, nil
		}

		user := response.OK.User
		_, err = session.CreateNewSession(options.Res, user.ID, map[string]interface{}{}, map[string]interface{}{}, userContext)
		if err != nil {
			return epmodels.SignInResponse{}, err
		}

		return response, nil
	}

	signUpPOST := func(formFields []epmodels.TypeFormField, options epmodels.APIOptions, userContext supertokens.UserContext) (epmodels.SignUpResponse, error) {
		var email string
		for _, formField := range formFields {
			if formField.ID == "email" {
				email = formField.Value
			}
		}
		var password string
		for _, formField := range formFields {
			if formField.ID == "password" {
				password = formField.Value
			}
		}

		response, err := (*options.RecipeImplementation.SignUp)(email, password, userContext)
		if err != nil {
			return epmodels.SignUpResponse{}, err
		}
		if response.EmailAlreadyExistsError != nil {
			return response, nil
		}

		user := response.OK.User

		_, err = session.CreateNewSession(options.Res, user.ID, map[string]interface{}{}, map[string]interface{}{}, userContext)
		if err != nil {
			return epmodels.SignUpResponse{}, err
		}

		return response, nil
	}
	return epmodels.APIInterface{
		EmailExistsGET:                 &emailExistsGET,
		GeneratePasswordResetTokenPOST: &generatePasswordResetTokenPOST,
		PasswordResetPOST:              &passwordResetPOST,
		SignInPOST:                     &signInPOST,
		SignUpPOST:                     &signUpPOST,
	}
}
