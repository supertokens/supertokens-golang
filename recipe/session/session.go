/* Copyright (c) 2021, VRAI Labs and/or its affiliates. All rights reserved.
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

package session

import (
	defaultErrors "errors"
	"net/http"
	"reflect"

	"github.com/supertokens/supertokens-golang/recipe/session/errors"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

type SessionContainerInput struct {
	sessionHandle         string
	userID                string
	userDataInAccessToken map[string]interface{}
	res                   http.ResponseWriter
	accessToken           string
}

func makeSessionContainerInput(accessToken string, sessionHandle string, userID string, userDataInAccessToken map[string]interface{}, res http.ResponseWriter) SessionContainerInput {
	return SessionContainerInput{
		sessionHandle:         sessionHandle,
		userID:                userID,
		userDataInAccessToken: userDataInAccessToken,
		res:                   res,
		accessToken:           accessToken,
	}
}

func newSessionContainer(querier supertokens.Querier, config sessmodels.TypeNormalisedInput, session *SessionContainerInput) sessmodels.SessionContainer {

	return sessmodels.SessionContainer{
		RevokeSession: func(userContext supertokens.UserContext) error {
			success, err := revokeSessionHelper(querier, session.sessionHandle)
			if err != nil {
				return err
			}
			if success {
				clearSessionFromCookie(config, session.res)
			}
			return nil
		},

		GetSessionData: func(userContext supertokens.UserContext) (map[string]interface{}, error) {
			sessionInformation, err := getSessionInformationHelper(querier, session.sessionHandle)
			if err != nil {
				if defaultErrors.As(err, &errors.UnauthorizedError{}) {
					clearSessionFromCookie(config, session.res)
				}
				return nil, err
			}
			return sessionInformation.SessionData, nil
		},

		UpdateSessionData: func(newSessionData map[string]interface{}, userContext supertokens.UserContext) error {
			err := updateSessionDataHelper(querier, session.sessionHandle, newSessionData)
			if err != nil {
				if defaultErrors.As(err, &errors.UnauthorizedError{}) {
					clearSessionFromCookie(config, session.res)
				}
				return err
			}
			return nil
		},

		UpdateAccessTokenPayload: func(newAccessTokenPayload map[string]interface{}, userContext supertokens.UserContext) error {
			if newAccessTokenPayload == nil {
				newAccessTokenPayload = map[string]interface{}{}
			}

			resp, err := regenerateAccessTokenHelper(querier, session.sessionHandle, newAccessTokenPayload, session.accessToken)

			if err != nil {
				return err
			}

			session.userDataInAccessToken = resp.GetSessionResponse.Session.UserDataInAccessToken
			if !reflect.DeepEqual(resp.GetSessionResponse.AccessToken, sessmodels.CreateOrRefreshAPIResponseToken{}) {
				session.accessToken = resp.GetSessionResponse.AccessToken.Token
				setFrontTokenInHeaders(session.res, resp.GetSessionResponse.Session.UserID, resp.GetSessionResponse.AccessToken.Expiry, resp.GetSessionResponse.Session.UserDataInAccessToken)
				attachAccessTokenToCookie(config, session.res, resp.GetSessionResponse.AccessToken.Token, resp.GetSessionResponse.AccessToken.Expiry)
			}
			return nil
		},
		GetUserID: func() string {
			return session.userID
		},
		GetAccessTokenPayload: func() map[string]interface{} {
			return session.userDataInAccessToken
		},
		GetHandle: func() string {
			return session.sessionHandle
		},
		GetAccessToken: func() string {
			return session.accessToken
		},
		GetTimeCreated: func(userContext supertokens.UserContext) (uint64, error) {
			sessionInformation, err := getSessionInformationHelper(querier, session.sessionHandle)
			if err != nil {
				if defaultErrors.As(err, &errors.UnauthorizedError{}) {
					clearSessionFromCookie(config, session.res)
				}
				return 0, err
			}
			return sessionInformation.TimeCreated, nil
		},
		GetExpiry: func(userContext supertokens.UserContext) (uint64, error) {
			sessionInformation, err := getSessionInformationHelper(querier, session.sessionHandle)
			if err != nil {
				if defaultErrors.As(err, &errors.UnauthorizedError{}) {
					clearSessionFromCookie(config, session.res)
				}
				return 0, err
			}
			return sessionInformation.Expiry, nil
		},
	}
}
