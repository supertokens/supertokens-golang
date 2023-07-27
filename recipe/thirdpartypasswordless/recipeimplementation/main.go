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

func MakeRecipeImplementation(passwordlessQuerier supertokens.Querier, thirdPartyQuerier *supertokens.Querier, providers []tpmodels.ProviderInput) tplmodels.RecipeInterface {
	result := tplmodels.RecipeInterface{}

	passwordlessImplementation := passwordless.MakeRecipeImplementation(passwordlessQuerier)
	var thirdPartyImplementation *tpmodels.RecipeInterface
	if thirdPartyQuerier != nil {
		thirdPartyImplementationTemp := thirdparty.MakeRecipeImplementation(*thirdPartyQuerier, providers)
		thirdPartyImplementation = &thirdPartyImplementationTemp
	}

	var ogSignInUp func(thirdPartyID string, thirdPartyUserID string, email string, oAuthTokens tpmodels.TypeOAuthTokens, rawUserInfoFromProvider tpmodels.TypeRawUserInfoFromProvider, tenantId string, userContext supertokens.UserContext) (tpmodels.SignInUpResponse, error) = nil
	if thirdPartyImplementation != nil {
		ogSignInUp = *thirdPartyImplementation.SignInUp
	}
	thirPartySignInUp := func(thirdPartyID string, thirdPartyUserID string, email string, oAuthTokens tpmodels.TypeOAuthTokens, rawUserInfoFromProvider tpmodels.TypeRawUserInfoFromProvider, tenantId string, userContext supertokens.UserContext) (tplmodels.ThirdPartySignInUp, error) {
		if ogSignInUp == nil {
			return tplmodels.ThirdPartySignInUp{}, errors.New("no thirdparty provider configured")
		}
		result, err := ogSignInUp(thirdPartyID, thirdPartyUserID, email, oAuthTokens, rawUserInfoFromProvider, tenantId, userContext)
		if err != nil {
			return tplmodels.ThirdPartySignInUp{}, err
		}
		return tplmodels.ThirdPartySignInUp{
			OK: &struct {
				CreatedNewUser          bool
				User                    tplmodels.User
				OAuthTokens             map[string]interface{}
				RawUserInfoFromProvider tpmodels.TypeRawUserInfoFromProvider
			}{
				CreatedNewUser: result.OK.CreatedNewUser,
				User: tplmodels.User{
					ID:         result.OK.User.ID,
					TimeJoined: result.OK.User.TimeJoined,
					Email:      &result.OK.User.Email,
					ThirdParty: &result.OK.User.ThirdParty,
				},
			},
		}, nil
	}

	var ogManuallyCreateOrUpdateUser func(thirdPartyID string, thirdPartyUserID string, email string, tenantId string, userContext supertokens.UserContext) (tpmodels.ManuallyCreateOrUpdateUserResponse, error) = nil
	if thirdPartyImplementation != nil {
		ogManuallyCreateOrUpdateUser = *thirdPartyImplementation.ManuallyCreateOrUpdateUser
	}

	thirdPartyManuallyCreateOrUpdateUser := func(thirdPartyID string, thirdPartyUserID string, email string, tenantId string, userContext supertokens.UserContext) (tplmodels.ManuallyCreateOrUpdateUserResponse, error) {
		if ogManuallyCreateOrUpdateUser == nil {
			return tplmodels.ManuallyCreateOrUpdateUserResponse{}, errors.New("no thirdparty provider configured")
		}
		result, err := ogManuallyCreateOrUpdateUser(thirdPartyID, thirdPartyUserID, email, tenantId, userContext)
		if err != nil {
			return tplmodels.ManuallyCreateOrUpdateUserResponse{}, err
		}
		return tplmodels.ManuallyCreateOrUpdateUserResponse{
			OK: &struct {
				CreatedNewUser bool
				User           tplmodels.User
			}{
				CreatedNewUser: result.OK.CreatedNewUser,
				User: tplmodels.User{
					ID:         result.OK.User.ID,
					TimeJoined: result.OK.User.TimeJoined,
					Email:      &result.OK.User.Email,
					ThirdParty: &result.OK.User.ThirdParty,
				},
			},
		}, nil
	}

	var ogGetProvider func(thirdPartyID string, clientType *string, tenantId string, userContext supertokens.UserContext) (*tpmodels.TypeProvider, error) = nil
	if thirdPartyImplementation != nil {
		ogGetProvider = *thirdPartyImplementation.GetProvider
	}

	thirdPartyGetProvider := func(thirdPartyID string, clientType *string, tenantId string, userContext supertokens.UserContext) (*tpmodels.TypeProvider, error) {
		if ogGetProvider == nil {
			return nil, errors.New("no thirdparty provider configured")
		}
		return ogGetProvider(thirdPartyID, clientType, tenantId, userContext)
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
				TenantId:    userinfo.TenantId,
				ThirdParty:  &userinfo.ThirdParty,
			}, nil
		}
		return nil, nil
	}

	ogPlessGetUserByEmail := *passwordlessImplementation.GetUserByEmail
	var ogTPGetUsersByEmail func(email string, tenantId string, userContext supertokens.UserContext) ([]tpmodels.User, error) = nil
	if thirdPartyImplementation != nil {
		ogTPGetUsersByEmail = *thirdPartyImplementation.GetUsersByEmail
	}
	getUsersByEmail := func(email string, tenantId string, userContext supertokens.UserContext) ([]tplmodels.User, error) {
		fromPless, err := ogPlessGetUserByEmail(email, tenantId, userContext)
		if err != nil {
			return []tplmodels.User{}, err
		}

		fromTP := []tpmodels.User{}
		if ogTPGetUsersByEmail != nil {
			fromTP, err = ogTPGetUsersByEmail(email, tenantId, userContext)
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

	var ogGetUserByThirdPartyInfo func(thirdPartyID string, thirdPartyUserID string, tenantId string, userContext supertokens.UserContext) (*tpmodels.User, error) = nil
	if thirdPartyImplementation != nil {
		ogGetUserByThirdPartyInfo = *thirdPartyImplementation.GetUserByThirdPartyInfo
	}
	getUserByThirdPartyInfo := func(thirdPartyID string, thirdPartyUserID string, tenantId string, userContext supertokens.UserContext) (*tplmodels.User, error) {
		if ogGetUserByThirdPartyInfo == nil {
			return nil, nil
		}

		userinfo, err := ogGetUserByThirdPartyInfo(thirdPartyID, thirdPartyUserID, tenantId, userContext)
		if err != nil {
			return nil, err
		}

		if userinfo != nil {
			return &tplmodels.User{
				ID:          userinfo.ID,
				Email:       &userinfo.Email,
				PhoneNumber: nil,
				TimeJoined:  userinfo.TimeJoined,
				TenantId:    userinfo.TenantId,
				ThirdParty:  &userinfo.ThirdParty,
			}, nil
		}
		return nil, nil
	}

	ogCreateCode := *passwordlessImplementation.CreateCode
	createCode := func(email *string, phoneNumber *string, userInputCode *string, tenantId string, userContext supertokens.UserContext) (plessmodels.CreateCodeResponse, error) {
		return ogCreateCode(email, phoneNumber, userInputCode, tenantId, userContext)
	}

	ogConsumeCode := *passwordlessImplementation.ConsumeCode
	consumeCode := func(userInput *plessmodels.UserInputCodeWithDeviceID, linkCode *string, preAuthSessionID string, tenantId string, userContext supertokens.UserContext) (tplmodels.ConsumeCodeResponse, error) {
		response, err := ogConsumeCode(userInput, linkCode, preAuthSessionID, tenantId, userContext)
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
	createNewCodeForDevice := func(deviceID string, userInputCode *string, tenantId string, userContext supertokens.UserContext) (plessmodels.ResendCodeResponse, error) {
		return ogCreateNewCodeForDevice(deviceID, userInputCode, tenantId, userContext)
	}

	ogGetUserByPhoneNumber := *passwordlessImplementation.GetUserByPhoneNumber
	getUserByPhoneNumber := func(phoneNumber string, tenantId string, userContext supertokens.UserContext) (*tplmodels.User, error) {
		resp, err := ogGetUserByPhoneNumber(phoneNumber, tenantId, userContext)
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
	listCodesByDeviceID := func(deviceID string, tenantId string, userContext supertokens.UserContext) (*plessmodels.DeviceType, error) {
		return ogListCodesByDeviceID(deviceID, tenantId, userContext)
	}

	ogListCodesByEmail := *passwordlessImplementation.ListCodesByEmail
	listCodesByEmail := func(email string, tenantId string, userContext supertokens.UserContext) ([]plessmodels.DeviceType, error) {
		return ogListCodesByEmail(email, tenantId, userContext)
	}

	ogListCodesByPhoneNumber := *passwordlessImplementation.ListCodesByPhoneNumber
	listCodesByPhoneNumber := func(phoneNumber string, tenantId string, userContext supertokens.UserContext) ([]plessmodels.DeviceType, error) {
		return ogListCodesByPhoneNumber(phoneNumber, tenantId, userContext)
	}

	ogListCodesByPreAuthSessionID := *passwordlessImplementation.ListCodesByPreAuthSessionID
	listCodesByPreAuthSessionID := func(preAuthSessionID string, tenantId string, userContext supertokens.UserContext) (*plessmodels.DeviceType, error) {
		return ogListCodesByPreAuthSessionID(preAuthSessionID, tenantId, userContext)
	}

	ogRevokeAllCodes := *passwordlessImplementation.RevokeAllCodes
	revokeAllCodes := func(email *string, phoneNumber *string, tenantId string, userContext supertokens.UserContext) error {
		return ogRevokeAllCodes(email, phoneNumber, tenantId, userContext)
	}

	ogRevokeCode := *passwordlessImplementation.RevokeCode
	revokeCode := func(codeID string, tenantId string, userContext supertokens.UserContext) error {
		return ogRevokeCode(codeID, tenantId, userContext)
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
	result.ThirdPartyManuallyCreateOrUpdateUser = &thirdPartyManuallyCreateOrUpdateUser
	result.ThirdPartyGetProvider = &thirdPartyGetProvider
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
		(*thirdPartyImplementation.GetProvider) = *modifiedTp.GetProvider
		(*thirdPartyImplementation.SignInUp) = *modifiedTp.SignInUp
		(*thirdPartyImplementation.ManuallyCreateOrUpdateUser) = *modifiedTp.ManuallyCreateOrUpdateUser
		(*thirdPartyImplementation.GetProvider) = *modifiedTp.GetProvider
	}

	return result
}
