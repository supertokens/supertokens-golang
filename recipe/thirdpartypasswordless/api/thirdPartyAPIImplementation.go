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
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/recipe/thirdpartypasswordless/tplmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func GetThirdPartyIterfaceImpl(apiImplmentation tplmodels.APIInterface) tpmodels.APIInterface {
	if apiImplmentation.ThirdPartySignInUpPOST == nil || (*apiImplmentation.ThirdPartySignInUpPOST) == nil {
		return tpmodels.APIInterface{
			AuthorisationUrlGET:      apiImplmentation.AuthorisationUrlGET,
			AppleRedirectHandlerPOST: apiImplmentation.AppleRedirectHandlerPOST,
			SignInUpPOST:             nil,
		}
	}

	signInUpPOST := func(provider *tpmodels.TypeProvider, input tpmodels.TypeSignInUpInput, tenantId string, options tpmodels.APIOptions, userContext supertokens.UserContext) (tpmodels.SignInUpPOSTResponse, error) {
		result, err := (*apiImplmentation.ThirdPartySignInUpPOST)(provider, input, tenantId, options, userContext)
		if err != nil {
			return tpmodels.SignInUpPOSTResponse{}, err
		}

		if result.OK != nil {
			return tpmodels.SignInUpPOSTResponse{
				OK: &struct {
					CreatedNewUser          bool
					User                    tpmodels.User
					Session                 *sessmodels.TypeSessionContainer
					OAuthTokens             map[string]interface{}
					RawUserInfoFromProvider tpmodels.TypeRawUserInfoFromProvider
				}{
					CreatedNewUser: result.OK.CreatedNewUser,
					User: tpmodels.User{
						ID:         result.OK.User.ID,
						TimeJoined: result.OK.User.TimeJoined,
						Email:      *result.OK.User.Email,
						TenantIds:  result.OK.User.TenantIds,
						ThirdParty: *result.OK.User.ThirdParty,
					},
					Session:                 result.OK.Session,
					OAuthTokens:             result.OK.OAuthTokens,
					RawUserInfoFromProvider: result.OK.RawUserInfoFromProvider,
				},
			}, nil
		} else if result.NoEmailGivenByProviderError != nil {
			return tpmodels.SignInUpPOSTResponse{
				NoEmailGivenByProviderError: &struct{}{},
			}, nil
		} else {
			return tpmodels.SignInUpPOSTResponse{
				GeneralError: result.GeneralError,
			}, nil
		}
	}

	return tpmodels.APIInterface{
		AuthorisationUrlGET:      apiImplmentation.AuthorisationUrlGET,
		AppleRedirectHandlerPOST: apiImplmentation.AppleRedirectHandlerPOST,
		SignInUpPOST:             &signInUpPOST,
	}
}
