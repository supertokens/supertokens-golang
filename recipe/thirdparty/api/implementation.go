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
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/supertokens/supertokens-golang/recipe/emailverification"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func MakeAPIImplementation() tpmodels.APIInterface {

	authorisationUrlGET := func(provider *tpmodels.TypeProvider, redirectURIOnProviderDashboard string, tenantId string, options tpmodels.APIOptions, userContext supertokens.UserContext) (tpmodels.AuthorisationUrlGETResponse, error) {
		authRedirect, err := provider.GetAuthorisationRedirectURL(redirectURIOnProviderDashboard, userContext)
		if err != nil {
			return tpmodels.AuthorisationUrlGETResponse{}, err
		}

		return tpmodels.AuthorisationUrlGETResponse{
			OK: &authRedirect,
		}, nil
	}

	signInUpPOST := func(provider *tpmodels.TypeProvider, input tpmodels.TypeSignInUpInput, tenantId string, options tpmodels.APIOptions, userContext supertokens.UserContext) (tpmodels.SignInUpPOSTResponse, error) {
		var oAuthTokens map[string]interface{}
		var err error

		if input.RedirectURIInfo != nil {
			oAuthTokens, err = provider.ExchangeAuthCodeForOAuthTokens(*input.RedirectURIInfo, userContext)
			if err != nil {
				return tpmodels.SignInUpPOSTResponse{}, err
			}
		} else {
			oAuthTokens = *input.OAuthTokens
		}

		userInfo, err := provider.GetUserInfo(oAuthTokens, userContext)
		if err != nil {
			return tpmodels.SignInUpPOSTResponse{}, err
		}

		if userInfo.Email == nil && provider.Config.RequireEmail != nil && !*provider.Config.RequireEmail {
			userInfo.Email = &tpmodels.EmailStruct{
				ID:         provider.Config.GenerateFakeEmail(userInfo.ThirdPartyUserId, tenantId, userContext),
				IsVerified: true,
			}
		}

		emailInfo := userInfo.Email
		if emailInfo == nil {
			return tpmodels.SignInUpPOSTResponse{
				NoEmailGivenByProviderError: &struct{}{},
			}, nil
		}

		response, err := (*options.RecipeImplementation.SignInUp)(provider.ID, userInfo.ThirdPartyUserId, emailInfo.ID, oAuthTokens, userInfo.RawUserInfoFromProvider, tenantId, userContext)
		if err != nil {
			return tpmodels.SignInUpPOSTResponse{}, err
		}

		if emailInfo.IsVerified {
			evInstance := emailverification.GetRecipeInstance()
			if evInstance != nil {
				tokenResponse, err := (*evInstance.RecipeImpl.CreateEmailVerificationToken)(response.OK.User.ID, response.OK.User.Email, tenantId, userContext)
				if err != nil {
					return tpmodels.SignInUpPOSTResponse{}, err
				}
				if tokenResponse.OK != nil {
					_, err := (*evInstance.RecipeImpl.VerifyEmailUsingToken)(tokenResponse.OK.Token, tenantId, userContext)
					if err != nil {
						return tpmodels.SignInUpPOSTResponse{}, err
					}
				}
			}
		}

		session, err := session.CreateNewSession(options.Req, options.Res, tenantId, response.OK.User.ID, nil, nil, userContext)
		if err != nil {
			return tpmodels.SignInUpPOSTResponse{}, err
		}
		return tpmodels.SignInUpPOSTResponse{
			OK: &struct {
				CreatedNewUser          bool
				User                    tpmodels.User
				Session                 sessmodels.SessionContainer
				OAuthTokens             tpmodels.TypeOAuthTokens
				RawUserInfoFromProvider tpmodels.TypeRawUserInfoFromProvider
			}{
				CreatedNewUser:          response.OK.CreatedNewUser,
				User:                    response.OK.User,
				Session:                 session,
				OAuthTokens:             oAuthTokens,
				RawUserInfoFromProvider: userInfo.RawUserInfoFromProvider,
			},
		}, nil
	}

	appleRedirectHandlerPOST := func(formPostInfoFromProvider map[string]interface{}, options tpmodels.APIOptions, userContext supertokens.UserContext) error {
		invalidStateError := supertokens.BadInputError{Msg: "Invalid state in the Apple callback"}
		invalidRedirectError := supertokens.BadInputError{Msg: "Invalid frontendRedirectURI in the Apple callback state"}

		state, ok := formPostInfoFromProvider["state"].(string)
		if !ok || state == "" {
			return invalidStateError
		}

		stateBytes, err := base64.StdEncoding.DecodeString(state)
		if err != nil {
			return invalidStateError
		}

		var stateObj map[string]interface{}
		err = json.Unmarshal(stateBytes, &stateObj)
		if err != nil || stateObj == nil {
			return invalidStateError
		}

		redirectURL, ok := stateObj["frontendRedirectURI"].(string)
		if !ok {
			return invalidRedirectError
		}

		parsedRedirectURL, err := url.Parse(redirectURL)
		if err != nil {
			return invalidRedirectError
		}

		redirectOrigin, err := getComparableOrigin(parsedRedirectURL)
		if err != nil {
			supertokens.LogDebugMessage("appleRedirectHandlerPOST: rejecting frontendRedirectURI: " + err.Error())
			return invalidRedirectError
		}

		websiteOrigin, err := options.AppInfo.GetOrigin(options.Req, userContext)
		if err != nil {
			return err
		}
		parsedWebsiteOrigin, err := url.Parse(websiteOrigin.GetAsStringDangerous())
		if err != nil {
			return err
		}
		expectedOrigin, err := getComparableOrigin(parsedWebsiteOrigin)
		if err != nil {
			supertokens.LogDebugMessage("appleRedirectHandlerPOST: rejecting frontendRedirectURI because the configured origin can not be compared: " + err.Error())
			return invalidRedirectError
		}

		if redirectOrigin != expectedOrigin {
			supertokens.LogDebugMessage(fmt.Sprintf("appleRedirectHandlerPOST: rejecting frontendRedirectURI because its origin (%s) does not match the configured origin (%s)", redirectOrigin, expectedOrigin))
			return invalidRedirectError
		}

		query := parsedRedirectURL.Query()

		for k, v := range formPostInfoFromProvider {
			query.Add(k, fmt.Sprint(v))
		}

		parsedRedirectURL.RawQuery = query.Encode()

		options.Res.Header().Set("Location", parsedRedirectURL.String())
		options.Res.WriteHeader(http.StatusSeeOther)

		return nil
	}

	return tpmodels.APIInterface{
		AuthorisationUrlGET:      &authorisationUrlGET,
		SignInUpPOST:             &signInUpPOST,
		AppleRedirectHandlerPOST: &appleRedirectHandlerPOST,
	}
}

type comparableOrigin struct {
	scheme string
	host   string
	port   string
}

func (o comparableOrigin) String() string {
	return o.scheme + "://" + o.host + ":" + o.port
}

// getComparableOrigin returns the origin of an absolute http(s) URL in a form that can be compared
// directly: the scheme and host are lowercased and the default port is filled in.
func getComparableOrigin(parsedURL *url.URL) (comparableOrigin, error) {
	scheme := strings.ToLower(parsedURL.Scheme)
	if scheme != "http" && scheme != "https" {
		return comparableOrigin{}, errors.New("the URL must be an absolute http or https URL")
	}

	host := strings.ToLower(parsedURL.Hostname())
	if host == "" {
		return comparableOrigin{}, errors.New("the URL must have a host")
	}

	port := parsedURL.Port()
	if port == "" {
		if scheme == "https" {
			port = "443"
		} else {
			port = "80"
		}
	}

	return comparableOrigin{
		scheme: scheme,
		host:   host,
		port:   port,
	}, nil
}
