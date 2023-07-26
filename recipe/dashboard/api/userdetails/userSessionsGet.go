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
	"sync"

	"github.com/supertokens/supertokens-golang/recipe/dashboard/dashboardmodels"
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/supertokens"
)

type SessionType struct {
	SessionDataInDatabase            interface{} `json:"sessionDataInDatabase"`
	CustomClaimsInAccessTokenPayload interface{} `json:"accessTokenPayload"`
	UserId                           string      `json:"userId"`
	Expiry                           uint64      `json:"expiry"`
	TimeCreated                      uint64      `json:"timeCreated"`
	SessionHandle                    string      `json:"sessionHandle"`
}

type userSessionsGetResponse struct {
	Status   string        `json:"status"`
	Sessions []SessionType `json:"sessions"`
}

func UserSessionsGet(apiInterface dashboardmodels.APIInterface, options dashboardmodels.APIOptions, userContext supertokens.UserContext) (userSessionsGetResponse, error) {
	req := options.Req
	userId := req.URL.Query().Get("userId")

	if userId == "" {
		return userSessionsGetResponse{}, supertokens.BadInputError{
			Msg: "Missing required parameter 'userId'",
		}
	}

	response, err := session.GetAllSessionHandlesForUserWithContext(userId, userContext)

	if err != nil {
		return userSessionsGetResponse{}, err
	}

	sessions := []SessionType{}

	var processingGroup sync.WaitGroup
	processingGroup.Add(len(response))
	var errInBackground error

	for i, sessionHandle := range response {
		if errInBackground != nil {
			return userSessionsGetResponse{}, errInBackground
		}

		go func(i int, handle string) {
			sessionResponse, sessionError := session.GetSessionInformationWithContext(handle, userContext)

			if sessionError != nil {
				errInBackground = sessionError
				return
			}

			if sessionResponse != nil {
				sessions = append(sessions, SessionType{
					SessionDataInDatabase:            sessionResponse.SessionDataInDatabase,
					CustomClaimsInAccessTokenPayload: sessionResponse.CustomClaimsInAccessTokenPayload,
					UserId:                           sessionResponse.UserId,
					Expiry:                           sessionResponse.Expiry,
					TimeCreated:                      sessionResponse.TimeCreated,
					SessionHandle:                    sessionResponse.SessionHandle,
				})
			}

			defer processingGroup.Done()
		}(i, sessionHandle)
	}

	if errInBackground != nil {
		return userSessionsGetResponse{}, errInBackground
	}

	processingGroup.Wait()

	return userSessionsGetResponse{
		Status:   "OK",
		Sessions: sessions,
	}, nil
}
