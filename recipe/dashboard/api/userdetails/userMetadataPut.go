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
	UserId string
	Data string
}

func UserMetaDataPut(apiInterface dashboardmodels.APIInterface, options dashboardmodels.APIOptions)(userMetadataPutResponse, error) {
	body, err := supertokens.ReadFromRequest(options.Req)

	if err != nil {
		return userMetadataPutResponse{}, err
	}

	var readBody userMetaDataRequestBody
	err = json.Unmarshal(body, &readBody)
	if err != nil {
		return userMetadataPutResponse{}, err
	}

	var parsedMetaData map[string]interface{}
	parseErr := json.Unmarshal([]byte(readBody.Data), &parsedMetaData)

	if parseErr != nil {
		return userMetadataPutResponse{}, supertokens.BadInputError{
			Msg: "'data' must be a valid JSON body",
		}
	}

	clearErr := usermetadata.ClearUserMetadata(readBody.UserId)

	if clearErr != nil {
		return userMetadataPutResponse{}, clearErr
	}

	_, updateErr := usermetadata.UpdateUserMetadata(readBody.UserId, parsedMetaData)

	if updateErr != nil {
		return userMetadataPutResponse{}, updateErr
	}

	return userMetadataPutResponse{
		Status: "OK",
	}, nil
}