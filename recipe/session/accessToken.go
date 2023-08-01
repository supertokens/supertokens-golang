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
	"errors"
	"fmt"
	"strings"

	"github.com/MicahParks/keyfunc"
	"github.com/golang-jwt/jwt/v4"
	sterrors "github.com/supertokens/supertokens-golang/recipe/session/errors"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

type AccessTokenInfoStruct struct {
	SessionHandle           string
	UserID                  string
	RefreshTokenHash1       string
	ParentRefreshTokenHash1 *string
	UserData                map[string]interface{}
	AntiCsrfToken           *string
	ExpiryTime              uint64
	TimeCreated             uint64
	TenantId                string
}

func GetInfoFromAccessToken(jwtInfo sessmodels.ParsedJWTInfo, jwks *keyfunc.JWKS, doAntiCsrfCheck bool) (*AccessTokenInfoStruct, error) {
	var payload map[string]interface{}

	if jwtInfo.Version >= 3 {
		parsedToken, parseError := jwt.Parse(jwtInfo.RawTokenString, jwks.Keyfunc)
		if parseError != nil {
			supertokens.LogDebugMessage(fmt.Sprintf("GetInfoFromAccessToken: Returning TryRefreshTokenError because access token parsing failed - %s", parseError))
			return nil, sterrors.TryRefreshTokenError{
				Msg: parseError.Error(),
			}
		}

		if parsedToken.Valid {
			claims, ok := parsedToken.Claims.(jwt.MapClaims)
			if !ok {
				supertokens.LogDebugMessage("GetInfoFromAccessToken: Returning TryRefreshTokenError because access token claims are invalid")
				return nil, sterrors.TryRefreshTokenError{
					Msg: "Invalid JWT claims",
				}
			}

			// Convert the claims to a key-value pair
			claimsMap := make(map[string]interface{})
			for key, value := range claims {
				claimsMap[key] = value
			}

			payload = claimsMap
		}
	} else {
		keys := []interface{}{}

		// Read only key returns all public keys that can be used for JWT verification
		for _, value := range jwks.ReadOnlyKeys() {
			keys = append(keys, value)
		}

		for _, key := range keys {
			parsedToken, parseErr := jwt.Parse(jwtInfo.RawTokenString, func(token *jwt.Token) (interface{}, error) {
				// The key returned here is used by Parse to verify the JWT
				return key, nil
			})

			if parseErr != nil && errors.Is(parseErr, jwt.ErrSignatureInvalid) {
				continue
			}

			if parseErr != nil {
				supertokens.LogDebugMessage(fmt.Sprintf("GetInfoFromAccessToken: Returning TryRefreshTokenError because access token parsing failed - %s", parseErr))
				return nil, sterrors.TryRefreshTokenError{
					Msg: parseErr.Error(),
				}
			}

			if parsedToken.Valid {
				claims, ok := parsedToken.Claims.(jwt.MapClaims)
				if !ok {
					supertokens.LogDebugMessage("GetInfoFromAccessToken: Returning TryRefreshTokenError because access token claims are invalid")
					return nil, sterrors.TryRefreshTokenError{
						Msg: "Invalid JWT claims",
					}
				}

				// Convert the claims to a key-value pair
				claimsMap := make(map[string]interface{})
				for key, value := range claims {
					claimsMap[key] = value
				}

				payload = claimsMap
				break
			}
		}
	}

	if payload == nil {
		supertokens.LogDebugMessage("GetInfoFromAccessToken: Returning TryRefreshTokenError because access token JWT has no payload")
		return nil, sterrors.TryRefreshTokenError{
			Msg: "Invalid JWT",
		}
	}

	err := ValidateAccessTokenStructure(payload, jwtInfo.Version)
	if err != nil {
		supertokens.LogDebugMessage("GetInfoFromAccessToken: Returning TryRefreshTokenError because ValidateAccessTokenStructure returned an error")
		return nil, sterrors.TryRefreshTokenError{
			Msg: err.Error(),
		}
	}

	// We can assume these as defined, since ValidateAccessTokenStructure checks this
	var userID string
	var expiryTime uint64
	var timeCreated uint64
	var userData map[string]interface{}
	if jwtInfo.Version >= 3 {
		userID = *sanitizeStringInput(payload["sub"])
		expiryTime = *sanitizeNumberInputAsUint64(payload["exp"]) * uint64(1000)
		timeCreated = *sanitizeNumberInputAsUint64(payload["iat"]) * uint64(1000)
		userData = payload
	} else {
		userID = *sanitizeStringInput(payload["userId"])
		expiryTime = *sanitizeNumberInputAsUint64(payload["expiryTime"])
		timeCreated = *sanitizeNumberInputAsUint64(payload["timeCreated"])
		userData = payload["userData"].(map[string]interface{})
	}

	sessionHandle := sanitizeStringInput(payload["sessionHandle"])
	refreshTokenHash1 := sanitizeStringInput(payload["refreshTokenHash1"])
	parentRefreshTokenHash1 := sanitizeStringInput(payload["parentRefreshTokenHash1"])
	antiCsrfToken := sanitizeStringInput(payload["antiCsrfToken"])

	tenantId := supertokens.DefaultTenantId
	if jwtInfo.Version >= 4 {
		tenantId = *sanitizeStringInput(payload["tId"])
	}

	if antiCsrfToken == nil && doAntiCsrfCheck {
		supertokens.LogDebugMessage("GetInfoFromAccessToken: Returning TryRefreshTokenError because access does not contain the anti-csrf token.")
		return nil, sterrors.TryRefreshTokenError{
			Msg: "Access token does not contain the anti-csrf token.",
		}
	}

	if expiryTime < GetCurrTimeInMS() {
		supertokens.LogDebugMessage("GetInfoFromAccessToken: Returning TryRefreshTokenError because access is expired")
		return nil, sterrors.TryRefreshTokenError{
			Msg: "Access token expired",
		}
	}

	return &AccessTokenInfoStruct{
		SessionHandle:           *sessionHandle,
		UserID:                  userID,
		RefreshTokenHash1:       *refreshTokenHash1,
		ParentRefreshTokenHash1: parentRefreshTokenHash1,
		UserData:                userData,
		AntiCsrfToken:           antiCsrfToken,
		ExpiryTime:              expiryTime,
		TimeCreated:             timeCreated,
		TenantId:                tenantId,
	}, nil
}

