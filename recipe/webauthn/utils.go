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

package webauthn

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"

	"github.com/supertokens/supertokens-golang/ingredients/emaildelivery"
	"github.com/supertokens/supertokens-golang/recipe/webauthn/webauthnmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

const SUPERTOKENS_WEBAUTHN_RECOVER_ACCOUNT_URL = "https://api.supertokens.com/0/st/auth/webauthn/recover"

func makeRecoverAccountRequest(appInfo supertokens.NormalisedAppinfo, input emaildelivery.WebauthnRecoverAccountType) (*http.Request, error) {
	data := map[string]any{
		"email":             input.User.Email,
		"appName":           appInfo.AppName,
		"recoverAccountURL": input.RecoverAccountLink,
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", SUPERTOKENS_WEBAUTHN_RECOVER_ACCOUNT_URL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}
	req.Header.Set("content-type", "application/json; charset=utf-8")
	req.Header.Set("api-version", "0")
	return req, nil
}

func createAndSendEmailUsingSupertokensService(appInfo supertokens.NormalisedAppinfo, input emaildelivery.WebauthnRecoverAccountType) error {
	if supertokens.IsRunningInTestMode() {
		return nil
	}

	req, err := makeRecoverAccountRequest(appInfo, input)
	if err != nil {
		return err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		supertokens.LogDebugMessage(fmt.Sprintf("Error sending webauthn recover account email: %s", err.Error()))
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 300 {
		supertokens.LogDebugMessage(fmt.Sprintf("Email sent to %s", input.User.Email))
		return nil
	}

	body, readErr := io.ReadAll(resp.Body)
	if readErr == nil {
		supertokens.LogDebugMessage(fmt.Sprintf("Error sending webauthn recover account email. API returned %d status. Response: %s", resp.StatusCode, string(body)))
	} else {
		supertokens.LogDebugMessage(fmt.Sprintf("Error sending webauthn recover account email. API returned %d status.", resp.StatusCode))
	}
	return fmt.Errorf("error sending webauthn recover account email. API returned %d status", resp.StatusCode)
}

func makeDefaultEmailService(appInfo supertokens.NormalisedAppinfo) emaildelivery.EmailDeliveryInterface {
	sendEmail := func(input emaildelivery.EmailType, userContext supertokens.UserContext) error {
		if input.WebauthnRecoverAccount != nil {
			return createAndSendEmailUsingSupertokensService(appInfo, *input.WebauthnRecoverAccount)
		}
		return fmt.Errorf("should never come here")
	}
	return emaildelivery.EmailDeliveryInterface{
		SendEmail: &sendEmail,
	}
}

var emailRegex = regexp.MustCompile(
	`^(([^<>()\[\]\\.,;:\s@"]+(\.[^<>()\[\]\\.,;:\s@"]+)*)|(".+"))@` +
		`((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$`,
)

func defaultValidateEmail(email string, _ string, _ supertokens.UserContext) *string {
	if !emailRegex.MatchString(email) {
		msg := "Email is not valid"
		return &msg
	}
	return nil
}

func validateAndNormaliseUserInput(appInfo supertokens.NormalisedAppinfo, config *webauthnmodels.TypeInput) webauthnmodels.TypeNormalisedInput {
	inputConfig := webauthnmodels.TypeInput{}
	if config != nil {
		inputConfig = *config
	}

	typeNormalisedInput := makeTypeNormalisedInput(appInfo, inputConfig)

	getEmailDeliveryConfig := func() emaildelivery.TypeInputWithService {
		emailService := makeDefaultEmailService(appInfo)
		if inputConfig.EmailDelivery != nil && inputConfig.EmailDelivery.Service != nil {
			emailService = *inputConfig.EmailDelivery.Service
		}

		result := emaildelivery.TypeInputWithService{
			Service: emailService,
		}
		if inputConfig.EmailDelivery != nil && inputConfig.EmailDelivery.Override != nil {
			result.Override = inputConfig.EmailDelivery.Override
		}

		return result
	}

	typeNormalisedInput.GetEmailDeliveryConfig = getEmailDeliveryConfig

	if inputConfig.Override != nil {
		if inputConfig.Override.Functions != nil {
			typeNormalisedInput.Override.Functions = inputConfig.Override.Functions
		}
		if inputConfig.Override.APIs != nil {
			typeNormalisedInput.Override.APIs = inputConfig.Override.APIs
		}
	}

	return typeNormalisedInput
}

func makeTypeNormalisedInput(appInfo supertokens.NormalisedAppinfo, inputConfig webauthnmodels.TypeInput) webauthnmodels.TypeNormalisedInput {
	getRelyingPartyId := func(tenantId string, req *http.Request, userContext supertokens.UserContext) (string, error) {
		if inputConfig.GetRelyingPartyId != nil {
			return inputConfig.GetRelyingPartyId(tenantId, req, userContext)
		}

		apiDomainStr := appInfo.APIDomain.GetAsStringDangerous()
		parsed, err := url.Parse(apiDomainStr)
		if err != nil {
			return "", err
		}
		if parsed.Hostname() == "" {
			return "", fmt.Errorf("could not parse hostname from APIDomain: %s", apiDomainStr)
		}
		return parsed.Hostname(), nil
	}

	getRelyingPartyName := func(tenantId string, _ *http.Request, userContext supertokens.UserContext) (string, error) {
		if inputConfig.GetRelyingPartyName != nil {
			return inputConfig.GetRelyingPartyName(tenantId, userContext)
		}
		return appInfo.AppName, nil
	}

	getOrigin := func(tenantId string, req *http.Request, userContext supertokens.UserContext) (string, error) {
		if inputConfig.GetOrigin != nil {
			return inputConfig.GetOrigin(tenantId, req, userContext)
		}

		origin, err := appInfo.GetOrigin(req, userContext)
		if err != nil {
			return "", err
		}
		return origin.GetAsStringDangerous(), nil
	}

	validateEmailAddress := defaultValidateEmail
	if inputConfig.ValidateEmailAddress != nil {
		validateEmailAddress = inputConfig.ValidateEmailAddress
	}

	return webauthnmodels.TypeNormalisedInput{
		GetRelyingPartyId:    getRelyingPartyId,
		GetRelyingPartyName:  getRelyingPartyName,
		GetOrigin:            getOrigin,
		ValidateEmailAddress: validateEmailAddress,
		Override: webauthnmodels.OverrideStruct{
			Functions: func(originalImplementation webauthnmodels.RecipeInterface) webauthnmodels.RecipeInterface {
				return originalImplementation
			},
			APIs: func(originalImplementation webauthnmodels.APIInterface) webauthnmodels.APIInterface {
				return originalImplementation
			},
		},
	}
}

func GetRecoverAccountLink(appInfo supertokens.NormalisedAppinfo, token string, tenantId string, req *http.Request, userContext supertokens.UserContext) (string, error) {
	origin, err := appInfo.GetOrigin(req, userContext)
	if err != nil {
		return "", err
	}
	websiteBasePath := appInfo.WebsiteBasePath.GetAsStringDangerous()
	return fmt.Sprintf("%s%s/webauthn/recover?token=%s&tenantId=%s",
		origin.GetAsStringDangerous(),
		websiteBasePath,
		token,
		tenantId,
	), nil
}
