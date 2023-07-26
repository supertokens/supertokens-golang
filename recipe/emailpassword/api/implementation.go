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
	"fmt"

	"github.com/supertokens/supertokens-golang/ingredients/emaildelivery"
	"github.com/supertokens/supertokens-golang/recipe/emailpassword/epmodels"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
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
			supertokens.LogDebugMessage(fmt.Sprintf("Password reset email not sent, unknown user id: %s", user.ID))
			return epmodels.GeneratePasswordResetTokenPOSTResponse{
				OK: &struct{}{},
			}, nil
		}

		passwordResetLink := fmt.Sprintf(
			"%s%s/reset-password?token=%s&rid=%s",
			options.AppInfo.WebsiteDomain.GetAsStringDangerous(),
			options.AppInfo.WebsiteBasePath.GetAsStringDangerous(),
			response.OK.Token,
			options.RecipeID,
		)

		if err != nil {
			return epmodels.GeneratePasswordResetTokenPOSTResponse{}, err
		}

		supertokens.LogDebugMessage(fmt.Sprintf("Sending password reset email to %s", user.Email))
		err = (*options.EmailDelivery.IngredientInterfaceImpl.SendEmail)(emaildelivery.EmailType{
			PasswordReset: &emaildelivery.PasswordResetType{
				User: emaildelivery.User{
					ID:    user.ID,
					Email: user.Email,
				},
				PasswordResetLink: passwordResetLink,
			},
		}, userContext)
		if err != nil {
			return epmodels.GeneratePasswordResetTokenPOSTResponse{}, err
		}

		return epmodels.GeneratePasswordResetTokenPOSTResponse{
			OK: &struct{}{},
		}, nil
	}

	passwordResetPOST := func(formFields []epmodels.TypeFormField, token string, options epmodels.APIOptions, userContext supertokens.UserContext) (epmodels.ResetPasswordPOSTResponse, error) {
		var newPassword string
		for _, formField := range formFields {
			if formField.ID == "password" {
				newPassword = formField.Value
			}
		}

		response, err := (*options.RecipeImplementation.ResetPasswordUsingToken)(token, newPassword, userContext)
		if err != nil {
			return epmodels.ResetPasswordPOSTResponse{}, err
		}

		if response.OK != nil {
			return epmodels.ResetPasswordPOSTResponse{
				OK: response.OK,
			}, nil
		} else {
			return epmodels.ResetPasswordPOSTResponse{
				ResetPasswordInvalidTokenError: response.ResetPasswordInvalidTokenError,
			}, nil
		}
	}

	signInPOST := func(formFields []epmodels.TypeFormField, options epmodels.APIOptions, userContext supertokens.UserContext) (epmodels.SignInPOSTResponse, error) {
		var email string
		var password string
		for _, formField := range formFields {
			if formField.ID == "email" {
				email = formField.Value
			} else if formField.ID == "password" {
				password = formField.Value
			}
		}

		response, err := (*options.RecipeImplementation.SignIn)(email, password, userContext)
		if err != nil {
			return epmodels.SignInPOSTResponse{}, err
		}
		if response.WrongCredentialsError != nil {
			return epmodels.SignInPOSTResponse{
				WrongCredentialsError: &struct{}{},
			}, nil
		}

		user := response.OK.User
		session, err := session.CreateNewSession(options.Req, options.Res, user.ID, map[string]interface{}{}, map[string]interface{}{}, userContext)
		if err != nil {
			return epmodels.SignInPOSTResponse{}, err
		}

		return epmodels.SignInPOSTResponse{
			OK: &struct {
				User    epmodels.User
				Session sessmodels.SessionContainer
			}{
				User:    response.OK.User,
				Session: session,
			},
		}, nil
	}

	signUpPOST := func(formFields []epmodels.TypeFormField, options epmodels.APIOptions, userContext supertokens.UserContext) (epmodels.SignUpPOSTResponse, error) {
		var email string
		var password string
		for _, formField := range formFields {
			if formField.ID == "email" {
				email = formField.Value
			} else if formField.ID == "password" {
				password = formField.Value
			}
		}

		response, err := (*options.RecipeImplementation.SignUp)(email, password, userContext)
		if err != nil {
			return epmodels.SignUpPOSTResponse{}, err
		}
		if response.EmailAlreadyExistsError != nil {
			return epmodels.SignUpPOSTResponse{
				EmailAlreadyExistsError: &struct{}{},
			}, nil
		}

		user := response.OK.User

		session, err := session.CreateNewSession(options.Req, options.Res, user.ID, map[string]interface{}{}, map[string]interface{}{}, userContext)
		if err != nil {
			return epmodels.SignUpPOSTResponse{}, err
		}

		return epmodels.SignUpPOSTResponse{
			OK: &struct {
				User    epmodels.User
				Session sessmodels.SessionContainer
			}{
				User:    response.OK.User,
				Session: session,
			},
		}, nil
	}
	return epmodels.APIInterface{
		EmailExistsGET:                 &emailExistsGET,
		GeneratePasswordResetTokenPOST: &generatePasswordResetTokenPOST,
		PasswordResetPOST:              &passwordResetPOST,
		SignInPOST:                     &signInPOST,
		SignUpPOST:                     &signUpPOST,
	}
}
