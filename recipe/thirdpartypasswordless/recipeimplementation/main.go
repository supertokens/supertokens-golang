// /* Copyright (c) 2021, VRAI Labs and/or its affiliates. All rights reserved.
//  *
//  * This software is licensed under the Apache License, Version 2.0 (the
//  * "License") as published by the Apache Software Foundation.
//  *
//  * You may not use this file except in compliance with the License. You may
//  * obtain a copy of the License at http://www.apache.org/licenses/LICENSE-2.0
//  *
//  * Unless required by applicable law or agreed to in writing, software
//  * distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
//  * WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
//  * License for the specific language governing permissions and limitations
//  * under the License.
//  */

package recipeimplementation

import (
	"errors"

	"github.com/supertokens/supertokens-golang/recipe/passwordless"
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

	var ogSignInUp func(thirdPartyID string, thirdPartyUserID string, email tpmodels.EmailStruct, userContext supertokens.UserContext) (tpmodels.SignInUpResponse, error) = nil
	if thirdPartyImplementation != nil {
		ogSignInUp = *thirdPartyImplementation.SignInUp
	}
	thirPartySignInUp := func(thirdPartyID, thirdPartyUserID string, email tplmodels.EmailStruct, userContext supertokens.UserContext) (tplmodels.ThirdPartySignInUp, error) {
		if ogSignInUp == nil {
			return tplmodels.ThirdPartySignInUp{}, errors.New("no thirdparty provider configured")
		}
		result, err := ogSignInUp(thirdPartyID, thirdPartyUserID, tpmodels.EmailStruct{
			ID:         email.ID,
			IsVerified: email.IsVerified,
		}, userContext)
		if err != nil {
			return tplmodels.ThirdPartySignInUp{}, err
		}
		if result.FieldError != nil {
			return tplmodels.ThirdPartySignInUp{
				FieldError: &struct{ ErrorMsg string }{
					ErrorMsg: result.FieldError.ErrorMsg,
				},
			}, nil
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

	result.GetUserByID = &getUserByID
	result.GetUsersByEmail = &getUsersByEmail
	result.GetUserByThirdPartyInfo = &getUserByThirdPartyInfo
	result.ThirdPartySignInUp = &thirPartySignInUp

	// TODO: modified Passwordless

	if thirdPartyImplementation != nil {
		modifiedTp := MakeThirdPartyRecipeImplementation(result)
		(*thirdPartyImplementation.GetUserByID) = *modifiedTp.GetUserByID
		(*thirdPartyImplementation.GetUserByThirdPartyInfo) = *modifiedTp.GetUserByThirdPartyInfo
		(*thirdPartyImplementation.GetUsersByEmail) = *modifiedTp.GetUsersByEmail
		(*thirdPartyImplementation.SignInUp) = *modifiedTp.SignInUp
	}

	return result
}
