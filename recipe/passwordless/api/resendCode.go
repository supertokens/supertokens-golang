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
	"reflect"

	"github.com/supertokens/supertokens-golang/recipe/passwordless/plessmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func ResendCode(apiImplementation plessmodels.APIInterface, options plessmodels.APIOptions) error {
	if apiImplementation.ResendCodePOST == nil || (*apiImplementation.ResendCodePOST) == nil {
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

	preAuthSessionID, okPreAuthSessionID := readBody["preAuthSessionId"]
	deviceID, okDeviceID := readBody["deviceId"]

	if !okPreAuthSessionID {
		return supertokens.BadInputError{Msg: "Please provide preAuthSessionId"}
	}

	if !okDeviceID {
		return supertokens.BadInputError{Msg: "Please provide deviceId"}
	}

	if reflect.ValueOf(preAuthSessionID).Kind() != reflect.String {
		return supertokens.BadInputError{Msg: "Please make sure that preAuthSessionId is a string"}
	}

	if reflect.ValueOf(deviceID).Kind() != reflect.String {
		return supertokens.BadInputError{Msg: "Please make sure that deviceId is a string"}
	}

	response, err := (*apiImplementation.ResendCodePOST)(deviceID.(string), preAuthSessionID.(string), options, supertokens.MakeDefaultUserContextFromAPI(options.Req))
	if err != nil {
		return err
	}

	var result map[string]interface{}

	if response.OK != nil {
		result = map[string]interface{}{
			"status": "OK",
		}
	} else if response.ResetFlowError != nil {
		result = map[string]interface{}{
			"status": "RESTART_FLOW_ERROR",
		}
	} else if response.GeneralError != nil {
		result = map[string]interface{}{
			"status":  "GENERAL_ERROR",
			"message": response.GeneralError.Message,
		}
	} else {
		return supertokens.ErrorIfNoResponse(options.Res)
	}

	return supertokens.Send200Response(options.Res, result)
}
