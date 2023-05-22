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
	"fmt"
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
	accessToken           string
	frontToken            string
	refreshToken          *sessmodels.CreateOrRefreshAPIResponseToken
	recipeImpl            sessmodels.RecipeInterface
	antiCSRFToken         *string
	requestResponseInfo   *sessmodels.RequestResponseInfo
	accessTokenUpdated    bool
}

func makeSessionContainerInput(
	accessToken string,
	sessionHandle string,
	userID string,
	userDataInAccessToken map[string]interface{},
	recipeImpl sessmodels.RecipeInterface,
	frontToken string,
	antiCSRFToken *string,
	info *sessmodels.RequestResponseInfo,
	refreshToken *sessmodels.CreateOrRefreshAPIResponseToken,
	accessTokenUpdated bool,
) SessionContainerInput {
	return SessionContainerInput{
		sessionHandle:         sessionHandle,
		userID:                userID,
		userDataInAccessToken: userDataInAccessToken,
		accessToken:           accessToken,
		frontToken:            frontToken,
		refreshToken:          refreshToken,
		recipeImpl:            recipeImpl,
		antiCSRFToken:         antiCSRFToken,
		requestResponseInfo:   info,
		accessTokenUpdated:    accessTokenUpdated,
	}
}

func newSessionContainer(config sessmodels.TypeNormalisedInput, session *SessionContainerInput) sessmodels.SessionContainer {

	sessionContainer := &sessmodels.TypeSessionContainer{}
	sessionContainer.RevokeSessionWithContext = func(userContext supertokens.UserContext) error {
		_, err := (*session.recipeImpl.RevokeSession)(session.sessionHandle, userContext)
		if err != nil {
			return err
		}

		if session.requestResponseInfo != nil {
			// we do not check the output of calling revokeSession
			// before clearing the cookies because we are revoking the
			// current API request's session.
			// If we instead clear the cookies only when revokeSession
			// returns true, it can cause this kind of a bug:
			// https://github.com/supertokens/supertokens-node/issues/343
			ClearSession(config, session.requestResponseInfo.Res, session.requestResponseInfo.TokenTransferMethod)
		}
		return nil
	}

	sessionContainer.GetSessionDataInDatabaseWithContext = func(userContext supertokens.UserContext) (map[string]interface{}, error) {
		sessionInformation, err := (*session.recipeImpl.GetSessionInformation)(session.sessionHandle, userContext)
		if err != nil {
			return nil, err
		}
		if sessionInformation == nil {
			supertokens.LogDebugMessage("GetSessionDataInDatabaseWithContext: Returning UnauthorizedError because session does not exist anymore")
			return nil, errors.UnauthorizedError{Msg: "session does not exist anymore"}
		}
		return sessionInformation.SessionDataInDatabase, nil
	}

	sessionContainer.UpdateSessionDataInDatabaseWithContext = func(newSessionData map[string]interface{}, userContext supertokens.UserContext) error {
		updated, err := (*session.recipeImpl.UpdateSessionDataInDatabase)(session.sessionHandle, newSessionData, userContext)
		if err != nil {
			return err
		}
		if !updated {
			supertokens.LogDebugMessage("UpdateSessionDataInDatabaseWithContext: Returning UnauthorizedError because session does not exist anymore")
			return errors.UnauthorizedError{Msg: "session does not exist anymore"}
		}
		return nil
	}

	sessionContainer.GetTimeCreatedWithContext = func(userContext supertokens.UserContext) (uint64, error) {
		sessionInformation, err := (*session.recipeImpl.GetSessionInformation)(session.sessionHandle, userContext)
		if err != nil {
			return 0, err
		}
		if sessionInformation == nil {
			supertokens.LogDebugMessage("GetTimeCreatedWithContext: Returning UnauthorizedError because session does not exist anymore")
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
			supertokens.LogDebugMessage("GetExpiryWithContext: Returning UnauthorizedError because session does not exist anymore")
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

		for k, _ := range accessTokenPayload {
			if supertokens.DoesSliceContainString(k, protectedProps) {
				delete(accessTokenPayload, k)
			}
		}

		for k, v := range accessTokenPayloadUpdate {
			if v == nil {
				delete(accessTokenPayload, k)
			} else {
				accessTokenPayload[k] = v
			}
		}

		querier, err := supertokens.GetNewQuerierInstanceOrThrowError("")
		if err != nil {
			return err
		}

		response, err := regenerateAccessTokenHelper(*querier, &accessTokenPayload, sessionContainer.GetAccessToken())

		if err != nil {
			supertokens.LogDebugMessage(fmt.Sprintf("MergeIntoAccessTokenPayloadWithContext: Returning UnauthorizedError because we could not regenerate the session - %s", err))
			return errors.UnauthorizedError{
				Msg: errors.UnauthorizedErrorStr,
			}
		}

		if !reflect.DeepEqual(response.AccessToken, sessmodels.CreateOrRefreshAPIResponseToken{}) {
			responseToken, parseError := ParseJWTWithoutSignatureVerification(response.AccessToken.Token)
			if parseError != nil {
				return parseError
			}

			var payload map[string]interface{}
			if responseToken.Version < 3 {
				payload = response.Session.UserDataInAccessToken
			} else {
				payload = responseToken.Payload
			}

			session.userDataInAccessToken = payload
			session.accessToken = response.AccessToken.Token
			session.frontToken = BuildFrontToken(session.userID, response.AccessToken.Expiry, payload)
			session.accessTokenUpdated = true

			if session.requestResponseInfo != nil {
				setTokenErr := SetAccessTokenInResponse(config, session.requestResponseInfo.Res, session.accessToken, session.frontToken, session.requestResponseInfo.TokenTransferMethod)
				if setTokenErr != nil {
					return setTokenErr
				}
			}
		} else {
			// This case means that the access token has expired between the validation and this update
			// We can't update the access token on the FE, as it will need to call refresh anyway but we handle this as a successful update during this request.
			// the changes will be reflected on the FE after refresh is called
			currentAccessTokenPayload := sessionContainer.GetAccessTokenPayload()
			userDataInJWT := response.Session.UserDataInAccessToken

			userDataInAccessToken := map[string]interface{}{}

			for k, v := range currentAccessTokenPayload {
				userDataInAccessToken[k] = v
			}

			for k, v := range userDataInJWT {
				userDataInAccessToken[k] = v
			}

			session.userDataInAccessToken = userDataInAccessToken
		}

		return nil
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
			for _, protectedKey := range protectedProps {
				delete(validateClaimResponse.AccessTokenPayloadUpdate, protectedKey)
			}

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
	sessionContainer.GetSessionDataInDatabase = func() (map[string]interface{}, error) {
		return sessionContainer.GetSessionDataInDatabaseWithContext(&map[string]interface{}{})
	}
	sessionContainer.UpdateSessionDataInDatabase = func(newSessionData map[string]interface{}) error {
		return sessionContainer.UpdateSessionDataInDatabaseWithContext(newSessionData, &map[string]interface{}{})
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

	sessionContainer.GetAllSessionTokensDangerously = func() sessmodels.SessionTokens {
		var refreshToken *string = nil

		if session.refreshToken != nil {
			refreshToken = &session.refreshToken.Token
		}

		return sessmodels.SessionTokens{
			AccessToken:                   session.accessToken,
			RefreshToken:                  refreshToken,
			AntiCsrfToken:                 session.antiCSRFToken,
			FrontToken:                    session.frontToken,
			AccessAndFrontendTokenUpdated: session.accessTokenUpdated,
		}
	}

	sessionContainer.AttachToRequestResponse = func(info sessmodels.RequestResponseInfo) error {
		session.requestResponseInfo = &info

		if session.accessTokenUpdated {
			err := SetAccessTokenInResponse(config, info.Res, session.accessToken, session.frontToken, info.TokenTransferMethod)

			if err != nil {
				return err
			}

			if session.refreshToken != nil {
				err = setToken(config, info.Res, sessmodels.RefreshToken, session.refreshToken.Token, session.refreshToken.Expiry, info.TokenTransferMethod)

				if err != nil {
					return err
				}
			}

			if session.antiCSRFToken != nil {
				setAntiCsrfTokenInHeaders(info.Res, *session.antiCSRFToken)
			}
		}

		return nil
	}

	return sessionContainer
}
