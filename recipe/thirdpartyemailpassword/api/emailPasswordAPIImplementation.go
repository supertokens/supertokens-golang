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
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"github.com/supertokens/supertokens-golang/recipe/thirdpartyemailpassword/tpepmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func GetEmailPasswordIterfaceImpl(apiImplmentation tpepmodels.APIInterface) epmodels.APIInterface {

	result := epmodels.APIInterface{
		EmailExistsGET:                 apiImplmentation.EmailExistsGET,
		GeneratePasswordResetTokenPOST: apiImplmentation.GeneratePasswordResetTokenPOST,
		PasswordResetPOST:              apiImplmentation.PasswordResetPOST,
		SignInPOST:                     nil,
		SignUpPOST:                     nil,
	}

	if apiImplmentation.EmailPasswordSignInPOST != nil && (*apiImplmentation.EmailPasswordSignInPOST) != nil {
		signInPOST := func(formFields []epmodels.TypeFormField, options epmodels.APIOptions, userContext supertokens.UserContext) (epmodels.SignInPOSTResponse, error) {
			result, err := (*apiImplmentation.EmailPasswordSignInPOST)(formFields, options, userContext)
			if err != nil {
				return epmodels.SignInPOSTResponse{}, err
			}
			if result.OK != nil {
				return epmodels.SignInPOSTResponse{
					OK: &struct {
						User    epmodels.User
						Session sessmodels.SessionContainer
					}{

						User: epmodels.User{
							ID:         result.OK.User.ID,
							Email:      result.OK.User.Email,
							TimeJoined: result.OK.User.TimeJoined,
						},
						Session: result.OK.Session,
					},
				}, nil
			} else {
				return epmodels.SignInPOSTResponse{
					WrongCredentialsError: &struct{}{},
				}, nil
			}
		}
		result.SignInPOST = &signInPOST
	}

	if apiImplmentation.EmailPasswordSignUpPOST != nil && (*apiImplmentation.EmailPasswordSignUpPOST) != nil {
		signUpPOST := func(formFields []epmodels.TypeFormField, options epmodels.APIOptions, userContext supertokens.UserContext) (epmodels.SignUpPOSTResponse, error) {
			result, err := (*apiImplmentation.EmailPasswordSignUpPOST)(formFields, options, userContext)
			if err != nil {
				return epmodels.SignUpPOSTResponse{}, err
			}
			if result.OK != nil {
				return epmodels.SignUpPOSTResponse{
					OK: &struct {
						User    epmodels.User
						Session sessmodels.SessionContainer
					}{
						User: epmodels.User{
							ID:         result.OK.User.ID,
							Email:      result.OK.User.Email,
							TimeJoined: result.OK.User.TimeJoined,
						},
						Session: result.OK.Session,
					},
				}, nil
			} else {
				return epmodels.SignUpPOSTResponse{
					EmailAlreadyExistsError: &struct{}{},
				}, nil
			}
		}
		result.SignUpPOST = &signUpPOST
	}

	return result
}
