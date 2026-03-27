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
	"strings"

	"github.com/supertokens/supertokens-golang/recipe/webauthn/webauthnmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func RegisterOptions(apiImplementation webauthnmodels.APIInterface, tenantId string, options webauthnmodels.APIOptions, userContext supertokens.UserContext) error {
	if apiImplementation.RegisterOptionsPOST == nil {
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

	var email *string
	var recoverAccountToken *string

	if emailVal, ok := bodyMap["email"]; ok {
		if emailStr, ok := emailVal.(string); ok {
			emailStr = strings.TrimSpace(emailStr)
			email = &emailStr
		}
	}
	if tokenVal, ok := bodyMap["recoverAccountToken"]; ok {
		if tokenStr, ok := tokenVal.(string); ok {
			recoverAccountToken = &tokenStr
		}
	}

	if email == nil && recoverAccountToken == nil {
		return supertokens.BadInputError{Msg: "Please provide either email or recoverAccountToken"}
	}

	if email != nil {
		if validateErr := options.Config.ValidateEmailAddress(*email, tenantId, userContext); validateErr != nil {
			return supertokens.Send200Response(options.Res, map[string]interface{}{
				"status": "INVALID_EMAIL_ERROR",
				"err":    *validateErr,
			})
		}
	}

	resp, err := (*apiImplementation.RegisterOptionsPOST)(email, recoverAccountToken, tenantId, options, userContext)
	if err != nil {
		return err
	}

	if resp.OK != nil {
		ok := resp.OK
		excludeCredentials := make([]map[string]interface{}, 0, len(ok.ExcludeCredentials))
		for _, credential := range ok.ExcludeCredentials {
			transports := make([]string, 0, len(credential.Transports))
			for _, transport := range credential.Transports {
				transports = append(transports, string(transport))
			}
			excludeCredentials = append(excludeCredentials, map[string]interface{}{
				"id":         credential.ID,
				"type":       credential.Type,
				"transports": transports,
			})
		}
		pubKeyCredParams := make([]map[string]interface{}, 0, len(ok.PubKeyCredParams))
		for _, param := range ok.PubKeyCredParams {
			pubKeyCredParams = append(pubKeyCredParams, map[string]interface{}{
				"alg":  param.Alg,
				"type": param.Type,
			})
		}

		return supertokens.Send200Response(options.Res, map[string]interface{}{
			"status":                     "OK",
			"webauthnGeneratedOptionsId": ok.WebauthnGeneratedOptionsId,
			"createdAt":                  ok.CreatedAt,
			"expiresAt":                  ok.ExpiresAt,
			"rp": map[string]interface{}{
				"id":   ok.RP.ID,
				"name": ok.RP.Name,
			},
			"user": map[string]interface{}{
				"id":          ok.User.ID,
				"name":        ok.User.Name,
				"displayName": ok.User.DisplayName,
			},
			"challenge":          ok.Challenge,
			"timeout":            ok.Timeout,
			"attestation":        string(ok.Attestation),
			"excludeCredentials": excludeCredentials,
			"pubKeyCredParams":   pubKeyCredParams,
			"authenticatorSelection": map[string]interface{}{
				"requireResidentKey": ok.AuthenticatorSelection.RequireResidentKey,
				"residentKey":        string(ok.AuthenticatorSelection.ResidentKey),
				"userVerification":   string(ok.AuthenticatorSelection.UserVerification),
			},
		})
	}
	if resp.RecoverAccountTokenInvalidError != nil {
		return supertokens.Send200Response(options.Res, map[string]interface{}{
			"status": "RECOVER_ACCOUNT_TOKEN_INVALID_ERROR",
		})
	}
	if resp.InvalidEmailError != nil {
		return supertokens.Send200Response(options.Res, map[string]interface{}{
			"status": "INVALID_EMAIL_ERROR",
			"err":    resp.InvalidEmailError.Err,
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
