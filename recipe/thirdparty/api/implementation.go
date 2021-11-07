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
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"

	"github.com/derekstavis/go-qs"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
)

func MakeAPIImplementation() tpmodels.APIInterface {
	authorisationUrlGET := func(provider tpmodels.TypeProvider, options tpmodels.APIOptions) (tpmodels.AuthorisationUrlGETResponse, error) {
		providerInfo := provider.Get(nil, nil)
		params := map[string]string{}
		for key, value := range providerInfo.AuthorisationRedirect.Params {
			if reflect.ValueOf(value).Kind() == reflect.String {
				params[key] = value.(string)
			} else {
				call, ok := value.(func(req *http.Request) string)
				if ok {
					params[key] = call(options.Req)
				} else {
					return tpmodels.AuthorisationUrlGETResponse{}, errors.New("type of value in params must be a string or a function")
				}
			}
		}

		if providerInfo.GetRedirectURI != nil && !isUsingDevelopmentClientId(providerInfo.GetClientId()) {
			// the backend wants to set the redirectURI - so we set that here.

			// we add the not development keys because the oauth provider will
			// redirect to supertokens.io's URL which will redirect the app
			// to the the user's website, which will handle the callback as usual.
			// If we add this, then instead, the supertokens' site will redirect
			// the user to this API layer, which is not needed.
			rU, err := providerInfo.GetRedirectURI()
			if err != nil {
				return tpmodels.AuthorisationUrlGETResponse{}, err
			}
			params["redirect_uri"] = rU
		}

		if isUsingDevelopmentClientId(providerInfo.GetClientId()) {
			params["actual_redirect_uri"] = providerInfo.AuthorisationRedirect.URL

			for key, value := range params {
				if value == providerInfo.GetClientId() {
					params[key] = GetActualClientIdFromDevelopmentClientId(providerInfo.GetClientId())
				}
			}

		}

		paramsString, err := getParamString(params)
		if err != nil {
			return tpmodels.AuthorisationUrlGETResponse{}, err
		}
		url := providerInfo.AuthorisationRedirect.URL + "?" + paramsString

		if isUsingDevelopmentClientId(providerInfo.GetClientId()) {
			url = DevOauthAuthorisationUrl + "?" + paramsString
		}

		return tpmodels.AuthorisationUrlGETResponse{
			OK: &struct{ Url string }{
				Url: url,
			},
		}, nil
	}

	signInUpPOST := func(provider tpmodels.TypeProvider, code string, authCodeResponse interface{}, redirectURI string, options tpmodels.APIOptions) (tpmodels.SignInUpPOSTResponse, error) {
		{
			providerInfo := provider.Get(nil, nil)
			if isUsingDevelopmentClientId(providerInfo.GetClientId()) {
				redirectURI = DevOauthRedirectUrl
			} else if providerInfo.GetRedirectURI != nil {
				// we overwrite the redirectURI provided by the frontend
				// since the backend wants to take charge of setting this.
				rU, err := providerInfo.GetRedirectURI()
				if err != nil {
					return tpmodels.SignInUpPOSTResponse{}, err
				}
				redirectURI = rU
			}
		}

		providerInfo := provider.Get(&redirectURI, &code)

		var accessTokenAPIResponse map[string]interface{} = nil

		if authCodeResponse != nil && len(authCodeResponse.(map[string]interface{})) != 0 {
			accessTokenAPIResponse = authCodeResponse.(map[string]interface{})
		} else {
			if isUsingDevelopmentClientId(providerInfo.GetClientId()) {

				for key, value := range providerInfo.AccessTokenAPI.Params {
					if value == providerInfo.GetClientId() {
						providerInfo.AccessTokenAPI.Params[key] = GetActualClientIdFromDevelopmentClientId(providerInfo.GetClientId())
					}
				}
			}

			accessTokenAPIResponseTemp, err := postRequest(providerInfo)
			if err != nil {
				return tpmodels.SignInUpPOSTResponse{}, err
			}
			accessTokenAPIResponse = accessTokenAPIResponseTemp
		}

		userInfo, err := providerInfo.GetProfileInfo(accessTokenAPIResponse)
		if err != nil {
			return tpmodels.SignInUpPOSTResponse{}, err
		}

		emailInfo := userInfo.Email
		if emailInfo == nil {
			return tpmodels.SignInUpPOSTResponse{
				NoEmailGivenByProviderError: &struct{}{},
			}, nil
		}

		response, err := (*options.RecipeImplementation.SignInUp)(provider.ID, userInfo.ID, *emailInfo)
		if err != nil {
			return tpmodels.SignInUpPOSTResponse{}, err
		}
		if response.FieldError != nil {
			return tpmodels.SignInUpPOSTResponse{
				FieldError: &struct{ Error string }{
					Error: response.FieldError.Error,
				},
			}, nil
		}

		if emailInfo.IsVerified {
			tokenResponse, err := (*options.EmailVerificationRecipeImplementation.CreateEmailVerificationToken)(response.OK.User.ID, response.OK.User.Email)
			if err != nil {
				return tpmodels.SignInUpPOSTResponse{}, err
			}
			if tokenResponse.OK != nil {
				_, err := (*options.EmailVerificationRecipeImplementation.VerifyEmailUsingToken)(tokenResponse.OK.Token)
				if err != nil {
					return tpmodels.SignInUpPOSTResponse{}, err
				}
			}
		}

		_, err = session.CreateNewSession(options.Res, response.OK.User.ID, nil, nil)
		if err != nil {
			return tpmodels.SignInUpPOSTResponse{}, err
		}
		return tpmodels.SignInUpPOSTResponse{
			OK: &struct {
				CreatedNewUser   bool
				User             tpmodels.User
				AuthCodeResponse interface{}
			}{
				CreatedNewUser:   response.OK.CreatedNewUser,
				User:             response.OK.User,
				AuthCodeResponse: accessTokenAPIResponse,
			},
		}, nil
	}

	appleRedirectHandlerPOST := func(code string, state string, options tpmodels.APIOptions) error {
		redirectURL := options.AppInfo.WebsiteDomain.GetAsStringDangerous() +
			options.AppInfo.WebsiteBasePath.GetAsStringDangerous() + "/callback/apple?state=" + state + "&code=" + code

		options.Res.Header().Set("Content-Type", "text/html; charset=utf-8")

		fmt.Fprint(options.Res, "<html><head><script>window.location.replace(\""+redirectURL+"\");</script></head></html>")
		return nil
	}

	return tpmodels.APIInterface{
		AuthorisationUrlGET:      &authorisationUrlGET,
		SignInUpPOST:             &signInUpPOST,
		AppleRedirectHandlerPOST: &appleRedirectHandlerPOST,
	}
}

