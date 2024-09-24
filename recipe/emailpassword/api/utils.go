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
	"fmt"
	"net/http"
	"strings"

	"github.com/supertokens/supertokens-golang/recipe/emailpassword/epmodels"
	"github.com/supertokens/supertokens-golang/recipe/emailpassword/errors"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func withValueAsString(emailValue interface{}, errorStr string) (string, error) {
	// Throw error if the value is not a string
	valueAsString, asStrOk := emailValue.(string)
	if !asStrOk {
		return "", fmt.Errorf(errorStr)
	}

	return valueAsString, nil
}

func validateFormFieldsOrThrowError(configFormFields []epmodels.NormalisedFormField, formFieldsRaw interface{}, tenantId string) ([]epmodels.TypeFormField, error) {
	if formFieldsRaw == nil {
		return nil, supertokens.BadInputError{
			Msg: "Missing input param: formFields",
		}
	}

	if _, ok := formFieldsRaw.([]interface{}); !ok {
		return nil, supertokens.BadInputError{
			Msg: "formFields must be an array",
		}
	}

	var formFields []epmodels.TypeFormField
	for _, rawFormField := range formFieldsRaw.([]interface{}) {

		if _, ok := rawFormField.(map[string]interface{}); !ok {
			return nil, supertokens.BadInputError{
				Msg: "formFields must be an array of objects containing id and value of type string",
			}
		}

		if rawFormField.(map[string]interface{})["id"] != nil {
			if _, ok := rawFormField.(map[string]interface{})["id"].(string); !ok {
				return nil, supertokens.BadInputError{
					Msg: "formFields must be an array of objects containing id and value of type string",
				}
			}
		}

		jsonformField, err := json.Marshal(rawFormField)
		if err != nil {
			return nil, err
		}
		var formField epmodels.TypeFormField
		err = json.Unmarshal(jsonformField, &formField)
		if err != nil {
			return nil, err
		}

		if formField.ID == "email" || formField.ID == "password" {
			valueAsString, parseErr := withValueAsString(formField.Value, fmt.Sprintf("%s value must be a string", formField.ID))
			if parseErr != nil {
				return nil, supertokens.BadInputError{
					Msg: parseErr.Error(),
				}
			}

			formFields = append(formFields, epmodels.TypeFormField{
				ID:    formField.ID,
				Value: strings.TrimSpace(valueAsString),
			})
		} else {
			formFields = append(formFields, epmodels.TypeFormField{
				ID:    formField.ID,
				Value: formField.Value,
			})
		}
	}

	return formFields, validateFormOrThrowError(configFormFields, formFields, tenantId)
}

func validateFormOrThrowError(configFormFields []epmodels.NormalisedFormField, inputs []epmodels.TypeFormField, tenantId string) error {
	var validationErrors []errors.ErrorPayload
	if len(configFormFields) < len(inputs) {
		return supertokens.BadInputError{
			Msg: "Are you sending too many formFields?",
		}
	}
	for _, field := range configFormFields {
		var input epmodels.TypeFormField
		for _, inputField := range inputs {
			if inputField.ID == field.ID {
				input = inputField
				break
			}
		}

		isValidInput := true
		if input.Value == nil {
			isValidInput = false
		} else {
			// If it is a string, it shouldn't be empty.
			valueAsStr, err := withValueAsString(input.Value, "")
			if err == nil && strings.TrimSpace(valueAsStr) == "" {
				isValidInput = false
			}
		}

		// If the field is not option and input is invalid, we should
		// throw a validation error.
		if !isValidInput && !field.Optional {
			validationErrors = append(validationErrors, errors.ErrorPayload{ID: field.ID, ErrorMsg: "Field is not optional"})
		}

		// If the input is invalid, we don't need to do anything
		// as execution will reach here if field is optional.
		if !isValidInput {
			continue
		}

		err := field.Validate(input.Value, tenantId)
		if err != nil {
			validationErrors = append(validationErrors, errors.ErrorPayload{
				ID:       field.ID,
				ErrorMsg: *err,
			})
		}
	}
	if len(validationErrors) != 0 {
		return errors.FieldError{
			Msg:     "Error in input formFields",
			Payload: validationErrors,
		}
	}
	return nil
}

func GetPasswordResetLink(appInfo supertokens.NormalisedAppinfo, token string, tenantId string, request *http.Request, userContext supertokens.UserContext) (string, error) {
	websiteDomain, err := appInfo.GetOrigin(request, userContext)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf(
		"%s%s/reset-password?token=%s&tenantId=%s",
		websiteDomain.GetAsStringDangerous(),
		appInfo.WebsiteBasePath.GetAsStringDangerous(),
		token,
		tenantId,
	), nil
}
