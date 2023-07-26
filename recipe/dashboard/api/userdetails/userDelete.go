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
	"github.com/supertokens/supertokens-golang/recipe/dashboard/dashboardmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

type userDeleteResponse struct {
	Status string `json:"status"`
}

func UserDelete(apiInterface dashboardmodels.APIInterface, options dashboardmodels.APIOptions, userContext supertokens.UserContext) (userDeleteResponse, error) {
	req := options.Req
	userId := req.URL.Query().Get("userId")

	if userId == "" {
		return userDeleteResponse{}, supertokens.BadInputError{
			Msg: "Missing required parameter 'userId'",
		}
	}

	deleteError := supertokens.DeleteUser(userId)

	if deleteError != nil {
		return userDeleteResponse{}, deleteError
	}

	return userDeleteResponse{
		Status: "OK",
	}, nil
}
