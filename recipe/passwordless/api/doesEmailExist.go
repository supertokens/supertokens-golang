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
	"github.com/supertokens/supertokens-golang/recipe/passwordless/plessmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func DoesEmailExist(apiImplementation plessmodels.APIInterface, options plessmodels.APIOptions) error {
	if apiImplementation.EmailExistsGET == nil || (*apiImplementation.EmailExistsGET) == nil {
		options.OtherHandler(options.Res, options.Req)
		return nil
	}
	email := options.Req.URL.Query().Get("email")
	if email == "" {
		return supertokens.BadInputError{Msg: "Please provide the email as a GET param"}
	}
	result, err := (*apiImplementation.EmailExistsGET)(email, options, &map[string]interface{}{})
	if err != nil {
		return err
	}
	if result.OK != nil {
		return supertokens.Send200Response(options.Res, map[string]interface{}{
			"status": "OK",
			"exists": result.OK.Exists,
		})
	} else if result.GeneralError != nil {
		return supertokens.Send200Response(options.Res, supertokens.ConvertGeneralErrorToJsonResponse(*result.GeneralError))
	}
	return supertokens.ErrorIfNoResponse(options.Res)
}
