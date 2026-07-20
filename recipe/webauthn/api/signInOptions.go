/* Copyright (c) 2025, VRAI Labs and/or its affiliates. All rights reserved.
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
	"github.com/supertokens/supertokens-golang/recipe/webauthn/webauthnmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func SignInOptions(apiImplementation webauthnmodels.APIInterface, tenantId string, options webauthnmodels.APIOptions, userContext supertokens.UserContext) error {
	if apiImplementation.SignInOptionsPOST == nil {
		options.OtherHandler(options.Res, options.Req)
		return nil
	}

	resp, err := (*apiImplementation.SignInOptionsPOST)(tenantId, options, userContext)
	if err != nil {
		return err
	}

	if resp.OK != nil {
		ok := resp.OK
		return supertokens.Send200Response(options.Res, map[string]interface{}{
			"status":                     "OK",
			"webauthnGeneratedOptionsId": ok.WebauthnGeneratedOptionsId,
			"createdAt":                  ok.CreatedAt,
			"expiresAt":                  ok.ExpiresAt,
			"rpId":                       ok.RpId,
			"challenge":                  ok.Challenge,
			"timeout":                    ok.Timeout,
			"userVerification":           string(ok.UserVerification),
		})
	}
	if resp.InvalidOptionsError != nil {
		return supertokens.Send200Response(options.Res, map[string]interface{}{
			"status": "INVALID_OPTIONS_ERROR",
		})
	}
	if resp.GeneralError != nil {
		return supertokens.Send200Response(options.Res, supertokens.ConvertGeneralErrorToJsonResponse(*resp.GeneralError))
	}
	return supertokens.ErrorIfNoResponse(options.Res)
}
