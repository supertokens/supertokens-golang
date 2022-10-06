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
	"encoding/json"

	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

type bodyParams struct {
	ThirdPartyId    string                        `json:"thirdPartyId"`
	ClientId        string                        `json:"clientId"`
	RedirectURIInfo *tpmodels.TypeRedirectURIInfo `json:"redirectURIInfo"`
	OAuthTokens     *tpmodels.TypeOAuthTokens     `json:"oAuthTokens"`
}

func SignInUpAPI(apiImplementation tpmodels.APIInterface, options tpmodels.APIOptions) error {
	if apiImplementation.SignInUpPOST == nil || (*apiImplementation.SignInUpPOST) == nil {
		options.OtherHandler(options.Res, options.Req)
		return nil
	}

	body, err := supertokens.ReadFromRequest(options.Req)
	if err != nil {
		return err
	}
	var bodyParams bodyParams
	err = json.Unmarshal(body, &bodyParams)
	if err != nil {
		return err
	}

	var clientId *string = nil
	if bodyParams.ClientId != "" {
		clientId = &bodyParams.ClientId
	}

	if bodyParams.ThirdPartyId == "" {
		return supertokens.BadInputError{Msg: "Please provide the thirdPartyId in request body"}
	}

	input := tpmodels.TypeSignInUpInput{}
	if bodyParams.RedirectURIInfo != nil {
		input.RedirectURIInfo = bodyParams.RedirectURIInfo
		if bodyParams.RedirectURIInfo.RedirectURIOnProviderDashboard == "" {
			return supertokens.BadInputError{Msg: "Please provide the redirectURIOnProviderDashboard in request body"}
		}

	} else if bodyParams.OAuthTokens != nil {
		input.OAuthTokens = bodyParams.OAuthTokens
	} else {
		return supertokens.BadInputError{Msg: "Please provide one of redirectURIInfo or oAuthTokens in the request body"}
	}

	provider, err := findProvider(options, bodyParams.ThirdPartyId)
	if err != nil {
		return err
	}

	result, err := (*apiImplementation.SignInUpPOST)(*provider, clientId, input, options, supertokens.MakeDefaultUserContextFromAPI(options.Req))

	if err != nil {
		return err
	}

	if result.OK != nil {
		return supertokens.Send200Response(options.Res, map[string]interface{}{
			"status":         "OK",
			"user":           result.OK.User,
			"createdNewUser": result.OK.CreatedNewUser,
		})
	} else if result.NoEmailGivenByProviderError != nil {
		return supertokens.Send200Response(options.Res, map[string]interface{}{
			"status": "NO_EMAIL_GIVEN_BY_PROVIDER",
		})
	} else if result.GeneralError != nil {
		return supertokens.Send200Response(options.Res, supertokens.ConvertGeneralErrorToJsonResponse(*result.GeneralError))
	}
	return supertokens.ErrorIfNoResponse(options.Res)
}
