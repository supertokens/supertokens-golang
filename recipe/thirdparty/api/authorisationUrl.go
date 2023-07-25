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
	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func AuthorisationUrlAPI(apiImplementation tpmodels.APIInterface, options tpmodels.APIOptions) error {
	if apiImplementation.AuthorisationUrlGET == nil || (*apiImplementation.AuthorisationUrlGET) == nil {
		options.OtherHandler(options.Res, options.Req)
		return nil
	}

	queryParams := options.Req.URL.Query()
	thirdPartyId := queryParams.Get("thirdPartyId")
	redirectURIOnProviderDashboard := queryParams.Get("redirectURIOnProviderDashboard")

	var clientType *string
	if clientTypeStr := queryParams.Get("clientType"); clientTypeStr != "" {
		clientType = &clientTypeStr
	}

	if len(thirdPartyId) == 0 {
		return supertokens.BadInputError{Msg: "Please provide the thirdPartyId as a GET param"}
	}

	userContext := supertokens.MakeDefaultUserContextFromAPI(options.Req)

	providerResponse, err := (*options.RecipeImplementation.GetProvider)(thirdPartyId, clientType, "public", userContext) // TODO multitenancy pass tenantId
	if err != nil {
		return err
	}

	provider := providerResponse

	if provider == nil {
		return supertokens.BadInputError{Msg: "the provider " + thirdPartyId + " could not be found in the configuration"}
	}

	result, err := (*apiImplementation.AuthorisationUrlGET)(provider, redirectURIOnProviderDashboard, options, userContext)
	if err != nil {
		return err
	}
	if result.OK != nil {
		respBody := map[string]interface{}{
			"status":             "OK",
			"urlWithQueryParams": result.OK.URLWithQueryParams,
		}
		if result.OK.PKCECodeVerifier != nil {
			respBody["pkceCodeVerifier"] = *result.OK.PKCECodeVerifier
		}
		return supertokens.Send200Response(options.Res, respBody)
	} else if result.GeneralError != nil {
		return supertokens.Send200Response(options.Res, supertokens.ConvertGeneralErrorToJsonResponse(*result.GeneralError))
	}
	return supertokens.ErrorIfNoResponse(options.Res)
}
