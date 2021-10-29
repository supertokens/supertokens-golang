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

package recipeimplementation

import (
	"errors"

	"github.com/supertokens/supertokens-golang/recipe/emailpassword"
	"github.com/supertokens/supertokens-golang/recipe/emailpassword/epmodels"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/recipe/thirdpartyemailpassword/tpepmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func MakeRecipeImplementation(emailPasswordQuerier supertokens.Querier, thirdPartyQuerier *supertokens.Querier) tpepmodels.RecipeInterface {
	emailPasswordImplementation := emailpassword.MakeRecipeImplementation(emailPasswordQuerier)
	var thirdPartyImplementation *tpmodels.RecipeInterface
	if thirdPartyQuerier != nil {
		thirdPartyImplementationTemp := thirdparty.MakeRecipeImplementation(*thirdPartyQuerier)
		thirdPartyImplementation = &thirdPartyImplementationTemp
	}
	signUp := func(email, password string) (tpepmodels.SignUpResponse, error) {
		response, err := (*emailPasswordImplementation.SignUp)(email, password)
		if err != nil {
			return tpepmodels.SignUpResponse{}, err
		}
		if response.EmailAlreadyExistsError != nil {
			return tpepmodels.SignUpResponse{
				EmailAlreadyExistsError: &struct{}{},
			}, nil
		}
		return tpepmodels.SignUpResponse{
			OK: &struct{ User tpepmodels.User }{
				User: tpepmodels.User{
					ID:         response.OK.User.ID,
					Email:      response.OK.User.Email,
					TimeJoined: response.OK.User.TimeJoined,
					ThirdParty: nil,
				},
			},
		}, nil
	}

	signIn := func(email, password string) (tpepmodels.SignInResponse, error) {
		response, err := (*emailPasswordImplementation.SignIn)(email, password)
		if err != nil {
			return tpepmodels.SignInResponse{}, err
		}
		if response.WrongCredentialsError != nil {
			return tpepmodels.SignInResponse{
				WrongCredentialsError: &struct{}{},
			}, nil
		}
		return tpepmodels.SignInResponse{
			OK: &struct{ User tpepmodels.User }{
				User: tpepmodels.User{
					ID:         response.OK.User.ID,
					Email:      response.OK.User.Email,
					TimeJoined: response.OK.User.TimeJoined,
					ThirdParty: nil,
				},
			},
		}, nil
	}

	signInUp := func(thirdPartyID, thirdPartyUserID string, email tpepmodels.EmailStruct) (tpepmodels.SignInUpResponse, error) {
		if thirdPartyImplementation == nil {
			return tpepmodels.SignInUpResponse{}, errors.New("no thirdparty provider configured")
		}
		result, err := (*thirdPartyImplementation.SignInUp)(thirdPartyID, thirdPartyUserID, tpmodels.EmailStruct{
			ID:         email.ID,
			IsVerified: email.IsVerified,
		})
		if err != nil {
			return tpepmodels.SignInUpResponse{}, err
		}
		if result.FieldError != nil {
			return tpepmodels.SignInUpResponse{
				FieldError: &struct{ Error string }{
					Error: result.FieldError.Error,
				},
			}, nil
		}
		return tpepmodels.SignInUpResponse{
			OK: &struct {
				CreatedNewUser bool
				User           tpepmodels.User
			}{
				CreatedNewUser: result.OK.CreatedNewUser,
				User: tpepmodels.User{
					ID:         result.OK.User.ID,
					Email:      result.OK.User.Email,
					TimeJoined: result.OK.User.TimeJoined,
					ThirdParty: &result.OK.User.ThirdParty,
				},
			},
		}, nil
	}

	getUserByID := func(userID string) (*tpepmodels.User, error) {
		user, err := (*emailPasswordImplementation.GetUserByID)(userID)
		if err != nil {
			return nil, err
		}
		if user != nil {
			return &tpepmodels.User{
				ID:         user.ID,
				Email:      user.Email,
				TimeJoined: user.TimeJoined,
				ThirdParty: nil,
			}, nil
		}
		if thirdPartyImplementation == nil {
			return nil, nil
		}

		userinfo, err := (*thirdPartyImplementation.GetUserByID)(userID)
		if err != nil {
			return nil, err
		}

		if userinfo != nil {
			return &tpepmodels.User{
				ID:         userinfo.ID,
				Email:      userinfo.Email,
				TimeJoined: userinfo.TimeJoined,
				ThirdParty: &userinfo.ThirdParty,
			}, nil
		}
		return nil, nil
	}

	getUsersByEmail := func(email string) ([]tpepmodels.User, error) {
		fromEP, err := (*emailPasswordImplementation.GetUserByEmail)(email)
		if err != nil {
			return []tpepmodels.User{}, err
		}

		fromTP := []tpmodels.User{}
		if thirdPartyImplementation != nil {
			fromTP, err = (*thirdPartyImplementation.GetUsersByEmail)(email)
			if err != nil {
				return []tpepmodels.User{}, err
			}
		}
		finalResult := []tpepmodels.User{}

		if fromEP != nil {
			finalResult = append(finalResult, tpepmodels.User{
				ID:         fromEP.ID,
				TimeJoined: fromEP.TimeJoined,
				Email:      fromEP.Email,
				ThirdParty: nil,
			})
		}

		for _, tpUser := range fromTP {
			finalResult = append(finalResult, tpepmodels.User{
				ID:         tpUser.ID,
				TimeJoined: tpUser.TimeJoined,
				Email:      tpUser.Email,
				ThirdParty: &tpUser.ThirdParty,
			})
		}

		return finalResult, nil
	}

	getUserByThirdPartyInfo := func(thirdPartyID string, thirdPartyUserID string) (*tpepmodels.User, error) {
		if thirdPartyImplementation == nil {
			return nil, nil
		}

		userinfo, err := (*thirdPartyImplementation.GetUserByThirdPartyInfo)(thirdPartyID, thirdPartyUserID)
		if err != nil {
			return nil, err
		}

		if userinfo != nil {
			return &tpepmodels.User{
				ID:         userinfo.ID,
				Email:      userinfo.Email,
				TimeJoined: userinfo.TimeJoined,
				ThirdParty: &userinfo.ThirdParty,
			}, nil
		}
		return nil, nil
	}

	createResetPasswordToken := func(userID string) (epmodels.CreateResetPasswordTokenResponse, error) {
		return (*emailPasswordImplementation.CreateResetPasswordToken)(userID)
	}

	resetPasswordUsingToken := func(token, newPassword string) (epmodels.ResetPasswordUsingTokenResponse, error) {
		return (*emailPasswordImplementation.ResetPasswordUsingToken)(token, newPassword)
	}

	updateEmailOrPassword := func(userId string, email, password *string) (epmodels.UpdateEmailOrPasswordResponse, error) {
		return (*emailPasswordImplementation.UpdateEmailOrPassword)(userId, email, password)
	}

	return tpepmodels.RecipeInterface{
		GetUserByID:              &getUserByID,
		GetUsersByEmail:          &getUsersByEmail,
		GetUserByThirdPartyInfo:  &getUserByThirdPartyInfo,
		SignInUp:                 &signInUp,
		SignUp:                   &signUp,
		SignIn:                   &signIn,
		CreateResetPasswordToken: &createResetPasswordToken,
		ResetPasswordUsingToken:  &resetPasswordUsingToken,
		UpdateEmailOrPassword:    &updateEmailOrPassword,
	}
}
