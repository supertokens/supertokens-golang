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
	"github.com/supertokens/supertokens-golang/recipe/jwt/jwtmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func GetJWKS(apiImplementation jwtmodels.APIInterface, options jwtmodels.APIOptions, userContext supertokens.UserContext) error {
	if apiImplementation.GetJWKSGET == nil || (*apiImplementation.GetJWKSGET) == nil {
		options.OtherHandler(options.Res, options.Req)
		return nil
	}

	response, err := (*apiImplementation.GetJWKSGET)(options, userContext)
	if err != nil {
		return err
	}

	if response.GeneralError != nil {
		return supertokens.Send200Response(options.Res, supertokens.ConvertGeneralErrorToJsonResponse(*response.GeneralError))
	} else if response.OK != nil {
		options.Res.Header().Set("Access-Control-Allow-Origin", "*")
		return supertokens.Send200Response(options.Res, map[string]interface{}{
			"keys": response.OK.Keys,
		})
	}
	return supertokens.ErrorIfNoResponse(options.Res)
}
