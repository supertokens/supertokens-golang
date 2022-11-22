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
	"fmt"

	"github.com/supertokens/supertokens-golang/recipe/dashboard/dashboardmodels"
	"github.com/supertokens/supertokens-golang/recipe/emailverification"
	"github.com/supertokens/supertokens-golang/supertokens"
)

type userEmailVerifyGetResponse struct {
	Status string `json:"status"`
	IsVerified bool `json:"isVerified,omitempty"`
}

func UserEmailVerifyGet(apiImplementation dashboardmodels.APIInterface, options dashboardmodels.APIOptions)(userEmailVerifyGetResponse, error) {
	req := options.Req
	userId := req.URL.Query().Get("userId")

	if userId == "" {
		return userEmailVerifyGetResponse{}, supertokens.BadInputError {
			Msg: "Missing required parameter 'userId'",
		}
	}

	emailverificationInstance := emailverification.GetRecipeInstance()

	if emailverificationInstance == nil {
		return userEmailVerifyGetResponse{
			Status: "FEATURE_NOT_ENABLED_ERROR",
		}, nil
	}

	response, verificationError := emailverification.IsEmailVerified(userId, nil)

	if verificationError != nil {
		return userEmailVerifyGetResponse{}, verificationError
	}

	fmt.Println(response);

	return userEmailVerifyGetResponse{
		Status: "OK",
		IsVerified: response,
	}, nil
}