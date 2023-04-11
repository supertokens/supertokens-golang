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
	"strings"

	"github.com/supertokens/supertokens-golang/recipe/session/errors"
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

func getInfoFromAccessToken(jwtInfo ParsedJWTInfo, jwtSigningPublicKey string, doAntiCsrfCheck bool) (*accessTokenInfoStruct, error) {

	err := verifyJWT(jwtInfo, jwtSigningPublicKey)
	if err != nil {
		return nil, errors.TryRefreshTokenError{
			Msg: err.Error(),
		}
	}

	payload := jwtInfo.Payload
	// This should be called before this function, but the check is very quick, so we can also do them here
	err = validateAccessTokenStructure(payload)
	if err != nil {
		return nil, errors.TryRefreshTokenError{
			Msg: err.Error(),
		}
	}

	// We can assume these as defined, since validateAccessTokenPayload checks this
	sessionHandle := sanitizeStringInput(payload["sessionHandle"])
	userID := sanitizeStringInput(payload["userId"])
	refreshTokenHash1 := sanitizeStringInput(payload["refreshTokenHash1"])
	parentRefreshTokenHash1 := sanitizeStringInput(payload["parentRefreshTokenHash1"])
	userData := payload["userData"].(map[string]interface{})
	antiCsrfToken := sanitizeStringInput(payload["antiCsrfToken"])
	expiryTime := sanitizeNumberInputAsUint64(payload["expiryTime"])
	timeCreated := sanitizeNumberInputAsUint64(payload["timeCreated"])

	if antiCsrfToken == nil && doAntiCsrfCheck {
		return nil, defaultErrors.New("Access token does not contain the anti-csrf token.")
	}

	if *expiryTime < getCurrTimeInMS() {
		return nil, errors.TryRefreshTokenError{
			Msg: "Access token expired",
		}
	}

	return &accessTokenInfoStruct{
		sessionHandle:           *sessionHandle,
		userID:                  *userID,
		refreshTokenHash1:       *refreshTokenHash1,
		parentRefreshTokenHash1: parentRefreshTokenHash1,
		userData:                userData,
		antiCsrfToken:           antiCsrfToken,
		expiryTime:              *expiryTime,
		timeCreated:             *timeCreated,
	}, nil
}

func validateAccessTokenStructure(payload map[string]interface{}) error {
	err := defaultErrors.New("Access token does not contain all the information. Maybe the structure has changed?")

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
