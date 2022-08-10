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

package errors

import "github.com/supertokens/supertokens-golang/recipe/session/sessmodels"

const (
	UnauthorizedErrorStr       = "UNAUTHORISED"
	TryRefreshTokenErrorStr    = "TRY_REFRESH_TOKEN"
	TokenTheftDetectedErrorStr = "TOKEN_THEFT_DETECTED"
)

// TryRefreshTokenError used for when the refresh API needs to be called
type TryRefreshTokenError struct {
	Msg string
}

func (err TryRefreshTokenError) Error() string {
	return err.Msg
}

// TokenTheftDetectedError used for when token theft has happened for a session
type TokenTheftDetectedError struct {
	Msg     string
	Payload TokenTheftDetectedErrorPayload
}

type TokenTheftDetectedErrorPayload struct {
	SessionHandle string
	UserID        string
}

func (err TokenTheftDetectedError) Error() string {
	return err.Msg
}

// UnauthorizedError used for when the user has been logged out
type UnauthorizedError struct {
	Msg          string
	ClearCookies *bool
}

func (err UnauthorizedError) Error() string {
	return err.Msg
}

type InvalidClaimError struct {
	Msg           string
	InvalidClaims []sessmodels.ClaimValidationError
}

func (err InvalidClaimError) Error() string {
	return err.Msg
}
