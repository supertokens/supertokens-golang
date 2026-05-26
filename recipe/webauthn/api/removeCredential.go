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

func RemoveCredential(apiImplementation webauthnmodels.APIInterface, tenantId string, options webauthnmodels.APIOptions, userContext supertokens.UserContext) error {
	if apiImplementation.RemoveCredentialPOST == nil {
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

	webauthnCredentialId, ok := bodyMap["webauthnCredentialId"].(string)
	if !ok || webauthnCredentialId == "" {
		return supertokens.BadInputError{Msg: "A valid webauthnCredentialId is required"}
	}

	sess, err := getSessionWithoutClaimValidation(options, userContext)
	if err != nil {
		return err
	}

	resp, err := (*apiImplementation.RemoveCredentialPOST)(webauthnCredentialId, sess, tenantId, options, userContext)
	if err != nil {
		return err
	}

	if resp.OK != nil {
		return supertokens.Send200Response(options.Res, map[string]interface{}{
			"status": "OK",
		})
	}
	if resp.CredentialNotFoundError != nil {
		return supertokens.Send200Response(options.Res, map[string]interface{}{
			"status": "CREDENTIAL_NOT_FOUND_ERROR",
		})
	}
	if resp.GeneralError != nil {
		return supertokens.Send200Response(options.Res, supertokens.ConvertGeneralErrorToJsonResponse(*resp.GeneralError))
	}

	return supertokens.ErrorIfNoResponse(options.Res)
}