func postRequest(providerInfo tpmodels.TypeProviderGetResponse) (map[string]interface{}, error) {
	querystring, err := getParamString(providerInfo.AccessTokenAPI.Params)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", providerInfo.AccessTokenAPI.URL, bytes.NewBuffer([]byte(querystring)))
	if err != nil {
		return nil, err
	}
	req.Header.Set("content-type", "application/x-www-form-urlencoded")
	req.Header.Set("accept", "application/json") // few providers like github don't send back json response by default

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}
	var result map[string]interface{}
	err = json.Unmarshal(body, &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func getParamString(paramsMap map[string]string) (string, error) {
	params := map[string]interface{}{}
	for key, value := range paramsMap {
		params[key] = value
	}
	return qs.Marshal(params)
}

// If Third Party login is used with one of the following development keys, then the dev authorization url and the redirect url will be used.

var DevOauthClientIds = [...]string{
	"1060725074195-kmeum4crr01uirfl2op9kd5acmi9jutn.apps.googleusercontent.com", // google
	"467101b197249757c71f", // github
}

const (
	DevOauthAuthorisationUrl = "https://supertokens.io/dev/oauth/redirect-to-provider"
	DevOauthRedirectUrl      = "https://supertokens.io/dev/oauth/redirect-to-app"
	DevKeyIdentifier         = "4398792-"
)

func isUsingDevelopmentClientId(clientId string) bool {

	if strings.HasPrefix(clientId, DevKeyIdentifier) {
		return true
	} else {
		for _, devClientId := range DevOauthClientIds {
			if devClientId == clientId {
				return true
			}
		}
		return false
	}
}

func GetActualClientIdFromDevelopmentClientId(clientId string) string {
	if strings.HasPrefix(clientId, DevKeyIdentifier) {
		return strings.Split(clientId, DevKeyIdentifier)[1]
	}
	return clientId
}
