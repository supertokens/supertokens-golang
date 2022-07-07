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

	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

type SessionContainerInput struct {
	sessionHandle         string
	userID                string
	userDataInAccessToken map[string]interface{}
	res                   http.ResponseWriter
	accessToken           string
	recipeImpl            sessmodels.RecipeInterface
}

func makeSessionContainerInput(accessToken string, sessionHandle string, userID string, userDataInAccessToken map[string]interface{}, res http.ResponseWriter, recipeImpl sessmodels.RecipeInterface) SessionContainerInput {
	return SessionContainerInput{
		sessionHandle:         sessionHandle,
		userID:                userID,
		userDataInAccessToken: userDataInAccessToken,
		res:                   res,
		accessToken:           accessToken,
		recipeImpl:            recipeImpl,
	}
}

func newSessionContainer(config sessmodels.TypeNormalisedInput, session *SessionContainerInput) sessmodels.SessionContainer {

	revokeSessionWithContext := func(userContext supertokens.UserContext) error {
		success, err := (*session.recipeImpl.RevokeSession)(session.sessionHandle, userContext)
		if err != nil {
			return err
		}
		if success {
			clearSessionFromCookie(config, session.res)
		}
		return nil
	}

	getSessionDataWithContext := func(userContext supertokens.UserContext) (map[string]interface{}, error) {
		sessionInformation, err := (*session.recipeImpl.GetSessionInformation)(session.sessionHandle, userContext)
		if err != nil {
			return nil, err
		}
		if sessionInformation == nil {
			clearSessionFromCookie(config, session.res)
			return nil, defaultErrors.New("session does not exist anymore")
		}
		return sessionInformation.SessionData, nil
	}

	updateSessionDataWithContext := func(newSessionData map[string]interface{}, userContext supertokens.UserContext) error {
		updated, err := (*session.recipeImpl.UpdateSessionData)(session.sessionHandle, newSessionData, userContext)
		if err != nil {
			return err
		}
		if !updated {
			clearSessionFromCookie(config, session.res)
			return defaultErrors.New("session does not exist anymore")
		}
		return nil
	}

	updateAccessTokenPayloadWithContext := func(newAccessTokenPayload map[string]interface{}, userContext supertokens.UserContext) error {
		if newAccessTokenPayload == nil {
			newAccessTokenPayload = map[string]interface{}{}
		}

		resp, err := (*session.recipeImpl.RegenerateAccessToken)(session.accessToken, &newAccessTokenPayload, userContext)

		if err != nil {
			return err
		}

		if resp == nil {
			clearSessionFromCookie(config, session.res)
			return defaultErrors.New("session does not exist anymore")
		}

		session.userDataInAccessToken = resp.Session.UserDataInAccessToken

		if !reflect.DeepEqual(resp.AccessToken, sessmodels.CreateOrRefreshAPIResponseToken{}) {
			session.accessToken = resp.AccessToken.Token
			setFrontTokenInHeaders(session.res, resp.Session.UserID, resp.AccessToken.Expiry, resp.Session.UserDataInAccessToken)
			attachAccessTokenToCookie(config, session.res, resp.AccessToken.Token, resp.AccessToken.Expiry)
		}
		return nil
	}

	getTimeCreatedWithContext := func(userContext supertokens.UserContext) (uint64, error) {
		sessionInformation, err := (*session.recipeImpl.GetSessionInformation)(session.sessionHandle, userContext)
		if err != nil {
			return 0, err
		}
		if sessionInformation == nil {
			clearSessionFromCookie(config, session.res)
			return 0, defaultErrors.New("session does not exist anymore")
		}
		return sessionInformation.TimeCreated, nil
	}

	getExpiryWithContext := func(userContext supertokens.UserContext) (uint64, error) {
		sessionInformation, err := (*session.recipeImpl.GetSessionInformation)(session.sessionHandle, userContext)
		if err != nil {
			return 0, err
		}
		if sessionInformation == nil {
			clearSessionFromCookie(config, session.res)
			return 0, defaultErrors.New("session does not exist anymore")
		}
		return sessionInformation.Expiry, nil
	}

	getUserIDWithContext := func(userContext supertokens.UserContext) string {
		return session.userID
	}
	getAccessTokenPayloadWithContext := func(userContext supertokens.UserContext) map[string]interface{} {
		return session.userDataInAccessToken
	}

	getHandleWithContext := func(userContext supertokens.UserContext) string {
		return session.sessionHandle
	}
	getAccessTokenWithContext := func(userContext supertokens.UserContext) string {
		return session.accessToken
	}

	return sessmodels.SessionContainer{
		RevokeSessionWithContext:            revokeSessionWithContext,
		GetSessionDataWithContext:           getSessionDataWithContext,
		UpdateSessionDataWithContext:        updateSessionDataWithContext,
		UpdateAccessTokenPayloadWithContext: updateAccessTokenPayloadWithContext,
		GetUserIDWithContext:                getUserIDWithContext,
		GetAccessTokenPayloadWithContext:    getAccessTokenPayloadWithContext,
		GetHandleWithContext:                getHandleWithContext,
		GetAccessTokenWithContext:           getAccessTokenWithContext,
		GetTimeCreatedWithContext:           getTimeCreatedWithContext,
		GetExpiryWithContext:                getExpiryWithContext,
		RevokeSession: func() error {
			return revokeSessionWithContext(&map[string]interface{}{})
		},
		GetSessionData: func() (map[string]interface{}, error) {
			return getSessionDataWithContext(&map[string]interface{}{})
		},
		UpdateSessionData: func(newSessionData map[string]interface{}) error {
			return updateSessionDataWithContext(newSessionData, &map[string]interface{}{})
		},
		UpdateAccessTokenPayload: func(newAccessTokenPayload map[string]interface{}) error {
			return updateAccessTokenPayloadWithContext(newAccessTokenPayload, &map[string]interface{}{})
		},
		GetUserID: func() string {
			return getUserIDWithContext(&map[string]interface{}{})
		},
		GetAccessTokenPayload: func() map[string]interface{} {
			return getAccessTokenPayloadWithContext(&map[string]interface{}{})
		},
		GetHandle: func() string {
			return getHandleWithContext(&map[string]interface{}{})
		},
		GetAccessToken: func() string {
			return getAccessTokenWithContext(&map[string]interface{}{})
		},
		GetTimeCreated: func() (uint64, error) {
			return getTimeCreatedWithContext(&map[string]interface{}{})
		},
		GetExpiry: func() (uint64, error) {
			return getExpiryWithContext(&map[string]interface{}{})
		},
	}
}
