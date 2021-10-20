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
	"github.com/supertokens/supertokens-golang/recipe/thirdpartyemailpassword/tpepmodels"
)

func GetEmailPasswordIterfaceImpl(apiImplmentation tpepmodels.APIInterface) epmodels.APIInterface {

	result := epmodels.APIInterface{
		EmailExistsGET:                 apiImplmentation.EmailExistsGET,
		GeneratePasswordResetTokenPOST: apiImplmentation.GeneratePasswordResetTokenPOST,
		PasswordResetPOST:              apiImplmentation.PasswordResetPOST,
		SignInPOST:                     nil,
		SignUpPOST:                     nil,
	}

	if apiImplmentation.EmailPasswordSignInPOST != nil {
		result.SignInPOST = func(formFields []epmodels.TypeFormField, options epmodels.APIOptions) (epmodels.SignInResponse, error) {
			result, err := apiImplmentation.EmailPasswordSignInPOST(formFields, options)
			if err != nil {
				return epmodels.SignInResponse{}, err
			}
			if result.OK != nil {
				return epmodels.SignInResponse{
					OK: &struct{ User epmodels.User }{
						User: epmodels.User{
							ID:         result.OK.User.ID,
							Email:      result.OK.User.Email,
							TimeJoined: result.OK.User.TimeJoined,
						},
					},
				}, nil
			} else {
				return epmodels.SignInResponse{
					WrongCredentialsError: &struct{}{},
				}, nil
			}
		}
	}

	if apiImplmentation.EmailPasswordSignUpPOST != nil {
		result.SignUpPOST = func(formFields []epmodels.TypeFormField, options epmodels.APIOptions) (epmodels.SignUpResponse, error) {
			result, err := apiImplmentation.EmailPasswordSignUpPOST(formFields, options)
			if err != nil {
				return epmodels.SignUpResponse{}, err
			}
			if result.OK != nil {
				return epmodels.SignUpResponse{
					OK: &struct{ User epmodels.User }{
						User: epmodels.User{
							ID:         result.OK.User.ID,
							Email:      result.OK.User.Email,
							TimeJoined: result.OK.User.TimeJoined,
						},
					},
				}, nil
			} else {
				return epmodels.SignUpResponse{
					EmailAlreadyExistsError: &struct{}{},
				}, nil
			}
		}
	}

	return result
}
