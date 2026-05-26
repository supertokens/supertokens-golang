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
	"encoding/json"

	"github.com/supertokens/supertokens-golang/recipe/webauthn/webauthnmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func SignUp(apiImplementation webauthnmodels.APIInterface, tenantId string, options webauthnmodels.APIOptions, userContext supertokens.UserContext) error {
	if apiImplementation.SignUpPOST == nil {
		options.OtherHandler(options.Res, options.Req)
		return nil
	}

	body, err := supertokens.ReadFromRequest(options.Req)
	if err != nil {
		return err
	}
	var bodyMap map[string]interface{}
	if err := json.Unmarshal(body, &bodyMap); err != nil {
		return err
	}

	optionsId, ok := bodyMap["webauthnGeneratedOptionsId"].(string)
	if !ok || optionsId == "" {
		return supertokens.BadInputError{Msg: "webauthnGeneratedOptionsId is required"}
	}

	credentialRaw, ok := bodyMap["credential"]
	if !ok {
		return supertokens.BadInputError{Msg: "credential is required"}
	}

	credentialBytes, err := json.Marshal(credentialRaw)
	if err != nil {
		return err
	}
	var credential webauthnmodels.RegistrationPayload
	if err := json.Unmarshal(credentialBytes, &credential); err != nil {
		return supertokens.Send200Response(options.Res, map[string]interface{}{
			"status": "INVALID_CREDENTIALS_ERROR",
		})
	}

	resp, err := (*apiImplementation.SignUpPOST)(optionsId, credential, nil, tenantId, options, userContext)
	if err != nil {
		return err
	}

	if resp.OK != nil {
		return supertokens.Send200Response(options.Res, map[string]interface{}{
			"status": "OK",
			"user":   resp.OK.User,
		})
	}
	if resp.EmailAlreadyExistsError != nil {
		return supertokens.Send200Response(options.Res, map[string]interface{}{"status": "EMAIL_ALREADY_EXISTS_ERROR"})
	}
	if resp.OptionsNotFoundError != nil {
		return supertokens.Send200Response(options.Res, map[string]interface{}{"status": "OPTIONS_NOT_FOUND_ERROR"})
	}
	if resp.InvalidOptionsError != nil {
		return supertokens.Send200Response(options.Res, map[string]interface{}{"status": "INVALID_OPTIONS_ERROR"})
	}
	if resp.InvalidCredentialsError != nil {
		return supertokens.Send200Response(options.Res, map[string]interface{}{"status": "INVALID_CREDENTIALS_ERROR"})
	}
	if resp.InvalidAuthenticatorError != nil {
		return supertokens.Send200Response(options.Res, map[string]interface{}{
			"status": "INVALID_AUTHENTICATOR_ERROR",
			"reason": resp.InvalidAuthenticatorError.Reason,
		})
	}
	if resp.GeneralError != nil {
		return supertokens.Send200Response(options.Res, supertokens.ConvertGeneralErrorToJsonResponse(*resp.GeneralError))
	}
	return supertokens.ErrorIfNoResponse(options.Res)
}
