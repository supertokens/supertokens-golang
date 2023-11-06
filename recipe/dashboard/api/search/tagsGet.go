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

package search

import (
	"github.com/supertokens/supertokens-golang/recipe/dashboard/dashboardmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

type searchTagsResponse struct {
	Status string        `json:"status"`
	Tags   []interface{} `json:"tags"`
}

func SearchTagsGet(apiImplementation dashboardmodels.APIInterface, tenantId string, options dashboardmodels.APIOptions, userContext supertokens.UserContext) (searchTagsResponse, error) {
	querier, querierErr := supertokens.GetNewQuerierInstanceOrThrowError("dashboard")

	if querierErr != nil {
		return searchTagsResponse{}, querierErr
	}

	apiResponse, apiErr := querier.SendGetRequest("/user/search/tags", nil, userContext)
	if apiErr != nil {
		return searchTagsResponse{}, apiErr
	}

	return searchTagsResponse{
		Status: "OK",
		Tags:   apiResponse["tags"].([]interface{}),
	}, nil
}
