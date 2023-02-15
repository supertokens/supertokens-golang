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
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/supertokens"
)

type userSessionsPostResponse struct {
	Status string `json:"status"`
}

type userSessionsPostRequestBody struct {
	SessionHandles *[]string `json:"sessionHandles"`
}

func UserSessionsRevoke(apiInterface dashboardmodels.APIInterface, options dashboardmodels.APIOptions) (userSessionsPostResponse, error) {
	body, err := supertokens.ReadFromRequest(options.Req)

	if err != nil {
		return userSessionsPostResponse{}, err
	}

	var readBody userSessionsPostRequestBody
	err = json.Unmarshal(body, &readBody)
	if err != nil {
		return userSessionsPostResponse{}, err
	}

	sessionHandles := readBody.SessionHandles

	if sessionHandles == nil {
		return userSessionsPostResponse{}, supertokens.BadInputError{
			Msg: "Required parameter 'sessionHandles' is missing or has an invalid type",
		}
	}

	session.RevokeMultipleSessions(*sessionHandles, nil) // TODO pass tenant ID

	return userSessionsPostResponse{
		Status: "OK",
	}, nil
}
