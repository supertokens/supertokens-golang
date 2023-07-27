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

	"github.com/supertokens/supertokens-golang/recipe/emailpassword/epmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func GeneratePasswordResetToken(apiImplementation epmodels.APIInterface, tenantId string, options epmodels.APIOptions, userContext supertokens.UserContext) error {
	if apiImplementation.GeneratePasswordResetTokenPOST == nil ||
		(*apiImplementation.GeneratePasswordResetTokenPOST) == nil {
		options.OtherHandler(options.Res, options.Req)
		return nil
	}

	body, err := supertokens.ReadFromRequest(options.Req)
	if err != nil {
		return err
	}
	var formFieldsRaw map[string]interface{}
	err = json.Unmarshal(body, &formFieldsRaw)
	if err != nil {
		return err
	}

	formFields, err := validateFormFieldsOrThrowError(options.Config.ResetPasswordUsingTokenFeature.FormFieldsForGenerateTokenForm, formFieldsRaw["formFields"])
	if err != nil {
		return err
	}

	resp, err := (*apiImplementation.GeneratePasswordResetTokenPOST)(formFields, tenantId, options, userContext)
	if err != nil {
		return err
	}
	if resp.OK != nil {
		return supertokens.Send200Response(options.Res, map[string]interface{}{
			"status": "OK",
		})
	} else if resp.GeneralError != nil {
		return supertokens.Send200Response(options.Res, supertokens.ConvertGeneralErrorToJsonResponse(*resp.GeneralError))
	}
	return supertokens.ErrorIfNoResponse(options.Res)
}
