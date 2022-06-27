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
	"net/http"
	"reflect"

	"github.com/supertokens/supertokens-golang/recipe/emailverification/evmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func EmailVerify(apiImplementation evmodels.APIInterface, options evmodels.APIOptions) error {
	var result map[string]interface{}
	if options.Req.Method == http.MethodPost {
		if apiImplementation.VerifyEmailPOST == nil ||
			(*apiImplementation.VerifyEmailPOST) == nil {
			options.OtherHandler(options.Res, options.Req)
			return nil
		}

		body, err := supertokens.ReadFromRequest(options.Req)

		if err != nil {
			return err
		}
		var readBody map[string]interface{}
		err = json.Unmarshal(body, &readBody)
		if err != nil {
			return err
		}
		token, ok := readBody["token"]
		if !ok {
			return supertokens.BadInputError{Msg: "Please provide the email verification token"}
		}
		if reflect.ValueOf(token).Kind() != reflect.String {
			return supertokens.BadInputError{Msg: "The email verification token must be a string"}
		}

		response, err := (*apiImplementation.VerifyEmailPOST)(token.(string), options, &map[string]interface{}{})
		if err != nil {
			return err
		}
		if response.EmailVerificationInvalidTokenError != nil {
			result = map[string]interface{}{
				"status": "EMAIL_VERIFICATION_INVALID_TOKEN_ERROR",
			}
		} else if response.OK != nil {
			result = map[string]interface{}{
				"status": "OK",
				"user":   response.OK.User,
			}
		} else if response.GeneralError != nil {
			result = supertokens.ConvertGeneralErrorToJsonResponse(*response.GeneralError)
		}
	} else {
		if apiImplementation.IsEmailVerifiedGET == nil ||
			(*apiImplementation.IsEmailVerifiedGET) == nil {
			options.OtherHandler(options.Res, options.Req)
			return nil
		}

		isVerified, err := (*apiImplementation.IsEmailVerifiedGET)(options, &map[string]interface{}{})
		if err != nil {
			return err
		}

		if isVerified.OK != nil {
			result = map[string]interface{}{
				"status":     "OK",
				"isVerified": isVerified.OK.IsVerified,
			}
		} else if isVerified.GeneralError != nil {
			result = supertokens.ConvertGeneralErrorToJsonResponse(*isVerified.GeneralError)
		}
	}

	return supertokens.Send200Response(options.Res, result)
}
