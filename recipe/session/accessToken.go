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
	"encoding/json"
	"errors"
	"github.com/MicahParks/keyfunc"
	"github.com/golang-jwt/jwt/v4"
	sterrors "github.com/supertokens/supertokens-golang/recipe/session/errors"
	"strings"
)

type accessTokenInfoStruct struct {
	sessionHandle           string
	userID                  string
	refreshTokenHash1       string
	parentRefreshTokenHash1 *string
	userData                map[string]interface{}
	antiCsrfToken           *string
	expiryTime              uint64
	timeCreated             uint64
}

func getInfoFromAccessToken(jwtInfo ParsedJWTInfo, jwks keyfunc.JWKS, doAntiCsrfCheck bool) (*accessTokenInfoStruct, error) {
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

		for _, value := range jwks.ReadOnlyKeys() {
			keys = append(keys, value)
		}

		for _, key := range keys {
			/**
			For each key we need to create a jwks structure to use with the verification library
			{
				keys: [...]
			}
			*/
			keysTemp := [1]interface{}{key}
			jwksTemp := map[string]interface{}{
				"keys": keysTemp,
			}

			jsonString, marshalError := json.Marshal(jwksTemp)
			if marshalError != nil {
				return nil, sterrors.TryRefreshTokenError{
					Msg: "Invalid JWK response",
				}
			}

			jwksToUse, jwksError := keyfunc.NewJSON(jsonString)
			if jwksError != nil {
				return nil, sterrors.TryRefreshTokenError{
					Msg: "Invalid JWT response",
				}
			}

			parsedToken, parseErr := jwt.Parse(jwtInfo.RawTokenString, jwksToUse.Keyfunc)

			if parseErr != nil && (errors.Is(parseErr, jwt.ErrSignatureInvalid) || errors.Is(parseErr, keyfunc.ErrKIDNotFound)) {
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

	err := validateAccessTokenStructure(payload, jwtInfo.Version)
	if err != nil {
		return nil, sterrors.TryRefreshTokenError{
			Msg: err.Error(),
		}
	}

	// We can assume these as defined, since validateAccessTokenStructure checks this
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

	if antiCsrfToken == nil && doAntiCsrfCheck {
		return nil, sterrors.TryRefreshTokenError{
			Msg: "Access token does not contain the anti-csrf token.",
		}
	}

	if expiryTime < getCurrTimeInMS() {
		return nil, sterrors.TryRefreshTokenError{
			Msg: "Access token expired",
		}
	}

	return &accessTokenInfoStruct{
		sessionHandle:           *sessionHandle,
		userID:                  userID,
		refreshTokenHash1:       *refreshTokenHash1,
		parentRefreshTokenHash1: parentRefreshTokenHash1,
		userData:                userData,
		antiCsrfToken:           antiCsrfToken,
		expiryTime:              expiryTime,
		timeCreated:             timeCreated,
	}, nil
}

func validateAccessTokenStructure(payload map[string]interface{}, version int) error {
	err := errors.New("Access token does not contain all the information. Maybe the structure has changed?")

	if version >= 3 {
		if _, ok := payload["sessionHandle"].(string); !ok {
			return err
		}
		if _, ok := payload["sub"].(string); !ok {
			return err
		}
		if _, ok := payload["refreshTokenHash1"].(string); !ok {
			return err
		}
		if _, ok := payload["exp"].(float64); !ok {
			return err
		}
		if _, ok := payload["iat"].(float64); !ok {
			return err
		}
	} else {
		if _, ok := payload["sessionHandle"].(string); !ok {
			return err
		}
		if _, ok := payload["userId"].(string); !ok {
			return err
		}
		if _, ok := payload["refreshTokenHash1"].(string); !ok {
			return err
		}
		if payload["userData"] == nil {
			return err
		}
		if _, ok := payload["userData"].(map[string]interface{}); !ok {
			return err
		}
		if _, ok := payload["expiryTime"].(float64); !ok {
			return err
		}
		if _, ok := payload["timeCreated"].(float64); !ok {
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
