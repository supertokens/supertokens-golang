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
	"errors"

	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/recipe/thirdpartyemailpassword/tpepmodels"
)

func GetThirdPartyIterfaceImpl(apiImplmentation tpepmodels.APIInterface) tpmodels.APIInterface {
	signInUpPOST := apiImplmentation.SignInUpPOST
	if signInUpPOST == nil {
		return tpmodels.APIInterface{
			AuthorisationUrlGET: apiImplmentation.AuthorisationUrlGET,
			SignInUpPOST:        nil,
		}
	}
	return tpmodels.APIInterface{

		AuthorisationUrlGET: apiImplmentation.AuthorisationUrlGET,

		SignInUpPOST: func(provider tpmodels.TypeProvider, code, redirectURI string, options tpmodels.APIOptions) (tpmodels.SignInUpPOSTResponse, error) {
			resp, err := signInUpPOST(tpepmodels.SignInUpAPIInput{
				ThirdPartyInput: &tpepmodels.ThirdPartyInput{
					Provider:    provider,
					Code:        code,
					RedirectURI: redirectURI,
					Options:     options,
				},
			})
			if err != nil {
				return tpmodels.SignInUpPOSTResponse{}, err
			}
			result := resp.ThirdPartyOutput
			if result != nil {
				if result.OK != nil {
					return tpmodels.SignInUpPOSTResponse{
						OK: &struct {
							CreatedNewUser   bool
							User             tpmodels.User
							AuthCodeResponse interface{}
						}{
							CreatedNewUser: result.OK.CreatedNewUser,
							User: tpmodels.User{
								ID:         result.OK.User.ID,
								TimeJoined: result.OK.User.TimeJoined,
								Email:      result.OK.User.Email,
								ThirdParty: *result.OK.User.ThirdParty,
							},
						},
					}, nil
				} else if result.NoEmailGivenByProviderError != nil {
					return tpmodels.SignInUpPOSTResponse{
						NoEmailGivenByProviderError: &struct{}{},
					}, nil
				} else if result.FieldError != nil {
					return tpmodels.SignInUpPOSTResponse{
						FieldError: &struct{ Error string }{
							Error: result.FieldError.Error,
						},
					}, nil
				}
			}
			return tpmodels.SignInUpPOSTResponse{}, errors.New("should never come here")
		},
	}
}
