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

	"github.com/supertokens/supertokens-golang/recipe/dashboard/dashboardmodels"
	"github.com/supertokens/supertokens-golang/recipe/emailverification"
	"github.com/supertokens/supertokens-golang/supertokens"
)

type userEmailVerifyTokenPost struct {
	Status string `json:"status"`
}

type userEmailverifyTokenPostRequestBody struct {
	UserId *string `json:"userId"`
}

func UserEmailVerifyTokenPost(apiInterface dashboardmodels.APIInterface, tenantId string, options dashboardmodels.APIOptions, userContext supertokens.UserContext) (userEmailVerifyTokenPost, error) {
	body, err := supertokens.ReadFromRequest(options.Req)

	if err != nil {
		return userEmailVerifyTokenPost{}, err
	}

	var readBody userEmailverifyTokenPostRequestBody
	err = json.Unmarshal(body, &readBody)
	if err != nil {
		return userEmailVerifyTokenPost{}, err
	}

	if readBody.UserId == nil {
		return userEmailVerifyTokenPost{}, supertokens.BadInputError{
			Msg: "Required parameter 'userId' is missing",
		}
	}

	resp, err := emailverification.SendEmailVerificationEmail(tenantId, *readBody.UserId, nil, userContext)
	if err != nil {
		return userEmailVerifyTokenPost{}, err
	}

	if resp.EmailAlreadyVerifiedError != nil {
		return userEmailVerifyTokenPost{
			Status: "EMAIL_ALREADY_VERIFIED_ERROR",
		}, nil
	}

	return userEmailVerifyTokenPost{
		Status: "OK",
	}, nil
}
