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

	"github.com/supertokens/supertokens-golang/recipe/passwordless"
	"github.com/supertokens/supertokens-golang/recipe/passwordless/plessmodels"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/recipe/thirdpartypasswordless/tplmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func MakeRecipeImplementation(passwordlessQuerier supertokens.Querier, thirdPartyQuerier *supertokens.Querier) tplmodels.RecipeInterface {
	result := tplmodels.RecipeInterface{}

	passwordlessImplementation := passwordless.MakeRecipeImplementation(passwordlessQuerier)
	var thirdPartyImplementation *tpmodels.RecipeInterface
	if thirdPartyQuerier != nil {
		thirdPartyImplementationTemp := thirdparty.MakeRecipeImplementation(*thirdPartyQuerier)
		thirdPartyImplementation = &thirdPartyImplementationTemp
	}

	var ogSignInUp func(thirdPartyID string, thirdPartyUserID string, email string, responsesFromProvider tpmodels.TypeResponsesFromProvider, userContext supertokens.UserContext) (tpmodels.SignInUpResponse, error) = nil
	if thirdPartyImplementation != nil {
		ogSignInUp = *thirdPartyImplementation.SignInUp
	}
	thirPartySignInUp := func(thirdPartyID, thirdPartyUserID string, email string, responsesFromProvider tpmodels.TypeResponsesFromProvider, userContext supertokens.UserContext) (tplmodels.ThirdPartySignInUp, error) {
		if ogSignInUp == nil {
			return tplmodels.ThirdPartySignInUp{}, errors.New("no thirdparty provider configured")
		}
		result, err := ogSignInUp(thirdPartyID, thirdPartyUserID, email, responsesFromProvider, userContext)
		if err != nil {
			return tplmodels.ThirdPartySignInUp{}, err
		}
		return tplmodels.ThirdPartySignInUp{
			OK: &struct {
				CreatedNewUser bool
				User           tplmodels.User
			}{
				CreatedNewUser: result.OK.CreatedNewUser,
				User: tplmodels.User{
					ID:         result.OK.User.ID,
					Email:      &result.OK.User.Email,
					TimeJoined: result.OK.User.TimeJoined,
					ThirdParty: &result.OK.User.ThirdParty,
				},
			},
		}, nil
	}

	ogPlessGetUserByID := *passwordlessImplementation.GetUserByID
	var ogTPGetUserById func(userID string, userContext supertokens.UserContext) (*tpmodels.User, error) = nil
	if thirdPartyImplementation != nil {
		ogTPGetUserById = *thirdPartyImplementation.GetUserByID
	}
	getUserByID := func(userID string, userContext supertokens.UserContext) (*tplmodels.User, error) {
		user, err := ogPlessGetUserByID(userID, userContext)
		if err != nil {
			return nil, err
		}
		if user != nil {
			return &tplmodels.User{
				ID:          user.ID,
				Email:       user.Email,
				PhoneNumber: user.PhoneNumber,
				TimeJoined:  user.TimeJoined,
				ThirdParty:  nil,
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
			return &tplmodels.User{
				ID:          userinfo.ID,
				Email:       &userinfo.Email,
				PhoneNumber: nil,
				TimeJoined:  userinfo.TimeJoined,
				ThirdParty:  &userinfo.ThirdParty,
			}, nil
		}
		return nil, nil
	}

	ogPlessGetUserByEmail := *passwordlessImplementation.GetUserByEmail
	var ogTPGetUsersByEmail func(email string, userContext supertokens.UserContext) ([]tpmodels.User, error) = nil
	if thirdPartyImplementation != nil {
		ogTPGetUsersByEmail = *thirdPartyImplementation.GetUsersByEmail
	}
	getUsersByEmail := func(email string, userContext supertokens.UserContext) ([]tplmodels.User, error) {
		fromPless, err := ogPlessGetUserByEmail(email, userContext)
		if err != nil {
			return []tplmodels.User{}, err
		}

		fromTP := []tpmodels.User{}
		if ogTPGetUsersByEmail != nil {
			fromTP, err = ogTPGetUsersByEmail(email, userContext)
			if err != nil {
				return []tplmodels.User{}, err
			}
		}
		finalResult := []tplmodels.User{}

		if fromPless != nil {
			finalResult = append(finalResult, tplmodels.User{
				ID:          fromPless.ID,
				TimeJoined:  fromPless.TimeJoined,
				Email:       fromPless.Email,
				PhoneNumber: fromPless.PhoneNumber,
				ThirdParty:  nil,
			})
		}

		for _, tpUser := range fromTP {
			finalResult = append(finalResult, tplmodels.User{
				ID:          tpUser.ID,
				TimeJoined:  tpUser.TimeJoined,
				Email:       &tpUser.Email,
				PhoneNumber: nil,
				ThirdParty:  &tpUser.ThirdParty,
			})
		}

		return finalResult, nil
	}

	var ogGetUserByThirdPartyInfo func(thirdPartyID string, thirdPartyUserID string, userContext supertokens.UserContext) (*tpmodels.User, error) = nil
	if thirdPartyImplementation != nil {
		ogGetUserByThirdPartyInfo = *thirdPartyImplementation.GetUserByThirdPartyInfo
	}
	getUserByThirdPartyInfo := func(thirdPartyID string, thirdPartyUserID string, userContext supertokens.UserContext) (*tplmodels.User, error) {
		if ogGetUserByThirdPartyInfo == nil {
			return nil, nil
		}

		userinfo, err := ogGetUserByThirdPartyInfo(thirdPartyID, thirdPartyUserID, userContext)
		if err != nil {
			return nil, err
		}

		if userinfo != nil {
			return &tplmodels.User{
				ID:          userinfo.ID,
				Email:       &userinfo.Email,
				PhoneNumber: nil,
				TimeJoined:  userinfo.TimeJoined,
				ThirdParty:  &userinfo.ThirdParty,
			}, nil
		}
		return nil, nil
	}

	ogCreateCode := *passwordlessImplementation.CreateCode
	createCode := func(email *string, phoneNumber *string, userInputCode *string, userContext supertokens.UserContext) (plessmodels.CreateCodeResponse, error) {
		return ogCreateCode(email, phoneNumber, userInputCode, userContext)
	}

	ogConsumeCode := *passwordlessImplementation.ConsumeCode
	consumeCode := func(userInput *plessmodels.UserInputCodeWithDeviceID, linkCode *string, preAuthSessionID string, userContext supertokens.UserContext) (tplmodels.ConsumeCodeResponse, error) {
		response, err := ogConsumeCode(userInput, linkCode, preAuthSessionID, userContext)
		if err != nil {
			return tplmodels.ConsumeCodeResponse{}, err
		}

		if response.ExpiredUserInputCodeError != nil {
			return tplmodels.ConsumeCodeResponse{
				ExpiredUserInputCodeError: response.ExpiredUserInputCodeError,
			}, nil
		} else if response.IncorrectUserInputCodeError != nil {
			return tplmodels.ConsumeCodeResponse{
				IncorrectUserInputCodeError: response.IncorrectUserInputCodeError,
			}, nil
		} else if response.RestartFlowError != nil {
			return tplmodels.ConsumeCodeResponse{
				RestartFlowError: &struct{}{},
			}, nil
		}

		return tplmodels.ConsumeCodeResponse{
			OK: &struct {
				CreatedNewUser bool
				User           tplmodels.User
			}{
				CreatedNewUser: response.OK.CreatedNewUser,
				User: tplmodels.User{
					ID:          response.OK.User.ID,
					TimeJoined:  response.OK.User.TimeJoined,
					Email:       response.OK.User.Email,
					PhoneNumber: response.OK.User.PhoneNumber,
					ThirdParty:  nil,
				},
			},
		}, nil
	}

	ogCreateNewCodeForDevice := *passwordlessImplementation.CreateNewCodeForDevice
	createNewCodeForDevice := func(deviceID string, userInputCode *string, userContext supertokens.UserContext) (plessmodels.ResendCodeResponse, error) {
		return ogCreateNewCodeForDevice(deviceID, userInputCode, userContext)
	}

	ogGetUserByPhoneNumber := *passwordlessImplementation.GetUserByPhoneNumber
	getUserByPhoneNumber := func(phoneNumber string, userContext supertokens.UserContext) (*tplmodels.User, error) {
		resp, err := ogGetUserByPhoneNumber(phoneNumber, userContext)
		if err != nil {
			return &tplmodels.User{}, err
		}

		if resp == nil {
			return nil, nil
		}

		return &tplmodels.User{
			ID:          resp.ID,
			TimeJoined:  resp.TimeJoined,
			Email:       resp.Email,
			PhoneNumber: resp.PhoneNumber,
			ThirdParty:  nil,
		}, nil
	}

	ogListCodesByDeviceID := *passwordlessImplementation.ListCodesByDeviceID
	listCodesByDeviceID := func(deviceID string, userContext supertokens.UserContext) (*plessmodels.DeviceType, error) {
		return ogListCodesByDeviceID(deviceID, userContext)
	}

	ogListCodesByEmail := *passwordlessImplementation.ListCodesByEmail
	listCodesByEmail := func(email string, userContext supertokens.UserContext) ([]plessmodels.DeviceType, error) {
		return ogListCodesByEmail(email, userContext)
	}

	ogListCodesByPhoneNumber := *passwordlessImplementation.ListCodesByPhoneNumber
	listCodesByPhoneNumber := func(phoneNumber string, userContext supertokens.UserContext) ([]plessmodels.DeviceType, error) {
		return ogListCodesByPhoneNumber(phoneNumber, userContext)
	}

	ogListCodesByPreAuthSessionID := *passwordlessImplementation.ListCodesByPreAuthSessionID
	listCodesByPreAuthSessionID := func(preAuthSessionID string, userContext supertokens.UserContext) (*plessmodels.DeviceType, error) {
		return ogListCodesByPreAuthSessionID(preAuthSessionID, userContext)
	}

	ogRevokeAllCodes := *passwordlessImplementation.RevokeAllCodes
	revokeAllCodes := func(email *string, phoneNumber *string, userContext supertokens.UserContext) error {
		return ogRevokeAllCodes(email, phoneNumber, userContext)
	}

	ogRevokeCode := *passwordlessImplementation.RevokeCode
	revokeCode := func(codeID string, userContext supertokens.UserContext) error {
		return ogRevokeCode(codeID, userContext)
	}

	ogUpdateUser := *passwordlessImplementation.UpdateUser
	updatePasswordlessUser := func(userID string, email *string, phoneNumber *string, userContext supertokens.UserContext) (plessmodels.UpdateUserResponse, error) {
		user, err := (*result.GetUserByID)(userID, userContext)
		if err != nil {
			return plessmodels.UpdateUserResponse{}, err
		}

		if user == nil {
			return plessmodels.UpdateUserResponse{
				UnknownUserIdError: &struct{}{},
			}, nil
		} else if user.ThirdParty != nil {
			return plessmodels.UpdateUserResponse{}, errors.New("cannot update passwordless user info for those who signed up using third party login")
		}
		return ogUpdateUser(userID, email, phoneNumber, userContext)
	}

	ogDeleteEmailForUser := *passwordlessImplementation.DeleteEmailForUser
	deleteEmailForPasswordlessUser := func(userID string, userContext supertokens.UserContext) (plessmodels.DeleteUserResponse, error) {
		user, err := (*result.GetUserByID)(userID, userContext)
		if err != nil {
			return plessmodels.DeleteUserResponse{}, err
		}

		if user == nil {
			return plessmodels.DeleteUserResponse{
				UnknownUserIdError: &struct{}{},
			}, nil
		} else if user.ThirdParty != nil {
			return plessmodels.DeleteUserResponse{}, errors.New("cannot update passwordless user info for those who signed up using third party login")
		}
		return ogDeleteEmailForUser(userID, userContext)
	}

	ogDeletePhoneNumberForUser := *passwordlessImplementation.DeletePhoneNumberForUser
	deletePhoneNumberForPasswordlessUser := func(userID string, userContext supertokens.UserContext) (plessmodels.DeleteUserResponse, error) {
		user, err := (*result.GetUserByID)(userID, userContext)
		if err != nil {
			return plessmodels.DeleteUserResponse{}, err
		}

		if user == nil {
			return plessmodels.DeleteUserResponse{
				UnknownUserIdError: &struct{}{},
			}, nil
		} else if user.ThirdParty != nil {
			return plessmodels.DeleteUserResponse{}, errors.New("cannot update passwordless user info for those who signed up using third party login")
		}
		return ogDeletePhoneNumberForUser(userID, userContext)
	}

	result.GetUserByID = &getUserByID
	result.GetUsersByEmail = &getUsersByEmail
	result.GetUserByThirdPartyInfo = &getUserByThirdPartyInfo
	result.ThirdPartySignInUp = &thirPartySignInUp
	result.ConsumeCode = &consumeCode
	result.CreateCode = &createCode
	result.CreateNewCodeForDevice = &createNewCodeForDevice
	result.GetUserByPhoneNumber = &getUserByPhoneNumber
	result.ListCodesByDeviceID = &listCodesByDeviceID
	result.ListCodesByEmail = &listCodesByEmail
	result.ListCodesByPhoneNumber = &listCodesByPhoneNumber
	result.ListCodesByPreAuthSessionID = &listCodesByPreAuthSessionID
	result.RevokeAllCodes = &revokeAllCodes
	result.RevokeCode = &revokeCode
	result.UpdatePasswordlessUser = &updatePasswordlessUser
	result.DeleteEmailForPasswordlessUser = &deleteEmailForPasswordlessUser
	result.DeletePhoneNumberForUser = &deletePhoneNumberForPasswordlessUser

	modifiedPwdless := MakePasswordlessRecipeImplementation(result)
	(*passwordlessImplementation.ConsumeCode) = *modifiedPwdless.ConsumeCode
	(*passwordlessImplementation.CreateCode) = *modifiedPwdless.CreateCode
	(*passwordlessImplementation.CreateNewCodeForDevice) = *modifiedPwdless.CreateNewCodeForDevice
	(*passwordlessImplementation.GetUserByEmail) = *modifiedPwdless.GetUserByEmail
	(*passwordlessImplementation.GetUserByID) = *modifiedPwdless.GetUserByID
	(*passwordlessImplementation.GetUserByPhoneNumber) = *modifiedPwdless.GetUserByPhoneNumber
	(*passwordlessImplementation.ListCodesByDeviceID) = *modifiedPwdless.ListCodesByDeviceID
	(*passwordlessImplementation.ListCodesByEmail) = *modifiedPwdless.ListCodesByEmail
	(*passwordlessImplementation.ListCodesByPhoneNumber) = *modifiedPwdless.ListCodesByPhoneNumber
	(*passwordlessImplementation.ListCodesByPreAuthSessionID) = *modifiedPwdless.ListCodesByPreAuthSessionID
	(*passwordlessImplementation.RevokeAllCodes) = *modifiedPwdless.RevokeAllCodes
	(*passwordlessImplementation.RevokeCode) = *modifiedPwdless.RevokeCode
	(*passwordlessImplementation.UpdateUser) = *modifiedPwdless.UpdateUser
	(*passwordlessImplementation.DeleteEmailForUser) = *modifiedPwdless.DeleteEmailForUser
	(*passwordlessImplementation.DeletePhoneNumberForUser) = *modifiedPwdless.DeletePhoneNumberForUser

	if thirdPartyImplementation != nil {
		modifiedTp := MakeThirdPartyRecipeImplementation(result)
		(*thirdPartyImplementation.GetUserByID) = *modifiedTp.GetUserByID
		(*thirdPartyImplementation.GetUserByThirdPartyInfo) = *modifiedTp.GetUserByThirdPartyInfo
		(*thirdPartyImplementation.GetUsersByEmail) = *modifiedTp.GetUsersByEmail
		(*thirdPartyImplementation.SignInUp) = *modifiedTp.SignInUp
	}

	return result
}