func ValidateAccessTokenStructure(payload map[string]interface{}, version int) error {
	err := errors.New("Access token does not contain all the information. Maybe the structure has changed?")

	if version >= 3 {
		supertokens.LogDebugMessage("ValidateAccessTokenStructure: Access token is using version >= 3")
		if _, ok := payload["sessionHandle"].(string); !ok {
			supertokens.LogDebugMessage("ValidateAccessTokenStructure: sessionHandle not found in JWT payload")
			return err
		}
		if _, ok := payload["sub"].(string); !ok {
			supertokens.LogDebugMessage("ValidateAccessTokenStructure: sub claim not found in JWT payload")
			return err
		}
		if _, ok := payload["refreshTokenHash1"].(string); !ok {
			supertokens.LogDebugMessage("ValidateAccessTokenStructure: refreshTokenHash1 not found in JWT payload")
			return err
		}
		if _, ok := payload["exp"].(float64); !ok {
			supertokens.LogDebugMessage("ValidateAccessTokenStructure: exp claim not found in JWT payload")
			return err
		}
		if _, ok := payload["iat"].(float64); !ok {
			supertokens.LogDebugMessage("ValidateAccessTokenStructure: iat claim not found in JWT payload")
			return err
		}
		if version >= 4 {
			if _, ok := payload["tId"].(string); !ok {
				supertokens.LogDebugMessage("ValidateAccessTokenStructure: tId claim not found in JWT payload")
				return err
			}
		}
	} else {
		supertokens.LogDebugMessage("ValidateAccessTokenStructure: Access token is using version < 3")
		if _, ok := payload["sessionHandle"].(string); !ok {
			supertokens.LogDebugMessage("ValidateAccessTokenStructure: sessionHandle not found in JWT payload")
			return err
		}
		if _, ok := payload["userId"].(string); !ok {
			supertokens.LogDebugMessage("ValidateAccessTokenStructure: userId not found in JWT payload")
			return err
		}
		if _, ok := payload["refreshTokenHash1"].(string); !ok {
			supertokens.LogDebugMessage("ValidateAccessTokenStructure: refreshTokenHash1 not found in JWT payload")
			return err
		}
		if payload["userData"] == nil {
			supertokens.LogDebugMessage("ValidateAccessTokenStructure: userData not found in JWT payload")
			return err
		}
		if _, ok := payload["userData"].(map[string]interface{}); !ok {
			supertokens.LogDebugMessage("ValidateAccessTokenStructure: userData is invalid in JWT payload")
			return err
		}
		if _, ok := payload["expiryTime"].(float64); !ok {
			supertokens.LogDebugMessage("ValidateAccessTokenStructure: expiryTime not found in JWT payload")
			return err
		}
		if _, ok := payload["timeCreated"].(float64); !ok {
			supertokens.LogDebugMessage("ValidateAccessTokenStructure: timeCreated not found in JWT payload")
			return err
		}
	}

	return nil
}

func sanitizeStringInput(field interface{}) *string {
	if field != nil {
		str, ok := field.(string)
		if ok {
			temp := strings.TrimSpace(str)
			return &temp
		}
	}
	return nil
}

func sanitizeNumberInputAsUint64(field interface{}) *uint64 {
	if field != nil {
		num, ok := field.(float64)
		if ok {
			temp := uint64(num)
			return &temp
		}
	}
	return nil
}
