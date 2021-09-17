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

func getInfoFromAccessToken(token string, jwtSigningPublicKey string, doAntiCsrfCheck bool) (*accessTokenInfoStruct, error) {
	payload, err := verifyJWTAndGetPayload(token, jwtSigningPublicKey)
	if err != nil {
		return nil, errors.TryRefreshTokenError{
			Msg: err.Error(),
		}
	}

	sessionHandle := sanitizeStringInput(payload["sessionHandle"])
	userID := sanitizeStringInput(payload["userId"])
	refreshTokenHash1 := sanitizeStringInput(payload["refreshTokenHash1"])
	parentRefreshTokenHash1 := sanitizeStringInput(payload["parentRefreshTokenHash1"])

	var userData *map[string]interface{} = nil
	if payload["userData"] != nil {
		temp := payload["userData"].(map[string]interface{})
		userData = &temp
	}

	antiCsrfToken := sanitizeStringInput(payload["antiCsrfToken"])

	var expiryTime *uint64 = nil
	if payload["expiryTime"] != nil {
		temp := uint64(payload["expiryTime"].(float64))
		expiryTime = &temp
	}

	var timeCreated *uint64 = nil
	if payload["timeCreated"] != nil {
		temp := uint64(payload["timeCreated"].(float64))
		timeCreated = &temp
	}

	if sessionHandle == nil ||
		userID == nil ||
		refreshTokenHash1 == nil ||
		userData == nil ||
		(antiCsrfToken == nil && doAntiCsrfCheck) ||
		expiryTime == nil ||
		timeCreated == nil {
		return nil, errors.TryRefreshTokenError{
			Msg: "Access token does not contain all the information. Maybe the structure has changed?",
		}
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
		userData:                *userData,
		antiCsrfToken:           antiCsrfToken,
		expiryTime:              *expiryTime,
		timeCreated:             *timeCreated,
	}, nil
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
