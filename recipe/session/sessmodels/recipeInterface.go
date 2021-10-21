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

import "net/http"

type RecipeInterface struct {
	CreateNewSession            func(res http.ResponseWriter, userID string, accessTokenPayload map[string]interface{}, sessionData map[string]interface{}) (SessionContainer, error)
	GetSession                  func(req *http.Request, res http.ResponseWriter, options *VerifySessionOptions) (*SessionContainer, error)
	RefreshSession              func(req *http.Request, res http.ResponseWriter) (SessionContainer, error)
	GetSessionInformation       func(sessionHandle string) (SessionInformation, error)
	RevokeAllSessionsForUser    func(userID string) ([]string, error)
	GetAllSessionHandlesForUser func(userID string) ([]string, error)
	RevokeSession               func(sessionHandle string) (bool, error)
	RevokeMultipleSessions      func(sessionHandles []string) ([]string, error)
	UpdateSessionData           func(sessionHandle string, newSessionData map[string]interface{}) error
	UpdateAccessTokenPayload    func(sessionHandle string, newAccessTokenPayload map[string]interface{}) error
	GetAccessTokenLifeTimeMS    func() (uint64, error)
	GetRefreshTokenLifeTimeMS   func() (uint64, error)
}
