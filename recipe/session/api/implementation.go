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
	"net/http"

	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func MakeAPIImplementation() sessmodels.APIInterface {
	refreshPOST := func(options sessmodels.APIOptions, userContext supertokens.UserContext) error {
		_, err := (*options.RecipeImplementation.RefreshSession)(options.Req, options.Res, userContext)
		return err
	}

	verifySession := func(verifySessionOptions *sessmodels.VerifySessionOptions, options sessmodels.APIOptions, userContext supertokens.UserContext) (*sessmodels.SessionContainer, error) {
		method := options.Req.Method
		if method == http.MethodOptions || method == http.MethodTrace {
			return nil, nil
		}

		incomingPath, err := supertokens.NewNormalisedURLPath(options.Req.RequestURI)
		if err != nil {
			return nil, err
		}

		refreshTokenPath := options.Config.RefreshTokenPath
		if incomingPath.Equals(refreshTokenPath) && method == http.MethodPost {
			session, err := (*options.RecipeImplementation.RefreshSession)(options.Req, options.Res, userContext)
			return &session, err
		} else {
			sessionContainer, err := (*options.RecipeImplementation.GetSession)(options.Req, options.Res, verifySessionOptions, userContext)
			if err != nil {
				return nil, err
			}

			if sessionContainer == nil {
				return nil, nil
			}

			claimValidators, err := getRequiredClaimValidators(options.ClaimValidatorsAddedByOtherRecipes, options.RecipeImplementation, sessionContainer, verifySessionOptions.OverrideGlobalClaimValidators, userContext)
			if err != nil {
				return nil, err
			}
			err = sessionContainer.AssertClaimsWithContext(claimValidators, userContext)
			if err != nil {
				return nil, err
			}

			return sessionContainer, nil
		}
	}

	signOutPOST := func(options sessmodels.APIOptions, sessionContainer *sessmodels.SessionContainer, userContext supertokens.UserContext) (sessmodels.SignOutPOSTResponse, error) {
		if sessionContainer != nil {
			err := sessionContainer.RevokeSessionWithContext(userContext)
			if err != nil {
				return sessmodels.SignOutPOSTResponse{}, err
			}
		}

		return sessmodels.SignOutPOSTResponse{
			OK: &struct{}{},
		}, nil
	}

	return sessmodels.APIInterface{
		RefreshPOST:   &refreshPOST,
		VerifySession: &verifySession,
		SignOutPOST:   &signOutPOST,
	}
}
