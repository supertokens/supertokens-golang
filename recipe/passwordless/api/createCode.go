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
	"io/ioutil"
	"reflect"
	"strings"

	"github.com/nyaruka/phonenumbers"
	"github.com/supertokens/supertokens-golang/recipe/passwordless/plessmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func CreateCode(apiImplementation plessmodels.APIInterface, options plessmodels.APIOptions) error {
	if apiImplementation.CreateCodePOST == nil || (*apiImplementation.CreateCodePOST) == nil {
		options.OtherHandler(options.Res, options.Req)
		return nil
	}

	body, err := ioutil.ReadAll(options.Req.Body)
	if err != nil {
		return err
	}
	var readBody map[string]interface{}
	err = json.Unmarshal(body, &readBody)
	if err != nil {
		return err
	}

	email, okEmail := readBody["email"]
	phoneNumber, okPhoneNumber := readBody["phoneNumber"]

	if (!okEmail && !okPhoneNumber) || (okEmail && okPhoneNumber) {
		return supertokens.BadInputError{Msg: "Please provide exactly one of email or phoneNumber"}
	}

	if okEmail && reflect.ValueOf(email).Kind() != reflect.String {
		return supertokens.BadInputError{Msg: "Please make sure that email is a string"}
	}

	if okPhoneNumber && reflect.ValueOf(phoneNumber).Kind() != reflect.String {
		return supertokens.BadInputError{Msg: "Please make sure that phoneNumber is a string"}
	}

	if !okEmail && options.Config.ContactMethodEmail.Enabled {
		return supertokens.BadInputError{Msg: "Please provide an email since you enabled ContactMethodEmail"}
	}

	if !okPhoneNumber && options.Config.ContactMethodPhone.Enabled {
		return supertokens.BadInputError{Msg: "Please provide a phoneNumber since you have enabled ContactMethodPhone"}
	}

	if okEmail {
		// normalize and validate email
		email = strings.TrimSpace(email.(string))
		validateErr := options.Config.ContactMethodEmail.ValidateEmailAddress(email)
		if validateErr != nil {
			return supertokens.Send200Response(options.Res, map[string]interface{}{
				"status":  "GENERAL_ERROR",
				"message": *validateErr,
			})
		}
	}

	if okPhoneNumber {
		validateErr := options.Config.ContactMethodPhone.ValidatePhoneNumber(phoneNumber)
		if validateErr != nil {
			return supertokens.Send200Response(options.Res, map[string]interface{}{
				"status":  "GENERAL_ERROR",
				"message": *validateErr,
			})
		}

		parsedPhoneNumber, err := phonenumbers.Parse(phoneNumber.(string), "")
		if err != nil {
			// this can come here if the user has provided their own impl of ValidatePhoneNumber and
			// the phone number is valid according to their impl, but not according to the phonenumbers lib.
			phoneNumber = strings.TrimSpace(phoneNumber.(string))
		} else {
			phoneNumber = phonenumbers.Format(parsedPhoneNumber, phonenumbers.INTERNATIONAL)
			fmt.Println(phoneNumber)
		}
	}

	var emailStrPointer *string
	var phoneNumberStrPointer *string
	if okEmail {
		t := email.(string)
		emailStrPointer = &t
	}
	if okPhoneNumber {
		t := phoneNumber.(string)
		phoneNumberStrPointer = &t
	}

	response, err := (*apiImplementation.CreateCodePOST)(emailStrPointer, phoneNumberStrPointer, options, &map[string]interface{}{})
	if err != nil {
		return err
	}

	var result map[string]interface{}

	if response.OK != nil {
		result = map[string]interface{}{
			"status":           "OK",
			"deviceId":         response.OK.DeviceID,
			"preAuthSessionId": response.OK.PreAuthSessionID,
			"flowType":         response.OK.FlowType,
		}
	} else {
		result = map[string]interface{}{
			"status":  "GENERAL_ERROR",
			"message": response.GeneralError.Message,
		}
	}

	return supertokens.Send200Response(options.Res, result)
}
