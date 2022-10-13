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
	result := tpepmodels.RecipeInterface{}

	emailPasswordImplementation := emailpassword.MakeRecipeImplementation(emailPasswordQuerier)
	var thirdPartyImplementation *tpmodels.RecipeInterface
	if thirdPartyQuerier != nil {
		thirdPartyImplementationTemp := thirdparty.MakeRecipeImplementation(*thirdPartyQuerier)
		thirdPartyImplementation = &thirdPartyImplementationTemp
	}

	ogSignUp := *emailPasswordImplementation.SignUp
	signUp := func(email, password string, userContext supertokens.UserContext) (tpepmodels.SignUpResponse, error) {
		response, err := ogSignUp(email, password, userContext)
		if err != nil {
			return tpepmodels.SignUpResponse{}, err
		}
		if response.EmailAlreadyExistsError != nil {
			return tpepmodels.SignUpResponse{
				EmailAlreadyExistsError: &struct{}{},
			}, nil
		}
		return tpepmodels.SignUpResponse{
			OK: &struct {
				User tpepmodels.User
			}{
				User: tpepmodels.User{
					ID:         response.OK.User.ID,
					Email:      response.OK.User.Email,
					TimeJoined: response.OK.User.TimeJoined,
					ThirdParty: nil,
				},
			},
		}, nil
	}

	ogSignIn := *emailPasswordImplementation.SignIn
	signIn := func(email, password string, userContext supertokens.UserContext) (tpepmodels.SignInResponse, error) {
		response, err := ogSignIn(email, password, userContext)
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

	var ogSignInUp func(thirdPartyID string, thirdPartyUserID string, email string, oAuthTokens tpmodels.TypeOAuthTokens, rawUserInfoFromProvider map[string]interface{}, userContext supertokens.UserContext) (tpmodels.SignInUpResponse, error) = nil
	if thirdPartyImplementation != nil {
		ogSignInUp = *thirdPartyImplementation.SignInUp
	}
	signInUp := func(thirdPartyID string, thirdPartyUserID string, email string, oAuthTokens tpmodels.TypeOAuthTokens, rawUserInfoFromProvider map[string]interface{}, userContext supertokens.UserContext) (tpepmodels.SignInUpResponse, error) {
		if ogSignInUp == nil {
			return tpepmodels.SignInUpResponse{}, errors.New("no thirdparty provider configured")
		}
		result, err := ogSignInUp(thirdPartyID, thirdPartyUserID, email, oAuthTokens, rawUserInfoFromProvider, userContext)
		if err != nil {
			return tpepmodels.SignInUpResponse{}, err
		}

		return tpepmodels.SignInUpResponse{
			OK: &struct {
				CreatedNewUser          bool
				User                    tpepmodels.User
				OAuthTokens             tpmodels.TypeOAuthTokens
				RawUserInfoFromProvider map[string]interface{}
			}{
				CreatedNewUser: result.OK.CreatedNewUser,
				User: tpepmodels.User{
					ID:         result.OK.User.ID,
					Email:      result.OK.User.Email,
					TimeJoined: result.OK.User.TimeJoined,
					ThirdParty: &result.OK.User.ThirdParty,
				},
				OAuthTokens:             result.OK.OAuthTokens,
				RawUserInfoFromProvider: result.OK.RawUserInfoFromProvider,
			},
		}, nil
	}

	var ogCreateOrUpdateUser func(thirdPartyID string, thirdPartyUserID string, email string, userContext supertokens.UserContext) (tpmodels.CreateOrUpdateUserResponse, error) = nil
	if thirdPartyImplementation != nil {
		ogCreateOrUpdateUser = *thirdPartyImplementation.CreateOrUpdateUser
	}
	createOrUpdateUser := func(thirdPartyID string, thirdPartyUserID string, email string, userContext supertokens.UserContext) (tpepmodels.ThirdPartyCreateOrUpdateUserResponse, error) {
		if ogCreateOrUpdateUser == nil {
			return tpepmodels.ThirdPartyCreateOrUpdateUserResponse{}, errors.New("no thirdparty provider configured")
		}
		result, err := ogCreateOrUpdateUser(thirdPartyID, thirdPartyUserID, email, userContext)
		if err != nil {
			return tpepmodels.ThirdPartyCreateOrUpdateUserResponse{}, err
		}

		return tpepmodels.ThirdPartyCreateOrUpdateUserResponse{
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

	ogEPGetUserByID := *emailPasswordImplementation.GetUserByID
	var ogTPGetUserById func(userID string, userContext supertokens.UserContext) (*tpmodels.User, error) = nil
	if thirdPartyImplementation != nil {
		ogTPGetUserById = *thirdPartyImplementation.GetUserByID
	}
	getUserByID := func(userID string, userContext supertokens.UserContext) (*tpepmodels.User, error) {
		user, err := ogEPGetUserByID(userID, userContext)
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
		if ogTPGetUserById == nil {
			return nil, nil
		}

		userinfo, err := ogTPGetUserById(userID, userContext)
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

	ogEPGetUserByEmail := *emailPasswordImplementation.GetUserByEmail
	var ogTPGetUsersByEmail func(email string, userContext supertokens.UserContext) ([]tpmodels.User, error) = nil
	if thirdPartyImplementation != nil {
		ogTPGetUsersByEmail = *thirdPartyImplementation.GetUsersByEmail
	}
	getUsersByEmail := func(email string, userContext supertokens.UserContext) ([]tpepmodels.User, error) {
		fromEP, err := ogEPGetUserByEmail(email, userContext)
		if err != nil {
			return []tpepmodels.User{}, err
		}

		fromTP := []tpmodels.User{}
		if ogTPGetUsersByEmail != nil {
			fromTP, err = ogTPGetUsersByEmail(email, userContext)
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

	var ogGetUserByThirdPartyInfo func(thirdPartyID string, thirdPartyUserID string, userContext supertokens.UserContext) (*tpmodels.User, error) = nil
	if thirdPartyImplementation != nil {
		ogGetUserByThirdPartyInfo = *thirdPartyImplementation.GetUserByThirdPartyInfo
	}
	getUserByThirdPartyInfo := func(thirdPartyID string, thirdPartyUserID string, userContext supertokens.UserContext) (*tpepmodels.User, error) {
		if ogGetUserByThirdPartyInfo == nil {
			return nil, nil
		}

		userinfo, err := ogGetUserByThirdPartyInfo(thirdPartyID, thirdPartyUserID, userContext)
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

	ogCreateResetPasswordToken := *emailPasswordImplementation.CreateResetPasswordToken
	createResetPasswordToken := func(userID string, userContext supertokens.UserContext) (epmodels.CreateResetPasswordTokenResponse, error) {
		return ogCreateResetPasswordToken(userID, userContext)
	}

	ogResetPasswordUsingToken := *emailPasswordImplementation.ResetPasswordUsingToken
	resetPasswordUsingToken := func(token, newPassword string, userContext supertokens.UserContext) (epmodels.ResetPasswordUsingTokenResponse, error) {
		return ogResetPasswordUsingToken(token, newPassword, userContext)
	}

	ogUpdateEmailOrPassword := *emailPasswordImplementation.UpdateEmailOrPassword
	updateEmailOrPassword := func(userId string, email, password *string, userContext supertokens.UserContext) (epmodels.UpdateEmailOrPasswordResponse, error) {
		user, err := (*result.GetUserByID)(userId, userContext)
		if err != nil {
			return epmodels.UpdateEmailOrPasswordResponse{}, err
		}

		if user == nil {
			return epmodels.UpdateEmailOrPasswordResponse{
				UnknownUserIdError: &struct{}{},
			}, nil
		} else if user.ThirdParty != nil {
			return epmodels.UpdateEmailOrPasswordResponse{}, errors.New("cannot update email or password of a user who signed up using third party login")
		}

		return ogUpdateEmailOrPassword(userId, email, password, userContext)
	}

	result.GetUserByID = &getUserByID
	result.GetUsersByEmail = &getUsersByEmail
	result.GetUserByThirdPartyInfo = &getUserByThirdPartyInfo
	result.ThirdPartySignInUp = &signInUp
	result.ThirdPartyCreateOrUpdateUser = &createOrUpdateUser
	result.EmailPasswordSignUp = &signUp
	result.EmailPasswordSignIn = &signIn
	result.CreateResetPasswordToken = &createResetPasswordToken
	result.ResetPasswordUsingToken = &resetPasswordUsingToken
	result.UpdateEmailOrPassword = &updateEmailOrPassword

	modifiedEp := MakeEmailPasswordRecipeImplementation(result)
	(*emailPasswordImplementation.CreateResetPasswordToken) = *modifiedEp.CreateResetPasswordToken
	(*emailPasswordImplementation.GetUserByEmail) = *modifiedEp.GetUserByEmail
	(*emailPasswordImplementation.GetUserByID) = *modifiedEp.GetUserByID
	(*emailPasswordImplementation.ResetPasswordUsingToken) = *modifiedEp.ResetPasswordUsingToken
	(*emailPasswordImplementation.SignIn) = *modifiedEp.SignIn
	(*emailPasswordImplementation.SignUp) = *modifiedEp.SignUp
	(*emailPasswordImplementation.UpdateEmailOrPassword) = *modifiedEp.UpdateEmailOrPassword

	if thirdPartyImplementation != nil {
		modifiedTp := MakeThirdPartyRecipeImplementation(result)
		(*thirdPartyImplementation.GetUserByID) = *modifiedTp.GetUserByID
		(*thirdPartyImplementation.GetUserByThirdPartyInfo) = *modifiedTp.GetUserByThirdPartyInfo
		(*thirdPartyImplementation.GetUsersByEmail) = *modifiedTp.GetUsersByEmail
		(*thirdPartyImplementation.SignInUp) = *modifiedTp.SignInUp
		(*thirdPartyImplementation.CreateOrUpdateUser) = *modifiedTp.CreateOrUpdateUser
	}

	return result
}
