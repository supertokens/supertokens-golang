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
	"github.com/supertokens/supertokens-golang/recipe/emailpassword/epmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func MakeRecipeImplementation(querier supertokens.Querier) epmodels.RecipeInterface {
	signUp := func(email, password string) (epmodels.SignUpResponse, error) {
		response, err := querier.SendPostRequest("/recipe/signup", map[string]interface{}{
			"email":    email,
			"password": password,
		})
		if err != nil {
			return epmodels.SignUpResponse{}, err
		}
		status, ok := response["status"]
		if ok && status.(string) == "OK" {
			user, err := parseUser(response["user"])
			if err != nil {
				return epmodels.SignUpResponse{}, err
			}
			return epmodels.SignUpResponse{
				OK: &struct{ User epmodels.User }{User: *user},
			}, nil
		}
		return epmodels.SignUpResponse{
			EmailAlreadyExistsError: &struct{}{},
		}, nil
	}

	signIn := func(email, password string) (epmodels.SignInResponse, error) {
		response, err := querier.SendPostRequest("/recipe/signin", map[string]interface{}{
			"email":    email,
			"password": password,
		})
		if err != nil {
			return epmodels.SignInResponse{}, err
		}
		status, ok := response["status"]
		if ok && status.(string) == "OK" {
			user, err := parseUser(response["user"])
			if err != nil {
				return epmodels.SignInResponse{}, err
			}
			return epmodels.SignInResponse{
				OK: &struct{ User epmodels.User }{User: *user},
			}, nil
		}
		return epmodels.SignInResponse{
			WrongCredentialsError: &struct{}{},
		}, nil
	}

	getUserByID := func(userID string) (*epmodels.User, error) {
		response, err := querier.SendGetRequest("/recipe/user", map[string]string{
			"userId": userID,
		})
		if err != nil {
			return nil, err
		}
		status, ok := response["status"]
		if ok && status.(string) == "OK" {
			user, err := parseUser(response["user"])
			if err != nil {
				return nil, err
			}
			return user, nil
		}
		return nil, nil
	}

	getUserByEmail := func(email string) (*epmodels.User, error) {
		response, err := querier.SendGetRequest("/recipe/user", map[string]string{
			"email": email,
		})
		if err != nil {
			return nil, err
		}
		status, ok := response["status"]
		if ok && status.(string) == "OK" {
			user, err := parseUser(response["user"])
			if err != nil {
				return nil, err
			}
			return user, nil
		}
		return nil, nil
	}

	createResetPasswordToken := func(userID string) (epmodels.CreateResetPasswordTokenResponse, error) {
		response, err := querier.SendPostRequest("/recipe/user/password/reset/token", map[string]interface{}{
			"userId": userID,
		})
		if err != nil {
			return epmodels.CreateResetPasswordTokenResponse{}, err
		}
		status, ok := response["status"]
		if ok && status.(string) == "OK" {
			return epmodels.CreateResetPasswordTokenResponse{
				OK: &struct{ Token string }{Token: response["token"].(string)},
			}, nil
		}
		return epmodels.CreateResetPasswordTokenResponse{
			UnknownUserIdError: &struct{}{},
		}, nil
	}

	resetPasswordUsingToken := func(token, newPassword string) (epmodels.ResetPasswordUsingTokenResponse, error) {
		response, err := querier.SendPostRequest("/recipe/user/password/reset", map[string]interface{}{
			"method":      "token",
			"token":       token,
			"newPassword": newPassword,
		})
		if err != nil {
			return epmodels.ResetPasswordUsingTokenResponse{}, nil
		}

		if response["status"].(string) == "OK" {
			userId, ok := response["userId"]
			if ok {
				// using CDI >= 2.12
				userIdStr := userId.(string)
				return epmodels.ResetPasswordUsingTokenResponse{
					OK: &struct {
						UserId *string
					}{
						UserId: &userIdStr,
					},
				}, nil
			} else {
				// using CDI < 2.12
				return epmodels.ResetPasswordUsingTokenResponse{
					OK: &struct {
						UserId *string
					}{},
				}, nil
			}
		} else {
			return epmodels.ResetPasswordUsingTokenResponse{
				ResetPasswordInvalidTokenError: &struct{}{},
			}, nil
		}
	}

	updateEmailOrPassword := func(userId string, email, password *string) (epmodels.UpdateEmailOrPasswordResponse, error) {
		requestBody := map[string]interface{}{
			"userId": userId,
		}
		if email != nil {
			requestBody["email"] = email
		}
		if password != nil {
			requestBody["password"] = password
		}
		response, err := querier.SendPutRequest("/recipe/user", requestBody)
		if err != nil {
			return epmodels.UpdateEmailOrPasswordResponse{}, nil
		}

		if response["status"].(string) == "OK" {
			return epmodels.UpdateEmailOrPasswordResponse{
				OK: &struct{}{},
			}, nil
		} else if response["status"].(string) == "EMAIL_ALREADY_EXISTS_ERROR" {
			return epmodels.UpdateEmailOrPasswordResponse{
				EmailAlreadyExistsError: &struct{}{},
			}, nil
		} else {
			return epmodels.UpdateEmailOrPasswordResponse{
				UnknownUserIdError: &struct{}{},
			}, nil
		}
	}
	return epmodels.RecipeInterface{
		SignUp:                   &signUp,
		SignIn:                   &signIn,
		GetUserByID:              &getUserByID,
		GetUserByEmail:           &getUserByEmail,
		CreateResetPasswordToken: &createResetPasswordToken,
		ResetPasswordUsingToken:  &resetPasswordUsingToken,
		UpdateEmailOrPassword:    &updateEmailOrPassword,
	}
}
