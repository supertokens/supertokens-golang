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

	if len(thirdPartyId) == 0 {
		return supertokens.BadInputError{Msg: "Please provide the thirdPartyId as a GET param"}
	}

	var provider *tpmodels.TypeProvider = findRightProvider(options.Providers, thirdPartyId, nil)

	if provider == nil {
		return supertokens.BadInputError{Msg: "The third party provider " + thirdPartyId + " seems to be missing from the backend configs."}
	}

	result, err := (*apiImplementation.AuthorisationUrlGET)(*provider, options, &map[string]interface{}{})
	if err != nil {
		return err
	}
	if result.OK != nil {
		return supertokens.Send200Response(options.Res, map[string]interface{}{
			"status": "OK",
			"url":    result.OK.Url,
		})
	} else {
		return supertokens.Send200Response(options.Res, supertokens.ConvertGeneralErrorToJsonResponse(*result.GeneralError))
	}
}
