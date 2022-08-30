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

package sessionwithjwt

import (
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/supertokens/supertokens-golang/recipe/openid/openidmodels"
	"github.com/supertokens/supertokens-golang/recipe/session/claims"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func newSessionWithJWTContainer(originalSessionClass sessmodels.SessionContainer, openidRecipeImplementation openidmodels.RecipeInterface) sessmodels.SessionContainer {

	updateAccessTokenPayloadWithContext := func(newAccessTokenPayload map[string]interface{}, userContext supertokens.UserContext) error {
		if newAccessTokenPayload == nil {
			newAccessTokenPayload = map[string]interface{}{}
		}
		accessTokenPayload := originalSessionClass.GetAccessTokenPayloadWithContext(userContext)
		jwtPropertyName, ok := accessTokenPayload[ACCESS_TOKEN_PAYLOAD_JWT_PROPERTY_NAME_KEY]

		if !ok {
			return originalSessionClass.UpdateAccessTokenPayloadWithContext(newAccessTokenPayload, userContext)
		}

		existingJWT := accessTokenPayload[jwtPropertyName.(string)].(string)

		currentTimeInSeconds := uint64(time.Now().UnixNano() / 1000000000) // time in seconds

		claims := jwt.MapClaims{}
		decodedPayload := map[string]interface{}{}
		_, _, err := new(jwt.Parser).ParseUnverified(existingJWT, claims)
		if err != nil {
			return err
		}
		for key, val := range claims {
			decodedPayload[key] = val
		}

		jwtExpiry := uint64(decodedPayload["exp"].(float64)) - currentTimeInSeconds

		if jwtExpiry <= 0 {
			// it can come here if someone calls this function well after
			// the access token and the jwt payload have expired (which can happen if an API takes a VERY long time). In this case, we still want the jwt payload to update, but the resulting JWT should
			// not be alive for too long (since it's expired already). So we set it to
			// 1 second lifetime.
			jwtExpiry = 1
		}

		newAccessTokenPayload, err = addJWTToAccessTokenPayload(newAccessTokenPayload, jwtExpiry, originalSessionClass.GetUserIDWithContext(userContext), jwtPropertyName.(string), openidRecipeImplementation, userContext)
		if err != nil {
			return err
		}

		return originalSessionClass.UpdateAccessTokenPayloadWithContext(newAccessTokenPayload, userContext)
	}

	mergeIntoAccessTokenPayloadWithContext := func(accessTokenPayloadUpdate map[string]interface{}, userContext supertokens.UserContext) error {
		return originalSessionClass.MergeIntoAccessTokenPayloadWithContext(accessTokenPayloadUpdate, userContext)
	}

	assertClaimsWithContext := func(claimValidators []claims.SessionClaimValidator, userContext supertokens.UserContext) error {
		return originalSessionClass.AssertClaimsWithContext(claimValidators, userContext)
	}

	fetchAndSetClaimWithContext := func(claim claims.TypeSessionClaim, userContext supertokens.UserContext) error {
		return originalSessionClass.FetchAndSetClaimWithContext(claim, userContext)
	}

	setClaimValueWithContext := func(claim claims.TypeSessionClaim, value interface{}, userContext supertokens.UserContext) error {
		return originalSessionClass.SetClaimValueWithContext(claim, value, userContext)
	}

	getClaimValueWithContext := func(claim claims.TypeSessionClaim, userContext supertokens.UserContext) interface{} {
		return originalSessionClass.GetClaimValueWithContext(claim, userContext)
	}

	removeClaimWithContext := func(claim claims.TypeSessionClaim, userContext supertokens.UserContext) error {
		return originalSessionClass.RemoveClaimWithContext(claim, userContext)
	}

	return sessmodels.SessionContainer{
		RevokeSessionWithContext:         originalSessionClass.RevokeSessionWithContext,
		GetSessionDataWithContext:        originalSessionClass.GetSessionDataWithContext,
		UpdateSessionDataWithContext:     originalSessionClass.UpdateSessionDataWithContext,
		GetUserIDWithContext:             originalSessionClass.GetUserIDWithContext,
		GetAccessTokenPayloadWithContext: originalSessionClass.GetAccessTokenPayloadWithContext,
		GetHandleWithContext:             originalSessionClass.GetHandleWithContext,
		GetAccessTokenWithContext:        originalSessionClass.GetAccessTokenWithContext,
		GetTimeCreatedWithContext:        originalSessionClass.GetTimeCreatedWithContext,
		GetExpiryWithContext:             originalSessionClass.GetExpiryWithContext,
		RevokeSession:                    originalSessionClass.RevokeSession,
		GetSessionData:                   originalSessionClass.GetSessionData,
		UpdateSessionData:                originalSessionClass.UpdateSessionData,
		GetUserID:                        originalSessionClass.GetUserID,
		GetAccessTokenPayload:            originalSessionClass.GetAccessTokenPayload,
		GetHandle:                        originalSessionClass.GetHandle,
		GetAccessToken:                   originalSessionClass.GetAccessToken,
		GetTimeCreated:                   originalSessionClass.GetTimeCreated,
		GetExpiry:                        originalSessionClass.GetExpiry,

		UpdateAccessTokenPayloadWithContext: updateAccessTokenPayloadWithContext,
		UpdateAccessTokenPayload: func(newAccessTokenPayload map[string]interface{}) error {
			return updateAccessTokenPayloadWithContext(newAccessTokenPayload, &map[string]interface{}{})
		},

		MergeIntoAccessTokenPayloadWithContext: mergeIntoAccessTokenPayloadWithContext,
		MergeIntoAccessTokenPayload: func(accessTokenPayloadUpdate map[string]interface{}) error {
			return mergeIntoAccessTokenPayloadWithContext(accessTokenPayloadUpdate, &map[string]interface{}{})
		},

		AssertClaimsWithContext:     assertClaimsWithContext,
		FetchAndSetClaimWithContext: fetchAndSetClaimWithContext,
		SetClaimValueWithContext:    setClaimValueWithContext,
		GetClaimValueWithContext:    getClaimValueWithContext,
		RemoveClaimWithContext:      removeClaimWithContext,

		AssertClaims: func(claimValidators []claims.SessionClaimValidator) error {
			return assertClaimsWithContext(claimValidators, &map[string]interface{}{})
		},
		FetchAndSetClaim: func(claim claims.TypeSessionClaim) error {
			return fetchAndSetClaimWithContext(claim, &map[string]interface{}{})
		},
		SetClaimValue: func(claim claims.TypeSessionClaim, value interface{}) error {
			return setClaimValueWithContext(claim, value, &map[string]interface{}{})
		},
		GetClaimValue: func(claim claims.TypeSessionClaim) interface{} {
			return getClaimValueWithContext(claim, &map[string]interface{}{})
		},
		RemoveClaim: func(claim claims.TypeSessionClaim) error {
			return removeClaimWithContext(claim, &map[string]interface{}{})
		},
	}
}
