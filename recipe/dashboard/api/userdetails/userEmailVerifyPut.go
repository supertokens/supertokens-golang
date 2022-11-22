/* Copyright (c) 2022, VRAI Labs and/or its affiliates. All rights reserved.
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

package userdetails

import (
	"encoding/json"
	"errors"

	"github.com/supertokens/supertokens-golang/recipe/dashboard/dashboardmodels"
	"github.com/supertokens/supertokens-golang/recipe/emailverification"
	"github.com/supertokens/supertokens-golang/supertokens"
)

type userEmailVerifyPutResponse struct {
	Status string `json:"status"`
}

type userEmailVerifyPutRequestBody struct {
	UserID string
	Verified bool
}

func UserEmailVerifyPut(apiInterface dashboardmodels.APIInterface, options dashboardmodels.APIOptions)(userEmailVerifyPutResponse, error) {
	body, err := supertokens.ReadFromRequest(options.Req)

	if err != nil {
		return userEmailVerifyPutResponse{}, err
	}

	var readBody userEmailVerifyPutRequestBody
	err = json.Unmarshal(body, &readBody)
	if err != nil {
		return userEmailVerifyPutResponse{}, err
	}

	if readBody.Verified {
		tokenResponse, tokenErr := emailverification.CreateEmailVerificationToken(readBody.UserID, nil)

		if tokenErr != nil {
			return userEmailVerifyPutResponse{}, tokenErr
		}

		if tokenResponse.EmailAlreadyVerifiedError != nil {
			return userEmailVerifyPutResponse{
				Status: "OK",
			}, nil
		}

		verifyResponse, verifyErr := emailverification.VerifyEmailUsingToken(tokenResponse.OK.Token)

		if verifyErr != nil {
			return userEmailVerifyPutResponse{}, verifyErr
		}

		if verifyResponse.EmailVerificationInvalidTokenError != nil {
			return userEmailVerifyPutResponse{}, errors.New("Should never come here")
		}
	} else {
		_, unverifyErr := emailverification.UnverifyEmail(readBody.UserID, nil)

		if unverifyErr != nil {
			return userEmailVerifyPutResponse{}, unverifyErr
		}
	}

	return userEmailVerifyPutResponse{
		Status: "OK",
	}, nil
}