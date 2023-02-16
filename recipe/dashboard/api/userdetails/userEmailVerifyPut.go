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
	UserID   *string `json:"userId"`
	Verified *bool   `json:"verified"`
}

func UserEmailVerifyPut(apiInterface dashboardmodels.APIInterface, options dashboardmodels.APIOptions) (userEmailVerifyPutResponse, error) {
	body, err := supertokens.ReadFromRequest(options.Req)

	if err != nil {
		return userEmailVerifyPutResponse{}, err
	}

	var readBody userEmailVerifyPutRequestBody
	err = json.Unmarshal(body, &readBody)
	if err != nil {
		return userEmailVerifyPutResponse{}, err
	}

	if readBody.UserID == nil {
		return userEmailVerifyPutResponse{}, supertokens.BadInputError{
			Msg: "Required parameter 'userId' is missing",
		}
	}

	if readBody.Verified == nil {
		return userEmailVerifyPutResponse{}, supertokens.BadInputError{
			Msg: "Required parameter 'verified' is missing",
		}
	}

	if *readBody.Verified {
		tokenResponse, tokenErr := emailverification.CreateEmailVerificationToken(*readBody.UserID, nil, nil) // TODO tenantId

		if tokenErr != nil {
			return userEmailVerifyPutResponse{}, tokenErr
		}

		if tokenResponse.EmailAlreadyVerifiedError != nil {
			return userEmailVerifyPutResponse{
				Status: "OK",
			}, nil
		}

		verifyResponse, verifyErr := emailverification.VerifyEmailUsingToken(tokenResponse.OK.Token, nil) // TODO tenantId

		if verifyErr != nil {
			return userEmailVerifyPutResponse{}, verifyErr
		}

		// It should never come here because we generate the token immediately before this step
		if verifyResponse.EmailVerificationInvalidTokenError != nil {
			return userEmailVerifyPutResponse{}, errors.New("Should never come here")
		}
	} else {
		_, unverifyErr := emailverification.UnverifyEmail(*readBody.UserID, nil, nil) // TODO tenantId

		if unverifyErr != nil {
			return userEmailVerifyPutResponse{}, unverifyErr
		}
	}

	return userEmailVerifyPutResponse{
		Status: "OK",
	}, nil
}
