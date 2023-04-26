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
	"github.com/MicahParks/keyfunc"
	"github.com/golang-jwt/jwt/v4"
	sterrors "github.com/supertokens/supertokens-golang/recipe/session/errors"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"strings"
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
}

func GetInfoFromAccessToken(jwtInfo sessmodels.ParsedJWTInfo, jwks keyfunc.JWKS, doAntiCsrfCheck bool) (*AccessTokenInfoStruct, error) {
	var payload map[string]interface{}

	if jwtInfo.Version >= 3 {
		parsedToken, parseError := jwt.Parse(jwtInfo.RawTokenString, jwks.Keyfunc)
		if parseError != nil {
			return nil, sterrors.TryRefreshTokenError{
				Msg: parseError.Error(),
			}
		}

		if parsedToken.Valid {
			claims, ok := parsedToken.Claims.(jwt.MapClaims)
			if !ok {
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
				return nil, sterrors.TryRefreshTokenError{
					Msg: parseErr.Error(),
				}
			}

			if parsedToken.Valid {
				claims, ok := parsedToken.Claims.(jwt.MapClaims)
				if !ok {
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
		}
	}

	if payload == nil {
		return nil, sterrors.TryRefreshTokenError{
			Msg: "Invalid JWT",
		}
	}

	err := ValidateAccessTokenStructure(payload, jwtInfo.Version)
	if err != nil {
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
		expiryTime = *sanitizeNumberInputAsUint64(payload["ExpiryTime"])
		timeCreated = *sanitizeNumberInputAsUint64(payload["TimeCreated"])
		userData = payload["UserData"].(map[string]interface{})
	}

	sessionHandle := sanitizeStringInput(payload["SessionHandle"])
	refreshTokenHash1 := sanitizeStringInput(payload["RefreshTokenHash1"])
	parentRefreshTokenHash1 := sanitizeStringInput(payload["ParentRefreshTokenHash1"])
	antiCsrfToken := sanitizeStringInput(payload["AntiCsrfToken"])

	if antiCsrfToken == nil && doAntiCsrfCheck {
		return nil, sterrors.TryRefreshTokenError{
			Msg: "Access token does not contain the anti-csrf token.",
		}
	}

	if expiryTime < GetCurrTimeInMS() {
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
	}, nil
}

func ValidateAccessTokenStructure(payload map[string]interface{}, version int) error {
	err := errors.New("Access token does not contain all the information. Maybe the structure has changed?")

	if version >= 3 {
		if _, ok := payload["SessionHandle"].(string); !ok {
			return err
		}
		if _, ok := payload["sub"].(string); !ok {
			return err
		}
		if _, ok := payload["RefreshTokenHash1"].(string); !ok {
			return err
		}
		if _, ok := payload["exp"].(float64); !ok {
			return err
		}
		if _, ok := payload["iat"].(float64); !ok {
			return err
		}
	} else {
		if _, ok := payload["SessionHandle"].(string); !ok {
			return err
		}
		if _, ok := payload["userId"].(string); !ok {
			return err
		}
		if _, ok := payload["RefreshTokenHash1"].(string); !ok {
			return err
		}
		if payload["UserData"] == nil {
			return err
		}
		if _, ok := payload["UserData"].(map[string]interface{}); !ok {
			return err
		}
		if _, ok := payload["ExpiryTime"].(float64); !ok {
			return err
		}
		if _, ok := payload["TimeCreated"].(float64); !ok {
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
