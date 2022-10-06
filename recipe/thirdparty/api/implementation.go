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

	"github.com/derekstavis/go-qs"
	"github.com/supertokens/supertokens-golang/recipe/emailverification"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func MakeAPIImplementation() tpmodels.APIInterface {
	authorisationUrlGET := func(provider tpmodels.TypeProvider, clientID *string, redirectURIOnProviderDashboard string, options tpmodels.APIOptions, userContext supertokens.UserContext) (tpmodels.AuthorisationUrlGETResponse, error) {
		authRedirect, err := provider.GetAuthorisationRedirectURL(clientID, redirectURIOnProviderDashboard, userContext)
		if err != nil {
			return tpmodels.AuthorisationUrlGETResponse{}, err
		}

		return tpmodels.AuthorisationUrlGETResponse{
			OK: &authRedirect,
		}, nil
	}

	signInUpPOST := func(provider tpmodels.TypeProvider, clientID *string, input tpmodels.TypeSignInUpInput, options tpmodels.APIOptions, userContext supertokens.UserContext) (tpmodels.SignInUpPOSTResponse, error) {
		var oAuthTokens map[string]interface{} = nil
		var err error

		if input.RedirectURIInfo != nil {
			oAuthTokens, err = provider.ExchangeAuthCodeForOAuthTokens(clientID, *input.RedirectURIInfo, userContext)
			if err != nil {
				return tpmodels.SignInUpPOSTResponse{}, err
			}
		} else {
			oAuthTokens = *input.OAuthTokens
		}

		userInfo, err := provider.GetUserInfo(clientID, oAuthTokens, userContext)
		if err != nil {
			return tpmodels.SignInUpPOSTResponse{}, err
		}

		emailInfo := userInfo.EmailInfo
		if emailInfo == nil {
			return tpmodels.SignInUpPOSTResponse{
				NoEmailGivenByProviderError: &struct{}{},
			}, nil
		}

		response, err := (*options.RecipeImplementation.SignInUp)(provider.ID, userInfo.ThirdPartyUserId, emailInfo.ID, tpmodels.TypeResponsesFromProvider{
			OAuthTokens: oAuthTokens,
			UserInfo:    userInfo.ResponseFromProvider,
		}, userContext)
		if err != nil {
			return tpmodels.SignInUpPOSTResponse{}, err
		}

		if emailInfo.IsVerified {
			evInstance := emailverification.GetRecipeInstance()
			if evInstance != nil {
				tokenResponse, err := (*evInstance.RecipeImpl.CreateEmailVerificationToken)(response.OK.User.ID, response.OK.User.Email, userContext)
				if err != nil {
					return tpmodels.SignInUpPOSTResponse{}, err
				}
				if tokenResponse.OK != nil {
					_, err := (*evInstance.RecipeImpl.VerifyEmailUsingToken)(tokenResponse.OK.Token, userContext)
					if err != nil {
						return tpmodels.SignInUpPOSTResponse{}, err
					}
				}
			}
		}

		session, err := session.CreateNewSessionWithContext(options.Res, response.OK.User.ID, nil, nil, userContext)
		if err != nil {
			return tpmodels.SignInUpPOSTResponse{}, err
		}
		return tpmodels.SignInUpPOSTResponse{
			OK: &struct {
				CreatedNewUser        bool
				User                  tpmodels.User
				Session               sessmodels.SessionContainer
				ResponsesFromProvider tpmodels.TypeResponsesFromProvider
			}{
				CreatedNewUser: response.OK.CreatedNewUser,
				User:           response.OK.User,
				Session:        session,
				ResponsesFromProvider: tpmodels.TypeResponsesFromProvider{
					OAuthTokens: oAuthTokens,
					UserInfo:    userInfo.ResponseFromProvider,
				},
			},
		}, nil
	}

	appleRedirectHandlerPOST := func(formPostInfoFromProvider map[string]interface{}, options tpmodels.APIOptions, userContext supertokens.UserContext) error {
		queryParams, err := qs.Marshal(formPostInfoFromProvider)
		if err != nil {
			return err
		}
		// TODO extract redirect url from state
		redirectURL := options.AppInfo.WebsiteDomain.GetAsStringDangerous() +
			options.AppInfo.WebsiteBasePath.GetAsStringDangerous() + "/callback/apple?" + queryParams

		options.Res.Header().Set("Location", redirectURL)
		options.Res.WriteHeader(http.StatusFound)

		return nil
	}

	return tpmodels.APIInterface{
		AuthorisationUrlGET:      &authorisationUrlGET,
		SignInUpPOST:             &signInUpPOST,
		AppleRedirectHandlerPOST: &appleRedirectHandlerPOST,
	}
}
