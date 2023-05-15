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

package sessmodels

import (
	"github.com/supertokens/supertokens-golang/recipe/session/claims"
	"github.com/supertokens/supertokens-golang/supertokens"
)

type RecipeInterface struct {
	CreateNewSession            *func(userID string, accessTokenPayload map[string]interface{}, sessionDataInDatabase map[string]interface{}, disableAntiCsrf *bool, userContext supertokens.UserContext) (SessionContainer, error)
	GetSession                  *func(accessToken *string, antiCSRFToken *string, options *VerifySessionOptions, userContext supertokens.UserContext) (SessionContainer, error)
	RefreshSession              *func(refreshToken string, antiCSRFToken *string, disableAntiCSRF bool, userContext supertokens.UserContext) (SessionContainer, error)
	GetSessionInformation       *func(sessionHandle string, userContext supertokens.UserContext) (*SessionInformation, error)
	RevokeAllSessionsForUser    *func(userID string, userContext supertokens.UserContext) ([]string, error)
	GetAllSessionHandlesForUser *func(userID string, userContext supertokens.UserContext) ([]string, error)
	RevokeSession               *func(sessionHandle string, userContext supertokens.UserContext) (bool, error)
	RevokeMultipleSessions      *func(sessionHandles []string, userContext supertokens.UserContext) ([]string, error)
	UpdateSessionDataInDatabase *func(sessionHandle string, newSessionData map[string]interface{}, userContext supertokens.UserContext) (bool, error)
	MergeIntoAccessTokenPayload *func(sessionHandle string, accessTokenPayloadUpdate map[string]interface{}, userContext supertokens.UserContext) (bool, error)
	RegenerateAccessToken       *func(accessToken string, newAccessTokenPayload *map[string]interface{}, userContext supertokens.UserContext) (*RegenerateAccessTokenResponse, error)

	GetGlobalClaimValidators   *func(userId string, claimValidatorsAddedByOtherRecipes []claims.SessionClaimValidator, userContext supertokens.UserContext) ([]claims.SessionClaimValidator, error)
	ValidateClaims             *func(userId string, accessTokenPayload map[string]interface{}, claimValidators []claims.SessionClaimValidator, userContext supertokens.UserContext) (ValidateClaimsResult, error)
	ValidateClaimsInJWTPayload *func(userId string, jwtPayload map[string]interface{}, claimValidators []claims.SessionClaimValidator, userContext supertokens.UserContext) ([]claims.ClaimValidationError, error)
	FetchAndSetClaim           *func(sessionHandle string, claim *claims.TypeSessionClaim, userContext supertokens.UserContext) (bool, error)
	SetClaimValue              *func(sessionHandle string, claim *claims.TypeSessionClaim, value interface{}, userContext supertokens.UserContext) (bool, error)
	GetClaimValue              *func(sessionHandle string, claim *claims.TypeSessionClaim, userContext supertokens.UserContext) (GetClaimValueResult, error)
	RemoveClaim                *func(sessionHandle string, claim *claims.TypeSessionClaim, userContext supertokens.UserContext) (bool, error)
}
