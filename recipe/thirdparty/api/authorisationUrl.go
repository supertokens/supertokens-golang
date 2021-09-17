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
	"reflect"

	"github.com/supertokens/supertokens-golang/recipe/thirdparty/tpmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func AuthorisationUrlAPI(apiImplementation tpmodels.APIInterface, options tpmodels.APIOptions) error {
	if apiImplementation.AuthorisationUrlGET == nil {
		options.OtherHandler(options.Res, options.Req)
		return nil
	}

	queryParams := options.Req.URL.Query()
	thirdPartyId := queryParams.Get("thirdPartyId")

	if len(thirdPartyId) == 0 {
		return supertokens.BadInputError{Msg: "Please provide the thirdPartyId as a GET param"}
	}

	var provider tpmodels.TypeProvider
	for _, prov := range options.Providers {
		if prov.ID == thirdPartyId {
			provider = prov
		}
	}
	if reflect.DeepEqual(provider, tpmodels.TypeProvider{}) {
		return supertokens.BadInputError{Msg: "The third party provider " + thirdPartyId + " seems to not be configured on the backend. Please check your frontend and backend configs."}
	}

	result, err := apiImplementation.AuthorisationUrlGET(provider, options)
	if err != nil {
		return err
	}
	return supertokens.Send200Response(options.Res, map[string]interface{}{
		"status": "OK",
		"url":    result.OK.Url,
	})
}
