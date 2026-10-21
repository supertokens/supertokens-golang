/* Copyright (c) 2025, VRAI Labs and/or its affiliates. All rights reserved.
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

package api

import (
	"github.com/supertokens/supertokens-golang/recipe/session"
	"github.com/supertokens/supertokens-golang/recipe/session/claims"
	"github.com/supertokens/supertokens-golang/recipe/session/sessmodels"
	"github.com/supertokens/supertokens-golang/recipe/webauthn/webauthnmodels"
	"github.com/supertokens/supertokens-golang/supertokens"
)

func getSessionWithoutClaimValidation(options webauthnmodels.APIOptions, userContext supertokens.UserContext) (sessmodels.SessionContainer, error) {
	// These APIs manage the user's credentials, so we check that the session hasn't been revoked
	// (e.g. by account recovery) instead of only verifying the access token.
	checkDatabase := true
	return session.GetSession(
		options.Req,
		options.Res,
		&sessmodels.VerifySessionOptions{
			CheckDatabase: &checkDatabase,
			OverrideGlobalClaimValidators: func(_ []claims.SessionClaimValidator, _ sessmodels.SessionContainer, _ supertokens.UserContext) ([]claims.SessionClaimValidator, error) {
				return []claims.SessionClaimValidator{}, nil
			},
		},
		userContext,
	)
}
