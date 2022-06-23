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
	ThirdPartyId     string                 `json:"thirdPartyId"`
	Code             string                 `json:"code"`
	RedirectURI      string                 `json:"redirectURI"`
	AuthCodeResponse map[string]interface{} `json:"authCodeResponse"`
	ClientId         string                 `json:"clientId"`
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

	if bodyParams.Code == "" && bodyParams.AuthCodeResponse == nil {
		return supertokens.BadInputError{Msg: "Please provide one of code or authCodeResponse in the request body"}
	}

	if bodyParams.AuthCodeResponse != nil && bodyParams.AuthCodeResponse["access_token"] == nil {
		return supertokens.BadInputError{Msg: "Please provide the access_token inside the authCodeResponse request param"}
	}

	if bodyParams.RedirectURI == "" {
		return supertokens.BadInputError{Msg: "Please provide the redirectURI in request body"}
	}

	var provider *tpmodels.TypeProvider = findRightProvider(options.Providers, bodyParams.ThirdPartyId, clientId)

	if provider == nil {
		if clientId == nil {
			return supertokens.BadInputError{Msg: "The third party provider " + bodyParams.ThirdPartyId + " seems to be missing from the backend configs."}
		} else {
			return supertokens.BadInputError{Msg: "The third party provider " + bodyParams.ThirdPartyId + " seems to be missing from the backend configs. If it is configured, then please make sure that you are passing the correct clientId from the frontend."}
		}
	}

	result, err := (*apiImplementation.SignInUpPOST)(*provider, bodyParams.Code, bodyParams.AuthCodeResponse, bodyParams.RedirectURI, options, &map[string]interface{}{})

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
	} else {
		return supertokens.Send200Response(options.Res, supertokens.ConvertGeneralErrorToJsonResponse(*result.GeneralError))
	}
}
