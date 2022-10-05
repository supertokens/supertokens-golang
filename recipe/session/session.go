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
	"net/http"
	"reflect"

	"github.com/supertokens/supertokens-golang/recipe/session/claims"
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

	sessionContainer := &sessmodels.TypeSessionContainer{}
	sessionContainer.RevokeSessionWithContext = func(userContext supertokens.UserContext) error {
		_, err := (*session.recipeImpl.RevokeSession)(session.sessionHandle, userContext)
		if err != nil {
			return err
		}
		clearSessionFromCookie(config, session.res)
		return nil
	}

	sessionContainer.GetSessionDataWithContext = func(userContext supertokens.UserContext) (map[string]interface{}, error) {
		sessionInformation, err := (*session.recipeImpl.GetSessionInformation)(session.sessionHandle, userContext)
		if err != nil {
			return nil, err
		}
		if sessionInformation == nil {
			return nil, errors.UnauthorizedError{Msg: "session does not exist anymore"}
		}
		return sessionInformation.SessionData, nil
	}

	sessionContainer.UpdateSessionDataWithContext = func(newSessionData map[string]interface{}, userContext supertokens.UserContext) error {
		updated, err := (*session.recipeImpl.UpdateSessionData)(session.sessionHandle, newSessionData, userContext)
		if err != nil {
			return err
		}
		if !updated {
			return errors.UnauthorizedError{Msg: "session does not exist anymore"}
		}
		return nil
	}

	sessionContainer.UpdateAccessTokenPayloadWithContext = func(newAccessTokenPayload map[string]interface{}, userContext supertokens.UserContext) error {
		if newAccessTokenPayload == nil {
			newAccessTokenPayload = map[string]interface{}{}
		}

		resp, err := (*session.recipeImpl.RegenerateAccessToken)(session.accessToken, &newAccessTokenPayload, userContext)

		if err != nil {
			return err
		}

		if resp == nil {
			return errors.UnauthorizedError{Msg: "session does not exist anymore"}
		}

		session.userDataInAccessToken = resp.Session.UserDataInAccessToken

		if !reflect.DeepEqual(resp.AccessToken, sessmodels.CreateOrRefreshAPIResponseToken{}) {
			session.accessToken = resp.AccessToken.Token
			setFrontTokenInHeaders(session.res, resp.Session.UserID, resp.AccessToken.Expiry, resp.Session.UserDataInAccessToken)
			attachAccessTokenToCookie(config, session.res, resp.AccessToken.Token, resp.AccessToken.Expiry)
		}
		return nil
	}

	sessionContainer.GetTimeCreatedWithContext = func(userContext supertokens.UserContext) (uint64, error) {
		sessionInformation, err := (*session.recipeImpl.GetSessionInformation)(session.sessionHandle, userContext)
		if err != nil {
			return 0, err
		}
		if sessionInformation == nil {
			return 0, errors.UnauthorizedError{Msg: "session does not exist anymore"}
		}
		return sessionInformation.TimeCreated, nil
	}

	sessionContainer.GetExpiryWithContext = func(userContext supertokens.UserContext) (uint64, error) {
		sessionInformation, err := (*session.recipeImpl.GetSessionInformation)(session.sessionHandle, userContext)
		if err != nil {
			return 0, err
		}
		if sessionInformation == nil {
			return 0, errors.UnauthorizedError{Msg: "session does not exist anymore"}
		}
		return sessionInformation.Expiry, nil
	}

	sessionContainer.GetUserIDWithContext = func(userContext supertokens.UserContext) string {
		return session.userID
	}
	sessionContainer.GetAccessTokenPayloadWithContext = func(userContext supertokens.UserContext) map[string]interface{} {
		return session.userDataInAccessToken
	}

	sessionContainer.MergeIntoAccessTokenPayloadWithContext = func(accessTokenPayloadUpdate map[string]interface{}, userContext supertokens.UserContext) error {
		accessTokenPayload := sessionContainer.GetAccessTokenPayloadWithContext(userContext)
		for k, v := range accessTokenPayloadUpdate {
			if v == nil {
				delete(accessTokenPayload, k)
			} else {
				accessTokenPayload[k] = v
			}
		}
		return sessionContainer.UpdateAccessTokenPayloadWithContext(accessTokenPayload, userContext)
	}

	sessionContainer.GetHandleWithContext = func(userContext supertokens.UserContext) string {
		return session.sessionHandle
	}
	sessionContainer.GetAccessTokenWithContext = func(userContext supertokens.UserContext) string {
		return session.accessToken
	}

	sessionContainer.AssertClaimsWithContext = func(claimValidators []claims.SessionClaimValidator, userContext supertokens.UserContext) error {
		validateClaimResponse, err := (*session.recipeImpl.ValidateClaims)(session.userID, sessionContainer.GetAccessTokenPayloadWithContext(userContext), claimValidators, userContext)
		if err != nil {
			return err
		}

		if validateClaimResponse.AccessTokenPayloadUpdate != nil {
			err := sessionContainer.MergeIntoAccessTokenPayloadWithContext(validateClaimResponse.AccessTokenPayloadUpdate, userContext)
			if err != nil {
				return err
			}
		}

		if len(validateClaimResponse.InvalidClaims) > 0 {
			return errors.InvalidClaimError{
				Msg:           "invalid claims",
				InvalidClaims: validateClaimResponse.InvalidClaims,
			}
		}

		return nil
	}

	sessionContainer.FetchAndSetClaimWithContext = func(claim *claims.TypeSessionClaim, userContext supertokens.UserContext) error {
		update, err := claim.Build(sessionContainer.GetUserIDWithContext(userContext), nil, userContext)
		if err != nil {
			return err
		}
		return sessionContainer.MergeIntoAccessTokenPayloadWithContext(update, userContext)
	}

	sessionContainer.SetClaimValueWithContext = func(claim *claims.TypeSessionClaim, value interface{}, userContext supertokens.UserContext) error {
		update := claim.AddToPayload_internal(map[string]interface{}{}, value, userContext)
		return sessionContainer.MergeIntoAccessTokenPayloadWithContext(update, userContext)
	}

	sessionContainer.GetClaimValueWithContext = func(claim *claims.TypeSessionClaim, userContext supertokens.UserContext) interface{} {
		return claim.GetValueFromPayload(sessionContainer.GetAccessTokenPayloadWithContext(userContext), userContext)
	}

	sessionContainer.RemoveClaimWithContext = func(claim *claims.TypeSessionClaim, userContext supertokens.UserContext) error {
		update := claim.RemoveFromPayloadByMerge_internal(map[string]interface{}{}, userContext)
		return sessionContainer.MergeIntoAccessTokenPayloadWithContext(update, userContext)
	}

	sessionContainer.RevokeSession = func() error {
		return sessionContainer.RevokeSessionWithContext(&map[string]interface{}{})
	}
	sessionContainer.GetSessionData = func() (map[string]interface{}, error) {
		return sessionContainer.GetSessionDataWithContext(&map[string]interface{}{})
	}
	sessionContainer.UpdateSessionData = func(newSessionData map[string]interface{}) error {
		return sessionContainer.UpdateSessionDataWithContext(newSessionData, &map[string]interface{}{})
	}
	sessionContainer.UpdateAccessTokenPayload = func(newAccessTokenPayload map[string]interface{}) error {
		return sessionContainer.UpdateAccessTokenPayloadWithContext(newAccessTokenPayload, &map[string]interface{}{})
	}
	sessionContainer.GetUserID = func() string {
		return sessionContainer.GetUserIDWithContext(&map[string]interface{}{})
	}
	sessionContainer.GetAccessTokenPayload = func() map[string]interface{} {
		return sessionContainer.GetAccessTokenPayloadWithContext(&map[string]interface{}{})
	}
	sessionContainer.GetHandle = func() string {
		return sessionContainer.GetHandleWithContext(&map[string]interface{}{})
	}
	sessionContainer.GetAccessToken = func() string {
		return sessionContainer.GetAccessTokenWithContext(&map[string]interface{}{})
	}
	sessionContainer.GetTimeCreated = func() (uint64, error) {
		return sessionContainer.GetTimeCreatedWithContext(&map[string]interface{}{})
	}
	sessionContainer.GetExpiry = func() (uint64, error) {
		return sessionContainer.GetExpiryWithContext(&map[string]interface{}{})
	}

	sessionContainer.MergeIntoAccessTokenPayload = func(accessTokenPayloadUpdate map[string]interface{}) error {
		return sessionContainer.MergeIntoAccessTokenPayloadWithContext(accessTokenPayloadUpdate, &map[string]interface{}{})
	}

	sessionContainer.AssertClaims = func(claimValidators []claims.SessionClaimValidator) error {
		return sessionContainer.AssertClaimsWithContext(claimValidators, &map[string]interface{}{})
	}
	sessionContainer.FetchAndSetClaim = func(claim *claims.TypeSessionClaim) error {
		return sessionContainer.FetchAndSetClaimWithContext(claim, &map[string]interface{}{})
	}
	sessionContainer.SetClaimValue = func(claim *claims.TypeSessionClaim, value interface{}) error {
		return sessionContainer.SetClaimValueWithContext(claim, value, &map[string]interface{}{})
	}
	sessionContainer.GetClaimValue = func(claim *claims.TypeSessionClaim) interface{} {
		return sessionContainer.GetClaimValueWithContext(claim, &map[string]interface{}{})
	}
	sessionContainer.RemoveClaim = func(claim *claims.TypeSessionClaim) error {
		return sessionContainer.RemoveClaimWithContext(claim, &map[string]interface{}{})
	}

	return sessionContainer
}
