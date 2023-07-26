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
	"github.com/supertokens/supertokens-golang/recipe/usermetadata"
	"github.com/supertokens/supertokens-golang/supertokens"
)

type userMetadataPutResponse struct {
	Status string `json:"string"`
}

type userMetaDataRequestBody struct {
	UserId *string `json:"userId"`
	Data   *string `json:"data"`
}

func UserMetaDataPut(apiInterface dashboardmodels.APIInterface, options dashboardmodels.APIOptions, userContext supertokens.UserContext) (userMetadataPutResponse, error) {
	body, err := supertokens.ReadFromRequest(options.Req)

	if err != nil {
		return userMetadataPutResponse{}, err
	}

	var readBody userMetaDataRequestBody
	err = json.Unmarshal(body, &readBody)
	if err != nil {
		return userMetadataPutResponse{}, err
	}

	if readBody.UserId == nil {
		return userMetadataPutResponse{}, supertokens.BadInputError{
			Msg: "Required parameter 'userId' is missing",
		}
	}

	if readBody.Data == nil {
		return userMetadataPutResponse{}, supertokens.BadInputError{
			Msg: "Required parameter 'data' is missing",
		}
	}

	_, instanceError := usermetadata.GetRecipeInstanceOrThrowError()

	// This is so that the API exists early if the recipe has not been initialised
	if instanceError != nil {
		return userMetadataPutResponse{}, instanceError
	}

	var parsedMetaData map[string]interface{}
	parseErr := json.Unmarshal([]byte(*readBody.Data), &parsedMetaData)

	if parseErr != nil {
		return userMetadataPutResponse{}, supertokens.BadInputError{
			Msg: "'data' must be a valid JSON body",
		}
	}

	/**
	 * This API is meant to set the user metadata of a user. We delete the existing data
	 * before updating it because we want to make sure that shallow merging does not result
	 * in the data being incorrect
	 *
	 * For example if the old data is {test: "test", test2: "test2"} and the user wants to delete
	 * test2 from the data simply calling updateUserMetadata with {test: "test"} would not remove
	 * test2 because of shallow merging.
	 *
	 * Removing first ensures that the final data is exactly what the user wanted it to be
	 */
	clearErr := usermetadata.ClearUserMetadataWithContext(*readBody.UserId, userContext)

	if clearErr != nil {
		return userMetadataPutResponse{}, clearErr
	}

	_, updateErr := usermetadata.UpdateUserMetadataWithContext(*readBody.UserId, parsedMetaData, userContext)

	if updateErr != nil {
		return userMetadataPutResponse{}, updateErr
	}

	return userMetadataPutResponse{
		Status: "OK",
	}, nil
}
