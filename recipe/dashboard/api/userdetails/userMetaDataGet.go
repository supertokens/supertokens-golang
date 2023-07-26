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
	"github.com/supertokens/supertokens-golang/recipe/usermetadata"
	"github.com/supertokens/supertokens-golang/supertokens"
)

type userMetaDataGetResponse struct {
	Status string      `json:"status"`
	Data   interface{} `json:"data,omitempty"`
}

func UserMetaDataGet(apiInterface dashboardmodels.APIInterface, options dashboardmodels.APIOptions, userContext supertokens.UserContext) (userMetaDataGetResponse, error) {
	req := options.Req
	userId := req.URL.Query().Get("userId")

	if userId == "" {
		return userMetaDataGetResponse{}, supertokens.BadInputError{
			Msg: "Missing required parameter 'userId'",
		}
	}

	_, instanceError := usermetadata.GetRecipeInstanceOrThrowError()

	if instanceError != nil {
		return userMetaDataGetResponse{
			Status: "FEATURE_NOT_ENABLED_ERROR",
		}, nil
	}

	metadata, err := usermetadata.GetUserMetadata(userId, userContext)

	if err != nil {
		return userMetaDataGetResponse{}, err
	}

	return userMetaDataGetResponse{
		Status: "OK",
		Data:   metadata,
	}, nil
}
